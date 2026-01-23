package model

import (
	"time"

	"gorm.io/datatypes"
)

type SysConfig struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	ConfigKey    string         `gorm:"column:config_key;size:45;not null"`
	Value        string         `gorm:"column:value;type:text;not null"`
	IsSensitive  bool           `gorm:"column:is_sensitive;not null"`
	GroupPath    string         `gorm:"column:group_path;type:text"`
	Label        string         `gorm:"column:label;size:100"`
	Description  string         `gorm:"column:description;type:text"`
	ValueType    string         `gorm:"column:value_type;size:20;not null;default:string"`
	EnumOptions  datatypes.JSON `gorm:"column:enum_options;type:jsonb"`
	DefaultValue *string        `gorm:"column:default_value;type:text"`
	VisibleWhen  datatypes.JSON `gorm:"column:visible_when;type:jsonb"`
	Sort         int            `gorm:"column:sort;not null;default:0"`
	Meta         datatypes.JSON `gorm:"column:meta;type:jsonb"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (SysConfig) TableName() string { return "sys_config" }

type WebsiteInfo struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	InfoKey   string    `gorm:"column:info_key;size:45;not null"`
	Value     string    `gorm:"column:value;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (WebsiteInfo) TableName() string { return "website_info" }
