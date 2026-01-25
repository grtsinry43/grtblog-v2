package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type FriendLinkRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.FriendLink]
}

func NewFriendLinkRepository(db *gorm.DB) *FriendLinkRepository {
	return &FriendLinkRepository{
		db:   db,
		repo: NewGormRepository[model.FriendLink](db),
	}
}

func (r *FriendLinkRepository) FindByURL(ctx context.Context, url string) (*social.FriendLink, error) {
	rec, err := r.repo.First(ctx, "url = ?", url)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, social.ErrFriendLinkNotFound
		}
		return nil, err
	}
	entity := mapFriendLinkToDomain(*rec)
	return &entity, nil
}

func (r *FriendLinkRepository) Create(ctx context.Context, link *social.FriendLink) error {
	rec := mapFriendLinkToModel(link)
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	link.ID = rec.ID
	link.CreatedAt = rec.CreatedAt
	link.UpdatedAt = rec.UpdatedAt
	return nil
}

func mapFriendLinkToDomain(rec model.FriendLink) social.FriendLink {
	return social.FriendLink{
		ID:               rec.ID,
		Name:             rec.Name,
		URL:              rec.URL,
		Logo:             stringToPtr(rec.Logo),
		Description:      stringToPtr(rec.Description),
		RSSURL:           stringToPtr(rec.RSSURL),
		Kind:             rec.Kind,
		SyncMode:         rec.SyncMode,
		InstanceID:       rec.InstanceID,
		LastSyncAt:       rec.LastSyncAt,
		LastSyncStatus:   rec.LastSyncStatus,
		SyncInterval:     rec.SyncInterval,
		TotalPostsCached: rec.TotalPostsCached,
		UserID:           rec.UserID,
		IsActive:         rec.IsActive,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        rec.UpdatedAt,
		DeletedAt:        deletedAtToPtr(rec.DeletedAt),
	}
}

func mapFriendLinkToModel(link *social.FriendLink) model.FriendLink {
	return model.FriendLink{
		ID:               link.ID,
		Name:             link.Name,
		URL:              link.URL,
		Logo:             optionalString(link.Logo),
		Description:      optionalString(link.Description),
		RSSURL:           optionalString(link.RSSURL),
		Kind:             link.Kind,
		SyncMode:         link.SyncMode,
		InstanceID:       link.InstanceID,
		LastSyncAt:       link.LastSyncAt,
		LastSyncStatus:   link.LastSyncStatus,
		SyncInterval:     link.SyncInterval,
		TotalPostsCached: link.TotalPostsCached,
		UserID:           link.UserID,
		IsActive:         link.IsActive,
	}
}
