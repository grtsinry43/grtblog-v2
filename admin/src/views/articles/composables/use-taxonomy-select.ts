import { ref, onMounted, reactive, type Ref, computed } from 'vue'

import { listCategories, listTags, createTag, createCategory } from '@/services/taxonomy'

import type { ArticleTag } from '@/services/articles'
import type { SelectOption, useMessage } from 'naive-ui'

export function useTaxonomySelect(
  formTagIds: Ref<number[]>,
  formCategoryId: Ref<number | null>,
  message: ReturnType<typeof useMessage>,
) {
  const categoryOptions = ref<SelectOption[]>([])
  const tagOptions = ref<SelectOption[]>([])
  const dynamicTags = ref<string[]>([]) // 仅用于 UI 显示的标签名列表
  const tagSearchValue = ref('')

  // 1. 初始化选项数据
  onMounted(async () => {
    try {
      const [cats, tgs] = await Promise.all([listCategories(), listTags()])
      categoryOptions.value = cats.map((c) => ({ label: c.name, value: c.id }))
      tagOptions.value = tgs.map((t) => ({ label: t.name, value: t.id }))
    } catch (e) {
      console.error('Fetch taxonomy failed', e)
    }
  })

  // 2. 核心：设置初始标签（由 fetchArticle 的结果调用）
  function setInitialTags(tags: ArticleTag[]) {
    // 既然接口返回了 {id, name}，直接用 name 初始化 UI
    dynamicTags.value = tags.map((t) => t.name)
    // 确保 ID 也同步（虽然通常 fetch 已经设置了）
    formTagIds.value = tags.map((t) => t.id)
  }

  // 3. 处理标签变更（UI -> Logic）
  // 当用户在 NDynamicTags 里增加或删除标签时触发
  async function handleTagsChange(newTags: string[]) {
    const ids: number[] = []
    const nextDynamicTags: string[] = []

    for (const tagStr of newTags) {
      const trimmed = tagStr.trim()
      if (!trimmed) continue

      // 在现有选项中查找
      const existing = tagOptions.value.find((t) => t.label === trimmed)

      if (existing) {
        // 已存在：直接使用 ID
        ids.push(existing.value as number)
        nextDynamicTags.push(trimmed)
      } else {
        // 不存在：自动创建
        try {
          const created = await createTag(trimmed)
          // 将新标签加入选项池，避免下次重复创建
          tagOptions.value.push({ label: created.name, value: created.id })
          ids.push(created.id)
          nextDynamicTags.push(created.name)
        } catch (e) {
          message.error(`创建标签 "${trimmed}" 失败`)
        }
      }
    }

    // 更新 UI 和 表单 ID
    dynamicTags.value = nextDynamicTags
    formTagIds.value = ids
  }

  // 4. 处理从搜索框添加标签
  function addTagFromSearch(value: string) {
    if (!value) return
    if (!dynamicTags.value.includes(value)) {
      // 触发 handleTagsChange 来处理 ID 查找或创建
      handleTagsChange([...dynamicTags.value, value])
    }
    tagSearchValue.value = ''
  }

  // 5. 新建分类逻辑
  const newCatModal = reactive({ show: false, name: '', slug: '', loading: false })

  async function createNewCategory() {
    if (!newCatModal.name || !newCatModal.slug) return message.error('请填写完整')

    newCatModal.loading = true
    try {
      const res = await createCategory({ name: newCatModal.name, shortUrl: newCatModal.slug })
      categoryOptions.value.push({ label: res.name, value: res.id })
      formCategoryId.value = res.id // 自动选中
      message.success('分类创建成功')
      newCatModal.show = false
      newCatModal.name = ''
      newCatModal.slug = ''
    } catch (e) {
      message.error('分类创建失败')
    } finally {
      newCatModal.loading = false
    }
  }

  // 转换 TagOptions 为 AutoComplete 需要的格式
  const autoCompleteOptions = computed(() => {
    return tagOptions.value.map((t) => ({ label: t.label as string, value: t.label as string }))
  })

  return {
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
  }
}
