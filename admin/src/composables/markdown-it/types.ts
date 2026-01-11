import MarkdownIt from 'markdown-it'

import type {Options} from 'markdown-it'

export type MarkdownItInstance = MarkdownIt

// 插件/扩展函数的签名：接收 md 实例和可选的 options
export type MarkdownExtension = (md: MarkdownIt, options?: any) => void

export interface MarkdownConfig {
  /** 是否启用代码高亮 */
  highlight?: boolean
  /** 其他 markdown-it 原生配置 */
  options?: Options
}
