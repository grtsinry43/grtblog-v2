package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

// FederationInstanceRepository handles remote instance records.
type FederationInstanceRepository struct {
	db   *gorm.DB
	repo *GormRepository[model.FederationInstance]
}

func NewFederationInstanceRepository(db *gorm.DB) *FederationInstanceRepository {
	return &FederationInstanceRepository{
		db:   db,
		repo: NewGormRepository[model.FederationInstance](db),
	}
}

func (r *FederationInstanceRepository) GetByBaseURL(ctx context.Context, baseURL string) (*federation.FederationInstance, error) {
	rec, err := r.repo.First(ctx, "base_url = ?", baseURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, federation.ErrFederationInstanceNotFound
		}
		return nil, err
	}
	instance := mapFederationInstanceToDomain(*rec)
	return &instance, nil
}

func (r *FederationInstanceRepository) Create(ctx context.Context, instance *federation.FederationInstance) error {
	rec := mapFederationInstanceToModel(instance)
	if err := r.repo.Create(ctx, &rec); err != nil {
		return err
	}
	instance.ID = rec.ID
	instance.CreatedAt = rec.CreatedAt
	instance.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *FederationInstanceRepository) Update(ctx context.Context, instance *federation.FederationInstance) error {
	rec := mapFederationInstanceToModel(instance)
	return r.db.WithContext(ctx).Model(&model.FederationInstance{}).
		Where("id = ?", instance.ID).
		Updates(&rec).Error
}

func (r *FederationInstanceRepository) ListActive(ctx context.Context) ([]federation.FederationInstance, error) {
	recs, err := r.repo.List(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active").Order("updated_at DESC")
	})
	if err != nil {
		return nil, err
	}
	result := make([]federation.FederationInstance, len(recs))
	for i, rec := range recs {
		result[i] = mapFederationInstanceToDomain(rec)
	}
	return result, nil
}

// FederatedPostCacheRepository stores cached timeline posts.
type FederatedPostCacheRepository struct {
	db *gorm.DB
}

func NewFederatedPostCacheRepository(db *gorm.DB) *FederatedPostCacheRepository {
	return &FederatedPostCacheRepository{db: db}
}

func (r *FederatedPostCacheRepository) UpsertBatch(ctx context.Context, posts []federation.FederatedPostCache) error {
	if len(posts) == 0 {
		return nil
	}
	recs := make([]model.FederatedPostCache, len(posts))
	for i := range posts {
		recs[i] = mapFederatedPostCacheToModel(&posts[i])
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		UpdateAll: true,
	}).Create(&recs).Error
}

func (r *FederatedPostCacheRepository) ListByInstance(ctx context.Context, instanceID int64, since *time.Time, limit int) ([]federation.FederatedPostCache, error) {
	query := r.db.WithContext(ctx).Where("instance_id = ?", instanceID)
	if since != nil {
		query = query.Where("published_at >= ?", *since)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	query = query.Order("published_at DESC")
	var recs []model.FederatedPostCache
	if err := query.Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]federation.FederatedPostCache, len(recs))
	for i, rec := range recs {
		result[i] = mapFederatedPostCacheToDomain(rec)
	}
	return result, nil
}

func (r *FederatedPostCacheRepository) ListRecent(ctx context.Context, limit int) ([]federation.FederatedPostCache, error) {
	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	query = query.Order("published_at DESC")
	var recs []model.FederatedPostCache
	if err := query.Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]federation.FederatedPostCache, len(recs))
	for i, rec := range recs {
		result[i] = mapFederatedPostCacheToDomain(rec)
	}
	return result, nil
}

// FederatedCitationRepository stores citation workflows.
type FederatedCitationRepository struct {
	db *gorm.DB
}

func NewFederatedCitationRepository(db *gorm.DB) *FederatedCitationRepository {
	return &FederatedCitationRepository{db: db}
}

func (r *FederatedCitationRepository) Create(ctx context.Context, citation *federation.FederatedCitation) error {
	rec := mapFederatedCitationToModel(citation)
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	citation.ID = rec.ID
	citation.RequestedAt = rec.RequestedAt
	return nil
}

