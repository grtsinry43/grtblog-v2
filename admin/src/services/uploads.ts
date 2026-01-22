import { request } from './http'

export type FileType = 'picture' | 'file'

export interface UploadFileResponse {
  id: number
  name: string
  path: string
  publicUrl: string
  type: FileType
  size: number
  createdAt: string
  duplicated: boolean
}

export interface UploadFileListResponse {
  items: UploadFileResponse[]
  total: number
  page: number
  size: number
}

export interface ListUploadsParams {
  page?: number
  pageSize?: number
}

export interface RenameFilePayload {
  name: string
}

/**
 * Upload a file to the server
 * @param file - The file to upload
 * @param type - File type: 'picture' or 'file'
 */
export async function uploadFile(file: File, type: FileType): Promise<UploadFileResponse> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('type', type)

  return request<UploadFileResponse>('/upload', {
    method: 'POST',
    body: formData,
  })
}

/**
 * List uploaded files with pagination
 * @param params - Query parameters for listing files
 */
export function listUploads(params: ListUploadsParams = {}): Promise<UploadFileListResponse> {
  const query: Record<string, string> = {}
  
  if (params.page !== undefined) {
    query.page = String(params.page)
  }
  if (params.pageSize !== undefined) {
    query.pageSize = String(params.pageSize)
  }

  return request<UploadFileListResponse>('/uploads', {
    method: 'GET',
    query,
  })
}

/**
 * Rename a file (display name only)
 * @param id - File ID
 * @param payload - New name
 */
export function renameFile(id: number, payload: RenameFilePayload): Promise<UploadFileResponse> {
  return request<UploadFileResponse>(`/upload/${id}`, {
    method: 'PUT',
    body: payload,
  })
}

/**
 * Delete a file
 * @param id - File ID
 */
export function deleteFile(id: number): Promise<null> {
  return request<null>(`/upload/${id}`, {
    method: 'DELETE',
  })
}

/**
 * Download a file
 * @param id - File ID
 * @returns Blob URL for download
 */
export async function downloadFile(id: number, fileName: string): Promise<void> {
  // For download, we need to handle this differently since it returns binary data
  // We'll need to use a direct fetch approach
  const token = localStorage.getItem('token')
  const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v2').replace(/\/$/, '')
  
  const response = await fetch(`${API_BASE_URL}/upload/${id}/download`, {
    headers: {
      Authorization: token ? `Bearer ${token}` : '',
    },
  })

  if (!response.ok) {
    throw new Error('Download failed')
  }

  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  document.body.appendChild(a)
  a.click()
  window.URL.revokeObjectURL(url)
  document.body.removeChild(a)
}

/**
 * Get the full public URL for a file
 * @param path - Virtual path from the upload response
 */
export function getPublicUrl(path: string): string {
  const baseUrl = window.location.origin
  return `${baseUrl}${path}`
}
