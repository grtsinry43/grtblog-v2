import { RangeSet } from '@codemirror/state'
import { Decoration, type DecorationSet, EditorView, ViewPlugin, ViewUpdate } from '@codemirror/view'

const lineHint = Decoration.line({ class: 'cm-slash-hint-line' })
const lineHintLast = Decoration.line({ class: 'cm-slash-hint-line cm-slash-hint-line-last' })

function buildDecorations(view: EditorView): DecorationSet {
  const { state } = view
  if (!view.hasFocus) {
    return RangeSet.empty
  }
  const cursor = state.selection.main.head
  const line = state.doc.lineAt(cursor)
  if (line.text.trim() !== '') {
    return RangeSet.empty
  }
  const isLastLine = line.number === state.doc.lines

  const decorations: { from: number; to: number; value: Decoration }[] = []

  decorations.push({
    from: line.from,
    to: line.from,
    value: isLastLine ? lineHintLast : lineHint,
  })

  return RangeSet.of(decorations, true)
}

export const slashHintPlugin = ViewPlugin.fromClass(
  class {
    decorations: DecorationSet
    constructor(view: EditorView) {
      this.decorations = buildDecorations(view)
    }
    update(update: ViewUpdate) {
      if (update.docChanged || update.selectionSet) {
        this.decorations = buildDecorations(update.view)
      }
    }
  },
  {
    decorations: (v) => v.decorations,
  },
)

export const slashHintExtension = [slashHintPlugin]
