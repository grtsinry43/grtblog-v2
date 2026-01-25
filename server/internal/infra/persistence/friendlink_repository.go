package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type FriendLinkApplicationRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.FriendLinkApplication]
}

func NewFriendLinkApplicationRepository(db *gorm.DB) *FriendLinkApplicationRepository {
	return &FriendLinkApplicationRepository{
		db:   db,
		repo: NewGormRepository[model.FriendLinkApplication](db),
	}
}

func (r *FriendLinkApplicationRepository) FindByURL(ctx context.Context, url string) (*social.FriendLinkApplication, error) {
	rec, err := r.repo.First(ctx, "url = ?", url)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, social.ErrFriendLinkApplicationNotFound
		}
		return nil, err
	}
	entity := mapFriendLinkApplicationToDomain(*rec)
	return &entity, nil
}

func (r *FriendLinkApplicationRepository) Create(ctx context.Context, app *social.FriendLinkApplication) error {
	rec := mapFriendLinkApplicationToModel(app)
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	app.ID = rec.ID
	app.CreatedAt = rec.CreatedAt
	app.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *FriendLinkApplicationRepository) Update(ctx context.Context, app *social.FriendLinkApplication) error {
	rec := mapFriendLinkApplicationToModel(app)
	return r.db.WithContext(ctx).Model(&model.FriendLinkApplication{}).
		Where("url = ?", app.URL).
		Updates(map[string]any{
			"name":                rec.Name,
			"logo":                rec.Logo,
			"description":         rec.Description,
			"apply_channel":       rec.ApplyChannel,
			"requested_sync_mode": rec.RequestedSyncMode,
			"rss_url":             rec.RSSURL,
			"instance_url":        rec.InstanceURL,
			"manifest":            rec.Manifest,
			"signature_key_id":    rec.SignatureKeyID,
			"signature_verified":  rec.SignatureVerified,
			"user_id":             rec.UserID,
			"message":             rec.Message,
			"status":              rec.Status,
		}).Error
}

func mapFriendLinkApplicationToDomain(rec model.FriendLinkApplication) social.FriendLinkApplication {
	return social.FriendLinkApplication{
		ID:                rec.ID,
		Name:              rec.Name,
		URL:               rec.URL,
		Logo:              rec.Logo,
		Description:       rec.Description,
		ApplyChannel:      rec.ApplyChannel,
		RequestedSyncMode: rec.RequestedSyncMode,
		RSSURL:            rec.RSSURL,
		InstanceURL:       rec.InstanceURL,
		Manifest:          json.RawMessage(rec.Manifest),
		SignatureKeyID:    rec.SignatureKeyID,
		SignatureVerified: rec.SignatureVerified,
		UserID:            rec.UserID,
		Message:           rec.Message,
		Status:            rec.Status,
		CreatedAt:         rec.CreatedAt,
		UpdatedAt:         rec.UpdatedAt,
	}
}

func mapFriendLinkApplicationToModel(app *social.FriendLinkApplication) model.FriendLinkApplication {
	return model.FriendLinkApplication{
		ID:                app.ID,
		Name:              app.Name,
		URL:               app.URL,
		Logo:              app.Logo,
		Description:       app.Description,
		ApplyChannel:      app.ApplyChannel,
		RequestedSyncMode: app.RequestedSyncMode,
		RSSURL:            app.RSSURL,
		InstanceURL:       app.InstanceURL,
		Manifest:          datatypes.JSON(app.Manifest),
		SignatureKeyID:    app.SignatureKeyID,
		SignatureVerified: app.SignatureVerified,
		UserID:            app.UserID,
		Message:           app.Message,
		Status:            app.Status,
	}
}
