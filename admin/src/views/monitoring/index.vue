<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'SystemMonitor',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="系统监控">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于查看系统健康、任务与运行状态，方便监测整体稳定性。
      </p>
    </NCard>
  </ScrollContainer>
</template>
