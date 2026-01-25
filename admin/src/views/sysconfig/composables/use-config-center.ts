import { useMessage } from 'naive-ui'
import { ref, reactive, computed, onMounted } from 'vue'

import type { SysConfigTreeResponse, SysConfigItem, SysConfigUpdateItem, SysConfigGroup } from '@/services/sysconfig'

export type ConfigListFn = (keys?: string[]) => Promise<SysConfigTreeResponse>
export type ConfigUpdateFn = (items: SysConfigUpdateItem[]) => Promise<SysConfigTreeResponse>

export function useConfigCenter(listFn: ConfigListFn, updateFn: ConfigUpdateFn) {
  const message = useMessage()

  const loading = ref(false)
  const saving = ref(false)
  const tree = ref<SysConfigTreeResponse | null>(null)

  // 核心状态映射
  const valueMap = reactive<Record<string, unknown>>({})
  const originalMap = reactive<Record<string, unknown>>({})
  const jsonBufferMap = reactive<Record<string, string>>({})

  // 展开的折叠面板
  const expandedGroups = ref<string[]>([])

  // --- 辅助函数：扁平化遍历 ---
  const getAllItems = (data: SysConfigTreeResponse | null) => {
    const result: SysConfigItem[] = []
    if (!data) return result
    if (data.items) result.push(...data.items)

    const walk = (groups?: SysConfigGroup[]) => {
      if (!groups) return
      groups.forEach(g => {
        if (g.items) result.push(...g.items)
        walk(g.children)
      })
    }
    walk(data.groups)
    return result
  }

  // --- 核心逻辑：初始化数据 ---
  function seedMaps(data: SysConfigTreeResponse) {
    // 清空现有数据
    Object.keys(valueMap).forEach(k => delete valueMap[k])
    Object.keys(originalMap).forEach(k => delete originalMap[k])
    Object.keys(jsonBufferMap).forEach(k => delete jsonBufferMap[k])

    getAllItems(data).forEach((item) => {
      const key = item.key
      // 敏感数据处理：值置空，原始值置 undefined
      if (item.isSensitive) {
        valueMap[key] = ''
        originalMap[key] = undefined
      } else {
        const val = resolveInitialValue(item)
        valueMap[key] = val
        originalMap[key] = val
      }

      // JSON 特殊处理
      if (item.valueType === 'json') {
        const current = item.isSensitive ? undefined : valueMap[key]
        jsonBufferMap[key] = formatJSON(current)
      }
    })
  }

  function resolveInitialValue(item: SysConfigItem) {
    if (item.value !== undefined) return item.value
    if (item.defaultValue !== undefined) return item.defaultValue
    switch (item.valueType) {
      case 'bool': return false
      case 'number': return null
      case 'json': return null
      default: return ''
    }
  }

  // --- 核心逻辑：构建更新包 ---
  function buildUpdateItems(): SysConfigUpdateItem[] {
    const updates: SysConfigUpdateItem[] = []
    const allItems = getAllItems(tree.value)

    allItems.forEach((item) => {
      const key = item.key

      // 敏感字段：只有输入了内容才更新
      if (item.isSensitive) {
        const input = String(valueMap[key] ?? '').trim()
        if (input !== '') {
          updates.push({ key, value: input })
        }
        return
      }

      // JSON 字段：解析后对比
      if (item.valueType === 'json') {
        const text = String(jsonBufferMap[key] ?? '').trim()
        if (!text) return

        let parsed: unknown
        try {
          parsed = JSON.parse(text)
        } catch {
          throw new Error(`配置 ${item.label || key} 的 JSON 格式错误`)
        }

        if (!isSameValue(parsed, originalMap[key])) {
          updates.push({ key, value: parsed })
        }
        return
      }

      // 普通字段
      const nextValue = valueMap[key]
      if (!isSameValue(nextValue, originalMap[key])) {
        updates.push({ key, value: nextValue })
      }
    })
    return updates
  }

  // --- 业务逻辑：可见性判断 ---
  function isItemVisible(item: SysConfigItem): boolean {
    if (!Array.isArray(item.visibleWhen) || item.visibleWhen.length === 0) return true
    return item.visibleWhen.every((condition: any) => {
      if (!condition || typeof condition !== 'object') return true
      const { key, op, value } = condition
      if (!key || !op) return true
      const current = valueMap[key]
      if (op === 'eq') return current === value
      if (op === 'neq') return current !== value
      return true
    })
  }

  // API 交互
  async function fetch() {
    loading.value = true
    try {
      const data = await listFn()
      tree.value = data
      seedMaps(data)

      // 自动展开所有组
      const paths: string[] = []
      const walk = (gs?: SysConfigGroup[]) => {
        gs?.forEach(g => {
          paths.push(g.path)
          walk(g.children)
        })
      }
      walk(data.groups)
      expandedGroups.value = paths
    } catch (e: any) {
      message.error(e.message || '加载配置失败')
    } finally {
      loading.value = false
    }
  }

  async function save() {
    let updates: SysConfigUpdateItem[] = []
    try {
      updates = buildUpdateItems()
    } catch (e: any) {
      message.error(e.message)
      return
    }

    if (updates.length === 0) {
      message.warning('没有检测到更改')
      return
    }

    saving.value = true
    try {
      const updated = await updateFn(updates)
      tree.value = updated
      seedMaps(updated)
      message.success('保存成功')
    } catch (e: any) {
      message.error(e.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  // 挂载自动加载
  onMounted(fetch)

  return {
    loading,
    saving,
    tree,
    valueMap,
    jsonBufferMap,
    expandedGroups,
    isItemVisible, // 暴露给组件使用
    fetch,
    save,
    // 导出用于脏检查的计算属性
    pendingCount: computed(() => {
      try { return buildUpdateItems().length } catch { return 0 }
    })
  }
}

// Utils
function isSameValue(a: unknown, b: unknown) {
  if (a === b) return true
  if (a == null || b == null) return a === b // 处理 null vs undefined
  return JSON.stringify(a) === JSON.stringify(b)
}

function formatJSON(value: unknown) {
  if (value === undefined || value === null) return ''
  if (typeof value === 'string') return value
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return ''
  }
}
