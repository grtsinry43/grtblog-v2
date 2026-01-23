<script setup lang="ts">
import {
  NAlert,
  NButton,
  NButtonGroup,
  NCard,
  NCheckbox,
  NCheckboxGroup,
  NCode,
  NDataTable,
  NDivider,
  NDrawer,
  NDrawerContent,
  NEmpty,
  NForm,
  NFormItem,
  NGi,
  NGrid,
  NInput,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NStatistic,
  NSpace,
  NSwitch,
  NTabPane,
  NTable,
  NTag,
  NTabs,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'

import { ScrollContainer } from '@/components'
import TemplateEditor from '@/components/template-editor/TemplateEditor.vue'
import { formatTemplateJson } from '@/composables/template-editor/json-lint'
import {
  createWebhook,
  deleteWebhook,
  listWebhookEvents,
  listWebhooks,
  listWebhookHistory,
  replayWebhookHistory,
  testWebhook,
  updateWebhook,
} from '@/services/webhooks'

import type { DataTableColumns, SelectOption } from 'naive-ui'
import type { WebhookHistoryItem, WebhookItem } from '@/services/webhooks'

type HeaderRow = {
  key: string
  value: string
}

defineOptions({
  name: 'WebhookList',
})

const message = useMessage()

const activeTab = ref('list')
const webhooks = ref<WebhookItem[]>([])
const events = ref<string[]>([])
const loading = ref(false)
const historyLoading = ref(false)

const history = ref<WebhookHistoryItem[]>([])
const historyPage = ref(1)
const historyPageSize = ref(10)
const historyTotal = ref(0)

const formDrawerVisible = ref(false)
const saving = ref(false)
const editingWebhook = ref<WebhookItem | null>(null)

const testModalVisible = ref(false)
const testingWebhook = ref<WebhookItem | null>(null)
const testEventName = ref<string | null>(null)

const historyDrawerVisible = ref(false)
const activeHistory = ref<WebhookHistoryItem | null>(null)

const form = reactive({
  name: '',
  url: '',
  events: [] as string[],
  headers: [] as HeaderRow[],
  payloadTemplate: '',
  isEnabled: true,
})

const listFilters = reactive({
  keyword: '',
  status: 'all' as 'all' | 'enabled' | 'disabled',
  event: null as string | null,
})

const historyFilters = reactive({
  webhookId: null as number | null,
  eventName: null as string | null,
  isTest: null as boolean | null,
})

const eventOptions = computed<SelectOption[]>(() =>
  events.value.map((item) => ({
    label: item,
    value: item,
  })),
)

const webhookOptions = computed<SelectOption[]>(() =>
  webhooks.value.map((item) => ({
    label: item.name,
    value: item.id,
  })),
)

const statusOptions: SelectOption[] = [
  { label: '全部状态', value: 'all' },
  { label: '启用中', value: 'enabled' },
  { label: '已停用', value: 'disabled' },
]

const webhookMap = computed(() => {
  return new Map(webhooks.value.map((item) => [item.id, item.name]))
})

const formTitle = computed(() => (editingWebhook.value ? '编辑 Webhook' : '新建 Webhook'))
const formActionLabel = computed(() => (editingWebhook.value ? '保存' : '创建'))

const totalWebhooks = computed(() => webhooks.value.length)
const enabledCount = computed(() => webhooks.value.filter((item) => item.isEnabled).length)
const disabledCount = computed(() => totalWebhooks.value - enabledCount.value)
const historyFailureCount = computed(
  () => history.value.filter((item) => item.responseStatus < 200 || item.responseStatus >= 300).length,
)
const latestHistory = computed(() => history.value[0] ?? null)
const latestHistoryStatus = computed(() => {
  const entry = latestHistory.value
  if (!entry) {
    return { label: '暂无', type: 'default' as const }
  }
  const success = entry.responseStatus >= 200 && entry.responseStatus < 300
  return {
    label: success ? '成功' : '失败',
    type: success ? 'success' : 'error',
  }
})
const latestHistoryMeta = computed(() => {
  const entry = latestHistory.value
  if (!entry) return '暂无投递记录'
  const hookName = webhookMap.value.get(entry.webhookId) || `#${entry.webhookId}`
  return `${hookName} · ${formatDate(entry.createdAt)}`
})
const isTestOnly = computed({
  get: () => historyFilters.isTest === true,
  set: (value) => {
    historyFilters.isTest = value ? true : null
  },
})
const detailStatus = computed(() => {
  const entry = activeHistory.value
  if (!entry) {
    return { label: '-', type: 'default' as const }
  }
  const status = entry.responseStatus
  if (!status) {
    return { label: '未知', type: 'default' as const }
  }
  const success = status >= 200 && status < 300
  return {
    label: success ? `成功 ${status}` : `失败 ${status}`,
    type: success ? 'success' : 'error',
  }
})

const filteredWebhooks = computed(() => {
  const keyword = listFilters.keyword.trim().toLowerCase()
  return webhooks.value.filter((item) => {
    if (listFilters.status === 'enabled' && !item.isEnabled) return false
    if (listFilters.status === 'disabled' && item.isEnabled) return false
    if (listFilters.event && !item.events?.includes(listFilters.event)) return false
    if (!keyword) return true
    return item.name.toLowerCase().includes(keyword) || item.url.toLowerCase().includes(keyword)
  })
})

const columns = computed<DataTableColumns<WebhookItem>>(() => [
  {
    title: '名称',
    key: 'name',
    width: 160,
    render: (row) => h('div', { class: 'font-medium' }, row.name),
  },
  {
    title: '地址',
    key: 'url',
    render: (row) => h('div', { class: 'text-xs text-[var(--text-color-3)]' }, row.url),
  },
  {
    title: '事件',
    key: 'events',
    render: (row) => {
      if (!row.events || row.events.length === 0) {
        return '-'
      }
      return h(
        'div',
        { class: 'flex flex-wrap gap-1' },
        row.events.map((item) =>
          h(
            NTag,
            { size: 'small', type: 'info', bordered: false },
            { default: () => item },
          ),
        ),
      )
    },
  },
  {
    title: '状态',
    key: 'isEnabled',
    width: 90,
    render: (row) =>
      h(
        NTag,
        { size: 'small', type: row.isEnabled ? 'success' : 'warning', bordered: false },
        { default: () => (row.isEnabled ? '启用' : '停用') },
      ),
  },
  {
    title: '更新时间',
    key: 'updatedAt',
    width: 180,
    render: (row) => formatDate(row.updatedAt),
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    render: (row) =>
      h(
        NButtonGroup,
        { size: 'small' },
        {
          default: () => [
            h(
              NButton,
              {
                type: 'primary',
                secondary: true,
                onClick: () => openEdit(row),
              },
              { default: () => '编辑' },
            ),
            h(
              NButton,
              {
                tertiary: true,
                onClick: () => openTest(row),
              },
              { default: () => '测试' },
            ),
            h(
              NPopconfirm,
              {
                positiveText: '删除',
                negativeText: '取消',
                onPositiveClick: () => handleDelete(row),
              },
              {
                trigger: () =>
                  h(
                    NButton,
                    {
                      type: 'error',
                      secondary: true,
                    },
                    { default: () => '删除' },
                  ),
                default: () => '确认删除该 Webhook？',
              },
            ),
          ],
        },
      ),
  },
])

const historyColumns = computed<DataTableColumns<WebhookHistoryItem>>(() => [
  {
    title: '事件',
    key: 'eventName',
    width: 220,
  },
  {
    title: 'Webhook',
    key: 'webhookId',
    width: 180,
    render: (row) => webhookMap.value.get(row.webhookId) || `#${row.webhookId}`,
  },
  {
    title: '状态',
    key: 'responseStatus',
    width: 120,
    render: (row) => renderHistoryStatus(row),
  },
  {
    title: '测试',
    key: 'isTest',
    width: 90,
    render: (row) =>
      h(
        NTag,
        { size: 'small', type: row.isTest ? 'success' : 'default', bordered: false },
        { default: () => (row.isTest ? '是' : '否') },
      ),
  },
  {
    title: '时间',
    key: 'createdAt',
    width: 180,
    render: (row) => formatDate(row.createdAt),
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render: (row) =>
      h('div', { class: 'flex gap-2' }, [
        h(
          NButton,
          {
            size: 'small',
            secondary: true,
            onClick: () => openHistory(row),
          },
          { default: () => '详情' },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => handleReplay(row),
          },
          { default: () => '重放' },
        ),
      ]),
  },
])

function formatDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

function renderHistoryStatus(row: WebhookHistoryItem) {
  const success = row.responseStatus >= 200 && row.responseStatus < 300
  const label = success ? '成功' : '失败'
  return h(
    NTag,
    { size: 'small', type: success ? 'success' : 'error', bordered: false },
    { default: () => (row.responseStatus ? `${label} ${row.responseStatus}` : label) },
  )
}

function formatHeaders(headers?: Record<string, string>) {
  if (!headers || Object.keys(headers).length === 0) return '-'
  return Object.entries(headers)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n')
}

function formatBody(body?: string) {
  if (!body) return '-'
  try {
    const parsed = JSON.parse(body)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return body
  }
}

function resetForm() {
  form.name = ''
  form.url = ''
  form.events = []
  form.headers = []
  form.payloadTemplate = ''
  form.isEnabled = true
}

function ensureHeaderRow() {
  if (form.headers.length === 0) {
    form.headers.push({ key: '', value: '' })
  }
}

function openCreate() {
  editingWebhook.value = null
  resetForm()
  ensureHeaderRow()
  formDrawerVisible.value = true
}

function openEdit(item: WebhookItem) {
  editingWebhook.value = item
  form.name = item.name
  form.url = item.url
  form.events = [...(item.events || [])]
  form.payloadTemplate = item.payloadTemplate || ''
  form.isEnabled = item.isEnabled
  form.headers = Object.entries(item.headers || {}).map(([key, value]) => ({ key, value }))
  ensureHeaderRow()
  formDrawerVisible.value = true
}

