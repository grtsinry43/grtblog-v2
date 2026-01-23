<script setup lang="ts">
import { NCard, NSplit, NButton, NScrollbar, NTag } from 'naive-ui'
import { onMounted, ref } from 'vue'

import packageJson from '@/../package.json'
import { ScrollContainer } from '@/components'
import { useInjection } from '@/composables'
import { mediaQueryInjectionKey } from '@/injection'

defineOptions({
  name: 'About',
})

let codeToHtml: any

const { isMaxMd } = useInjection(mediaQueryInjectionKey)

const APP_NAME = import.meta.env.VITE_APP_NAME

const { dependencies, devDependencies } = packageJson

const dependenciesCodeHighlight = ref('')
const devDependenciesCodeHighlight = ref('')

onMounted(async () => {
  if (!codeToHtml) {
    // @ts-ignore
    const shiki = await import('https://cdn.jsdelivr.net/npm/shiki@3.7.0/+esm')
    codeToHtml = shiki.codeToHtml
  }

  codeToHtml(JSON.stringify(dependencies, null, 2), {
    lang: 'json',
    themes: {
      light: 'min-light',
      dark: 'dark-plus',
    },
  })
    .then((result: string) => (dependenciesCodeHighlight.value = result))
    .catch(() => (dependenciesCodeHighlight.value = JSON.stringify(dependencies, null, 2)))

  codeToHtml(JSON.stringify(devDependencies, null, 2), {
    lang: 'json',
    themes: {
      light: 'min-light',
      dark: 'dark-plus',
    },
  })
    .then((result: string) => (devDependenciesCodeHighlight.value = result))
    .catch(() => (devDependenciesCodeHighlight.value = JSON.stringify(devDependencies, null, 2)))
})
</script>
<template>
  <ScrollContainer wrapper-class="flex flex-col gap-y-2">
    <NCard
      :title="`关于 ${APP_NAME}`"
      :size="isMaxMd ? 'small' : undefined"
    >
      <p class="text-base">
        {{ APP_NAME }} 是一个轻盈而优雅的后台管理模板，主要技术栈由
        <a
          href="https://vuejs.org/"
          target="_blank"
        >
          <NButton
            strong
            secondary
            size="small"
            color="#42b883"
          >
            Vue3
          </NButton>
        </a>
        <a
          href="https://www.naiveui.com/"
          target="_blank"
        >
          <NButton
            strong
            secondary
            color="#75B93F"
            size="small"
            style="margin-left: 4px"
          >
            Naive UI
          </NButton>
        </a>
        <a
          href="https://vitejs.dev/"
          target="_blank"
        >
          <NButton
            strong
            secondary
            color="#9499ff"
            size="small"
            style="margin-left: 4px"
          >
            Vite7
          </NButton>
        </a>
        <a
          href="https://tailwindcss.com/"
          target="_blank"
        >
          <NButton
            strong
            secondary
            color="#00bcff"
            size="small"
            class="ml-1!"
          >
            TailwindCSS4
          </NButton>
        </a>
        和
        <NButton
          strong
          secondary
          size="small"
        >
          TypeScript
        </NButton>
        构建。
      </p>
      <p class="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
        grtblog-v2 是对 v1 的系统性重构：回到单体结构、减少依赖与复杂度，以默认 SSG 为主、按需引入 SSR / API，
        面向创作者与读者打造一个可持续维护的内容平台。
      </p>
      <p class="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
        本后台为 lithe admin 的二次开发版本，专为 grtblog-v2 的内容管理、发布与运营流程定制。
      </p>
      <p class="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
        项目由 Go API、SvelteKit 前台、Vue 后台与共享 Markdown 组件能力组成。
      </p>
      <div class="mt-3 text-sm text-neutral-600 dark:text-neutral-400">
        <p>面向用户的核心能力：</p>
        <ul class="mt-1 list-disc pl-5">
          <li>Markdown 写作与组件块：相册、提示框、时间轴、链接卡片、年终卡片等。</li>
          <li>文章与元信息管理：摘要/导语/封面/短链、置顶/热门/原创、分类与标签。</li>
          <li>媒体资源管理：图片与文件上传、预览、重命名与下载。</li>
          <li>账号与安全：JWT + RBAC 权限控制、OAuth 绑定、登录限流与人机校验。</li>
        </ul>
      </div>
      <div class="mt-3 text-sm text-neutral-600 dark:text-neutral-400">
        <p>渲染与更新机制：</p>
        <ul class="mt-1 list-disc pl-5">
          <li>SvelteKit 作为渲染器输出 SSR 页面，后端抓取生成静态 HTML 快照对外发布。</li>
          <li>文章变更触发事件驱动的异步刷新，并通过 WebSocket 推送内容更新。</li>
          <li>前台基于内容哈希校验版本，必要时拉取最新内容。</li>
          <li>规划引入脏路径计算，仅刷新受影响的页面与列表。</li>
        </ul>
      </div>
      <p class="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
        另外提供 WebSocket 房间能力，用于内容互动场景。
      </p>
    </NCard>
    <div class="flex gap-x-2 max-lg:flex-col">
      <NCard
        title="依赖信息"
        :size="isMaxMd ? 'small' : undefined"
      >
        <NSplit
          direction="vertical"
          pane1-class="pb-4"
          pane2-class="pt-4"
          default-size="2"
        >
          <template #1>
            <NTag
              class="mb-4"
              :bordered="false"
              type="info"
              size="small"
              >dependencies</NTag
            >
            <NScrollbar>
              <div v-html="dependenciesCodeHighlight"></div>
            </NScrollbar>
          </template>

          <template #2>
            <NTag
              class="mb-4"
              :bordered="false"
              type="info"
              size="small"
              >devDependencies</NTag
            >
            <NScrollbar>
              <div v-html="devDependenciesCodeHighlight"></div>
            </NScrollbar>
          </template>
          <template #resize-trigger>
            <div
              class="h-px w-full cursor-col-resize bg-neutral-200 transition-[background-color] dark:bg-neutral-700"
            ></div>
          </template>
        </NSplit>
      </NCard>
    </div>
  </ScrollContainer>
</template>
