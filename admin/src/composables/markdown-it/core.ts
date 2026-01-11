import MarkdownIt from 'markdown-it'

import { componentBlockExtension } from './extensions/comp-block.ts'
import { highlightPlugin } from './plugins/highlight'

import type { MarkdownConfig } from './types'

export function createMarkdownIt(config: MarkdownConfig = {}) {
  const md = new MarkdownIt({
    html: true,
    linkify: true,
    typographer: true,
    ...config.options,
  })

  // 1. 加载 Plugins
  if (config.highlight !== false) {
    highlightPlugin(md)
  }

  // 可以在这里引入其他社区插件，如 markdown-it-anchor 等
  // md.use(anchorPlugin)

  // 2. 加载自定义业务扩展 Extensions
  componentBlockExtension(md)

  return md
}
