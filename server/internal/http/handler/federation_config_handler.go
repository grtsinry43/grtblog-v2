package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federationconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// FederationConfigHandler provides settings-center style APIs for federation_config.
type FederationConfigHandler struct {
	svc *federationconfig.Service
}

func NewFederationConfigHandler(svc *federationconfig.Service) *FederationConfigHandler {
	return &FederationConfigHandler{svc: svc}
}

// ListFederationConfig lists federation config items.
// @Summary 联合配置列表
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param keys query string false "指定配置 key（逗号分隔）"
// @Success 200 {object} contract.SysConfigTreeResp
// @Security BearerAuth
// @Router /admin/federation/config [get]
// @Security JWTAuth
func (h *FederationConfigHandler) ListFederationConfig(c *fiber.Ctx) error {
	var keys []string
	if raw := c.Query("keys"); raw != "" {
		parts := strings.Split(raw, ",")
		keys = make([]string, 0, len(parts))
		seen := make(map[string]struct{}, len(parts))
		for _, item := range parts {
			key := strings.TrimSpace(item)
			if key == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	items, err := h.svc.ListConfigs(c.Context(), keys)
	if err != nil {
		return err
	}
	tree, err := buildSysConfigTree(items)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "配置解析失败", err)
	}
	return response.Success(c, tree)
}

// UpdateFederationConfig updates federation config items.
// @Summary 更新联合配置
// @Tags FederationAdmin
// @Accept json
// @Produce json
// @Param request body contract.SysConfigBatchUpdateReq true "配置更新参数"
// @Success 200 {object} contract.SysConfigTreeResp
// @Security BearerAuth
// @Router /admin/federation/config [put]
// @Security JWTAuth
func (h *FederationConfigHandler) UpdateFederationConfig(c *fiber.Ctx) error {
	var req contract.SysConfigBatchUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if len(req.Items) == 0 {
		return response.NewBizErrorWithMsg(response.ParamsError, "items 不能为空")
	}

	updates := make([]sysconfig.UpdateItem, 0, len(req.Items))
	for _, item := range req.Items {
		key := strings.TrimSpace(item.Key)
		if key == "" {
			return response.NewBizErrorWithMsg(response.ParamsError, "key 不能为空")
		}
		updates = append(updates, sysconfig.UpdateItem{
			Key:          key,
			Value:        item.Value,
			IsSensitive:  item.IsSensitive,
			GroupPath:    item.GroupPath,
			Label:        item.Label,
			Description:  item.Description,
			ValueType:    item.ValueType,
			EnumOptions:  item.EnumOptions,
			DefaultValue: item.DefaultValue,
			VisibleWhen:  item.VisibleWhen,
			Sort:         item.Sort,
			Meta:         item.Meta,
		})
	}

	updated, err := h.svc.UpdateConfigs(c.Context(), updates)
	if err != nil {
		var validationErr *sysconfig.UpdateValidationError
		if errors.As(err, &validationErr) {
			return response.NewBizErrorWithMsg(response.ParamsError, validationErr.Error())
		}
		return err
	}
	tree, err := buildSysConfigTree(updated)
	if err != nil {
		return response.NewBizErrorWithCause(response.ServerError, "配置解析失败", err)
	}
	return response.SuccessWithMessage(c, tree, "更新成功")
}
