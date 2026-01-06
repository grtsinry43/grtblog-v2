package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type TagRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.Tag]
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		db:   db,
		repo: NewGormRepository[model.Tag](db),
	}
}

func (r *TagRepository) List(ctx context.Context) ([]*content.Tag, error) {
	records, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	})
	if err != nil {
		return nil, err
	}
	result := make([]*content.Tag, len(records))
	for i, rec := range records {
		item := mapTagToDomain(rec)
		result[i] = item
	}
	return result, nil
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (*content.Tag, error) {
	rec, err := r.repo.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrTagNotFound
		}
		return nil, err
	}
	return mapTagToDomain(*rec), nil
}

func (r *TagRepository) Create(ctx context.Context, tag *content.Tag) error {
	rec := model.Tag{
		Name: tag.Name,
	}
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	tag.ID = rec.ID
	tag.CreatedAt = rec.CreatedAt
	tag.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *TagRepository) Update(ctx context.Context, tag *content.Tag) error {
	updates := map[string]any{
		"name":       tag.Name,
		"updated_at": time.Now(),
	}
	result := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", tag.ID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return content.ErrTagNotFound
	}
	return nil
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	affected, err := r.repo.DeleteWhere(ctx, "id = ?", id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return content.ErrTagNotFound
	}
	return nil
}
