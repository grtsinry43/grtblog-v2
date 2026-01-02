package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FriendLinkHandler struct {
	svc *friendlink.Service
}

func NewFriendLinkHandler(svc *friendlink.Service) *FriendLinkHandler {
	return &FriendLinkHandler{svc: svc}
}

// FriendLinkApplicationRequest 展示友链申请请求体验证结构
type FriendLinkApplicationRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

// SubmitApplication godoc
// @Summary 提交或更新友链申请
// @Tags FriendLink
// @Accept json
// @Produce json
// @Param request body FriendLinkApplicationRequest true "友链申请"
// @Success 200 {object} FriendLinkApplicationResponse
// @Security BearerAuth
// @Router /friend-links/applications [post]
func (h *FriendLinkHandler) SubmitApplication(c *fiber.Ctx) error {
	var req FriendLinkApplicationRequest
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
	cmd := friendlink.SubmitCommand{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Message:     req.Message,
		UserID:      &userID,
	}
	result, err := h.svc.Submit(c.Context(), cmd)
	if err != nil {
		return err
	}
	msg := "友链申请已提交成功，我们会尽快审核"
	if !result.Created {
		msg = "你之前的友链申请已更新，感谢耐心等待"
	}
	Audit(c, "friend_link.submit", map[string]any{"url": req.URL, "name": req.Name, "created": result.Created})
	return response.SuccessWithMessage(c, toFriendLinkApplicationVO(result.Application), msg)
}
