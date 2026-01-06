package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type WebsiteInfoRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.WebsiteInfo]
}

func NewWebsiteInfoRepository(db *gorm.DB) *WebsiteInfoRepository {
	return &WebsiteInfoRepository{
		db:   db,
		repo: NewGormRepository[model.WebsiteInfo](db),
	}
}

func (r *WebsiteInfoRepository) List(ctx context.Context) ([]config.WebsiteInfo, error) {
	records, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("info_key ASC")
	})
	if err != nil {
		return nil, err
	}
	result := make([]config.WebsiteInfo, len(records))
	for i, rec := range records {
		result[i] = mapWebsiteInfoToDomain(rec)
	}
	return result, nil
}

func (r *WebsiteInfoRepository) GetByKey(ctx context.Context, key string) (*config.WebsiteInfo, error) {
	rec, err := r.repo.First(ctx, "info_key = ?", key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, config.ErrWebsiteInfoNotFound
		}
		return nil, err
	}
	info := mapWebsiteInfoToDomain(*rec)
	return &info, nil
}

func (r *WebsiteInfoRepository) Create(ctx context.Context, info *config.WebsiteInfo) error {
	rec := mapWebsiteInfoToModel(*info)
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	info.ID = rec.ID
	info.CreatedAt = rec.CreatedAt
	info.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *WebsiteInfoRepository) Update(ctx context.Context, info *config.WebsiteInfo) error {
	updates := map[string]any{
		"value": info.Value,
	}
	result := r.db.WithContext(ctx).
		Model(&model.WebsiteInfo{}).
		Where("info_key = ?", info.Key).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return config.ErrWebsiteInfoNotFound
	}
	return nil
}

func (r *WebsiteInfoRepository) Delete(ctx context.Context, key string) error {
	affected, err := r.repo.DeleteWhere(ctx, "info_key = ?", key)
	if err != nil {
		return err
	}
	if affected == 0 {
		return config.ErrWebsiteInfoNotFound
	}
	return nil
}

func mapWebsiteInfoToDomain(rec model.WebsiteInfo) config.WebsiteInfo {
	return config.WebsiteInfo{
		ID:        rec.ID,
		Key:       rec.InfoKey,
		Value:     rec.Value,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

func mapWebsiteInfoToModel(info config.WebsiteInfo) model.WebsiteInfo {
	return model.WebsiteInfo{
		ID:      info.ID,
		InfoKey: info.Key,
		Value:   info.Value,
	}
}
