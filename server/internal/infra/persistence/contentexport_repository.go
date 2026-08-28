package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	exportdomain "github.com/grtsinry43/grtblog-v2/server/internal/domain/contentexport"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type ContentExportRepository struct{ db *gorm.DB }

func NewContentExportRepository(db *gorm.DB) *ContentExportRepository {
	return &ContentExportRepository{db: db}
}

func (r *ContentExportRepository) Create(ctx context.Context, item *exportdomain.Record) error {
	rec := exportRecordToModel(*item)
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *ContentExportRepository) Update(ctx context.Context, item *exportdomain.Record) error {
	rec := exportRecordToModel(*item)
	return r.db.WithContext(ctx).Save(&rec).Error
}

func (r *ContentExportRepository) Get(ctx context.Context, id string) (*exportdomain.Record, error) {
	var rec model.ExportRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exportdomain.ErrNotFound
		}
		return nil, err
	}
	item := exportRecordFromModel(rec)
	return &item, nil
}

func (r *ContentExportRepository) List(ctx context.Context) ([]exportdomain.Record, error) {
	var records []model.ExportRecord
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	items := make([]exportdomain.Record, len(records))
	for i, rec := range records {
		items[i] = exportRecordFromModel(rec)
	}
	return items, nil
}

func (r *ContentExportRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ExportRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exportdomain.ErrNotFound
	}
	return nil
}

func (r *ContentExportRepository) MarkInterrupted(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&model.ExportRecord{}).
		Where("status IN ?", []string{string(exportdomain.StatusQueued), string(exportdomain.StatusRunning)}).
		Updates(map[string]any{
			"status": string(exportdomain.StatusFailed), "stage": "interrupted",
			"error_message": "导出任务因服务重启而中断", "completed_at": time.Now().UTC(),
		}).Error
}

func (r *ContentExportRepository) CreateTicket(ctx context.Context, ticket exportdomain.DownloadTicket) error {
	rec := model.ExportDownloadTicket{TokenHash: ticket.TokenHash, ExportID: ticket.ExportID, ExpiresAt: ticket.ExpiresAt, CreatedAt: ticket.CreatedAt}
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *ContentExportRepository) ResolveTicket(ctx context.Context, tokenHash string) (*exportdomain.Record, error) {
	var rec model.ExportRecord
	err := r.db.WithContext(ctx).Table("export_ops.export_record AS e").
		Select("e.*").
		Joins("JOIN export_ops.download_ticket AS t ON t.export_id = e.id").
		Where("t.token_hash = ? AND t.expires_at > ? AND e.status = ?", tokenHash, time.Now().UTC(), exportdomain.StatusCompleted).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exportdomain.ErrInvalidTicket
		}
		return nil, err
	}
	item := exportRecordFromModel(rec)
	return &item, nil
}

func (r *ContentExportRepository) DeleteExpiredTickets(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at <= ?", time.Now().UTC()).Delete(&model.ExportDownloadTicket{}).Error
}

func exportRecordToModel(item exportdomain.Record) model.ExportRecord {
	return model.ExportRecord{
		ID: item.ID, Filename: item.Filename, Status: string(item.Status), Stage: item.Stage,
		TriggerType: item.TriggerType, Mode: item.Mode, SizeBytes: item.SizeBytes, SHA256: item.SHA256,
		AppVersion: item.AppVersion, SiteName: item.SiteName, SiteURL: item.SiteURL,
		ArticleCount: item.ArticleCount, MomentsCount: item.MomentsCount, PagesCount: item.PagesCount,
		ThinkingsCount: item.ThinkingsCount, ImageCount: item.ImageCount, FailedImageCount: item.FailedImageCount,
		ErrorMessage: item.ErrorMessage, CreatedAt: item.CreatedAt, StartedAt: item.StartedAt, CompletedAt: item.CompletedAt,
	}
}

func exportRecordFromModel(rec model.ExportRecord) exportdomain.Record {
	return exportdomain.Record{
		ID: rec.ID, Filename: rec.Filename, Status: exportdomain.Status(rec.Status), Stage: rec.Stage,
		TriggerType: rec.TriggerType, Mode: rec.Mode, SizeBytes: rec.SizeBytes, SHA256: rec.SHA256,
		AppVersion: rec.AppVersion, SiteName: rec.SiteName, SiteURL: rec.SiteURL,
		ArticleCount: rec.ArticleCount, MomentsCount: rec.MomentsCount, PagesCount: rec.PagesCount,
		ThinkingsCount: rec.ThinkingsCount, ImageCount: rec.ImageCount, FailedImageCount: rec.FailedImageCount,
		ErrorMessage: rec.ErrorMessage, CreatedAt: rec.CreatedAt, StartedAt: rec.StartedAt, CompletedAt: rec.CompletedAt,
	}
}
