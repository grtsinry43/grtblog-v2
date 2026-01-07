import { autocompletion } from '@codemirror/autocomplete'
import { markdown } from '@codemirror/lang-markdown'
import { EditorState, Compartment } from '@codemirror/state'
import { oneDark } from '@codemirror/theme-one-dark' // 或者你自己的主题
import { EditorView, basicSetup } from 'codemirror'
import { type Ref, shallowRef, onMounted, onUnmounted, watchEffect } from 'vue'

import { customBlockExtension } from './extensions/custom-block'
import { slashHintExtension } from './extensions/slash-hint'
import './editor.css'
// 引入我们刚才写的扩展
import { slashCommandSource } from './extensions/slash-command'

interface UseCodeMirrorProps {
  initialDoc?: string
  onChange?: (doc: string) => void
  darkMode?: boolean
  readonly?: boolean
}

export function useCodeMirror(container: Ref<HTMLElement | undefined>, props: UseCodeMirrorProps) {
  const view = shallowRef<EditorView>()
  const themeConfig = new Compartment() // 用于动态切换主题
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
        themeConfig.of(props.darkMode ? oneDark : []),
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

  // 监听 DarkMode 变化
  watchEffect(() => {
    if (view.value) {
      view.value.dispatch({
        effects: themeConfig.reconfigure(props.darkMode ? oneDark : []),
      })
    }
  })

  onUnmounted(() => {
    view.value?.destroy()
  })

  return {
    view,
  }
}
