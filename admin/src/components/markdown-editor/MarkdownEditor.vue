<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { computed, ref, watch } from 'vue'

import { useCodeMirror } from '@/composables/markdown-editor/use-codemirror.ts'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

const editorRef = ref<HTMLElement>()
const themeVars = useThemeVars()
const editorStyle = computed(() => ({
  '--cm-bg': themeVars.value.cardColor,
  '--cm-fg': themeVars.value.textColor1,
  '--cm-fg-muted': themeVars.value.textColor3,
  '--cm-border': themeVars.value.borderColor,
  '--cm-gutter-fg': themeVars.value.textColor3,
  '--cm-gutter-bg': themeVars.value.cardColor,
  '--cm-selection': themeVars.value.pressedColor,
  '--cm-active-line': themeVars.value.hoverColor,
  '--cm-cursor': themeVars.value.primaryColor,
  '--cm-tooltip-bg': themeVars.value.popoverColor,
  '--cm-tooltip-border': themeVars.value.borderColor,
  '--cm-tooltip-shadow': themeVars.value.boxShadow2,
  '--cm-tooltip-hover': themeVars.value.buttonColor2Hover,
  '--cm-tooltip-selected': themeVars.value.buttonColor2Pressed,
  '--cm-block-highlight': `color-mix(in srgb, ${themeVars.value.primaryColor} 18%, transparent)`,
  '--cm-accent': themeVars.value.primaryColor,
  '--cm-syntax-keyword': themeVars.value.primaryColor,
  '--cm-syntax-string': themeVars.value.successColor,
  '--cm-syntax-number': themeVars.value.warningColor,
  '--cm-syntax-property': themeVars.value.infoColor,
  '--cm-syntax-function': themeVars.value.infoColor,
  '--cm-syntax-type': themeVars.value.primaryColor,
  '--cm-syntax-tag': themeVars.value.errorColor,
  '--cm-syntax-invalid': themeVars.value.errorColor,
  '--cm-syntax-comment': themeVars.value.textColor3,
  '--cm-syntax-variable': themeVars.value.textColor2,
  '--cm-syntax-operator': themeVars.value.textColor2,
  '--cm-syntax-punctuation': themeVars.value.textColor3,
}))

// 优雅的调用
const { view } = useCodeMirror(editorRef, {
  initialDoc: props.modelValue,
  onChange: (val) => {
    emit('update:modelValue', val)
  },
})

// 如果父组件通过代码修改了 modelValue (非输入触发)，需要同步回编辑器
// 注意：要避免光标跳动，需要判断内容是否一致
watch(
  () => props.modelValue,
  (newVal) => {
    if (view.value && view.value.state.doc.toString() !== newVal) {
      view.value.dispatch({
        changes: { from: 0, to: view.value.state.doc.length, insert: newVal },
      })
    }
  },
)
</script>

<template>
  <div
    ref="editorRef"
    class="codemirror-wrapper h-full w-full overflow-visible rounded-md"
    :style="editorStyle"
  ></div>
</template>

<style scoped>
:deep(.cm-editor) {
  height: 100%;
}

:deep(.cm-content) {
  padding-bottom: 2rem;
}
</style>
