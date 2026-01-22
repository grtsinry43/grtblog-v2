<script setup lang="ts">
import { PreviewLink20Regular } from '@vicons/fluent'
import { PaperPlaneOutline, SaveOutline } from '@vicons/ionicons5'
import {
  NButton,
  NCard,
  NDivider,
  NDrawer,
  NDrawerContent,
  NDynamicTags,
  NForm,
  NFormItem,
  NAutoComplete,
  NInput,
  NModal,
  NPopover,
  NSelect,
  NSwitch,
  useMessage,
  NButtonGroup,
} from 'naive-ui'
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import MarkdownEditor from '@/components/markdown-editor/MarkdownEditor.vue'
import MarkdownPreview from '@/components/markdown-editor/MarkdownPreview.vue'
import { useLeaveConfirm } from '@/composables'
import { createArticle, getArticle, updateArticle } from '@/services/articles'
import { createCategory, createTag, listCategories, listTags } from '@/services/taxonomy'

import type { SelectOption } from 'naive-ui'

defineOptions({
  name: 'ArticleEdit',
})

const route = useRoute()
const router = useRouter()
const message = useMessage()

const isCreating = computed(() => route.name === 'articleCreate')

const form = reactive({
  title: '',
  summary: '',
  leadIn: '',
  content: '',
  cover: '',
  categoryId: null as number | null,
  tagIds: [] as number[],
  shortUrl: '',
  isPublished: false,
  isTop: false,
  isHot: false,
  isOriginal: true,
})

const actionLabel = computed(() => {
  if (!form.isPublished) return '保存'
  return isCreating.value ? '发布' : '发布新版本'
})
const actionIcon = computed(() => (form.isPublished ? PaperPlaneOutline : SaveOutline))

const categories = ref<SelectOption[]>([])
const tags = ref<SelectOption[]>([])
const dynamicTags = ref<string[]>([])
const tagSearch = ref('')
const tagsReady = ref(false)
const syncingTags = ref(false)
const loading = ref(false)
const saving = ref(false)
const showMeta = ref(false)
const showPreview = ref(false)
const previewMode = ref<'markdown' | 'page'>('markdown')
const initialSnapshot = ref('')
const showCategoryModal = ref(false)
const creatingCategory = ref(false)
const newCategory = reactive({
  name: '',
  shortUrl: '',
})

const previewUrl = computed(() => {
  const slug = form.shortUrl.trim()
  if (!slug) return ''
  return `/posts/${slug}`
})

