import { autocompletion } from '@codemirror/autocomplete'
import { markdown } from '@codemirror/lang-markdown'
import { EditorState, Compartment, type Extension } from '@codemirror/state'
import { GFM, Subscript, Superscript } from '@lezer/markdown'
import { EditorView, basicSetup } from 'codemirror'
import { type Ref, shallowRef, onMounted, onUnmounted } from 'vue'

import { codeMirrorTheme } from './codemirror-theme'
import { componentAttributeSource } from './extensions/component-attrs'
import {
  createComponentEditorExtension,
  type ComponentEditPayload,
} from './extensions/component-editor'
import { customBlockExtension } from './extensions/custom-block'
import { slashCommandSource } from './extensions/slash-command'
import './editor.css'
import { slashHintExtension } from './extensions/slash-hint'

import type { ViewUpdate } from '@codemirror/view'

interface UseCodeMirrorProps {
  initialDoc?: string
  onChange?: (doc: string) => void
  readonly?: boolean
  onComponentEdit?: (payload: ComponentEditPayload) => void
  // 允许传入额外的 extensions (如果需要)
  extensions?: Extension[]
}

// 定义钩子类型
export type UpdateHook = (update: ViewUpdate) => void

export function useCodeMirror(container: Ref<HTMLElement | undefined>, props: UseCodeMirrorProps) {
  const view = shallowRef<EditorView>()
  const readonlyConfig = new Compartment()

  const updateCallbacks = new Set<UpdateHook>()

  // 对外暴露的注册函数
  const onViewUpdate = (callback: UpdateHook) => {
    updateCallbacks.add(callback)
    return () => updateCallbacks.delete(callback) // 返回清理函数
  }

  // 内部扩展：统一分发更新事件
  const eventBusExtension = EditorView.updateListener.of((update) => {
    // 触发所有订阅者
    updateCallbacks.forEach((cb) => cb(update))

    // 处理原有的 onChange
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
        markdown({
          extensions: [GFM, Subscript, Superscript],
        }),
        EditorView.lineWrapping,
        codeMirrorTheme,
        readonlyConfig.of(EditorState.readOnly.of(!!props.readonly)),

        // 功能扩展
        autocompletion({
          override: [slashCommandSource, componentAttributeSource],
          icons: false,
          defaultKeymap: true,
        }),
        customBlockExtension,
        slashHintExtension,
        createComponentEditorExtension({ onEdit: props.onComponentEdit }),

        // 注册事件总线扩展
        eventBusExtension,

        // 合并外部传入的扩展
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
    onViewUpdate, // 导出钩子
  }
}
