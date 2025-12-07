package model

import "time"

type SysConfig struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	ConfigKey string    `gorm:"column:config_key;size:45;not null"`
	Value     string    `gorm:"column:value;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
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
