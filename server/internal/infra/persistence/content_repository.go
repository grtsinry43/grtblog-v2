package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// CreateArticle 创建文章
func (r *ContentRepository) CreateArticle(ctx context.Context, article *content.Article) error {
	tocBytes, err := tocToBytes(article.TOC)
	if err != nil {
		return err
	}

	articleModel := &model.Article{
		Title:       article.Title,
		Summary:     article.Summary,
		AISummary:   article.AISummary,
		LeadIn:      article.LeadIn,
		TOC:         tocBytes,
		Content:     article.Content,
		AuthorID:    article.AuthorID,
		Cover:       article.Cover,
		CategoryID:  article.CategoryID,
		ShortURL:    article.ShortURL,
		IsPublished: article.IsPublished,
		IsTop:       article.IsTop,
		IsHot:       article.IsHot,
		IsOriginal:  article.IsOriginal,
		CreatedAt:   article.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(articleModel).Error; err != nil {
		return err
	}

	article.ID = articleModel.ID
	article.UpdatedAt = articleModel.UpdatedAt
	return nil
}

// GetArticleByID 根据ID获取文章
func (r *ContentRepository) GetArticleByID(ctx context.Context, id int64) (*content.Article, error) {
	var articleModel model.Article
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&articleModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrArticleNotFound
		}
		return nil, err
	}

	return r.modelToArticle(&articleModel), nil
}

// GetArticleByShortURL 根据短链接获取文章
func (r *ContentRepository) GetArticleByShortURL(ctx context.Context, shortURL string) (*content.Article, error) {
	var articleModel model.Article
	if err := r.db.WithContext(ctx).Where("short_url = ?", shortURL).First(&articleModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrArticleNotFound
		}
		return nil, err
	}

	return r.modelToArticle(&articleModel), nil
}

// UpdateArticle 更新文章
func (r *ContentRepository) UpdateArticle(ctx context.Context, article *content.Article) error {
	tocBytes, err := tocToBytes(article.TOC)
	if err != nil {
		return err
	}

	articleModel := &model.Article{
		ID:          article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		AISummary:   article.AISummary,
		LeadIn:      article.LeadIn,
		TOC:         tocBytes,
		Content:     article.Content,
		CategoryID:  article.CategoryID,
		Cover:       article.Cover,
		ShortURL:    article.ShortURL,
		IsPublished: article.IsPublished,
		IsTop:       article.IsTop,
		IsHot:       article.IsHot,
		IsOriginal:  article.IsOriginal,
	}

	if err := r.db.WithContext(ctx).Save(articleModel).Error; err != nil {
		return err
	}

	article.UpdatedAt = articleModel.UpdatedAt
	return nil
}

// DeleteArticle 删除文章（软删除）
func (r *ContentRepository) DeleteArticle(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Article{}).Error
}

// ListArticles 获取文章列表（内部使用，包含未发布）
func (r *ContentRepository) ListArticles(ctx context.Context, options content.ArticleListOptionsInternal) ([]*content.Article, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Article{})

	// 应用过滤条件
	if options.CategoryID != nil {
		query = query.Where("category_id = ?", *options.CategoryID)
	}
	if options.AuthorID != nil {
		query = query.Where("author_id = ?", *options.AuthorID)
	}
	if options.Published != nil {
		query = query.Where("is_published = ?", *options.Published)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR summary ILIKE ? OR content ILIKE ?", search, search, search)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (options.Page - 1) * options.PageSize
	var articleModels []*model.Article
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&articleModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域对象
	articles := make([]*content.Article, len(articleModels))
	for i, am := range articleModels {
		articles[i] = r.modelToArticle(am)
	}

	return articles, total, nil
}

// ListPublicArticles 获取公开文章列表
func (r *ContentRepository) ListPublicArticles(ctx context.Context, options content.ArticleListOptions) ([]*content.Article, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Article{}).Where("is_published = ?", true)

	// 应用过滤条件
	if options.CategoryID != nil {
		query = query.Where("category_id = ?", *options.CategoryID)
	}
	if options.AuthorID != nil {
		query = query.Where("author_id = ?", *options.AuthorID)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR summary ILIKE ?", search, search)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	offset := (options.Page - 1) * options.PageSize
	var articleModels []*model.Article
	if err := query.Order("is_top DESC, created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&articleModels).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域对象
	articles := make([]*content.Article, len(articleModels))
	for i, am := range articleModels {
		articles[i] = r.modelToArticle(am)
	}

	return articles, total, nil
}

// modelToArticle 将数据库模型转换为领域对象
func (r *ContentRepository) modelToArticle(am *model.Article) *content.Article {
	toc, err := bytesToToc(am.TOC)
	if err != nil {
		// 如果解析失败，使用空的map
		toc = make(map[string]any)
	}

	return &content.Article{
		ID:          am.ID,
		Title:       am.Title,
		Summary:     am.Summary,
		AISummary:   am.AISummary,
		LeadIn:      am.LeadIn,
		TOC:         toc,
		Content:     am.Content,
		AuthorID:    am.AuthorID,
		Cover:       am.Cover,
		CategoryID:  am.CategoryID,
		CommentID:   am.CommentID,
		ShortURL:    am.ShortURL,
		IsPublished: am.IsPublished,
		IsTop:       am.IsTop,
		IsHot:       am.IsHot,
		IsOriginal:  am.IsOriginal,
		CreatedAt:   am.CreatedAt,
		UpdatedAt:   am.UpdatedAt,
		DeletedAt:   timeToTimePtr(am.DeletedAt.Time),
	}
}

// 辅助函数
func timeToTimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// JSONB 转换辅助函数
func tocToBytes(toc map[string]any) ([]byte, error) {
	if toc == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(toc)
}

func bytesToToc(data []byte) (map[string]any, error) {
	var toc map[string]any
	if len(data) == 0 {
		return make(map[string]any), nil
	}
	err := json.Unmarshal(data, &toc)
	return toc, err
}
