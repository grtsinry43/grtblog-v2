<script setup lang="ts">
import { useFullscreen } from '@vueuse/core'

import { ButtonAnimation } from '@/components'

const { isFullscreen, enter, exit } = useFullscreen()
</script>
<template>
  <div>
    <ButtonAnimation
      @click="isFullscreen ? exit() : enter()"
      :title="isFullscreen ? '退出全屏' : '全屏'"
    >
      <span :class="isFullscreen ? 'iconify ph--arrows-in' : 'iconify ph--arrows-out'" />
    </ButtonAnimation>
  </div>
</template>
