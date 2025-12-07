package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type WebsiteInfoRepository struct {
	db *gorm.DB
}

func NewWebsiteInfoRepository(db *gorm.DB) *WebsiteInfoRepository {
	return &WebsiteInfoRepository{db: db}
}

func (r *WebsiteInfoRepository) List(ctx context.Context) ([]config.WebsiteInfo, error) {
	var records []model.WebsiteInfo
	if err := r.db.WithContext(ctx).Order("info_key ASC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]config.WebsiteInfo, len(records))
	for i, rec := range records {
		result[i] = mapWebsiteInfoToDomain(rec)
	}
	return result, nil
}

func (r *WebsiteInfoRepository) GetByKey(ctx context.Context, key string) (*config.WebsiteInfo, error) {
	var rec model.WebsiteInfo
	if err := r.db.WithContext(ctx).Where("info_key = ?", key).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, config.ErrWebsiteInfoNotFound
		}
		return nil, err
	}
	info := mapWebsiteInfoToDomain(rec)
	return &info, nil
}

func (r *WebsiteInfoRepository) Create(ctx context.Context, info *config.WebsiteInfo) error {
	rec := mapWebsiteInfoToModel(*info)
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
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
	result := r.db.WithContext(ctx).Where("info_key = ?", key).Delete(&model.WebsiteInfo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
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
