<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount, render, h, type Component } from 'vue'

import { useMarkdownIt } from '@/composables/markdown-it/use-markdown-it.ts'
import { getMarkdownComponent } from '@/composables/markdown/shared/components'
import { toRefsPreferencesStore } from '@/stores/preferences'

// --- Props 定义 ---
interface Props {
  /** Markdown 源文本 */
  source: string
  /** * 自定义组件映射表
   * key: 对应 markdown 里的 component name (如 UserCard)
   * value: 对应的 Vue 组件对象
   */
  components?: Record<string, Component>
}

const props = withDefaults(defineProps<Props>(), {
  source: '',
  components: () => ({}),
})

// --- 核心逻辑 ---
const containerRef = ref<HTMLElement | null>(null)
// 使用我们封装好的 hook
const { isDark } = toRefsPreferencesStore()
const { html } = useMarkdownIt(() => props.source, {}, [() => isDark.value])
const renderKey = ref(0)

const escapeHtml = (value = '') =>
  value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

const parseProps = (raw?: string) => {
  if (!raw) return {}
  try {
    const parsed = JSON.parse(raw) as Record<string, string>
    return parsed && typeof parsed === 'object' ? parsed : {}
  } catch {
    return {}
  }
}

const buildPropsHtml = (componentName: string | undefined, props: Record<string, string>) => {
  const meta = getMarkdownComponent(componentName)
  const entries = meta?.attrs?.length
    ? meta.attrs.map((attr) => ({
        key: attr.key,
        value: props[attr.key] ?? attr.defaultValue ?? attr.placeholder ?? '',
        isEmpty: props[attr.key] === undefined || props[attr.key] === '',
      }))
    : Object.entries(props).map(([key, value]) => ({
        key,
        value,
        isEmpty: value === '',
      }))

  if (entries.length === 0) return ''

  const tags = entries
    .map((entry) => {
      const safeKey = escapeHtml(entry.key)
      const safeValue = escapeHtml(entry.value)
      const emptyClass = entry.isEmpty ? ' is-empty' : ''
      return `<span class="md-component-fallback__prop${emptyClass}">
        <span class="md-component-fallback__prop-key">${safeKey}</span>
        <span class="md-component-fallback__prop-sep">=</span>
        <span class="md-component-fallback__prop-value">${safeValue}</span>
      </span>`
    })
    .join('')

  return `<div class="md-component-fallback__props">${tags}</div>`
}

// 存储已挂载的组件容器，用于清理，防止内存泄漏
let mountedApps: HTMLElement[] = []

// 清理函数：销毁手动挂载的 Vue 实例
const cleanupComponents = () => {
  mountedApps.forEach((el) => {
    // 使用 Vue 的 render 函数传 null 来卸载组件
    render(null, el)
  })
  mountedApps = []
}

// 挂载动态组件的核心函数
const mountDynamicComponents = async () => {
  await nextTick() // 等待 v-html 渲染完成
  if (!containerRef.value) return

  // 1. 找到所有的占位符 (由 extensions/componentBlock.ts 生成)
  const placeholders = containerRef.value.querySelectorAll<HTMLElement>('.md-component-placeholder')

  placeholders.forEach((el) => {
    const componentName = el.dataset.component
    const componentProps = parseProps(el.dataset.props)
    const ComponentDef = props.components[componentName || '']

    if (ComponentDef) {
      // 2. 清空占位符内容 (如果有的话)
      el.innerHTML = ''

      // 3. 手动渲染 Vue 组件挂载到这个 DOM 节点上
      // 这里可以传递 props，比如把 dataset 里的其他参数传进去
      const vnode = h(ComponentDef, {
        // 如果需要，可以在这里把 el.dataset 传给组件 props
        ...componentProps,
      })

      render(vnode, el)

      // 4. 记录下来以便后续销毁
      mountedApps.push(el)
    } else {
      // 如果找不到组件，显示占位提示
      const safeName = escapeHtml(componentName || 'unknown')
      const propsHtml = buildPropsHtml(componentName || 'unknown', componentProps)
      el.innerHTML = `<div class="md-component-fallback">
        <span class="md-component-fallback__label">自定义组件</span>
        <span class="md-component-fallback__name">${safeName}</span>
        ${propsHtml}
        <span class="md-component-fallback__hint">将在最终页面预览时展示</span>
      </div>`
    }
  })
}

// --- 监听变化 ---
watch(
  html,
  () => {
    // HTML 更新前，先清理旧的组件实例
    cleanupComponents()
    // HTML 更新后，重新挂载
    mountDynamicComponents()
  },
  { immediate: true },
)

watch(
  isDark,
  async () => {
    cleanupComponents()
    renderKey.value += 1
    await nextTick()
    mountDynamicComponents()
  },
  { flush: 'post' },
)

// 组件销毁时彻底清理
onBeforeUnmount(() => {
  cleanupComponents()
})
</script>

<template>
  <div
    class="markdown-preview markdown-body prose max-w-none"
    ref="containerRef"
    :data-theme="isDark ? 'dark' : 'light'"
    :class="{ 'prose-invert': isDark }"
    :key="renderKey"
    v-html="html"
  ></div>
</template>