function addHeaderRow() {
  form.headers.push({ key: '', value: '' })
}

function removeHeaderRow(index: number) {
  form.headers.splice(index, 1)
  ensureHeaderRow()
}

function buildHeaderPayload() {
  const headers: Record<string, string> = {}
  form.headers.forEach((row) => {
    const key = row.key.trim()
    if (!key) return
    headers[key] = row.value
  })
  return headers
}

async function fetchWebhooks() {
  loading.value = true
  try {
    webhooks.value = await listWebhooks()
  } finally {
    loading.value = false
  }
}

async function fetchEvents() {
  const result = await listWebhookEvents()
  events.value = result.events || []
}

async function fetchHistory() {
  historyLoading.value = true
  try {
    const response = await listWebhookHistory({
      page: historyPage.value,
      pageSize: historyPageSize.value,
      webhookId: historyFilters.webhookId ?? undefined,
      eventName: historyFilters.eventName ?? undefined,
      isTest: historyFilters.isTest ?? undefined,
    })
    history.value = response.items
    historyTotal.value = response.total
  } finally {
    historyLoading.value = false
  }
}

function applyHistoryFilters() {
  historyPage.value = 1
  fetchHistory()
}

function resetHistoryFilters() {
  historyFilters.webhookId = null
  historyFilters.eventName = null
  historyFilters.isTest = null
  historyPage.value = 1
  fetchHistory()
}

function resetListFilters() {
  listFilters.keyword = ''
  listFilters.status = 'all'
  listFilters.event = null
}

function handleHistoryPageChange(value: number) {
  historyPage.value = value
  fetchHistory()
}

