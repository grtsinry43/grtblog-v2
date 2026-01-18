<script setup lang="ts">
import { TextBold20Regular as IconBold, TextItalic20Regular as IconItalic } from '@vicons/fluent'
import { CodeSlash as IconCode, RemoveCircleOutline as IconStrike } from '@vicons/ionicons5'
import { NButton, NIcon, useThemeVars } from 'naive-ui'
import { computed } from 'vue'

import type { FormatType } from '@/composables/markdown-editor/utils/ast'

interface Props {
  visible: boolean
  pos: { top: number; left: number }
  activeFormats: Set<FormatType>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'command', type: FormatType): void
}>()

const themeVars = useThemeVars()

const style = computed(() => ({
  top: `${props.pos.top}px`,
  left: `${props.pos.left}px`,
  transform: 'translate(-50%, -120%)',
  position: 'fixed' as const,
  zIndex: 9999,
}))

const buttons = [
  { type: 'bold', icon: IconBold, tip: '加粗' },
  { type: 'italic', icon: IconItalic, tip: '斜体' },
  { type: 'strike', icon: IconStrike, tip: '删除线' },
  { type: 'code', icon: IconCode, tip: '行内代码' },
] as const

const handleMousedown = () => {
  // 阻止默认行为
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade-scale">
      <div
        v-if="visible"
        class="floating-menu"
        :style="style"
        @mousedown.prevent="handleMousedown"
      >
        <div class="menu-content">
          <n-button
            v-for="btn in buttons"
            :key="btn.type"
            quaternary
            size="small"
            class="action-btn"
            :class="{ 'is-active': activeFormats.has(btn.type as FormatType) }"
            @click="emit('command', btn.type as FormatType)"
          >
            <template #icon>
              <n-icon :component="btn.icon" />
            </template>
          </n-button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.floating-menu {
  pointer-events: auto;
}

.menu-content {
  /* * 核心改动：回归纯色
   * 使用 popoverColor 确保背景是不透明的实色（深色模式下是深灰，浅色下是纯白）
   */
  background-color: v-bind('themeVars.popoverColor');

  /* 严格使用 Naive UI 系统变量 */
  border-radius: v-bind('themeVars.borderRadius');
  box-shadow: v-bind('themeVars.boxShadow2');
  border: 1px solid v-bind('themeVars.dividerColor'); /* 加上边框防止和背景融为一体 */

  display: flex;
  align-items: center;
  padding: 4px;
  gap: 2px;

  transition:
    background-color 0.3s var(--n-bezier),
    box-shadow 0.3s var(--n-bezier);
}

.action-btn {
  /* 按钮圆角跟随全局设置 */
  border-radius: v-bind('themeVars.borderRadius') !important;
  color: v-bind('themeVars.textColor2');
  transition: all 0.2s var(--n-bezier);
}

/* * 激活状态
 * 依然使用 color-mix 是为了让激活色看起来“不刺眼”，
 * 这里是把主色调和背景色混合，而不是和透明混合，保证了没有透明度。
 */
.action-btn.is-active {
  /* 在 srgb 空间混合：15% 的主色 + 85% 的 Popover 背景色 = 看起来像淡色背景，其实是实色 */
  background-color: color-mix(
    in srgb,
    v-bind('themeVars.primaryColor'),
    v-bind('themeVars.popoverColor') 85%
  );
  color: v-bind('themeVars.primaryColor');
}

/* 悬停状态 */
.action-btn:not(.is-active):hover {
  background-color: v-bind('themeVars.hoverColor');
  color: v-bind('themeVars.textColor1');
}

/* 动画保持不变，依旧使用 Naive 标准曲线 */
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition:
    opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: translate(-50%, -100%) scale(0.95);
}
</style>
