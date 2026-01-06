package identity

import "context"

// Repository 定义用户及其权限相关的持久化操作。
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByCredential(ctx context.Context, credential string) (*User, error)
	UpdateProfile(ctx context.Context, userID int64, nickname, avatar, email string) (*User, error)
	UpdatePassword(ctx context.Context, userID int64, hashed string) error
	ListOAuthBindings(ctx context.Context, userID int64) ([]UserOAuthBinding, error)
	FindByOAuth(ctx context.Context, providerKey, oauthID string) (*User, error)
	BindOAuth(ctx context.Context, link UserOAuth) error
	CountUsers(ctx context.Context) (int64, error)
}

// OAuthProviderRepository 提供 OAuth/OIDC 提供方配置。
type OAuthProviderRepository interface {
	ListEnabled(ctx context.Context) ([]OAuthProvider, error)
	GetByKey(ctx context.Context, key string) (*OAuthProvider, error)
	ListAll(ctx context.Context) ([]OAuthProvider, error)
	Create(ctx context.Context, provider *OAuthProvider) error
	Update(ctx context.Context, provider *OAuthProvider) error
	Delete(ctx context.Context, key string) error
}
