import { request } from './http'

export interface FederationRemoteCheckResp {
  manifest?: unknown
  public_key?: unknown
  endpoints?: unknown
}

export interface FederationProxyResp {
  status_code: number
  body: string
}

export interface FederationFriendLinkRequest {
  target_url: string
  message?: string
  rss_url?: string
}

export function checkFederationRemote(targetUrl: string) {
  return request<FederationRemoteCheckResp>('/admin/federation/remote/check', {
    method: 'GET',
    query: { target_url: targetUrl },
  })
}

export function requestFederationFriendlink(payload: FederationFriendLinkRequest) {
  return request<FederationProxyResp>('/admin/federation/friendlinks/request', {
    method: 'POST',
    body: payload,
  })
}
