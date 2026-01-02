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

func (h *WebsiteInfoHandler) listAll(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}
	return response.Success(c, toWebsiteInfoListResponse(items))
}

// PublicList godoc
// @Summary 公开获取网站信息
// @Tags WebsiteInfo
// @Produce json
// @Success 200 {object} WebsiteInfoListEnvelope
// @Router /public/website-info [get]
func (h *WebsiteInfoHandler) PublicList(c *fiber.Ctx) error {
	return h.listAll(c)
}

// List godoc
// @Summary 获取全部网站信息（需要 config:read 权限）
// @Tags WebsiteInfo
// @Produce json
// @Success 200 {object} WebsiteInfoListEnvelope
// @Security BearerAuth
// @Router /website-info [get]
func (h *WebsiteInfoHandler) List(c *fiber.Ctx) error {
	return h.listAll(c)
}

// WebsiteInfoRequest 用于 swagger 展示请求体
type WebsiteInfoRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Create godoc
// @Summary 新增网站信息
// @Tags WebsiteInfo
// @Accept json
// @Produce json
// @Param request body WebsiteInfoRequest true "网站配置"
// @Success 200 {object} WebsiteInfoDetailEnvelope
// @Security BearerAuth
// @Router /website-info [post]
func (h *WebsiteInfoHandler) Create(c *fiber.Ctx) error {
	var req WebsiteInfoRequest
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
	Audit(c, "website_info.create", map[string]any{"key": info.Key})
	return response.SuccessWithMessage(c, toWebsiteInfoResponse(*info), "created")
}

// Update godoc
// @Summary 更新网站信息
// @Tags WebsiteInfo
// @Accept json
// @Produce json
// @Param key path string true "配置键"
// @Param request body WebsiteInfoRequest true "网站配置"
// @Success 200 {object} WebsiteInfoDetailEnvelope
// @Security BearerAuth
// @Router /website-info/{key} [put]
func (h *WebsiteInfoHandler) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
	}
	var req WebsiteInfoRequest
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
	Audit(c, "website_info.update", map[string]any{"key": info.Key})
	return response.SuccessWithMessage(c, toWebsiteInfoResponse(*info), "updated")
}

// Delete godoc
// @Summary 删除网站信息
// @Tags WebsiteInfo
// @Produce json
// @Param key path string true "配置键"
// @Success 200 {object} GenericMessageEnvelope
// @Security BearerAuth
// @Router /website-info/{key} [delete]
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
	Audit(c, "website_info.delete", map[string]any{"key": key})
	return response.SuccessWithMessage[any](c, nil, "删除成功")
}
