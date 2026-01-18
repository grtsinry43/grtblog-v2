<script setup lang="ts">
import { ref } from 'vue'

import MarkdownEditor from '@/components/markdown-editor/MarkdownEditor.vue'
import MarkdownPreview from '@/components/markdown-editor/MarkdownPreview.vue'

const value = ref('# Hello Markdown Editor\n\n试着输入一些内容...')
</script>

<template>
  <div class="editor-container">
    <div class="pane editor-pane">
      <markdown-editor
        v-model="value"
        class="h-full"
      />
    </div>

    <div class="divider"></div>

    <div class="pane preview-pane">
      <markdown-preview :source="value" />
    </div>
  </div>
</template>

<style scoped>
.editor-container {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 1px minmax(0, 1fr);
  height: 100%;
  min-height: 0;
  width: 100%;
  overflow: hidden;
}

.pane {
  height: 100%;
  min-height: 0;
  overflow: auto;
  position: relative;
}

.preview-pane {
  padding: 0;
}

@media (max-width: 768px) {
  .editor-container {
    grid-template-columns: 1fr;
    grid-template-rows: minmax(0, 1fr) 1px minmax(0, 1fr);
  }
}

.pane::-webkit-scrollbar,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar) {
  width: 3px;
  height: 3px;
}
.pane::-webkit-scrollbar-track,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-track),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-track) {
  background: transparent;
}
.pane::-webkit-scrollbar-thumb,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb) {
  background: color-mix(
    in srgb,
    var(--primary-color),
    var(--popover-color, var(--card-color)) 85%
  );
  border-radius: 0;
  transition: background-color 160ms ease;
}
.pane::-webkit-scrollbar-thumb:hover,
.editor-pane :deep(.cm-scroller::-webkit-scrollbar-thumb:hover),
.preview-pane :deep(.markdown-preview::-webkit-scrollbar-thumb:hover) {
  background: color-mix(
    in srgb,
    var(--primary-color),
    var(--popover-color, var(--card-color)) 70%
  );
}
.pane {
  scrollbar-width: thin;
  scrollbar-color: color-mix(
      in srgb,
      var(--primary-color),
      var(--popover-color, var(--card-color)) 85%
    )
    transparent;
}

.editor-pane :deep(.cm-editor) {
  height: auto !important;
  min-height: 100%;
}

.editor-pane :deep(.cm-scroller) {
  height: auto !important;
  overflow: visible !important;
}

.preview-pane :deep(.markdown-preview) {
  height: auto !important;
  overflow: visible !important;
}
</style>
