package model

import (
	"time"

	"gorm.io/datatypes"
)

type FederationInstance struct {
	ID              int64          `gorm:"column:id;primaryKey"`
	BaseURL         string         `gorm:"column:base_url;size:255;not null"`
	Name            *string        `gorm:"column:name;size:255"`
	Description     *string        `gorm:"column:description;type:text"`
	ProtocolVersion *string        `gorm:"column:protocol_version;size:20"`
	PublicKey       *string        `gorm:"column:public_key;type:text"`
	KeyID           *string        `gorm:"column:key_id;type:text"`
	Features        datatypes.JSON `gorm:"column:features;type:jsonb;not null"`
	Policies        datatypes.JSON `gorm:"column:policies;type:jsonb;not null"`
	Endpoints       datatypes.JSON `gorm:"column:endpoints;type:jsonb;not null"`
	Status          string         `gorm:"column:status;size:20;not null"`
	LastSeenAt      *time.Time     `gorm:"column:last_seen_at"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (FederationInstance) TableName() string { return "federation_instance" }

type FederatedPostCache struct {
	ID             int64          `gorm:"column:id;primaryKey"`
	InstanceID     int64          `gorm:"column:instance_id;not null"`
	RemotePostID   *string        `gorm:"column:remote_post_id;size:255"`
	URL            string         `gorm:"column:url;size:500;not null"`
	Title          string         `gorm:"column:title;size:500;not null"`
	Summary        string         `gorm:"column:summary;type:text;not null"`
	ContentPreview *string        `gorm:"column:content_preview;type:text"`
	Author         datatypes.JSON `gorm:"column:author;type:jsonb;not null"`
	Tags           datatypes.JSON `gorm:"column:tags;type:jsonb;not null"`
	Categories     datatypes.JSON `gorm:"column:categories;type:jsonb;not null"`
	PublishedAt    time.Time      `gorm:"column:published_at;not null"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at"`
	CoverImage     *string        `gorm:"column:cover_image;size:500"`
	Language       *string        `gorm:"column:language;size:20"`
	AllowCitation  bool           `gorm:"column:allow_citation;not null"`
	AllowComment   bool           `gorm:"column:allow_comment;not null"`
	ETag           *string        `gorm:"column:etag;size:255"`
	LastModified   *string        `gorm:"column:last_modified;size:255"`
	CachedAt       time.Time      `gorm:"column:cached_at;autoCreateTime"`
}

func (FederatedPostCache) TableName() string { return "federated_post_cache" }

type FederatedCitation struct {
	ID               int64      `gorm:"column:id;primaryKey"`
	SourceInstanceID int64      `gorm:"column:source_instance_id;not null"`
	SourcePostURL    string     `gorm:"column:source_post_url;size:500;not null"`
	SourcePostTitle  *string    `gorm:"column:source_post_title;size:500"`
	TargetArticleID  int64      `gorm:"column:target_article_id;not null"`
	CitationContext  *string    `gorm:"column:citation_context;type:text"`
	CitationType     string     `gorm:"column:citation_type;size:20;not null"`
	Status           string     `gorm:"column:status;size:20;not null"`
	RequestedAt      time.Time  `gorm:"column:requested_at;autoCreateTime"`
	ApprovedAt       *time.Time `gorm:"column:approved_at"`
	RejectedAt       *time.Time `gorm:"column:rejected_at"`
	RejectReason     *string    `gorm:"column:reject_reason;type:text"`
}

func (FederatedCitation) TableName() string { return "federated_citation" }

type FederatedMention struct {
	ID               int64      `gorm:"column:id;primaryKey"`
	SourceInstanceID int64      `gorm:"column:source_instance_id;not null"`
	SourcePostURL    string     `gorm:"column:source_post_url;size:500;not null"`
	SourcePostTitle  *string    `gorm:"column:source_post_title;size:500"`
	MentionedUserID  int64      `gorm:"column:mentioned_user_id;not null"`
	MentionContext   string     `gorm:"column:mention_context;type:text;not null"`
	MentionType      string     `gorm:"column:mention_type;size:20;not null"`
	IsRead           bool       `gorm:"column:is_read;not null"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime"`
	ReadAt           *time.Time `gorm:"column:read_at"`
}

func (FederatedMention) TableName() string { return "federated_mention" }
