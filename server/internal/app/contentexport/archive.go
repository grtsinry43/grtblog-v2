package contentexport

// 归档工具复刻自 internal/app/backup/archive.go（其 helper 均为 unexported，
// 无法跨包引用）。保持相同模式：O_CREATE|O_EXCL 0600、failed-flag 清理、
// SHA-256 校验、checksums.sha256、context 可取消的拷贝。

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type archiveEntry struct {
	archivePath string
	diskPath    string
	sum         string
}

// writeExportArchive 把 workdir 整棵树打成 tar.gz：manifest.json 先行，
// 随后按排序后的路径写入全部文件，末尾写 checksums.sha256。
func writeExportArchive(ctx context.Context, outputPath, workdir string, manifest *Manifest) error {
	entries, err := collectArchiveEntries(workdir)
	if err != nil {
		return err
	}
	if manifest.Checksums == nil {
		manifest.Checksums = make(map[string]string, len(entries))
	}
	for _, entry := range entries {
		manifest.Checksums[entry.archivePath] = entry.sum
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	failed := true
	defer func() {
		_ = file.Close()
		if failed {
			_ = os.Remove(outputPath)
		}
	}()
	gz := gzip.NewWriter(file)
	tw := tar.NewWriter(gz)

	manifestRaw, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := writeTarBytes(tw, "manifest.json", manifestRaw, 0o600); err != nil {
		return err
	}
	for _, entry := range entries {
		if err := writeTarFile(ctx, tw, entry.archivePath, entry.diskPath); err != nil {
			return err
		}
	}
	var checksumText strings.Builder
	for _, entry := range entries {
		fmt.Fprintf(&checksumText, "%s  %s\n", entry.sum, entry.archivePath)
	}
	if err := writeTarBytes(tw, "checksums.sha256", []byte(checksumText.String()), 0o600); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	failed = false
	return nil
}

// collectArchiveEntries 排序遍历 workdir，跳过符号链接与非常规文件。
func collectArchiveEntries(workdir string) ([]archiveEntry, error) {
	var entries []archiveEntry
	err := filepath.WalkDir(workdir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(workdir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("导出目录包含不支持的符号链接: %s", rel)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("导出目录包含不支持的文件: %s", rel)
		}
		sum, err := hashPath(path)
		if err != nil {
			return err
		}
		entries = append(entries, archiveEntry{
			archivePath: filepath.ToSlash(rel),
			diskPath:    path,
			sum:         sum,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].archivePath < entries[j].archivePath })
	return entries, nil
}

func writeTarBytes(tw *tar.Writer, name string, raw []byte, mode int64) error {
	header := &tar.Header{Name: name, Mode: mode, Size: int64(len(raw))}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err := tw.Write(raw)
	return err
}

func writeTarFile(ctx context.Context, tw *tar.Writer, archivePath, diskPath string) error {
	info, err := os.Stat(diskPath)
	if err != nil {
		return err
	}
	header := &tar.Header{Name: archivePath, Mode: 0o600, Size: info.Size()}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	file, err := os.Open(diskPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(tw, &contextReader{ctx: ctx, reader: file})
	return err
}

func hashPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
