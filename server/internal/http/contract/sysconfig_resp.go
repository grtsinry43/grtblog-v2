package contract

import (
	"encoding/json"
	"time"
)

// SysConfigTreeResp 系统配置树响应。
type SysConfigTreeResp struct {
	Groups []SysConfigGroupResp `json:"groups"`
	Items  []SysConfigItemResp  `json:"items,omitempty"`
}

// SysConfigGroupResp 配置分组信息。
type SysConfigGroupResp struct {
	Key      string               `json:"key"`
	Path     string               `json:"path"`
	Label    string               `json:"label"`
	Children []SysConfigGroupResp `json:"children,omitempty"`
	Items    []SysConfigItemResp  `json:"items,omitempty"`
}

// SysConfigItemResp 系统配置项响应。
type SysConfigItemResp struct {
	Key          string           `json:"key"`
	GroupPath    string           `json:"groupPath,omitempty"`
	Label        string           `json:"label,omitempty"`
	Description  string           `json:"description,omitempty"`
	ValueType    string           `json:"valueType"`
	EnumOptions  json.RawMessage  `json:"enumOptions"`
	DefaultValue *json.RawMessage `json:"defaultValue,omitempty"`
	VisibleWhen  json.RawMessage  `json:"visibleWhen"`
	Sort         int              `json:"sort"`
	Meta         json.RawMessage  `json:"meta"`
	IsSensitive  bool             `json:"isSensitive"`
	Value        *json.RawMessage `json:"value,omitempty"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}
