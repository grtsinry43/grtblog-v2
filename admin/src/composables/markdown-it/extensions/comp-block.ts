import container from 'markdown-it-container'

import type { MarkdownExtension } from '../types'

/**
 * 自定义语法扩展
 * 语法: ::: component <name> :::
 * 作用: 不渲染内容，只输出一个带有 data-component="name" 的占位符 div
 */
const INLINE_COMPONENTS = ['gallery', 'callout', 'timeline'] as const

export const componentBlockExtension: MarkdownExtension = (md) => {
  const registerContainer = (
    containerName: string,
    resolveComponentName: (tokenInfo: string) => string,
    validate?: (params: string) => boolean,
  ) => {
    md.use(container as any, containerName, {
      validate,
      render: (tokens: any[], idx: number) => {
        if (tokens[idx].nesting === 1) {
          const componentName = resolveComponentName(tokens[idx].info)
          return `<div class="md-component-placeholder" data-component="${md.utils.escapeHtml(componentName)}">`
        }
        return '</div>\n'
      },
    })
  }

  // ::: component <name>
  registerContainer(
    'component',
    (info) => {
      const match = info.trim().match(/^component\s+(.*)$/)
      return (match?.[1] || 'unknown').trim()
    },
    (params) => {
      return /^component\s+/.test(params.trim())
    },
  )

  // ::: gallery / ::: callout / ::: timeline
  INLINE_COMPONENTS.forEach((name) => {
    registerContainer(name, () => name)
  })
}
