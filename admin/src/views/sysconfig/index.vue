<script setup lang="ts">
import {
  NButton,
  NCard,
  NCollapse,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NSpin,
  NSwitch,
  NTag,
} from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { ScrollContainer } from '@/components'
import { useDiscreteApi, useLeaveConfirm } from '@/composables'
import { listSysConfigs, updateSysConfigs } from '@/services/sysconfig'

import GroupPanel from './GroupPanel.vue'

import type { SysConfigItem, SysConfigTreeResponse, SysConfigUpdateItem } from '@/services/sysconfig'

defineOptions({
  name: 'SystemSettings',
})

const { message } = useDiscreteApi()

const loading = ref(false)
const saving = ref(false)
const tree = ref<SysConfigTreeResponse | null>(null)
const expandedGroups = ref<string[]>([])

const valueMap = reactive<Record<string, unknown>>({})
const originalMap = reactive<Record<string, unknown>>({})
const jsonBufferMap = reactive<Record<string, string>>({})

const allItems = computed<SysConfigItem[]>(() => {
  const result: SysConfigItem[] = []
  if (tree.value?.items?.length) {
    result.push(...tree.value.items)
  }
  const walk = (groups?: SysConfigTreeResponse['groups']) => {
    if (!groups || groups.length === 0) return
    groups.forEach((group) => {
      if (group.items?.length) {
        result.push(...group.items)
      }
      walk(group.children)
    })
  }
  walk(tree.value?.groups)
  return result
})

const pendingUpdates = computed(() => {
  try {
    return buildUpdateItems().length
  } catch {
    return 0
  }
})

useLeaveConfirm({
  when: computed(() => pendingUpdates.value > 0),
})

onMounted(() => {
  void loadConfigs()
})

async function loadConfigs() {
  loading.value = true
  try {
    const data = await listSysConfigs()
    tree.value = data
    seedMaps(data)
    expandedGroups.value = collectGroupPaths(data.groups)
  } catch (err) {
    message.error(err instanceof Error ? err.message : '加载配置失败')
  } finally {
    loading.value = false
  }
}

async function saveConfigs() {
  let updates: SysConfigUpdateItem[] = []
  try {
    updates = buildUpdateItems()
  } catch (err) {
    message.error(err instanceof Error ? err.message : '配置校验失败')
    return
  }
  if (updates.length === 0) {
    message.warning('没有需要保存的更改')
    return
  }
  saving.value = true
  try {
    const updated = await updateSysConfigs(updates)
    tree.value = updated
    seedMaps(updated)
    message.success('配置已更新')
  } catch (err) {
    message.error(err instanceof Error ? err.message : '保存失败')
  } finally {
    saving.value = false
  }
}

function seedMaps(data: SysConfigTreeResponse) {
  clearObject(valueMap)
  clearObject(originalMap)
  clearObject(jsonBufferMap)
  flattenItems(data).forEach((item) => {
    const key = item.key
    if (item.isSensitive) {
      valueMap[key] = ''
      originalMap[key] = undefined
    } else {
      const resolved = resolveInitialValue(item)
      valueMap[key] = resolved
      originalMap[key] = resolved
    }
    if (item.valueType === 'json') {
      const currentValue = item.isSensitive ? undefined : valueMap[key]
      jsonBufferMap[key] = formatJSON(currentValue)
    }
  })
}

function resolveInitialValue(item: SysConfigItem) {
  if (item.value !== undefined) return item.value
  if (item.defaultValue !== undefined) return item.defaultValue
  switch (item.valueType) {
    case 'bool':
      return false
    case 'number':
      return null
    case 'json':
      return null
    case 'enum':
    case 'string':
    default:
      return ''
  }
}

function buildUpdateItems(): SysConfigUpdateItem[] {
  const updates: SysConfigUpdateItem[] = []
  allItems.value.forEach((item) => {
    const key = item.key
    if (item.isSensitive) {
      const input = String(valueMap[key] ?? '').trim()
      if (input !== '') {
        updates.push({ key, value: input })
      }
      return
    }

    if (item.valueType === 'json') {
      const text = String(jsonBufferMap[key] ?? '').trim()
      if (text === '') {
        return
      }
      let parsed: unknown
      try {
        parsed = JSON.parse(text)
      } catch {
        throw new Error(`配置 ${key} 的 JSON 无法解析`)
      }
      if (!isSameValue(parsed, originalMap[key])) {
        updates.push({ key, value: parsed })
      }
      return
    }

    const nextValue = valueMap[key]
    if (!isSameValue(nextValue, originalMap[key])) {
      updates.push({ key, value: nextValue })
    }
  })
  return updates
}

function isItemVisible(item: SysConfigItem) {
  if (!Array.isArray(item.visibleWhen) || item.visibleWhen.length === 0) return true
  return item.visibleWhen.every((condition) => {
    if (!condition || typeof condition !== 'object') return true
    const { key, op, value } = condition as { key?: string; op?: string; value?: unknown }
    if (!key || !op) return true
    const current = valueMap[key]
    if (op === 'eq') return current === value
    if (op === 'neq') return current !== value
    return true
  })
}

function getInputType(item: SysConfigItem) {
  const inputType = (item.meta as { inputType?: string } | undefined)?.inputType
  if (inputType === 'password') return 'password'
  return 'text'
}

