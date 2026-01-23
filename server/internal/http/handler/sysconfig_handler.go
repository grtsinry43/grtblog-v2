package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	domainconfig "github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type SysConfigHandler struct {
	svc *sysconfig.Service
}

func NewSysConfigHandler(svc *sysconfig.Service) *SysConfigHandler {
	return &SysConfigHandler{svc: svc}
}

// ListSysConfig godoc
// @Summary 获取系统配置
// @Tags SysConfig
// @Produce json
// @Param keys query string false "配置 key 列表，逗号分隔"
// @Success 200 {object} contract.SysConfigTreeResp
// @Security BearerAuth
// @Router /admin/sysconfig [get]
// @Security JWTAuth
func (h *SysConfigHandler) ListSysConfig(c *fiber.Ctx) error {
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

// UpdateSysConfig godoc
// @Summary 批量更新系统配置
// @Tags SysConfig
// @Accept json
// @Produce json
// @Param request body contract.SysConfigBatchUpdateReq true "批量更新参数"
// @Success 200 {object} contract.SysConfigTreeResp
// @Security BearerAuth
// @Router /admin/sysconfig [put]
// @Security JWTAuth
func (h *SysConfigHandler) UpdateSysConfig(c *fiber.Ctx) error {
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

type sysConfigGroupNode struct {
	Key      string
	Path     string
	Label    string
	Children []*sysConfigGroupNode
	Items    []contract.SysConfigItemResp
}

var (
	emptyJSONArray  = json.RawMessage("[]")
	emptyJSONObject = json.RawMessage("{}")
)

const (
	valueTypeString = "string"
	valueTypeNumber = "number"
	valueTypeBool   = "bool"
	valueTypeEnum   = "enum"
	valueTypeJSON   = "json"
)

func buildSysConfigTree(items []domainconfig.SysConfig) (contract.SysConfigTreeResp, error) {
	rootItems := make([]contract.SysConfigItemResp, 0)
	rootGroups := make([]*sysConfigGroupNode, 0)
	groupIndex := make(map[string]*sysConfigGroupNode)

	for _, cfg := range items {
		item, err := mapSysConfigItem(cfg)
		if err != nil {
			return contract.SysConfigTreeResp{}, err
		}
		groupPath := strings.Trim(strings.TrimSpace(cfg.GroupPath), "/")
		if groupPath == "" {
			rootItems = append(rootItems, item)
			continue
		}
		parts := strings.Split(groupPath, "/")
		var parent *sysConfigGroupNode
		currentPath := ""
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = currentPath + "/" + part
			}
			node, exists := groupIndex[currentPath]
			if !exists {
				node = &sysConfigGroupNode{
					Key:   part,
					Path:  currentPath,
					Label: part,
				}
				groupIndex[currentPath] = node
				if parent == nil {
					rootGroups = append(rootGroups, node)
				} else {
					parent.Children = append(parent.Children, node)
				}
			}
			parent = node
		}
		if parent == nil {
			rootItems = append(rootItems, item)
			continue
		}
		parent.Items = append(parent.Items, item)
	}

	return contract.SysConfigTreeResp{
		Groups: convertSysConfigGroups(rootGroups),
		Items:  rootItems,
	}, nil
}

func convertSysConfigGroups(nodes []*sysConfigGroupNode) []contract.SysConfigGroupResp {
	if len(nodes) == 0 {
		return []contract.SysConfigGroupResp{}
	}
	result := make([]contract.SysConfigGroupResp, len(nodes))
	for i, node := range nodes {
		result[i] = contract.SysConfigGroupResp{
			Key:      node.Key,
			Path:     node.Path,
			Label:    node.Label,
			Children: convertSysConfigGroups(node.Children),
			Items:    node.Items,
		}
	}
	return result
}

func mapSysConfigItem(cfg domainconfig.SysConfig) (contract.SysConfigItemResp, error) {
	valueType, err := normalizeValueType(cfg.ValueType)
	if err != nil {
		return contract.SysConfigItemResp{}, err
	}
	valueRaw, err := valueToJSON(valueType, cfg.Value)
	if err != nil {
		return contract.SysConfigItemResp{}, err
	}
	defaultRaw, err := defaultValueToJSON(valueType, cfg.DefaultValue)
	if err != nil {
		return contract.SysConfigItemResp{}, err
	}

	resp := contract.SysConfigItemResp{
		Key:          cfg.Key,
		GroupPath:    cfg.GroupPath,
		Label:        cfg.Label,
		Description:  cfg.Description,
		ValueType:    valueType,
		EnumOptions:  normalizeRaw(cfg.EnumOptions, emptyJSONArray),
		DefaultValue: defaultRaw,
		VisibleWhen:  normalizeRaw(cfg.VisibleWhen, emptyJSONArray),
		Sort:         cfg.Sort,
		Meta:         normalizeRaw(cfg.Meta, emptyJSONObject),
		IsSensitive:  cfg.IsSensitive,
		CreatedAt:    cfg.CreatedAt,
		UpdatedAt:    cfg.UpdatedAt,
	}
	if !cfg.IsSensitive && valueRaw != nil {
		resp.Value = valueRaw
	}
	return resp, nil
}

func normalizeValueType(valueType string) (string, error) {
	valueType = strings.TrimSpace(strings.ToLower(valueType))
	if valueType == "" {
		return valueTypeString, nil
	}
	switch valueType {
	case valueTypeString, valueTypeNumber, valueTypeBool, valueTypeEnum, valueTypeJSON:
		return valueType, nil
	default:
		return "", fmt.Errorf("unsupported valueType: %s", valueType)
	}
}

func valueToJSON(valueType string, value string) (*json.RawMessage, error) {
	switch valueType {
	case valueTypeString, valueTypeEnum:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		raw := json.RawMessage(encoded)
		return &raw, nil
	case valueTypeNumber:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return nil, nil
		}
		if _, err := strconv.ParseFloat(trimmed, 64); err != nil {
			return nil, err
		}
		raw := json.RawMessage(trimmed)
		return &raw, nil
	case valueTypeBool:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return nil, nil
		}
		val, err := strconv.ParseBool(trimmed)
		if err != nil {
			return nil, err
		}
		raw := json.RawMessage(strconv.FormatBool(val))
		return &raw, nil
	case valueTypeJSON:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return nil, nil
		}
		if !json.Valid([]byte(trimmed)) {
			return nil, fmt.Errorf("invalid json")
		}
		raw := json.RawMessage(trimmed)
		return &raw, nil
	default:
		return nil, fmt.Errorf("unsupported valueType: %s", valueType)
	}
}

func defaultValueToJSON(valueType string, value *string) (*json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" && valueType != valueTypeString && valueType != valueTypeEnum {
		return nil, nil
	}
	raw, err := valueToJSON(valueType, *value)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func normalizeRaw(raw json.RawMessage, defaultValue json.RawMessage) json.RawMessage {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return defaultValue
	}
	return trimmed
}
