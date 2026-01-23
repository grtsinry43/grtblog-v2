import { reactive, computed, type Ref } from 'vue'

import { getMarkdownComponent } from '@/composables/markdown/shared/components'

import type { EditorView } from 'codemirror'

// 与 extension 定义保持一致
export interface ComponentEditPayload {
  name: string
  lineFrom: number
  lineTo: number
  isComponentSyntax: boolean
  attrs: Record<string, string>
}

export function useComponentInserter(view: Ref<EditorView | undefined>) {
  const state = reactive({
    show: false,
    name: '',
    lineFrom: 0,
    lineTo: 0,
    isComponentSyntax: false,
    attrs: {} as Record<string, string>,
  })

  const touchedKeys = new Set<string>()
  const componentMeta = computed(() => getMarkdownComponent(state.name))

  const open = (payload: ComponentEditPayload) => {
    state.name = payload.name
    state.lineFrom = payload.lineFrom
    state.lineTo = payload.lineTo
    state.isComponentSyntax = payload.isComponentSyntax
    state.attrs = { ...payload.attrs }
    touchedKeys.clear()
    state.show = true
  }

  const formatLine = () => {
    const meta = componentMeta.value
    const definedKeys = meta?.attrs?.map((attr) => attr.key) ?? []
    const extraKeys = Object.keys(state.attrs).filter((key) => !definedKeys.includes(key))

    const keys = [...definedKeys, ...extraKeys].filter(
      (key) => key in state.attrs || touchedKeys.has(key),
    )

    const attrsString = keys.map((key) => `${key}="${state.attrs[key] ?? ''}"`).join(' ')

    const prefix = state.isComponentSyntax ? `::: component ${state.name}` : `::: ${state.name}`

    return attrsString ? `${prefix}{${attrsString}}` : prefix
  }

  const apply = () => {
    if (!view.value || !state.show) return
    const newLine = formatLine()
    view.value.dispatch({
      changes: { from: state.lineFrom, to: state.lineTo, insert: newLine },
    })
    // 更新选中范围，防止连续编辑位置错乱
    state.lineTo = state.lineFrom + newLine.length
  }

  const updateAttr = (key: string, value: string | boolean) => {
    touchedKeys.add(key)
    state.attrs[key] = String(value)
    apply()
  }

  return {
    state,
    componentMeta,
    open,
    updateAttr,
    close: () => (state.show = false),
  }
}
