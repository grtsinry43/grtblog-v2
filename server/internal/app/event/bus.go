package event

import (
	"context"
	"time"
)

// Event represents a domain event published by application services.
type Event interface {
	Name() string
	OccurredAt() time.Time
}

// Handler consumes events by name.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Bus publishes events to subscribed handlers.
type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(name string, handler Handler)
}

// NopBus is a safe default when no event bus is configured.
type NopBus struct{}

func (NopBus) Publish(context.Context, Event) error { return nil }
func (NopBus) Subscribe(string, Handler)            {}
