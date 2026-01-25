package contract

import "time"

// MomentResp 手记响应。
type MomentResp struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	AISummary   *string      `json:"aiSummary,omitempty"`
	TOC         []TOCNode    `json:"toc,omitempty"`
	Content     string       `json:"content"`
	ContentHash string       `json:"contentHash"`
	AuthorID    int64        `json:"authorId"`
	Image       []string     `json:"image,omitempty"`
	ColumnID    *int64       `json:"columnId,omitempty"`
	CommentID   *int64       `json:"commentAreaId,omitempty"`
	ShortURL    string       `json:"shortUrl"`
	IsPublished bool         `json:"isPublished"`
	IsTop       bool         `json:"isTop"`
	IsHot       bool         `json:"isHot"`
	IsOriginal  bool         `json:"isOriginal"`
	Topics      []TagResp    `json:"topics,omitempty"`
	Metrics     *MetricsResp `json:"metrics,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// MomentListItemResp 手记列表项响应。
type MomentListItemResp struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	ShortURL       string    `json:"shortUrl"`
	AuthorName     string    `json:"authorName,omitempty"`
	Summary        string    `json:"summary"`
	Avatar         string    `json:"avatar,omitempty"`
	Image          []string  `json:"image,omitempty"`
	Views          int64     `json:"views"`
	ColumnName     string    `json:"columnName,omitempty"`
	ColumnShortURL string    `json:"columnShortUrl,omitempty"`
	CommentID      *int64    `json:"commentAreaId,omitempty"`
	Topics         []string  `json:"topics"`
	Likes          int       `json:"likes"`
	Comments       int       `json:"comments"`
	IsTop          bool      `json:"isTop"`
	IsHot          bool      `json:"isHot"`
	IsOriginal     bool      `json:"isOriginal"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// MomentListResp 手记列表响应。
type MomentListResp struct {
	Items []MomentListItemResp `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// MomentContentPayload 手记内容推送数据。
type MomentContentPayload struct {
	ContentHash string    `json:"contentHash"`
	Title       string    `json:"title,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	TOC         []TOCNode `json:"toc"`
	Content     string    `json:"content,omitempty"`
}

// CheckMomentLatestResp 手记版本校验响应。
type CheckMomentLatestResp struct {
	Latest bool `json:"latest"`
	MomentContentPayload
}
