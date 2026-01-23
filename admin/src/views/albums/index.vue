<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'AlbumList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="相册管理">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于管理相册与图片集内容，便于统一展示与维护。
      </p>
    </NCard>
  </ScrollContainer>
</template>
