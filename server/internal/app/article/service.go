package article

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type Service struct {
	repo content.Repository
}

func NewService(repo content.Repository) *Service {
	return &Service{repo: repo}
}

// CreateArticle 创建文章
func (s *Service) CreateArticle(ctx context.Context, authorID int64, cmd CreateArticleCommand) (*content.Article, error) {
	// 生成短链接如果未提供
	shortURL := ""
	if cmd.ShortURL != nil {
		shortURL = *cmd.ShortURL
	} else {
		shortURL = s.generateShortURL()
	}

	// 设置创建时间
	createdAt := time.Now()
	if cmd.CreatedAt != nil {
		createdAt = *cmd.CreatedAt
	}

	// 处理 TOC
	toc := cmd.TOC
	if toc == nil {
		toc = make(map[string]any)
	}

	article := &content.Article{
		Title:       cmd.Title,
		Summary:     cmd.Summary,
		AISummary:   nil, // 后续可以添加 AI 摘要功能
		LeadIn:      cmd.LeadIn,
		TOC:         toc,
		Content:     cmd.Content,
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

	return article, nil
}

// UpdateArticle 更新文章
func (s *Service) UpdateArticle(ctx context.Context, cmd UpdateArticleCommand) (*content.Article, error) {
	// 先获取现有文章
	existing, err := s.repo.GetArticleByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 处理 TOC
	toc := cmd.TOC
	if toc == nil {
		toc = make(map[string]any)
	}

	// 更新字段
	existing.Title = cmd.Title
	existing.Summary = cmd.Summary
	existing.LeadIn = cmd.LeadIn
	existing.TOC = toc
	existing.Content = cmd.Content
	existing.Cover = cmd.Cover
	existing.CategoryID = cmd.CategoryID
	existing.ShortURL = cmd.ShortURL
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
	return s.repo.DeleteArticle(ctx, id)
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

// ToResponse 将领域对象转换为响应 DTO
func (s *Service) ToResponse(ctx context.Context, article *content.Article) (*ViewArticleResponse, error) {
	// 获取标签
	tags, err := s.repo.GetTagsByArticleID(ctx, article.ID)
	if err != nil {
		return nil, err
	}

	tagResponses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	// 获取指标
	metrics, err := s.repo.GetArticleMetrics(ctx, article.ID)
	var metricsResponse *MetricsResponse
	if err == nil && metrics != nil {
		metricsResponse = &MetricsResponse{
			Views:    metrics.Views,
			Likes:    metrics.Likes,
			Comments: metrics.Comments,
		}
	}

	return &ViewArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		AISummary:   article.AISummary,
		LeadIn:      article.LeadIn,
		TOC:         article.TOC,
		Content:     article.Content,
		AuthorID:    article.AuthorID,
		Cover:       article.Cover,
		CategoryID:  article.CategoryID,
		ShortURL:    article.ShortURL,
		IsPublished: article.IsPublished,
		IsTop:       article.IsTop,
		IsHot:       article.IsHot,
		IsOriginal:  article.IsOriginal,
		Tags:        tagResponses,
		Metrics:     metricsResponse,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}

// generateShortURL 生成短链接
func (s *Service) generateShortURL() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
