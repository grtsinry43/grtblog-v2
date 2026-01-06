package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Username  string         `gorm:"column:username;size:45;not null"`
	Nickname  string         `gorm:"column:nickname;size:45;not null"`
	Email     string         `gorm:"column:email;size:255"`
	Password  string         `gorm:"column:password;size:60"`
	Avatar    string         `gorm:"column:avatar;size:255"`
	IsActive  bool           `gorm:"column:is_active;default:true"`
	IsAdmin   bool           `gorm:"column:is_admin;default:false"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string { return "app_user" }

type AdminToken struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	Token       string    `gorm:"column:token;size:255;not null"`
	UserID      int64     `gorm:"column:user_id;not null"`
	Description string    `gorm:"column:description"`
	ExpireAt    time.Time `gorm:"column:expire_at;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AdminToken) TableName() string { return "admin_token" }

type OAuthProvider struct {
	ID                    int64     `gorm:"column:id;primaryKey"`
	ProviderKey           string    `gorm:"column:provider_key;size:64;not null"`
	DisplayName           string    `gorm:"column:display_name;size:128;not null"`
	ClientID              string    `gorm:"column:client_id;size:255;not null"`
	ClientSecret          string    `gorm:"column:client_secret;size:255;not null"`
	AuthorizationEndpoint string    `gorm:"column:authorization_endpoint;size:512;not null"`
	TokenEndpoint         string    `gorm:"column:token_endpoint;size:512;not null"`
	UserinfoEndpoint      string    `gorm:"column:userinfo_endpoint;size:512"`
	RedirectURITemplate   string    `gorm:"column:redirect_uri_template;size:512;not null"`
	Scopes                string    `gorm:"column:scopes;size:512"`
	Issuer                string    `gorm:"column:issuer;size:512"`
	JWKSURI               string    `gorm:"column:jwks_uri;size:512"`
	PKCERequired          bool      `gorm:"column:pkce_required"`
	Enabled               bool      `gorm:"column:enabled"`
	ExtraParams           []byte    `gorm:"column:extra_params;type:jsonb"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (OAuthProvider) TableName() string { return "oauth_provider" }

type UserOAuth struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	UserID       int64      `gorm:"column:user_id;not null"`
	ProviderKey  string     `gorm:"column:provider_key;size:64;not null"`
	OAuthID      string     `gorm:"column:oauth_id;size:255;not null"`
	AccessToken  string     `gorm:"column:access_token;size:512"`
	RefreshToken string     `gorm:"column:refresh_token;size:512"`
	ExpiresAt    *time.Time `gorm:"column:expires_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserOAuth) TableName() string { return "user_oauth" }
