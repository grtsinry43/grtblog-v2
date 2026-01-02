<script setup lang="ts">
import { provide, useAttrs } from 'vue'

import { buttonAnimationInjectionKey } from './injection'

import type { ButtonAnimationProps } from './interface'

defineProps<ButtonAnimationProps>()

defineOptions({
  inheritAttrs: false,
})

provide(buttonAnimationInjectionKey, useAttrs() as ButtonAnimationProps)
</script>
<template>
  <slot />
</template>
