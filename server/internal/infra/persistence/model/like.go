package model

import "time"

type ContentLike struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	TargetType string    `gorm:"column:target_type;type:like_target_type;not null"`
	TargetID   int64     `gorm:"column:target_id;not null"`
	UserID     *int64    `gorm:"column:user_id"`
	SessionID  string    `gorm:"column:session_id;size:255"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ContentLike) TableName() string { return "content_like" }
