import { useComponentModifier, useDiscreteApi } from '@/composables'

import type { MessageType, ModalOptions, NotificationType } from 'naive-ui'

const { message, modal, notification } = useDiscreteApi()

const { getModalModifier } = useComponentModifier()

export const useMessageDiscrete = (type: MessageType = 'info') => {
  message.create(`${type} - 与Setup里调用单独分开显示`, {
    type,
    duration: 5000,
    closable: true,
  })
}

export const useModalDiscrete = (type: ModalOptions['type'] = 'info') => {
  modal.create({
    ...getModalModifier(),
    title: `${type}`,
    content: `Setup 外的 ${type} Dialog`,
    preset: 'dialog',
    type,
    positiveText: '确认',
    negativeText: '算了',
  })
}

export const useNotificationDiscrete = (type: NotificationType = 'info') => {
  notification.create({
    type,
    content: '又要说点啥呢',
    meta: '我是 Setup 外的 Notification， 出现位置当然可以不一样',
    duration: 5000,
  })
}
