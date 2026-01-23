<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'AdvancedInfo',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="高级信息">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于查看站点高级配置与运行信息，集中展示核心参数与状态。
      </p>
    </NCard>
  </ScrollContainer>
</template>
