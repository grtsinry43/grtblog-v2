package federation

import (
	"context"
	"time"
)

// FederationInstanceRepository manages remote instance records.
type FederationInstanceRepository interface {
	GetByBaseURL(ctx context.Context, baseURL string) (*FederationInstance, error)
	Create(ctx context.Context, instance *FederationInstance) error
	Update(ctx context.Context, instance *FederationInstance) error
	ListActive(ctx context.Context) ([]FederationInstance, error)
}

// FederatedPostCacheRepository stores cached timeline posts.
type FederatedPostCacheRepository interface {
	UpsertBatch(ctx context.Context, posts []FederatedPostCache) error
	ListByInstance(ctx context.Context, instanceID int64, since *time.Time, limit int) ([]FederatedPostCache, error)
	ListRecent(ctx context.Context, limit int) ([]FederatedPostCache, error)
}

// FederatedCitationRepository stores citation workflows.
type FederatedCitationRepository interface {
	Create(ctx context.Context, citation *FederatedCitation) error
	UpdateStatus(ctx context.Context, id int64, status string, reason *string) error
	ListByTarget(ctx context.Context, articleID int64, status string) ([]FederatedCitation, error)
}

// FederatedMentionRepository stores mentions delivered to local users.
type FederatedMentionRepository interface {
	Create(ctx context.Context, mention *FederatedMention) error
	MarkRead(ctx context.Context, id int64) error
	ListByUser(ctx context.Context, userID int64, unreadOnly bool) ([]FederatedMention, error)
}
