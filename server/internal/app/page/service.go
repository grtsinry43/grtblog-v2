package page

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/contentutil"
	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type Service struct {
	repo   content.Repository
	events appEvent.Bus
}

func NewService(repo content.Repository, events appEvent.Bus) *Service {
	if events == nil {
		events = appEvent.NopBus{}
	}
	return &Service{repo: repo, events: events}
}

// CreatePage 创建页面
func (s *Service) CreatePage(ctx context.Context, cmd CreatePageCmd) (*content.Page, error) {
	shortURL := ""
	if cmd.ShortURL != nil {
		shortURL = strings.TrimSpace(*cmd.ShortURL)
	}
	if shortURL == "" {
		shortURL = contentutil.GenerateShortURLFromTitle(cmd.Title)
	}
	shortURL, err := s.ensureShortURLAvailable(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	if cmd.CreatedAt != nil {
		createdAt = *cmd.CreatedAt
	}

	description := trimPtr(cmd.Description)
	toc := contentutil.GenerateTOC(cmd.Content)

	page := &content.Page{
		Title:       cmd.Title,
		Description: description,
		AISummary:   nil,
		TOC:         toc,
		Content:     cmd.Content,
		ContentHash: content.PageContentHash(cmd.Title, description, cmd.Content),
		ShortURL:    shortURL,
		IsEnabled:   cmd.IsEnabled,
		IsBuiltin:   cmd.IsBuiltin,
		CreatedAt:   createdAt,
	}

	if err := s.repo.CreatePage(ctx, page); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, PageCreated{
		ID:       page.ID,
		Title:    page.Title,
		ShortURL: page.ShortURL,
		Enabled:  page.IsEnabled,
		At:       now,
	})

	return page, nil
}

// UpdatePage 更新页面
func (s *Service) UpdatePage(ctx context.Context, cmd UpdatePageCmd) (*content.Page, error) {
	existing, err := s.repo.GetPageByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	description := trimPtr(cmd.Description)
	toc := contentutil.GenerateTOC(cmd.Content)

	existing.Title = cmd.Title
	existing.Description = description
	existing.TOC = toc
	existing.Content = cmd.Content
	existing.ContentHash = content.PageContentHash(cmd.Title, description, cmd.Content)
	shortURL := strings.TrimSpace(cmd.ShortURL)
	if shortURL == "" {
		shortURL = existing.ShortURL
	}
	if shortURL != existing.ShortURL {
		shortURL, err = s.ensureShortURLAvailable(ctx, shortURL)
		if err != nil {
			return nil, err
		}
	}
	existing.ShortURL = shortURL
	existing.IsEnabled = cmd.IsEnabled
	existing.IsBuiltin = cmd.IsBuiltin

	if err := s.repo.UpdatePage(ctx, existing); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, PageUpdated{
		ID:          existing.ID,
		Title:       existing.Title,
		ShortURL:    existing.ShortURL,
		Enabled:     existing.IsEnabled,
		ContentHash: existing.ContentHash,
		Description: existing.Description,
		TOC:         existing.TOC,
		Content:     existing.Content,
		At:          now,
	})

	return existing, nil
}

// GetPageByID 根据 ID 获取页面
func (s *Service) GetPageByID(ctx context.Context, id int64) (*content.Page, error) {
	return s.repo.GetPageByID(ctx, id)
}

// GetPageByShortURL 根据短链接获取页面
func (s *Service) GetPageByShortURL(ctx context.Context, shortURL string) (*content.Page, error) {
	page, err := s.repo.GetPageByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	_ = s.repo.UpdatePageViews(ctx, page.ID)

	return page, nil
}

// ListPages 获取页面列表
func (s *Service) ListPages(ctx context.Context, options content.PageListOptionsInternal) ([]*content.Page, int64, error) {
	return s.repo.ListPages(ctx, options)
}

// DeletePage 删除页面
func (s *Service) DeletePage(ctx context.Context, id int64) error {
	page, err := s.repo.GetPageByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.DeletePage(ctx, id); err != nil {
		return err
	}
	_ = s.events.Publish(ctx, PageDeleted{
		ID:       page.ID,
		Title:    page.Title,
		ShortURL: page.ShortURL,
		At:       time.Now(),
	})
	return nil
}

// GetPageMetrics 获取页面指标
func (s *Service) GetPageMetrics(ctx context.Context, pageID int64) (*content.PageMetrics, error) {
	return s.repo.GetPageMetrics(ctx, pageID)
}

func (s *Service) ensureShortURLAvailable(ctx context.Context, shortURL string) (string, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		for i := 0; i < 5; i++ {
			candidate := contentutil.GenerateRandomShortURL()
			_, err := s.repo.GetPageByShortURL(ctx, candidate)
			if err != nil {
				if errors.Is(err, content.ErrPageNotFound) {
					return candidate, nil
				}
				return "", err
			}
		}
		return "", content.ErrPageShortURLExists
	}

	existing, err := s.repo.GetPageByShortURL(ctx, shortURL)
	if err != nil && !errors.Is(err, content.ErrPageNotFound) {
		return "", err
	}
	if err == nil && existing != nil {
		return "", content.ErrPageShortURLExists
	}
	return shortURL, nil
}

func trimPtr(val *string) *string {
	if val == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*val)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
