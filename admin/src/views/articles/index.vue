<script setup lang="ts">
import { ref } from 'vue'

import MarkdownEditor from '@/components/markdown-editor/MarkdownEditor.vue'
const value = ref('# Hello Markdown Editor')
</script>

<template>
  <markdown-editor :model-value="value" :is-dark="false" />
</template>

<style scoped></style>
