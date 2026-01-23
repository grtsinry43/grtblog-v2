package webhook

import "time"

type Webhook struct {
	ID              int64
	Name            string
	URL             string
	Events          []string
	Headers         map[string]string
	PayloadTemplate string
	IsEnabled       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type DeliveryHistory struct {
	ID              int64
	WebhookID       int64
	EventName       string
	RequestURL      string
	RequestHeaders  map[string]string
	RequestBody     string
	ResponseStatus  int
	ResponseHeaders map[string]string
	ResponseBody    string
	ErrorMessage    string
	IsTest          bool
	CreatedAt       time.Time
}
