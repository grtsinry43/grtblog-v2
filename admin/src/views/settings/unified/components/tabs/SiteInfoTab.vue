<script setup lang="ts">
import { listSysConfigs, updateSysConfigs } from '@/services/sysconfig'

import ConfigPanel from '../ConfigPanel'

const emit = defineEmits<{ 'dirty-change': [dirty: boolean] }>()
</script>

<template>
  <ConfigPanel
    :list-fn="listSysConfigs"
    :update-fn="updateSysConfigs"
    title="站点信息"
    description="网站基础信息、OG 默认值与搜索支持。搜索与 AI 阅读分组可设置 robots.txt 爬取策略和 llms.txt 阅读指南；站点地图和正文目录自动生成。"
    :filter-groups="['site']"
    :exclude-item-keys="['site.theme_extend_info']"
    :on-dirty-change="(dirty: boolean) => emit('dirty-change', dirty)"
  />
</template>
