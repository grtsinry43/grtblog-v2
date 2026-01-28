package contract

import (
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type NavMenuResp struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Icon      *string       `json:"icon,omitempty"`
	Sort      int           `json:"sort"`
	ParentID  *int64        `json:"parentId,omitempty"`
	Children  []NavMenuResp `json:"children,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// NavMenuListRespEnvelope 用于 swagger 展示。
type NavMenuListRespEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   []NavMenuResp `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// NavMenuDetailRespEnvelope 用于 swagger 展示。
type NavMenuDetailRespEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   NavMenuResp   `json:"data"`
	Meta   response.Meta `json:"meta"`
}
