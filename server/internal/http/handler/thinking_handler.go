package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/thinking"
	domainthinking "github.com/grtsinry43/grtblog-v2/server/internal/domain/thinking"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type ThinkingHandler struct {
	svc *thinking.Service
}

func NewThinkingHandler(svc *thinking.Service) *ThinkingHandler {
	return &ThinkingHandler{svc: svc}
}

// CreateThinking godoc
// @Summary 发布回响
// @Tags Thinking
// @Accept json
// @Produce json
// @Param request body contract.CreateThinkingReq true "创建参数"
// @Success 200 {object} contract.ThinkingResp
// @Security JWTAuth
// @Router /thinkings [post]
func (h *ThinkingHandler) CreateThinking(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var req contract.CreateThinkingReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	author := req.Author
	if author == "" {
		author = "Admin"
	}

	created, err := h.svc.Create(c.Context(), thinking.CreateThinkingCmd{
		Content: req.Content,
		Author:  author,
	})
	if err != nil {
		return h.mapError(c, err)
	}

	Audit(c, "thinking.create", map[string]any{
		"thinkingId":     created.ID,
		"thinkingAuthor": created.Author,
		"userId":         claims.UserID,
	})

	return response.Success(c, toThinkingResp(created))
}

// UpdateThinking godoc
// @Summary 更新回响
// @Tags Thinking
// @Accept json
// @Produce json
// @Param request body contract.UpdateThinkingReq true "更新参数"
// @Success 200
// @Security JWTAuth
// @Router /thinkings [put]
func (h *ThinkingHandler) UpdateThinking(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var req contract.UpdateThinkingReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "无效的请求体", err)
	}

	_, err := h.svc.Update(c.Context(), thinking.UpdateThinkingCmd{
		ID:      req.ID,
		Content: req.Content,
	})
	if err != nil {
		return h.mapError(c, err)
	}

	Audit(c, "thinking.update", map[string]any{
		"thinkingId": req.ID,
		"userId":     claims.UserID,
	})

	return response.SuccessWithMessage[any](c, nil, "更新回响成功")
}

// ListThinkings godoc
// @Summary 获取回响列表
// @Tags Thinking
// @Param page query int false "页码"
// @Param pageSize query int false "页大小"
// @Success 200 {object} contract.ListThinkingResp
// @Router /thinkings [get]
func (h *ThinkingHandler) ListThinkings(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	items, total, err := h.svc.List(c.Context(), pageSize, offset)
	if err != nil {
		return h.mapError(c, err)
	}

	resItems := make([]*contract.ThinkingResp, len(items))
	for i, item := range items {
		resItems[i] = toThinkingResp(item)
	}

	return response.Success(c, contract.ListThinkingResp{
		Items: resItems,
		Total: total,
	})
}

// DeleteThinking godoc
// @Summary 删除回响
// @Tags Thinking
// @Param id path int true "回响ID"
// @Success 200
// @Security JWTAuth
// @Router /thinkings/{id} [delete]
func (h *ThinkingHandler) DeleteThinking(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的ID")
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return h.mapError(c, err)
	}

	Audit(c, "thinking.delete", map[string]any{
		"thinkingId": id,
		"userId":     claims.UserID,
	})

	return response.SuccessWithMessage[any](c, nil, "回响删除成功")
}

func (h *ThinkingHandler) mapError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domainthinking.ErrThinkingNotFound):
		return response.NewBizErrorWithMsg(response.NotFound, "回响不存在")
	case errors.Is(err, domainthinking.ErrThinkingContentEmpty):
		return response.NewBizErrorWithMsg(response.ParamsError, "内容不能为空")
	default:
		return err
	}
}

func toThinkingResp(t *domainthinking.Thinking) *contract.ThinkingResp {
	return &contract.ThinkingResp{
		ID:        t.ID,
		CommentID: t.CommentID,
		Content:   t.Content,
		Author:    t.Author,
		Metrics: contract.ThinkingMetrics{
			Views:    t.Metrics.Views,
			Likes:    t.Metrics.Likes,
			Comments: t.Metrics.Comments,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
