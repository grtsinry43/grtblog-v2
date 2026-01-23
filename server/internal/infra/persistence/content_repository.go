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

func (r *ContentRepository) CreateColumn(ctx context.Context, column *content.MomentColumn) error {
	rec := model.MomentColumn{
		Name:     column.Name,
		ShortURL: optionalString(column.ShortURL),
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	column.ID = rec.ID
	column.CreatedAt = rec.CreatedAt
	column.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *ContentRepository) GetColumnByID(ctx context.Context, id int64) (*content.MomentColumn, error) {
	var rec model.MomentColumn
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrColumnNotFound
		}
		return nil, err
	}
	return mapColumnToDomain(rec), nil
}

func (r *ContentRepository) ListColumns(ctx context.Context) ([]*content.MomentColumn, error) {
	var records []model.MomentColumn
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*content.MomentColumn, len(records))
	for i, rec := range records {
		result[i] = mapColumnToDomain(rec)
	}
	return result, nil
}

func (r *ContentRepository) UpdateColumn(ctx context.Context, column *content.MomentColumn) error {
	updates := map[string]any{
		"name":      column.Name,
		"short_url": optionalString(column.ShortURL),
	}
	result := r.db.WithContext(ctx).
		Model(&model.MomentColumn{}).
		Where("id = ?", column.ID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrColumnNotFound
	}
	return nil
}

func (r *ContentRepository) DeleteColumn(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.MomentColumn{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrColumnNotFound
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

func (r *ContentRepository) AddTopicsToMoment(ctx context.Context, momentID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	records := make([]model.MomentTopic, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		records = append(records, model.MomentTopic{
			MomentID: momentID,
			TagID:    tagID,
		})
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "moment_id"}, {Name: "tag_id"}},
			DoNothing: true,
		}).
		Create(&records).Error
}

func (r *ContentRepository) SyncTopicsToMoment(ctx context.Context, momentID int64, tagIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("moment_id = ?", momentID).Delete(&model.MomentTopic{}).Error; err != nil {
			return err
		}
		if len(tagIDs) == 0 {
			return nil
		}
		records := make([]model.MomentTopic, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			records = append(records, model.MomentTopic{
				MomentID: momentID,
				TagID:    tagID,
			})
		}
		return tx.Create(&records).Error
	})
}

func (r *ContentRepository) GetTopicsByMomentID(ctx context.Context, momentID int64) ([]*content.Tag, error) {
	var records []model.Tag
	err := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Joins("JOIN moment_topic ON moment_topic.tag_id = tag.id").
		Where("moment_topic.moment_id = ?", momentID).
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

func (r *ContentRepository) UpdateMomentViews(ctx context.Context, momentID int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.MomentMetrics{}).
		Where("moment_id = ?", momentID).
		UpdateColumn("views", gorm.Expr("views + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}
	rec := model.MomentMetrics{
		MomentID: momentID,
		Views:    1,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "moment_id"}},
			DoUpdates: clause.Assignments(map[string]any{"views": gorm.Expr("views + 1"), "updated_at": time.Now()}),
		}).
		Create(&rec).Error
}

func (r *ContentRepository) GetMomentMetrics(ctx context.Context, momentID int64) (*content.MomentMetrics, error) {
	var rec model.MomentMetrics
	result := r.db.WithContext(ctx).Where("moment_id = ?", momentID).Limit(1).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &content.MomentMetrics{
		MomentID:  rec.MomentID,
		Views:     rec.Views,
		Likes:     rec.Likes,
		Comments:  rec.Comments,
		UpdatedAt: rec.UpdatedAt,
	}, nil
}

func (r *ContentRepository) UpdatePageViews(ctx context.Context, pageID int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.PageMetrics{}).
		Where("page_id = ?", pageID).
		UpdateColumn("views", gorm.Expr("views + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}
	rec := model.PageMetrics{
		PageID: pageID,
		Views:  1,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "page_id"}},
			DoUpdates: clause.Assignments(map[string]any{"views": gorm.Expr("views + 1"), "updated_at": time.Now()}),
		}).
		Create(&rec).Error
}

