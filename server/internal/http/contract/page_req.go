package contract

import (
	"encoding/json"
	"strings"
	"time"
)

// CreatePageReq 创建页面请求。
type CreatePageReq struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Description *string    `json:"description,omitempty"`
	Content     string     `json:"content" validate:"required"`
	ShortURL    *string    `json:"shortUrl"`
	IsEnabled   bool       `json:"isEnabled"`
	IsBuiltin   bool       `json:"isBuiltin"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type createPageReqJSON struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Content     string  `json:"content"`
	ShortURL    *string `json:"shortUrl"`
	IsEnabled   bool    `json:"isEnabled"`
	IsBuiltin   bool    `json:"isBuiltin"`
	CreatedAt   *string `json:"createdAt"`
}

func (r *CreatePageReq) UnmarshalJSON(data []byte) error {
	var aux createPageReqJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Title = aux.Title
	r.Description = aux.Description
	r.Content = aux.Content
	r.ShortURL = aux.ShortURL
	r.IsEnabled = aux.IsEnabled
	r.IsBuiltin = aux.IsBuiltin

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

// UpdatePageReq 更新页面请求。
type UpdatePageReq struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
	Content     string  `json:"content" validate:"required"`
	ShortURL    string  `json:"shortUrl" validate:"required"`
	IsEnabled   bool    `json:"isEnabled"`
	IsBuiltin   bool    `json:"isBuiltin"`
}

// ListPagesReq 页面列表查询请求。
type ListPagesReq struct {
	Page     int     `json:"page" validate:"min=1"`
	PageSize int     `json:"pageSize" validate:"min=1,max=100"`
	Enabled  *bool   `json:"enabled,omitempty"`
	Builtin  *bool   `json:"builtin,omitempty"`
	Search   *string `json:"search,omitempty"`
}

// CheckPageLatestReq 页面版本校验请求。
type CheckPageLatestReq struct {
	Hash string `json:"hash"`
}