function handleHistoryPageSizeChange(value: number) {
  historyPageSize.value = value
  historyPage.value = 1
  fetchHistory()
}

async function handleSave() {
  if (!form.name.trim()) {
    message.error('请填写名称')
    return
  }
  if (!form.url.trim()) {
    message.error('请填写 URL')
    return
  }
  if (form.events.length === 0) {
    message.error('请选择订阅事件')
    return
  }

  saving.value = true
  try {
    const payload = {
      name: form.name.trim(),
      url: form.url.trim(),
      events: form.events,
      headers: buildHeaderPayload(),
      payloadTemplate: form.payloadTemplate,
      isEnabled: form.isEnabled,
    }

    if (editingWebhook.value) {
      await updateWebhook(editingWebhook.value.id, payload)
    } else {
      await createWebhook(payload)
    }
    formDrawerVisible.value = false
    await fetchWebhooks()
  } finally {
    saving.value = false
  }
}

function handleFormatPayload() {
  try {
    form.payloadTemplate = formatTemplateJson(form.payloadTemplate)
    message.success('已格式化')
  } catch (err) {
    const reason = err instanceof Error ? err.message : 'JSON 格式不正确'
    message.error(`格式化失败：${reason}`)
  }
}

async function handleDelete(item: WebhookItem) {
  await deleteWebhook(item.id)
  fetchWebhooks()
}

function openTest(item: WebhookItem) {
  testingWebhook.value = item
  testEventName.value = item.events?.[0] || null
  testModalVisible.value = true
}

async function handleTest() {
  if (!testingWebhook.value) return
  await testWebhook(testingWebhook.value.id, testEventName.value)
  testModalVisible.value = false
  fetchHistory()
}

function openHistory(item: WebhookHistoryItem) {
  activeHistory.value = item
  historyDrawerVisible.value = true
}

async function handleReplay(item: WebhookHistoryItem) {
  await replayWebhookHistory(item.id)
  fetchHistory()
}

onMounted(async () => {
  await Promise.all([fetchWebhooks(), fetchEvents(), fetchHistory()])
})
</script>

