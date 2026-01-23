package contract

import (
	"encoding/json"
	"strings"
	"time"
)

// CreateMomentReq 创建手记请求。
type CreateMomentReq struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content" validate:"required"`
	Image       []string   `json:"image,omitempty"`
	ColumnID    *int64     `json:"columnId,omitempty"`
	TopicIDs    []int64    `json:"topicIds,omitempty"`
	ShortURL    *string    `json:"shortUrl"`
	IsPublished bool       `json:"isPublished" validate:"required"`
	IsTop       bool       `json:"isTop"`
	IsHot       bool       `json:"isHot"`
	IsOriginal  bool       `json:"isOriginal"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"` // 可以自定义发布时间
}

type createMomentReqJSON struct {
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content"`
	Image       []string `json:"image"`
	ColumnID    *int64   `json:"columnId"`
	TopicIDs    []int64  `json:"topicIds"`
	ShortURL    *string  `json:"shortUrl"`
	IsPublished bool     `json:"isPublished"`
	IsTop       bool     `json:"isTop"`
	IsHot       bool     `json:"isHot"`
	IsOriginal  bool     `json:"isOriginal"`
	CreatedAt   *string  `json:"createdAt"`
}

func (r *CreateMomentReq) UnmarshalJSON(data []byte) error {
	var aux createMomentReqJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Title = aux.Title
	r.Summary = aux.Summary
	r.Content = aux.Content
	r.Image = aux.Image
	r.ColumnID = aux.ColumnID
	r.TopicIDs = aux.TopicIDs
	r.ShortURL = aux.ShortURL
	r.IsPublished = aux.IsPublished
	r.IsTop = aux.IsTop
	r.IsHot = aux.IsHot
	r.IsOriginal = aux.IsOriginal

	if aux.CreatedAt == nil {
		r.CreatedAt = nil
		return nil
	}
	if strings.TrimSpace(*aux.CreatedAt) == "" {
		now := time.Now()
		r.CreatedAt = &now
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, *aux.CreatedAt)
	if err != nil {
		return err
	}
	r.CreatedAt = &parsed
	return nil
}

// UpdateMomentReq 更新手记请求。
type UpdateMomentReq struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content" validate:"required"`
	Image       []string `json:"image,omitempty"`
	ColumnID    *int64   `json:"columnId,omitempty"`
	TopicIDs    []int64  `json:"topicIds,omitempty"`
	ShortURL    string   `json:"shortUrl" validate:"required"`
	IsPublished bool     `json:"isPublished"`
	IsTop       bool     `json:"isTop"`
	IsHot       bool     `json:"isHot"`
	IsOriginal  bool     `json:"isOriginal"`
}

// ListMomentsReq 手记列表查询请求。
type ListMomentsReq struct {
	Page      int     `json:"page" validate:"min=1"`
	PageSize  int     `json:"pageSize" validate:"min=1,max=100"`
	ColumnID  *int64  `json:"columnId,omitempty"`
	TopicID   *int64  `json:"topicId,omitempty"`
	AuthorID  *int64  `json:"authorId,omitempty"`
	Published *bool   `json:"published,omitempty"`
	Search    *string `json:"search,omitempty"`
}

// CheckMomentLatestReq 手记版本校验请求。
type CheckMomentLatestReq struct {
	Hash string `json:"hash"`
}
