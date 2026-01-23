<script setup lang="ts">
import { PreviewLink20Regular } from '@vicons/fluent'
import { PaperPlaneOutline, SaveOutline } from '@vicons/ionicons5'
import {
  NButton,
  NButtonGroup,
  NCard,
  NDivider,
  NDrawer,
  NDrawerContent,
  NDynamicTags,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NPopover,
  NSelect,
  NSwitch,
  useMessage,
  NAutoComplete,
} from 'naive-ui'
import { computed, onMounted, ref, toRef } from 'vue'

// 组件
import MarkdownEditor from '@/components/markdown-editor/MarkdownEditor.vue'
import MarkdownPreview from '@/components/markdown-editor/MarkdownPreview.vue'

// 逻辑 Hooks
import { useArticleForm } from './composables/use-article-form'
import { useEditorStats } from './composables/use-editor-stats'
import { useTaxonomySelect } from './composables/use-taxonomy-select'

defineOptions({ name: 'ArticleEdit' })

const message = useMessage()

// 1. 初始化表单核心逻辑
const { form, loading, saving, isCreating, fetch, save } = useArticleForm()

// 2. 初始化分类与标签逻辑
// 将表单中的响应式属性传给 Hook，实现双向绑定
const {
  categoryOptions,
  tagOptions,
  dynamicTags,
  tagSearchValue,
  autoCompleteOptions,
  newCatModal,
  setInitialTags,
  handleTagsChange,
  addTagFromSearch,
  createNewCategory,
} = useTaxonomySelect(toRef(form, 'tagIds'), toRef(form, 'categoryId'), message)

// 3. 初始化编辑器统计逻辑
const { cursorPos, selectionStats, statsIdle, markActivity, handleCursorChange, getStats } =
  useEditorStats()

// 4. 视图状态管理
const showMeta = ref(false)
const showPreview = ref(false)
const previewMode = ref<'markdown' | 'page'>('markdown')

// 5. 计算属性
const stats = computed(() => getStats(form.content))
const actionLabel = computed(() => {
  if (!form.isPublished) return '保存'
  return isCreating.value ? '发布' : '发布新版本'
})
const actionIcon = computed(() => (form.isPublished ? PaperPlaneOutline : SaveOutline))
const previewUrl = computed(() => {
  const slug = form.shortUrl?.trim()
  return slug ? `/posts/${slug}` : ''
})

