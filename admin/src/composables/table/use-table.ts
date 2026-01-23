import { ref, reactive, onMounted } from 'vue'


interface ApiResponse<T> {
  items: T[]
  total: number
}

type ApiFunction<T> = (params: { page: number; pageSize: number; [key: string]: any }) => Promise<ApiResponse<T>>

export function useTable<T>(api: ApiFunction<T>, initialParams: Record<string, any> = {}) {
  const loading = ref(false)
  const data = ref<T[]>([])

  // 适配下 Naive UI Pagination 的响应式对象
  const pagination = reactive({
    page: 1,
    pageSize: 10,
    itemCount: 0,
    showSizePicker: true,
    pageSizes: [10, 20, 50],
    // 把 pageChange 和 sizeChange 的逻辑内置
    onChange: (page: number) => {
      pagination.page = page
      fetchData()
    },
    onUpdatePageSize: (pageSize: number) => {
      pagination.pageSize = pageSize
      pagination.page = 1 // 改变每页大小时重置回第一页
      fetchData()
    }
  })

  async function fetchData() {
    loading.value = true
    try {
      const res = await api({
        page: pagination.page,
        pageSize: pagination.pageSize,
        ...initialParams // 允许传入额外的搜索参数
      })

      // 这里的 items 兼容你的 ArticleListResponse 结构
      data.value = res.items as any
      pagination.itemCount = res.total
    } catch (error) {
      console.error('Fetch table data failed:', error)
    } finally {
      loading.value = false
    }
  }

  // 默认挂载时请求一次
  onMounted(() => {
    fetchData()
  })

  // 返回给组件用的东西
  return {
    loading,
    data,
    pagination, // 直接绑定给 NPagination
    refresh: fetchData // 暴露刷新方法（比如删除后调用）
  }
}
