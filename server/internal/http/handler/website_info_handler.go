package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/websiteinfo"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type WebsiteInfoHandler struct {
	svc *websiteinfo.Service
}

func NewWebsiteInfoHandler(svc *websiteinfo.Service) *WebsiteInfoHandler {
	return &WebsiteInfoHandler{svc: svc}
}

func (h *WebsiteInfoHandler) List(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}
	return response.Success(c, items)
}

type websiteInfoRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *WebsiteInfoHandler) Create(c *fiber.Ctx) error {
	var req websiteInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	if req.Key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	info, err := h.svc.Create(c.Context(), websiteinfo.CreateCommand{
		Key:   req.Key,
		Value: req.Value,
	})
	if err != nil {
		return err
	}
	return response.SuccessWithMessage(c, info, "created")
}

func (h *WebsiteInfoHandler) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	var req websiteInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "请求体解析失败")
	}
	cmd := websiteinfo.UpdateCommand{
		Key:   key,
		Value: req.Value,
	}
	info, err := h.svc.Update(c.Context(), cmd)
	if err != nil {
		if errors.Is(err, config.ErrWebsiteInfoNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return err
	}
	return response.SuccessWithMessage(c, info, "updated")
}

func (h *WebsiteInfoHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	if err := h.svc.Delete(c.Context(), key); err != nil {
		if errors.Is(err, config.ErrWebsiteInfoNotFound) {
			return response.NewBizError(response.NotFound)
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "删除成功")
}
