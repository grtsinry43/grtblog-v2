<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSelect,
  NSpace,
  NTag,
  NTreeSelect,
  useMessage,
} from 'naive-ui'

import { ScrollContainer } from '@/components'
import { navMenuIconOptions } from '@/constants/nav-menu-icons'
import {
  createNavMenu,
  deleteNavMenu,
  listNavMenus,
  reorderNavMenus,
  updateNavMenu,
  type NavMenuItem,
  type NavMenuOrderItem,
} from '@/services/navigation'

import NavMenuTree from './components/NavMenuTree.vue'

defineOptions({
  name: 'NavMenuManagement',
})

const message = useMessage()

const menuTree = ref<NavMenuItem[]>([])
const loading = ref(false)
const orderDirty = ref(false)
const savingOrder = ref(false)

const formRef = ref<InstanceType<typeof NForm> | null>(null)
const modalOpen = ref(false)
const formSubmitting = ref(false)
const editingItem = ref<NavMenuItem | null>(null)

const formState = reactive({
  name: '',
  url: '',
  parentId: 0,
  icon: null as string | null,
})

const formRules = {
  name: {
    required: true,
    message: '请输入菜单名称',
    trigger: ['blur', 'input'],
  },
  url: {
    required: true,
    message: '请输入菜单链接',
    trigger: ['blur', 'input'],
  },
}

const resetForm = () => {
  formState.name = ''
  formState.url = ''
  formState.parentId = 0
  formState.icon = null
  editingItem.value = null
}

const fetchMenus = async () => {
  loading.value = true
  try {
    const data = await listNavMenus()
    menuTree.value = normalizeMenus(data)
    orderDirty.value = false
  } catch (error) {
    message.error(error instanceof Error ? error.message : '获取菜单失败')
  } finally {
    loading.value = false
  }
}

const normalizeMenus = (items: NavMenuItem[] = []): NavMenuItem[] => {
  return items.map((item) => ({
    ...item,
    children: item.children ? normalizeMenus(item.children) : [],
  }))
}

const markDirty = () => {
  orderDirty.value = true
}

const collectDescendantIds = (item: NavMenuItem | null, set: Set<number>) => {
  if (!item) return
  set.add(item.id)
  item.children?.forEach((child) => collectDescendantIds(child, set))
}

const buildTreeOptions = (items: NavMenuItem[], blocked: Set<number>) => {
  return items
    .filter((item) => !blocked.has(item.id))
    .map((item) => ({
      label: item.name,
      key: item.id,
      children: item.children ? buildTreeOptions(item.children, blocked) : undefined,
    }))
}

const treeOptions = computed(() => {
  const blocked = new Set<number>()
  collectDescendantIds(editingItem.value, blocked)
  const options = buildTreeOptions(menuTree.value, blocked)
  return [{ label: '顶级菜单', key: 0, children: options }]
})

const openCreate = (parent?: NavMenuItem | null) => {
  resetForm()
  if (parent) {
    formState.parentId = parent.id
  }
  modalOpen.value = true
}

const openEdit = (item: NavMenuItem) => {
  editingItem.value = item
  formState.name = item.name
  formState.url = item.url
  formState.parentId = item.parentId ?? 0
  formState.icon = item.icon ?? null
  modalOpen.value = true
}

const handleDelete = async (item: NavMenuItem) => {
  if (!window.confirm(`确认删除菜单「${item.name}」及其子项吗？`)) return
  try {
    await deleteNavMenu(item.id)
    message.success('已删除')
    await fetchMenus()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '删除失败')
  }
}

const buildOrderPayload = (items: NavMenuItem[], parentId: number | null = null, acc: NavMenuOrderItem[] = []) => {
  items.forEach((item, index) => {
    acc.push({
      id: item.id,
      parentId,
      sort: index + 1,
    })
    if (item.children?.length) {
      buildOrderPayload(item.children, item.id, acc)
    }
  })
  return acc
}

const saveOrder = async () => {
  const payload = buildOrderPayload(menuTree.value)
  if (!payload.length) {
    message.warning('没有可保存的排序')
    return
  }
  savingOrder.value = true
  try {
    await reorderNavMenus(payload)
    message.success('排序已保存')
    orderDirty.value = false
  } catch (error) {
    message.error(error instanceof Error ? error.message : '保存排序失败')
  } finally {
    savingOrder.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  formSubmitting.value = true
  const parentId = formState.parentId === 0 ? null : formState.parentId
  const payload = {
    name: formState.name.trim(),
    url: formState.url.trim(),
    parentId,
    icon: formState.icon || null,
  }

  try {
    if (editingItem.value) {
      await updateNavMenu(editingItem.value.id, payload)
      message.success('菜单已更新')
    } else {
      await createNavMenu(payload)
      message.success('菜单已创建')
    }
    modalOpen.value = false
    await fetchMenus()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    formSubmitting.value = false
  }
}

onMounted(() => {
  fetchMenus()
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="导航菜单">
      <div class="space-y-6">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="space-y-2">
            <div class="text-sm text-neutral-600 dark:text-neutral-400">
              导航菜单用于前台侧边栏与移动端导航，支持拖拽调整层级与顺序。
            </div>
            <NTag size="small" type="info">图标字段使用 lucide-svelte 的图标名称</NTag>
          </div>
          <NSpace>
            <NButton secondary :loading="loading" @click="fetchMenus">刷新</NButton>
            <NButton type="primary" @click="openCreate()">新增菜单</NButton>
            <NButton
              type="success"
              :disabled="!orderDirty"
              :loading="savingOrder"
              @click="saveOrder"
            >
              保存排序
            </NButton>
          </NSpace>
        </div>

        <NavMenuTree
          v-model:items="menuTree"
          @edit="openEdit"
          @delete="handleDelete"
          @add-child="openCreate"
          @drag="markDirty"
        />

        <div v-if="orderDirty" class="text-xs text-orange-600 dark:text-orange-400">
          当前排序有改动，请点击“保存排序”同步到服务端。
        </div>
      </div>
    </NCard>

    <NModal v-model:show="modalOpen" preset="card" title="菜单配置" class="w-full max-w-[440px]">
      <NForm ref="formRef" :model="formState" :rules="formRules" label-placement="left" label-width="90">
        <NFormItem label="名称" path="name">
          <NInput v-model:value="formState.name" placeholder="例如：首页" />
        </NFormItem>
        <NFormItem label="链接" path="url">
          <NInput v-model:value="formState.url" placeholder="例如：/ 或 https://..." />
        </NFormItem>
        <NFormItem label="父级菜单">
          <NTreeSelect
            v-model:value="formState.parentId"
            :options="treeOptions"
            clearable
            placeholder="顶级菜单"
          />
        </NFormItem>
        <NFormItem label="图标">
          <NSelect
            v-model:value="formState.icon"
            :options="navMenuIconOptions"
            placeholder="选择图标（可选）"
            clearable
            filterable
          />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end gap-3">
          <NButton @click="modalOpen = false">取消</NButton>
          <NButton type="primary" :loading="formSubmitting" @click="handleSubmit">
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </ScrollContainer>
</template>
