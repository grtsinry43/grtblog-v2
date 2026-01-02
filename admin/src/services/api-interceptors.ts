import { useDiscreteApi } from '@/composables'
import { toRefsUserStore, useUserStore, pinia } from '@/stores'

import { addErrorInterceptor, addResponseInterceptor, setAuthTokenProvider } from './http'

export function setupApiInterceptors() {
  const { message } = useDiscreteApi()
  const userStore = useUserStore(pinia)
  const { token } = toRefsUserStore()

  setAuthTokenProvider(() => token.value)

  addResponseInterceptor(({ envelope }) => {
    if (envelope.msg && envelope.msg !== 'success' && envelope.msg !== 'ok') {
      message.success(envelope.msg)
    }
  })

  addErrorInterceptor((error) => {
    message.error(error.message || '网络异常，请稍后重试')
  })
}
