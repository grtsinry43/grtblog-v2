<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'NoteList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="手记管理">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于整理手记内容与草稿，后续可扩展发布与归档能力。
      </p>
    </NCard>
  </ScrollContainer>
</template>