// 6. 生命周期
onMounted(async () => {
  // 加载文章数据
  const data = await fetch()
  // 如果加载成功且有 tags 数据（带名字），初始化标签 UI
  if (data?.tags) {
    setInitialTags(data.tags)
  }
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <header
      class="z-10 flex shrink-0 flex-col gap-3 px-10 py-8 backdrop-blur sm:h-24 sm:flex-row sm:items-center sm:justify-between sm:py-0"
    >
      <div class="flex w-full items-center gap-4 sm:flex-1">
        <NInput
          v-model:value="form.title"
          placeholder="在这里开始你的写作吧..."
          :bordered="false"
          class="flex-1 text-xl! leading-tight font-bold sm:text-2xl!"
          style="--n-caret-color: var(--primary-color); background-color: transparent"
        />
      </div>

      <div class="flex w-full flex-wrap items-center gap-3 sm:w-auto sm:flex-nowrap sm:gap-4">
        <div class="flex items-baseline gap-1">
          <div class="iconify self-center ph--link-simple" />
          <span class="text-xs leading-none">/posts/</span>
          <input
            v-model="form.shortUrl"
            placeholder="请填写短链接"
            class="w-24 border-b border-current/30 p-0 pb-0.5 text-[11px] leading-none focus:border-primary focus:outline-none sm:w-32"
          />
        </div>

        <NButtonGroup>
          <NButton
            :type="!form.isPublished ? 'primary' : 'default'"
            :ghost="form.isPublished"
            @click="form.isPublished = false"
          >
            草稿
          </NButton>
          <NButton
            :type="form.isPublished ? 'primary' : 'default'"
            :ghost="!form.isPublished"
            @click="form.isPublished = true"
          >
            发布
          </NButton>
        </NButtonGroup>

        <div class="flex items-center gap-2">
          <NButton
            quaternary
            circle
            size="small"
            @click="showMeta = true"
          >
            <template #icon><div class="iconify text-xl ph--sliders-horizontal" /></template>
          </NButton>

          <NButton
            quaternary
            circle
            size="small"
            :type="showPreview ? 'primary' : 'default'"
            @click="showPreview = !showPreview"
          >
            <template #icon><PreviewLink20Regular /></template>
          </NButton>

          <NButton
            type="primary"
            size="medium"
            :loading="saving"
            @click="save"
            class="px-5 font-medium shadow-sm active:scale-95"
          >
            <template #icon><component :is="actionIcon" /></template>
            {{ actionLabel }}
          </NButton>
        </div>
      </div>
    </header>

    <main class="flex min-h-0 flex-1 overflow-hidden">
      <div
        class="editor-container grid h-full min-h-0 w-full"
        :class="showPreview ? 'grid-cols-1 lg:grid-cols-2' : 'grid-cols-1'"
      >
        <div
          class="pane editor-pane relative h-full overflow-auto"
          @scroll="markActivity"
          @wheel="markActivity"
        >
          <MarkdownEditor
            v-model="form.content"
            class="h-full min-h-full"
            @cursor-change="handleCursorChange"
          />

          <div
            class="pointer-events-none absolute right-3 bottom-3 z-10 transition-opacity duration-200"
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
                  <template #trigger
                    ><span class="cursor-help">字数 {{ stats.charCount }}</span></template
                  >
                  <div class="flex flex-col gap-0.5 text-[11px] leading-tight">
                    <span v-if="selectionStats.total">选中 {{ selectionStats.chars }}</span>
                    <span>中文 {{ stats.chineseCharCount }}</span>
                    <span>英文词 {{ stats.wordCount }}</span>
                    <span>字符 {{ stats.totalCharCount }}</span>
                    <span>段落 {{ stats.paragraphCount }}</span>
                  </div>
                </NPopover>
                <span v-if="selectionStats.total">选中 {{ selectionStats.chars }} 字</span>
                <span>{{ cursorPos.line }}:{{ cursorPos.column }}</span>
                <span>预计阅读 {{ stats.readingMinutes }} 分钟</span>
              </div>
            </NCard>
          </div>
        </div>

        <div
          v-if="showPreview"
          class="pane preview-pane relative h-full overflow-auto"
          @scroll="markActivity"
        >
          <div class="absolute top-3 right-3 z-10">
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
                  <template #icon><div class="iconify text-lg ph--dots-three-vertical" /></template>
                </NButton>
              </template>
              <div class="flex flex-col gap-1 p-1">
                <NButton
                  :type="previewMode === 'markdown' ? 'primary' : 'default'"
                  quaternary
                  size="small"
                  class="w-full justify-start px-2"
                  @click="previewMode = 'markdown'"
                  >Markdown 预览</NButton
                >
                <NButton
                  :type="previewMode === 'page' ? 'primary' : 'default'"
                  quaternary
                  size="small"
                  class="w-full justify-start px-2"
                  @click="previewMode = 'page'"
                  >网页预览</NButton
                >
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
              class="h-full w-full border-0"
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
                <div class="flex w-full items-center gap-2">
                  <NSelect
                    v-model:value="form.categoryId"
                    :options="categoryOptions"
                    placeholder="选择分类"
                    clearable
                    filterable
                    class="flex-1"
                  />
                  <NButton
                    quaternary
                    size="small"
                    @click="newCatModal.show = true"
                    >新建</NButton
                  >
                </div>
              </NFormItem>
              <NFormItem
                label="标签"
                :show-feedback="false"
              >
                <div class="flex w-full flex-col gap-2">
                  <NDynamicTags
                    :value="dynamicTags"
                    @update:value="handleTagsChange"
                  />
                  <div class="flex items-center gap-2">
                    <NAutoComplete
                      v-model:value="tagSearchValue"
                      :options="autoCompleteOptions"
                      placeholder="搜索或创建标签"
                      class="flex-1"
                      @select="addTagFromSearch"
                      :input-props="{
                        onKeydown: (e: KeyboardEvent) => {
                          if (e.key === 'Enter') addTagFromSearch(tagSearchValue)
                        },
                      }"
                    />
                    <NButton
                      quaternary
                      size="small"
                      @click="addTagFromSearch(tagSearchValue)"
                      >添加</NButton
                    >
                  </div>
                </div>
              </NFormItem>
            </NForm>
          </div>

          <NDivider style="margin: 0" />

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
                  <template #prefix><div class="iconify ph--image" /></template>
                </NInput>
              </NFormItem>
            </NForm>
          </div>

          <NDivider style="margin: 0" />

          <div class="space-y-4">
            <div class="flex items-center gap-2 text-sm font-medium">
              <div class="iconify ph--toggle-left" />
              <span>属性</span>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div
                class="flex items-center justify-between rounded-lg px-4 py-3"
              >
                <span class="text-sm">置顶</span
                ><NSwitch
                  v-model:value="form.isTop"
                  size="small"
                />
              </div>
              <div
                class="flex items-center justify-between rounded-lg  px-4 py-3"
              >
                <span class="text-sm">热门</span
                ><NSwitch
                  v-model:value="form.isHot"
                  size="small"
                />
              </div>
              <div
                class="flex items-center justify-between rounded-lg px-4 py-3"
              >
                <span class="text-sm">原创</span
                ><NSwitch
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
      v-model:show="newCatModal.show"
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
            <NInput
              v-model:value="newCatModal.name"
              placeholder="例如：随笔"
            />
          </NFormItem>
          <NFormItem
            label="短链接"
            :show-feedback="false"
          >
            <NInput
              v-model:value="newCatModal.slug"
              placeholder="例如：notes"
            />
          </NFormItem>
        </NForm>
        <div class="mt-4 flex justify-end gap-2">
          <NButton
            quaternary
            @click="newCatModal.show = false"
            >取消</NButton
          >
          <NButton
            type="primary"
            :loading="newCatModal.loading"
            @click="createNewCategory"
            >创建并选择</NButton
          >
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
