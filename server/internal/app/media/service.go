package media

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/media"
)

type Service struct {
	repo      media.Repository
	uploadDir string
}

func NewService(repo media.Repository, uploadDir string) *Service {
	trimmed := strings.TrimSpace(uploadDir)
	if trimmed == "" {
		trimmed = filepath.Join("storage", "uploads")
	}
	return &Service{
		repo:      repo,
		uploadDir: trimmed,
	}
}

type UploadResult struct {
	File    media.UploadFile
	Created bool
}

func (s *Service) Upload(ctx context.Context, file *multipart.FileHeader, fileType string) (*UploadResult, error) {
	if file == nil {
		return nil, errors.New("file is required")
	}

	dir, err := dirForType(fileType)
	if err != nil {
		return nil, err
	}

	hash, err := hashFile(file)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByHash(ctx, hash)
	if err != nil && !errors.Is(err, media.ErrUploadFileNotFound) {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := s.buildFilename(dir, ext)
	storedPath := "/" + dir + "/" + filename
	diskPath := s.diskPathFromStored(storedPath)

	if existing != nil {
		existingDisk := s.diskPathFromStored(existing.Path)
		if fileExists(existingDisk) {
			return &UploadResult{File: *existing, Created: false}, nil
		}
		if err := s.saveFile(file, diskPath); err != nil {
			return nil, err
		}
		if existing.Path != storedPath {
			if err := s.repo.UpdatePath(ctx, existing.ID, storedPath); err != nil {
				return nil, err
			}
			existing.Path = storedPath
		}
		return &UploadResult{File: *existing, Created: false}, nil
	}

	if err := s.saveFile(file, diskPath); err != nil {
		return nil, err
	}

	record := &media.UploadFile{
		Name: file.Filename,
		Path: storedPath,
		Type: strings.ToLower(strings.TrimSpace(fileType)),
		Size: file.Size,
		Hash: hash,
	}
	if err := s.repo.Create(ctx, record); err != nil {
		return nil, err
	}
	return &UploadResult{File: *record, Created: true}, nil
}

type ListResult struct {
	Items []media.UploadFile
	Total int64
	Page  int
	Size  int
}

func (s *Service) List(ctx context.Context, page int, size int) (*ListResult, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	items, total, err := s.repo.List(ctx, offset, size)
	if err != nil {
		return nil, err
	}
	return &ListResult{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func (s *Service) Rename(ctx context.Context, id int64, name string) (*media.UploadFile, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, errors.New("name is required")
	}
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if file.Name == trimmed {
		return file, nil
	}
	if err := s.repo.UpdateName(ctx, id, trimmed); err != nil {
		return nil, err
	}
	file.Name = trimmed
	return file, nil
}

func (s *Service) Delete(ctx context.Context, id int64) (*media.UploadFile, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	diskPath := s.diskPathFromStored(file.Path)
	if err := removeFile(diskPath); err != nil {
		return nil, err
	}
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*media.UploadFile, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ResolveDiskPath(storedPath string) (string, error) {
	diskPath := s.diskPathFromStored(storedPath)
	if diskPath == "" {
		return "", errors.New("empty stored path")
	}
	return diskPath, nil
}

func (s *Service) saveFile(file *multipart.FileHeader, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if fileExists(path) {
		return nil
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *Service) diskPathFromStored(storedPath string) string {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return ""
	}
	clean := filepath.Clean(trimmed)
	clean = strings.TrimPrefix(clean, string(filepath.Separator))
	uploadDir := filepath.Clean(s.uploadDir)
	if strings.HasPrefix(clean, uploadDir+string(filepath.Separator)) || clean == uploadDir {
		return clean
	}
	return filepath.Join(uploadDir, clean)
}

func (s *Service) buildFilename(dir string, ext string) string {
	base := time.Now().Format("2006-01-02-15:04:05")
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	for i := 0; i < 5; i++ {
		suffix := randomHex(2)
		filename := base + "-" + suffix + ext
		if !fileExists(filepath.Join(s.uploadDir, dir, filename)) {
			return filename
		}
	}
	suffix := randomHex(4)
	return base + "-" + suffix + ext
}

func randomHex(n int) string {
	if n <= 0 {
		return ""
	}
	byteLen := (n + 1) / 2
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)[:n]
	}
	fallback := hex.EncodeToString([]byte(time.Now().Format("150405.000")))
	if len(fallback) >= n {
		return fallback[:n]
	}
	return fallback
}

func dirForType(fileType string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(fileType)) {
	case "picture":
		return "pictures", nil
	case "file":
		return "files", nil
	default:
		return "", media.ErrInvalidUploadType
	}
}

func hashFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, src); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func removeFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
