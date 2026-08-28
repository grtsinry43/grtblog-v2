<script setup lang="ts">
import { NButton, NEmpty, NPopconfirm, NSpin, NTag, NTooltip } from 'naive-ui'

import { formatDate } from '@/utils/format'

import type { ExportMode, ExportRecord, ExportStatus } from '../model/types'

defineProps<{
  records: ExportRecord[]
  loading: boolean
  deletingId: string | null
}>()

const emit = defineEmits<{
  download: [id: string]
  delete: [id: string]
}>()

const statusMeta: Record<
  ExportStatus,
  { label: string; type: 'default' | 'info' | 'success' | 'error' }
> = {
  queued: { label: '等待中', type: 'default' },
  running: { label: '导出中', type: 'info' },
  completed: { label: '已完成', type: 'success' },
  failed: { label: '失败', type: 'error' },
}

const modeLabels: Record<ExportMode, string> = {
  structured: '结构化',
  flatten: '单文件',
  both: '两者',
}

const stageLabels: Record<string, string> = {
  queued: '等待执行',
  collecting_content: '采集内容',
  resolving_images: '处理站内图片',
  downloading_external: '下载外链图片',
  packing_archive: '打包归档',
  completed: '完成',
  failed: '失败',
  interrupted: '被服务重启中断',
}

function formatSize(size: number) {
  if (!size) return '—'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const index = Math.min(Math.floor(Math.log(size) / Math.log(1024)), units.length - 1)
  return `${(size / 1024 ** index).toFixed(index === 0 ? 0 : 1)} ${units[index]}`
}
</script>

<template>
  <NSpin :show="loading">
    <NEmpty
      v-if="!records.length && !loading"
      description="还没有内容导出，创建第一份内容包吧"
      class="py-16"
    />
    <div
      v-else
      class="space-y-3"
    >
      <article
        v-for="item in records"
        :key="item.id"
        class="rounded-xl border border-neutral-200 bg-white p-4 transition-shadow hover:shadow-sm dark:border-neutral-700 dark:bg-neutral-900"
      >
        <div class="flex gap-4 max-md:flex-col">
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="iconify text-xl text-[var(--primary-color)] ph--package" />
              <h3 class="truncate font-medium">{{ item.filename }}</h3>
              <NTag
                size="small"
                :type="statusMeta[item.status].type"
                :bordered="false"
              >
                {{ statusMeta[item.status].label }}
              </NTag>
              <NTag
                size="small"
                :bordered="false"
              >
                {{ modeLabels[item.mode] || item.mode }}
              </NTag>
            </div>

            <div
              class="mt-3 grid grid-cols-2 gap-x-8 gap-y-2 text-xs text-neutral-500 md:grid-cols-4 dark:text-neutral-400"
            >
              <div>
                <span class="block text-neutral-400">创建时间</span>{{ formatDate(item.createdAt) }}
              </div>
              <div>
                <span class="block text-neutral-400">归档大小</span>{{ formatSize(item.sizeBytes) }}
              </div>
              <div>
                <span class="block text-neutral-400">内容条目</span>
                {{ item.articleCount + item.momentsCount + item.pagesCount + item.thinkingsCount }}
                条
              </div>
              <div>
                <span class="block text-neutral-400">打包图片</span>
                <NTooltip v-if="item.failedImageCount > 0">
                  <template #trigger>
                    <span class="cursor-help text-orange-500">
                      {{ item.imageCount }} 张（{{ item.failedImageCount }} 张失败）
                    </span>
                  </template>
                  部分外链图片下载失败，对应位置保留了原始链接，详见归档内 manifest.json
                </NTooltip>
                <span v-else>{{ item.imageCount }} 张</span>
              </div>
            </div>

            <div
              v-if="item.status === 'queued' || item.status === 'running'"
              class="mt-3 flex items-center gap-2 text-xs text-blue-600 dark:text-blue-400"
            >
              <span class="iconify animate-spin ph--spinner" />
              {{ stageLabels[item.stage] || item.stage }}
            </div>
            <div
              v-if="item.errorMessage"
              class="mt-3 rounded-md bg-red-50 px-3 py-2 text-xs text-red-600 dark:bg-red-950/30 dark:text-red-400"
            >
              {{ item.errorMessage }}
            </div>
          </div>

          <div class="flex shrink-0 items-start gap-2">
            <NButton
              size="small"
              type="primary"
              secondary
              :disabled="item.status !== 'completed'"
              @click="emit('download', item.id)"
            >
              <template #icon><span class="iconify ph--download-simple" /></template>
              下载
            </NButton>
            <NPopconfirm
              :disabled="item.status === 'queued' || item.status === 'running'"
              @positive-click="emit('delete', item.id)"
            >
              <template #trigger>
                <NButton
                  size="small"
                  tertiary
                  type="error"
                  :loading="deletingId === item.id"
                  :disabled="item.status === 'queued' || item.status === 'running'"
                >
                  <template #icon><span class="iconify ph--trash" /></template>
                </NButton>
              </template>
              删除后归档文件也会被永久移除，确认继续？
            </NPopconfirm>
          </div>
        </div>
      </article>
    </div>
  </NSpin>
</template>
