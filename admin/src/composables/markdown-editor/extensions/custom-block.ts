import { RangeSet } from '@codemirror/state'
import { Decoration, type DecorationSet, EditorView, ViewPlugin, ViewUpdate } from '@codemirror/view'

const blockLine = Decoration.line({ class: 'cm-custom-block' })
const blockLineStart = Decoration.line({ class: 'cm-custom-block cm-custom-block-start' })
const blockLineEnd = Decoration.line({ class: 'cm-custom-block cm-custom-block-end' })
const labelMark = Decoration.mark({ class: 'cm-custom-block-label' })

const VALID_COMPONENTS = ['gallery', 'callout', 'timeline']

// 2. 构建逻辑：使用行扫描代替全文正则
function buildDecorations(view: EditorView): DecorationSet {
  const decorations: { from: number; to: number; value: Decoration }[] = []
  const doc = view.state.doc

  // 状态机变量
  let inBlock = false
  let startLineNumber = -1 // 记录块开始行号

  // 遍历每一行 (性能远优于 regex 全文匹配)
  for (let i = 1; i <= doc.lines; i++) {
    const line = doc.line(i)
    const text = line.text.trim() // 去除首尾空格方便判断

    // 检查是否是开始标记 ::: name 或 ::: component name
    if (text.startsWith(':::')) {
      // 提取组件名，处理 ":::gallery" / "::: gallery" / "::: component gallery"
      const content = line.text.trim().slice(3).trim()
      const parts = content.split(/\s+/)
      const componentName = parts[0] === 'component' ? parts[1] : parts[0]

      if (VALID_COMPONENTS.includes(componentName)) {
        // 这是一个合法的开始标记
        if (!inBlock) {
          inBlock = true
          startLineNumber = i

          // 【视觉优化】给关键字 "gallery" 加高亮
          // 重新定位关键字在行内的确切位置
          const headerText = line.text
          const nameIndex = headerText.indexOf(componentName)
          if (nameIndex !== -1) {
            decorations.push({
              from: line.from + nameIndex,
              to: line.from + nameIndex + componentName.length,
              value: labelMark,
            })
          }
        } else {
          // 如果已经在 block 里又遇到了 :::，说明上一个 block 结束了（或者嵌套了）
          // 这里我们做简单处理：遇到新的 ::: 视为结束上一个，开始下一个，或者作为结束标记
          // 简单的 ::: 单独一行通常是结束
          if (content === '') {
            // 结束标记
            inBlock = false
            for (let lineIndex = startLineNumber; lineIndex <= i; lineIndex++) {
              const targetLine = doc.line(lineIndex)
              const isStart = lineIndex === startLineNumber
              const isEnd = lineIndex === i
              const lineDecoration = isStart
                ? blockLineStart
                : isEnd
                  ? blockLineEnd
                  : blockLine
              decorations.push({
                from: targetLine.from,
                to: targetLine.from,
                value: lineDecoration,
              })
            }
          }
        }
      } else if (text === ':::' && inBlock) {
        // 这是一个结束标记
        inBlock = false
        for (let lineIndex = startLineNumber; lineIndex <= i; lineIndex++) {
          const targetLine = doc.line(lineIndex)
          const isStart = lineIndex === startLineNumber
          const isEnd = lineIndex === i
          const lineDecoration = isStart
            ? blockLineStart
            : isEnd
              ? blockLineEnd
              : blockLine
          decorations.push({
            from: targetLine.from,
            to: targetLine.from,
            value: lineDecoration,
          })
        }
      }
    }
  }

  // 3. 关键修复：使用 RangeSet.of(arr, true)
  // 第二个参数 true 表示“请帮我排序并处理重叠”。
  // 这允许 blockMark (大范围) 和 labelMark (小范围) 同时存在。
  return RangeSet.of(decorations, true)
}

// 4. 插件定义
export const customBlockPlugin = ViewPlugin.fromClass(
  class {
    decorations: DecorationSet
    constructor(view: EditorView) {
      this.decorations = buildDecorations(view)
    }
    update(update: ViewUpdate) {
      // 只有文档内容变了才重新计算 (视口变化不需要重算，因为我们扫描了全文档)
      if (update.docChanged) {
        this.decorations = buildDecorations(update.view)
      }
    }
  },
  {
    decorations: (v) => v.decorations,
  },
)

export const customBlockExtension = [customBlockPlugin]
