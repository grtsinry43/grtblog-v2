package webhook

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domainwebhook "github.com/grtsinry43/grtblog-v2/server/internal/domain/webhook"
)

type Service struct {
	repo   domainwebhook.Repository
	sender *Sender
}

func NewService(repo domainwebhook.Repository, sender *Sender) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
	}
}

func (s *Service) ListEvents() []string {
	return append([]string(nil), AvailableEventNames...)
}

func (s *Service) List(ctx context.Context) ([]*domainwebhook.Webhook, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domainwebhook.Webhook, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, hook *domainwebhook.Webhook) error {
	if hook == nil {
		return errors.New("webhook is nil")
	}
	if err := normalizeAndValidate(hook); err != nil {
		return err
	}
	return s.repo.Create(ctx, hook)
}

func (s *Service) Update(ctx context.Context, hook *domainwebhook.Webhook) error {
	if hook == nil {
		return errors.New("webhook is nil")
	}
	if err := normalizeAndValidate(hook); err != nil {
		return err
	}
	return s.repo.Update(ctx, hook)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListHistory(ctx context.Context, options domainwebhook.DeliveryHistoryListOptions) ([]*domainwebhook.DeliveryHistory, int64, error) {
	return s.repo.ListHistory(ctx, options)
}

func (s *Service) Test(ctx context.Context, id int64, eventName *string) error {
	hook, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	selected, err := pickEventName(hook, eventName)
	if err != nil {
		return err
	}
	event, err := SampleEvent(selected)
	if err != nil {
		return err
	}
	if err := s.sender.Send(ctx, hook, selected, event, true); err != nil {
		return fmt.Errorf("%w: %v", domainwebhook.ErrWebhookDeliveryFailed, err)
	}
	return nil
}

func (s *Service) Replay(ctx context.Context, historyID int64) error {
	history, err := s.repo.GetHistoryByID(ctx, historyID)
	if err != nil {
		return err
	}

	hook, err := s.repo.GetByID(ctx, history.WebhookID)
	if err != nil {
		if !errors.Is(err, domainwebhook.ErrWebhookNotFound) {
			return err
		}
		hook = &domainwebhook.Webhook{
			ID:  history.WebhookID,
			URL: history.RequestURL,
		}
	}
	if hook.URL == "" {
		hook.URL = history.RequestURL
	}
	headers := history.RequestHeaders
	if headers == nil {
		headers = map[string]string{}
	}
	if err := s.sender.SendRaw(ctx, hook, history.EventName, history.RequestBody, headers, history.IsTest); err != nil {
		return fmt.Errorf("%w: %v", domainwebhook.ErrWebhookDeliveryFailed, err)
	}
	return nil
}

func normalizeAndValidate(hook *domainwebhook.Webhook) error {
	hook.Name = strings.TrimSpace(hook.Name)
	hook.URL = strings.TrimSpace(hook.URL)
	if hook.PayloadTemplate != "" {
		hook.PayloadTemplate = strings.TrimSpace(hook.PayloadTemplate)
	}
	if hook.Headers == nil {
		hook.Headers = map[string]string{}
	}
	if len(hook.Events) == 0 {
		return domainwebhook.ErrWebhookInvalidEvents
	}
	unique := make(map[string]struct{}, len(hook.Events))
	filtered := make([]string, 0, len(hook.Events))
	for _, name := range hook.Events {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if !IsValidEventName(name) {
			return domainwebhook.ErrWebhookInvalidEvents
		}
		if _, ok := unique[name]; ok {
			continue
		}
		unique[name] = struct{}{}
		filtered = append(filtered, name)
	}
	if len(filtered) == 0 {
		return domainwebhook.ErrWebhookInvalidEvents
	}
	hook.Events = filtered
	if strings.TrimSpace(hook.PayloadTemplate) == "" {
		hook.PayloadTemplate = defaultPayloadTemplate
	}
	return nil
}

func pickEventName(hook *domainwebhook.Webhook, requested *string) (string, error) {
	if requested != nil {
		name := strings.TrimSpace(*requested)
		if name == "" {
			return "", domainwebhook.ErrWebhookInvalidEvents
		}
		if !IsValidEventName(name) {
			return "", domainwebhook.ErrWebhookInvalidEvents
		}
		return name, nil
	}
	if len(hook.Events) == 0 {
		return "", domainwebhook.ErrWebhookInvalidEvents
	}
	return hook.Events[0], nil
}
