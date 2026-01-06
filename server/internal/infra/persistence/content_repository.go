package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type ContentRepository struct {
	db *gorm.DB
}

func (r *ContentRepository) CreateCategory(ctx context.Context, category *content.ArticleCategory) error {
	rec := model.ArticleCategory{
		Name:     category.Name,
		ShortURL: optionalString(category.ShortURL),
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	category.ID = rec.ID
	category.CreatedAt = rec.CreatedAt
	category.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *ContentRepository) GetCategoryByID(ctx context.Context, id int64) (*content.ArticleCategory, error) {
	var rec model.ArticleCategory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrCategoryNotFound
		}
		return nil, err
	}
	return mapCategoryToDomain(rec), nil
}

func (r *ContentRepository) ListCategories(ctx context.Context) ([]*content.ArticleCategory, error) {
	var records []model.ArticleCategory
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*content.ArticleCategory, len(records))
	for i, rec := range records {
		result[i] = mapCategoryToDomain(rec)
	}
	return result, nil
}

func (r *ContentRepository) UpdateCategory(ctx context.Context, category *content.ArticleCategory) error {
	updates := map[string]any{
		"name":      category.Name,
		"short_url": optionalString(category.ShortURL),
	}
	result := r.db.WithContext(ctx).
		Model(&model.ArticleCategory{}).
		Where("id = ?", category.ID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrCategoryNotFound
	}
	return nil
}

func (r *ContentRepository) DeleteCategory(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ArticleCategory{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrCategoryNotFound
	}
	return nil
}

func (r *ContentRepository) CreateTag(ctx context.Context, tag *content.Tag) error {
	rec := model.Tag{Name: tag.Name}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	tag.ID = rec.ID
	tag.CreatedAt = rec.CreatedAt
	tag.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *ContentRepository) GetTagByID(ctx context.Context, id int64) (*content.Tag, error) {
	var rec model.Tag
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrTagNotFound
		}
		return nil, err
	}
	return mapTagToDomain(rec), nil
}

func (r *ContentRepository) GetTagByName(ctx context.Context, name string) (*content.Tag, error) {
	var rec model.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrTagNotFound
		}
		return nil, err
	}
	return mapTagToDomain(rec), nil
}

func (r *ContentRepository) ListTags(ctx context.Context) ([]*content.Tag, error) {
	var records []model.Tag
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*content.Tag, len(records))
	for i, rec := range records {
		result[i] = mapTagToDomain(rec)
	}
	return result, nil
}

func (r *ContentRepository) UpdateTag(ctx context.Context, tag *content.Tag) error {
	result := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", tag.ID).
		Update("name", tag.Name)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrTagNotFound
	}
	return nil
}

func (r *ContentRepository) DeleteTag(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrTagNotFound
	}
	return nil
}

func (r *ContentRepository) TagIDsExist(ctx context.Context, ids []int64) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id IN ?", ids).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count == int64(len(ids)), nil
}

func (r *ContentRepository) AddTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	records := make([]model.ArticleTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		records = append(records, model.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "article_id"}, {Name: "tag_id"}},
			DoNothing: true,
		}).
		Create(&records).Error
}

func (r *ContentRepository) SyncTagsToArticle(ctx context.Context, articleID int64, tagIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		if len(tagIDs) == 0 {
			return nil
		}
		records := make([]model.ArticleTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			records = append(records, model.ArticleTag{
				ArticleID: articleID,
				TagID:     tagID,
			})
		}
		return tx.Create(&records).Error
	})
}

func (r *ContentRepository) GetTagsByArticleID(ctx context.Context, articleID int64) ([]*content.Tag, error) {
	var records []model.Tag
	err := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Joins("JOIN article_tag ON article_tag.tag_id = tag.id").
		Where("article_tag.article_id = ?", articleID).
		Order("tag.name ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	result := make([]*content.Tag, len(records))
	for i, rec := range records {
		result[i] = mapTagToDomain(rec)
	}
	return result, nil
}

func (r *ContentRepository) UpdateArticleViews(ctx context.Context, articleID int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.ArticleMetrics{}).
		Where("article_id = ?", articleID).
		UpdateColumn("views", gorm.Expr("views + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}
	rec := model.ArticleMetrics{
		ArticleID: articleID,
		Views:     1,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "article_id"}},
			DoUpdates: clause.Assignments(map[string]any{"views": gorm.Expr("views + 1"), "updated_at": time.Now()}),
		}).
		Create(&rec).Error
}

