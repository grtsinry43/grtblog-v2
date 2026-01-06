package contract

import (
	"time"

	"github.com/jinzhu/copier"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// UserResp 描述对外暴露的用户字段（小驼峰）。
type UserResp struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	Avatar    string  `json:"avatar"`
	IsActive  bool    `json:"isActive"`
	IsAdmin   bool    `json:"isAdmin"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

func ToUserResp(u identity.User) UserResp {
	var deleted *string
	if u.DeletedAt != nil {
		val := u.DeletedAt.Format(time.RFC3339)
		deleted = &val
	}
	var resp UserResp
	_ = copier.Copy(&resp, u)
	resp.CreatedAt = u.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = u.UpdatedAt.Format(time.RFC3339)
	resp.DeletedAt = deleted
	return resp
}

// LoginResp 供登录接口返回数据使用。
type LoginResp struct {
	Token string   `json:"token"`
	User  UserResp `json:"user"`
}

// AccessInfoResp 返回当前登录用户的权限信息。
type AccessInfoResp struct {
	User UserResp `json:"user"`
}

// InitStateResp 返回是否需要初始化。
type InitStateResp struct {
	Initialized bool `json:"initialized"`
}

// RegisterRespEnvelope 仅用于 swagger 展示。
type RegisterRespEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   UserResp      `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// LoginRespEnvelope 仅用于 swagger 展示。
type LoginRespEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   LoginResp     `json:"data"`
	Meta   response.Meta `json:"meta"`
}
