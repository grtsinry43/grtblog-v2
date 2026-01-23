package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	domainwebhook "github.com/grtsinry43/grtblog-v2/server/internal/domain/webhook"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(ctx context.Context, hook *domainwebhook.Webhook) error {
	rec, err := mapWebhookToModel(hook)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	hook.ID = rec.ID
	hook.CreatedAt = rec.CreatedAt
	hook.UpdatedAt = rec.UpdatedAt
	return nil
}

func (r *WebhookRepository) Update(ctx context.Context, hook *domainwebhook.Webhook) error {
	rec, err := mapWebhookToModel(hook)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).
		Model(&model.Webhook{}).
		Where("id = ?", hook.ID).
		Updates(map[string]any{
			"name":             rec.Name,
			"url":              rec.URL,
			"events":           rec.Events,
			"headers":          rec.Headers,
			"payload_template": rec.PayloadTemplate,
			"is_enabled":       rec.IsEnabled,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainwebhook.ErrWebhookNotFound
	}
	return nil
}

func (r *WebhookRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Webhook{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainwebhook.ErrWebhookNotFound
	}
	return nil
}

func (r *WebhookRepository) GetByID(ctx context.Context, id int64) (*domainwebhook.Webhook, error) {
	var rec model.Webhook
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainwebhook.ErrWebhookNotFound
		}
		return nil, err
	}
	return mapWebhookToDomain(rec)
}

func (r *WebhookRepository) List(ctx context.Context) ([]*domainwebhook.Webhook, error) {
	var records []model.Webhook
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*domainwebhook.Webhook, len(records))
	for i, rec := range records {
		item, err := mapWebhookToDomain(rec)
		if err != nil {
			return nil, err
		}
		result[i] = item
	}
	return result, nil
}

func (r *WebhookRepository) ListEnabledByEvent(ctx context.Context, eventName string) ([]*domainwebhook.Webhook, error) {
	var records []model.Webhook
	query := r.db.WithContext(ctx).Where("is_enabled = ?", true)
	if eventName != "" {
		query = query.Where("events @> ?", fmt.Sprintf("[\"%s\"]", eventName))
	}
	if err := query.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	result := make([]*domainwebhook.Webhook, len(records))
	for i, rec := range records {
		item, err := mapWebhookToDomain(rec)
		if err != nil {
			return nil, err
		}
		result[i] = item
	}
	return result, nil
}

func (r *WebhookRepository) CreateHistory(ctx context.Context, history *domainwebhook.DeliveryHistory) error {
	rec, err := mapHistoryToModel(history)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return err
	}
	history.ID = rec.ID
	history.CreatedAt = rec.CreatedAt
	return nil
}

func (r *WebhookRepository) GetHistoryByID(ctx context.Context, id int64) (*domainwebhook.DeliveryHistory, error) {
	var rec model.WebhookHistory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainwebhook.ErrDeliveryHistoryNotFound
		}
		return nil, err
	}
	return mapHistoryToDomain(rec)
}

func (r *WebhookRepository) ListHistory(ctx context.Context, options domainwebhook.DeliveryHistoryListOptions) ([]*domainwebhook.DeliveryHistory, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.WebhookHistory{})
	if options.WebhookID != nil {
		query = query.Where("webhook_id = ?", *options.WebhookID)
	}
	if options.EventName != nil && *options.EventName != "" {
		query = query.Where("event_name = ?", *options.EventName)
	}
	if options.IsTest != nil {
		query = query.Where("is_test = ?", *options.IsTest)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (options.Page - 1) * options.PageSize
	var records []model.WebhookHistory
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(options.PageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*domainwebhook.DeliveryHistory, len(records))
	for i, rec := range records {
		item, err := mapHistoryToDomain(rec)
		if err != nil {
			return nil, 0, err
		}
		result[i] = item
	}
	return result, total, nil
}

func mapWebhookToModel(hook *domainwebhook.Webhook) (model.Webhook, error) {
	eventsBytes, err := json.Marshal(hook.Events)
	if err != nil {
		return model.Webhook{}, err
	}
	headersBytes, err := json.Marshal(hook.Headers)
	if err != nil {
		return model.Webhook{}, err
	}
	return model.Webhook{
		ID:              hook.ID,
		Name:            hook.Name,
		URL:             hook.URL,
		Events:          eventsBytes,
		Headers:         headersBytes,
		PayloadTemplate: hook.PayloadTemplate,
		IsEnabled:       hook.IsEnabled,
	}, nil
}

func mapWebhookToDomain(rec model.Webhook) (*domainwebhook.Webhook, error) {
	var events []string
	if err := json.Unmarshal(rec.Events, &events); err != nil {
		return nil, err
	}
	headers := map[string]string{}
	if len(rec.Headers) > 0 {
		if err := json.Unmarshal(rec.Headers, &headers); err != nil {
			return nil, err
		}
	}
	return &domainwebhook.Webhook{
		ID:              rec.ID,
		Name:            rec.Name,
		URL:             rec.URL,
		Events:          events,
		Headers:         headers,
		PayloadTemplate: rec.PayloadTemplate,
		IsEnabled:       rec.IsEnabled,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       rec.UpdatedAt,
		DeletedAt:       deletedAtToPtr(rec.DeletedAt),
	}, nil
}

func mapHistoryToModel(history *domainwebhook.DeliveryHistory) (model.WebhookHistory, error) {
	if history.RequestHeaders == nil {
		history.RequestHeaders = map[string]string{}
	}
	headersBytes, err := json.Marshal(history.RequestHeaders)
	if err != nil {
		return model.WebhookHistory{}, err
	}
	if history.ResponseHeaders == nil {
		history.ResponseHeaders = map[string]string{}
	}
	responseHeadersBytes, err := json.Marshal(history.ResponseHeaders)
	if err != nil {
		return model.WebhookHistory{}, err
	}
	return model.WebhookHistory{
		ID:              history.ID,
		WebhookID:       history.WebhookID,
		EventName:       history.EventName,
		RequestURL:      history.RequestURL,
		RequestHeaders:  headersBytes,
		RequestBody:     history.RequestBody,
		ResponseStatus:  history.ResponseStatus,
		ResponseHeaders: responseHeadersBytes,
		ResponseBody:    history.ResponseBody,
		ErrorMessage:    history.ErrorMessage,
		IsTest:          history.IsTest,
	}, nil
}

func mapHistoryToDomain(rec model.WebhookHistory) (*domainwebhook.DeliveryHistory, error) {
	headers := map[string]string{}
	if len(rec.RequestHeaders) > 0 {
		if err := json.Unmarshal(rec.RequestHeaders, &headers); err != nil {
			return nil, err
		}
	}
	responseHeaders := map[string]string{}
	if len(rec.ResponseHeaders) > 0 {
		if err := json.Unmarshal(rec.ResponseHeaders, &responseHeaders); err != nil {
			return nil, err
		}
	}
	return &domainwebhook.DeliveryHistory{
		ID:              rec.ID,
		WebhookID:       rec.WebhookID,
		EventName:       rec.EventName,
		RequestURL:      rec.RequestURL,
		RequestHeaders:  headers,
		RequestBody:     rec.RequestBody,
		ResponseStatus:  rec.ResponseStatus,
		ResponseHeaders: responseHeaders,
		ResponseBody:    rec.ResponseBody,
		ErrorMessage:    rec.ErrorMessage,
		IsTest:          rec.IsTest,
		CreatedAt:       rec.CreatedAt,
	}, nil
}
