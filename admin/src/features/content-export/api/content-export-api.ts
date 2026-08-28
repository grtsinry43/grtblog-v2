import { request } from '@/services/http'

import type { ExportDownloadTicket, ExportMode, ExportRecord } from '../model/types'

const basePath = '/admin/exports'

export function listExports() {
  return request<ExportRecord[]>(basePath)
}

export function createExport(mode: ExportMode) {
  return request<ExportRecord>(basePath, { method: 'POST', body: { mode } })
}

export function deleteExport(id: string) {
  return request<{ id: string }>(`${basePath}/${encodeURIComponent(id)}`, { method: 'DELETE' })
}

export function createExportDownloadTicket(id: string) {
  return request<ExportDownloadTicket>(`${basePath}/${encodeURIComponent(id)}/download-ticket`, {
    method: 'POST',
  })
}
