<script setup lang="ts">
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSpace,
  NSwitch,
  useThemeVars,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'

import EditorFloatingMenu from '@/components/markdown-editor/EditorFloatingMenu.vue'
import { useCodeMirror } from '@/composables/markdown-editor/use-codemirror'
import { useComponentInserter } from '@/composables/markdown-editor/use-component-inserter.ts'
import { useFloatingMenu } from '@/composables/markdown-editor/use-floating-menu'
import { cah } from '@/utils/chromaHelper'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (
    e: 'cursor-change',
    value: {
      line: number
      column: number
      selectionChars: number
      selectionTotal: number
    },
  ): void
}>()

const editorRef = ref<HTMLElement>()
const themeVars = useThemeVars()

// 样式定义
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
}))

// 1. 初始化 CodeMirror
const { view, onViewUpdate } = useCodeMirror(editorRef, {
  initialDoc: props.modelValue,
  onChange: (val) => emit('update:modelValue', val),
  // 点击图标时，触发 inserter
  onComponentEdit: (payload) => inserter.open(payload),
})

// 2. 挂载组件插入逻辑
const inserter = useComponentInserter(view)

// 3. 挂载浮动菜单逻辑
const { isVisible, menuPos, activeFormats, executeCommand } = useFloatingMenu({
  view,
  onViewUpdate,
})

// 4. 光标与选择更新事件
onViewUpdate((update) => {
  const pos = update.state.selection.main.head
  const selection = update.state.selection.main
  const selectionText =
    selection.from === selection.to
      ? ''
      : update.state.doc.sliceString(selection.from, selection.to)

  const line = update.state.doc.lineAt(pos)

  emit('cursor-change', {
    line: line.number,
    column: pos - line.from + 1,
    selectionChars: selectionText.replace(/\s/g, '').length,
    selectionTotal: selectionText.length,
  })
})

// 5. 外部 modelValue 变更同步
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
  <div class="relative h-full">
    <div
      ref="editorRef"
      class="codemirror-wrapper h-full w-full overflow-visible rounded-md"
      :style="editorStyle"
    ></div>

    <NModal
      v-model:show="inserter.state.show"
      style="width: 520px; max-width: 90vw"
    >
      <NCard
        :title="inserter.componentMeta.value?.label ?? '组件参数'"
        size="small"
      >
        <NForm
          label-placement="left"
          label-width="110px"
        >
          <NFormItem
            v-for="attr in inserter.componentMeta.value?.attrs ?? []"
            :key="attr.key"
            :label="attr.label"
          >
            <NSwitch
              v-if="attr.inputType === 'switch'"
              :value="(inserter.state.attrs[attr.key] ?? attr.defaultValue ?? 'false') === 'true'"
              @update:value="(val) => inserter.updateAttr(attr.key, val)"
            />
            <NInput
              v-else
              :value="inserter.state.attrs[attr.key] ?? attr.defaultValue ?? ''"
              :placeholder="attr.placeholder"
              @update:value="(val) => inserter.updateAttr(attr.key, val)"
            />
          </NFormItem>
        </NForm>
        <NSpace justify="end">
          <NButton @click="inserter.close">关闭</NButton>
        </NSpace>
      </NCard>
    </NModal>

    <EditorFloatingMenu
      :visible="isVisible"
      :pos="menuPos"
      :active-formats="activeFormats"
      @command="executeCommand"
    />
  </div>
</template>

<style scoped>
:deep(.cm-editor) {
  height: 100%;
}
:deep(.cm-content) {
  padding-bottom: 2rem;
}
</style>
