import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { EditorView } from '@codemirror/view'
import { tags, tags as t } from '@lezer/highlight'

const editorTheme = EditorView.theme({
  '&': {
    color: 'var(--cm-fg, #1f2328)',
    backgroundColor: 'var(--cm-bg, #ffffff)',
  },
  '.cm-content': {
    caretColor: 'var(--cm-cursor, #3b82f6)',
    padding: '20px 0px 40px 4px',
  },
  '.cm-cursor, .cm-dropCursor': {
    borderLeftColor: 'var(--cm-cursor, #3b82f6)',
  },
  '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection':
    {
      backgroundColor: 'var(--cm-selection, rgba(59, 130, 246, 0.2))',
    },
  '.cm-activeLine': {
    backgroundColor: 'var(--cm-active-line, rgba(100, 116, 139, 0.08))',
  },
  '.cm-activeLineGutter': {
    backgroundColor: 'var(--cm-active-line, rgba(100, 116, 139, 0.08))',
  },
  '.cm-gutters': {
    backgroundColor: 'var(--cm-gutter-bg, transparent)',
    color: 'var(--cm-gutter-fg, #94a3b8)',
    border: 'none',
  },
  '.cm-panels': {
    backgroundColor: 'var(--cm-bg, #ffffff)',
    color: 'var(--cm-fg, #1f2328)',
  },
  '.cm-tooltip': {
    backgroundColor: 'var(--cm-tooltip-bg, #ffffff)',
    color: 'var(--cm-fg, #1f2328)',
    border: '1px solid var(--cm-tooltip-border, #e5e7eb)',
    boxShadow: 'var(--cm-tooltip-shadow, 0 6px 16px rgba(15, 23, 42, 0.12))',
  },
  '.cm-tooltip .cm-tooltip-arrow:after': {
    borderTopColor: 'var(--cm-tooltip-bg, #ffffff)',
    borderBottomColor: 'var(--cm-tooltip-bg, #ffffff)',
  },
})

const highlightStyle = HighlightStyle.define([
  { tag: [t.keyword, t.modifier], color: 'var(--cm-syntax-keyword, #7c3aed)' },
  {
    tag: [t.name, t.deleted, t.character, t.macroName],
    color: 'var(--cm-syntax-variable, #1f2328)',
  },
  { tag: [t.propertyName, t.attributeName], color: 'var(--cm-syntax-property, #0ea5e9)' },
  { tag: [t.string, t.special(t.string)], color: 'var(--cm-syntax-string, #16a34a)' },
  {
    tag: [t.number, t.changed, t.annotation, t.constant(t.name)],
    color: 'var(--cm-syntax-number, #f59e0b)',
  },
  {
    tag: [t.function(t.variableName), t.function(t.name)],
    color: 'var(--cm-syntax-function, #0ea5e9)',
  },
  { tag: [t.operator, t.operatorKeyword], color: 'var(--cm-syntax-operator, #64748b)' },
  { tag: [t.className, t.typeName], color: 'var(--cm-syntax-type, #6366f1)' },
  { tag: [t.tagName], color: 'var(--cm-syntax-tag, #ef4444)' },
  { tag: [t.invalid], color: 'var(--cm-syntax-invalid, #ef4444)' },
  { tag: [t.comment], color: 'var(--cm-syntax-comment, #94a3b8)', fontStyle: 'italic' },
  { tag: [t.punctuation, t.bracket], color: 'var(--cm-syntax-punctuation, #94a3b8)' },
  { tag: tags.strikethrough, textDecoration: 'line-through', color: 'var(--cm-fg-muted)' },
])

export const codeMirrorTheme = [editorTheme, syntaxHighlighting(highlightStyle)]
