package webhook

import (
	"context"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	domainwebhook "github.com/grtsinry43/grtblog-v2/server/internal/domain/webhook"
)

type dispatchTask struct {
	hook      *domainwebhook.Webhook
	eventName string
	event     appEvent.Event
}

type Dispatcher struct {
	repo   domainwebhook.Repository
	sender *Sender
	queue  chan dispatchTask
}

func NewDispatcher(repo domainwebhook.Repository, sender *Sender, workers int, queueSize int) *Dispatcher {
	if workers <= 0 {
		workers = 1
	}
	if queueSize <= 0 {
		queueSize = 1
	}
	d := &Dispatcher{
		repo:   repo,
		sender: sender,
		queue:  make(chan dispatchTask, queueSize),
	}
	for i := 0; i < workers; i++ {
		go d.worker()
	}
	return d
}

func (d *Dispatcher) Handle(ctx context.Context, event appEvent.Event) error {
	if event == nil {
		return nil
	}
	hooks, err := d.repo.ListEnabledByEvent(ctx, event.Name())
	if err != nil {
		return err
	}
	if len(hooks) == 0 {
		return nil
	}

	for _, hook := range hooks {
		task := dispatchTask{
			hook:      hook,
			eventName: event.Name(),
			event:     event,
		}
		select {
		case d.queue <- task:
		default:
			go d.sender.RecordHistoryFromEvent(context.Background(), hook, event.Name(), event, "queue full", false)
		}
	}
	return nil
}

func (d *Dispatcher) worker() {
	for task := range d.queue {
		_ = d.sender.Send(context.Background(), task.hook, task.eventName, task.event, false)
	}
}

func RegisterSubscribers(bus appEvent.Bus, handler appEvent.Handler) {
	if bus == nil || handler == nil {
		return
	}
	for _, name := range AvailableEventNames {
		bus.Subscribe(name, handler)
	}
}
