package federation

import (
	"encoding/json"
	"time"
)

// FederationConfig stores local instance federation settings.
type FederationConfig struct {
	ID              int64
	Enabled         bool
	InstanceName    *string
	InstanceURL     *string
	PublicKey       *string
	PrivateKey      *string
	SignatureAlg    string
	RequireHTTPS    bool
	AllowInbound    bool
	AllowOutbound   bool
	DefaultPolicies json.RawMessage
	RateLimits      json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FederationInstance represents a remote instance discovered via well-known.
type FederationInstance struct {
	ID              int64
	BaseURL         string
	Name            *string
	Description     *string
	ProtocolVersion *string
	PublicKey       *string
	KeyID           *string
	Features        json.RawMessage
	Policies        json.RawMessage
	Endpoints       json.RawMessage
	Status          string
	LastSeenAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FederatedPostCache stores cached remote posts for timeline/recommendations.
type FederatedPostCache struct {
	ID             int64
	InstanceID     int64
	RemotePostID   *string
	URL            string
	Title          string
	Summary        string
	ContentPreview *string
	Author         json.RawMessage
	Tags           json.RawMessage
	Categories     json.RawMessage
	PublishedAt    time.Time
	UpdatedAt      *time.Time
	CoverImage     *string
	Language       *string
	AllowCitation  bool
	AllowComment   bool
	ETag           *string
	LastModified   *string
	CachedAt       time.Time
}

// FederatedCitation tracks cross-site citation requests.
type FederatedCitation struct {
	ID               int64
	SourceInstanceID int64
	SourcePostURL    string
	SourcePostTitle  *string
	TargetArticleID  int64
	CitationContext  *string
	CitationType     string
	Status           string
	RequestedAt      time.Time
	ApprovedAt       *time.Time
	RejectedAt       *time.Time
	RejectReason     *string
}

// FederatedMention stores cross-site mentions delivered to local users.
type FederatedMention struct {
	ID               int64
	SourceInstanceID int64
	SourcePostURL    string
	SourcePostTitle  *string
	MentionedUserID  int64
	MentionContext   string
	MentionType      string
	IsRead           bool
	CreatedAt        time.Time
	ReadAt           *time.Time
}
