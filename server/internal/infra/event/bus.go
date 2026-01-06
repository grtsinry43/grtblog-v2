package event

import (
	"context"
	"sync"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
)

// InMemoryBus is a synchronous event bus for in-process handlers.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]appEvent.Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]appEvent.Handler),
	}
}

func (b *InMemoryBus) Subscribe(name string, handler appEvent.Handler) {
	if handler == nil || name == "" {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], handler)
}

func (b *InMemoryBus) Publish(ctx context.Context, event appEvent.Event) error {
	if event == nil {
		return nil
	}
	b.mu.RLock()
	handlers := append([]appEvent.Handler(nil), b.handlers[event.Name()]...)
	b.mu.RUnlock()

	var firstErr error
	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
