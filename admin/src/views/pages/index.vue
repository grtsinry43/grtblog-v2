<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'PageList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="页面管理">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于管理独立页面与内容路由，支持后续扩展模板与布局配置。
      </p>
    </NCard>
  </ScrollContainer>
</template>
