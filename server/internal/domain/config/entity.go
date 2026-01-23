package config

import (
	"encoding/json"
	"time"
)

type SysConfig struct {
	ID           int64
	Key          string
	Value        string
	IsSensitive  bool
	GroupPath    string
	Label        string
	Description  string
	ValueType    string
	EnumOptions  json.RawMessage
	DefaultValue *string
	VisibleWhen  json.RawMessage
	Sort         int
	Meta         json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WebsiteInfo struct {
	ID        int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
