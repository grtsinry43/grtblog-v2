import { request } from './http'

export interface CategoryItem {
  id: number
  name: string
  shortUrl: string
  createdAt: string
  updatedAt: string
}

export interface TagItem {
  id: number
  name: string
  createdAt: string
  updatedAt: string
}

export function listCategories() {
  return request<CategoryItem[]>('/categories', {
    method: 'GET',
  })
}

export function listTags() {
  return request<TagItem[]>('/tags', {
    method: 'GET',
  })
}

export function createCategory(payload: { name: string; shortUrl: string }) {
  return request<CategoryItem>('/admin/categories', {
    method: 'POST',
    body: payload,
  })
}

export function createTag(name: string) {
  return request<TagItem>('/admin/tags', {
    method: 'POST',
    body: { name },
  })
}
