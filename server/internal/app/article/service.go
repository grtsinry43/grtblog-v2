package article

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

// CreateArticle 创建文章
func (s *Service) CreateArticle(ctx context.Context, authorID int64, cmd CreateArticleCmd) (*content.Article, error) {
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

	if cmd.CategoryID != nil {
		if _, err := s.repo.GetCategoryByID(ctx, *cmd.CategoryID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureTagsExist(ctx, cmd.TagIDs); err != nil {
		return nil, err
	}

	// 设置创建时间
	createdAt := time.Now()
	if cmd.CreatedAt != nil {
		createdAt = *cmd.CreatedAt
	}

	toc := contentutil.GenerateTOC(cmd.Content)
	summary := contentutil.BuildSummary(cmd.Summary, cmd.Content)

	article := &content.Article{
		Title:       cmd.Title,
		Summary:     summary,
		AISummary:   nil, // 后续可以添加 AI 摘要功能
		LeadIn:      cmd.LeadIn,
		TOC:         toc,
		Content:     cmd.Content,
		ContentHash: content.ArticleContentHash(cmd.Title, cmd.LeadIn, cmd.Content),
		AuthorID:    authorID,
		Cover:       cmd.Cover,
		CategoryID:  cmd.CategoryID,
		ShortURL:    shortURL,
		IsPublished: cmd.IsPublished,
		IsTop:       cmd.IsTop,
		IsHot:       cmd.IsHot,
		IsOriginal:  cmd.IsOriginal,
		CreatedAt:   createdAt,
	}

	if err := s.repo.CreateArticle(ctx, article); err != nil {
		return nil, err
	}

	// 如果有标签，则关联标签
	if len(cmd.TagIDs) > 0 {
		if err := s.repo.SyncTagsToArticle(ctx, article.ID, cmd.TagIDs); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	_ = s.events.Publish(ctx, ArticleCreated{
		ID:        article.ID,
		AuthorID:  article.AuthorID,
		Title:     article.Title,
		ShortURL:  article.ShortURL,
		Published: article.IsPublished,
		At:        now,
	})
	if article.IsPublished {
		_ = s.events.Publish(ctx, ArticlePublished{
			ID:       article.ID,
			AuthorID: article.AuthorID,
			Title:    article.Title,
			ShortURL: article.ShortURL,
			At:       now,
		})
		publishFederationSignals(ctx, s.events, article, cmd.Content)
	}

	return article, nil
}

// UpdateArticle 更新文章
func (s *Service) UpdateArticle(ctx context.Context, cmd UpdateArticleCmd) (*content.Article, error) {
	// 先获取现有文章
	existing, err := s.repo.GetArticleByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	prevPublished := existing.IsPublished
	prevContentHash := existing.ContentHash

	if cmd.CategoryID != nil {
		if _, err := s.repo.GetCategoryByID(ctx, *cmd.CategoryID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureTagsExist(ctx, cmd.TagIDs); err != nil {
		return nil, err
	}

	toc := contentutil.GenerateTOC(cmd.Content)
	summary := contentutil.BuildSummary(cmd.Summary, cmd.Content)

	// 更新字段
	existing.Title = cmd.Title
	existing.Summary = summary
	existing.LeadIn = cmd.LeadIn
	existing.TOC = toc
	existing.Content = cmd.Content
	existing.ContentHash = content.ArticleContentHash(cmd.Title, cmd.LeadIn, cmd.Content)
	existing.Cover = cmd.Cover
	existing.CategoryID = cmd.CategoryID
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

	if err := s.repo.UpdateArticle(ctx, existing); err != nil {
		return nil, err
	}

	// 同步标签
	if err := s.repo.SyncTagsToArticle(ctx, existing.ID, cmd.TagIDs); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, ArticleUpdated{
		ID:          existing.ID,
		AuthorID:    existing.AuthorID,
		Title:       existing.Title,
		ShortURL:    existing.ShortURL,
		Published:   existing.IsPublished,
		ContentHash: existing.ContentHash,
		LeadIn:      existing.LeadIn,
		TOC:         existing.TOC,
		Content:     existing.Content,
		At:          now,
	})
	if !prevPublished && existing.IsPublished {
		_ = s.events.Publish(ctx, ArticlePublished{
			ID:       existing.ID,
			AuthorID: existing.AuthorID,
			Title:    existing.Title,
			ShortURL: existing.ShortURL,
			At:       now,
		})
	}
	if prevPublished && !existing.IsPublished {
		_ = s.events.Publish(ctx, ArticleUnpublished{
			ID:       existing.ID,
			AuthorID: existing.AuthorID,
			Title:    existing.Title,
			ShortURL: existing.ShortURL,
			At:       now,
		})
	}
	if existing.IsPublished && (!prevPublished || prevContentHash != existing.ContentHash) {
		publishFederationSignals(ctx, s.events, existing, existing.Content)
	}

	return existing, nil
}

// GetArticleByID 根据 ID 获取文章
func (s *Service) GetArticleByID(ctx context.Context, id int64) (*content.Article, error) {
	return s.repo.GetArticleByID(ctx, id)
}

// GetArticleByShortURL 根据短链接获取文章
func (s *Service) GetArticleByShortURL(ctx context.Context, shortURL string) (*content.Article, error) {
	article, err := s.repo.GetArticleByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	_ = s.repo.UpdateArticleViews(ctx, article.ID)

	return article, nil
}

// ListArticles 获取文章列表
func (s *Service) ListArticles(ctx context.Context, options content.ArticleListOptionsInternal) ([]*content.Article, int64, error) {
	return s.repo.ListArticles(ctx, options)
}

// DeleteArticle 删除文章
func (s *Service) DeleteArticle(ctx context.Context, id int64) error {
	article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteArticle(ctx, id); err != nil {
		return err
	}
	_ = s.events.Publish(ctx, ArticleDeleted{
		ID:       article.ID,
		AuthorID: article.AuthorID,
		Title:    article.Title,
		ShortURL: article.ShortURL,
		At:       time.Now(),
	})
	return nil
}

// GetArticleWithTags 获取文章及其标签
func (s *Service) GetArticleWithTags(ctx context.Context, id int64) (*content.Article, []*content.Tag, error) {
	article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	tags, err := s.repo.GetTagsByArticleID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return article, tags, nil
}

// GetArticleMetrics 获取文章指标
func (s *Service) GetArticleMetrics(ctx context.Context, articleID int64) (*content.ArticleMetrics, error) {
	return s.repo.GetArticleMetrics(ctx, articleID)
}

// GetArticleTags 获取文章标签。
func (s *Service) GetArticleTags(ctx context.Context, articleID int64) ([]*content.Tag, error) {
	return s.repo.GetTagsByArticleID(ctx, articleID)
}

// generateShortURL 生成短链接
func (s *Service) ensureShortURLAvailable(ctx context.Context, shortURL string) (string, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		for i := 0; i < 5; i++ {
			candidate := contentutil.GenerateRandomShortURL()
			_, err := s.repo.GetArticleByShortURL(ctx, candidate)
			if err != nil {
				if errors.Is(err, content.ErrArticleNotFound) {
					return candidate, nil
				}
				return "", err
			}
		}
		return "", content.ErrArticleShortURLExists
	}

	existing, err := s.repo.GetArticleByShortURL(ctx, shortURL)
	if err != nil && !errors.Is(err, content.ErrArticleNotFound) {
		return "", err
	}
	if err == nil && existing != nil {
		return "", content.ErrArticleShortURLExists
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
