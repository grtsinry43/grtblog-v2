package contract

import "time"

type ThinkingResp struct {
	ID        int64           `json:"id"`
	CommentID int64           `json:"commentId"`
	Content   string          `json:"content"`
	Author    string          `json:"author"`
	Metrics   ThinkingMetrics `json:"metrics"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type ThinkingMetrics struct {
	Views    int64 `json:"views"`
	Likes    int   `json:"likes"`
	Comments int   `json:"comments"`
}

type ListThinkingResp struct {
	Items []*ThinkingResp `json:"items"`
	Total int64           `json:"total"`
}
