<script setup lang="ts">
import {
  CloudArrowUp24Regular,
  Copy24Regular,
  Delete24Regular,
  Document24Regular,
  Edit24Regular,
  Image24Regular,
} from '@vicons/fluent'
import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NIcon,
  NImage,
  NInput,
  NModal,
  NPagination,
  NRadioButton,
  NRadioGroup,
  NSpace,
  NTag,
  NTree,
  NUpload,
  useMessage,
} from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'

import { deleteFile, downloadFile, listUploads, renameFile, uploadFile } from '@/services/uploads'

import type { FileType, UploadFileResponse } from '@/services/uploads'
import type { DataTableColumns, UploadFileInfo } from 'naive-ui'

const message = useMessage()

// State
const files = ref<UploadFileResponse[]>([])
const loading = ref(false)
const uploading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const uploadType = ref<FileType>('picture')
const treeSelectedKeys = ref<string[]>(['all'])

// Rename modal
const renameModalVisible = ref(false)
const renamingFile = ref<UploadFileResponse | null>(null)
const newFileName = ref('')

// Delete confirmation
const deleteModalVisible = ref(false)
const deletingFile = ref<UploadFileResponse | null>(null)

// Image preview
const previewVisible = ref(false)
const previewImageUrl = ref('')

// Computed
const isEmpty = computed(() => files.value.length === 0 && !loading.value)
const filteredFiles = computed(() => {
  const selected = treeSelectedKeys.value[0]
  if (selected === 'picture' || selected === 'file') {
    return files.value.filter((item) => item.type === selected)
  }
  return files.value
})
const treeData = computed(() => {
  const pictureCount = files.value.filter((item) => item.type === 'picture').length
  const fileCount = files.value.filter((item) => item.type === 'file').length
  return [
    {
      key: 'all',
      label: `全部 (${files.value.length})`,
      children: [
        {
          key: 'picture',
          label: `图片 (${pictureCount})`,
        },
        {
          key: 'file',
          label: `文件 (${fileCount})`,
        },
      ],
    },
  ]
})

// Format file size
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

// Format date
function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Table columns
const columns = computed<DataTableColumns<UploadFileResponse>>(() => [
  {
    title: '预览',
    key: 'preview',
    width: 80,
    render: (row) => {
      if (row.type === 'picture') {
        return h(
          'div',
          {
            style: {
              cursor: 'pointer',
              width: '50px',
              height: '50px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            },
            onClick: () => {
              previewImageUrl.value = row.publicUrl
              previewVisible.value = true
            },
          },
          h(NImage, {
            src: row.publicUrl,
            width: 50,
            height: 50,
            objectFit: 'cover',
            style: 'border-radius: 4px',
            previewDisabled: true,
          }),
        )
      }
      return h(
        NIcon,
        { size: 32, color: '#18a058' },
        {
          default: () => h(Document24Regular),
        },
      )
    },
  },
  {
    title: '文件名',
    key: 'name',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) => {
      return h(
        NTag,
        {
          type: row.type === 'picture' ? 'success' : 'info',
          size: 'small',
        },
        {
          default: () => (row.type === 'picture' ? '图片' : '文件'),
        },
      )
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 120,
    render: (row) => formatFileSize(row.size),
  },
  {
    title: '上传时间',
    key: 'createdAt',
    width: 180,
    render: (row) => formatDate(row.createdAt),
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render: (row) => {
      return h(
        NSpace,
        { size: 'small' },
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                onClick: () => handleCopyUrl(row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(Copy24Regular) }),
                default: () => '复制链接',
              },
            ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                onClick: () => openRenameModal(row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(Edit24Regular) }),
                default: () => '重命名',
              },
            ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                onClick: () => handleDownload(row),
              },
              {
                default: () => '下载',
              },
            ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                type: 'error',
                onClick: () => openDeleteModal(row),
              },
              {
                icon: () => h(NIcon, null, { default: () => h(Delete24Regular) }),
                default: () => '删除',
              },
            ),
          ],
        },
      )
    },
  },
])