func (r *FederatedCitationRepository) UpdateStatus(ctx context.Context, id int64, status string, reason *string) error {
	updates := map[string]any{
		"status": status,
	}
	if status == "approved" {
		updates["approved_at"] = time.Now().UTC()
	}
	if status == "rejected" {
		updates["rejected_at"] = time.Now().UTC()
		updates["reject_reason"] = reason
	}
	return r.db.WithContext(ctx).Model(&model.FederatedCitation{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *FederatedCitationRepository) ListByTarget(ctx context.Context, articleID int64, status string) ([]federation.FederatedCitation, error) {
	query := r.db.WithContext(ctx).Where("target_article_id = ?", articleID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query = query.Order("requested_at DESC")
	var recs []model.FederatedCitation
	if err := query.Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]federation.FederatedCitation, len(recs))
	for i, rec := range recs {
		result[i] = mapFederatedCitationToDomain(rec)
	}
	return result, nil
}

// FederatedMentionRepository stores mentions delivered to local users.
type FederatedMentionRepository struct {
	db *gorm.DB
}

func NewFederatedMentionRepository(db *gorm.DB) *FederatedMentionRepository {
	return &FederatedMentionRepository{db: db}
}

func (r *FederatedMentionRepository) Create(ctx context.Context, mention *federation.FederatedMention) error {
	rec := mapFederatedMentionToModel(mention)
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	mention.ID = rec.ID
	mention.CreatedAt = rec.CreatedAt
	return nil
}

func (r *FederatedMentionRepository) MarkRead(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.FederatedMention{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_read": true,
			"read_at": time.Now().UTC(),
		}).Error
}

func (r *FederatedMentionRepository) ListByUser(ctx context.Context, userID int64, unreadOnly bool) ([]federation.FederatedMention, error) {
	query := r.db.WithContext(ctx).Where("mentioned_user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	query = query.Order("created_at DESC")
	var recs []model.FederatedMention
	if err := query.Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]federation.FederatedMention, len(recs))
	for i, rec := range recs {
		result[i] = mapFederatedMentionToDomain(rec)
	}
	return result, nil
}

func mapFederationInstanceToDomain(rec model.FederationInstance) federation.FederationInstance {
	return federation.FederationInstance{
		ID:              rec.ID,
		BaseURL:         rec.BaseURL,
		Name:            rec.Name,
		Description:     rec.Description,
		ProtocolVersion: rec.ProtocolVersion,
		PublicKey:       rec.PublicKey,
		KeyID:           rec.KeyID,
		Features:        json.RawMessage(rec.Features),
		Policies:        json.RawMessage(rec.Policies),
		Endpoints:       json.RawMessage(rec.Endpoints),
		Status:          rec.Status,
		LastSeenAt:      rec.LastSeenAt,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       rec.UpdatedAt,
	}
}

func mapFederationInstanceToModel(instance *federation.FederationInstance) model.FederationInstance {
	return model.FederationInstance{
		ID:              instance.ID,
		BaseURL:         instance.BaseURL,
		Name:            instance.Name,
		Description:     instance.Description,
		ProtocolVersion: instance.ProtocolVersion,
		PublicKey:       instance.PublicKey,
		KeyID:           instance.KeyID,
		Features:        datatypes.JSON(instance.Features),
		Policies:        datatypes.JSON(instance.Policies),
		Endpoints:       datatypes.JSON(instance.Endpoints),
		Status:          instance.Status,
		LastSeenAt:      instance.LastSeenAt,
	}
}

func mapFederatedPostCacheToDomain(rec model.FederatedPostCache) federation.FederatedPostCache {
	return federation.FederatedPostCache{
		ID:             rec.ID,
		InstanceID:     rec.InstanceID,
		RemotePostID:   rec.RemotePostID,
		URL:            rec.URL,
		Title:          rec.Title,
		Summary:        rec.Summary,
		ContentPreview: rec.ContentPreview,
		Author:         json.RawMessage(rec.Author),
		Tags:           json.RawMessage(rec.Tags),
		Categories:     json.RawMessage(rec.Categories),
		PublishedAt:    rec.PublishedAt,
		UpdatedAt:      rec.UpdatedAt,
		CoverImage:     rec.CoverImage,
		Language:       rec.Language,
		AllowCitation:  rec.AllowCitation,
		AllowComment:   rec.AllowComment,
		ETag:           rec.ETag,
		LastModified:   rec.LastModified,
		CachedAt:       rec.CachedAt,
	}
}

