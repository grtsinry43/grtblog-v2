package article

import "time"

// CreateArticleCmd 创建文章命令。
type CreateArticleCmd struct {
	Title       string
	Summary     string
	LeadIn      *string
	Content     string
	Cover       *string
	CategoryID  *int64
	TagIDs      []int64
	ShortURL    *string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
	CreatedAt   *time.Time // 可选：因为可能会有自定义发布时间的需求
}

// UpdateArticleCmd 更新文章命令。
type UpdateArticleCmd struct {
	ID          int64
	Title       string
	Summary     string
	LeadIn      *string
	Content     string
	Cover       *string
	CategoryID  *int64
	TagIDs      []int64
	ShortURL    string
	IsPublished bool
	IsTop       bool
	IsHot       bool
	IsOriginal  bool
}
