package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type SysConfigRepository struct {
	db *gorm.DB
}

func NewSysConfigRepository(db *gorm.DB) *SysConfigRepository {
	return &SysConfigRepository{db: db}
}

func (r *SysConfigRepository) GetByKey(ctx context.Context, key string) (*config.SysConfig, error) {
	var rec model.SysConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, config.ErrSysConfigNotFound
		}
		return nil, err
	}
	return &config.SysConfig{
		ID:        rec.ID,
		Key:       rec.ConfigKey,
		Value:     rec.Value,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}, nil
}
