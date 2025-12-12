package identity

import "context"

// Repository 定义用户及其权限相关的持久化操作。
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByCredential(ctx context.Context, credential string) (*User, error)
	GetRoles(ctx context.Context, userID int64) ([]string, error)
	GetPermissions(ctx context.Context, userID int64) ([]string, error)
	AssignRoles(ctx context.Context, userID int64, roles []string) error
}
