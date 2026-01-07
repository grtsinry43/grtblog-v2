<script setup lang="ts">
import { ref, watch } from 'vue'

import { useCodeMirror } from '@/composables/markdown-editor/use-codemirror.ts'

const props = defineProps<{
  modelValue: string
  isDark?: boolean
}>()

const emit = defineEmits(['update:modelValue'])

const editorRef = ref<HTMLElement>()

// 优雅的调用
const { view } = useCodeMirror(editorRef, {
  initialDoc: props.modelValue,
  darkMode: props.isDark,
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
    class="codemirror-wrapper h-full w-full overflow-visible rounded-md border border-gray-200 dark:border-gray-700"
  ></div>
</template>

<style scoped>
/* 也可以在这里微调 CM 的全局字体 */
:deep(.cm-editor) {
  height: 100%;
  font-family: 'Fira Code', monospace;
}
:deep(.cm-content) {
  padding-bottom: 2rem;
}
</style>
