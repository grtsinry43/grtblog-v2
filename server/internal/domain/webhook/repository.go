package webhook

import "context"

type Repository interface {
	Create(ctx context.Context, hook *Webhook) error
	Update(ctx context.Context, hook *Webhook) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Webhook, error)
	List(ctx context.Context) ([]*Webhook, error)
	ListEnabledByEvent(ctx context.Context, eventName string) ([]*Webhook, error)

	CreateHistory(ctx context.Context, history *DeliveryHistory) error
	GetHistoryByID(ctx context.Context, id int64) (*DeliveryHistory, error)
	ListHistory(ctx context.Context, options DeliveryHistoryListOptions) ([]*DeliveryHistory, int64, error)
}
