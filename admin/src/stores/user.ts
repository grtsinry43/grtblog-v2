import { useStorage } from '@vueuse/core'
import { acceptHMRUpdate, defineStore, storeToRefs } from 'pinia'
import { ref } from 'vue'

import router from '@/router'
import { resolveMenu, resolveRoute } from '@/router/helper'
import { routeRecordRaw } from '@/router/record'

import { pinia } from '.'

import { getAccessInfo } from '@/services/auth'

import type { MenuMixedOptions } from '@/router/interface'
import type { MenuOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'

interface User {
  avatar: string
  email: string
  id: number | null
  nickname: string
  roles: string[]
  permissions: string[]
  username: string
  createdAt: string
  updatedAt: string
}

const createEmptyUser = (): User => ({
  avatar: '',
  email: '',
  id: null,
  nickname: '',
  roles: [],
  permissions: [],
  username: '',
  createdAt: '',
  updatedAt: '',
})

export const useUserStore = defineStore('userStore', () => {
  const user = useStorage<User>('user', createEmptyUser())

  const token = useStorage<string | null>('token', null)

  const menuList = ref<MenuOption[]>([])

  const routeList = ref<RouteRecordRaw[]>([])

  function payloadFromToken(value: string | null) {
    if (!value) return null
    const parts = value.split('.')
    if (parts.length < 2 || !parts[1]) return null
    try {
      const normalized = parts[1].replace(/-/g, '+').replace(/_/g, '/')
      const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=')
      const decoded = typeof atob === 'function' ? atob(padded) : ''
      return JSON.parse(decoded)
    } catch {
      return null
    }
  }

  function rolesFromToken(value: string | null) {
    const payload = payloadFromToken(value)
    if (payload && Array.isArray(payload.roles)) {
      return payload.roles.filter((role: unknown) => typeof role === 'string')
    }
    return []
  }

  function permissionsFromToken(value: string | null) {
    const payload = payloadFromToken(value)
    if (payload && Array.isArray(payload.perms)) {
      return payload.perms.filter((perm: unknown) => typeof perm === 'string')
    }
    return []
  }

  function setAuth(session: { token: string; user: Partial<User> }) {
    token.value = session.token
    const roles = session.user.roles?.length ? session.user.roles : rolesFromToken(session.token)
    const permissions = session.user.permissions?.length
      ? session.user.permissions
      : permissionsFromToken(session.token)
    user.value = {
      ...createEmptyUser(),
      ...session.user,
      roles,
      permissions,
      createdAt: session.user.createdAt || '',
      updatedAt: session.user.updatedAt || '',
    }
  }

  async function resolveMenuRoute() {
    const roles = user.value.roles.length ? user.value.roles : rolesFromToken(token.value)
    const isAdmin = roles.includes('admin')

    const res = await new Promise<MenuMixedOptions[]>((resolve) => {
      if (isAdmin) {
        resolve(routeRecordRaw)
      } else {
        const allowedRoutes = ['dashboard', 'dataShow', 'notfoundPage', 'about', 'userCenter']
        const filteredRoutes = routeRecordRaw.filter((route) => {
          return !route.type && route.name && allowedRoutes.includes(route.name as string)
        })
        resolve(filteredRoutes)
      }
    })

    const resolvedMenu = resolveMenu(res) || []
    const resolvedRoute = resolveRoute(res) || []

    menuList.value = resolvedMenu
    routeList.value = resolvedRoute

    return {
      menuList: resolvedMenu,
      routeList: resolvedRoute,
    }
  }

  async function refreshAccessInfo() {
    if (!token.value) return
    const result = await getAccessInfo()
    setAuth({
      token: token.value,
      user: {
        id: result.user.id,
        username: result.user.username,
        nickname: result.user.nickname,
        email: result.user.email,
        avatar: result.user.avatar,
        roles: result.roles,
        permissions: result.permissions,
      },
    })
  }

  function cleanup(redirectPath?: string) {
    router.replace({
      name: 'signIn',
      ...(redirectPath ? { query: { r: redirectPath } } : {}),
    })

    token.value = null
    user.value = createEmptyUser()

    if (router.hasRoute('layout')) {
      router.removeRoute('layout')
    }

    menuList.value = []

    routeList.value = []
  }

  return {
    user,
    token,
    menuList,
    routeList,
    setAuth,
    refreshAccessInfo,
    resolveMenuRoute,
    cleanup,
  }
})

export function toRefsUserStore() {
  return {
    ...storeToRefs(useUserStore(pinia)),
  }
}

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
}
