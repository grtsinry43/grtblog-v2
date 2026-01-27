package model

import (
	"time"
)

type Thinking struct {
	ID        int64           `gorm:"column:id;primaryKey;autoIncrement"`
	CommentID int64           `gorm:"column:comment_id;not null"`
	Content   string          `gorm:"column:content;type:text;not null"`
	AuthorID  int64           `gorm:"column:author_id"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	Metrics   ThinkingMetrics `gorm:"foreignKey:ThinkingID;references:ID"`
}

func (Thinking) TableName() string { return "thinking" }

type ThinkingMetrics struct {
	ThinkingID int64     `gorm:"column:thinking_id;primaryKey"`
	Views      int64     `gorm:"column:views;not null;default:0"`
	Likes      int       `gorm:"column:likes;not null;default:0"`
	Comments   int       `gorm:"column:comments;not null;default:0"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ThinkingMetrics) TableName() string { return "thinking_metrics" }
