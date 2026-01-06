package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"gorm.io/gorm"
)

// AdminOAuthHandler 提供 OAuth Provider 管理接口（仅后台）。
type AdminOAuthHandler struct {
	repo *persistence.OAuthProviderRepository
}

func NewAdminOAuthHandler(repo *persistence.OAuthProviderRepository) *AdminOAuthHandler {
	return &AdminOAuthHandler{repo: repo}
}

type OAuthProviderPayload struct {
	Key                   string         `json:"key"`
	DisplayName           string         `json:"displayName"`
	ClientID              string         `json:"clientId"`
	ClientSecret          string         `json:"clientSecret"`
	AuthorizationEndpoint string         `json:"authorizationEndpoint"`
	TokenEndpoint         string         `json:"tokenEndpoint"`
	UserinfoEndpoint      string         `json:"userinfoEndpoint"`
	RedirectURITemplate   string         `json:"redirectUriTemplate"`
	Scopes                string         `json:"scopes"`
	Issuer                string         `json:"issuer"`
	JWKSURI               string         `json:"jwksUri"`
	PKCERequired          bool           `json:"pkceRequired"`
	Enabled               bool           `json:"enabled"`
	ExtraParams           map[string]any `json:"extraParams"`
}

func payloadToDomain(p OAuthProviderPayload) identity.OAuthProvider {
	return identity.OAuthProvider{
		ProviderKey:           p.Key,
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
		ExtraParams:           p.ExtraParams,
	}
}

// Swagger envelopes
type OAuthProviderListEnvelope struct {
	Code   int                      `json:"code"`
	BizErr string                   `json:"bizErr"`
	Msg    string                   `json:"msg"`
	Data   []identity.OAuthProvider `json:"data"`
	Meta   response.Meta            `json:"meta"`
}

type OAuthProviderEnvelope struct {
	Code   int                    `json:"code"`
	BizErr string                 `json:"bizErr"`
	Msg    string                 `json:"msg"`
	Data   identity.OAuthProvider `json:"data"`
	Meta   response.Meta          `json:"meta"`
}

type AdminOAuthProviderResp struct {
	Key                   string         `json:"key"`
	DisplayName           string         `json:"displayName"`
	ClientID              string         `json:"clientId"`
	AuthorizationEndpoint string         `json:"authorizationEndpoint"`
	TokenEndpoint         string         `json:"tokenEndpoint"`
	UserinfoEndpoint      string         `json:"userinfoEndpoint"`
	RedirectURITemplate   string         `json:"redirectUriTemplate"`
	Scopes                string         `json:"scopes"`
	Issuer                string         `json:"issuer"`
	JWKSURI               string         `json:"jwksUri"`
	PKCERequired          bool           `json:"pkceRequired"`
	Enabled               bool           `json:"enabled"`
	ExtraParams           map[string]any `json:"extraParams"`
	CreatedAt             string         `json:"createdAt"`
	UpdatedAt             string         `json:"updatedAt"`
}

type AdminOAuthProviderListEnvelope struct {
	Code   int                      `json:"code"`
	BizErr string                   `json:"bizErr"`
	Msg    string                   `json:"msg"`
	Data   []AdminOAuthProviderResp `json:"data"`
	Meta   response.Meta            `json:"meta"`
}

type AdminOAuthProviderEnvelope struct {
	Code   int                    `json:"code"`
	BizErr string                 `json:"bizErr"`
	Msg    string                 `json:"msg"`
	Data   AdminOAuthProviderResp `json:"data"`
	Meta   response.Meta          `json:"meta"`
}

type AdminGenericMessageEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   any           `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// List godoc
// @Summary 列出全部 OAuth Providers
// @Tags Admin-OAuth
// @Produce json
// @Success 200 {object} AdminOAuthProviderListEnvelope
// @Security BearerAuth
// @Router /admin/oauth-providers [get]
func (h *AdminOAuthHandler) List(c *fiber.Ctx) error {
	items, err := h.repo.ListAll(c.Context())
	if err != nil {
		return err
	}
	resp := make([]AdminOAuthProviderResp, 0, len(items))
	for _, item := range items {
		resp = append(resp, sanitizeAdminProvider(item))
	}
	return response.Success(c, resp)
}

// Create godoc
// @Summary 创建 OAuth Provider
// @Tags Admin-OAuth
// @Accept json
// @Produce json
// @Param request body OAuthProviderPayload true "provider 配置"
// @Success 200 {object} AdminOAuthProviderEnvelope
// @Security BearerAuth
// @Router /admin/oauth-providers [post]
func (h *AdminOAuthHandler) Create(c *fiber.Ctx) error {
	var req OAuthProviderPayload
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	p := payloadToDomain(req)
	if p.ProviderKey == "" || p.AuthorizationEndpoint == "" || p.TokenEndpoint == "" || p.RedirectURITemplate == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key/auth/token/redirect 不能为空")
	}
	if err := h.repo.Create(c.Context(), &p); err != nil {
		return err
	}
	Audit(c, "admin.oauth.create", map[string]any{"key": p.ProviderKey})
	return response.SuccessWithMessage(c, sanitizeAdminProvider(p), "created")
}

// Update godoc
// @Summary 更新 OAuth Provider
// @Tags Admin-OAuth
// @Accept json
// @Produce json
// @Param key path string true "provider key"
// @Param request body OAuthProviderPayload true "provider 配置"
// @Success 200 {object} AdminOAuthProviderEnvelope
// @Security BearerAuth
// @Router /admin/oauth-providers/{key} [put]
func (h *AdminOAuthHandler) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	var req OAuthProviderPayload
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	p := payloadToDomain(req)
	p.ProviderKey = key
	if err := h.repo.Update(c.Context(), &p); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return err
	}
	Audit(c, "admin.oauth.update", map[string]any{"key": key})
	return response.SuccessWithMessage(c, sanitizeAdminProvider(p), "updated")
}

// Delete godoc
// @Summary 删除 OAuth Provider
// @Tags Admin-OAuth
// @Produce json
// @Param key path string true "provider key"
// @Success 200 {object} AdminGenericMessageEnvelope
// @Security BearerAuth
// @Router /admin/oauth-providers/{key} [delete]
func (h *AdminOAuthHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	if err := h.repo.Delete(c.Context(), key); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return err
	}
	Audit(c, "admin.oauth.delete", map[string]any{"key": key})
	return response.SuccessWithMessage[any](c, nil, "deleted")
}

func sanitizeAdminProvider(p identity.OAuthProvider) AdminOAuthProviderResp {
	return AdminOAuthProviderResp{
		Key:                   p.ProviderKey,
		DisplayName:           p.DisplayName,
		ClientID:              p.ClientID,
		AuthorizationEndpoint: p.AuthorizationEndpoint,
		TokenEndpoint:         p.TokenEndpoint,
		UserinfoEndpoint:      p.UserinfoEndpoint,
		RedirectURITemplate:   p.RedirectURITemplate,
		Scopes:                p.Scopes,
		Issuer:                p.Issuer,
		JWKSURI:               p.JWKSURI,
		PKCERequired:          p.PKCERequired,
		Enabled:               p.Enabled,
		ExtraParams:           p.ExtraParams,
		CreatedAt:             p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             p.UpdatedAt.Format(time.RFC3339),
	}
}