<template>
  <ScrollContainer wrapper-class="flex flex-col gap-4">
    <NCard>
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div class="space-y-1">
          <div class="text-lg font-semibold">Webhook 管理</div>
          <div class="text-xs text-[var(--text-color-3)]">配置事件推送、测试与投递记录。</div>
        </div>
        <NSpace align="center">
          <NButton
            secondary
            @click="fetchWebhooks"
          >
            刷新
          </NButton>
          <NButton
            type="primary"
            @click="openCreate"
          >
            新建 Webhook
          </NButton>
        </NSpace>
      </div>
      <NDivider class="my-4" />
      <NGrid
        cols="1 640:2 900:4"
        x-gap="16"
        y-gap="12"
      >
        <NGi>
          <NStatistic
            label="Webhook 总数"
            tabular-nums
          >
            {{ totalWebhooks }}
          </NStatistic>
        </NGi>
        <NGi>
          <NStatistic
            label="启用中"
            tabular-nums
          >
            {{ enabledCount }}
          </NStatistic>
        </NGi>
        <NGi>
          <NStatistic
            label="已停用"
            tabular-nums
          >
            {{ disabledCount }}
          </NStatistic>
        </NGi>
        <NGi>
          <NStatistic label="最近投递">
            <NTag
              size="small"
              :bordered="false"
              :type="latestHistoryStatus.type === 'default' ? undefined : latestHistoryStatus.type"
            >
              {{ latestHistoryStatus.label }}
            </NTag>
          </NStatistic>
          <div class="mt-1 text-xs text-[var(--text-color-3)]">{{ latestHistoryMeta }}</div>
        </NGi>
      </NGrid>
    </NCard>

    <NTabs
      v-model:value="activeTab"
      type="line"
      animated
    >
      <NTabPane
        name="list"
        tab="Webhook 列表"
      >
        <NCard title="Webhook 列表">
          <template #header-extra>
            <NTag
              size="small"
              type="info"
              :bordered="false"
            >
              共 {{ filteredWebhooks.length }} 条
            </NTag>
          </template>
          <NForm
            label-placement="left"
            label-width="60"
            :show-feedback="false"
          >
            <NGrid
              cols="1 640:2 900:4"
              x-gap="16"
              y-gap="8"
            >
              <NGi>
                <NFormItem label="关键词">
                  <NInput
                    v-model:value="listFilters.keyword"
                    clearable
                    placeholder="名称 / URL"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="状态">
                  <NSelect
                    v-model:value="listFilters.status"
                    :options="statusOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="事件">
                  <NSelect
                    v-model:value="listFilters.event"
                    :options="eventOptions"
                    clearable
                    placeholder="全部"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="操作">
                  <NSpace>
                    <NButton
                      secondary
                      @click="resetListFilters"
                    >
                      重置
                    </NButton>
                  </NSpace>
                </NFormItem>
              </NGi>
            </NGrid>
          </NForm>
          <NDivider class="my-3" />
          <div
            v-if="filteredWebhooks.length === 0 && !loading"
            class="py-10"
          >
            <NEmpty description="暂无 Webhook" />
          </div>
          <NDataTable
            v-else
            :columns="columns"
            :data="filteredWebhooks"
            :loading="loading"
            :row-key="(row) => row.id"
            striped
          />
        </NCard>
      </NTabPane>

      <NTabPane
        name="history"
        tab="投递历史"
      >
        <NCard title="投递历史">
          <template #header-extra>
            <NSpace size="small">
              <NTag
                size="small"
                type="warning"
                :bordered="false"
              >
                本页失败 {{ historyFailureCount }}
              </NTag>
              <NTag
                size="small"
                type="info"
                :bordered="false"
              >
                共 {{ historyTotal }} 条
              </NTag>
            </NSpace>
          </template>
          <NForm
            label-placement="left"
            label-width="60"
            :show-feedback="false"
          >
            <NGrid
              cols="1 640:2 900:4"
              x-gap="16"
              y-gap="8"
            >
              <NGi>
                <NFormItem label="Webhook">
                  <NSelect
                    v-model:value="historyFilters.webhookId"
                    :options="webhookOptions"
                    clearable
                    placeholder="全部"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="事件">
                  <NSelect
                    v-model:value="historyFilters.eventName"
                    :options="eventOptions"
                    clearable
                    placeholder="全部"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="测试">
                  <NCheckbox v-model:checked="isTestOnly">仅测试</NCheckbox>
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="操作">
                  <NSpace>
                    <NButton
                      type="primary"
                      secondary
                      @click="applyHistoryFilters"
                    >
                      筛选
                    </NButton>
                    <NButton @click="resetHistoryFilters">重置</NButton>
                    <NButton
                      tertiary
                      @click="fetchHistory"
                    >
                      刷新
                    </NButton>
                  </NSpace>
                </NFormItem>
              </NGi>
            </NGrid>
          </NForm>
          <NDivider class="my-3" />
          <div
            v-if="history.length === 0 && !historyLoading"
            class="py-10"
          >
            <NEmpty description="暂无投递记录" />
          </div>
          <NDataTable
            v-else
            :columns="historyColumns"
            :data="history"
            :loading="historyLoading"
            :row-key="(row) => row.id"
            striped
          />
          <div class="flex justify-end pt-4">
            <NPagination
              :page="historyPage"
              :page-size="historyPageSize"
              :item-count="historyTotal"
              show-size-picker
              :page-sizes="[10, 20, 50]"
              @update:page="handleHistoryPageChange"
              @update:page-size="handleHistoryPageSizeChange"
            />
          </div>
        </NCard>
      </NTabPane>
    </NTabs>

    <NDrawer
      v-model:show="formDrawerVisible"
      placement="right"
      width="680"
    >
      <NDrawerContent
        :title="formTitle"
        closable
        header-style="padding: 20px 24px"
        body-style="padding: 0"
      >
        <ScrollContainer wrapper-class="flex flex-col gap-5">
          <NForm label-placement="top">
            <NGrid
              cols="1 640:2"
              x-gap="16"
              y-gap="12"
            >
              <NGi>
                <NFormItem label="名称">
                  <NInput
                    v-model:value="form.name"
                    placeholder="如：联合站点推送"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="URL">
                  <NInput
                    v-model:value="form.url"
                    placeholder="https://example.com/webhook"
                  />
                </NFormItem>
              </NGi>
            </NGrid>

            <NDivider>事件订阅</NDivider>
            <NFormItem label="订阅事件">
              <NCheckboxGroup v-model:value="form.events">
                <NGrid
                  cols="1 640:2"
                  x-gap="16"
                  y-gap="8"
                >
                  <NGi
                    v-for="item in events"
                    :key="item"
                  >
                    <NCheckbox
                      :value="item"
                      :label="item"
                    />
                  </NGi>
                </NGrid>
              </NCheckboxGroup>
            </NFormItem>
            <NFormItem label="启用">
              <NSwitch v-model:value="form.isEnabled" />
            </NFormItem>

            <NDivider>请求配置</NDivider>
            <NFormItem label="Headers">
              <div class="flex w-full flex-col gap-2">
                <div
                  v-for="(row, index) in form.headers"
                  :key="`${row.key}-${index}`"
                  class="flex items-center gap-2"
                >
                  <div class="w-40">
                    <NInput
                      v-model:value="row.key"
                      placeholder="Header Key"
                      class="w-full"
                    />
                  </div>
                  <div class="min-w-0 flex-1">
                    <NInput
                      v-model:value="row.value"
                      placeholder="Header Value"
                      class="w-full"
                    />
                  </div>
                  <NButton
                    tertiary
                    size="small"
                    @click="removeHeaderRow(index)"
                  >
                    删除
                  </NButton>
                </div>
                <div>
                  <NButton
                    size="small"
                    @click="addHeaderRow"
                  >
                    添加 Header
                  </NButton>
                </div>
              </div>
            </NFormItem>
            <NFormItem label="Payload 模板">
              <div class="flex w-full flex-col gap-2">
                <div class="flex justify-end">
                  <NButton
                    size="small"
                    tertiary
                    @click="handleFormatPayload"
                  >
                    格式化
                  </NButton>
                </div>
                <TemplateEditor v-model="form.payloadTemplate" />
                <NAlert
                  type="info"
                  :show-icon="false"
                >
                  可用变量：<span class="template-code" v-pre>{{.Name}}</span>、
                  <span class="template-code" v-pre>{{.OccurredAt}}</span>、
                  <span class="template-code" v-pre>{{ toJSON .Event }}</span>。
                </NAlert>
              </div>
            </NFormItem>
          </NForm>
          <div class="flex justify-end gap-2">
            <NButton @click="formDrawerVisible = false">取消</NButton>
            <NButton
              type="primary"
              :loading="saving"
              @click="handleSave"
            >
              {{ formActionLabel }}
            </NButton>
          </div>
        </ScrollContainer>
      </NDrawerContent>
    </NDrawer>

    <NModal
      v-model:show="testModalVisible"
      preset="dialog"
      title="测试 Webhook"
      positive-text="发送"
      negative-text="取消"
      @positive-click="handleTest"
    >
      <div class="flex flex-col gap-3 py-2">
        <div class="text-xs text-[var(--text-color-3)]">
          选择一个事件用于测试投递。
        </div>
        <NSelect
          v-model:value="testEventName"
          :options="eventOptions"
          placeholder="选择事件"
          clearable
        />
      </div>
    </NModal>

    <NDrawer
      v-model:show="historyDrawerVisible"
      placement="right"
      width="640"
    >
      <NDrawerContent
        title="投递详情"
        closable
        header-style="padding: 20px 24px"
        body-style="padding: 0"
      >
        <ScrollContainer wrapper-class="flex flex-col gap-4">
          <NEmpty
            v-if="!activeHistory"
            description="暂无投递详情"
          />
          <template v-else>
            <NCard title="概览">
              <NTable
                size="small"
                :bordered="false"
                :single-line="false"
              >
                <tbody>
                  <tr>
                    <th class="w-24 text-xs text-[var(--text-color-3)]">事件</th>
                    <td class="font-medium">{{ activeHistory.eventName || '-' }}</td>
                  </tr>
                  <tr>
                    <th class="text-xs text-[var(--text-color-3)]">Webhook</th>
                    <td>{{ webhookMap.get(activeHistory.webhookId) || `#${activeHistory.webhookId}` }}</td>
                  </tr>
                  <tr>
                    <th class="text-xs text-[var(--text-color-3)]">请求 URL</th>
                    <td class="break-words font-mono text-xs">{{ activeHistory.requestUrl || '-' }}</td>
                  </tr>
                  <tr>
                    <th class="text-xs text-[var(--text-color-3)]">状态</th>
                    <td>
                      <NTag
                        size="small"
                        :bordered="false"
                        :type="detailStatus.type === 'default' ? undefined : detailStatus.type"
                      >
                        {{ detailStatus.label }}
                      </NTag>
                    </td>
                  </tr>
                  <tr>
                    <th class="text-xs text-[var(--text-color-3)]">测试</th>
                    <td>
                      <NTag
                        v-if="activeHistory.isTest"
                        size="small"
                        type="warning"
                        :bordered="false"
                      >
                        是
                      </NTag>
                      <span v-else>否</span>
                    </td>
                  </tr>
                  <tr>
                    <th class="text-xs text-[var(--text-color-3)]">时间</th>
                    <td>{{ formatDate(activeHistory.createdAt) }}</td>
                  </tr>
                </tbody>
              </NTable>
            </NCard>

            <NCard>
              <NTabs
                type="segment"
                animated
              >
                <NTabPane
                  name="request"
                  tab="请求"
                >
                  <NSpace
                    vertical
                    size="large"
                  >
                    <NCard
                      size="small"
                      title="Headers"
                    >
                      <NCode
                        :code="formatHeaders(activeHistory.requestHeaders)"
                        word-wrap
                      />
                    </NCard>
                    <NCard
                      size="small"
                      title="Body"
                    >
                      <NCode
                        :code="formatBody(activeHistory.requestBody)"
                        language="json"
                        word-wrap
                      />
                    </NCard>
                  </NSpace>
                </NTabPane>
                <NTabPane
                  name="response"
                  tab="响应"
                >
                  <NSpace
                    vertical
                    size="large"
                  >
                    <NCard
                      size="small"
                      title="Headers"
                    >
                      <NCode
                        :code="formatHeaders(activeHistory.responseHeaders)"
                        word-wrap
                      />
                    </NCard>
                    <NCard
                      size="small"
                      title="Body"
                    >
                      <NCode
                        :code="formatBody(activeHistory.responseBody)"
                        language="json"
                        word-wrap
                      />
                    </NCard>
                    <NAlert
                      v-if="activeHistory.errorMessage"
                      type="error"
                      :show-icon="false"
                    >
                      {{ activeHistory.errorMessage }}
                    </NAlert>
                    <NAlert
                      v-else
                      type="info"
                      :show-icon="false"
                    >
                      无错误信息
                    </NAlert>
                  </NSpace>
                </NTabPane>
              </NTabs>
            </NCard>
          </template>
        </ScrollContainer>
      </NDrawerContent>
    </NDrawer>
  </ScrollContainer>
</template>

<style scoped>
.template-code {
  font-family: 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
    'Liberation Mono', 'Courier New', monospace;
}
</style>
