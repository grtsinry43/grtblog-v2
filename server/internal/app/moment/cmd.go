package moment

import "time"

// CreateMomentCmd 创建手记命令。
type CreateMomentCmd struct {
	Title       string
	Summary     string
	Content     string
	Image       *string
	ColumnID    *int64
	TopicIDs    []int64
	ShortURL    *string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
	CreatedAt   *time.Time // 可选：因为可能会有自定义发布时间的需求
}

// UpdateMomentCmd 更新手记命令。
type UpdateMomentCmd struct {
	ID          int64
	Title       string
	Summary     string
	Content     string
	Image       *string
	ColumnID    *int64
	TopicIDs    []int64
	ShortURL    string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
}
