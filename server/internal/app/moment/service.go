package moment

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

// CreateMoment 创建手记
func (s *Service) CreateMoment(ctx context.Context, authorID int64, cmd CreateMomentCmd) (*content.Moment, error) {
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

	if cmd.ColumnID != nil {
		if _, err := s.repo.GetColumnByID(ctx, *cmd.ColumnID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureTagsExist(ctx, cmd.TopicIDs); err != nil {
		return nil, err
	}

	createdAt := time.Now()
	if cmd.CreatedAt != nil {
		createdAt = *cmd.CreatedAt
	}

	toc := contentutil.GenerateTOC(cmd.Content)
	summary := contentutil.BuildSummary(cmd.Summary, cmd.Content)

	moment := &content.Moment{
		Title:       cmd.Title,
		Summary:     summary,
		AISummary:   nil,
		TOC:         toc,
		Content:     cmd.Content,
		ContentHash: content.MomentContentHash(cmd.Title, summary, cmd.Content),
		AuthorID:    authorID,
		Image:       cmd.Image,
		ColumnID:    cmd.ColumnID,
		ShortURL:    shortURL,
		IsPublished: cmd.IsPublished,
		IsTop:       cmd.IsTop,
		IsHot:       cmd.IsHot,
		IsOriginal:  cmd.IsOriginal,
		CreatedAt:   createdAt,
	}

	if err := s.repo.CreateMoment(ctx, moment); err != nil {
		return nil, err
	}

	if len(cmd.TopicIDs) > 0 {
		if err := s.repo.SyncTopicsToMoment(ctx, moment.ID, cmd.TopicIDs); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	_ = s.events.Publish(ctx, MomentCreated{
		ID:        moment.ID,
		AuthorID:  moment.AuthorID,
		Title:     moment.Title,
		ShortURL:  moment.ShortURL,
		Published: moment.IsPublished,
		At:        now,
	})
	if moment.IsPublished {
		_ = s.events.Publish(ctx, MomentPublished{
			ID:       moment.ID,
			AuthorID: moment.AuthorID,
			Title:    moment.Title,
			ShortURL: moment.ShortURL,
			At:       now,
		})
	}

	return moment, nil
}

// UpdateMoment 更新手记
func (s *Service) UpdateMoment(ctx context.Context, cmd UpdateMomentCmd) (*content.Moment, error) {
	existing, err := s.repo.GetMomentByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	prevPublished := existing.IsPublished

	if cmd.ColumnID != nil {
		if _, err := s.repo.GetColumnByID(ctx, *cmd.ColumnID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureTagsExist(ctx, cmd.TopicIDs); err != nil {
		return nil, err
	}

	toc := contentutil.GenerateTOC(cmd.Content)
	summary := contentutil.BuildSummary(cmd.Summary, cmd.Content)

	existing.Title = cmd.Title
	existing.Summary = summary
	existing.TOC = toc
	existing.Content = cmd.Content
	existing.ContentHash = content.MomentContentHash(cmd.Title, summary, cmd.Content)
	existing.Image = cmd.Image
	existing.ColumnID = cmd.ColumnID
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
	existing.IsPublished = cmd.IsPublished
	existing.IsTop = cmd.IsTop
	existing.IsHot = cmd.IsHot
	existing.IsOriginal = cmd.IsOriginal

	if err := s.repo.UpdateMoment(ctx, existing); err != nil {
		return nil, err
	}

	if err := s.repo.SyncTopicsToMoment(ctx, existing.ID, cmd.TopicIDs); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, MomentUpdated{
		ID:          existing.ID,
		AuthorID:    existing.AuthorID,
		Title:       existing.Title,
		ShortURL:    existing.ShortURL,
		Published:   existing.IsPublished,
		ContentHash: existing.ContentHash,
		Summary:     existing.Summary,
		TOC:         existing.TOC,
		Content:     existing.Content,
		At:          now,
	})
	if !prevPublished && existing.IsPublished {
		_ = s.events.Publish(ctx, MomentPublished{
			ID:       existing.ID,
			AuthorID: existing.AuthorID,
			Title:    existing.Title,
			ShortURL: existing.ShortURL,
			At:       now,
		})
	}
	if prevPublished && !existing.IsPublished {
		_ = s.events.Publish(ctx, MomentUnpublished{
			ID:       existing.ID,
			AuthorID: existing.AuthorID,
			Title:    existing.Title,
			ShortURL: existing.ShortURL,
			At:       now,
		})
	}

	return existing, nil
}

// GetMomentByID 根据 ID 获取手记
func (s *Service) GetMomentByID(ctx context.Context, id int64) (*content.Moment, error) {
	return s.repo.GetMomentByID(ctx, id)
}

// GetMomentByShortURL 根据短链接获取手记
func (s *Service) GetMomentByShortURL(ctx context.Context, shortURL string) (*content.Moment, error) {
	moment, err := s.repo.GetMomentByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	_ = s.repo.UpdateMomentViews(ctx, moment.ID)

	return moment, nil
}

// ListMoments 获取手记列表
func (s *Service) ListMoments(ctx context.Context, options content.MomentListOptionsInternal) ([]*content.Moment, int64, error) {
	return s.repo.ListMoments(ctx, options)
}

// DeleteMoment 删除手记
func (s *Service) DeleteMoment(ctx context.Context, id int64) error {
	momentItem, err := s.repo.GetMomentByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteMoment(ctx, id); err != nil {
		return err
	}
	_ = s.events.Publish(ctx, MomentDeleted{
		ID:       momentItem.ID,
		AuthorID: momentItem.AuthorID,
		Title:    momentItem.Title,
		ShortURL: momentItem.ShortURL,
		At:       time.Now(),
	})
	return nil
}

// GetMomentTopics 获取手记话题。
func (s *Service) GetMomentTopics(ctx context.Context, momentID int64) ([]*content.Tag, error) {
	return s.repo.GetTopicsByMomentID(ctx, momentID)
}

// GetMomentMetrics 获取手记指标
func (s *Service) GetMomentMetrics(ctx context.Context, momentID int64) (*content.MomentMetrics, error) {
	return s.repo.GetMomentMetrics(ctx, momentID)
}

func (s *Service) ensureShortURLAvailable(ctx context.Context, shortURL string) (string, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		for i := 0; i < 5; i++ {
			candidate := contentutil.GenerateRandomShortURL()
			_, err := s.repo.GetMomentByShortURL(ctx, candidate)
			if err != nil {
				if errors.Is(err, content.ErrMomentNotFound) {
					return candidate, nil
				}
				return "", err
			}
		}
		return "", content.ErrMomentShortURLExists
	}

	existing, err := s.repo.GetMomentByShortURL(ctx, shortURL)
	if err != nil && !errors.Is(err, content.ErrMomentNotFound) {
		return "", err
	}
	if err == nil && existing != nil {
		return "", content.ErrMomentShortURLExists
	}
	return shortURL, nil
}

func (s *Service) ensureTagsExist(ctx context.Context, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	unique := make(map[int64]struct{}, len(tagIDs))
	for _, id := range tagIDs {
		if id <= 0 {
			return content.ErrTagNotFound
		}
		unique[id] = struct{}{}
	}
	ids := make([]int64, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	ok, err := s.repo.TagIDsExist(ctx, ids)
	if err != nil {
		return err
	}
	if !ok {
		return content.ErrTagNotFound
	}
	return nil
}
