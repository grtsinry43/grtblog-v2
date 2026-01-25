package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	appfed "github.com/grtsinry43/grtblog-v2/server/internal/app/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

type FederationAdminHandler struct {
	cfgSvc      *federationconfig.Service
	contentRepo content.Repository
	outbound    *appfed.OutboundService
	resolver    *fedinfra.Resolver
}

func NewFederationAdminHandler(cfgSvc *federationconfig.Service, contentRepo content.Repository, outbound *appfed.OutboundService, resolver *fedinfra.Resolver) *FederationAdminHandler {
	return &FederationAdminHandler{
		cfgSvc:      cfgSvc,
		contentRepo: contentRepo,
		outbound:    outbound,
		resolver:    resolver,
	}
}

// RequestFriendLink 由后台发起对外友链申请。
// @Summary 后台发起友链申请
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param request body contract.FederationAdminFriendLinkRequestReq true "友链申请参数"
// @Success 200 {object} contract.FederationAdminProxyResp
// @Security BearerAuth
// @Router /admin/federation/friendlinks/request [post]
// @Security JWTAuth
func (h *FederationAdminHandler) RequestFriendLink(c *fiber.Ctx) error {
	var req contract.FederationAdminFriendLinkRequestReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	target := strings.TrimSpace(req.TargetURL)
	if target == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_url 不能为空")
	}
	if h.outbound == nil {
		return response.NewBizErrorWithMsg(response.ServerError, "联邦服务未初始化")
	}
	resp, raw, err := h.outbound.SendFriendLinkRequest(c.Context(), target, req.Message, req.RSSURL)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "请求失败", err)
	}
	return response.Success(c, contract.FederationAdminProxyResp{
		StatusCode: resp.StatusCode,
		Body:       string(raw),
	})
}

// SendCitation 由后台发起对外引用请求。
// @Summary 后台发起引用请求
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param request body contract.FederationAdminCitationReq true "引用请求参数"
// @Success 200 {object} contract.FederationAdminProxyResp
// @Security BearerAuth
// @Router /admin/federation/citations/request [post]
// @Security JWTAuth
func (h *FederationAdminHandler) SendCitation(c *fiber.Ctx) error {
	var req contract.FederationAdminCitationReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	target := strings.TrimSpace(req.TargetInstanceURL)
	if target == "" || strings.TrimSpace(req.TargetPostID) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_instance_url/target_post_id 不能为空")
	}
	if h.outbound == nil {
		return response.NewBizErrorWithMsg(response.ServerError, "联邦服务未初始化")
	}
	article, err := h.resolveArticle(c, req.SourceArticleID, req.SourceShortURL)
	if err != nil {
		if errors.Is(err, content.ErrArticleNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return response.NewBizErrorWithCause(response.ServerError, "文章获取失败", err)
	}
	context := strings.TrimSpace(req.CitationContext)
	if context == "" {
		context = article.Summary
	}
	citationType := strings.TrimSpace(req.CitationType)
	ev := appfed.CitationDetected{
		ArticleID:      article.ID,
		AuthorID:       article.AuthorID,
		Title:          article.Title,
		ShortURL:       article.ShortURL,
		TargetInstance: target,
		TargetPostID:   strings.TrimSpace(req.TargetPostID),
		Context:        context,
		CitationType:   citationType,
	}
	resp, raw, err := h.outbound.SendCitation(c.Context(), ev)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "请求失败", err)
	}
	return response.Success(c, contract.FederationAdminProxyResp{
		StatusCode: resp.StatusCode,
		Body:       string(raw),
	})
}

// SendMention 由后台发起对外提及通知。
// @Summary 后台发起提及通知
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param request body contract.FederationAdminMentionReq true "提及通知参数"
// @Success 200 {object} contract.FederationAdminProxyResp
// @Security BearerAuth
// @Router /admin/federation/mentions/notify [post]
// @Security JWTAuth
func (h *FederationAdminHandler) SendMention(c *fiber.Ctx) error {
	var req contract.FederationAdminMentionReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	target := strings.TrimSpace(req.TargetInstanceURL)
	if target == "" || strings.TrimSpace(req.MentionedUser) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_instance_url/mentioned_user 不能为空")
	}
	if h.outbound == nil {
		return response.NewBizErrorWithMsg(response.ServerError, "联邦服务未初始化")
	}
	article, err := h.resolveArticle(c, req.SourceArticleID, req.SourceShortURL)
	if err != nil {
		if errors.Is(err, content.ErrArticleNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return response.NewBizErrorWithCause(response.ServerError, "文章获取失败", err)
	}
	context := strings.TrimSpace(req.MentionContext)
	if context == "" {
		context = article.Summary
	}
	mentionType := strings.TrimSpace(req.MentionType)
	ev := appfed.MentionDetected{
		ArticleID:      article.ID,
		AuthorID:       article.AuthorID,
		Title:          article.Title,
		ShortURL:       article.ShortURL,
		TargetUser:     strings.TrimSpace(req.MentionedUser),
		TargetInstance: target,
		Context:        context,
		MentionType:    mentionType,
	}
	resp, raw, err := h.outbound.SendMention(c.Context(), ev)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "请求失败", err)
	}
	return response.Success(c, contract.FederationAdminProxyResp{
		StatusCode: resp.StatusCode,
		Body:       string(raw),
	})
}

