package contract

// CreateWebhookReq Webhook 创建请求。
type CreateWebhookReq struct {
	Name            string            `json:"name" validate:"required,max=100"`
	URL             string            `json:"url" validate:"required,max=512"`
	Events          []string          `json:"events" validate:"required"`
	Headers         map[string]string `json:"headers,omitempty"`
	PayloadTemplate string            `json:"payloadTemplate,omitempty"`
	IsEnabled       bool              `json:"isEnabled"`
}

// UpdateWebhookReq Webhook 更新请求。
type UpdateWebhookReq struct {
	Name            string            `json:"name" validate:"required,max=100"`
	URL             string            `json:"url" validate:"required,max=512"`
	Events          []string          `json:"events" validate:"required"`
	Headers         map[string]string `json:"headers,omitempty"`
	PayloadTemplate string            `json:"payloadTemplate,omitempty"`
	IsEnabled       bool              `json:"isEnabled"`
}

// WebhookTestReq Webhook 测试请求。
type WebhookTestReq struct {
	EventName *string `json:"eventName,omitempty"`
}

// WebhookHistoryListReq Webhook 历史列表查询请求。
type WebhookHistoryListReq struct {
	Page      int     `json:"page" validate:"min=1"`
	PageSize  int     `json:"pageSize" validate:"min=1,max=100"`
	WebhookID *int64  `json:"webhookId,omitempty"`
	EventName *string `json:"eventName,omitempty"`
	IsTest    *bool   `json:"isTest,omitempty"`
}