const cursorPos = ref({ line: 1, column: 1 })
const selectionStats = ref({ chars: 0, total: 0 })
const charCount = computed(() => form.content.replace(/\s/g, '').length)
const totalCharCount = computed(() => form.content.length)
const chineseCharCount = computed(() => (form.content.match(/[\u4e00-\u9fff]/g) ?? []).length)
const wordCount = computed(() => (form.content.match(/[A-Za-z0-9]+/g) ?? []).length)
const whitespaceCount = computed(() => (form.content.match(/\s/g) ?? []).length)
const paragraphCount = computed(() => {
  return form.content
    .split(/\n+/)
    .map((line) => line.trim())
    .filter(Boolean).length
})
const lineCount = computed(() => {
  if (!form.content) return 0
  return form.content.split('\n').length
})
const headingCount = computed(() => {
  return form.content
    .split('\n')
    .filter((line) => /^#{1,6}\s+\S+/.test(line.trim())).length
})
const tagOptions = computed(() =>
  tags.value.map((item) => ({ label: item.label as string, value: item.label as string })),
)
const readingMinutes = computed(() => {
  if (!charCount.value) return 0
  return Math.max(1, Math.ceil(charCount.value / 300))
})
const statsIdle = ref(true)
let statsIdleTimer: ReturnType<typeof setTimeout> | null = null

const isDirty = computed(() => {
  const snapshot = JSON.stringify({
    title: form.title,
    summary: form.summary,
    leadIn: form.leadIn,
    content: form.content,
    cover: form.cover,
    categoryId: form.categoryId,
    tagIds: form.tagIds,
    shortUrl: form.shortUrl,
    isPublished: form.isPublished,
    isTop: form.isTop,
    isHot: form.isHot,
    isOriginal: form.isOriginal,
  })
  return initialSnapshot.value !== '' && snapshot !== initialSnapshot.value
})

const articleId = computed(() => {
  if (isCreating.value) return null
  const raw = route.params.id
  const id = typeof raw === 'string' ? Number(raw) : Array.isArray(raw) ? Number(raw[0]) : NaN
  return Number.isFinite(id) ? id : null
})

function resetForm() {
  form.title = ''
  form.summary = ''
  form.leadIn = ''
  form.content = ''
  form.cover = ''
  form.categoryId = null
  form.tagIds = []
  dynamicTags.value = []
  form.shortUrl = ''
  form.isPublished = false
  form.isTop = false
  form.isHot = false
  form.isOriginal = true
}

async function fetchTaxonomy() {
  const [categoryList, tagList] = await Promise.all([listCategories(), listTags()])
  categories.value = categoryList.map((item) => ({
    label: item.name,
    value: item.id,
  }))
  tags.value = tagList.map((item) => ({
    label: item.name,
    value: item.id,
  }))
  tagsReady.value = true
}

async function fetchArticle() {
  if (isCreating.value) {
    resetForm()
    return
  }
  if (!articleId.value) {
    message.error('无效的文章ID')
    router.replace({ name: 'articleList' })
    return
  }
  loading.value = true
  try {
    const result = await getArticle(articleId.value)
    form.title = result.title
    form.summary = result.summary || ''
    form.leadIn = result.leadIn || ''
    form.content = result.content
    form.cover = result.cover || ''
    form.categoryId = result.categoryId ?? null
    form.tagIds = result.tags?.map((tag) => tag.id) ?? []
    dynamicTags.value = result.tags?.map((tag) => tag.name) ?? []
    form.shortUrl = result.shortUrl
    form.isPublished = result.isPublished
    form.isTop = result.isTop
    form.isHot = result.isHot
    form.isOriginal = result.isOriginal
  } finally {
    loading.value = false
  }
}

function captureSnapshot() {
  initialSnapshot.value = JSON.stringify({
    title: form.title,
    summary: form.summary,
    leadIn: form.leadIn,
    content: form.content,
    cover: form.cover,
    categoryId: form.categoryId,
    tagIds: form.tagIds,
    shortUrl: form.shortUrl,
    isPublished: form.isPublished,
    isTop: form.isTop,
    isHot: form.isHot,
    isOriginal: form.isOriginal,
  })
}

function toList() {
  router.push({ name: 'articleList' })
}

function normalizeText(value: string) {
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

async function saveArticle() {
  if (!form.title.trim()) {
    message.error('请输入标题')
    return
  }
  if (!form.content.trim()) {
    message.error('请输入正文内容')
    return
  }
  if (!isCreating.value && !form.shortUrl.trim()) {
    message.error('短链接不能为空')
    return
  }

  saving.value = true
  try {
    const payload = {
      title: form.title.trim(),
      summary: form.summary.trim(),
      leadIn: normalizeText(form.leadIn),
      content: form.content,
      cover: normalizeText(form.cover),
      categoryId: form.categoryId ?? null,
      tagIds: form.tagIds,
      shortUrl: normalizeText(form.shortUrl),
      isPublished: form.isPublished,
      isTop: form.isTop,
      isHot: form.isHot,
      isOriginal: form.isOriginal,
    }

    if (isCreating.value) {
      await createArticle(payload)
      captureSnapshot()
      toList()
    } else if (articleId.value) {
      await updateArticle(articleId.value, {
        ...payload,
        shortUrl: form.shortUrl.trim(),
      })
      captureSnapshot()
      toList()
    }
  } finally {
    saving.value = false
  }
}

function openMeta() {
  showMeta.value = true
}

function togglePreview() {
  showPreview.value = !showPreview.value
}

function setPreviewMode(mode: 'markdown' | 'page') {
  previewMode.value = mode
}

function updateCursor(value: {
  line: number
  column: number
  selectionChars: number
  selectionTotal: number
}) {
  cursorPos.value = { line: value.line, column: value.column }
  selectionStats.value = { chars: value.selectionChars, total: value.selectionTotal }
}

async function syncTags(next: string[]) {
  if (!tagsReady.value || syncingTags.value) return
  syncingTags.value = true
  try {
    const names = next.map((name) => name.trim()).filter(Boolean)
    if (!names.length) {
      form.tagIds = []
      return
    }
    const existingIds: number[] = []
    const toCreate: string[] = []
    names.forEach((name) => {
      const found = tags.value.find((option) => option.label === name)
      if (found && typeof found.value === 'number') {
        existingIds.push(found.value)
      } else {
        toCreate.push(name)
      }
    })

    const createdIds: number[] = []
    for (const name of toCreate) {
      const created = await createTag(name)
      tags.value = [...tags.value, { label: created.name, value: created.id }]
      createdIds.push(created.id)
    }

    form.tagIds = Array.from(new Set([...existingIds, ...createdIds]))
  } finally {
    syncingTags.value = false
  }
}

function addTagFromSearch(value: string) {
  if (!value) return
  if (!dynamicTags.value.includes(value)) {
    dynamicTags.value = [...dynamicTags.value, value]
  }
  tagSearch.value = ''
}

function openCategoryModal() {
  newCategory.name = ''
  newCategory.shortUrl = ''
  showCategoryModal.value = true
}

async function handleCreateCategory() {
  const name = newCategory.name.trim()
  const shortUrl = newCategory.shortUrl.trim()
  if (!name) {
    message.error('请输入分类名称')
    return
  }
  if (!shortUrl) {
    message.error('请输入分类短链接')
    return
  }
  creatingCategory.value = true
  try {
    const created = await createCategory({ name, shortUrl })
    const option = { label: created.name, value: created.id }
    categories.value = [...categories.value, option]
    form.categoryId = created.id
    message.success('分类已创建')
    showCategoryModal.value = false
  } finally {
    creatingCategory.value = false
  }
}

async function handleTagSelection(next: Array<number | string>) {
  if (!next.length) {
    form.tagIds = []
    tagSelection.value = []
    return
  }

  const existing = next.filter((val) => typeof val === 'number') as number[]
  const pending = next.filter((val) => typeof val === 'string') as string[]

  if (!pending.length) {
    form.tagIds = existing
    tagSelection.value = existing
    return
  }

  const createdIds: number[] = []
  for (const name of pending) {
    const trimmed = name.trim()
    if (!trimmed) continue
    const existed = tags.value.find((option) => option.label === trimmed)
    if (existed && typeof existed.value === 'number') {
      createdIds.push(existed.value)
      continue
    }
    const created = await createTag(trimmed)
    tags.value = [...tags.value, { label: created.name, value: created.id }]
    createdIds.push(created.id)
  }

  const merged = Array.from(new Set([...existing, ...createdIds]))
  form.tagIds = merged
  tagSelection.value = merged
}

function markStatsActivity() {
  statsIdle.value = false
  if (statsIdleTimer) clearTimeout(statsIdleTimer)
  statsIdleTimer = setTimeout(() => {
    statsIdle.value = true
  }, 800)
}

function setDraftMode() {
  form.isPublished = false
}

function setPublishMode() {
  form.isPublished = true
}

onMounted(async () => {
  await fetchTaxonomy()
  await fetchArticle()
  captureSnapshot()
  await syncTags(dynamicTags.value)
})

onBeforeUnmount(() => {
  if (statsIdleTimer) clearTimeout(statsIdleTimer)
})

watch(dynamicTags, (next) => {
  syncTags(next)
})


useLeaveConfirm({
  when: isDirty,
  title: '未保存的更改',
  content: '当前内容未保存，确定要离开吗？',
  positiveText: '离开',
  negativeText: '继续编辑',
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <!-- Toolbar -->
    <header
      class="z-10 flex shrink-0 flex-col gap-3 px-4 py-3 backdrop-blur sm:h-16 sm:flex-row sm:items-center sm:justify-between sm:py-0"
    >
      <!-- Left: Title -->
      <div class="flex w-full items-center gap-4 sm:flex-1">
        <NInput
          v-model:value="form.title"
          placeholder="在这里开始你的写作吧..."
          :bordered="false"
          class="flex-1 text-xl! leading-tight font-bold sm:text-2xl!"
          style="
            --n-caret-color: var(--primary-color);
            --n-placeholder-color: inherit;
            --n-border: none;
            --n-box-shadow-focus: none;
            background-color: transparent;
            padding: 0;
          "
        />
      </div>

      <!-- Right: Meta & Actions -->
      <div class="flex w-full flex-wrap items-center gap-3 sm:w-auto sm:flex-nowrap sm:gap-4">
        <div class="flex items-baseline gap-1">
          <div class="iconify ph--link-simple self-center" />
          <span class="text-xs leading-none">/posts/</span>
          <input
            v-model="form.shortUrl"
            placeholder="slug"
            class="w-24 border-b border-current/30 p-0 pb-0.5 text-[11px] leading-none focus:border-[var(--primary-color)] focus:outline-none sm:w-32"
          />
        </div>

          <n-button-group>
            <n-button
              :type="!form.isPublished ? 'primary' : 'default'"
              :ghost="form.isPublished"
              @click="setDraftMode"
            >
              草稿
            </n-button>
            <n-button
              :type="form.isPublished ? 'primary' : 'default'"
              :ghost="!form.isPublished"
              @click="setPublishMode"
            >
              发布
            </n-button>
          </n-button-group>

        <div class="flex items-center gap-2">
          <NButton
            quaternary
            circle
            size="small"
            @click="openMeta"
          >
            <template #icon>
              <div
                class="iconify text-xl ph--sliders-horizontal"
              />
            </template>
          </NButton>

          <NButton
            quaternary
            circle
            size="small"
            :type="showPreview ? 'primary' : 'default'"
            @click="togglePreview"
          >
            <template #icon>
              <PreviewLink20Regular />
            </template>
          </NButton>

          <NButton
            type="primary"
            size="medium"
            :loading="saving"
            @click="saveArticle"
            class="px-5 font-medium shadow-sm active:scale-95"
          >
            <template #icon>
              <component :is="actionIcon" />
            </template>
            {{ actionLabel }}
          </NButton>
        </div>
      </div>
    </header>

    <!-- Main Workspace -->
    <main class="flex min-h-0 flex-1 overflow-hidden">
      <!-- Editor Container -->
      <div
        class="editor-container grid h-full min-h-0 w-full"
        :class="showPreview ? 'grid-cols-1 lg:grid-cols-2' : 'grid-cols-1'"
      >
        <div
          class="pane editor-pane relative h-full overflow-auto"
          @scroll="markStatsActivity"
          @wheel="markStatsActivity"
        >
          <MarkdownEditor
            v-model="form.content"
            class="h-full min-h-full"
            @cursor-change="updateCursor"
          />
          <div
            class="pointer-events-none absolute bottom-3 right-3 z-10 transition-opacity duration-200"
            :class="statsIdle ? 'opacity-75 hover:opacity-100' : 'opacity-0'"
          >
            <NCard
              size="small"
              class="pointer-events-auto shadow-sm"
              content-style="padding: 6px 8px;"
            >
              <div class="flex items-center gap-3 text-[13px]">
                <NPopover
                  trigger="hover"
                  :disabled="!statsIdle"
                  content-style="padding: 4px 6px;"
                >
                  <template #trigger>
                    <span class="cursor-help">字数 {{ charCount }}</span>
                  </template>
                  <div class="flex flex-col gap-0.5 text-[11px] leading-tight">
                    <span v-if="selectionStats.total">选中 {{ selectionStats.chars }}</span>
                    <span>中文 {{ chineseCharCount }}</span>
                    <span>英文词 {{ wordCount }}</span>
                    <span>字符 {{ totalCharCount }}</span>
                    <span>非空白 {{ charCount }}</span>
                    <span>空白 {{ whitespaceCount }}</span>
                    <span>段落 {{ paragraphCount }}</span>
                    <span>行数 {{ lineCount }}</span>
                    <span>标题 {{ headingCount }}</span>
                  </div>
                </NPopover>
                <span v-if="selectionStats.total">选中 {{ selectionStats.chars }} 字</span>
                <span>{{ cursorPos.line }}:{{ cursorPos.column }}</span>
                <span>预计阅读 {{ readingMinutes }} 分钟</span>
              </div>
            </NCard>
          </div>
        </div>
        <div
          v-if="showPreview"
          class="pane preview-pane relative h-full overflow-auto"
          @scroll="markStatsActivity"
          @wheel="markStatsActivity"
        >
          <div class="absolute right-3 top-3 z-10">
            <NPopover
              trigger="click"
              placement="bottom-end"
            >
              <template #trigger>
                <NButton
                  tertiary
                  type="primary"
                  circle
                  size="small"
                  class="shadow-sm"
                >
                  <template #icon>
                    <div class="iconify text-lg ph--dots-three-vertical" />
                  </template>
                </NButton>
              </template>
              <div class="flex flex-col gap-1 p-1">
                <n-button
                  :type="previewMode === 'markdown' ? 'primary' : 'default'"
                  quaternary
                  size="small"
                  class="w-full justify-start px-2"
                  @click="setPreviewMode('markdown')"
                >
                  Markdown 预览
                </n-button>
                <n-button
                  :type="previewMode === 'page' ? 'primary' : 'default'"
                  quaternary
                  size="small"
                  class="w-full justify-start px-2"
                  @click="setPreviewMode('page')"
                >
                  网页预览
                </n-button>
              </div>
            </NPopover>
          </div>
          <MarkdownPreview
            v-if="previewMode === 'markdown'"
            :source="form.content"
            class="p-4 sm:p-8"
          />
          <div
            v-else
            class="h-full w-full"
          >
            <iframe
              v-if="previewUrl"
              :src="previewUrl"
              class="h-full w-full"
            />
            <div
              v-else
              class="flex h-full items-center justify-center text-sm opacity-60"
            >
              请先填写 slug
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Settings Drawer -->
    <NDrawer
      v-model:show="showMeta"
      placement="right"
      width="400"
    >
      <NDrawerContent
        title="文章设置"
        :native-scrollbar="false"
        closable
        header-style="padding: 24px;"
        body-style="padding: 24px;"
      >
        <div class="flex flex-col gap-6">
          <!-- Classification Section -->
          <div class="space-y-4">
            <div class="flex items-center gap-2 text-sm font-medium">
              <div class="iconify ph--tag" />
              <span>分类与标签</span>
            </div>
            <NForm
              label-placement="top"
              label-width="auto"
              class="space-y-4"
            >
                <NFormItem
                  label="分类"
                  :show-feedback="false"
                >
                  <div class="flex items-center gap-2">
                    <NSelect
                      v-model:value="form.categoryId"
                      :options="categories"
                      placeholder="选择分类"
                      clearable
                      filterable
                      class="flex-1"
                    />
                    <NButton
                      quaternary
                      size="small"
                      @click="openCategoryModal"
                    >
                      新建
                    </NButton>
                  </div>
                </NFormItem>
                <NFormItem
                  label="标签"
                  :show-feedback="false"
                >
                  <div class="flex flex-col gap-2">
                    <NDynamicTags v-model:value="dynamicTags" />
                    <div class="flex items-center gap-2">
                      <NAutoComplete
                        v-model:value="tagSearch"
                        :options="tagOptions"
                        placeholder="搜索或创建标签"
                        class="flex-1"
                        @select="addTagFromSearch"
                        :input-props="{
                          onKeydown: (e: KeyboardEvent) => {
                            if (e.key === 'Enter') addTagFromSearch(tagSearch)
                          },
                        }"
                      />
                      <NButton
                        quaternary
                        size="small"
                        @click="addTagFromSearch(tagSearch)"
                      >
                        添加
                      </NButton>
                    </div>
                  </div>
                </NFormItem>
              </NForm>
          </div>

          <NDivider style="margin: 0" />

          <!-- Meta Info Section -->
          <div class="space-y-4">
            <div class="flex items-center gap-2 text-sm font-medium">
              <div class="iconify ph--article" />
              <span>元信息</span>
            </div>
            <NForm
              label-placement="top"
              label-width="auto"
              class="space-y-4"
            >
              <NFormItem
                label="摘要"
                :show-feedback="false"
              >
                <NInput
                  v-model:value="form.summary"
                  type="textarea"
                  placeholder="简短的摘要..."
                  :autosize="{ minRows: 2, maxRows: 4 }"
                />
              </NFormItem>
              <NFormItem
                label="封面图"
                :show-feedback="false"
              >
                <NInput
                  v-model:value="form.cover"
                  placeholder="图片 URL"
                >
                  <template #prefix>
                      <div class="iconify ph--image" />
                    </template>
                  </NInput>
                </NFormItem>
            </NForm>
          </div>

          <NDivider style="margin: 0" />

          <!-- Attributes Section -->
          <div class="space-y-4">
            <div class="flex items-center gap-2 text-sm font-medium">
              <div class="iconify ph--toggle-left" />
              <span>属性</span>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="flex items-center justify-between rounded-lg px-4 py-3">
                <span class="text-sm">置顶</span>
                <NSwitch
                  v-model:value="form.isTop"
                  size="small"
                />
              </div>
              <div class="flex items-center justify-between rounded-lg px-4 py-3">
                <span class="text-sm">热门</span>
                <NSwitch
                  v-model:value="form.isHot"
                  size="small"
                />
              </div>
              <div class="flex items-center justify-between rounded-lg px-4 py-3">
                <span class="text-sm">原创</span>
                <NSwitch
                  v-model:value="form.isOriginal"
                  size="small"
                />
              </div>
            </div>
          </div>
        </div>
      </NDrawerContent>
    </NDrawer>

    <NModal
      v-model:show="showCategoryModal"
      style="width: 420px; max-width: 90vw"
    >
      <NCard
        title="新建分类"
        size="small"
      >
        <NForm
          label-placement="top"
          label-width="auto"
          class="space-y-3"
        >
          <NFormItem
            label="名称"
            :show-feedback="false"
          >
            <NInput v-model:value="newCategory.name" placeholder="例如：随笔" />
          </NFormItem>
          <NFormItem
            label="短链接"
            :show-feedback="false"
          >
            <NInput v-model:value="newCategory.shortUrl" placeholder="例如：notes" />
          </NFormItem>
        </NForm>
        <div class="mt-4 flex justify-end gap-2">
          <NButton
            quaternary
            @click="showCategoryModal = false"
          >
            取消
          </NButton>
          <NButton
            type="primary"
            :loading="creatingCategory"
            @click="handleCreateCategory"
          >
            创建并选择
          </NButton>
        </div>
      </NCard>
    </NModal>
  </div>
</template>

<style scoped>
.editor-container {
  /* Clean grid layout handled by Tailwind classes */
}

/* Custom scrollbar refinements for a cleaner look */
.pane::-webkit-scrollbar,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar) {
  width: 5px;
  height: 5px;
}

.pane::-webkit-scrollbar-track,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-track),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-track) {
  background: transparent;
}

:global(.dark) .pane::-webkit-scrollbar-thumb,
:global(.dark) .editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb),
:global(.dark) .preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb) {
  background-color: #374151;
}

.pane::-webkit-scrollbar-thumb:hover,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb:hover),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb:hover) {
  background-color: #d1d5db;
}

:global(.dark) .pane::-webkit-scrollbar-thumb:hover,
:global(.dark) .editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb:hover),
:global(.dark) .preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb:hover) {
  background-color: #4b5563;
}

.editor-pane :deep(.cm-editor) {
  height: 100% !important;
  font-family: inherit;
}

.editor-pane :deep(.cm-scroller) {
  padding-bottom: 50vh; /* Allow scrolling past end */
  font-family: 'JetBrains Mono', monospace; /* Optional: technical font for code */
  line-height: 1.6;
}

.preview-pane :deep(.markdown-preview) {
  padding-bottom: 50vh;
}
</style>
