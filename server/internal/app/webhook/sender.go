package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	domainwebhook "github.com/grtsinry43/grtblog-v2/server/internal/domain/webhook"
)

const (
	defaultPayloadTemplate = `{"name":"{{.Name}}","occurredAt":"{{.OccurredAt}}","event":{{ toJSON .Event }}}`
	maxResponseBodyBytes   = 256 * 1024
)

type TemplateData struct {
	Name       string
	OccurredAt time.Time
	Event      any
}

type Sender struct {
	repo    domainwebhook.Repository
	client  *http.Client
	timeout time.Duration
}

func NewSender(repo domainwebhook.Repository, timeout time.Duration) *Sender {
	return &Sender{
		repo:    repo,
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

func (s *Sender) Send(ctx context.Context, hook *domainwebhook.Webhook, eventName string, event appEvent.Event, isTest bool) error {
	if hook == nil {
		return errors.New("webhook is nil")
	}
	data := TemplateData{
		Name:       eventName,
		OccurredAt: event.OccurredAt(),
		Event:      event,
	}
	payload, err := renderTemplate(hook.PayloadTemplate, data)
	if err != nil {
		s.recordHistory(ctx, hook, eventName, payload, hook.Headers, 0, nil, "", err.Error(), isTest)
		return err
	}
	headers, err := renderHeaders(hook.Headers, data)
	if err != nil {
		s.recordHistory(ctx, hook, eventName, payload, hook.Headers, 0, nil, "", err.Error(), isTest)
		return err
	}
	return s.sendRaw(ctx, hook, eventName, payload, headers, isTest)
}

func (s *Sender) SendRaw(ctx context.Context, hook *domainwebhook.Webhook, eventName string, payload string, headers map[string]string, isTest bool) error {
	return s.sendRaw(ctx, hook, eventName, payload, headers, isTest)
}

func (s *Sender) RecordHistoryFromEvent(ctx context.Context, hook *domainwebhook.Webhook, eventName string, event appEvent.Event, reason string, isTest bool) {
	if hook == nil || event == nil {
		return
	}
	data := TemplateData{
		Name:       eventName,
		OccurredAt: event.OccurredAt(),
		Event:      event,
	}
	payload, err := renderTemplate(hook.PayloadTemplate, data)
	if err != nil {
		reason = err.Error()
	}
	headers, err := renderHeaders(hook.Headers, data)
	if err != nil {
		reason = err.Error()
	}
	s.recordHistory(ctx, hook, eventName, payload, headers, 0, nil, "", reason, isTest)
}

func (s *Sender) sendRaw(ctx context.Context, hook *domainwebhook.Webhook, eventName string, payload string, headers map[string]string, isTest bool) error {
	if hook == nil {
		return errors.New("webhook is nil")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hook.URL, strings.NewReader(payload))
	if err != nil {
		s.recordHistory(ctx, hook, eventName, payload, headers, 0, nil, "", err.Error(), isTest)
		return err
	}
	if headers == nil {
		headers = map[string]string{}
	}
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.recordHistory(ctx, hook, eventName, payload, headers, 0, nil, "", err.Error(), isTest)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	body := string(bodyBytes)
	responseHeaders := flattenHeaders(resp.Header)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		s.recordHistory(ctx, hook, eventName, payload, headers, resp.StatusCode, responseHeaders, body, err.Error(), isTest)
		return err
	}

	s.recordHistory(ctx, hook, eventName, payload, headers, resp.StatusCode, responseHeaders, body, "", isTest)
	return nil
}

func (s *Sender) recordHistory(ctx context.Context, hook *domainwebhook.Webhook, eventName string, payload string, headers map[string]string, status int, responseHeaders map[string]string, responseBody string, errMsg string, isTest bool) {
	if hook == nil {
		return
	}
	if headers == nil {
		headers = map[string]string{}
	}
	if responseHeaders == nil {
		responseHeaders = map[string]string{}
	}
	history := &domainwebhook.DeliveryHistory{
		WebhookID:       hook.ID,
		EventName:       eventName,
		RequestURL:      hook.URL,
		RequestHeaders:  headers,
		RequestBody:     payload,
		ResponseStatus:  status,
		ResponseHeaders: responseHeaders,
		ResponseBody:    responseBody,
		ErrorMessage:    errMsg,
		IsTest:          isTest,
	}
	_ = s.repo.CreateHistory(ctx, history)
}

func flattenHeaders(headers http.Header) map[string]string {
	if len(headers) == 0 {
		return map[string]string{}
	}
	result := make(map[string]string, len(headers))
	for key, values := range headers {
		result[key] = strings.Join(values, ",")
	}
	return result
}

func renderTemplate(tmpl string, data TemplateData) (string, error) {
	content := strings.TrimSpace(tmpl)
	if content == "" {
		content = defaultPayloadTemplate
	}
	rendered, err := executeTemplate(content, data)
	if err != nil {
		return "", err
	}
	return rendered, nil
}

func renderHeaders(headers map[string]string, data TemplateData) (map[string]string, error) {
	if len(headers) == 0 {
		return map[string]string{}, nil
	}
	out := make(map[string]string, len(headers))
	for key, value := range headers {
		rendered, err := executeTemplate(value, data)
		if err != nil {
			return nil, err
		}
		out[key] = rendered
	}
	return out, nil
}

func executeTemplate(tmpl string, data TemplateData) (string, error) {
	t, err := template.New("tpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"toJSON": func(v any) (string, error) {
				bytes, err := json.Marshal(v)
				if err != nil {
					return "", err
				}
				return string(bytes), nil
			},
		}).
		Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
