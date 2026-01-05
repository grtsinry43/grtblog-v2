package contract

import (
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type WebsiteInfoResp struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func ToWebsiteInfoResp(info config.WebsiteInfo) WebsiteInfoResp {
	return WebsiteInfoResp{
		Key:       info.Key,
		Value:     info.Value,
		CreatedAt: info.CreatedAt.Format(response.TimeLayout),
		UpdatedAt: info.UpdatedAt.Format(response.TimeLayout),
	}
}

func ToWebsiteInfoListResp(items []config.WebsiteInfo) []WebsiteInfoResp {
	result := make([]WebsiteInfoResp, 0, len(items))
	for _, item := range items {
		result = append(result, ToWebsiteInfoResp(item))
	}
	return result
}

// 用于 swagger 展示。
type WebsiteInfoListRespEnvelope struct {
	Code   int               `json:"code"`
	BizErr string            `json:"bizErr"`
	Msg    string            `json:"msg"`
	Data   []WebsiteInfoResp `json:"data"`
	Meta   response.Meta     `json:"meta"`
}

type WebsiteInfoDetailRespEnvelope struct {
	Code   int             `json:"code"`
	BizErr string          `json:"bizErr"`
	Msg    string          `json:"msg"`
	Data   WebsiteInfoResp `json:"data"`
	Meta   response.Meta   `json:"meta"`
}

type GenericMessageEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   interface{}   `json:"data"`
	Meta   response.Meta `json:"meta"`
}
