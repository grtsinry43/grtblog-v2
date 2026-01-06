package contract

import "time"

// ArticleResp 文章响应。
type ArticleResp struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	AISummary   *string      `json:"aiSummary,omitempty"`
	LeadIn      *string      `json:"leadIn,omitempty"`
	TOC         []TOCNode    `json:"toc,omitempty"`
	Content     string       `json:"content"`
	AuthorID    int64        `json:"authorId"`
	Cover       *string      `json:"cover,omitempty"`
	CategoryID  *int64       `json:"categoryId,omitempty"`
	ShortURL    string       `json:"shortUrl"`
	IsPublished bool         `json:"isPublished"`
	IsTop       bool         `json:"isTop"`
	IsHot       bool         `json:"isHot"`
	IsOriginal  bool         `json:"isOriginal"`
	Tags        []TagResp    `json:"tags,omitempty"`
	Metrics     *MetricsResp `json:"metrics,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// TOCNode 目录节点。
type TOCNode struct {
	Name     string    `json:"name"`
	Anchor   string    `json:"anchor"`
	Children []TOCNode `json:"children,omitempty"`
}

// TagResp 标签响应。
type TagResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// MetricsResp 指标响应。
type MetricsResp struct {
	Views    int64 `json:"views"`
	Likes    int   `json:"likes"`
	Comments int   `json:"comments"`
}

// ArticleListItemResp 文章列表项响应。
type ArticleListItemResp struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	ShortURL         string    `json:"shortUrl"`
	AuthorName       string    `json:"authorName,omitempty"`
	Summary          string    `json:"summary"`
	Avatar           string    `json:"avatar,omitempty"`
	Cover            string    `json:"cover,omitempty"`
	Views            int64     `json:"views"`
	CategoryName     string    `json:"categoryName,omitempty"`
	CategoryShortURL string    `json:"categoryShortUrl,omitempty"`
	Tags             []string  `json:"tags"`
	Likes            int       `json:"likes"`
	Comments         int       `json:"comments"`
	IsTop            bool      `json:"isTop"`
	IsHot            bool      `json:"isHot"`
	IsOriginal       bool      `json:"isOriginal"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ArticleListResp 文章列表响应。
type ArticleListResp struct {
	Items []ArticleListItemResp `json:"items"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}
