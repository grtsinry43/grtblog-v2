<script setup lang="ts">
import {
  NButton,
  NCard,
  NDrawer,
  NDrawerContent,
  NDynamicTags,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSwitch,
  useMessage,
} from 'naive-ui'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { ScrollContainer } from '@/components'
import MarkdownEditor from '@/components/markdown-editor/MarkdownEditor.vue'
import MarkdownPreview from '@/components/markdown-editor/MarkdownPreview.vue'
import { createArticle, getArticle, updateArticle } from '@/services/articles'
import { listCategories, listTags } from '@/services/taxonomy'

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

const categories = ref<SelectOption[]>([])
const tags = ref<SelectOption[]>([])
const dynamicTags = ref<string[]>([])
const loading = ref(false)
const saving = ref(false)
const showMeta = ref(false)

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
      toList()
    } else if (articleId.value) {
      await updateArticle(articleId.value, {
        ...payload,
        shortUrl: form.shortUrl.trim(),
      })
      toList()
    }
  } finally {
    saving.value = false
  }
}

function openMeta() {
  showMeta.value = true
}

onMounted(async () => {
  await fetchTaxonomy()
  await fetchArticle()
})

watch(dynamicTags, (newVal) => {
  // Map tag names back to IDs
  const newTagIds: number[] = []
  newVal.forEach((name) => {
    const found = tags.value.find((t) => t.label === name)
    if (found && typeof found.value === 'number') {
      newTagIds.push(found.value)
    }
  })
  form.tagIds = newTagIds
})
</script>

<template>
  <ScrollContainer wrapper-class="flex flex-col gap-y-4 px-4 py-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <NButton secondary circle @click="toList">
          <div class="iconify ph--arrow-left text-lg" />
        </NButton>
        <div class="flex flex-col">
          <span class="text-lg font-medium leading-tight">
            {{ isCreating ? '新建文章' : '编辑文章' }}
          </span>
          <span class="text-xs text-gray-500">
            {{ isCreating ? 'New Article' : form.title || 'Untitled' }}
          </span>
        </div>
      </div>
      
      <div class="flex items-center gap-3">
        <NButton secondary @click="openMeta">
          文章设置
        </NButton>
        <NButton type="primary" :loading="saving" @click="saveArticle">
          保存发布
        </NButton>
      </div>
    </div>

    <div class="editor-container h-[calc(100vh-140px)] rounded-lg border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-800">
      <div class="pane editor-pane">
        <MarkdownEditor v-model="form.content" class="h-full" />
      </div>
      <div class="divider border-r border-gray-200 dark:border-gray-700" />
      <div class="pane preview-pane">
        <MarkdownPreview :source="form.content" />
      </div>
    </div>
    <NDrawer
      v-model:show="showMeta"
      placement="right"
      width="420"
    >
      <NDrawerContent title="文章元信息">
        <NForm
          label-placement="left"
          label-width="80"
          class="space-y-4"
        >
          <NFormItem label="标题">
            <NInput
              v-model:value="form.title"
              placeholder="请输入标题"
            />
          </NFormItem>
          <NFormItem label="短链接">
            <NInput
              v-model:value="form.shortUrl"
              placeholder="比如: hello-world"
            />
          </NFormItem>
          <NFormItem label="分类">
            <NSelect
              v-model:value="form.categoryId"
              :options="categories"
              placeholder="选择分类"
              clearable
            />
          </NFormItem>
          <NFormItem label="标签">
            <NDynamicTags v-model:value="dynamicTags" />
          </NFormItem>
          <NFormItem label="摘要">
            <NInput
              v-model:value="form.summary"
              type="textarea"
              placeholder="文章摘要"
              :autosize="{ minRows: 2, maxRows: 4 }"
            />
          </NFormItem>
          <NFormItem label="导语">
            <NInput
              v-model:value="form.leadIn"
              type="textarea"
              placeholder="导语（可选）"
              :autosize="{ minRows: 2, maxRows: 4 }"
            />
          </NFormItem>
          <NFormItem label="封面">
            <NInput
              v-model:value="form.cover"
              placeholder="封面图片链接（可选）"
            />
          </NFormItem>
          <div class="grid gap-4 sm:grid-cols-2">
            <NFormItem label="发布">
              <NSwitch v-model:value="form.isPublished" />
            </NFormItem>
            <NFormItem label="置顶">
              <NSwitch v-model:value="form.isTop" />
            </NFormItem>
            <NFormItem label="热门">
              <NSwitch v-model:value="form.isHot" />
            </NFormItem>
            <NFormItem label="原创">
              <NSwitch v-model:value="form.isOriginal" />
            </NFormItem>
          </div>
        </NForm>
      </NDrawerContent>
    </NDrawer>
  </ScrollContainer>
</template>

<style scoped>
.editor-container {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 1px minmax(0, 1fr);
  width: 100%;
  overflow: hidden;
}

.pane {
  height: 100%;
  min-height: 0;
  overflow: auto;
  position: relative;
}

.preview-pane {
  padding: 0;
}

@media (max-width: 768px) {
  .editor-container {
    grid-template-columns: 1fr;
    grid-template-rows: minmax(0, 1fr) 1px minmax(0, 1fr);
  }
}

.pane::-webkit-scrollbar,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar) {
  width: 3px;
  height: 3px;
}
.pane::-webkit-scrollbar-track,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-track),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-track) {
  background: transparent;
}
.pane::-webkit-scrollbar-thumb,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb) {
  background: color-mix(
    in srgb,
    var(--primary-color),
    var(--popover-color, var(--card-color)) 85%
  );
  border-radius: 0;
  transition: background-color 160ms ease;
}
.pane::-webkit-scrollbar-thumb:hover,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb:hover),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb:hover) {
  background: color-mix(
    in srgb,
    var(--primary-color),
    var(--popover-color, var(--card-color)) 70%
  );
}
.pane {
  scrollbar-width: thin;
  scrollbar-color: color-mix(
      in srgb,
      var(--primary-color),
      var(--popover-color, var(--card-color)) 85%
    )
    transparent;
}

.editor-pane :deep(.cm-editor) {
  height: auto !important;
  min-height: 100%;
}

.editor-pane :deep(.cm-scroller) {
  height: auto !important;
  overflow: visible !important;
}

.preview-pane :deep(.markdown-preview) {
  height: auto !important;
  overflow: visible !important;
}
</style>
