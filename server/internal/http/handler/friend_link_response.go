package handler

import (
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FriendLinkApplicationVO struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name,omitempty"`
	URL         string  `json:"url"`
	Logo        *string `json:"logo,omitempty"`
	Description *string `json:"description,omitempty"`
	UserID      *int64  `json:"userId,omitempty"`
	Message     *string `json:"message,omitempty"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func toFriendLinkApplicationVO(app social.FriendLinkApplication) FriendLinkApplicationVO {
	return FriendLinkApplicationVO{
		ID:          app.ID,
		Name:        app.Name,
		URL:         app.URL,
		Logo:        app.Logo,
		Description: app.Description,
		UserID:      app.UserID,
		Message:     app.Message,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt.Format(response.TimeLayout),
		UpdatedAt:   app.UpdatedAt.Format(response.TimeLayout),
	}
}

// FriendLinkApplicationResponse 用于 swagger 展示友链申请操作结果。
type FriendLinkApplicationResponse struct {
	Code   int                     `json:"code"`
	BizErr string                  `json:"bizErr"`
	Msg    string                  `json:"msg"`
	Data   FriendLinkApplicationVO `json:"data"`
	Meta   response.Meta           `json:"meta"`
}
