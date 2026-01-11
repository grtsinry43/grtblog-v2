import hljs from 'highlight.js'

import type { MarkdownExtension } from '../types'

export const highlightPlugin: MarkdownExtension = (md) => {
  md.set({
    highlight: (str, lang) => {
      if (lang && hljs.getLanguage(lang)) {
        try {
          return `<pre class="hljs"><code>${
            hljs.highlight(str, { language: lang, ignoreIllegals: true }).value
          }</code></pre>`
        } catch (__) {}
      }
      return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>'
    },
  })
}
