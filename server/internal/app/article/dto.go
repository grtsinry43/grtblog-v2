package article

import "time"

// CreateArticleCommand 创建文章命令
type CreateArticleCommand struct {
	Title       string         `json:"title" validate:"required,max=255"`
	Summary     string         `json:"summary" validate:"required"`
	LeadIn      *string        `json:"leadIn,omitempty"`
	Content     string         `json:"content" validate:"required"`
	Cover       *string        `json:"cover,omitempty"`
	CategoryID  *int64         `json:"categoryId,omitempty"`
	TagIDs      []int64        `json:"tagIds,omitempty"`
	ShortURL    *string        `json:"shortUrl,omitempty"`
	IsPublished bool           `json:"isPublished"`
	IsTop       bool           `json:"isTop"`
	IsHot       bool           `json:"isHot"`
	IsOriginal  bool           `json:"isOriginal"`
	TOC         map[string]any `json:"toc,omitempty"`
	CreatedAt   *time.Time     `json:"createdAt,omitempty"` // 可选：因为可能会有自定义发布时间的需求
}

// UpdateArticleCommand 更新文章命令
type UpdateArticleCommand struct {
	ID          int64          `json:"id" validate:"required"`
	Title       string         `json:"title" validate:"required,max=255"`
	Summary     string         `json:"summary" validate:"required"`
	LeadIn      *string        `json:"leadIn,omitempty"`
	Content     string         `json:"content" validate:"required"`
	Cover       *string        `json:"cover,omitempty"`
	CategoryID  *int64         `json:"categoryId,omitempty"`
	TagIDs      []int64        `json:"tagIds,omitempty"`
	ShortURL    string         `json:"shortUrl" validate:"required"`
	IsPublished bool           `json:"isPublished"`
	IsTop       bool           `json:"isTop"`
	IsHot       bool           `json:"isHot"`
	IsOriginal  bool           `json:"isOriginal"`
	TOC         map[string]any `json:"toc,omitempty"`
}

// ViewArticleResponse 文章响应DTO
type ViewArticleResponse struct {
	ID          int64            `json:"id"`
	Title       string           `json:"title"`
	Summary     string           `json:"summary"`
	AISummary   *string          `json:"aiSummary,omitempty"`
	LeadIn      *string          `json:"leadIn,omitempty"`
	TOC         map[string]any   `json:"toc,omitempty"`
	Content     string           `json:"content"`
	AuthorID    int64            `json:"authorId"`
	Cover       *string          `json:"cover,omitempty"`
	CategoryID  *int64           `json:"categoryId,omitempty"`
	ShortURL    string           `json:"shortUrl"`
	IsPublished bool             `json:"isPublished"`
	IsTop       bool             `json:"isTop"`
	IsHot       bool             `json:"isHot"`
	IsOriginal  bool             `json:"isOriginal"`
	Tags        []TagResponse    `json:"tags,omitempty"`
	Metrics     *MetricsResponse `json:"metrics,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// TagResponse 标签响应DTO
type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// MetricsResponse 指标响应DTO
type MetricsResponse struct {
	Views    int64 `json:"views"`
	Likes    int   `json:"likes"`
	Comments int   `json:"comments"`
}

// ListArticleResponse 文章列表响应DTO
type ListArticleResponse struct {
	Items []ViewArticleResponse `json:"items"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// ListArticlesQuery 公开的文章列表查询选项
type ListArticlesQuery struct {
	Page       int     `json:"page" validate:"min=1"`
	PageSize   int     `json:"pageSize" validate:"min=1,max=100"`
	CategoryID *int64  `json:"categoryId,omitempty"`
	TagID      *int64  `json:"tagId,omitempty"`
	AuthorID   *int64  `json:"authorId,omitempty"`
	Published  *bool   `json:"published,omitempty"` // 仅管理员可用
	Search     *string `json:"search,omitempty"`
}
