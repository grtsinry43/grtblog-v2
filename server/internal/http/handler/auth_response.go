package handler

import (
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// UserResponse 描述对外暴露的用户字段（小驼峰）。
type UserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	Avatar    string  `json:"avatar"`
	IsActive  bool    `json:"isActive"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

func toUserResponse(u identity.User) UserResponse {
	var deleted *string
	if u.DeletedAt != nil {
		val := u.DeletedAt.Format(time.RFC3339)
		deleted = &val
	}
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deleted,
	}
}

// LoginResponse 供登录接口返回数据使用。
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterResponseEnvelope 仅用于 swagger 展示。
type RegisterResponseEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   UserResponse  `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// LoginResponseEnvelope 仅用于 swagger 展示。
type LoginResponseEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   LoginResponse `json:"data"`
	Meta   response.Meta `json:"meta"`
}
