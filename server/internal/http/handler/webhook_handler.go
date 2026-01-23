package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/webhook"
	domainwebhook "github.com/grtsinry43/grtblog-v2/server/internal/domain/webhook"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type WebhookHandler struct {
	svc *webhook.Service
}

func NewWebhookHandler(svc *webhook.Service) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

// ListEvents godoc
// @Summary 获取 Webhook 事件列表
// @Tags Webhook
// @Produce json
// @Success 200 {object} contract.WebhookEventListResp
// @Security BearerAuth
// @Router /admin/webhooks/events [get]
// @Security JWTAuth
func (h *WebhookHandler) ListEvents(c *fiber.Ctx) error {
	resp := contract.WebhookEventListResp{
		Events: h.svc.ListEvents(),
	}
	return response.Success(c, resp)
}

// ListWebhooks godoc
// @Summary 获取 Webhook 列表
// @Tags Webhook
// @Produce json
// @Success 200 {object} []contract.WebhookResp
// @Security BearerAuth
// @Router /admin/webhooks [get]
// @Security JWTAuth
func (h *WebhookHandler) ListWebhooks(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}
	resp := make([]contract.WebhookResp, len(items))
	for i, item := range items {
		resp[i] = mapWebhookResp(item)
	}
	return response.Success(c, resp)
}

// CreateWebhook godoc
// @Summary 创建 Webhook
// @Tags Webhook
// @Accept json
// @Produce json
// @Param request body contract.CreateWebhookReq true "创建 Webhook 参数"
// @Success 200 {object} contract.WebhookResp
// @Security BearerAuth
// @Router /admin/webhooks [post]
// @Security JWTAuth
func (h *WebhookHandler) CreateWebhook(c *fiber.Ctx) error {
	var req contract.CreateWebhookReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	hook := &domainwebhook.Webhook{
		Name:            req.Name,
		URL:             req.URL,
		Events:          req.Events,
		Headers:         req.Headers,
		PayloadTemplate: req.PayloadTemplate,
		IsEnabled:       req.IsEnabled,
	}
	if err := h.svc.Create(c.Context(), hook); err != nil {
		if errors.Is(err, domainwebhook.ErrWebhookInvalidEvents) {
			return response.NewBizErrorWithMsg(response.ParamsError, "事件列表无效")
		}
		return err
	}

	resp := mapWebhookResp(hook)
	return response.SuccessWithMessage(c, resp, "Webhook 创建成功")
}

// UpdateWebhook godoc
// @Summary 更新 Webhook
// @Tags Webhook
// @Accept json
// @Produce json
// @Param id path int true "Webhook ID"
// @Param request body contract.UpdateWebhookReq true "更新 Webhook 参数"
// @Success 200 {object} contract.WebhookResp
// @Security BearerAuth
// @Router /admin/webhooks/{id} [put]
// @Security JWTAuth
func (h *WebhookHandler) UpdateWebhook(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的 Webhook ID")
	}

	var req contract.UpdateWebhookReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	hook := &domainwebhook.Webhook{
		ID:              id,
		Name:            req.Name,
		URL:             req.URL,
		Events:          req.Events,
		Headers:         req.Headers,
		PayloadTemplate: req.PayloadTemplate,
		IsEnabled:       req.IsEnabled,
	}
	if err := h.svc.Update(c.Context(), hook); err != nil {
		if errors.Is(err, domainwebhook.ErrWebhookInvalidEvents) {
			return response.NewBizErrorWithMsg(response.ParamsError, "事件列表无效")
		}
		if errors.Is(err, domainwebhook.ErrWebhookNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "Webhook 不存在")
		}
		return err
	}

	updated, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	resp := mapWebhookResp(updated)
	return response.SuccessWithMessage(c, resp, "Webhook 更新成功")
}

// DeleteWebhook godoc
// @Summary 删除 Webhook
// @Tags Webhook
// @Produce json
// @Param id path int true "Webhook ID"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /admin/webhooks/{id} [delete]
// @Security JWTAuth
func (h *WebhookHandler) DeleteWebhook(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的 Webhook ID")
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		if errors.Is(err, domainwebhook.ErrWebhookNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "Webhook 不存在")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "Webhook 删除成功")
}

