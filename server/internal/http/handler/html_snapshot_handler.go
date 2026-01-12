package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

const (
	defaultHTMLSnapshotBaseURL = "http://localhost:3000"
	htmlSnapshotPageSize       = 100
)

type HTMLSnapshotHandler struct {
	contentRepo content.Repository
	baseURL     string
	client      *http.Client
}

func NewHTMLSnapshotHandler(contentRepo content.Repository) *HTMLSnapshotHandler {
	return &HTMLSnapshotHandler{
		contentRepo: contentRepo,
		baseURL:     defaultHTMLSnapshotBaseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// RefreshPostsHTML godoc
// @Summary 刷新文章 HTML 缓存
// @Tags Public
// @Produce json
// @Success 200 {object} any
// @Router /public/html/posts/refresh [post]
func (h *HTMLSnapshotHandler) RefreshPostsHTML(c *fiber.Ctx) error {
	go func() {
		if err := h.generatePostsHTML(context.Background()); err != nil {
			log.Printf("[html-snapshot] generate posts html failed: %v", err)
		}
	}()

	return response.SuccessWithMessage[any](c, nil, "ok")
}

func (h *HTMLSnapshotHandler) generatePostsHTML(ctx context.Context) error {
	start := time.Now()
	outputDir := filepath.Join("storage", "html", "posts")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	page := 1
	successCount := 0
	for {
		articles, total, err := h.contentRepo.ListPublicArticles(ctx, content.ArticleListOptions{
			Page:     page,
			PageSize: htmlSnapshotPageSize,
		})
		if err != nil {
			return err
		}

		for _, article := range articles {
			shortURL := strings.TrimSpace(article.ShortURL)
			if shortURL == "" {
				continue
			}
			if strings.Contains(shortURL, "/") || strings.Contains(shortURL, "\\") {
				log.Printf("[html-snapshot] skip invalid short url: %s", shortURL)
				continue
			}

			escaped := url.PathEscape(shortURL)
			pageURL := fmt.Sprintf("%s/posts/%s", h.baseURL, escaped)
			filePath := filepath.Join(outputDir, shortURL+".html")
			if err := h.fetchAndSave(ctx, pageURL, filePath); err != nil {
				log.Printf("[html-snapshot] fetch %s failed: %v", pageURL, err)
			} else {
				successCount++
			}
		}

		if len(articles) == 0 || int64(page*htmlSnapshotPageSize) >= total {
			break
		}
		page++
	}

	indexURL := fmt.Sprintf("%s/posts/", h.baseURL)
	indexPath := filepath.Join(outputDir, "index.html")
	if err := h.fetchAndSave(ctx, indexURL, indexPath); err != nil {
		return fmt.Errorf("fetch index: %w", err)
	}
	successCount++
	log.Printf("[html-snapshot] done success=%d duration=%s", successCount, time.Since(start))

	return nil
}

func (h *HTMLSnapshotHandler) fetchAndSave(ctx context.Context, pageURL, filePath string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, pageURL, nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0o644)
}
