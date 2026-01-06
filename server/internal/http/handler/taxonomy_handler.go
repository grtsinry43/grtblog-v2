package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/taxonomy"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type TaxonomyHandler struct {
	categories *taxonomy.CategoryService
	columns    *taxonomy.ColumnService
	tags       *taxonomy.TagService
}

func NewTaxonomyHandler(categories *taxonomy.CategoryService, columns *taxonomy.ColumnService, tags *taxonomy.TagService) *TaxonomyHandler {
	return &TaxonomyHandler{
		categories: categories,
		columns:    columns,
		tags:       tags,
	}
}

// ListCategories godoc
// @Summary 获取文章分类列表
// @Tags Category
// @Produce json
// @Success 200 {object} []contract.CategoryResp
// @Router /categories [get]
func (h *TaxonomyHandler) ListCategories(c *fiber.Ctx) error {
	items, err := h.categories.List(c.Context())
	if err != nil {
		return err
	}
	resp := make([]contract.CategoryResp, len(items))
	for i, item := range items {
		shortURL := ""
		if item.ShortURL != nil {
			shortURL = *item.ShortURL
		}
		resp[i] = contract.CategoryResp{
			ID:        item.ID,
			Name:      item.Name,
			ShortURL:  shortURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}
	return response.Success(c, resp)
}

// CreateCategory godoc
// @Summary 创建文章分类
// @Tags Category
// @Accept json
// @Produce json
// @Param request body contract.CategoryCreateReq true "分类参数"
// @Success 200 {object} contract.CategoryResp
// @Security BearerAuth
// @Router /admin/categories [post]
func (h *TaxonomyHandler) CreateCategory(c *fiber.Ctx) error {
	var req contract.CategoryCreateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分类名称不能为空")
	}
	if req.ShortURL == nil || strings.TrimSpace(*req.ShortURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分类短链接不能为空")
	}
	item, err := h.categories.Create(c.Context(), req.Name, req.ShortURL)
	if err != nil {
		return err
	}
	shortURL := ""
	if item.ShortURL != nil {
		shortURL = *item.ShortURL
	}
	resp := contract.CategoryResp{
		ID:        item.ID,
		Name:      item.Name,
		ShortURL:  shortURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return response.SuccessWithMessage(c, resp, "创建成功")
}

// UpdateCategory godoc
// @Summary 更新文章分类
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param request body contract.CategoryUpdateReq true "分类参数"
// @Success 200 {object} contract.CategoryResp
// @Security BearerAuth
// @Router /admin/categories/{id} [put]
func (h *TaxonomyHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的分类ID")
	}
	var req contract.CategoryUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分类名称不能为空")
	}
	if req.ShortURL == nil || strings.TrimSpace(*req.ShortURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分类短链接不能为空")
	}
	item, err := h.categories.Update(c.Context(), id, req.Name, req.ShortURL)
	if err != nil {
		if errors.Is(err, content.ErrCategoryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "分类不存在")
		}
		return err
	}
	shortURL := ""
	if item.ShortURL != nil {
		shortURL = *item.ShortURL
	}
	resp := contract.CategoryResp{
		ID:        item.ID,
		Name:      item.Name,
		ShortURL:  shortURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return response.SuccessWithMessage(c, resp, "更新成功")
}

// DeleteCategory godoc
// @Summary 删除文章分类
// @Tags Category
// @Param id path int true "分类ID"
// @Success 200 {object} contract.EmptyRespEnvelope
// @Security BearerAuth
// @Router /admin/categories/{id} [delete]
func (h *TaxonomyHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的分类ID")
	}
	if err := h.categories.Delete(c.Context(), id); err != nil {
		if errors.Is(err, content.ErrCategoryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "分类不存在")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "删除成功")
}

