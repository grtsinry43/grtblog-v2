import type { ButtonAnimationProps } from './interface'
import type { InjectionKey } from 'vue'

export const buttonAnimationInjectionKey: InjectionKey<Partial<ButtonAnimationProps>> = Symbol(
  'buttonAnimationInjectionKey',
)
