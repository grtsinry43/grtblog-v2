import { request } from './http'

export interface NavMenuItem {
  id: number
  name: string
  url: string
  icon?: string | null
  sort: number
  parentId?: number | null
  children?: NavMenuItem[]
  createdAt?: string
  updatedAt?: string
}

export interface CreateNavMenuPayload {
  name: string
  url: string
  parentId?: number | null
  icon?: string | null
}

export interface UpdateNavMenuPayload {
  name: string
  url: string
  parentId?: number | null
  icon?: string | null
  sort?: number
}

export interface NavMenuOrderItem {
  id: number
  parentId?: number | null
  sort: number
}

export function listNavMenus() {
  return request<NavMenuItem[]>('/admin/nav-menus')
}

export function createNavMenu(payload: CreateNavMenuPayload) {
  return request<NavMenuItem>('/admin/nav-menus', {
    method: 'POST',
    body: payload,
  })
}

export function updateNavMenu(id: number, payload: UpdateNavMenuPayload) {
  return request<NavMenuItem>(`/admin/nav-menus/${id}`, {
    method: 'PUT',
    body: payload,
  })
}

export function deleteNavMenu(id: number) {
  return request<void>(`/admin/nav-menus/${id}`, {
    method: 'DELETE',
  })
}

export function reorderNavMenus(items: NavMenuOrderItem[]) {
  return request<void>('/admin/nav-menus/reorder', {
    method: 'PUT',
    body: { items },
  })
}
