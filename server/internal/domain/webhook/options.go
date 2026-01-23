package webhook

type DeliveryHistoryListOptions struct {
	Page      int
	PageSize  int
	WebhookID *int64
	EventName *string
	IsTest    *bool
}
