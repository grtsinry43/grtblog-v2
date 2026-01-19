<script setup lang="ts">
import { NButton, NCard, NDataTable, NPagination, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { ScrollContainer } from '@/components'
import { listArticles } from '@/services/articles'

import type { ArticleListItem } from '@/services/articles'
import type { DataTableColumns } from 'naive-ui'

defineOptions({
  name: 'ArticleList',
})

const router = useRouter()

const list = ref<ArticleListItem[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const columns: DataTableColumns<ArticleListItem> = [
  {
    title: '标题',
    key: 'title',
    width: 260,
    render: (row) => h('div', { class: 'font-medium' }, row.title),
  },
  {
    title: '分类',
    key: 'categoryName',
    width: 140,
    render: (row) => row.categoryName || '-',
  },
  {
    title: '标签',
    key: 'tags',
    render: (row) => {
      if (!row.tags || row.tags.length === 0) {
        return '-'
      }
      return h(
        'div',
        { class: 'flex flex-wrap gap-1' },
        row.tags.map((tag) =>
          h(
            NTag,
            { size: 'small', type: 'info', bordered: false },
            { default: () => tag },
          ),
        ),
      )
    },
  },
  {
    title: '数据',
    key: 'metrics',
    width: 160,
    render: (row) => `${row.views} / ${row.likes} / ${row.comments}`,
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
    width: 120,
    render: (row) =>
      h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => toEdit(row.id),
        },
        { default: () => '编辑' },
      ),
  },
]

function formatDate(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

async function fetchList() {
  loading.value = true
  try {
    const response = await listArticles({
      page: page.value,
      pageSize: pageSize.value,
    })
    list.value = response.items
    total.value = response.total
  } finally {
    loading.value = false
  }
}

function handlePageChange(value: number) {
  page.value = value
  fetchList()
}

function handlePageSizeChange(value: number) {
  pageSize.value = value
  page.value = 1
  fetchList()
}

function toCreate() {
  router.push({ name: 'articleCreate' })
}

function toEdit(id: number) {
  router.push({ name: 'articleEdit', params: { id } })
}

onMounted(async () => {
  await fetchList()
})
</script>

<template>
  <ScrollContainer wrapper-class="flex flex-col gap-y-3">
    <NCard>
      <div class="flex items-center justify-between">
        <div class="text-sm text-[var(--text-color-3)]">文章列表</div>
        <NButton
          type="primary"
          @click="toCreate"
        >
          新建文章
        </NButton>
      </div>
    </NCard>
    <NCard>
      <NDataTable
        :columns="columns"
        :data="list"
        :loading="loading"
        :row-key="(row) => row.id"
      />
      <div class="flex justify-end pt-4">
        <NPagination
          :page="page"
          :page-size="pageSize"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </NCard>
  </ScrollContainer>
</template>
