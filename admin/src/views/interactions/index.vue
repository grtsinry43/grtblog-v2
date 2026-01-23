<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'CommentInteraction',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="评论与互动">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于查看评论与互动动态，帮助追踪内容反馈与更新节奏。
      </p>
    </NCard>
  </ScrollContainer>
</template>