// TestWebhook godoc
// @Summary 测试 Webhook
// @Tags Webhook
// @Accept json
// @Produce json
// @Param id path int true "Webhook ID"
// @Param request body contract.WebhookTestReq false "测试参数"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /admin/webhooks/{id}/test [post]
// @Security JWTAuth
func (h *WebhookHandler) TestWebhook(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的 Webhook ID")
	}

	var req contract.WebhookTestReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	if err := h.svc.Test(c.Context(), id, req.EventName); err != nil {
		if errors.Is(err, domainwebhook.ErrWebhookInvalidEvents) {
			return response.NewBizErrorWithMsg(response.ParamsError, "事件列表无效")
		}
		if errors.Is(err, domainwebhook.ErrWebhookNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "Webhook 不存在")
		}
		if errors.Is(err, domainwebhook.ErrWebhookDeliveryFailed) {
			return response.SuccessWithMessage[any](c, nil, "Webhook 测试失败，已记录投递历史")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "Webhook 测试成功")
}

// ListHistory godoc
// @Summary 获取 Webhook 投递历史
// @Tags Webhook
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param webhookId query int false "Webhook ID"
// @Param eventName query string false "事件名称"
// @Param isTest query bool false "是否测试"
// @Success 200 {object} contract.WebhookHistoryListResp
// @Security BearerAuth
// @Router /admin/webhooks/deliveries [get]
// @Security JWTAuth
func (h *WebhookHandler) ListHistory(c *fiber.Ctx) error {
	query := contract.WebhookHistoryListReq{
		Page:     1,
		PageSize: 10,
	}
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		query.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize", "10")); err == nil && pageSize > 0 && pageSize <= 100 {
		query.PageSize = pageSize
	}
	if webhookID, err := strconv.ParseInt(c.Query("webhookId"), 10, 64); err == nil {
		query.WebhookID = &webhookID
	}
	if eventName := c.Query("eventName"); eventName != "" {
		query.EventName = &eventName
	}
	if isTestStr := c.Query("isTest"); isTestStr != "" {
		if isTest, err := strconv.ParseBool(isTestStr); err == nil {
			query.IsTest = &isTest
		}
	}

	items, total, err := h.svc.ListHistory(c.Context(), domainwebhook.DeliveryHistoryListOptions{
		Page:      query.Page,
		PageSize:  query.PageSize,
		WebhookID: query.WebhookID,
		EventName: query.EventName,
		IsTest:    query.IsTest,
	})
	if err != nil {
		return err
	}

	resp := make([]contract.WebhookHistoryResp, len(items))
	for i, item := range items {
		resp[i] = mapHistoryResp(item)
	}

	listResp := contract.WebhookHistoryListResp{
		Items: resp,
		Total: total,
		Page:  query.Page,
		Size:  query.PageSize,
	}
	return response.Success(c, listResp)
}

// ReplayHistory godoc
// @Summary 重放 Webhook 投递历史
// @Tags Webhook
// @Produce json
// @Param id path int true "失败记录 ID"
// @Success 200 {object} any
// @Security BearerAuth
// @Router /admin/webhooks/deliveries/{id}/replay [post]
// @Security JWTAuth
func (h *WebhookHandler) ReplayHistory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的记录ID")
	}
	if err := h.svc.Replay(c.Context(), id); err != nil {
		if errors.Is(err, domainwebhook.ErrDeliveryHistoryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "记录不存在")
		}
		if errors.Is(err, domainwebhook.ErrWebhookDeliveryFailed) {
			return response.SuccessWithMessage[any](c, nil, "重放失败，已记录投递历史")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "重放成功")
}

func mapWebhookResp(hook *domainwebhook.Webhook) contract.WebhookResp {
	headers := hook.Headers
	if headers == nil {
		headers = map[string]string{}
	}
	return contract.WebhookResp{
		ID:              hook.ID,
		Name:            hook.Name,
		URL:             hook.URL,
		Events:          hook.Events,
		Headers:         headers,
		PayloadTemplate: hook.PayloadTemplate,
		IsEnabled:       hook.IsEnabled,
		CreatedAt:       hook.CreatedAt,
		UpdatedAt:       hook.UpdatedAt,
	}
}

func mapHistoryResp(history *domainwebhook.DeliveryHistory) contract.WebhookHistoryResp {
	headers := history.RequestHeaders
	if headers == nil {
		headers = map[string]string{}
	}
	responseHeaders := history.ResponseHeaders
	if responseHeaders == nil {
		responseHeaders = map[string]string{}
	}
	return contract.WebhookHistoryResp{
		ID:              history.ID,
		WebhookID:       history.WebhookID,
		EventName:       history.EventName,
		RequestURL:      history.RequestURL,
		RequestHeaders:  headers,
		RequestBody:     history.RequestBody,
		ResponseStatus:  history.ResponseStatus,
		ResponseHeaders: responseHeaders,
		ResponseBody:    history.ResponseBody,
		ErrorMessage:    history.ErrorMessage,
		IsTest:          history.IsTest,
		CreatedAt:       history.CreatedAt,
	}
}
