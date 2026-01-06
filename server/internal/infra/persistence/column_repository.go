package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type MomentColumnRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.MomentColumn]
}

func NewMomentColumnRepository(db *gorm.DB) *MomentColumnRepository {
	return &MomentColumnRepository{
		db:   db,
		repo: NewGormRepository[model.MomentColumn](db),
	}
}

func (r *MomentColumnRepository) List(ctx context.Context) ([]*content.MomentColumn, error) {
	records, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	})
	if err != nil {
		return nil, err
	}
	result := make([]*content.MomentColumn, len(records))
	for i, rec := range records {
		item := mapColumnToDomain(rec)
		result[i] = item
	}
	return result, nil
}

func (r *MomentColumnRepository) GetByID(ctx context.Context, id int64) (*content.MomentColumn, error) {
	rec, err := r.repo.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, content.ErrColumnNotFound
		}
		return nil, err
	}
	return mapColumnToDomain(*rec), nil
}

func (r *MomentColumnRepository) Create(ctx context.Context, column *content.MomentColumn) error {
	rec := model.MomentColumn{
		Name:     column.Name,
		ShortURL: optionalString(column.ShortURL),
	}
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	column.ID = rec.ID
	column.CreatedAt = rec.CreatedAt
	column.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *MomentColumnRepository) Update(ctx context.Context, column *content.MomentColumn) error {
	updates := map[string]any{
		"name":       column.Name,
		"short_url":  optionalString(column.ShortURL),
		"updated_at": time.Now(),
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

func (r *MomentColumnRepository) Delete(ctx context.Context, id int64) error {
	affected, err := r.repo.DeleteWhere(ctx, "id = ?", id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return content.ErrColumnNotFound
	}
	return nil
}

func mapColumnToDomain(rec model.MomentColumn) *content.MomentColumn {
	return &content.MomentColumn{
		ID:        rec.ID,
		Name:      rec.Name,
		ShortURL:  stringToPtr(rec.ShortURL),
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: deletedAtToPtr(rec.DeletedAt),
	}
}
