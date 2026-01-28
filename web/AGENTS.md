# AGENTS.md — GrtBlog v2 (Svelte 5 + SvelteKit) 工程规范

> 目标：把“优雅代码”的原则落到当前仓库结构与命名上，供任何 LLM / Agent 在本项目中遵循。

## 0. 项目定位
- 内容平台 / 博客系统
- 静态优先（SSG/预渲染） + 渐进式增强（hydration islands）
- 统一设计系统 + 可扩展功能模块

## 1. 目录与职责边界（按当前仓库结构）
- `src/routes/**`: 页面编排 + SvelteKit load 数据接入（SSR/SSG/SEO）。不写复杂业务逻辑。
- `src/lib/features/<feature>/**`: 业务域模块（API + types + components）。例如：`post`。
- `src/lib/shared/**`: 跨功能共享能力（clients、markdown、theme、token）。
- `src/lib/ui/**`: 通用 UI 组件（Button/Badge/Tag/...）。
- `src/lib/assets/**`: 资源或静态导入。
- `static/**`: 直接公开静态资源。

## 2. 组合优先：页面只做编排
- 页面组件只负责“拼装”，不要堆业务细节。
- 业务逻辑必须下沉到：
  - `src/lib/features/<feature>/api.ts`（请求与协议）
  - `src/lib/features/<feature>/types.ts`（类型定义）
  - `src/lib/shared/**`（跨域能力）
  - `src/lib/ui/**`（可复用 UI）

## 3. 浏览器 API 必须封装为能力
- 禁止在业务组件里散落 `new IntersectionObserver` / `addEventListener` / `ResizeObserver`。
- 统一封装为可复用能力（建议路径）：
  - `src/lib/shared/dom/*` 或 `src/lib/shared/actions/*`（Svelte actions）
- 必须处理：
  - cleanup（解除监听/observer）
  - SSR 安全（仅在浏览器环境执行）
  - 稳定引用（避免重复绑定）

## 4. 状态流：选择器/派生优先
- 页面级数据：优先通过 `+page.server.ts/+page.ts load` 提供。
- 跨层共享：`setContext/getContext` + store。
- 订阅粒度：优先 `derived store` 或 `$derived`，组件只读最小数据。
- 避免大对象 prop drilling。

## 5. Svelte 5 Runes 约束
- `$state`: 本地可变状态
- `$derived`: 派生状态（避免手写同步）
- `$effect`: 只做副作用（DOM/网络/订阅），不做纯计算
- `$props`: 统一 props 访问

## 6. 数据获取与 API 约定
- API 请求统一通过 `src/lib/shared/clients/api.ts` 的 `getApi(fetch)`。
- feature 模块 API 参考 `src/lib/features/post/api.ts`（优先接收 `fetch`）。
- 服务端 load 中使用 `fetch` 版本；客户端调用可省略。

## 7. Markdown/TOC 规范
- Markdown 渲染入口：`src/lib/shared/markdown`（`createMarkdownIt` / `renderMarkdown`）。
- 结构化 TOC：应由后端/构建期生成（SEO 友好）。
- 前端高亮：IntersectionObserver 只负责 active 交互（封装为 action）。
- Markdown 组件增强：使用 `src/lib/shared/markdown/components` 注册/挂载。

## 8. 性能策略（静态优先 + 渐进增强）
- 文章/列表/友链/归档：SSR/SSG 输出完整 HTML。
- 评论/点赞/TOC 高亮/图片缩放等：client-only islands。
- 重库只在 `onMount`/交互后动态 import。

## 9. SEO 与友链页面
- 友链列表是强 SEO 页：必须 SSR/SSG 输出完整列表。
- 客户端只做增强（shuffle、筛选、申请表单、动画）。

## 10. WebAuthn / Passkey
- `navigator.credentials.*` 仅可在客户端触发（onMount/点击）。
- 服务端只负责 challenge 签发与验证。

## 11. 设计系统与样式约束
- 全局主题/Token：`src/routes/layout.css`（Tailwind v4 + @theme）。
- 组件样式优先在 `<style>` 中使用 `@apply`，并用 `@reference` 引入 `layout.css`。
- 视觉风格保持“温暖灰 + jade 主色 + serif/sans/mono”体系，不引入新的风格系统。

## 12. 禁止事项（硬性）
- 不在页面组件里直接写 observer/event 绑定。
- 不引入散装全局状态，必须收敛到 store/context。
- 不在 SSR 阶段访问 `window/document/navigator`。
- 不在首屏引入重库；能动态加载就动态加载。
- 不改 `build/`、`node_modules/`、`design-system/`（除非明确要求）。

## 13. 可选落地路径建议（如需新增模块）
- DOM actions: `src/lib/shared/actions/*.ts`
- DOM helpers: `src/lib/shared/dom/*.ts`
- 业务模块：`src/lib/features/<feature>/{api.ts,types.ts,components/*}`
- UI 组件：`src/lib/ui/*.svelte`

---

如需扩展本规范，请先解释原因与改动范围，再调整本文件。
