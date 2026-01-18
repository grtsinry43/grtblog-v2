import { ref, type Ref } from 'vue'

import { getActiveStyles, toggleFormat, type FormatType } from './utils/ast'

import type { UpdateHook } from './use-codemirror'
import type { EditorView, ViewUpdate } from '@codemirror/view'

// 定义入参接口
interface EditorContext {
  view: Ref<EditorView | undefined>
  onViewUpdate: (callback: UpdateHook) => void
}

export function useFloatingMenu({ view, onViewUpdate }: EditorContext) {
  const isVisible = ref(false)
  const menuPos = ref({ top: 0, left: 0 })
  const activeFormats = ref<Set<FormatType>>(new Set())

  onViewUpdate((update: ViewUpdate) => {
    const currentView = update.view
    const { selection } = currentView.state

    // 1. 基础显隐判断
    if (selection.main.empty) {
      isVisible.value = false
      return
    }

    // 2. 计算坐标
    const start = currentView.coordsAtPos(selection.main.from)
    const end = currentView.coordsAtPos(selection.main.to)

    if (!start || !end) {
      isVisible.value = false
      return
    }

    // 3. 更新 UI 位置
    menuPos.value = {
      left: (start.left + end.right) / 2,
      top: start.top - 6,
    }

    // 4. 状态回填 (查询 AST)
    activeFormats.value = getActiveStyles(currentView)
    isVisible.value = true
  })

  // 执行命令
  const executeCommand = (type: FormatType) => {
    if (view.value) {
      toggleFormat(view.value, type)
      view.value.focus() // 保持焦点
    }
  }

  return {
    isVisible,
    menuPos,
    activeFormats,
    executeCommand,
  }
}