func (r *ContentRepository) GetPageMetrics(ctx context.Context, pageID int64) (*content.PageMetrics, error) {
	var rec model.PageMetrics
	result := r.db.WithContext(ctx).Where("page_id = ?", pageID).Limit(1).Find(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &content.PageMetrics{
		PageID:    rec.PageID,
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
		ContentHash: article.ContentHash,
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
		"content_hash": article.ContentHash,
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

// CreateMoment 创建手记
func (r *ContentRepository) CreateMoment(ctx context.Context, moment *content.Moment) error {
	tocBytes, err := tocToBytes(moment.TOC)
	if err != nil {
		return err
	}

	momentModel := &model.Moment{
		Title:       moment.Title,
		Summary:     moment.Summary,
		AISummary:   moment.AISummary,
		TOC:         tocBytes,
		Content:     moment.Content,
		ContentHash: moment.ContentHash,
		AuthorID:    moment.AuthorID,
		Image:       moment.Image,
		ColumnID:    moment.ColumnID,
		ShortURL:    moment.ShortURL,
		IsPublished: moment.IsPublished,
		IsTop:       moment.IsTop,
		IsHot:       moment.IsHot,
		IsOriginal:  moment.IsOriginal,
		CreatedAt:   moment.CreatedAt,
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(momentModel).Error; err != nil {
			if isMomentShortURLConstraint(err) {
				return content.ErrMomentShortURLExists
			}
			return err
		}

		metrics := model.MomentMetrics{
			MomentID: momentModel.ID,
			Views:    0,
			Likes:    0,
			Comments: 0,
		}
		if err := tx.Create(&metrics).Error; err != nil {
			return err
		}

		moment.ID = momentModel.ID
		moment.UpdatedAt = momentModel.UpdatedAt
		return nil
	})
}

// GetMomentByID 根据ID获取手记
func (r *ContentRepository) GetMomentByID(ctx context.Context, id int64) (*content.Moment, error) {
	var momentModel model.Moment
	result := r.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&momentModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrMomentNotFound
	}

	return r.modelToMoment(&momentModel), nil
}

// GetMomentByShortURL 根据短链接获取手记
func (r *ContentRepository) GetMomentByShortURL(ctx context.Context, shortURL string) (*content.Moment, error) {
	var momentModel model.Moment
	result := r.db.WithContext(ctx).Where("short_url = ?", shortURL).Limit(1).Find(&momentModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrMomentNotFound
	}

	return r.modelToMoment(&momentModel), nil
}

// UpdateMoment 更新手记
func (r *ContentRepository) UpdateMoment(ctx context.Context, moment *content.Moment) error {
	tocBytes, err := tocToBytes(moment.TOC)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]any{
		"title":        moment.Title,
		"summary":      moment.Summary,
		"ai_summary":   moment.AISummary,
		"toc":          tocBytes,
		"content":      moment.Content,
		"content_hash": moment.ContentHash,
		"column_id":    moment.ColumnID,
		"img":          moment.Image,
		"short_url":    moment.ShortURL,
		"is_published": moment.IsPublished,
		"is_top":       moment.IsTop,
		"is_hot":       moment.IsHot,
		"is_original":  moment.IsOriginal,
		"updated_at":   now,
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Moment{}).
		Where("id = ?", moment.ID).
		Updates(updates).Error; err != nil {
		if isMomentShortURLConstraint(err) {
			return content.ErrMomentShortURLExists
		}
		return err
	}

	moment.UpdatedAt = now
	return nil
}

// DeleteMoment 删除手记（软删除）
func (r *ContentRepository) DeleteMoment(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&model.Moment{}).Error; err != nil {
			return err
		}
		return tx.Where("moment_id = ?", id).Delete(&model.MomentMetrics{}).Error
	})
}