function getEnumOptions(item: SysConfigItem) {
  if (!Array.isArray(item.enumOptions)) return []
  return item.enumOptions.map((option) => {
    if (typeof option === 'string') {
      return { label: option, value: option }
    }
    return {
      label: option.label,
      value: option.value,
    }
  })
}

function flattenItems(data: SysConfigTreeResponse) {
  const result: SysConfigItem[] = []
  if (data.items?.length) {
    result.push(...data.items)
  }
  const walk = (groups?: SysConfigTreeResponse['groups']) => {
    if (!groups || groups.length === 0) return
    groups.forEach((group) => {
      if (group.items?.length) result.push(...group.items)
      walk(group.children)
    })
  }
  walk(data.groups)
  return result
}

function collectGroupPaths(groups?: SysConfigTreeResponse['groups']) {
  const result: string[] = []
  const walk = (nodes?: SysConfigTreeResponse['groups']) => {
    if (!nodes || nodes.length === 0) return
    nodes.forEach((node) => {
      result.push(node.path)
      walk(node.children)
    })
  }
  walk(groups)
  return result
}

function isSameValue(a: unknown, b: unknown) {
  if (a === b) return true
  if (a == null || b == null) return a === b
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

function clearObject(target: Record<string, unknown>) {
  Object.keys(target).forEach((key) => {
    delete target[key]
  })
}
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard>
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <div class="text-base font-semibold">系统设置</div>
            <div class="text-xs text-neutral-500">一些设置因为安全性原因不会返回字段，留空即保持不变</div>
          </div>
          <div class="flex items-center gap-2">
            <NTag v-if="pendingUpdates > 0" type="warning">待保存 {{ pendingUpdates }}</NTag>
            <NButton size="small" secondary :loading="loading" @click="loadConfigs">刷新</NButton>
            <NButton size="small" type="primary" :loading="saving" @click="saveConfigs">保存</NButton>
          </div>
        </div>
      </template>

      <NSpin :show="loading">
        <div v-if="!tree || allItems.length === 0" class="py-8">
          <NEmpty description="暂无配置项" />
        </div>
        <div v-else class="space-y-6">
          <div v-if="tree.items?.length" class="rounded-lg border border-dashed border-neutral-200 p-4">
            <div class="mb-4 text-sm font-medium text-neutral-600">未分组</div>
            <NForm label-placement="left" label-width="160">
              <template v-for="item in tree.items" :key="item.key">
                <NFormItem v-if="isItemVisible(item)">
                  <template #label>
                    <div class="flex items-center gap-2">
                      <span>{{ item.label || item.key }}</span>
                      <!--<NTag size="small" type="info">{{ item.valueType }}</NTag>-->
                      <!--<NTag v-if="item.isSensitive" size="small" type="error">敏感</NTag>-->
                    </div>
                  </template>

                  <NSwitch v-if="item.valueType === 'bool'" v-model:value="valueMap[item.key]" />
                  <NInputNumber
                    v-else-if="item.valueType === 'number'"
                    v-model:value="valueMap[item.key]"
                    :show-button="false"
                    clearable
                  />
                  <NSelect
                    v-else-if="item.valueType === 'enum'"
                    v-model:value="valueMap[item.key]"
                    :options="getEnumOptions(item)"
                    clearable
                  />
                  <NInput
                    v-else-if="item.valueType === 'json'"
                    v-model:value="jsonBufferMap[item.key]"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 6 }"
                    placeholder="请输入 JSON"
                  />
                  <NInput
                    v-else
                    v-model:value="valueMap[item.key]"
                    :type="getInputType(item)"
                    clearable
                    :placeholder="item.isSensitive ? '留空不更新' : ''"
                  />
                </NFormItem>
              </template>
            </NForm>
          </div>

          <NCollapse v-if="tree.groups?.length" v-model:expanded-names="expandedGroups">
            <GroupPanel
              v-for="group in tree.groups"
              :key="group.path"
              :node="group"
              :show-item="isItemVisible"
            >
              <template #item="{ item }">
                <NFormItem>
                  <template #label>
                    <div class="flex items-center gap-2">
                      <span>{{ item.label || item.key }}</span>
                      <!--<NTag size="small" type="info">{{ item.valueType }}</NTag>-->
                      <!--<NTag v-if="item.isSensitive" size="small" type="error">敏感</NTag>-->
                    </div>
                  </template>

                  <NSwitch v-if="item.valueType === 'bool'" v-model:value="valueMap[item.key]" />
                  <NInputNumber
                    v-else-if="item.valueType === 'number'"
                    v-model:value="valueMap[item.key]"
                    :show-button="false"
                    clearable
                  />
                  <NSelect
                    v-else-if="item.valueType === 'enum'"
                    v-model:value="valueMap[item.key]"
                    :options="getEnumOptions(item)"
                    clearable
                  />
                  <NInput
                    v-else-if="item.valueType === 'json'"
                    v-model:value="jsonBufferMap[item.key]"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 6 }"
                    placeholder="请输入 JSON"
                  />
                  <NInput
                    v-else
                    v-model:value="valueMap[item.key]"
                    :type="getInputType(item)"
                    clearable
                    :placeholder="item.isSensitive ? '留空不更新' : ''"
                  />
                </NFormItem>
              </template>
            </GroupPanel>
          </NCollapse>
        </div>
      </NSpin>
    </NCard>
  </ScrollContainer>
</template>
