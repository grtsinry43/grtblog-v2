import { CompletionContext, type CompletionResult } from '@codemirror/autocomplete'
import { syntaxTree } from '@codemirror/language'

// ... options 定义保持不变 ...
const options = [
  { label: 'Heading 1', type: 'keyword', apply: '# ', detail: '一级标题' },
  { label: 'Heading 2', type: 'keyword', apply: '## ', detail: '二级标题' },
  { label: 'Code Block', type: 'keyword', apply: '```\n\n```', detail: '代码块' },
  { label: 'Quote', type: 'keyword', apply: '> ', detail: '引用' },
  { label: 'Gallery', type: 'variable', apply: '::: gallery\n\n:::', detail: '相册组件' },
  { label: 'Callout', type: 'variable', apply: '::: callout\n\n:::', detail: '提示框' },
  { label: 'Timeline', type: 'variable', apply: '::: timeline\n\n:::', detail: '时间轴' },
]

export const slashCommandSource = (context: CompletionContext): CompletionResult | null => {
  const { state, pos } = context

  // 1. 获取当前行
  const line = state.doc.lineAt(pos)
  const offset = pos - line.from
  const textBefore = line.text.slice(0, offset)

  // 2. 正则匹配：以 / 结尾，且前面是行首或空格
  // 捕获组 match[1] 是斜线后的内容
  const match = /(?:^|\s)\/(\w*)$/.exec(textBefore)

  if (!match) return null

  // 3. 排除代码块和注释区域
  const tree = syntaxTree(state)
  const node = tree.resolveInner(pos, -1)
  if (node.name.includes('Code') || node.name.includes('Comment')) {
    return null
  }

  // 4. 获取搜索词 (斜线后的部分)
  const query = match[1]

  // 5. 【关键步骤】手动筛选选项
  // 如果没有这一步，输入 /h 时，会显示所有选项，体验不好
  // 如果完全依赖 CM 的自动筛选，因为包含了 "/"，会导致所有选项都不匹配
  const filteredOptions = options.filter((option) => {
    const searchStr = query.toLowerCase()
    return (
      option.label.toLowerCase().includes(searchStr) ||
      option.detail.toLowerCase().includes(searchStr)
    )
  })

  // 如果筛选后没有结果，返回 null 关闭菜单
  if (filteredOptions.length === 0) return null

  return {
    // from 包含斜线，这样选中后会把 "/" 替换掉
    from: line.from + match.index + (match[0].startsWith(' ') ? 1 : 0),
    to: pos,
    options: filteredOptions,

    // 【核心修复】关闭 CM 默认过滤
    // 因为我们的 range 包含了 "/"，默认过滤器会拿 "/" 去匹配 "Heading"，导致匹配失败
    filter: false,

    // 只有当斜线后是字母数字时才继续触发此补全
    // 这行代码并不是必须的，因为 filter: false 后逻辑由我们控制，但加上性能更好
    // validFor: /^\/(\w*)?$/
  }
}
