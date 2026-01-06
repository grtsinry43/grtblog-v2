package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type ArticleCategoryRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.ArticleCategory]
}

func NewArticleCategoryRepository(db *gorm.DB) *ArticleCategoryRepository {
	return &ArticleCategoryRepository{
		db:   db,
		repo: NewGormRepository[model.ArticleCategory](db),
	}
}

func (r *ArticleCategoryRepository) List(ctx context.Context) ([]*content.ArticleCategory, error) {
	records, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	})
	if err != nil {
		return nil, err
	}
	result := make([]*content.ArticleCategory, len(records))
	for i, rec := range records {
		item := mapCategoryToDomain(rec)
		result[i] = item
	}
	return result, nil
}

func (r *ArticleCategoryRepository) GetByID(ctx context.Context, id int64) (*content.ArticleCategory, error) {
	rec, err := r.repo.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrCategoryNotFound
		}
		return nil, err
	}
	return mapCategoryToDomain(*rec), nil
}

func (r *ArticleCategoryRepository) Create(ctx context.Context, category *content.ArticleCategory) error {
	rec := model.ArticleCategory{
		Name:     category.Name,
		ShortURL: optionalString(category.ShortURL),
	}
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	category.ID = rec.ID
	category.CreatedAt = rec.CreatedAt
	category.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *ArticleCategoryRepository) Update(ctx context.Context, category *content.ArticleCategory) error {
	updates := map[string]any{
		"name":       category.Name,
		"short_url":  optionalString(category.ShortURL),
		"updated_at": time.Now(),
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

func (r *ArticleCategoryRepository) Delete(ctx context.Context, id int64) error {
	affected, err := r.repo.DeleteWhere(ctx, "id = ?", id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return content.ErrCategoryNotFound
	}
	return nil
}
