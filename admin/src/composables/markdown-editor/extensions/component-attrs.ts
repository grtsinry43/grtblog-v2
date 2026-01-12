import { type CompletionContext, type CompletionResult } from '@codemirror/autocomplete'

import {
  getMarkdownComponent,
  parseComponentAttributes,
  parseComponentInfo,
} from '@/composables/markdown/shared/components'

const COMPONENT_HEADER_RE = /^\s*:::\s*(.+)$/

export const componentAttributeSource = (context: CompletionContext): CompletionResult | null => {
  const { state, pos } = context
  const line = state.doc.lineAt(pos)
  const offset = pos - line.from
  const textBefore = line.text.slice(0, offset)

  const headerMatch = COMPONENT_HEADER_RE.exec(textBefore)
  if (!headerMatch) return null

  const openIndex = textBefore.indexOf('{')
  if (openIndex === -1) return null

  const closeIndex = textBefore.indexOf('}', openIndex)
  if (closeIndex !== -1) return null

  const { name } = parseComponentInfo(headerMatch[1])
  const component = getMarkdownComponent(name)
  if (!component || component.attrs.length === 0) return null

  const attrTextBefore = textBefore.slice(openIndex + 1)
  const tokenStart = Math.max(attrTextBefore.lastIndexOf(' '), attrTextBefore.lastIndexOf('{')) + 1
  const currentToken = attrTextBefore.slice(tokenStart)
  if (currentToken.includes('=')) return null

  const existingKeys = new Set(Object.keys(parseComponentAttributes(attrTextBefore)))
  const query = currentToken
  const insertPrefix = tokenStart > 0 && !/\s$/.test(attrTextBefore.slice(0, tokenStart)) ? ' ' : ''

  const queryLower = query.toLowerCase()
  const options = component.attrs
    .filter((attr) => !existingKeys.has(attr.key))
    .filter((attr) => (queryLower ? attr.key.toLowerCase().includes(queryLower) : true))
    .map((attr) => ({
      label: attr.key,
      type: 'property',
      detail: attr.label,
      apply: `${insertPrefix}${attr.key}=""`,
    }))

  if (!options.length) return null

  return {
    from: pos - query.length,
    to: pos,
    options,
    filter: false,
  }
}
