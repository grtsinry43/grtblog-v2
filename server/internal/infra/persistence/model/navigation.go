package model

import (
	"time"

	"gorm.io/gorm"
)

type NavMenu struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name;size:45;not null"`
	URL       string         `gorm:"column:url;size:255;not null"`
	Sort      int            `gorm:"column:sort;not null"`
	ParentID  *int64         `gorm:"column:parent_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (NavMenu) TableName() string { return "nav_menu" }

type FooterSection struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Title     string    `gorm:"column:title;size:255;not null"`
	Links     []byte    `gorm:"column:links;type:jsonb;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (FooterSection) TableName() string { return "footer_section" }