func mapFederatedPostCacheToModel(post *federation.FederatedPostCache) model.FederatedPostCache {
	return model.FederatedPostCache{
		ID:             post.ID,
		InstanceID:     post.InstanceID,
		RemotePostID:   post.RemotePostID,
		URL:            post.URL,
		Title:          post.Title,
		Summary:        post.Summary,
		ContentPreview: post.ContentPreview,
		Author:         datatypes.JSON(post.Author),
		Tags:           datatypes.JSON(post.Tags),
		Categories:     datatypes.JSON(post.Categories),
		PublishedAt:    post.PublishedAt,
		UpdatedAt:      post.UpdatedAt,
		CoverImage:     post.CoverImage,
		Language:       post.Language,
		AllowCitation:  post.AllowCitation,
		AllowComment:   post.AllowComment,
		ETag:           post.ETag,
		LastModified:   post.LastModified,
		CachedAt:       post.CachedAt,
	}
}

func mapFederatedCitationToDomain(rec model.FederatedCitation) federation.FederatedCitation {
	return federation.FederatedCitation{
		ID:               rec.ID,
		SourceInstanceID: rec.SourceInstanceID,
		SourcePostURL:    rec.SourcePostURL,
		SourcePostTitle:  rec.SourcePostTitle,
		TargetArticleID:  rec.TargetArticleID,
		CitationContext:  rec.CitationContext,
		CitationType:     rec.CitationType,
		Status:           rec.Status,
		RequestedAt:      rec.RequestedAt,
		ApprovedAt:       rec.ApprovedAt,
		RejectedAt:       rec.RejectedAt,
		RejectReason:     rec.RejectReason,
	}
}

func mapFederatedCitationToModel(citation *federation.FederatedCitation) model.FederatedCitation {
	return model.FederatedCitation{
		ID:               citation.ID,
		SourceInstanceID: citation.SourceInstanceID,
		SourcePostURL:    citation.SourcePostURL,
		SourcePostTitle:  citation.SourcePostTitle,
		TargetArticleID:  citation.TargetArticleID,
		CitationContext:  citation.CitationContext,
		CitationType:     citation.CitationType,
		Status:           citation.Status,
		RequestedAt:      citation.RequestedAt,
		ApprovedAt:       citation.ApprovedAt,
		RejectedAt:       citation.RejectedAt,
		RejectReason:     citation.RejectReason,
	}
}

func mapFederatedMentionToDomain(rec model.FederatedMention) federation.FederatedMention {
	return federation.FederatedMention{
		ID:               rec.ID,
		SourceInstanceID: rec.SourceInstanceID,
		SourcePostURL:    rec.SourcePostURL,
		SourcePostTitle:  rec.SourcePostTitle,
		MentionedUserID:  rec.MentionedUserID,
		MentionContext:   rec.MentionContext,
		MentionType:      rec.MentionType,
		IsRead:           rec.IsRead,
		CreatedAt:        rec.CreatedAt,
		ReadAt:           rec.ReadAt,
	}
}

func mapFederatedMentionToModel(mention *federation.FederatedMention) model.FederatedMention {
	return model.FederatedMention{
		ID:               mention.ID,
		SourceInstanceID: mention.SourceInstanceID,
		SourcePostURL:    mention.SourcePostURL,
		SourcePostTitle:  mention.SourcePostTitle,
		MentionedUserID:  mention.MentionedUserID,
		MentionContext:   mention.MentionContext,
		MentionType:      mention.MentionType,
		IsRead:           mention.IsRead,
		CreatedAt:        mention.CreatedAt,
		ReadAt:           mention.ReadAt,
	}
}
