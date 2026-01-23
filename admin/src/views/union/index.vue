<script setup lang="ts">
import { NCard } from 'naive-ui'

import { ScrollContainer } from '@/components'

defineOptions({
  name: 'UnionList',
})
</script>

<template>
  <ScrollContainer wrapper-class="p-4">
    <NCard title="联合">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        用于联合能力的配置与联动管理，后续可扩展聚合与协作能力。
      </p>
    </NCard>
  </ScrollContainer>
</template>
