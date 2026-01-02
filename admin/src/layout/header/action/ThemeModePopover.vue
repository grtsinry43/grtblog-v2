<script setup lang="ts">
import { NPopselect } from 'naive-ui'
import { computed, h } from 'vue'

import { ButtonAnimation } from '@/components'
import { toRefsPreferencesStore } from '@/stores'

import type { PopoverProps } from 'naive-ui'

interface ThemeModePopover extends /* @vue-ignore */ PopoverProps {}

defineProps<ThemeModePopover>()

defineOptions({
  inheritAttrs: false,
})

const { themeMode } = toRefsPreferencesStore()

const themeModeDropdownOptions = [
  {
    icon: () => h('span', { class: 'iconify ph--sun size-5' }),
    iconName: 'iconify ph--sun',
    key: 'light',
    value: 'light',
    label: '浅色模式',
  },
  {
    icon: () => h('span', { class: 'iconify ph--moon size-5' }),
    iconName: 'iconify ph--moon',
    key: 'dark',
    value: 'dark',
    label: '深色模式',
  },
  {
    icon: () => h('span', { class: 'iconify ph--desktop size-5' }),
    iconName: 'iconify ph--desktop',
    key: 'auto',
    value: 'auto',
    label: '跟随系统',
  },
]

const themeIconClasses = computed(
  () =>
    themeModeDropdownOptions.find((item) => item.key === themeMode.value)?.iconName ||
    'iconify ph--desktop size-5',
)

function renderSelectLabel(option: (typeof themeModeDropdownOptions)[number]) {
  return h(
    'div',
    {
      class: 'flex items-center gap-x-2',
    },
    [option.icon(), option.label],
  )
}
</script>
<template>
  <NPopselect
    class="p-0.5"
    trigger="click"
    v-bind="$attrs"
    v-model:value="themeMode"
    :options="themeModeDropdownOptions"
    :render-label="renderSelectLabel"
    :to="false"
  >
    <ButtonAnimation title="主题">
      <span :class="themeIconClasses" />
    </ButtonAnimation>
  </NPopselect>
</template>