func (r *ContentRepository) GetArticleMetrics(ctx context.Context, articleID int64) (*content.ArticleMetrics, error) {
	var rec model.ArticleMetrics
	result := r.db.WithContext(ctx).Where("article_id = ?", articleID).Limit(1).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &content.ArticleMetrics{
		ArticleID: rec.ArticleID,
		Views:     rec.Views,
		Likes:     rec.Likes,
		Comments:  rec.Comments,
		UpdatedAt: rec.UpdatedAt,
	}, nil
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

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(articleModel).Error; err != nil {
			if isArticleShortURLConstraint(err) {
				return content.ErrArticleShortURLExists
			}
			return err
		}

		metrics := model.ArticleMetrics{
			ArticleID: articleModel.ID,
			Views:     0,
			Likes:     0,
			Comments:  0,
		}
		if err := tx.Create(&metrics).Error; err != nil {
			return err
		}

		article.ID = articleModel.ID
		article.UpdatedAt = articleModel.UpdatedAt
		return nil
	})
}

// GetArticleByID 根据ID获取文章
func (r *ContentRepository) GetArticleByID(ctx context.Context, id int64) (*content.Article, error) {
	var articleModel model.Article
	result := r.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&articleModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrArticleNotFound
	}

	return r.modelToArticle(&articleModel), nil
}

// GetArticleByShortURL 根据短链接获取文章
func (r *ContentRepository) GetArticleByShortURL(ctx context.Context, shortURL string) (*content.Article, error) {
	var articleModel model.Article
	result := r.db.WithContext(ctx).Where("short_url = ?", shortURL).Limit(1).Find(&articleModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrArticleNotFound
	}

	return r.modelToArticle(&articleModel), nil
}

// UpdateArticle 更新文章
func (r *ContentRepository) UpdateArticle(ctx context.Context, article *content.Article) error {
	tocBytes, err := tocToBytes(article.TOC)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]any{
		"title":        article.Title,
		"summary":      article.Summary,
		"ai_summary":   article.AISummary,
		"lead_in":      article.LeadIn,
		"toc":          tocBytes,
		"content":      article.Content,
		"category_id":  article.CategoryID,
		"cover":        article.Cover,
		"short_url":    article.ShortURL,
		"is_published": article.IsPublished,
		"is_top":       article.IsTop,
		"is_hot":       article.IsHot,
		"is_original":  article.IsOriginal,
		"updated_at":   now,
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", article.ID).
		Updates(updates).Error; err != nil {
		if isArticleShortURLConstraint(err) {
			return content.ErrArticleShortURLExists
		}
		return err
	}

	article.UpdatedAt = now
	return nil
}

// DeleteArticle 删除文章（软删除）
func (r *ContentRepository) DeleteArticle(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&model.Article{}).Error; err != nil {
			return err
		}
		return tx.Where("article_id = ?", id).Delete(&model.ArticleMetrics{}).Error
	})
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
		// 如果解析失败，使用空列表
		toc = []content.TOCNode{}
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
func tocToBytes(toc []content.TOCNode) ([]byte, error) {
	if toc == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(toc)
}

func bytesToToc(data []byte) ([]content.TOCNode, error) {
	var toc []content.TOCNode
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return []content.TOCNode{}, nil
	}
	if err := json.Unmarshal(trimmed, &toc); err == nil {
		return toc, nil
	} else if len(trimmed) > 0 && trimmed[0] == '{' {
		return []content.TOCNode{}, nil
	} else {
		return nil, err
	}
}

func mapCategoryToDomain(rec model.ArticleCategory) *content.ArticleCategory {
	return &content.ArticleCategory{
		ID:        rec.ID,
		Name:      rec.Name,
		ShortURL:  stringToPtr(rec.ShortURL),
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: deletedAtToPtr(rec.DeletedAt),
	}
}

func mapTagToDomain(rec model.Tag) *content.Tag {
	return &content.Tag{
		ID:        rec.ID,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: deletedAtToPtr(rec.DeletedAt),
	}
}

func deletedAtToPtr(deleted gorm.DeletedAt) *time.Time {
	if !deleted.Valid {
		return nil
	}
	return &deleted.Time
}

func optionalString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func stringToPtr(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func isArticleShortURLConstraint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "uq_article_short_url")
}
