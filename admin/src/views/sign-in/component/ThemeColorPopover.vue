<script setup lang="ts">
import { NPopover } from 'naive-ui'

import { ButtonAnimation } from '@/components'
import { toRefsPreferencesStore } from '@/stores'
import twColors from '@/utils/tailwindColor'

const { themeColor } = toRefsPreferencesStore()

const colorSwatches = [
  twColors.red[500],
  twColors.orange[500],
  twColors.amber[500],
  twColors.yellow[500],
  twColors.lime[500],
  twColors.green[500],
  twColors.emerald[500],
  twColors.teal[500],
  twColors.cyan[500],
  twColors.sky[500],
  twColors.blue[500],
  twColors.indigo[500],
  twColors.violet[500],
  twColors.purple[500],
  twColors.fuchsia[500],
  twColors.pink[500],
]
</script>
<template>
  <NPopover
    trigger="click"
    placement="bottom"
    :to="false"
  >
    <template #trigger>
      <ButtonAnimation>
        <span class="iconify ph--palette" />
      </ButtonAnimation>
    </template>
    <div class="grid grid-cols-6 gap-1.5">
      <div
        v-for="color in colorSwatches"
        :key="color"
        class="size-5 cursor-pointer rounded-full"
        :class="{
          'ring-3 ring-primary/30': color === themeColor,
        }"
        :style="{ backgroundColor: color }"
        @click="themeColor = color"
      />
    </div>
  </NPopover>
</template>
