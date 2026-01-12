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
  min-width: 0;
  min-height: 0;
  overflow: hidden;
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

.pane::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
.pane::-webkit-scrollbar-track {
  background: transparent;
}
.pane::-webkit-scrollbar-thumb {
  background: rgba(128, 128, 128, 0.4);
  border-radius: 4px;
}
.pane::-webkit-scrollbar-thumb:hover {
  background: rgba(128, 128, 128, 0.6);
}
</style>
