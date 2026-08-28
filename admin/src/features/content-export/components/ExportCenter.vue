<script setup lang="ts">
import { NAlert, NButton, NCard, NRadioButton, NRadioGroup, useMessage } from 'naive-ui'
import { shallowRef } from 'vue'

import { useContentExports } from '../composables/use-content-exports'

import ExportList from './ExportList.vue'

import type { ExportMode } from '../model/types'

const message = useMessage()
const { records, loading, creating, deletingId, refresh, create, remove, download } =
  useContentExports()

const mode = shallowRef<ExportMode>('both')

async function handleCreate() {
  try {
    await create(mode.value)
    message.success('内容导出已开始，可留在此页查看进度')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '创建导出任务失败')
  }
}

async function handleDelete(id: string) {
  try {
    await remove(id)
    message.success('导出已删除')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除导出失败')
  }
}

async function handleDownload(id: string) {
  try {
    await download(id)
  } catch (error) {
    message.error(error instanceof Error ? error.message : '生成下载链接失败')
  }
}
</script>

<template>
  <div class="space-y-4">
    <NAlert
      type="info"
      :bordered="false"
    >
      内容导出通过内容服务层一次性打包全部文章、手记、思考与页面（含未发布草稿），输出 Markdown +
      完整元信息，并把站内与外链图片一并下载打包，解压即可离线阅读。
    </NAlert>

    <NCard :bordered="false">
      <div class="mb-5 flex items-start justify-between gap-4 max-sm:flex-col">
        <div>
          <h2 class="text-base font-semibold">内容导出包</h2>
          <p class="mt-1 text-sm text-neutral-500 dark:text-neutral-400">
            选择布局后创建导出任务，任务在后台运行。归档内含 manifest.json、taxonomy.json、uploads/
            与内容目录。
          </p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <NRadioGroup v-model:value="mode">
            <NRadioButton value="structured">结构化</NRadioButton>
            <NRadioButton value="flatten">单文件</NRadioButton>
            <NRadioButton value="both">两者</NRadioButton>
          </NRadioGroup>
          <NButton
            secondary
            :loading="loading"
            @click="refresh()"
          >
            <template #icon><span class="iconify ph--arrows-clockwise" /></template>
            刷新
          </NButton>
          <NButton
            type="primary"
            :loading="creating"
            :disabled="
              records.some((item) => item.status === 'queued' || item.status === 'running')
            "
            @click="handleCreate"
          >
            <template #icon><span class="iconify ph--export" /></template>
            立即导出
          </NButton>
        </div>
      </div>

      <ExportList
        :records="records"
        :loading="loading"
        :deleting-id="deletingId"
        @download="handleDownload"
        @delete="handleDelete"
      />
    </NCard>
  </div>
</template>