<style scoped>
/* 在这里引入 github-markdown-css 或者自定义样式
  .markdown-body 是标准类名
*/
.markdown-preview {
  width: 100%;
  height: 100%;
  padding: 20px;
  /* 基础排版优化 */
  line-height: 1.6;
  word-wrap: break-word;
  overflow-wrap: anywhere;
  background-color: transparent;
  color: inherit;
  --md-hljs-fg: #24292f;
  --md-hljs-bg: transparent;
  --md-hljs-comment: #6e7781;
  --md-hljs-keyword: #cf222e;
  --md-hljs-string: #0a3069;
  --md-hljs-number: #0550ae;
  --md-hljs-title: #116329;
  --md-hljs-variable: #953800;
  --md-hljs-attr: #8250df;
  --md-hljs-built-in: #0550ae;
  --md-hljs-literal: #0550ae;
  --md-hljs-type: #953800;
  --md-component-border: rgba(17, 24, 39, 0.2);
  --md-component-bg: rgba(17, 24, 39, 0.03);
  --md-component-name: #111827;
}

.markdown-preview[data-theme='dark'] {
  --md-hljs-fg: #c9d1d9;
  --md-hljs-bg: transparent;
  --md-hljs-comment: #8b949e;
  --md-hljs-keyword: #ff7b72;
  --md-hljs-string: #a5d6ff;
  --md-hljs-number: #79c0ff;
  --md-hljs-title: #7ee787;
  --md-hljs-variable: #ffa657;
  --md-hljs-attr: #d2a8ff;
  --md-hljs-built-in: #79c0ff;
  --md-hljs-literal: #79c0ff;
  --md-hljs-type: #ffa657;
  --md-component-border: rgba(255, 255, 255, 0.2);
  --md-component-bg: rgba(255, 255, 255, 0.04);
  --md-component-name: #e5e7eb;
}

.markdown-preview :deep(.hljs),
.markdown-preview :deep(pre.hljs) {
  color: var(--md-hljs-fg);
  background-color: var(--md-hljs-bg);
}

.markdown-preview :deep(pre.hljs) {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-preview :deep(.hljs-comment),
.markdown-preview :deep(.hljs-quote) {
  color: var(--md-hljs-comment);
}

.markdown-preview :deep(.hljs-keyword),
.markdown-preview :deep(.hljs-selector-tag),
.markdown-preview :deep(.hljs-literal),
.markdown-preview :deep(.hljs-section),
.markdown-preview :deep(.hljs-link) {
  color: var(--md-hljs-keyword);
}

.markdown-preview :deep(.hljs-string),
.markdown-preview :deep(.hljs-title),
.markdown-preview :deep(.hljs-name),
.markdown-preview :deep(.hljs-type),
.markdown-preview :deep(.hljs-attribute),
.markdown-preview :deep(.hljs-symbol),
.markdown-preview :deep(.hljs-bullet),
.markdown-preview :deep(.hljs-addition) {
  color: var(--md-hljs-string);
}

.markdown-preview :deep(.hljs-number),
.markdown-preview :deep(.hljs-meta),
.markdown-preview :deep(.hljs-built_in),
.markdown-preview :deep(.hljs-builtin-name) {
  color: var(--md-hljs-number);
}

.markdown-preview :deep(.hljs-title),
.markdown-preview :deep(.hljs-class .hljs-title) {
  color: var(--md-hljs-title);
}

.markdown-preview :deep(.hljs-variable),
.markdown-preview :deep(.hljs-template-variable) {
  color: var(--md-hljs-variable);
}

.markdown-preview :deep(.hljs-attr),
.markdown-preview :deep(.hljs-attribute) {
  color: var(--md-hljs-attr);
}

.markdown-preview :deep(.hljs-built_in),
.markdown-preview :deep(.hljs-builtin-name) {
  color: var(--md-hljs-built-in);
}

.markdown-preview :deep(.hljs-literal) {
  color: var(--md-hljs-literal);
}

.markdown-preview :deep(.hljs-type) {
  color: var(--md-hljs-type);
}

.markdown-preview :deep(.hljs-deletion) {
  color: var(--md-hljs-keyword);
}

/* 如果你想给组件占位符一些默认样式 */
:deep(.md-component-placeholder) {
  margin: 1rem 0;
}

:deep(.md-component-fallback) {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 14px;
  border: 1px dashed var(--md-component-border);
  border-radius: 10px;
  background: var(--md-component-bg);
}

:deep(.md-component-fallback__label) {
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: rgba(107, 114, 128, 0.9);
}

:deep(.md-component-fallback__name) {
  font-size: 14px;
  font-weight: 600;
  color: var(--md-component-name);
}

:deep(.md-component-fallback__props) {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 2px;
}

:deep(.md-component-fallback__prop) {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 8px;
  font-size: 11px;
  border-radius: 999px;
  border: 1px solid var(--md-component-border);
  background: var(--md-component-bg);
  color: var(--md-component-name);
}

:deep(.md-component-fallback__prop.is-empty) {
  opacity: 0.6;
}

:deep(.md-component-fallback__prop-key) {
  font-weight: 600;
}

:deep(.md-component-fallback__prop-sep) {
  opacity: 0.6;
}

:deep(.md-component-fallback__hint) {
  font-size: 12px;
  color: rgba(107, 114, 128, 0.9);
}
</style>
