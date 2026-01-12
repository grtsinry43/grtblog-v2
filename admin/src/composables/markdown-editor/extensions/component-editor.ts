import { EditorView } from '@codemirror/view'

import { getMarkdownComponent, parseComponentInfo } from '@/composables/markdown/shared/components'

export interface ComponentEditPayload {
  name: string
  attrs: Record<string, string>
  rawAttrs: string
  lineFrom: number
  lineTo: number
  lineText: string
  isComponentSyntax: boolean
}

interface ComponentEditorOptions {
  onEdit?: (payload: ComponentEditPayload) => void
}

export const createComponentEditorExtension = (options: ComponentEditorOptions) => {
  if (!options.onEdit) return []

  return EditorView.domEventHandlers({
    mousedown: (event, view) => {
      const pos = view.posAtCoords({ x: event.clientX, y: event.clientY })
      if (pos == null) return false

      const line = view.state.doc.lineAt(pos)
      const trimmed = line.text.trim()
      if (!trimmed.startsWith(':::')) return false

      const info = trimmed.slice(3).trim()
      const parsed = parseComponentInfo(info)
      const component = getMarkdownComponent(parsed.name)
      if (!component) return false

      const isComponentSyntax = /^:::\s*component\s+/.test(trimmed)
      options.onEdit?.({
        name: parsed.name,
        attrs: parsed.attrs,
        rawAttrs: parsed.rawAttrs,
        lineFrom: line.from,
        lineTo: line.to,
        lineText: line.text,
        isComponentSyntax,
      })

      return false
    },
  })
}
