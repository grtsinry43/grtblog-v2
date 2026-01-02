<script setup lang="ts">
import { toRefsPreferencesStore } from '@/stores'

import Action from './action/index.vue'
import Logo from './logo/index.vue'
import Navigation from './navigation/index.vue'

defineOptions({
  name: 'HeaderLayout',
})

const { navigationMode } = toRefsPreferencesStore()
</script>
<template>
  <header class="flex bg-naive-card transition-[background-color]">
    <Logo />
    <div
      class="flex flex-1 items-center border-l px-4 py-3.5 transition-[border-color]"
      :class="navigationMode === 'sidebar' ? 'border-naive-border' : 'border-transparent'"
    >
      <Navigation />
      <Action class="gap-x-3 pl-4" />
    </div>
  </header>
</template>
