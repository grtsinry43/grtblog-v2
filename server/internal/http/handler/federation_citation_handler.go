package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

type FederationCitationHandler struct {
	cfgSvc       *federationconfig.Service
	contentRepo  content.Repository
	instanceRepo federation.FederationInstanceRepository
	citationRepo federation.FederatedCitationRepository
	linkRepo     social.FriendLinkRepository
	resolver     *fedinfra.Resolver
	verifier     *fedinfra.Verifier
}

func NewFederationCitationHandler(
	cfgSvc *federationconfig.Service,
	contentRepo content.Repository,
	instanceRepo federation.FederationInstanceRepository,
	citationRepo federation.FederatedCitationRepository,
	linkRepo social.FriendLinkRepository,
	resolver *fedinfra.Resolver,
	verifier *fedinfra.Verifier,
) *FederationCitationHandler {
	return &FederationCitationHandler{
		cfgSvc:       cfgSvc,
		contentRepo:  contentRepo,
		instanceRepo: instanceRepo,
		citationRepo: citationRepo,
		linkRepo:     linkRepo,
		resolver:     resolver,
		verifier:     verifier,
	}
}

// RequestCitation handles signed citation requests from remote instances.
// @Summary 联合引用申请（入站）
// @Tags Federation
// @Accept json
// @Produce json
// @Param request body contract.FederationCitationRequestReq true "引用申请参数"
// @Success 200 {object} contract.FederationCitationResponseResp
// @Router /api/federation/citations/request [post]
func (h *FederationCitationHandler) RequestCitation(c *fiber.Ctx) error {
	body := c.Body()
	req, err := parseFederationRequest(c)
	if err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求解析失败", err)
	}

	signature, err := h.verifier.VerifyRequest(c.Context(), req, body)
	if err != nil {
		log.Printf("[federation] 入站 引用申请 校验失败 ip=%s err=%v", c.IP(), err)
		return response.NewBizErrorWithMsg(response.Unauthorized, "签名校验失败")
	}

	var payload contract.FederationCitationRequestReq
	if err := json.Unmarshal(body, &payload); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(payload.SourceInstanceURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "source_instance_url 不能为空")
	}
	if strings.TrimSpace(payload.SourcePost.URL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "source_post.url 不能为空")
	}
	if strings.TrimSpace(payload.TargetPostID) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "target_post_id 不能为空")
	}
	if signature != nil && signature.BaseURL != "" && !sameBaseURL(signature.BaseURL, payload.SourceInstanceURL) {
		return response.NewBizErrorWithMsg(response.Unauthorized, "签名来源与请求不一致")
	}

	settings, err := h.cfgSvc.Settings(c.Context())
	if err != nil || !settings.Enabled {
		return response.NewBizErrorWithMsg(response.Unauthorized, "联合未启用")
	}
	policy := parseFederationPolicy(settings)
	if !policyBool(policy.AllowCitation, true) {
		return response.NewBizErrorWithMsg(response.Unauthorized, "未允许引用请求")
	}
	if !settings.AllowInbound {
		return response.NewBizErrorWithMsg(response.Unauthorized, "已关闭入站请求")
	}

	article, err := h.resolveTargetArticle(c.Context(), payload.TargetPostID)
	if err != nil {
		if errors.Is(err, content.ErrArticleNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return response.NewBizErrorWithCause(response.ServerError, "目标文章解析失败", err)
	}
	if !article.IsPublished {
		return response.NewBizError(response.NotFound)
	}

	instance, err := ensureFederationInstance(c.Context(), payload.SourceInstanceURL, h.resolver, h.instanceRepo)
	if err != nil {
		return err
	}

	status := "pending"
	if policyBool(policy.AutoApproveFriendlinkCitation, false) && h.isFriendLink(c.Context(), payload.SourceInstanceURL) {
		status = "approved"
	}

	citationType := strings.TrimSpace(payload.CitationType)
	if citationType == "" {
		citationType = "reference"
	}

	citation := &federation.FederatedCitation{
		SourceInstanceID: instance.ID,
		SourcePostURL:    payload.SourcePost.URL,
		SourcePostTitle:  toOptionalString(payload.SourcePost.Title),
		TargetArticleID:  article.ID,
		CitationContext:  toOptionalString(payload.CitationContext),
		CitationType:     citationType,
		Status:           status,
		RequestedAt:      time.Now().UTC(),
	}
	if err := h.citationRepo.Create(c.Context(), citation); err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "创建引用记录失败", err)
	}

	// TODO: 当状态为 approved 时，写入特殊评论（需要评论模块支持）。

	resp := contract.FederationCitationResponseResp{
		CitationID: citation.ID,
		Status:     status,
	}
	log.Printf("[federation] 入站 引用申请 source=%s target_post=%s citation_id=%d status=%s key_id=%s", payload.SourceInstanceURL, payload.TargetPostID, citation.ID, status, signature.KeyID)
	return response.Success(c, resp)
}

func (h *FederationCitationHandler) resolveTargetArticle(ctx context.Context, targetID string) (*content.Article, error) {
	if numericID, err := strconv.ParseInt(targetID, 10, 64); err == nil {
		return h.contentRepo.GetArticleByID(ctx, numericID)
	}
	return h.contentRepo.GetArticleByShortURL(ctx, targetID)
}

func (h *FederationCitationHandler) isFriendLink(ctx context.Context, baseURL string) bool {
	if h.linkRepo == nil {
		return false
	}
	_, err := h.linkRepo.FindByURL(ctx, strings.TrimRight(baseURL, "/"))
	return err == nil
}
