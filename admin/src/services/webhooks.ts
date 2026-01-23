import { request } from './http'

export interface WebhookItem {
  id: number
  name: string
  url: string
  events: string[]
  headers: Record<string, string>
  payloadTemplate: string
  isEnabled: boolean
  createdAt: string
  updatedAt: string
}

export interface WebhookHistoryItem {
  id: number
  webhookId: number
  eventName: string
  requestUrl: string
  requestHeaders: Record<string, string>
  requestBody: string
  responseStatus: number
  responseHeaders: Record<string, string>
  responseBody?: string
  errorMessage?: string
  isTest: boolean
  createdAt: string
}

export interface WebhookHistoryListResponse {
  items: WebhookHistoryItem[]
  total: number
  page: number
  size: number
}

export interface WebhookEventListResponse {
  events: string[]
}

export interface CreateWebhookPayload {
  name: string
  url: string
  events: string[]
  headers?: Record<string, string>
  payloadTemplate?: string
  isEnabled: boolean
}

export interface UpdateWebhookPayload {
  name: string
  url: string
  events: string[]
  headers?: Record<string, string>
  payloadTemplate?: string
  isEnabled: boolean
}

export interface ListWebhookHistoryParams {
  page?: number
  pageSize?: number
  webhookId?: number
  eventName?: string
  isTest?: boolean
}

function stripEmpty<T extends Record<string, unknown>>(value: T) {
  return Object.fromEntries(
    Object.entries(value).filter(([, entry]) => entry !== undefined && entry !== null && entry !== ''),
  ) as T
}

export function listWebhooks() {
  return request<WebhookItem[]>('/admin/webhooks', {
    method: 'GET',
  })
}

export function listWebhookEvents() {
  return request<WebhookEventListResponse>('/admin/webhooks/events', {
    method: 'GET',
  })
}

export function createWebhook(payload: CreateWebhookPayload) {
  return request<WebhookItem>('/admin/webhooks', {
    method: 'POST',
    body: payload,
  })
}

export function updateWebhook(id: number, payload: UpdateWebhookPayload) {
  return request<WebhookItem>(`/admin/webhooks/${id}`, {
    method: 'PUT',
    body: payload,
  })
}

export function deleteWebhook(id: number) {
  return request<void>(`/admin/webhooks/${id}`, {
    method: 'DELETE',
  })
}

export function testWebhook(id: number, eventName?: string | null) {
  return request<void>(`/admin/webhooks/${id}/test`, {
    method: 'POST',
    body: eventName ? { eventName } : {},
  })
}

export function listWebhookHistory(params: ListWebhookHistoryParams) {
  return request<WebhookHistoryListResponse>('/admin/webhooks/deliveries', {
    method: 'GET',
    query: stripEmpty(params),
  })
}

export function replayWebhookHistory(id: number) {
  return request<void>(`/admin/webhooks/deliveries/${id}/replay`, {
    method: 'POST',
  })
}
