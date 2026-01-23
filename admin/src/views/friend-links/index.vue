<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'FriendLinkList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="友链">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于处理友链申请与展示配置，维护站点之间的推荐关系。
      </p>
    </NCard>
  </ScrollContainer>
</template>
