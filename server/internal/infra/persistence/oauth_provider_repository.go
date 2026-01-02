package persistence

import (
	"context"
	"encoding/json"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type OAuthProviderRepository struct {
	db *gorm.DB
}

func NewOAuthProviderRepository(db *gorm.DB) *OAuthProviderRepository {
	return &OAuthProviderRepository{db: db}
}

func (r *OAuthProviderRepository) ListEnabled(ctx context.Context) ([]identity.OAuthProvider, error) {
	var recs []model.OAuthProvider
	if err := r.db.WithContext(ctx).Where("enabled = true").Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]identity.OAuthProvider, 0, len(recs))
	for _, rec := range recs {
		result = append(result, mapOAuthProvider(rec))
	}
	return result, nil
}

func (r *OAuthProviderRepository) GetByKey(ctx context.Context, key string) (*identity.OAuthProvider, error) {
	var rec model.OAuthProvider
	if err := r.db.WithContext(ctx).Where("provider_key = ? AND enabled = true", key).First(&rec).Error; err != nil {
		return nil, err
	}
	provider := mapOAuthProvider(rec)
	return &provider, nil
}

func (r *OAuthProviderRepository) ListAll(ctx context.Context) ([]identity.OAuthProvider, error) {
	var recs []model.OAuthProvider
	if err := r.db.WithContext(ctx).Find(&recs).Error; err != nil {
		return nil, err
	}
	result := make([]identity.OAuthProvider, 0, len(recs))
	for _, rec := range recs {
		result = append(result, mapOAuthProvider(rec))
	}
	return result, nil
}

func (r *OAuthProviderRepository) Create(ctx context.Context, provider *identity.OAuthProvider) error {
	rec := toOAuthModel(*provider)
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	provider.ID = rec.ID
	provider.CreatedAt = rec.CreatedAt
	provider.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *OAuthProviderRepository) Update(ctx context.Context, provider *identity.OAuthProvider) error {
	rec := toOAuthModel(*provider)
	result := r.db.WithContext(ctx).
		Model(&model.OAuthProvider{}).
		Where("provider_key = ?", provider.ProviderKey).
		Updates(rec)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OAuthProviderRepository) Delete(ctx context.Context, key string) error {
	result := r.db.WithContext(ctx).Where("provider_key = ?", key).Delete(&model.OAuthProvider{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func mapOAuthProvider(rec model.OAuthProvider) identity.OAuthProvider {
	var extras map[string]any
	if len(rec.ExtraParams) > 0 {
		_ = json.Unmarshal(rec.ExtraParams, &extras)
	}
	return identity.OAuthProvider{
		ID:                    rec.ID,
		ProviderKey:           rec.ProviderKey,
		DisplayName:           rec.DisplayName,
		ClientID:              rec.ClientID,
		ClientSecret:          rec.ClientSecret,
		AuthorizationEndpoint: rec.AuthorizationEndpoint,
		TokenEndpoint:         rec.TokenEndpoint,
		UserinfoEndpoint:      rec.UserinfoEndpoint,
		RedirectURITemplate:   rec.RedirectURITemplate,
		Scopes:                rec.Scopes,
		Issuer:                rec.Issuer,
		JWKSURI:               rec.JWKSURI,
		PKCERequired:          rec.PKCERequired,
		Enabled:               rec.Enabled,
		ExtraParams:           extras,
		CreatedAt:             rec.CreatedAt,
		UpdatedAt:             rec.UpdatedAt,
	}
}

func toOAuthModel(p identity.OAuthProvider) model.OAuthProvider {
	var extras []byte
	if p.ExtraParams != nil {
		extras, _ = json.Marshal(p.ExtraParams)
	}
	return model.OAuthProvider{
		ID:                    p.ID,
		ProviderKey:           p.ProviderKey,
		DisplayName:           p.DisplayName,
		ClientID:              p.ClientID,
		ClientSecret:          p.ClientSecret,
		AuthorizationEndpoint: p.AuthorizationEndpoint,
		TokenEndpoint:         p.TokenEndpoint,
		UserinfoEndpoint:      p.UserinfoEndpoint,
		RedirectURITemplate:   p.RedirectURITemplate,
		Scopes:                p.Scopes,
		Issuer:                p.Issuer,
		JWKSURI:               p.JWKSURI,
		PKCERequired:          p.PKCERequired,
		Enabled:               p.Enabled,
		ExtraParams:           extras,
	}
}
