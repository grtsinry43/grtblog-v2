package federation

import "time"

// Manifest mirrors .well-known/blog-federation/manifest.json.
type Manifest struct {
	ProtocolVersion string           `json:"protocol_version"`
	Instance        ManifestNode     `json:"instance"`
	Software        ManifestSoftware `json:"software"`
	Admin           *ManifestAdmin   `json:"admin,omitempty"`
	Features        []string         `json:"features"`
	Policies        ManifestPolicy   `json:"policies"`
	RateLimits      ManifestRate     `json:"rate_limits"`
	RSSFeeds        []RSSFeed        `json:"rss_feeds,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type ManifestNode struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

type ManifestSoftware struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository string `json:"repository,omitempty"`
}

type ManifestAdmin struct {
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
}

type ManifestPolicy struct {
	AllowCitation                 bool  `json:"allow_citation"`
	AllowMention                  bool  `json:"allow_mention"`
	AutoApproveFriendlinkCitation bool  `json:"auto_approve_friendlink_citation"`
	RequireHTTPS                  bool  `json:"require_https"`
	MaxCacheAge                   int64 `json:"max_cache_age"`
}

type ManifestRate struct {
	TimelineSync    int64 `json:"timeline_sync"`
	CitationRequest int64 `json:"citation_request"`
	MentionNotify   int64 `json:"mention_notify"`
}

type RSSFeed struct {
	URL   string `json:"url"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

// PublicKeyDoc mirrors .well-known/blog-federation/public-key.json.
type PublicKeyDoc struct {
	KeyID     string     `json:"key_id"`
	Algorithm string     `json:"algorithm"`
	PublicKey string     `json:"public_key"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// EndpointsDoc mirrors .well-known/blog-federation/endpoints.json.
type EndpointsDoc struct {
	BaseURL   string            `json:"base_url"`
	Endpoints map[string]string `json:"endpoints"`
}

// VerifiedSignature carries resolved signer context.
type VerifiedSignature struct {
	KeyID    string
	BaseURL  string
	DateTime time.Time
}
