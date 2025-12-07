package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/friendlink"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FriendLinkHandler struct {
	svc *friendlink.Service
}

func NewFriendLinkHandler(svc *friendlink.Service) *FriendLinkHandler {
	return &FriendLinkHandler{svc: svc}
}

type friendLinkApplicationRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Message     string `json:"message"`
	UserID      *int64 `json:"userId"`
}

func (h *FriendLinkHandler) SubmitApplication(c *fiber.Ctx) error {
	var req friendLinkApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	if req.URL == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "友链 URL 不能为空")
	}
	cmd := friendlink.SubmitCommand{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Message:     req.Message,
		UserID:      req.UserID,
	}
	result, err := h.svc.Submit(c.Context(), cmd)
	if err != nil {
		return err
	}
	msg := "友链申请已提交成功，我们会尽快审核"
	if !result.Created {
		msg = "你之前的友链申请已更新，感谢耐心等待"
	}
	return response.SuccessWithMessage(c, result.Application, msg)
}
