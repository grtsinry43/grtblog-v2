package article

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
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
func (s *Service) CreateArticle(ctx context.Context, authorID int64, cmd CreateArticleCmd) (*content.Article, error) {
	// 生成短链接如果未提供
	shortURL := ""
	if cmd.ShortURL != nil {
		shortURL = strings.TrimSpace(*cmd.ShortURL)
	}
	if shortURL == "" {
		for i := 0; i < 5; i++ {
			candidate := s.generateShortURL()
			_, err := s.repo.GetArticleByShortURL(ctx, candidate)
			if err != nil {
				if errors.Is(err, content.ErrArticleNotFound) {
					shortURL = candidate
					break
				}
				return nil, err
			}
		}
		if shortURL == "" {
			return nil, content.ErrArticleShortURLExists
		}
	} else {
		existing, err := s.repo.GetArticleByShortURL(ctx, shortURL)
		if err != nil && !errors.Is(err, content.ErrArticleNotFound) {
			return nil, err
		}
		if err == nil && existing != nil {
			return nil, content.ErrArticleShortURLExists
		}
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
func (s *Service) UpdateArticle(ctx context.Context, cmd UpdateArticleCmd) (*content.Article, error) {
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
	if cmd.ShortURL != "" && cmd.ShortURL != existing.ShortURL {
		other, err := s.repo.GetArticleByShortURL(ctx, cmd.ShortURL)
		if err != nil && !errors.Is(err, content.ErrArticleNotFound) {
			return nil, err
		}
		if err == nil && other != nil && other.ID != existing.ID {
			return nil, content.ErrArticleShortURLExists
		}
	}
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

// GetArticleTags 获取文章标签。
func (s *Service) GetArticleTags(ctx context.Context, articleID int64) ([]*content.Tag, error) {
	return s.repo.GetTagsByArticleID(ctx, articleID)
}

// generateShortURL 生成短链接
func (s *Service) generateShortURL() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
