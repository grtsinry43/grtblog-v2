<script setup lang="ts">
import { reactive, watch } from 'vue'

import { ButtonAnimation } from '@/components'
import router from '@/router'

const navigationState = reactive({
  canGoBack: true,
  canGoForward: true,
})

const stopWatch = watch(
  () => router.currentRoute.value,
  () => {
    if (!window.navigation) {
      stopWatch()
      return
    }
    navigationState.canGoBack = window.navigation.canGoBack
    navigationState.canGoForward = window.navigation.canGoForward
  },
  {
    immediate: true,
  },
)
</script>
<template>
  <div class="flex gap-x-0.5 pr-2">
    <ButtonAnimation
      size="medium"
      @click="router.back()"
      title="上一页"
      :disabled="!navigationState.canGoBack"
    >
      <span class="iconify size-4.5 ph--arrow-left" />
    </ButtonAnimation>
    <ButtonAnimation
      size="medium"
      @click="router.forward()"
      title="下一页"
      :disabled="!navigationState.canGoForward"
    >
      <span class="iconify size-4.5 ph--arrow-right" />
    </ButtonAnimation>
  </div>
</template>
