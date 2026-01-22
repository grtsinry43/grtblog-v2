import { onBeforeUnmount } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { useDialog } from 'naive-ui'

import { useComponentModifier } from '@/composables'

import type { Ref } from 'vue'

interface LeaveConfirmOptions {
  when: Ref<boolean> | (() => boolean)
  title?: string
  content?: string
  positiveText?: string
  negativeText?: string
}

function resolveWhen(when: LeaveConfirmOptions['when']) {
  return typeof when === 'function' ? when() : when.value
}

export function useLeaveConfirm(options: LeaveConfirmOptions) {
  const dialog = useDialog()
  const { getModalModifier } = useComponentModifier()

  const title = options.title ?? '未保存的更改'
  const content = options.content ?? '你有未保存的更改，确定要离开吗？'
  const positiveText = options.positiveText ?? '离开'
  const negativeText = options.negativeText ?? '取消'

  const handleBeforeUnload = (event: BeforeUnloadEvent) => {
    if (!resolveWhen(options.when)) return
    event.preventDefault()
    event.returnValue = ''
  }

  window.addEventListener('beforeunload', handleBeforeUnload)

  onBeforeRouteLeave(() => {
    if (!resolveWhen(options.when)) return true
    return new Promise<boolean>((resolve) => {
      dialog.warning({
        ...getModalModifier(),
        title,
        content,
        positiveText,
        negativeText,
        onPositiveClick: () => resolve(true),
        onNegativeClick: () => resolve(false),
        onClose: () => resolve(false),
      })
    })
  })

  onBeforeUnmount(() => {
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })
}
