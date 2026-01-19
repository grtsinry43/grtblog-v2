package htmlsnapshot

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

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

const (
	defaultBaseURL = "http://localhost:3000"
	pageSize       = 100
)

type Service struct {
	contentRepo content.Repository
	baseURL     string
	client      *http.Client
}

func NewService(contentRepo content.Repository, baseURL string) *Service {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Service{
		contentRepo: contentRepo,
		baseURL:     baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *Service) RefreshPostsHTML(ctx context.Context) error {
	start := time.Now()
	outputDir := filepath.Join("storage", "html", "posts")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	page := 1
	successCount := 0
	for {
		articles, total, err := s.contentRepo.ListPublicArticles(ctx, content.ArticleListOptions{
			Page:     page,
			PageSize: pageSize,
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
				continue
			}

			escaped := url.PathEscape(shortURL)
			pageURL := fmt.Sprintf("%s/posts/%s", s.baseURL, escaped)
			filePath := filepath.Join(outputDir, shortURL+".html")
			if err := s.fetchAndSave(ctx, pageURL, filePath); err == nil {
				successCount++
			}
		}

		if len(articles) == 0 || int64(page*pageSize) >= total {
			break
		}
		page++
	}

	indexURL := fmt.Sprintf("%s/posts/", s.baseURL)
	indexPath := filepath.Join(outputDir, "index.html")
	if err := s.fetchAndSave(ctx, indexURL, indexPath); err != nil {
		return fmt.Errorf("fetch index: %w", err)
	}
	successCount++
	log.Printf("[html-snapshot] done success=%d duration=%s", successCount, time.Since(start))

	return nil
}

func (s *Service) fetchAndSave(ctx context.Context, pageURL, filePath string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, pageURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
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
