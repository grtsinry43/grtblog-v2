package contentexport

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteExportArchive(t *testing.T) {
	workdir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workdir, "site-flatten"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(workdir, "uploads", "pictures"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workdir, "taxonomy.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workdir, "site-flatten", "a.md"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workdir, "uploads", "pictures", "x.png"), []byte("png"), 0o600); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), "export.tar.gz")
	manifest := &Manifest{FormatVersion: ArchiveFormatVersion, ExportID: "id", Mode: "flatten", ContainsSensitive: true}
	if err := writeExportArchive(context.Background(), out, workdir, manifest); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)

	names := make(map[string]bool)
	var checksums string
	contentHashes := make(map[string]string)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		names[hdr.Name] = true
		data, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if hdr.Name == "checksums.sha256" {
			checksums = string(data)
		}
		sum := sha256.Sum256(data)
		contentHashes[hdr.Name] = hex.EncodeToString(sum[:])
	}

	for _, want := range []string{"manifest.json", "taxonomy.json", "site-flatten/a.md", "uploads/pictures/x.png", "checksums.sha256"} {
		if !names[want] {
			t.Errorf("归档缺少 %q，实际: %v", want, names)
		}
	}
	// checksums.sha256 与 manifest.Checksums 覆盖 workdir 内全部文件且校验和正确
	for _, file := range []string{"taxonomy.json", "site-flatten/a.md", "uploads/pictures/x.png"} {
		if !strings.Contains(checksums, contentHashes[file]+"  "+file) {
			t.Errorf("checksums.sha256 缺少 %s 的正确条目:\n%s", file, checksums)
		}
		if manifest.Checksums[file] != contentHashes[file] {
			t.Errorf("manifest.Checksums[%s] = %q, want %q", file, manifest.Checksums[file], contentHashes[file])
		}
	}
}

func TestWriteExportArchiveRefusesOverwrite(t *testing.T) {
	workdir := t.TempDir()
	outDir := t.TempDir()
	out := filepath.Join(outDir, "export.tar.gz")
	if err := os.WriteFile(out, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := &Manifest{FormatVersion: ArchiveFormatVersion}
	if err := writeExportArchive(context.Background(), out, workdir, manifest); err == nil {
		t.Fatal("已存在的输出文件应拒绝覆盖 (O_EXCL)")
	}
	// 失败不应残留半成品覆盖原文件
	data, err := os.ReadFile(out)
	if err != nil || string(data) != "existing" {
		t.Fatalf("原文件应保持不变: %v %q", err, data)
	}
}
