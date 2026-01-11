import { computed, toValue, type MaybeRefOrGetter } from 'vue'

import { createMarkdownIt } from './core'

import type { MarkdownConfig } from './types'

export function useMarkdownIt(
  source: MaybeRefOrGetter<string>,
  config: MarkdownConfig = {},
  renderDeps: MaybeRefOrGetter<unknown>[] = [],
) {
  // 初始化实例，这里用每次创建吧，之后配置可能改变
  const _md = createMarkdownIt(config)

  const html = computed(() => {
    renderDeps.forEach((dep) => toValue(dep))
    const text = toValue(source)
    if (!text) return ''
    return _md.render(text)
  })

  return {
    html,
    _md, // 暴露实例以便调试或特殊操作
  }
}
