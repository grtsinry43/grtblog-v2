package model

import (
	"time"

	"gorm.io/gorm"
)

type CommentArea struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	AreaName  string    `gorm:"column:area_name;size:45;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CommentArea) TableName() string { return "comment_area" }

type Comment struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	AreaID    int64          `gorm:"column:area_id;not null"`
	Content   string         `gorm:"column:content;type:text;not null"`
	AuthorID  *int64         `gorm:"column:author_id"`
	NickName  string         `gorm:"column:nick_name;size:45"`
	IP        string         `gorm:"column:ip;size:45"`
	Location  string         `gorm:"column:location;size:45"`
	Platform  string         `gorm:"column:platform;size:45"`
	Browser   string         `gorm:"column:browser;size:45"`
	Email     string         `gorm:"column:email;size:255"`
	Website   string         `gorm:"column:website;size:255"`
	IsOwner   bool           `gorm:"column:is_owner"`
	IsFriend  bool           `gorm:"column:is_friend"`
	IsAuthor  bool           `gorm:"column:is_author"`
	IsViewed  bool           `gorm:"column:is_viewed"`
	IsTop     bool           `gorm:"column:is_top"`
	ParentID  *int64         `gorm:"column:parent_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Comment) TableName() string { return "comment" }
