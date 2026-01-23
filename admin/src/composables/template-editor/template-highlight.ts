import { RangeSetBuilder } from '@codemirror/state'
import {
  Decoration,
  EditorView,
  ViewPlugin,
  hoverTooltip,
  type DecorationSet,
  type ViewUpdate,
} from '@codemirror/view'

const templateToken = Decoration.mark({ class: 'cm-template-token' })
const templateInvalid = Decoration.mark({
  class: 'cm-template-invalid',
  attributes: { title: '未知模板变量' },
})

function normalizeExpr(expr: string) {
  return expr.replace(/\s+/g, ' ').trim()
}

function isValidTemplateExpression(expr: string) {
  const normalized = normalizeExpr(expr)
  if (normalized === '.Name' || normalized === '.OccurredAt') {
    return true
  }
  if (/^\.Event(\.[A-Za-z_][A-Za-z0-9_]*)*$/.test(normalized)) {
    return true
  }
  return /^toJSON\s+\.Event(\.[A-Za-z_][A-Za-z0-9_]*)*$/.test(normalized)
}

function buildDecorations(view: EditorView) {
  const builder = new RangeSetBuilder<Decoration>()
  const pattern = /{{[^}]*}}/g

  for (const range of view.visibleRanges) {
    const text = view.state.doc.sliceString(range.from, range.to)
    pattern.lastIndex = 0
    for (let match = pattern.exec(text); match; match = pattern.exec(text)) {
      const start = range.from + match.index
      const end = start + match[0].length
      const expr = match[0].slice(2, -2)
      const decoration = isValidTemplateExpression(expr) ? templateToken : templateInvalid
      builder.add(start, end, decoration)
    }
  }

  return builder.finish()
}

type TemplateMatch = {
  from: number
  to: number
  expr: string
}

function findTemplateMatch(lineText: string, lineFrom: number, pos: number): TemplateMatch | null {
  const pattern = /{{[^}]*}}/g
  let match: RegExpExecArray | null
  while ((match = pattern.exec(lineText))) {
    const from = lineFrom + match.index
    const to = from + match[0].length
    if (pos >= from && pos <= to) {
      return {
        from,
        to,
        expr: match[0].slice(2, -2),
      }
    }
  }
  return null
}

export const templateHighlightExtension = ViewPlugin.fromClass(
  class {
    decorations: DecorationSet

    constructor(view: EditorView) {
      this.decorations = buildDecorations(view)
    }

    update(update: ViewUpdate) {
      if (update.docChanged || update.viewportChanged) {
        this.decorations = buildDecorations(update.view)
      }
    }
  },
  {
    decorations: (instance) => instance.decorations,
  },
)

export const templateTooltipExtension = hoverTooltip((view, pos, side) => {
  const line = view.state.doc.lineAt(pos)
  const match = findTemplateMatch(line.text, line.from, pos)
  if (!match) return null
  if (match.from === pos && side < 0) return null
  if (match.to === pos && side > 0) return null
  if (isValidTemplateExpression(match.expr)) return null

  return {
    pos: match.from,
    end: match.to,
    above: true,
    create() {
      const dom = document.createElement('div')
      dom.className = 'cm-template-tooltip'
      dom.textContent = '未知模板变量'
      return { dom }
    },
  }
})
