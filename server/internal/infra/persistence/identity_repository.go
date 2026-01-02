package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type IdentityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

func (r *IdentityRepository) FindByID(ctx context.Context, id int64) (*identity.User, error) {
	var rec model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrUserNotFound
		}
		return nil, err
	}
	user := mapUserToDomain(rec)
	return &user, nil
}

func (r *IdentityRepository) Create(ctx context.Context, user *identity.User) error {
	rec := model.User{
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   user.Avatar,
		IsActive: user.IsActive,
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		if isUniqueConstraint(err) {
			return identity.ErrUserExists
		}
		return err
	}
	user.ID = rec.ID
	user.CreatedAt = rec.CreatedAt
	user.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *IdentityRepository) FindByUsername(ctx context.Context, username string) (*identity.User, error) {
	var rec model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrUserNotFound
		}
		return nil, err
	}
	user := mapUserToDomain(rec)
	return &user, nil
}

func (r *IdentityRepository) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	var rec model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrUserNotFound
		}
		return nil, err
	}
	user := mapUserToDomain(rec)
	return &user, nil
}

func (r *IdentityRepository) FindByCredential(ctx context.Context, credential string) (*identity.User, error) {
	var rec model.User
	if err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", credential, credential).
		First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrInvalidCredentials
		}
		return nil, err
	}
	user := mapUserToDomain(rec)
	return &user, nil
}

func (r *IdentityRepository) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	var roles []string
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Select("role.role_name").
		Joins("JOIN role ON role.id = user_role.role_id").
		Where("user_role.user_id = ?", userID).
		Scan(&roles).Error
	return roles, err
}

func (r *IdentityRepository) GetPermissions(ctx context.Context, userID int64) ([]string, error) {
	var perms []string
	err := r.db.WithContext(ctx).
		Table("role_permission").
		Select("DISTINCT permission.permission_name").
		Joins("JOIN permission ON permission.id = role_permission.permission_id").
		Joins("JOIN user_role ON user_role.role_id = role_permission.role_id").
		Where("user_role.user_id = ?", userID).
		Scan(&perms).Error
	return perms, err
}

func (r *IdentityRepository) AssignRoles(ctx context.Context, userID int64, roles []string) error {
	if len(roles) == 0 {
		return nil
	}
	var roleRecords []model.Role
	if err := r.db.WithContext(ctx).
		Where("role_name IN ?", roles).
		Find(&roleRecords).Error; err != nil {
		return err
	}
	if len(roleRecords) != len(roles) {
		return identity.ErrRoleNotFound
	}
	var mappings []model.UserRole
	for _, role := range roleRecords {
		mappings = append(mappings, model.UserRole{
			UserID: userID,
			RoleID: role.ID,
		})
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
		DoNothing: true,
	}).Create(&mappings).Error
}

func (r *IdentityRepository) FindByOAuth(ctx context.Context, providerKey, oauthID string) (*identity.User, error) {
	var rec model.User
	err := r.db.WithContext(ctx).
		Table("user_oauth").
		Select("app_user.*").
		Joins("JOIN app_user ON app_user.id = user_oauth.user_id").
		Where("user_oauth.provider_key = ? AND user_oauth.oauth_id = ?", providerKey, oauthID).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity.ErrUserNotFound
		}
		return nil, err
	}
	user := mapUserToDomain(rec)
	return &user, nil
}

func (r *IdentityRepository) BindOAuth(ctx context.Context, link identity.UserOAuth) error {
	rec := model.UserOAuth{
		UserID:       link.UserID,
		ProviderKey:  link.ProviderKey,
		OAuthID:      link.OAuthID,
		AccessToken:  link.AccessToken,
		RefreshToken: link.RefreshToken,
		ExpiresAt:    link.ExpiresAt,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_key"}, {Name: "oauth_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "access_token", "refresh_token", "expires_at", "updated_at"}),
	}).Create(&rec).Error
}

func (r *IdentityRepository) UpdateProfile(ctx context.Context, userID int64, nickname, avatar, email string) (*identity.User, error) {
	updates := map[string]any{
		"updated_at": time.Now(),
	}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	if email != "" {
		updates["email"] = email
	}
	if len(updates) == 1 { // only updated_at
		return r.FindByID(ctx, userID)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		if isUniqueConstraint(err) {
			return nil, identity.ErrUserExists
		}
		return nil, err
	}
	return r.FindByID(ctx, userID)
}

func (r *IdentityRepository) UpdatePassword(ctx context.Context, userID int64, hashed string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"password":   hashed,
			"updated_at": time.Now(),
		}).Error
}

func (r *IdentityRepository) ListOAuthBindings(ctx context.Context, userID int64) ([]identity.UserOAuthBinding, error) {
	type bindingRow struct {
		ProviderKey string
		DisplayName string
		OAuthID     string
		Scopes      string
		CreatedAt   time.Time
		ExpiresAt   *time.Time
	}
	var rows []bindingRow
	err := r.db.WithContext(ctx).
		Table("user_oauth").
		Select("user_oauth.provider_key, user_oauth.oauth_id, user_oauth.created_at, user_oauth.expires_at, oauth_provider.display_name, oauth_provider.scopes").
		Joins("JOIN oauth_provider ON oauth_provider.provider_key = user_oauth.provider_key").
		Where("user_oauth.user_id = ?", userID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]identity.UserOAuthBinding, 0, len(rows))
	for _, row := range rows {
		result = append(result, identity.UserOAuthBinding{
			ProviderKey:   row.ProviderKey,
			ProviderName:  row.DisplayName,
			OAuthID:       row.OAuthID,
			BoundAt:       row.CreatedAt,
			ExpiresAt:     row.ExpiresAt,
			ProviderScope: row.Scopes,
		})
	}
	return result, nil
}

func mapUserToDomain(rec model.User) identity.User {
	var deleted *time.Time
	if rec.DeletedAt.Valid {
		deleted = &rec.DeletedAt.Time
	}
	return identity.User{
		ID:        rec.ID,
		Username:  rec.Username,
		Nickname:  rec.Nickname,
		Email:     rec.Email,
		Password:  rec.Password,
		Avatar:    rec.Avatar,
		IsActive:  rec.IsActive,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		DeletedAt: deleted,
	}
}

func isUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "uq_app_user_username") || strings.Contains(err.Error(), "uq_app_user_email")
}
