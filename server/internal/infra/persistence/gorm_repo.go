package persistence

import (
	"context"

	"gorm.io/gorm"
)

// GormRepository 提供常见的 GORM CRUD 基础能力。
type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

// Create 持久化新实体。
func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// FirstByID 根据主键读取实体。
func (r *GormRepository[T]) FirstByID(ctx context.Context, id any) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// First 按条件读取单条记录。
func (r *GormRepository[T]) First(ctx context.Context, query any, args ...any) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where(query, args...).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 按条件读取多条记录。
func (r *GormRepository[T]) List(ctx context.Context, opts ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	query := r.db.WithContext(ctx)
	for _, opt := range opts {
		query = opt(query)
	}
	var entities []T
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Save 更新实体（根据主键）。
func (r *GormRepository[T]) Save(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// DeleteByID 根据主键删除实体。
func (r *GormRepository[T]) DeleteByID(ctx context.Context, id any) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// DeleteWhere 根据条件删除记录，返回影响行数。
func (r *GormRepository[T]) DeleteWhere(ctx context.Context, query any, args ...any) (int64, error) {
	var entity T
	result := r.db.WithContext(ctx).Where(query, args...).Delete(&entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
