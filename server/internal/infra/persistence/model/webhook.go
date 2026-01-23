package model

import (
	"time"

	"gorm.io/gorm"
)

type Webhook struct {
	ID              int64          `gorm:"column:id;primaryKey"`
	Name            string         `gorm:"column:name;size:100;not null"`
	URL             string         `gorm:"column:url;size:512;not null"`
	Events          []byte         `gorm:"column:events;type:jsonb;not null"`
	Headers         []byte         `gorm:"column:headers;type:jsonb;not null"`
	PayloadTemplate string         `gorm:"column:payload_template;type:text;not null"`
	IsEnabled       bool           `gorm:"column:is_enabled"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Webhook) TableName() string { return "webhook" }

type WebhookHistory struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	WebhookID       int64     `gorm:"column:webhook_id;not null"`
	EventName       string    `gorm:"column:event_name;size:128;not null"`
	RequestURL      string    `gorm:"column:request_url;size:512;not null"`
	RequestHeaders  []byte    `gorm:"column:request_headers;type:jsonb;not null"`
	RequestBody     string    `gorm:"column:request_body;type:text;not null"`
	ResponseStatus  int       `gorm:"column:response_status"`
	ResponseHeaders []byte    `gorm:"column:response_headers;type:jsonb;not null"`
	ResponseBody    string    `gorm:"column:response_body;type:text"`
	ErrorMessage    string    `gorm:"column:error_message;type:text"`
	IsTest          bool      `gorm:"column:is_test"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (WebhookHistory) TableName() string { return "webhook_history" }
