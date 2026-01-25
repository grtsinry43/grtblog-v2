package contract

import (
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/social"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type FriendLinkApplicationResp struct {
	ID                int64   `json:"id"`
	Name              *string `json:"name,omitempty"`
	URL               string  `json:"url"`
	Logo              *string `json:"logo,omitempty"`
	Description       *string `json:"description,omitempty"`
	ApplyChannel      string  `json:"applyChannel"`
	RequestedSyncMode string  `json:"requestedSyncMode"`
	RSSURL            *string `json:"rssUrl,omitempty"`
	InstanceURL       *string `json:"instanceUrl,omitempty"`
	SignatureVerified bool    `json:"signatureVerified"`
	UserID            *int64  `json:"userId,omitempty"`
	Message           *string `json:"message,omitempty"`
	Status            string  `json:"status"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

func ToFriendLinkApplicationResp(app social.FriendLinkApplication) FriendLinkApplicationResp {
	return FriendLinkApplicationResp{
		ID:                app.ID,
		Name:              app.Name,
		URL:               app.URL,
		Logo:              app.Logo,
		Description:       app.Description,
		ApplyChannel:      app.ApplyChannel,
		RequestedSyncMode: app.RequestedSyncMode,
		RSSURL:            app.RSSURL,
		InstanceURL:       app.InstanceURL,
		SignatureVerified: app.SignatureVerified,
		UserID:            app.UserID,
		Message:           app.Message,
		Status:            app.Status,
		CreatedAt:         app.CreatedAt.Format(response.TimeLayout),
		UpdatedAt:         app.UpdatedAt.Format(response.TimeLayout),
	}
}

// FriendLinkApplicationRespEnvelope 用于 swagger 展示友链申请操作结果。
type FriendLinkApplicationRespEnvelope struct {
	Code   int                       `json:"code"`
	BizErr string                    `json:"bizErr"`
	Msg    string                    `json:"msg"`
	Data   FriendLinkApplicationResp `json:"data"`
	Meta   response.Meta             `json:"meta"`
}