// CheckRemote 校验远端连通性（manifest/public-key/endpoints）。
// @Summary 远端联通性检查
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param target_url query string true "远端实例地址"
// @Success 200 {object} contract.FederationAdminRemoteCheckResp
// @Security BearerAuth
// @Router /admin/federation/remote/check [get]
// @Security JWTAuth
func (h *FederationAdminHandler) CheckRemote(c *fiber.Ctx) error {
	target := strings.TrimSpace(c.Query("target_url"))
	if target == "" {
		var req contract.FederationAdminRemoteCheckReq
		if err := c.BodyParser(&req); err == nil {
			target = strings.TrimSpace(req.TargetURL)
		}
	}
	if target == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_url 不能为空")
	}
	if h.resolver == nil {
		return response.NewBizErrorWithMsg(response.ServerError, "resolver 未初始化")
	}
	baseURL := normalizeInstanceURL(target)
	if baseURL == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_url 不能为空")
	}
	manifest, err := h.resolver.FetchManifest(c.Context(), baseURL)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "拉取 manifest 失败", err)
	}
	publicKey, err := h.resolver.FetchPublicKey(c.Context(), baseURL)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "拉取公钥失败", err)
	}
	endpoints, err := h.resolver.FetchEndpoints(c.Context(), baseURL)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "拉取 endpoints 失败", err)
	}
	return response.Success(c, contract.FederationAdminRemoteCheckResp{
		Manifest:  mapManifestResp(manifest),
		PublicKey: mapPublicKeyResp(publicKey),
		Endpoints: mapEndpointsResp(endpoints),
	})
}

func (h *FederationAdminHandler) resolveArticle(c *fiber.Ctx, id *int64, shortURL *string) (*content.Article, error) {
	if id != nil && *id > 0 {
		return h.contentRepo.GetArticleByID(c.Context(), *id)
	}
	if shortURL != nil && strings.TrimSpace(*shortURL) != "" {
		return h.contentRepo.GetArticleByShortURL(c.Context(), strings.TrimSpace(*shortURL))
	}
	return nil, content.ErrArticleNotFound
}

func mapManifestResp(manifest *fedinfra.Manifest) map[string]any {
	if manifest == nil {
		return nil
	}
	return map[string]any{
		"protocol_version": manifest.ProtocolVersion,
		"instance": map[string]any{
			"name":        manifest.Instance.Name,
			"url":         manifest.Instance.URL,
			"description": manifest.Instance.Description,
			"language":    manifest.Instance.Language,
			"timezone":    manifest.Instance.Timezone,
		},
		"software": map[string]any{
			"name":    manifest.Software.Name,
			"version": manifest.Software.Version,
		},
		"features": manifest.Features,
		"policies": map[string]any{
			"allow_citation":                   manifest.Policies.AllowCitation,
			"allow_mention":                    manifest.Policies.AllowMention,
			"auto_approve_friendlink_citation": manifest.Policies.AutoApproveFriendlinkCitation,
			"require_https":                    manifest.Policies.RequireHTTPS,
			"max_cache_age":                    manifest.Policies.MaxCacheAge,
		},
		"created_at": manifest.CreatedAt.Format(time.RFC3339),
		"updated_at": manifest.UpdatedAt.Format(time.RFC3339),
	}
}

func mapPublicKeyResp(doc *fedinfra.PublicKeyDoc) map[string]any {
	if doc == nil {
		return nil
	}
	return map[string]any{
		"key_id":     doc.KeyID,
		"algorithm":  doc.Algorithm,
		"public_key": doc.PublicKey,
	}
}

func mapEndpointsResp(doc *fedinfra.EndpointsDoc) map[string]any {
	if doc == nil {
		return nil
	}
	return map[string]any{
		"base_url":  doc.BaseURL,
		"endpoints": doc.Endpoints,
	}
}
