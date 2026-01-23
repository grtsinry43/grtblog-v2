import { type Diagnostic, linter } from '@codemirror/lint'

function isEscaped(input: string, index: number) {
  let count = 0
  for (let i = index - 1; i >= 0 && input[i] === '\\'; i -= 1) {
    count += 1
  }
  return count % 2 === 1
}

function buildReplacement(length: number, inString: boolean) {
  const base = inString ? '__T__' : 'null'
  if (length <= base.length) {
    return base.slice(0, length)
  }
  return `${base}${' '.repeat(length - base.length)}`
}

export function sanitizeTemplateJson(input: string) {
  let inString = false
  let index = 0
  let output = ''

  while (index < input.length) {
    const char = input[index]
    if (char === '"' && !isEscaped(input, index)) {
      inString = !inString
      output += char
      index += 1
      continue
    }

    if (char === '{' && input[index + 1] === '{') {
      const endOffset = input.indexOf('}}', index + 2)
      if (endOffset !== -1) {
        const tokenLength = endOffset + 2 - index
        output += buildReplacement(tokenLength, inString)
        index = endOffset + 2
        continue
      }
    }

    output += char
    index += 1
  }

  return output
}

type PlaceholderToken = {
  token: string
  placeholder: string
  inString: boolean
}

export function formatTemplateJson(input: string) {
  let inString = false
  let index = 0
  let output = ''
  const placeholders: PlaceholderToken[] = []

  while (index < input.length) {
    const char = input[index]
    if (char === '"' && !isEscaped(input, index)) {
      inString = !inString
      output += char
      index += 1
      continue
    }

    if (char === '{' && input[index + 1] === '{') {
      const endOffset = input.indexOf('}}', index + 2)
      if (endOffset !== -1) {
        const token = input.slice(index, endOffset + 2)
        const placeholder = `__TEMPLATE_${placeholders.length}__`
        placeholders.push({ token, placeholder, inString })
        output += inString ? placeholder : `"${placeholder}"`
        index = endOffset + 2
        continue
      }
    }

    output += char
    index += 1
  }

  const parsed = JSON.parse(output)
  let formatted = JSON.stringify(parsed, null, 2)
  for (const item of placeholders) {
    if (item.inString) {
      formatted = formatted.replaceAll(item.placeholder, item.token)
    } else {
      formatted = formatted.replaceAll(`"${item.placeholder}"`, item.token)
    }
  }
  return formatted
}

function parseErrorPosition(message: string) {
  const match = message.match(/position\s+(\d+)/)
  if (!match) return null
  const value = Number(match[1])
  if (!Number.isFinite(value)) return null
  return value
}

export const templateJsonLintExtension = linter((view) => {
  const raw = view.state.doc.toString()
  if (!raw.trim()) return []

  const sanitized = sanitizeTemplateJson(raw)
  try {
    JSON.parse(sanitized)
    return []
  } catch (err) {
    const message = err instanceof Error ? err.message : 'JSON 格式错误'
    const diagnostics: Diagnostic[] = []
    const position = parseErrorPosition(message)
    if (position !== null) {
      const from = Math.min(Math.max(position, 0), view.state.doc.length)
      const to = Math.min(from + 1, view.state.doc.length)
      diagnostics.push({
        from,
        to,
        severity: 'error',
        message,
      })
    } else {
      diagnostics.push({
        from: 0,
        to: Math.min(1, view.state.doc.length),
        severity: 'error',
        message,
      })
    }
    return diagnostics
  }
})
