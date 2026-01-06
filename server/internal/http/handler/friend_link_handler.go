package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FriendLinkHandler struct {
	svc *friendlink.Service
}

func NewFriendLinkHandler(svc *friendlink.Service) *FriendLinkHandler {
	return &FriendLinkHandler{svc: svc}
}

// SubmitApplication godoc
// @Summary 提交或更新友链申请
// @Tags FriendLink
// @Accept json
// @Produce json
// @Param request body contract.FriendLinkApplicationReq true "友链申请"
// @Success 200 {object} contract.FriendLinkApplicationRespEnvelope
// @Security BearerAuth
// @Router /friend-links/applications [post]
func (h *FriendLinkHandler) SubmitApplication(c *fiber.Ctx) error {
	var req contract.FriendLinkApplicationReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}
	if req.URL == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "友链 URL 不能为空")
	}
	userID := claims.UserID
	var cmd friendlink.SubmitCmd
	if err := copier.Copy(&cmd, req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体映射失败")
	}
	cmd.UserID = &userID
	result, err := h.svc.Submit(c.Context(), cmd)
	if err != nil {
		return err
	}
	msg := "友链申请已提交成功，我们会尽快审核"
	if !result.Created {
		msg = "你之前的友链申请已更新，感谢耐心等待"
	}
	Audit(c, "friend-link.submit", map[string]any{"url": req.URL, "name": req.Name, "created": result.Created})
	return response.SuccessWithMessage(c, contract.ToFriendLinkApplicationResp(result.Application), msg)
}
