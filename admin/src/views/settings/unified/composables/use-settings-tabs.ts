import { ref, reactive, watch, type Component, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AdvancedTab from '../components/tabs/AdvancedTab.vue'
import AiTab from '../components/tabs/AiTab.vue'
import ApiTokensTab from '../components/tabs/ApiTokensTab.vue'
import BackupTab from '../components/tabs/BackupTab.vue'
import ContentTab from '../components/tabs/ContentTab.vue'
import EmailTab from '../components/tabs/EmailTab.vue'
import ExportTab from '../components/tabs/ExportTab.vue'
import FederationTab from '../components/tabs/FederationTab.vue'
import SecurityTab from '../components/tabs/SecurityTab.vue'
import SiteInfoTab from '../components/tabs/SiteInfoTab.vue'
import TelemetryTab from '../components/tabs/TelemetryTab.vue'
import ThemeExtendTab from '../components/tabs/ThemeExtendTab.vue'
import WebhookTab from '../components/tabs/WebhookTab.vue'

export interface SettingsTab {
  key: string
  label: string
  icon: string
  component: Component
  /** When true, tab content fills the available height instead of scrolling externally. */
  fillHeight?: boolean
}

export const settingsTabs: SettingsTab[] = [
  {
    key: 'site-info',
    label: '基本信息',
    icon: 'iconify ph--globe-hemisphere-west',
    component: SiteInfoTab,
  },
  {
    key: 'theme-extend',
    label: '主题扩展',
    icon: 'iconify ph--paint-brush',
    component: ThemeExtendTab,
    fillHeight: true,
  },
  {
    key: 'security',
    label: '安全与登录',
    icon: 'iconify ph--shield-check',
    component: SecurityTab,
  },
  { key: 'content', label: '内容与评论', icon: 'iconify ph--article', component: ContentTab },
  { key: 'email', label: '邮件', icon: 'iconify ph--envelope', component: EmailTab },
  { key: 'webhook', label: 'Webhook', icon: 'iconify ph--webhooks-logo', component: WebhookTab },
  { key: 'federation', label: '联合', icon: 'iconify ph--circles-three', component: FederationTab },
  { key: 'ai', label: 'AI', icon: 'iconify ph--robot', component: AiTab },
  { key: 'api-tokens', label: 'API Tokens', icon: 'iconify ph--key', component: ApiTokensTab },
  { key: 'backup', label: '备份与恢复', icon: 'iconify ph--archive', component: BackupTab },
  {
    key: 'content-export',
    label: '内容导出',
    icon: 'iconify ph--package',
    component: ExportTab,
  },
  { key: 'telemetry', label: '遥测', icon: 'iconify ph--chart-line-up', component: TelemetryTab },
  { key: 'advanced', label: '高级', icon: 'iconify ph--gear', component: AdvancedTab },
]

const validKeys = new Set(settingsTabs.map((t) => t.key))

export function useSettingsTabs() {
  const route = useRoute()
  const router = useRouter()

  const initialTab = validKeys.has(route.query.tab as string)
    ? (route.query.tab as string)
    : 'site-info'

  const activeTab = ref(initialTab)
  const dirtyTabs = reactive(new Set<string>())

  const currentTabDef = computed(() => settingsTabs.find((t) => t.key === activeTab.value)!)

  watch(activeTab, (tab) => {
    if (route.query.tab !== tab) {
      const url = router.resolve({ query: { ...route.query, tab } }).href
      history.replaceState(history.state, '', url)
    }
  })

  function switchTab(tabKey: string) {
    if (tabKey === activeTab.value) return
    activeTab.value = tabKey
  }

  function setDirty(tabKey: string, dirty: boolean) {
    if (dirty) {
      dirtyTabs.add(tabKey)
    } else {
      dirtyTabs.delete(tabKey)
    }
  }

  return {
    tabs: settingsTabs,
    activeTab,
    switchTab,
    currentTabDef,
    dirtyTabs,
    setDirty,
  }
}
