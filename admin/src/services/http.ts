import { $fetch, FetchError } from 'ofetch'

import type { FetchOptions, FetchResponse } from 'ofetch'

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v2').replace(/\/$/, '')

type ResponseInterceptor<T> = (ctx: { data: T; envelope: ApiEnvelope<T> }) => void | Promise<void>
type ErrorInterceptor = (error: ApiError) => void | Promise<void>

export interface ApiMeta {
  requestId?: string
  timestamp?: string
}

export interface ApiEnvelope<T> {
  code: number
  bizErr: string
  msg: string
  data: T
  meta: ApiMeta
}

export class ApiError extends Error {
  code?: number
  bizErr?: string
  status?: number
  meta?: ApiMeta

  constructor(message: string, options: { code?: number; bizErr?: string; status?: number; meta?: ApiMeta } = {}) {
    super(message)
    this.name = 'ApiError'
    this.code = options.code
    this.bizErr = options.bizErr
    this.status = options.status
    this.meta = options.meta
    Object.setPrototypeOf(this, ApiError.prototype)
  }
}

let tokenProvider: (() => string | null) | null = null
const responseInterceptors: ResponseInterceptor<any>[] = []
const errorInterceptors: ErrorInterceptor[] = []

export function setAuthTokenProvider(provider: () => string | null) {
  tokenProvider = provider
}

export function addResponseInterceptor<T>(interceptor: ResponseInterceptor<T>) {
  responseInterceptors.push(interceptor as ResponseInterceptor<any>)
  return () => {
    const idx = responseInterceptors.indexOf(interceptor as ResponseInterceptor<any>)
    if (idx >= 0) responseInterceptors.splice(idx, 1)
  }
}

export function addErrorInterceptor(interceptor: ErrorInterceptor) {
  errorInterceptors.push(interceptor)
  return () => {
    const idx = errorInterceptors.indexOf(interceptor)
    if (idx >= 0) errorInterceptors.splice(idx, 1)
  }
}

function normalizePath(path: string) {
  if (!path.startsWith('/')) return `/${path}`
  return path
}

function shouldSetJsonContentType(body: FetchOptions['body']) {
  if (!body) return false
  if (typeof body === 'string') return true
  if (typeof FormData !== 'undefined' && body instanceof FormData) return false
  if (typeof URLSearchParams !== 'undefined' && body instanceof URLSearchParams) return false
  if (typeof Blob !== 'undefined' && body instanceof Blob) return false
  if (typeof ArrayBuffer !== 'undefined' && body instanceof ArrayBuffer) return false
  return true
}

export async function request<T>(path: string, options: FetchOptions = {}): Promise<T> {
  const headers = new Headers(options.headers || {})
  if (shouldSetJsonContentType(options.body) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  const bearer = tokenProvider ? tokenProvider() : null
  if (bearer && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${bearer}`)
  }

  try {
    let response: FetchResponse<string>
    try {
      response = await $fetch.raw(`${API_BASE_URL}${normalizePath(path)}`, {
        ...options,
        headers,
        responseType: 'text',
      })
    } catch (error) {
      if (error instanceof FetchError && error.response) {
        response = error.response as FetchResponse<string>
      } else {
        throw error
      }
    }

    const status = response.status
    const text = typeof response._data === 'string' ? response._data : ''

    let payload: ApiEnvelope<T> | null = null
    if (text) {
      try {
        payload = JSON.parse(text) as ApiEnvelope<T>
      } catch {
        // ignore json parse error, handled below
      }
    }

    if (status < 200 || status >= 300) {
      throw new ApiError(payload?.msg || `请求失败（${status}）`, {
        code: payload?.code,
        bizErr: payload?.bizErr,
        status,
        meta: payload?.meta,
      })
    }

    if (!payload) {
      throw new ApiError('无法解析服务端响应', { status })
    }

    if (payload.code !== 0) {
      throw new ApiError(payload.msg || payload.bizErr || '请求失败', {
        code: payload.code,
        bizErr: payload.bizErr,
        status,
        meta: payload.meta,
      })
    }

    const data = payload.data
    for (const interceptor of responseInterceptors) {
      await interceptor({ data, envelope: payload })
    }
    return data
  } catch (error) {
    if (error instanceof ApiError) {
      for (const interceptor of errorInterceptors) {
        await interceptor(error)
      }
      throw error
    }
    if (error instanceof Error) {
      const err = new ApiError(error.message)
      for (const interceptor of errorInterceptors) {
        await interceptor(err)
      }
      throw err
    }
    const err = new ApiError('网络异常，请稍后重试')
    for (const interceptor of errorInterceptors) {
      await interceptor(err)
    }
    throw err
  }
}