// ListMoments 获取手记列表（内部使用，包含未发布）
func (r *ContentRepository) ListMoments(ctx context.Context, options content.MomentListOptionsInternal) ([]*content.Moment, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Moment{})

	if options.ColumnID != nil {
		query = query.Where("column_id = ?", *options.ColumnID)
	}
	if options.AuthorID != nil {
		query = query.Where("author_id = ?", *options.AuthorID)
	}
	if options.Published != nil {
		query = query.Where("is_published = ?", *options.Published)
	}
	if options.TopicID != nil {
		subQuery := r.db.WithContext(ctx).
			Model(&model.MomentTopic{}).
			Select("moment_id").
			Where("tag_id = ?", *options.TopicID)
		query = query.Where("id IN (?)", subQuery)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR summary ILIKE ? OR content ILIKE ?", search, search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (options.Page - 1) * options.PageSize
	var momentModels []*model.Moment
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&momentModels).Error; err != nil {
		return nil, 0, err
	}

	moments := make([]*content.Moment, len(momentModels))
	for i, mm := range momentModels {
		moments[i] = r.modelToMoment(mm)
	}

	return moments, total, nil
}

// ListPublicMoments 获取公开手记列表
func (r *ContentRepository) ListPublicMoments(ctx context.Context, options content.MomentListOptions) ([]*content.Moment, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Moment{}).Where("is_published = ?", true)

	if options.ColumnID != nil {
		query = query.Where("column_id = ?", *options.ColumnID)
	}
	if options.AuthorID != nil {
		query = query.Where("author_id = ?", *options.AuthorID)
	}
	if options.TopicID != nil {
		subQuery := r.db.WithContext(ctx).
			Model(&model.MomentTopic{}).
			Select("moment_id").
			Where("tag_id = ?", *options.TopicID)
		query = query.Where("id IN (?)", subQuery)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR summary ILIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (options.Page - 1) * options.PageSize
	var momentModels []*model.Moment
	if err := query.Order("is_top DESC, created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&momentModels).Error; err != nil {
		return nil, 0, err
	}

	moments := make([]*content.Moment, len(momentModels))
	for i, mm := range momentModels {
		moments[i] = r.modelToMoment(mm)
	}

	return moments, total, nil
}

// CreatePage 创建页面
func (r *ContentRepository) CreatePage(ctx context.Context, page *content.Page) error {
	tocBytes, err := tocToBytes(page.TOC)
	if err != nil {
		return err
	}

	pageModel := &model.Page{
		Title:       page.Title,
		Description: optionalString(page.Description),
		AISummary:   optionalString(page.AISummary),
		TOC:         tocBytes,
		Content:     page.Content,
		ContentHash: page.ContentHash,
		ShortURL:    page.ShortURL,
		IsEnabled:   page.IsEnabled,
		IsBuiltin:   page.IsBuiltin,
		CreatedAt:   page.CreatedAt,
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pageModel).Error; err != nil {
			if isPageShortURLConstraint(err) {
				return content.ErrPageShortURLExists
			}
			return err
		}

		metrics := model.PageMetrics{
			PageID:   pageModel.ID,
			Views:    0,
			Likes:    0,
			Comments: 0,
		}
		if err := tx.Create(&metrics).Error; err != nil {
			return err
		}

		page.ID = pageModel.ID
		page.UpdatedAt = pageModel.UpdatedAt
		return nil
	})
}

// GetPageByID 根据ID获取页面
func (r *ContentRepository) GetPageByID(ctx context.Context, id int64) (*content.Page, error) {
	var pageModel model.Page
	result := r.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&pageModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrPageNotFound
	}

	return r.modelToPage(&pageModel), nil
}

// GetPageByShortURL 根据短链接获取页面
func (r *ContentRepository) GetPageByShortURL(ctx context.Context, shortURL string) (*content.Page, error) {
	var pageModel model.Page
	result := r.db.WithContext(ctx).Where("short_url = ?", shortURL).Limit(1).Find(&pageModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, content.ErrPageNotFound
	}

	return r.modelToPage(&pageModel), nil
}

// UpdatePage 更新页面
func (r *ContentRepository) UpdatePage(ctx context.Context, page *content.Page) error {
	tocBytes, err := tocToBytes(page.TOC)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]any{
		"title":        page.Title,
		"description":  optionalString(page.Description),
		"ai_summary":   optionalString(page.AISummary),
		"toc":          tocBytes,
		"content":      page.Content,
		"content_hash": page.ContentHash,
		"short_url":    page.ShortURL,
		"is_enabled":   page.IsEnabled,
		"is_builtin":   page.IsBuiltin,
		"updated_at":   now,
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Page{}).
		Where("id = ?", page.ID).
		Updates(updates).Error; err != nil {
		if isPageShortURLConstraint(err) {
			return content.ErrPageShortURLExists
		}
		return err
	}

	page.UpdatedAt = now
	return nil
}

