import { useMessage } from 'naive-ui'
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useLeaveConfirm } from '@/composables'
import { createArticle, getArticle, updateArticle, type ArticleDetail } from '@/services/articles'

export function useArticleForm() {
  const route = useRoute()
  const router = useRouter()
  const message = useMessage()

  // ID 解析：处理 'new' 或具体的数字 ID
  const articleId = computed(() => {
    const param = route.params.id
    if (!param || param === 'new') return null
    const id = Number(param)
    return Number.isFinite(id) ? id : null
  })

  const isCreating = computed(() => articleId.value === null)
  const loading = ref(false)
  const saving = ref(false)
  const initialSnapshot = ref('')

  // 表单数据模型
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

  // 脏检查逻辑
  const takeSnapshot = () => JSON.stringify(form)
  const isDirty = computed(
    () => initialSnapshot.value !== '' && takeSnapshot() !== initialSnapshot.value,
  )

  // 获取数据
  async function fetch() {
    if (isCreating.value) {
      initialSnapshot.value = takeSnapshot()
      return null
    }

    loading.value = true
    try {
      const data = await getArticle(articleId.value!)

      // 数据回填
      form.title = data.title
      form.summary = data.summary || ''
      form.leadIn = data.leadIn || ''
      form.content = data.content
      form.cover = data.cover || ''
      form.categoryId = data.categoryId ?? null
      // 注意：这里只负责回填 ID，Tags 的名字显示交给 useTaxonomySelect 处理
      form.tagIds = data.tags?.map((t) => t.id) ?? []
      form.shortUrl = data.shortUrl
      form.isPublished = data.isPublished
      form.isTop = data.isTop
      form.isHot = data.isHot
      form.isOriginal = data.isOriginal

      initialSnapshot.value = takeSnapshot()
      return data // 返回完整数据供外部使用（如初始化标签名）
    } catch (e) {
      console.error(e)
      message.error('无法加载文章数据')
      router.replace({ name: 'articleList' })
      return null
    } finally {
      loading.value = false
    }
  }

  // 保存数据
  async function save() {
    if (!form.title.trim()) return message.error('请输入标题')
    if (!form.content.trim()) return message.error('请输入正文内容')
    if (!isCreating.value && !form.shortUrl.trim()) return message.error('短链接不能为空')

    saving.value = true
    try {
      // 构造 payload，去除空字符串
      const payload = {
        ...form,
        leadIn: form.leadIn || null,
        cover: form.cover || null,
        shortUrl: form.shortUrl || undefined, //如果是新建，可能允许为空由后端生成？根据原逻辑保持一致
      }

      if (isCreating.value) {
        await createArticle(payload)
        message.success('创建成功')
      } else {
        await updateArticle(articleId.value!, payload)
        message.success('更新成功')
      }

      // 保存成功后更新快照，避免触发离开提示
      initialSnapshot.value = takeSnapshot()

      // 保存后跳转回列表 (保持原逻辑)
      router.push({ name: 'articleList' })
    } catch (e: any) {
      message.error(e.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  // 注册离开确认
  useLeaveConfirm({
    when: isDirty,
    title: '未保存的更改',
    content: '当前内容未保存，确定要离开吗？',
    positiveText: '离开',
    negativeText: '继续编辑',
  })

  // 挂载时自动获取
  onMounted(fetch) // 注意：这里 fetch 的返回值在 onMounted 里无法被 setup 直接拿到，所以我们在 edit.vue 里还要再显式调用一次或者由 edit.vue 接管 onMounted

  return {
    form,
    loading,
    saving,
    isCreating,
    isDirty,
    fetch,
    save,
  }
}
