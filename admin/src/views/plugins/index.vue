<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'PluginList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="插件与云函数">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于扩展站点能力与自动化任务接入，预留插件与云函数管理入口。
      </p>
    </NCard>
  </ScrollContainer>
</template>
