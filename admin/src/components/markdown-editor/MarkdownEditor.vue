<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NInput, NModal, NSpace, NSwitch, useThemeVars } from 'naive-ui'
import { computed, reactive, ref, watch } from 'vue'

import { useCodeMirror } from '@/composables/markdown-editor/use-codemirror.ts'
import { getMarkdownComponent } from '@/composables/markdown/shared/components'
import { cah } from '@/utils/chromaHelper'

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

const componentEditor = reactive({
  show: false,
  name: '',
  lineFrom: 0,
  lineTo: 0,
  isComponentSyntax: false,
  attrs: {} as Record<string, string>,
})
const componentTouchedKeys = new Set<string>()

const componentMeta = computed(() => getMarkdownComponent(componentEditor.name))

const formatComponentLine = () => {
  const meta = componentMeta.value
  const definedKeys = meta?.attrs?.map((attr) => attr.key) ?? []
  const extraKeys = Object.keys(componentEditor.attrs).filter((key) => !definedKeys.includes(key))
  const keys = [...definedKeys, ...extraKeys].filter(
    (key) => key in componentEditor.attrs || componentTouchedKeys.has(key),
  )
  const attrsString = keys
    .map((key) => `${key}="${componentEditor.attrs[key] ?? ''}"`)
    .join(' ')
  const prefix = componentEditor.isComponentSyntax
    ? `::: component ${componentEditor.name}`
    : `::: ${componentEditor.name}`
  return attrsString ? `${prefix}{${attrsString}}` : prefix
}

const applyComponentUpdate = () => {
  if (!view.value || !componentEditor.show) return
  const newLine = formatComponentLine()
  view.value.dispatch({
    changes: { from: componentEditor.lineFrom, to: componentEditor.lineTo, insert: newLine },
  })
  componentEditor.lineTo = componentEditor.lineFrom + newLine.length
}

const updateAttr = (key: string, value: string | boolean) => {
  componentTouchedKeys.add(key)
  componentEditor.attrs[key] = String(value)
  applyComponentUpdate()
}

// 优雅的调用
const { view } = useCodeMirror(editorRef, {
  initialDoc: props.modelValue,
  onChange: (val) => {
    emit('update:modelValue', val)
  },
  onComponentEdit: (payload) => {
    componentEditor.name = payload.name
    componentEditor.lineFrom = payload.lineFrom
    componentEditor.lineTo = payload.lineTo
    componentEditor.isComponentSyntax = payload.isComponentSyntax
    componentEditor.attrs = { ...payload.attrs }
    componentTouchedKeys.clear()
    componentEditor.show = true
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
  <div>
    <div
      ref="editorRef"
      class="codemirror-wrapper h-full w-full overflow-visible rounded-md"
      :style="editorStyle"
    ></div>

    <NModal v-model:show="componentEditor.show" style="width: 520px; max-width: 90vw;">
      <NCard
        :title="componentMeta?.label ?? '组件参数'"
        size="small"
      >
        <NForm label-placement="left" label-width="110px">
          <NFormItem
            v-for="attr in componentMeta?.attrs ?? []"
            :key="attr.key"
            :label="attr.label"
          >
            <NSwitch
              v-if="attr.inputType === 'switch'"
              :value="(componentEditor.attrs[attr.key] ?? attr.defaultValue ?? 'false') === 'true'"
              @update:value="(val) => updateAttr(attr.key, val)"
            />
            <NInput
              v-else
              :value="componentEditor.attrs[attr.key] ?? attr.defaultValue ?? ''"
              :placeholder="attr.placeholder"
              @update:value="(val) => updateAttr(attr.key, val)"
            />
          </NFormItem>
        </NForm>
        <NSpace justify="end">
          <NButton @click="componentEditor.show = false">关闭</NButton>
        </NSpace>
      </NCard>
    </NModal>
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
