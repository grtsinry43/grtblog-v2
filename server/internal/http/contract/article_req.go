package contract

import "time"

// CreateArticleReq 创建文章请求。
type CreateArticleReq struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Summary     string     `json:"summary" validate:"required"`
	LeadIn      *string    `json:"leadIn,omitempty"`
	Content     string     `json:"content" validate:"required"`
	Cover       *string    `json:"cover,omitempty"`
	CategoryID  *int64     `json:"categoryId,omitempty"`
	TagIDs      []int64    `json:"tagIds,omitempty"`
	ShortURL    *string    `json:"shortUrl"`
	IsPublished bool       `json:"isPublished" validate:"required"`
	IsTop       bool       `json:"isTop"`
	IsHot       bool       `json:"isHot"`
	IsOriginal  bool       `json:"isOriginal"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"` // 可以自定义发布时间
}

// UpdateArticleReq 更新文章请求。
type UpdateArticleReq struct {
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

// ListArticlesReq 文章列表查询请求。
type ListArticlesReq struct {
	Page       int     `json:"page" validate:"min=1"`
	PageSize   int     `json:"pageSize" validate:"min=1,max=100"`
	CategoryID *int64  `json:"categoryId,omitempty"`
	TagID      *int64  `json:"tagId,omitempty"`
	AuthorID   *int64  `json:"authorId,omitempty"`
	Published  *bool   `json:"published,omitempty"`
	Search     *string `json:"search,omitempty"`
}