// ListColumns godoc
// @Summary 获取手记分区列表
// @Tags Column
// @Produce json
// @Success 200 {object} []contract.ColumnResp
// @Router /columns [get]
func (h *TaxonomyHandler) ListColumns(c *fiber.Ctx) error {
	items, err := h.columns.List(c.Context())
	if err != nil {
		return err
	}
	resp := make([]contract.ColumnResp, len(items))
	for i, item := range items {
		shortURL := ""
		if item.ShortURL != nil {
			shortURL = *item.ShortURL
		}
		resp[i] = contract.ColumnResp{
			ID:        item.ID,
			Name:      item.Name,
			ShortURL:  shortURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}
	return response.Success(c, resp)
}

// CreateColumn godoc
// @Summary 创建手记分区
// @Tags Column
// @Accept json
// @Produce json
// @Param request body contract.ColumnCreateReq true "分区参数"
// @Success 200 {object} contract.ColumnResp
// @Security BearerAuth
// @Router /admin/columns [post]
func (h *TaxonomyHandler) CreateColumn(c *fiber.Ctx) error {
	var req contract.ColumnCreateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分区名称不能为空")
	}
	if req.ShortURL == nil || strings.TrimSpace(*req.ShortURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分区短链接不能为空")
	}
	item, err := h.columns.Create(c.Context(), req.Name, req.ShortURL)
	if err != nil {
		return err
	}
	shortURL := ""
	if item.ShortURL != nil {
		shortURL = *item.ShortURL
	}
	resp := contract.ColumnResp{
		ID:        item.ID,
		Name:      item.Name,
		ShortURL:  shortURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return response.SuccessWithMessage(c, resp, "创建成功")
}

// UpdateColumn godoc
// @Summary 更新手记分区
// @Tags Column
// @Accept json
// @Produce json
// @Param id path int true "分区ID"
// @Param request body contract.ColumnUpdateReq true "分区参数"
// @Success 200 {object} contract.ColumnResp
// @Security BearerAuth
// @Router /admin/columns/{id} [put]
func (h *TaxonomyHandler) UpdateColumn(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的分区ID")
	}
	var req contract.ColumnUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分区名称不能为空")
	}
	if req.ShortURL == nil || strings.TrimSpace(*req.ShortURL) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "分区短链接不能为空")
	}
	item, err := h.columns.Update(c.Context(), id, req.Name, req.ShortURL)
	if err != nil {
		if errors.Is(err, content.ErrColumnNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "手记分区不存在")
		}
		return err
	}
	shortURL := ""
	if item.ShortURL != nil {
		shortURL = *item.ShortURL
	}
	resp := contract.ColumnResp{
		ID:        item.ID,
		Name:      item.Name,
		ShortURL:  shortURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return response.SuccessWithMessage(c, resp, "更新成功")
}

// DeleteColumn godoc
// @Summary 删除手记分区
// @Tags Column
// @Param id path int true "分区ID"
// @Success 200 {object} contract.EmptyRespEnvelope
// @Security BearerAuth
// @Router /admin/columns/{id} [delete]
func (h *TaxonomyHandler) DeleteColumn(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的分区ID")
	}
	if err := h.columns.Delete(c.Context(), id); err != nil {
		if errors.Is(err, content.ErrColumnNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "手记分区不存在")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "删除成功")
}

// ListTags godoc
// @Summary 获取标签列表
// @Tags Tag
// @Produce json
// @Success 200 {object} []contract.TagItemResp
// @Router /tags [get]
func (h *TaxonomyHandler) ListTags(c *fiber.Ctx) error {
	items, err := h.tags.List(c.Context())
	if err != nil {
		return err
	}
	resp := make([]contract.TagItemResp, len(items))
	for i, item := range items {
		_ = copier.Copy(&resp[i], item)
	}
	return response.Success(c, resp)
}

// CreateTag godoc
// @Summary 创建标签
// @Tags Tag
// @Accept json
// @Produce json
// @Param request body contract.TagCreateReq true "标签参数"
// @Success 200 {object} contract.TagItemResp
// @Security BearerAuth
// @Router /admin/tags [post]
func (h *TaxonomyHandler) CreateTag(c *fiber.Ctx) error {
	var req contract.TagCreateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "标签名称不能为空")
	}
	item, err := h.tags.Create(c.Context(), req.Name)
	if err != nil {
		return err
	}
	var resp contract.TagItemResp
	_ = copier.Copy(&resp, item)
	return response.SuccessWithMessage(c, resp, "创建成功")
}

// UpdateTag godoc
// @Summary 更新标签
// @Tags Tag
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param request body contract.TagUpdateReq true "标签参数"
// @Success 200 {object} contract.TagItemResp
// @Security BearerAuth
// @Router /admin/tags/{id} [put]
func (h *TaxonomyHandler) UpdateTag(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的标签ID")
	}
	var req contract.TagUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "标签名称不能为空")
	}
	item, err := h.tags.Update(c.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "标签不存在")
		}
		return err
	}
	var resp contract.TagItemResp
	_ = copier.Copy(&resp, item)
	return response.SuccessWithMessage(c, resp, "更新成功")
}

// DeleteTag godoc
// @Summary 删除标签
// @Tags Tag
// @Param id path int true "标签ID"
// @Success 200 {object} contract.EmptyRespEnvelope
// @Security BearerAuth
// @Router /admin/tags/{id} [delete]
func (h *TaxonomyHandler) DeleteTag(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的标签ID")
	}
	if err := h.tags.Delete(c.Context(), id); err != nil {
		if errors.Is(err, content.ErrTagNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "标签不存在")
		}
		return err
	}
	return response.SuccessWithMessage[any](c, nil, "删除成功")
}
