import { syntaxTree } from '@codemirror/language'
import { EditorView } from '@codemirror/view'

import type { ChangeSpec, EditorState } from '@codemirror/state'
// 我们这里只处理 粗体，斜体，代码和删除线，其他样式可以类推添加
export type FormatType = 'bold' | 'italic' | 'code' | 'strike'

interface StyleConfig {
  mark: string
  node: string
}

const STYLES: Record<FormatType, StyleConfig> = {
  bold: { mark: '**', node: 'StrongEmphasis' },
  italic: { mark: '*', node: 'Emphasis' },
  code: { mark: '`', node: 'InlineCode' },
  strike: { mark: '~~', node: 'Strikethrough' },
}

// 需要避让的节点（黑名单）
const IGNORED_NODES = [
  'FencedCode',
  'Blockquote',
  'HorizontalRule',
  'URL',
  'Image',
  'ComponentBlock', // 这里是我们自定义的组件块节点
  'ComponentMarker', // 组件的 ::: 标记
]

/**
 * 获取当前选区激活的样式集合
 */
export function getActiveStyles(view: EditorView): Set<FormatType> {
  const { state } = view
  const { from, to } = state.selection.main
  const active = new Set<FormatType>()

  // 如果选区非空，探测中间；如果为空，探测光标左侧
  const pos = from === to ? from : from + 1

  // resolve(pos, -1) 表示向左偏向查找
  let node = syntaxTree(state).resolve(pos, -1)

  // 向上冒泡查找父节点
  while (node.parent) {
    for (const [key, config] of Object.entries(STYLES)) {
      if (node.name === config.node) {
        active.add(key as FormatType)
      }
    }
    node = node.parent
  }
  return active
}

/**
 * 切换样式
 */
export function toggleFormat(view: EditorView, type: FormatType) {
  const active = getActiveStyles(view)
  const { mark, node } = STYLES[type]

  if (active.has(type)) {
    unwrapStyle(view, mark, node)
  } else {
    wrapStyleSafe(view, mark)
  }
}

/**
 * 拆包逻辑：移除样式
 */
function unwrapStyle(view: EditorView, mark: string, nodeName: string) {
  const { state, dispatch } = view
  const changes: ChangeSpec[] = []
  const { from, to } = state.selection.main

  syntaxTree(state).iterate({
    from,
    to,
    enter: (node) => {
      if (node.name === nodeName) {
        const text = state.sliceDoc(node.from, node.to)
        // 确保确实是用这个标记包裹的
        if (text.startsWith(mark) && text.endsWith(mark)) {
          changes.push(
            { from: node.from, to: node.from + mark.length, insert: '' }, // 删头
            { from: node.to - mark.length, to: node.to, insert: '' }, // 删尾
          )
        }
      }
    },
  })

  if (changes.length) dispatch({ changes, userEvent: 'format.unwrap' })
}

/**
 * 智能包裹逻辑：避让组件和特殊块
 */
function wrapStyleSafe(view: EditorView, mark: string) {
  const { state, dispatch } = view
  const changes: ChangeSpec[] = []
  const range = state.selection.main

  if (range.empty) return // 空选区不处理

  let lastPos = range.from

  // 遍历 AST
  syntaxTree(state).iterate({
    from: range.from,
    to: range.to,
    enter: (node) => {
      if (IGNORED_NODES.includes(node.name)) {
        // 1. 处理黑名单节点之前的纯文本
        if (lastPos < node.from) {
          changes.push(...createWrapChanges(state, lastPos, node.from, mark))
        }
        // 2. 跳过黑名单节点
        lastPos = node.to
        return false // 不再进入子节点
      }
    },
  })

  // 3. 处理剩下的部分
  if (lastPos < range.to) {
    changes.push(...createWrapChanges(state, lastPos, range.to, mark))
  }

  if (changes.length) dispatch({ changes, userEvent: 'format.wrap' })
}

/**
 * 这个函数用于生成包裹变更 (也就是剔除首尾空格)
 * @param state 编辑器状态
 * @param from 选区开始绝对坐标
 * @param to 选区结束绝对坐标
 * @param mark 包裹符号 (如 "**")
 */
function createWrapChanges(
  state: EditorState,
  from: number,
  to: number,
  mark: string,
): ChangeSpec[] {
  // 1. 获取选区内的原始文本
  const text = state.sliceDoc(from, to)

  // 2. 边界检查：如果全是空白字符，直接忽略，不包裹
  // 避免出现 "** **" 这种无意义且可能破坏排版的内容
  if (!text.trim()) {
    return []
  }

  // 3. 计算前导空格长度
  // 比如 "  text" -> search 返回 2
  // search(/\S|$/) 意思是查找第一个非空字符，找不到就返回字符串末尾
  const leadingSpaceLen = text.search(/\S|$/)

  // 4. 计算尾部空格长度
  // 比如 "text  " -> 总长 6 - trimEnd后长度 4 = 2
  const trailingSpaceLen = text.length - text.trimEnd().length

  // 5. 计算实际插入点的绝对坐标
  // 开始插入点 = 选区起点 + 前面的空格数
  const insertFrom = from + leadingSpaceLen

  // 结束插入点 = 选区终点 - 后面的空格数
  const insertTo = to - trailingSpaceLen

  // 6. 返回变更数组
  return [
    { from: insertFrom, insert: mark },
    { from: insertTo, insert: mark },
  ]
}
