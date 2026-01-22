package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/media"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type UploadFileRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.UploadFile]
}

func NewUploadFileRepository(db *gorm.DB) *UploadFileRepository {
	return &UploadFileRepository{
		db:   db,
		repo: NewGormRepository[model.UploadFile](db),
	}
}

func (r *UploadFileRepository) FindByHash(ctx context.Context, hash string) (*media.UploadFile, error) {
	rec, err := r.repo.First(ctx, "hash = ?", hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, media.ErrUploadFileNotFound
		}
		return nil, err
	}
	entity := mapUploadFileToDomain(*rec)
	return &entity, nil
}

func (r *UploadFileRepository) FindByID(ctx context.Context, id int64) (*media.UploadFile, error) {
	rec, err := r.repo.FirstByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, media.ErrUploadFileNotFound
		}
		return nil, err
	}
	entity := mapUploadFileToDomain(*rec)
	return &entity, nil
}

func (r *UploadFileRepository) Create(ctx context.Context, file *media.UploadFile) error {
	rec := mapUploadFileToModel(file)
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	file.ID = rec.ID
	file.CreatedAt = rec.CreatedAt
	return nil
}

func (r *UploadFileRepository) UpdatePath(ctx context.Context, id int64, path string) error {
	return r.db.WithContext(ctx).
		Model(&model.UploadFile{}).
		Where("id = ?", id).
		Update("path", path).Error
}

func (r *UploadFileRepository) UpdateName(ctx context.Context, id int64, name string) error {
	return r.db.WithContext(ctx).
		Model(&model.UploadFile{}).
		Where("id = ?", id).
		Update("name", name).Error
}

func (r *UploadFileRepository) List(ctx context.Context, offset int, limit int) ([]media.UploadFile, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.UploadFile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []model.UploadFile
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	files := make([]media.UploadFile, len(records))
	for i, rec := range records {
		files[i] = mapUploadFileToDomain(rec)
	}
	return files, total, nil
}

func (r *UploadFileRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.UploadFile{}, id).Error
}

func mapUploadFileToDomain(rec model.UploadFile) media.UploadFile {
	return media.UploadFile{
		ID:        rec.ID,
		Name:      rec.Name,
		Path:      rec.Path,
		Type:      rec.Type,
		Size:      rec.Size,
		Hash:      rec.Hash,
		CreatedAt: rec.CreatedAt,
	}
}

func mapUploadFileToModel(file *media.UploadFile) model.UploadFile {
	return model.UploadFile{
		ID:   file.ID,
		Name: file.Name,
		Path: file.Path,
		Type: file.Type,
		Size: file.Size,
		Hash: file.Hash,
	}
}
