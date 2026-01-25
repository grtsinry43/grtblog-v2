package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FriendLink struct {
	ID               int64          `gorm:"column:id;primaryKey"`
	Name             string         `gorm:"column:name;size:255;not null"`
	URL              string         `gorm:"column:url;size:255;not null"`
	Logo             string         `gorm:"column:logo;size:255"`
	Description      string         `gorm:"column:description"`
	RSSURL           string         `gorm:"column:rss_url;size:255"`
	Kind             string         `gorm:"column:kind;size:20;not null"`
	SyncMode         string         `gorm:"column:sync_mode;size:20;not null"`
	InstanceID       *int64         `gorm:"column:instance_id"`
	LastSyncAt       *time.Time     `gorm:"column:last_sync_at"`
	LastSyncStatus   *string        `gorm:"column:last_sync_status;size:20"`
	SyncInterval     *int           `gorm:"column:sync_interval"`
	TotalPostsCached int            `gorm:"column:total_posts_cached;not null"`
	UserID           *int64         `gorm:"column:user_id"`
	IsActive         bool           `gorm:"column:is_active"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (FriendLink) TableName() string { return "friend_link" }

type FriendLinkApplication struct {
	ID                int64          `gorm:"column:id;primaryKey"`
	Name              *string        `gorm:"column:name;size:255"`
	URL               string         `gorm:"column:url;size:255;not null"`
	Logo              *string        `gorm:"column:logo;size:255"`
	Description       *string        `gorm:"column:description"`
	ApplyChannel      string         `gorm:"column:apply_channel;size:20;not null"`
	RequestedSyncMode string         `gorm:"column:requested_sync_mode;size:20;not null"`
	RSSURL            *string        `gorm:"column:rss_url;size:255"`
	InstanceURL       *string        `gorm:"column:instance_url;size:255"`
	Manifest          datatypes.JSON `gorm:"column:manifest;type:jsonb"`
	SignatureKeyID    *string        `gorm:"column:signature_key_id;type:text"`
	SignatureVerified bool           `gorm:"column:signature_verified;not null"`
	UserID            *int64         `gorm:"column:user_id"`
	Message           *string        `gorm:"column:message"`
	Status            string         `gorm:"column:status;size:20"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (FriendLinkApplication) TableName() string { return "friend_link_applications" }

type GlobalNotification struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	Content    string    `gorm:"column:content;type:text;not null"`
	PublishAt  time.Time `gorm:"column:publish_at;not null"`
	ExpireAt   time.Time `gorm:"column:expire_at;not null"`
	AllowClose bool      `gorm:"column:allow_close"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (GlobalNotification) TableName() string { return "global_notification" }
