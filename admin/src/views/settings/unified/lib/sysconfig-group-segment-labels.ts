/**
 * 配置树分组标题：API 只保证稳定的 path/key（与 DB group_path 段一致），展示文案由 UI 层维护。
 * 后台当前为中文界面；若以后做多语言，可改为 vue-i18n 或按 locale 加载的同类表。
 */
const SEGMENT_LABELS: Readonly<Record<string, string>> = {
  activitypub: 'ActivityPub',
  ai: 'AI',
  base: '基础',
  comment: '评论',
  discovery: '搜索与 AI 阅读',
  email: '邮件',
  federation: '联合',
  friendlink: '友链',
  interaction: '互动',
  limits: '限流',
  notification: '通知',
  policies: '策略',
  prompt: '提示词',
  security: '安全',
  send: '发送',
  site: '站点',
  basic: '基础',
  social: '社交',
  smtp: 'SMTP',
  storage: '存储',
  subscription: '订阅',
  task: '任务',
  telemetry: '遥测',
  turnstile: 'Turnstile 验证',
  upload: '上传',
  webhook: 'Webhook',
}

/** 折叠面板等处的分组标题：有映射用映射，否则回退到接口给的 label，再回退 key。 */
export function titleForSysconfigGroup(group: { key: string; label?: string }): string {
  const k = group.key.trim().toLowerCase()
  if (k && SEGMENT_LABELS[k]) return SEGMENT_LABELS[k]
  return (group.label?.trim() || group.key).trim() || group.key
}
