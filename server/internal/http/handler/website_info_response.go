package handler

import (
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type WebsiteInfoResponse struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func toWebsiteInfoResponse(info config.WebsiteInfo) WebsiteInfoResponse {
	return WebsiteInfoResponse{
		Key:       info.Key,
		Value:     info.Value,
		CreatedAt: info.CreatedAt.Format(response.TimeLayout),
		UpdatedAt: info.UpdatedAt.Format(response.TimeLayout),
	}
}

func toWebsiteInfoListResponse(items []config.WebsiteInfo) []WebsiteInfoResponse {
	result := make([]WebsiteInfoResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toWebsiteInfoResponse(item))
	}
	return result
}

// 用于 swagger 展示。
type WebsiteInfoListEnvelope struct {
	Code   int                   `json:"code"`
	BizErr string                `json:"bizErr"`
	Msg    string                `json:"msg"`
	Data   []WebsiteInfoResponse `json:"data"`
	Meta   response.Meta         `json:"meta"`
}

type WebsiteInfoDetailEnvelope struct {
	Code   int                 `json:"code"`
	BizErr string              `json:"bizErr"`
	Msg    string              `json:"msg"`
	Data   WebsiteInfoResponse `json:"data"`
	Meta   response.Meta       `json:"meta"`
}

type GenericMessageEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   interface{}   `json:"data"`
	Meta   response.Meta `json:"meta"`
}
