package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	appnav "github.com/grtsinry43/grtblog-v2/server/internal/app/navigation"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/navigation"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type NavMenuHandler struct {
	svc *appnav.Service
}

func NewNavMenuHandler(svc *appnav.Service) *NavMenuHandler {
	return &NavMenuHandler{svc: svc}
}

// ListPublic godoc
// @Summary 公开获取导航菜单
// @Tags NavMenu
// @Produce json
// @Success 200 {object} contract.NavMenuListRespEnvelope
// @Router /public/nav-menus [get]
func (h *NavMenuHandler) ListPublic(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}

	return response.Success(c, buildNavMenuTree(items))
}

// ListAdmin godoc
// @Summary 获取导航菜单（管理端）
// @Tags NavMenu
// @Produce json
// @Success 200 {object} contract.NavMenuListRespEnvelope
// @Security BearerAuth
// @Router /admin/nav-menus [get]
func (h *NavMenuHandler) ListAdmin(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}

	return response.Success(c, buildNavMenuTree(items))
}

// Create godoc
// @Summary 创建导航菜单
// @Tags NavMenu
// @Accept json
// @Produce json
// @Param request body contract.CreateNavMenuReq true "导航菜单"
// @Success 200 {object} contract.NavMenuDetailRespEnvelope
// @Security BearerAuth
// @Router /admin/nav-menus [post]
func (h *NavMenuHandler) Create(c *fiber.Ctx) error {
	var req contract.CreateNavMenuReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "菜单名称不能为空")
	}
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "菜单链接不能为空")
	}

	created, err := h.svc.Create(c.Context(), appnav.CreateNavMenuCmd{
		Name:     name,
		URL:      url,
		ParentID: req.ParentID,
		Icon:     req.Icon,
	})
	if err != nil {
		return err
	}

	return response.SuccessWithMessage[contract.NavMenuResp](c, toNavMenuResp(created), "菜单创建成功")
}

// Update godoc
// @Summary 更新导航菜单
// @Tags NavMenu
// @Accept json
// @Produce json
// @Param id path int64 true "菜单ID"
// @Param request body contract.UpdateNavMenuReq true "导航菜单"
// @Success 200 {object} contract.NavMenuDetailRespEnvelope
// @Security BearerAuth
// @Router /admin/nav-menus/{id} [put]
func (h *NavMenuHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的菜单ID")
	}

	var req contract.UpdateNavMenuReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "菜单名称不能为空")
	}
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "菜单链接不能为空")
	}

	updated, err := h.svc.Update(c.Context(), appnav.UpdateNavMenuCmd{
		ID:       id,
		Name:     name,
		URL:      url,
		ParentID: req.ParentID,
		Icon:     req.Icon,
		Sort:     req.Sort,
	})
	if err != nil {
		return err
	}

	return response.SuccessWithMessage[contract.NavMenuResp](c, toNavMenuResp(updated), "菜单更新成功")
}

// Delete godoc
// @Summary 删除导航菜单
// @Tags NavMenu
// @Produce json
// @Param id path int64 true "菜单ID"
// @Success 200 {object} contract.GenericMessageEnvelope
// @Security BearerAuth
// @Router /admin/nav-menus/{id} [delete]
func (h *NavMenuHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的菜单ID")
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return err
	}

	return response.SuccessWithMessage[any](c, nil, "菜单已删除")
}

// Reorder godoc
// @Summary 更新导航菜单排序
// @Tags NavMenu
// @Accept json
// @Produce json
// @Param request body contract.ReorderNavMenuReq true "排序信息"
// @Success 200 {object} contract.GenericMessageEnvelope
// @Security BearerAuth
// @Router /admin/nav-menus/reorder [put]
func (h *NavMenuHandler) Reorder(c *fiber.Ctx) error {
	var req contract.ReorderNavMenuReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}

	if len(req.Items) == 0 {
		return response.NewBizErrorWithMsg(response.ParamsError, "排序数据不能为空")
	}

	items := make([]appnav.NavMenuOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, appnav.NavMenuOrderItem{
			ID:       item.ID,
			ParentID: item.ParentID,
			Sort:     item.Sort,
		})
	}

	if err := h.svc.UpdateOrder(c.Context(), items); err != nil {
		return err
	}

	return response.SuccessWithMessage[any](c, nil, "菜单排序已更新")
}

func toNavMenuResp(item *navigation.NavMenu) contract.NavMenuResp {
	if item == nil {
		return contract.NavMenuResp{}
	}
	return contract.NavMenuResp{
		ID:        item.ID,
		Name:      item.Name,
		URL:       item.URL,
		Icon:      item.Icon,
		Sort:      item.Sort,
		ParentID:  item.ParentID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func buildNavMenuTree(items []*navigation.NavMenu) []contract.NavMenuResp {
	if len(items) == 0 {
		// 返回空数组，而不是nil
		return []contract.NavMenuResp{}
	}

	childrenMap := map[int64][]*contract.NavMenuResp{}
	nodeMap := map[int64]*contract.NavMenuResp{}

	for _, item := range items {
		resp := toNavMenuResp(item)
		respCopy := resp
		pid := int64(0)
		if item.ParentID != nil {
			pid = *item.ParentID
		}
		nodeMap[item.ID] = &respCopy
		childrenMap[pid] = append(childrenMap[pid], &respCopy)
	}

	var attach func(parentID int64)
	attach = func(parentID int64) {
		node := nodeMap[parentID]
		if node == nil {
			return
		}
		children := childrenMap[parentID]
		if len(children) == 0 {
			return
		}
		resultChildren := make([]contract.NavMenuResp, 0, len(children))
		for _, child := range children {
			attach(child.ID)
			resultChildren = append(resultChildren, *child)
		}
		node.Children = resultChildren
	}

	roots := childrenMap[0]
	if len(roots) == 0 {
		return nil
	}

	result := make([]contract.NavMenuResp, 0, len(roots))
	for _, root := range roots {
		attach(root.ID)
		result = append(result, *root)
	}

	return result
}
