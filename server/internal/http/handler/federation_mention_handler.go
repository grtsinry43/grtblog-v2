package handler

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

type FederationMentionHandler struct {
	cfgSvc       *federationconfig.Service
	instanceRepo federation.FederationInstanceRepository
	mentionRepo  federation.FederatedMentionRepository
	userRepo     identity.Repository
	resolver     *fedinfra.Resolver
	verifier     *fedinfra.Verifier
}

func NewFederationMentionHandler(
	cfgSvc *federationconfig.Service,
	instanceRepo federation.FederationInstanceRepository,
	mentionRepo federation.FederatedMentionRepository,
	userRepo identity.Repository,
	resolver *fedinfra.Resolver,
	verifier *fedinfra.Verifier,
) *FederationMentionHandler {
	return &FederationMentionHandler{
		cfgSvc:       cfgSvc,
		instanceRepo: instanceRepo,
		mentionRepo:  mentionRepo,
		userRepo:     userRepo,
		resolver:     resolver,
		verifier:     verifier,
	}
}

// NotifyMention handles cross-site mention notifications.
func (h *FederationMentionHandler) NotifyMention(c *fiber.Ctx) error {
	body := c.Body()
	req, err := parseFederationRequest(c)
	if err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求解析失败", err)
	}

	signature, err := h.verifier.VerifyRequest(c.Context(), req, body)
	if err != nil {
		return response.NewBizErrorWithMsg(response.Unauthorized, "签名校验失败")
	}

	var payload contract.FederationMentionNotifyReq
	if err := json.Unmarshal(body, &payload); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(payload.SourceInstanceURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "source_instance_url 不能为空")
	}
	if strings.TrimSpace(payload.SourcePost.URL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "source_post.url 不能为空")
	}
	if strings.TrimSpace(payload.MentionedUser) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "mentioned_user 不能为空")
	}
	if strings.TrimSpace(payload.MentionContext) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "mention_context 不能为空")
	}
	if signature != nil && signature.BaseURL != "" && !sameBaseURL(signature.BaseURL, payload.SourceInstanceURL) {
		return response.NewBizErrorWithMsg(response.Unauthorized, "签名来源与请求不一致")
	}

	settings, err := h.cfgSvc.Settings(c.Context())
	if err != nil || !settings.Enabled {
		return response.NewBizErrorWithMsg(response.Unauthorized, "联合未启用")
	}
	policy := parseFederationPolicy(settings)
	if !policyBool(policy.AllowMention, true) {
		return response.NewBizErrorWithMsg(response.Unauthorized, "未允许被提及")
	}
	if !settings.AllowInbound {
		return response.NewBizErrorWithMsg(response.Unauthorized, "已关闭入站请求")
	}

	user, err := h.userRepo.FindByUsername(c.Context(), payload.MentionedUser)
	if err != nil {
		if errors.Is(err, identity.ErrUserNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "用户不存在")
		}
		return response.NewBizErrorWithCause(response.ServerError, "用户查询失败", err)
	}

	instance, err := ensureFederationInstance(c.Context(), payload.SourceInstanceURL, h.resolver, h.instanceRepo)
	if err != nil {
		return err
	}

	mentionType := strings.TrimSpace(payload.MentionType)
	if mentionType == "" {
		mentionType = "discussion"
	}

	mention := &federation.FederatedMention{
		SourceInstanceID: instance.ID,
		SourcePostURL:    payload.SourcePost.URL,
		SourcePostTitle:  toOptionalString(payload.SourcePost.Title),
		MentionedUserID:  user.ID,
		MentionContext:   payload.MentionContext,
		MentionType:      mentionType,
		IsRead:           false,
		CreatedAt:        time.Now().UTC(),
	}
	if err := h.mentionRepo.Create(c.Context(), mention); err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "写入提及失败", err)
	}

	// TODO: 转为站内信或通知（需要消息模块支持）。

	resp := contract.FederationMentionNotifyResp{
		MentionID: mention.ID,
		Delivered: true,
	}
	return response.Success(c, resp)
}
