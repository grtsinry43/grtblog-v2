import { json } from '@codemirror/lang-json'
import { EditorState, Compartment, type Extension } from '@codemirror/state'
import { EditorView, basicSetup } from 'codemirror'
import { type Ref, shallowRef, onMounted, onUnmounted } from 'vue'

import { codeMirrorTheme } from '@/composables/markdown-editor/codemirror-theme'
import '@/composables/markdown-editor/editor.css'
import { templateJsonLintExtension } from './json-lint'
import { templateHighlightExtension, templateTooltipExtension } from './template-highlight'

import type { ViewUpdate } from '@codemirror/view'

interface UseTemplateCodeMirrorProps {
  initialDoc?: string
  onChange?: (doc: string) => void
  readonly?: boolean
  extensions?: Extension[]
}

export type UpdateHook = (update: ViewUpdate) => void

export function useTemplateCodeMirror(container: Ref<HTMLElement | undefined>, props: UseTemplateCodeMirrorProps) {
  const view = shallowRef<EditorView>()
  const readonlyConfig = new Compartment()
  const updateCallbacks = new Set<UpdateHook>()

  const onViewUpdate = (callback: UpdateHook) => {
    updateCallbacks.add(callback)
    return () => updateCallbacks.delete(callback)
  }

  const eventBusExtension = EditorView.updateListener.of((update) => {
    updateCallbacks.forEach((cb) => cb(update))
    if (update.docChanged) {
      props.onChange?.(update.state.doc.toString())
    }
  })

  onMounted(() => {
    if (!container.value) return

    const startState = EditorState.create({
      doc: props.initialDoc || '',
      extensions: [
        basicSetup,
        json(),
        EditorView.lineWrapping,
        codeMirrorTheme,
        readonlyConfig.of(EditorState.readOnly.of(!!props.readonly)),
        templateHighlightExtension,
        templateTooltipExtension,
        templateJsonLintExtension,
        eventBusExtension,
        ...(props.extensions || []),
      ],
    })

    view.value = new EditorView({
      state: startState,
      parent: container.value,
    })
  })

  onUnmounted(() => {
    view.value?.destroy()
    updateCallbacks.clear()
  })

  return {
    view,
    onViewUpdate,
  }
}
