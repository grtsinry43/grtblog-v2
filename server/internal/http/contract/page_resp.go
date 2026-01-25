package contract

import "time"

// PageResp 页面响应。
type PageResp struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description *string      `json:"description,omitempty"`
	AISummary   *string      `json:"aiSummary,omitempty"`
	TOC         []TOCNode    `json:"toc,omitempty"`
	Content     string       `json:"content"`
	ContentHash string       `json:"contentHash"`
	CommentID   *int64       `json:"commentAreaId,omitempty"`
	ShortURL    string       `json:"shortUrl"`
	IsEnabled   bool         `json:"isEnabled"`
	IsBuiltin   bool         `json:"isBuiltin"`
	Metrics     *MetricsResp `json:"metrics,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// PageListItemResp 页面列表项响应。
type PageListItemResp struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	ShortURL    string    `json:"shortUrl"`
	Description *string   `json:"description,omitempty"`
	Views       int64     `json:"views"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	CommentID   *int64    `json:"commentAreaId,omitempty"`
	IsEnabled   bool      `json:"isEnabled"`
	IsBuiltin   bool      `json:"isBuiltin"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// PageListResp 页面列表响应。
type PageListResp struct {
	Items []PageListItemResp `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// PageContentPayload 页面内容推送数据。
type PageContentPayload struct {
	ContentHash string    `json:"contentHash"`
	Title       string    `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	TOC         []TOCNode `json:"toc"`
	Content     string    `json:"content,omitempty"`
}

// CheckPageLatestResp 页面版本校验响应。
type CheckPageLatestResp struct {
	Latest bool `json:"latest"`
	PageContentPayload
}
