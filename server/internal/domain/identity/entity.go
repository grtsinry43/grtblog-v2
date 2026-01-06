package identity

import "time"

// User aggregates登录/权限相关信息。
type User struct {
	ID        int64
	Username  string
	Nickname  string
	Email     string
	Password  string
	Avatar    string
	IsActive  bool
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Role struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Permission struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole struct {
	UserID int64
	RoleID int64
}

type RolePermission struct {
	RoleID       int64
	PermissionID int64
}

type AdminToken struct {
	ID          int64
	Token       string
	UserID      int64
	Description string
	ExpireAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OAuthProvider struct {
	ID                    int64
	ProviderKey           string
	DisplayName           string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserinfoEndpoint      string
	RedirectURITemplate   string
	Scopes                string
	Issuer                string
	JWKSURI               string
	PKCERequired          bool
	Enabled               bool
	ExtraParams           map[string]any
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type UserOAuth struct {
	ID           int64
	UserID       int64
	ProviderKey  string
	OAuthID      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserOAuthBinding 描述用户已绑定的外部身份信息（用于展示）。
type UserOAuthBinding struct {
	ProviderKey   string
	ProviderName  string
	OAuthID       string
	BoundAt       time.Time
	ExpiresAt     *time.Time
	ProviderScope string
}
