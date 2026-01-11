import { autocompletion } from '@codemirror/autocomplete'
import { markdown } from '@codemirror/lang-markdown'
import { EditorState, Compartment } from '@codemirror/state'
import { EditorView, basicSetup } from 'codemirror'
import { type Ref, shallowRef, onMounted, onUnmounted } from 'vue'

import { codeMirrorTheme } from './codemirror-theme'
import { customBlockExtension } from './extensions/custom-block'
// 引入我们刚才写的扩展
import { slashCommandSource } from './extensions/slash-command'
import './editor.css'
import { slashHintExtension } from './extensions/slash-hint'

interface UseCodeMirrorProps {
  initialDoc?: string
  onChange?: (doc: string) => void
  readonly?: boolean
}

export function useCodeMirror(container: Ref<HTMLElement | undefined>, props: UseCodeMirrorProps) {
  const view = shallowRef<EditorView>()
  const readonlyConfig = new Compartment() // 用于动态切换只读

  onMounted(() => {
    if (!container.value) return

    const startState = EditorState.create({
      doc: props.initialDoc || '',
      extensions: [
        // 1. 基础配置
        basicSetup,
        markdown(),

        // 2. 状态管理 Compartments
        codeMirrorTheme,
        readonlyConfig.of(EditorState.readOnly.of(!!props.readonly)),

        // 3. 核心功能：自动补全 (含 Slash 指令)
        autocompletion({
          override: [slashCommandSource], // 注入我们的 slash 逻辑
          icons: false, // 纯文字列表更像 Notion
          defaultKeymap: true,
        }),

        // 4. 核心功能：自定义块高亮
        customBlockExtension,
        slashHintExtension,

        // 5. 事件监听
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            props.onChange?.(update.state.doc.toString())
          }
        }),
      ],
    })

    view.value = new EditorView({
      state: startState,
      parent: container.value,
    })
  })

  onUnmounted(() => {
    view.value?.destroy()
  })

  return {
    view,
  }
}
