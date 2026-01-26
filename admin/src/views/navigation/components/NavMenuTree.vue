<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NButton,
  NTag,
  NCard,
  NSpace,
  NThing,
  NPopconfirm,
  NEmpty,
  NButtonGroup,
  NIcon,
} from 'naive-ui'
import { VueDraggable } from 'vue-draggable-plus'
// 引入 Ionicons5 图标
import {
  ReorderTwo,
  CreateOutline,
  Add,
  TrashBinOutline,
  FolderOpenOutline,
  DocumentTextOutline,
} from '@vicons/ionicons5'

import type { NavMenuItem } from '@/services/navigation'

defineOptions({
  name: 'NavMenuTree',
})

const props = withDefaults(defineProps<{ items?: NavMenuItem[] }>(), {
  items: () => [],
})

const emit = defineEmits<{
  (e: 'update:items', value: NavMenuItem[]): void
  (e: 'edit', item: NavMenuItem): void
  (e: 'delete', item: NavMenuItem): void
  (e: 'add-child', item: NavMenuItem): void
  (e: 'drag'): void
}>()

const dragGroup = {
  name: 'nav-menu',
  pull: true,
  put: true,
}

// ================= 核心逻辑修改开始 =================

// 1. 本地维护一个响应式数组，避免直接修改 props 导致的 UI 抖动
const localList = ref<NavMenuItem[]>([])

// 2. 监听 props 变化，同步到本地（解决外部更新数据时的同步问题）
watch(
  () => props.items,
  (newVal) => {
    // 简单的引用对比或深度对比，防止死循环
    // 这里直接覆盖，因为通常 props 变更是由父组件完全刷新引起的
    localList.value = newVal || []
  },
  { immediate: true },
)

// 3. 拖拽结束或数据变更时，通知父组件
// 注意：vue-draggable-plus 会自动修改 v-model 绑定的 localList
const handleDragChange = () => {
  emit('update:items', localList.value)
  emit('drag')
}

// 4. 处理子级递归更新
// 当子组件 emit('update:items') 时，我们需要更新当前层级对应节点的 children
const handleChildUpdate = (element: NavMenuItem, newChildren: NavMenuItem[]) => {
  element.children = newChildren
  // 子级变了，也相当于当前层级的数据变了，需要向上冒泡
  emit('update:items', localList.value)
}

// ================= 核心逻辑修改结束 =================
</script>

<template>
  <VueDraggable
    v-if="localList.length"
    v-model="localList"
    :group="dragGroup"
    handle=".drag-handle"
    item-key="id"
    class="flex flex-col gap-3"
    @end="handleDragChange"
  >
    <!-- 使用 v-for 显式渲染 -->
    <div
      v-for="element in localList"
      :key="element.id"
    >
      <NCard
        size="small"
        hoverable
        content-style="padding: 10px 16px;"
      >
        <NThing>
          <!-- 拖拽手柄 -->
          <template #avatar>
            <div class="flex h-full items-center">
              <NIcon
                size="20"
                class="drag-handle cursor-move text-gray-400 transition-colors hover:text-gray-600 active:text-primary"
                :component="ReorderTwo"
              />
            </div>
          </template>

          <!-- 标题栏 -->
          <template #header>
            <NSpace
              align="center"
              :size="8"
            >
              <span class="text-sm font-medium">{{ element.name }}</span>
              <NTag
                v-if="element.icon"
                :bordered="false"
                type="info"
                size="small"
                round
              >
                {{ element.icon }}
              </NTag>
            </NSpace>
          </template>

          <!-- 描述/URL -->
          <template #description>
            <div class="truncate text-xs text-gray-400">
              {{ element.url || '无链接' }}
            </div>
          </template>

          <!-- 右侧操作栏 -->
          <template #header-extra>
            <NSpace
              align="center"
              :size="8"
            >
              <NButtonGroup size="tiny">
                <NButton
                  secondary
                  strong
                  @click="emit('add-child', element)"
                >
                  <template #icon><NIcon :component="Add" /></template>
                  子项
                </NButton>
                <NButton
                  secondary
                  strong
                  @click="emit('edit', element)"
                >
                  <template #icon><NIcon :component="CreateOutline" /></template>
                  编辑
                </NButton>
              </NButtonGroup>

              <NPopconfirm
                @positive-click="emit('delete', element)"
                negative-text="手滑了"
                positive-text="确定删除"
              >
                <template #trigger>
                  <NButton
                    size="tiny"
                    secondary
                    type="error"
                  >
                    <template #icon><NIcon :component="TrashBinOutline" /></template>
                  </NButton>
                </template>
                确定要删除“{{ element.name }}”吗？<br />
                <span class="text-xs text-gray-400">如果是目录，子菜单也会被删除。</span>
              </NPopconfirm>
            </NSpace>
          </template>
        </NThing>

        <!-- 递归子菜单 -->
        <!-- 注意：这里不直接使用 v-model:items，而是拆开写，为了通过 handleChildUpdate 拦截更新 -->
        <div
          v-if="element.children?.length"
          class="mt-3 pl-10"
        >
          <div class="border-l-2 border-gray-100 pl-3 dark:border-neutral-800">
            <NavMenuTree
              :items="element.children"
              @update:items="(val) => handleChildUpdate(element, val)"
              @edit="emit('edit', $event)"
              @delete="emit('delete', $event)"
              @add-child="emit('add-child', $event)"
              @drag="emit('drag')"
            />
          </div>
        </div>
      </NCard>
    </div>
  </VueDraggable>

  <!-- 空状态 -->
  <NEmpty
    v-else
    description="暂无菜单"
    class="rounded-lg border border-dashed border-gray-200 py-8 dark:border-neutral-800"
  >
    <template #extra>
      <div class="text-xs text-gray-400">还没有菜单呀，点击添加吧～</div>
    </template>
  </NEmpty>
</template>

<style scoped>
/* 微调 NThing 头部边距，使其更紧凑 */
:deep(.n-thing-header) {
  margin-bottom: 2px !important;
}
</style>
