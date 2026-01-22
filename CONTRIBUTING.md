# 👋 欢迎加入 grtblog-v2 开发！

嘿，欢迎来到 **grtblog-v2**！很高兴能和你一起开发这个项目。🎉

这个项目不仅仅是一个博客，更是我们对 **高性能 Web 架构** 的一次探索。我们用 Go 做了后端，Vue 做了后台，SvelteKit 做了前台，还自己搓了一套类似 Next.js ISR 的静态生成机制。

为了让你少踩坑，玩得开心，这里有一份“生存指南”，包含了启动项目、调试 ISR 和代码规范的黑话。

---

## 🛠 我们用的家伙什 (技术栈)

这是一个 Monorepo（虽然没强制 workspace），主要分成三块地盘：

* **Server (`/server`)**: **Go (1.24+)** + Fiber + GORM。
* *特点*：DDD 架构，负责数据管理和驱动静态页面生成。


* **Admin (`/admin`)**: **Vue 3** + Naive UI。
* *特点*：给作者（也就是我俩）用的，极致的 Markdown 写作体验。


* **Web (`/web`)**: **SvelteKit** + Tailwind CSS。
* *特点*：给读者看的。平时是 SSR 模式，但在生成静态快照时会被 Go 后端爬取。



---

## ⚡️ 快速启动 (Quick Start)

本地开发主要跑这三个服务。

### 1. 启动后端 (Server)

所有数据的源头，跑在 `:8080`。

```bash
cd server
# 第一次跑记得搞定 .env 和数据库，具体直接复制重命名 .env.example ，然后 make migrate-up
make run

```

### 2. 启动后台 (Admin)

写文章的地方，跑在 `:5799`。

```bash
cd admin
pnpm i && pnpm dev

```

### 3. 启动前台 (Web)

看效果的地方。

* **日常开发**：跑在 `:5173`（支持热更新，改完代码即刻生效）。
* **SSR 生产模拟**：跑在 `:3000`（用于测试 ISR 生成流程）。

```bash
cd web
pnpm i && pnpm dev  # 启动开发服务器 :5173

```

---

## 🔍 核心玩法：调试 ISR (静态生成)

这是本项目最硬核的地方。我们不直接用 SvelteKit 的 Adapter-Static，而是让 **Go 后端去爬 `:3000` (SSR 生产端口) 的页面**，然后存成静态 HTML。

如果你改了 **SvelteKit 的路由逻辑** 或者 **后端的生成代码**，一定要测一下这个流程！

我写好了一个自动化脚本，一键模拟生产环境的“构建 -> 生成 -> 预览”全过程：

```bash
# 确保后端 (:8080) 已经运行着，然后由根目录执行：
make preview-isr

```

**这个命令会自动做这些事：**

1. 🏗️ **Build**: 编译 Web 前端。
2. 🔌 **Serve**: 在后台悄悄启动 SSR 服务（端口 **:3000**）。
3. 🔄 **Trigger**: 调用后端 API 抓取 `:3000` 的页面（生成静态 HTML 到 `server/storage/html`）。
4. 🌍 **Preview**: 启动一个静态文件服务器（端口 **:5555**）并自动打开浏览器。

> **看到 `:5555` 里的页面正常，才说明 ISR 机制没挂！**

---

## 📐 避坑指南 (Guidelines)

为了防止代码打架，这里有几条铁律：

### 1. 后端：同步 vs 异步

我们在后端有两个很像的操作，别搞混了：

* **API 接口 (`RefreshPostsHTML` Handler)**：必须是 **同步 (Sync)** 的。脚本调用它时，要一直等到生成完才返回 200。这样脚本才知道什么时候去关掉 SSR 进程。
* **事件监听 (`ArticleUpdated` Subscriber)**：必须是 **异步 (Async)** 的。保存文章时，通过 EventBus 触发生成，千万别阻塞保存接口，不然前端会卡死。

### 2. 前端：URL 结尾斜杠

SvelteKit 里配置了 `trailingSlash: 'always'`。

* **一定要** 保证生成的 URL 结尾带 `/`（例如 `/posts/hello/`）。
* 这样才能配合文件系统里的 `posts/hello/index.html` 结构，否则静态部署后刷新页面会 404。

### 3. Markdown 组件

我们在 `shared/markdown/components.ts` 里定义了 Markdown 组件（比如相册、卡片）。

* 这块代码是 **Admin 和 Web 共用** 的。
* 如果你要加新组件，记得两边都要适配，别只改了一头，另一头渲染不出来。

---

## 📝 提交代码 (Commits)

咱们还是走标准那一套，看着整齐：

* `feat`: 加新功能
* `fix`: 修 Bug
* `docs`: 改文档
* `refactor`: 重构（没改功能）
* `chore`: 杂活（构建、依赖等）

例子：`feat(server): 增加文章字数统计功能`

---

💡 **有问题随时叫我！Happy Coding!** 🚀