// Fetch files list
async function fetchFiles() {
  loading.value = true
  try {
    const response = await listUploads({
      page: page.value,
      pageSize: pageSize.value,
    })
    files.value = response.items
    total.value = response.total
  } catch (error) {
    message.error('加载文件列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Handle file upload
async function handleUpload({ file }: { file: UploadFileInfo }) {
  if (!file.file) return

  uploading.value = true
  try {
    const response = await uploadFile(file.file, uploadType.value)
    message.success(response.duplicated ? '文件已存在，已复用' : '上传成功')
    await fetchFiles()
  } catch (error) {
    message.error('上传失败')
    console.error(error)
  } finally {
    uploading.value = false
  }
}

// Copy URL to clipboard
async function handleCopyUrl(file: UploadFileResponse) {
  try {
    await navigator.clipboard.writeText(file.publicUrl)
    message.success('链接已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
    console.error(error)
  }
}

// Open rename modal
function openRenameModal(file: UploadFileResponse) {
  renamingFile.value = file
  newFileName.value = file.name
  renameModalVisible.value = true
}

// Handle rename
async function handleRename() {
  if (!renamingFile.value || !newFileName.value.trim()) {
    message.warning('请输入文件名')
    return
  }

  try {
    await renameFile(renamingFile.value.id, { name: newFileName.value.trim() })
    message.success('重命名成功')
    renameModalVisible.value = false
    await fetchFiles()
  } catch (error) {
    message.error('重命名失败')
    console.error(error)
  }
}

// Open delete modal
function openDeleteModal(file: UploadFileResponse) {
  deletingFile.value = file
  deleteModalVisible.value = true
}

// Handle delete
async function handleDelete() {
  if (!deletingFile.value) return

  try {
    await deleteFile(deletingFile.value.id)
    message.success('删除成功')
    deleteModalVisible.value = false
    
    // If current page becomes empty after deletion, go to previous page
    if (files.value.length === 1 && page.value > 1) {
      page.value--
    }
    
    await fetchFiles()
  } catch (error) {
    message.error('删除失败')
    console.error(error)
  }
}

// Handle download
async function handleDownload(file: UploadFileResponse) {
  try {
    await downloadFile(file.id, file.name)
    message.success('下载开始')
  } catch (error) {
    message.error('下载失败')
    console.error(error)
  }
}

// Handle page change
function handlePageChange(newPage: number) {
  page.value = newPage
  fetchFiles()
}

// Handle page size change
function handlePageSizeChange(newPageSize: number) {
  pageSize.value = newPageSize
  page.value = 1
  fetchFiles()
}

// Initialize
onMounted(() => {
  fetchFiles()
})

function handleTreeSelect(keys: Array<string | number>) {
  const selectedKey = String(keys[0] ?? 'all')
  treeSelectedKeys.value = [selectedKey]
}
</script>

<template>
  <div class="uploads-container">
    <NCard title="文件管理" :bordered="false">
      <template #header-extra>
        <NSpace align="center">
          <NRadioGroup v-model:value="uploadType" size="small">
            <NRadioButton value="picture">
              <NIcon>
                <Image24Regular />
              </NIcon>
              图片
            </NRadioButton>
            <NRadioButton value="file">
              <NIcon>
                <Document24Regular />
              </NIcon>
              文件
            </NRadioButton>
          </NRadioGroup>
          <NUpload
            :show-file-list="false"
            :custom-request="handleUpload"
            :disabled="uploading"
          >
            <NButton type="primary" :loading="uploading">
              <template #icon>
                <NIcon>
                  <CloudArrowUp24Regular />
                </NIcon>
              </template>
              上传文件
            </NButton>
          </NUpload>
        </NSpace>
      </template>

      <div class="upload-area">
        <NUpload
          :show-file-list="false"
          :custom-request="handleUpload"
          :disabled="uploading"
          directory-dnd
        >
          <div class="upload-dragger">
            <div class="upload-icon">
              <NIcon size="48" :depth="3">
                <CloudArrowUp24Regular />
              </NIcon>
            </div>
            <div class="upload-text">
              <p class="upload-hint">点击或拖拽文件到此区域上传</p>
              <p class="upload-description">
                当前类型：{{ uploadType === 'picture' ? '图片' : '文件' }}
              </p>
            </div>
          </div>
        </NUpload>
      </div>

      <div v-if="isEmpty" class="empty-container">
        <NEmpty description="暂无文件" />
      </div>

      <div v-else class="content-layout">
        <div class="tree-panel">
          <div class="tree-title">文件树</div>
          <NTree
            :data="treeData"
            :selected-keys="treeSelectedKeys"
            block-line
            default-expand-all
            @update:selected-keys="handleTreeSelect"
          />
        </div>

        <div class="table-panel">
          <NDataTable
            :columns="columns"
            :data="filteredFiles"
            :loading="loading"
            :bordered="false"
            :single-line="false"
          />

          <div class="pagination-container">
            <NPagination
              v-model:page="page"
              v-model:page-size="pageSize"
              :page-count="Math.ceil(total / pageSize)"
              :page-sizes="[10, 20, 50, 100]"
              show-size-picker
              @update:page="handlePageChange"
              @update:page-size="handlePageSizeChange"
            />
          </div>
        </div>
      </div>
    </NCard>

    <!-- Rename Modal -->
    <NModal
      v-model:show="renameModalVisible"
      preset="dialog"
      title="重命名文件"
      positive-text="确认"
      negative-text="取消"
      @positive-click="handleRename"
    >
      <div style="padding: 16px 0">
        <NInput
          v-model:value="newFileName"
          placeholder="请输入新文件名"
          @keyup.enter="handleRename"
        />
      </div>
    </NModal>

    <!-- Delete Confirmation Modal -->
    <NModal
      v-model:show="deleteModalVisible"
      preset="dialog"
      title="确认删除"
      type="warning"
      positive-text="删除"
      negative-text="取消"
      @positive-click="handleDelete"
    >
      <p>确定要删除文件 "{{ deletingFile?.name }}" 吗？</p>
      <p style="color: #f5222d; margin-top: 8px">此操作将永久删除文件，无法恢复。</p>
    </NModal>

    <!-- Image Preview Modal -->
    <NModal v-model:show="previewVisible" preset="card" style="max-width: 800px">
      <template #header>
        <span>图片预览</span>
      </template>
      <div class="preview-container">
        <NImage :src="previewImageUrl" />
      </div>
    </NModal>
  </div>
</template>

<style scoped>
.uploads-container {
  padding: 16px;
  height: 100%;
}

.upload-area {
  margin-bottom: 24px;
}

.upload-dragger {
  padding: 40px 20px;
  background: transparent;
  border: 2px dashed var(--n-border-color);
  border-radius: 8px;
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.upload-dragger:hover {
  border-color: var(--n-border-color);
  background: rgba(0, 0, 0, 0.03);
}

.upload-icon {
  margin-bottom: 12px;
  color: var(--n-text-color-disabled);
}

.upload-hint {
  font-size: 16px;
  color: var(--n-text-color);
  margin: 0 0 8px;
}

.upload-description {
  font-size: 14px;
  color: var(--n-text-color-disabled);
  margin: 0;
}

.empty-container {
  padding: 60px 0;
}

.content-layout {
  display: flex;
  gap: 16px;
}

.tree-panel {
  width: 220px;
  min-width: 200px;
  padding: 12px;
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
  background: transparent;
  height: fit-content;
}

.tree-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 8px;
}

.table-panel {
  flex: 1;
  min-width: 0;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 16px 0;
}

.preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
}

@media (max-width: 900px) {
  .content-layout {
    flex-direction: column;
  }

  .tree-panel {
    width: 100%;
  }
}
</style>
