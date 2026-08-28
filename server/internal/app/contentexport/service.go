// Package contentexport 实现内置内容导出：通过内容服务层一次性加载
// 全量内容（admin 全集），生成 Markdown + 完整 meta + 打包图片的自包含
// 归档包。架构镜像 app/backup（异步任务 + ticket 下载）。
package contentexport

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/grtsinry43/grtblog-v2/server/internal/buildinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	exportdomain "github.com/grtsinry43/grtblog-v2/server/internal/domain/contentexport"
)

type Service struct {
	cfg       config.ExportConfig
	uploadDir string
	repo      exportdomain.Repository
	collector *Collector
	mapper    *Mapper
	rootCtx   context.Context
	mu        sync.Mutex
}

func NewService(cfg config.ExportConfig, uploadDir string, repo exportdomain.Repository, collector *Collector, mapper *Mapper) *Service {
	return &Service{
		cfg:       cfg,
		uploadDir: uploadDir,
		repo:      repo,
		collector: collector,
		mapper:    mapper,
		rootCtx:   context.Background(),
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	if err := os.MkdirAll(s.cfg.RootDir, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(s.cfg.RootDir, ".work"), 0o700); err != nil {
		return err
	}
	if err := s.repo.MarkInterrupted(ctx); err != nil {
		return err
	}
	return s.repo.DeleteExpiredTickets(ctx)
}

// Create 创建并异步执行一次导出任务。mode 为空时默认 both。
func (s *Service) Create(ctx context.Context, mode string) (*exportdomain.Record, error) {
	trimmed := strings.TrimSpace(mode)
	if trimmed == "" {
		trimmed = string(exportdomain.ModeBoth)
	}
	normalized, ok := exportdomain.ValidMode(trimmed)
	if !ok {
		return nil, errors.New("mode 只支持 structured、flatten 或 both")
	}
	if !s.mu.TryLock() {
		return nil, exportdomain.ErrExportRunning
	}
	now := time.Now().UTC()
	id := uuid.NewString()
	item := &exportdomain.Record{
		ID:          id,
		Filename:    fmt.Sprintf("grtblog-export-%s-%s.tar.gz", now.Format("20060102T150405Z"), id[:8]),
		Status:      exportdomain.StatusQueued,
		Stage:       "queued",
		TriggerType: "manual",
		Mode:        string(normalized),
		CreatedAt:   now,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	go func() {
		defer s.mu.Unlock()
		timeout := s.cfg.JobTimeout
		if timeout <= 0 {
			timeout = 30 * time.Minute
		}
		jobCtx, cancel := context.WithTimeout(s.rootCtx, timeout)
		defer cancel()
		s.run(jobCtx, item)
	}()
	return item, nil
}

func (s *Service) run(ctx context.Context, item *exportdomain.Record) {
	now := time.Now().UTC()
	item.Status, item.Stage, item.StartedAt = exportdomain.StatusRunning, "collecting_content", &now
	if err := s.repo.Update(context.Background(), item); err != nil {
		return
	}
	fail := func(err error) {
		completed := time.Now().UTC()
		item.Status, item.Stage, item.CompletedAt = exportdomain.StatusFailed, "failed", &completed
		item.ErrorMessage = err.Error()
		_ = s.repo.Update(context.Background(), item)
	}

	workdir := filepath.Join(s.cfg.RootDir, ".work", item.ID)
	if err := os.MkdirAll(workdir, 0o700); err != nil {
		fail(err)
		return
	}
	defer os.RemoveAll(workdir)

	exportedAt := time.Now().UTC()
	snap, err := s.collector.Collect(ctx)
	if err != nil {
		fail(fmt.Errorf("采集内容失败: %w", err))
		return
	}
	item.SiteName, item.SiteURL = snap.SiteName, snap.SiteURL
	item.AppVersion = buildinfo.Version()
	item.ArticleCount = int64(len(snap.Articles))
	item.MomentsCount = int64(len(snap.Moments))
	item.PagesCount = int64(len(snap.Pages))
	item.ThinkingsCount = int64(len(snap.Thinkings))

	item.Stage = "resolving_images"
	_ = s.repo.Update(context.Background(), item)
	resolver := NewResolver(s.cfg, s.uploadDir, workdir, snap.PublicHost)
	total := len(snap.Articles) + len(snap.Moments) + len(snap.Pages) + len(snap.Thinkings)
	items := make([]*preparedItem, 0, total)

	usedPostSlugs := make(map[string]bool)
	usedPageSlugs := make(map[string]bool)

	for _, art := range snap.Articles {
		resp, mapErr := s.mapper.ArticleResp(ctx, art)
		if mapErr != nil {
			fail(fmt.Errorf("映射文章 %d 失败: %w", art.ID, mapErr))
			return
		}
		slug := uniqueSlug(slugFor(art.ShortURL, "article", art.ID), usedPostSlugs, art.ID)
		it := &preparedItem{
			kind:        "article",
			id:          art.ID,
			routePath:   "/posts/" + slug,
			sourcePath:  fmt.Sprintf("/admin/articles/%d", art.ID),
			structDir:   []string{"posts", slug},
			flattenName: flattenFilename("article", resp.Title, art.ID),
			content:     resp.Content,
			meta:        resp,
			extInfo:     resp.ExtInfo,
		}
		referrer := fmt.Sprintf("article:%d", art.ID)
		it.spans = extractImageRefs(resp.Content)
		for _, sp := range it.spans {
			resolver.See(ctx, sp.ref, referrer)
		}
		if resp.Cover != nil && *resp.Cover != "" {
			resolver.See(ctx, *resp.Cover, referrer)
		}
		for _, ref := range extInfoRefs(resp.ExtInfo) {
			resolver.See(ctx, ref, referrer)
		}
		items = append(items, it)
	}

	for _, mo := range snap.Moments {
		resp, mapErr := s.mapper.MomentResp(ctx, snap.SiteTZ, mo)
		if mapErr != nil {
			fail(fmt.Errorf("映射手记 %d 失败: %w", mo.ID, mapErr))
			return
		}
		year, month, day := dateParts(mo.CreatedAt, snap.SiteTZ)
		slug := slugFor(mo.ShortURL, "moment", mo.ID)
		it := &preparedItem{
			kind:        "moment",
			id:          mo.ID,
			routePath:   fmt.Sprintf("/moments/%s/%s/%s/%s", year, month, day, slug),
			sourcePath:  fmt.Sprintf("/admin/moments/%d", mo.ID),
			structDir:   []string{"moments", year, month, day, slug},
			flattenName: flattenFilename("moment", resp.Title, mo.ID),
			content:     resp.Content,
			meta:        resp,
			extInfo:     resp.ExtInfo,
		}
		referrer := fmt.Sprintf("moment:%d", mo.ID)
		it.spans = extractImageRefs(resp.Content)
		for _, sp := range it.spans {
			resolver.See(ctx, sp.ref, referrer)
		}
		for _, img := range resp.Image {
			resolver.See(ctx, img, referrer)
		}
		for _, ref := range extInfoRefs(resp.ExtInfo) {
			resolver.See(ctx, ref, referrer)
		}
		items = append(items, it)
	}

	for _, t := range snap.Thinkings {
		resp, mapErr := s.mapper.ThinkingResp(ctx, t)
		if mapErr != nil {
			fail(fmt.Errorf("映射思考 %d 失败: %w", t.ID, mapErr))
			return
		}
		year, month, day := dateParts(t.CreatedAt, snap.SiteTZ)
		anchor := fmt.Sprintf("thinking-%d", t.ID)
		it := &preparedItem{
			kind:        "thinking",
			id:          t.ID,
			routePath:   "/thinkings#" + anchor,
			sourcePath:  fmt.Sprintf("/thinkings/%d", t.ID),
			structDir:   []string{"thinkings", year, month, day, anchor},
			flattenName: flattenFilename("thinking", firstLine(resp.Content), t.ID),
			content:     resp.Content,
			meta:        resp,
		}
		referrer := fmt.Sprintf("thinking:%d", t.ID)
		it.spans = extractImageRefs(resp.Content)
		for _, sp := range it.spans {
			resolver.See(ctx, sp.ref, referrer)
		}
		items = append(items, it)
	}

	for _, pg := range snap.Pages {
		resp, mapErr := s.mapper.PageResp(ctx, pg)
		if mapErr != nil {
			fail(fmt.Errorf("映射页面 %d 失败: %w", pg.ID, mapErr))
			return
		}
		slug := uniqueSlug(slugFor(pg.ShortURL, "page", pg.ID), usedPageSlugs, pg.ID)
		it := &preparedItem{
			kind:        "page",
			id:          pg.ID,
			routePath:   "/" + slug,
			sourcePath:  fmt.Sprintf("/admin/pages/%d", pg.ID),
			structDir:   []string{"pages", slug},
			flattenName: flattenFilename("page", resp.Title, pg.ID),
			content:     resp.Content,
			meta:        resp,
			extInfo:     resp.ExtInfo,
		}
		referrer := fmt.Sprintf("page:%d", pg.ID)
		it.spans = extractImageRefs(resp.Content)
		for _, sp := range it.spans {
			resolver.See(ctx, sp.ref, referrer)
		}
		for _, ref := range extInfoRefs(resp.ExtInfo) {
			resolver.See(ctx, ref, referrer)
		}
		items = append(items, it)
	}

	item.Stage = "downloading_external"
	_ = s.repo.Update(context.Background(), item)
	resolver.DownloadAll(ctx)

	mode := exportdomain.Mode(item.Mode)
	if mode == exportdomain.ModeStructured || mode == exportdomain.ModeBoth {
		if err := WriteStructured(items, workdir, resolver, exportedAt); err != nil {
			fail(fmt.Errorf("写入 structured 布局失败: %w", err))
			return
		}
	}
	if mode == exportdomain.ModeFlatten || mode == exportdomain.ModeBoth {
		if err := WriteFlatten(items, workdir, resolver, exportedAt); err != nil {
			fail(fmt.Errorf("写入 flatten 布局失败: %w", err))
			return
		}
	}
	if err := WriteTaxonomy(snap, workdir); err != nil {
		fail(fmt.Errorf("写入 taxonomy 失败: %w", err))
		return
	}

	selfCount, externalCount, failedDownloads := resolver.Stats()
	if failedDownloads == nil {
		failedDownloads = []FailedDownload{}
	}
	item.ImageCount = selfCount + externalCount
	item.FailedImageCount = int64(len(failedDownloads))

	manifest := &Manifest{
		FormatVersion: ArchiveFormatVersion,
		ExportID:      item.ID,
		Mode:          item.Mode,
		CreatedAt:     exportedAt,
		AppVersion:    item.AppVersion,
		SiteName:      item.SiteName,
		SiteURL:       item.SiteURL,
		Counts: ManifestCounts{
			Articles:   item.ArticleCount,
			Moments:    item.MomentsCount,
			Pages:      item.PagesCount,
			Thinkings:  item.ThinkingsCount,
			Categories: int64(len(snap.Categories)),
			Columns:    int64(len(snap.Columns)),
			Tags:       int64(len(snap.Tags)),
			Total:      int64(total),
		},
		Images: ManifestImages{
			SelfHosted:      selfCount,
			External:        externalCount,
			FailedCount:     int64(len(failedDownloads)),
			FailedDownloads: failedDownloads,
		},
		ContainsSensitive: true,
	}

	item.Stage = "packing_archive"
	_ = s.repo.Update(context.Background(), item)
	tempArchive := filepath.Join(s.cfg.RootDir, ".work", item.ID+".tar.gz")
	if err := writeExportArchive(ctx, tempArchive, workdir, manifest); err != nil {
		fail(fmt.Errorf("打包导出归档失败: %w", err))
		return
	}
	finalPath := s.archivePath(item.Filename)
	if err := os.Rename(tempArchive, finalPath); err != nil {
		fail(fmt.Errorf("发布导出归档失败: %w", err))
		return
	}
	stat, err := os.Stat(finalPath)
	if err != nil {
		fail(err)
		return
	}
	archiveSum, err := hashPath(finalPath)
	if err != nil {
		fail(err)
		return
	}
	completed := time.Now().UTC()
	item.Status, item.Stage, item.CompletedAt = exportdomain.StatusCompleted, "completed", &completed
	item.SizeBytes, item.SHA256, item.ErrorMessage = stat.Size(), archiveSum, ""
	_ = s.repo.Update(context.Background(), item)
}

func (s *Service) List(ctx context.Context) ([]exportdomain.Record, error) { return s.repo.List(ctx) }

func (s *Service) Get(ctx context.Context, id string) (*exportdomain.Record, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if item.Status == exportdomain.StatusRunning || item.Status == exportdomain.StatusQueued {
		return exportdomain.ErrExportRunning
	}
	if err := os.Remove(s.archivePath(item.Filename)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) IssueDownloadTicket(ctx context.Context, id string) (string, time.Time, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return "", time.Time{}, err
	}
	if item.Status != exportdomain.StatusCompleted {
		return "", time.Time{}, errors.New("导出任务尚未完成")
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	expires := time.Now().UTC().Add(s.cfg.TicketTTL)
	ticket := exportdomain.DownloadTicket{TokenHash: tokenHash(token), ExportID: id, ExpiresAt: expires, CreatedAt: time.Now().UTC()}
	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return "", time.Time{}, err
	}
	return token, expires, nil
}

func (s *Service) ResolveDownload(ctx context.Context, token string) (*exportdomain.Record, string, error) {
	item, err := s.repo.ResolveTicket(ctx, tokenHash(token))
	if err != nil {
		return nil, "", err
	}
	path := s.archivePath(item.Filename)
	if _, err := os.Stat(path); err != nil {
		return nil, "", err
	}
	return item, path, nil
}

func (s *Service) archivePath(filename string) string {
	return filepath.Join(s.cfg.RootDir, filepath.Base(filename))
}

func uniqueSlug(base string, used map[string]bool, id int64) string {
	slug := base
	if used[slug] {
		slug = fmt.Sprintf("%s-%d", base, id)
	}
	used[slug] = true
	return slug
}

func tokenHash(token string) string {
	h := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(h[:])
}
