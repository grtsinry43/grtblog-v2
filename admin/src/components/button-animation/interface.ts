import type { ButtonProps } from 'naive-ui'

type Animation = 'beat' | 'rotate' | 'shake'

export interface ButtonAnimationProps extends /* @vue-ignore */ ButtonProps {
  duration?: number
  animation?: Animation | boolean
}
