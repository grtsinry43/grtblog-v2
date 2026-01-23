package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type SysConfigRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.SysConfig]
}

func NewSysConfigRepository(db *gorm.DB) *SysConfigRepository {
	return &SysConfigRepository{
		db:   db,
		repo: NewGormRepository[model.SysConfig](db),
	}
}

func (r *SysConfigRepository) GetByKey(ctx context.Context, key string) (*config.SysConfig, error) {
	rec, err := r.repo.First(ctx, "config_key = ?", key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, config.ErrSysConfigNotFound
		}
		return nil, err
	}
	return &config.SysConfig{
		ID:           rec.ID,
		Key:          rec.ConfigKey,
		Value:        rec.Value,
		IsSensitive:  rec.IsSensitive,
		GroupPath:    rec.GroupPath,
		Label:        rec.Label,
		Description:  rec.Description,
		ValueType:    rec.ValueType,
		EnumOptions:  json.RawMessage(rec.EnumOptions),
		DefaultValue: rec.DefaultValue,
		VisibleWhen:  json.RawMessage(rec.VisibleWhen),
		Sort:         rec.Sort,
		Meta:         json.RawMessage(rec.Meta),
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}, nil
}

func (r *SysConfigRepository) List(ctx context.Context, keys []string) ([]config.SysConfig, error) {
	records, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		if len(keys) > 0 {
			return db.Where("config_key IN ?", keys).Order("group_path").Order("sort").Order("config_key")
		}
		return db.Order("group_path").Order("sort").Order("config_key")
	})
	if err != nil {
		return nil, err
	}
	result := make([]config.SysConfig, len(records))
	for i, rec := range records {
		result[i] = config.SysConfig{
			ID:           rec.ID,
			Key:          rec.ConfigKey,
			Value:        rec.Value,
			IsSensitive:  rec.IsSensitive,
			GroupPath:    rec.GroupPath,
			Label:        rec.Label,
			Description:  rec.Description,
			ValueType:    rec.ValueType,
			EnumOptions:  json.RawMessage(rec.EnumOptions),
			DefaultValue: rec.DefaultValue,
			VisibleWhen:  json.RawMessage(rec.VisibleWhen),
			Sort:         rec.Sort,
			Meta:         json.RawMessage(rec.Meta),
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
		}
	}
	return result, nil
}

func (r *SysConfigRepository) Upsert(ctx context.Context, configs []config.SysConfig) error {
	if len(configs) == 0 {
		return nil
	}
	now := time.Now()
	records := make([]model.SysConfig, len(configs))
	for i, cfg := range configs {
		records[i] = model.SysConfig{
			ConfigKey:    cfg.Key,
			Value:        cfg.Value,
			IsSensitive:  cfg.IsSensitive,
			GroupPath:    cfg.GroupPath,
			Label:        cfg.Label,
			Description:  cfg.Description,
			ValueType:    cfg.ValueType,
			EnumOptions:  datatypes.JSON(cfg.EnumOptions),
			DefaultValue: cfg.DefaultValue,
			VisibleWhen:  datatypes.JSON(cfg.VisibleWhen),
			Sort:         cfg.Sort,
			Meta:         datatypes.JSON(cfg.Meta),
			UpdatedAt:    now,
		}
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "config_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"value",
				"is_sensitive",
				"group_path",
				"label",
				"description",
				"value_type",
				"enum_options",
				"default_value",
				"visible_when",
				"sort",
				"meta",
				"updated_at",
			}),
		}).
		Create(&records).Error
}
