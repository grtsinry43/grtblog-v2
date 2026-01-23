package webhook

import "errors"

var ErrWebhookNotFound = errors.New("Webhook 不存在")
var ErrDeliveryHistoryNotFound = errors.New("Webhook 投递记录不存在")
var ErrWebhookInvalidEvents = errors.New("Webhook 事件列表无效")
var ErrWebhookDeliveryFailed = errors.New("Webhook 投递失败")
