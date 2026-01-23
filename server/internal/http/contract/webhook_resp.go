package contract

import "time"

// WebhookResp Webhook 响应。
type WebhookResp struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	URL             string            `json:"url"`
	Events          []string          `json:"events"`
	Headers         map[string]string `json:"headers"`
	PayloadTemplate string            `json:"payloadTemplate"`
	IsEnabled       bool              `json:"isEnabled"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

// WebhookHistoryResp Webhook 历史响应。
type WebhookHistoryResp struct {
	ID              int64             `json:"id"`
	WebhookID       int64             `json:"webhookId"`
	EventName       string            `json:"eventName"`
	RequestURL      string            `json:"requestUrl"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	RequestBody     string            `json:"requestBody"`
	ResponseStatus  int               `json:"responseStatus"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	ResponseBody    string            `json:"responseBody,omitempty"`
	ErrorMessage    string            `json:"errorMessage,omitempty"`
	IsTest          bool              `json:"isTest"`
	CreatedAt       time.Time         `json:"createdAt"`
}

// WebhookHistoryListResp Webhook 历史列表响应。
type WebhookHistoryListResp struct {
	Items []WebhookHistoryResp `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// WebhookEventListResp Webhook 事件列表响应。
type WebhookEventListResp struct {
	Events []string `json:"events"`
}
