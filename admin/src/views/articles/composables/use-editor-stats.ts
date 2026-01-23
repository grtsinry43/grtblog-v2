import { ref } from 'vue'

export function useEditorStats() {
  const cursorPos = ref({ line: 1, column: 1 })
  const selectionStats = ref({ chars: 0, total: 0 })
  const statsIdle = ref(true)
  let statsIdleTimer: ReturnType<typeof setTimeout> | null = null

  function handleCursorChange(payload: {
    line: number
    column: number
    selectionChars: number
    selectionTotal: number
  }) {
    cursorPos.value = { line: payload.line, column: payload.column }
    selectionStats.value = { chars: payload.selectionChars, total: payload.selectionTotal }
  }

  function markActivity() {
    statsIdle.value = false
    if (statsIdleTimer) clearTimeout(statsIdleTimer)
    statsIdleTimer = setTimeout(() => {
      statsIdle.value = true
    }, 800)
  }

  function getStats(content: string) {
    const charCount = content.replace(/\s/g, '').length
    const totalCharCount = content.length
    const chineseCharCount = (content.match(/[\u4e00-\u9fff]/g) ?? []).length
    const wordCount = (content.match(/[A-Za-z0-9]+/g) ?? []).length
    const paragraphCount = content
      .split(/\n+/)
      .map((l) => l.trim())
      .filter(Boolean).length
    const readingMinutes = Math.max(1, Math.ceil(charCount / 300))

    return {
      charCount,
      totalCharCount,
      chineseCharCount,
      wordCount,
      paragraphCount,
      readingMinutes,
    }
  }

  return {
    cursorPos,
    selectionStats,
    statsIdle,
    handleCursorChange,
    markActivity,
    getStats,
  }
}
