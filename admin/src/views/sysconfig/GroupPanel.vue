<script setup lang="ts">
import { NCollapse, NCollapseItem, NForm } from 'naive-ui'

import type { SysConfigGroup, SysConfigItem } from '@/services/sysconfig'

defineOptions({
  name: 'GroupPanel',
})

defineProps<{
  node: SysConfigGroup
  showItem: (item: SysConfigItem) => boolean
}>()
</script>

<template>
  <NCollapseItem :name="node.path" :title="node.label || node.key">
    <div class="space-y-4">
      <NForm v-if="node.items?.length" label-placement="left" label-width="160">
        <template v-for="item in node.items" :key="item.key">
          <slot v-if="showItem(item)" name="item" :item="item" />
        </template>
      </NForm>

      <NCollapse v-if="node.children?.length">
        <GroupPanel
          v-for="child in node.children"
          :key="child.path"
          :node="child"
          :show-item="showItem"
        >
          <template #item="slotProps">
            <slot name="item" v-bind="slotProps" />
          </template>
        </GroupPanel>
      </NCollapse>
    </div>
  </NCollapseItem>
</template>
