import { request } from './http'

export interface ArticleListItem {
  id: number
  title: string
  shortUrl: string
  authorName?: string
  summary: string
  avatar?: string
  cover?: string
  views: number
  categoryName?: string
  categoryShortUrl?: string
  tags: string[]
  likes: number
  comments: number
  isTop: boolean
  isHot: boolean
  isOriginal: boolean
  createdAt: string
  updatedAt: string
}

export interface ArticleListResponse {
  items: ArticleListItem[]
  total: number
  page: number
  size: number
}

export interface ArticleTag {
  id: number
  name: string
}

export interface ArticleDetail {
  id: number
  title: string
  summary: string
  aiSummary?: string | null
  leadIn?: string | null
  content: string
  contentHash: string
  authorId: number
  cover?: string | null
  categoryId?: number | null
  shortUrl: string
  isPublished: boolean
  isTop: boolean
  isHot: boolean
  isOriginal: boolean
  tags?: ArticleTag[]
  createdAt: string
  updatedAt: string
}

export interface ListArticlesParams {
  page?: number
  pageSize?: number
  categoryId?: number
  tagId?: number
  authorId?: number
  published?: boolean
  search?: string
}

export interface CreateArticlePayload {
  title: string
  summary: string
  leadIn?: string | null
  content: string
  cover?: string | null
  categoryId?: number | null
  tagIds?: number[]
  shortUrl?: string | null
  isPublished: boolean
  isTop: boolean
  isHot: boolean
  isOriginal: boolean
  createdAt?: string | null
}

export interface UpdateArticlePayload {
  title: string
  summary: string
  leadIn?: string | null
  content: string
  cover?: string | null
  categoryId?: number | null
  tagIds?: number[]
  shortUrl: string
  isPublished: boolean
  isTop: boolean
  isHot: boolean
  isOriginal: boolean
}

function stripEmpty<T extends Record<string, unknown>>(value: T) {
  return Object.fromEntries(
    Object.entries(value).filter(
      ([, entry]) => entry !== undefined && entry !== null && entry !== '',
    ),
  ) as T
}

export function listArticles(params: ListArticlesParams) {
  return request<ArticleListResponse>('/articles', {
    method: 'GET',
    query: stripEmpty(params),
  })
}

export function getArticle(id: number) {
  return request<ArticleDetail>(`/articles/${id}`, {
    method: 'GET',
  })
}

export function createArticle(payload: CreateArticlePayload) {
  return request<ArticleDetail>('/articles', {
    method: 'POST',
    body: payload,
  })
}

export function updateArticle(id: number, payload: UpdateArticlePayload) {
  return request<ArticleDetail>(`/articles/${id}`, {
    method: 'PUT',
    body: payload,
  })
}
