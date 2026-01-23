package contract

import "encoding/json"

// SysConfigBatchUpdateReq 系统配置批量更新请求。
type SysConfigBatchUpdateReq struct {
	Items []SysConfigUpdateItem `json:"items" validate:"required"`
}

// SysConfigUpdateItem 单条配置更新。
type SysConfigUpdateItem struct {
	Key          string           `json:"key" validate:"required,max=45"`
	Value        *json.RawMessage `json:"value,omitempty"`
	IsSensitive  *bool            `json:"isSensitive,omitempty"`
	GroupPath    *string          `json:"groupPath,omitempty"`
	Label        *string          `json:"label,omitempty"`
	Description  *string          `json:"description,omitempty"`
	ValueType    *string          `json:"valueType,omitempty"`
	EnumOptions  *json.RawMessage `json:"enumOptions,omitempty"`
	DefaultValue *json.RawMessage `json:"defaultValue,omitempty"`
	VisibleWhen  *json.RawMessage `json:"visibleWhen,omitempty"`
	Sort         *int             `json:"sort,omitempty"`
	Meta         *json.RawMessage `json:"meta,omitempty"`
}
