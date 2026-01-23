<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { computed, ref, watch } from 'vue'

import { useTemplateCodeMirror } from '@/composables/template-editor/use-codemirror'
import { cah } from '@/utils/chromaHelper'

const props = defineProps<{
  modelValue: string
  readonly?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const editorRef = ref<HTMLElement>()
const themeVars = useThemeVars()
const editorStyle = computed(() => ({
  '--cm-bg': themeVars.value.cardColor,
  '--cm-fg': themeVars.value.textColor1,
  '--cm-fg-muted': themeVars.value.textColor3,
  '--cm-border': themeVars.value.borderColor,
  '--cm-gutter-fg': themeVars.value.textColor3,
  '--cm-gutter-bg': themeVars.value.cardColor,
  '--cm-selection': cah(themeVars.value.primaryColor, 0.22),
  '--cm-active-line': cah(themeVars.value.primaryColor, 0.02),
  '--cm-cursor': themeVars.value.primaryColor,
  '--cm-tooltip-bg': themeVars.value.popoverColor,
  '--cm-tooltip-border': themeVars.value.borderColor,
  '--cm-tooltip-shadow': themeVars.value.boxShadow2,
  '--cm-tooltip-hover': themeVars.value.buttonColor2Hover,
  '--cm-tooltip-selected': themeVars.value.buttonColor2Pressed,
  '--cm-block-highlight': `color-mix(in srgb, ${themeVars.value.primaryColor} 10%, transparent)`,
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
  borderColor: themeVars.value.borderColor,
}))

const { view } = useTemplateCodeMirror(editorRef, {
  initialDoc: props.modelValue,
  readonly: props.readonly,
  onChange: (val) => {
    emit('update:modelValue', val)
  },
})

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
    class="codemirror-wrapper min-h-[220px] w-full overflow-hidden rounded-md border"
    :style="editorStyle"
  ></div>
</template>

<style scoped>
:deep(.cm-editor) {
  min-height: 220px;
}

:deep(.cm-scroller),
:deep(.cm-content),
:deep(.cm-line) {
  font-family: 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas,
    'Liberation Mono', 'Courier New', monospace;
}

:deep(.cm-template-token) {
  color: var(--cm-syntax-keyword, #7c3aed) !important;
  font-weight: 600;
  background-color: color-mix(in srgb, var(--cm-syntax-keyword, #7c3aed) 10%, transparent);
  border-radius: 3px;
  padding: 0 1px;
}

:deep(.cm-template-invalid) {
  color: var(--cm-syntax-invalid, #ef4444) !important;
  text-decoration: underline wavy var(--cm-syntax-invalid, #ef4444);
  background-color: color-mix(in srgb, var(--cm-syntax-invalid, #ef4444) 10%, transparent);
  border-radius: 3px;
  padding: 0 1px;
}

:deep(.cm-lintRange) {
  text-decoration: underline wavy var(--cm-syntax-invalid, #ef4444);
}

:deep(.cm-template-tooltip) {
  font-size: 12px;
  color: var(--cm-fg, #1f2328);
}
</style>
