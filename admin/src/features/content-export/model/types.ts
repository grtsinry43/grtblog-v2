export type ExportMode = 'structured' | 'flatten' | 'both'

export type ExportStatus = 'queued' | 'running' | 'completed' | 'failed'

export interface ExportRecord {
  id: string
  filename: string
  status: ExportStatus
  stage: string
  triggerType: string
  mode: ExportMode
  sizeBytes: number
  sha256?: string
  appVersion?: string
  siteName?: string
  siteUrl?: string
  articleCount: number
  momentsCount: number
  pagesCount: number
  thinkingsCount: number
  imageCount: number
  failedImageCount: number
  errorMessage?: string
  createdAt: string
  startedAt?: string
  completedAt?: string
}

export interface ExportDownloadTicket {
  url: string
  expiresAt: string
}
