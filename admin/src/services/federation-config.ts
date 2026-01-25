import { request } from './http'

import type { SysConfigTreeResponse, SysConfigUpdateItem } from './sysconfig'

export function listFederationConfigs(keys?: string[]) {
  const query = keys && keys.length > 0 ? { keys: keys.join(',') } : undefined
  return request<SysConfigTreeResponse>('/admin/federation/config', {
    method: 'GET',
    query,
  })
}

export function updateFederationConfigs(items: SysConfigUpdateItem[]) {
  return request<SysConfigTreeResponse>('/admin/federation/config', {
    method: 'PUT',
    body: { items },
  })
}
