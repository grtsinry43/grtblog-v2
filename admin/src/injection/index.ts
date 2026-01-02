import type { MediaQueryProvider, LayoutProvider } from './interface'
import type { InjectionKey } from 'vue'

export const mediaQueryInjectionKey: InjectionKey<MediaQueryProvider> =
  Symbol('mediaQueryInjectionKey')

export const layoutInjectionKey: InjectionKey<LayoutProvider> = Symbol('layoutInjectionKey')

export * from './interface'
