import { request } from './http'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  isActive: boolean
  createdAt: string
  updatedAt: string
  deletedAt?: string | null
}

export interface LoginResponse {
  token: string
  user: UserInfo
  roles: string[]
  permissions: string[]
}

export interface LoginPayload {
  credential: string
  password: string
  turnstileToken?: string
}

export function login(payload: LoginPayload) {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export interface AccessInfoResponse {
  user: UserInfo
  roles: string[]
  permissions: string[]
}

export function getAccessInfo() {
  return request<AccessInfoResponse>('/auth/access-info', {
    method: 'GET',
  })
}

export interface UpdateProfilePayload {
  nickname?: string
  avatar?: string
  email?: string
}

export function updateProfile(payload: UpdateProfilePayload) {
  return request<UserInfo>('/auth/profile', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export interface ChangePasswordPayload {
  oldPassword: string
  newPassword: string
}

export function changePassword(payload: ChangePasswordPayload) {
  return request<null>('/auth/password', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export interface OAuthBinding {
  providerKey: string
  providerName: string
  oauthID: string
  boundAt: string
  expiresAt?: string | null
  providerScope?: string
}

export function getOAuthBindings() {
  return request<OAuthBinding[]>('/auth/oauth-bindings', {
    method: 'GET',
  })
}
