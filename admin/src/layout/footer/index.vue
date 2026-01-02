<script setup lang="ts">
defineOptions({
  name: 'FooterLayout',
})

const APP_NAME = import.meta.env.VITE_APP_NAME
</script>
<template>
  <footer
    class="min-h-0 border-t border-naive-border bg-naive-card transition-[background-color,border-color]"
  >
    <div class="flex items-center justify-center overflow-hidden py-1.5 text-xs">
      <span>{{ APP_NAME }}</span>
    </div>
  </footer>
</template>