// DeletePage 删除页面（软删除）
func (r *ContentRepository) DeletePage(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&model.Page{}).Error; err != nil {
			return err
		}
		return tx.Where("page_id = ?", id).Delete(&model.PageMetrics{}).Error
	})
}

// ListPages 获取页面列表（内部使用，包含管理功能）
func (r *ContentRepository) ListPages(ctx context.Context, options content.PageListOptionsInternal) ([]*content.Page, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Page{})

	if options.Enabled != nil {
		query = query.Where("is_enabled = ?", *options.Enabled)
	}
	if options.Builtin != nil {
		query = query.Where("is_builtin = ?", *options.Builtin)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (options.Page - 1) * options.PageSize
	var pageModels []*model.Page
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&pageModels).Error; err != nil {
		return nil, 0, err
	}

	pages := make([]*content.Page, len(pageModels))
	for i, pm := range pageModels {
		pages[i] = r.modelToPage(pm)
	}

	return pages, total, nil
}

// ListPublicPages 获取公开页面列表
func (r *ContentRepository) ListPublicPages(ctx context.Context, options content.PageListOptions) ([]*content.Page, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Page{})

	if options.Enabled != nil {
		query = query.Where("is_enabled = ?", *options.Enabled)
	}
	if options.Builtin != nil {
		query = query.Where("is_builtin = ?", *options.Builtin)
	}
	if options.Search != nil && *options.Search != "" {
		search := "%" + *options.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (options.Page - 1) * options.PageSize
	var pageModels []*model.Page
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&pageModels).Error; err != nil {
		return nil, 0, err
	}

	pages := make([]*content.Page, len(pageModels))
	for i, pm := range pageModels {
		pages[i] = r.modelToPage(pm)
	}

	return pages, total, nil
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
		ContentHash: am.ContentHash,
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

// modelToMoment 将数据库模型转换为领域对象
func (r *ContentRepository) modelToMoment(mm *model.Moment) *content.Moment {
	toc, err := bytesToToc(mm.TOC)
	if err != nil {
		toc = []content.TOCNode{}
	}

	return &content.Moment{
		ID:          mm.ID,
		Title:       mm.Title,
		Summary:     mm.Summary,
		AISummary:   mm.AISummary,
		Content:     mm.Content,
		ContentHash: mm.ContentHash,
		AuthorID:    mm.AuthorID,
		TOC:         toc,
		Image:       mm.Image,
		ColumnID:    mm.ColumnID,
		CommentID:   mm.CommentID,
		ShortURL:    mm.ShortURL,
		IsPublished: mm.IsPublished,
		IsTop:       mm.IsTop,
		IsHot:       mm.IsHot,
		IsOriginal:  mm.IsOriginal,
		CreatedAt:   mm.CreatedAt,
		UpdatedAt:   mm.UpdatedAt,
		DeletedAt:   timeToTimePtr(mm.DeletedAt.Time),
	}
}

// modelToPage 将数据库模型转换为领域对象
func (r *ContentRepository) modelToPage(pm *model.Page) *content.Page {
	toc, err := bytesToToc(pm.TOC)
	if err != nil {
		toc = []content.TOCNode{}
	}

	return &content.Page{
		ID:          pm.ID,
		Title:       pm.Title,
		Description: stringToPtr(pm.Description),
		AISummary:   stringToPtr(pm.AISummary),
		TOC:         toc,
		Content:     pm.Content,
		ContentHash: pm.ContentHash,
		CommentID:   pm.CommentID,
		ShortURL:    pm.ShortURL,
		IsEnabled:   pm.IsEnabled,
		IsBuiltin:   pm.IsBuiltin,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
		DeletedAt:   timeToTimePtr(pm.DeletedAt.Time),
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

func isMomentShortURLConstraint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "uq_moment_short_url")
}

func isPageShortURLConstraint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "uq_page_short_url")
}
