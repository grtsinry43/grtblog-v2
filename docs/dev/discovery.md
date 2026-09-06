# 搜索与 AI 阅读入口

本阶段按请求生成，不生成静态产物，也不订阅 ISR 内容事件。`/sitemap/` 已绕过 Nginx 静态文件优先规则和前台自动回填。

## 入口

| 地址 | 内容 |
| --- | --- |
| `/sitemap/` | 复用博客主题的导航页面，支持标题与摘要筛选 |
| `/sitemap.xml` | XML 地图，超过 10,000 条时返回分片索引 |
| `/sitemaps/1.xml` | XML 分片 |
| `/robots.txt` | 后台配置的爬取策略，加自动生成的 Sitemap 声明 |
| `/llms.txt` | 博客简介、作者背景、阅读说明、精选内容和目录入口 |
| `/posts/llms.txt`、`/moments/llms.txt`、`/pages/llms.txt` | 每页 100 条的 Markdown 阅读目录，使用 `?page=2` 翻页 |
| `原网页路径/index.md` | 公开文章、手记和独立页面的 Markdown 正文 |

HTML head 使用 `rel="describedby"` 指向阅读指南，内容详情页使用 `rel="alternate" type="text/markdown"` 指向正文。文本正文标注原网页地址，并通过 HTTP Link 声明 canonical；文本版本返回 `X-Robots-Tag: noindex`。

## 后台与部署

运行 `0070_add_discovery_settings.sql` 迁移后，在“设置 → 站点信息 → 搜索与 AI 阅读”编辑：博客阅读简介、作者背景、阅读与引用说明、精选阅读路径、robots.txt 爬取策略。配置复用现有系统配置表单和保存 API。简介留空使用网站描述，robots 策略默认及留空时禁止抓取 `/admin`，其余路径允许抓取。后台 HTML 和 Nginx 响应另外固定声明 `noindex, nofollow`，不随 robots 配置改变。

必须配置有效的 `site.public_url`。自动链接使用该域名，不使用请求 Host。robots 策略可填写多组 User-agent / Allow / Disallow；本站 Sitemap 地址由系统自动附加。爬取规则不改变公开内容的访问权限，也不改变 sitemap 的收录集合。

部署时需要同步更新 Go、前台和 `deploy/nginx/nginx.conf`。旧静态 robots 文件已替换为动态路由。现有文章 HTML 中新增的 head 与页脚链接要等页面重新渲染后才会出现；本次不触发全站 ISR 重建。

## 内容与缓存

- 列表查询只读取元数据，正文入口单独读取内容；草稿、软删除和禁用页面不会输出。
- 手记路径按站点时区生成；错误日期返回 404。内置页面不输出数据库中的占位正文。
- 地图包含首页及主要导航、文章、手记、独立页面、分类、专栏、公开相册与照片。重复分页和标签弹窗不单独收录。
- 内容 lastmod 使用 `ContentUpdatedAt`，不使用浏览量或渲染时间；无可靠内容修改时间的导航、分类、相册等省略 lastmod。
- 精选阅读只从当前公开目录匹配路径，撤稿后自动消失。
- XML、robots 和 Markdown 使用 `Cache-Control: no-cache` 与 ETag，先读取当前状态再判断 304；HTML 导航使用 no-store。
- 数据或站点配置不可用时返回 503，避免生成空的成功地图。robots 只依赖配置，不读取内容目录。
- Markdown 保留代码，转换展示容器为文字或链接，并解析相对图片与链接；相册与思考目前通过原网页阅读。
- year-card 支持当前的 URL 属性和旧版 `{slug=...}` 写法：有 URL 时输出 Markdown 链接，只有 slug 时保留年度总结引用说明，不推测地址。联邦 mention 复用后端识别规则，输出普通 `@用户@实例` 文本；代码示例和 `<!--more-->` 保持原样。

## 验证

```sh
cd server
go test ./internal/app/discovery
DISCOVERY_TEST_DSN='host=... dbname=... user=... sslmode=disable' go test ./internal/http/handler -run TestDiscoveryPostgresLifecycle -count=1
```

集成测试须指定临时 PostgreSQL，自动创建和清理独立 schema，覆盖配置迁移与保存、公开过滤、撤稿与删除、相册可见性、日期路径、条件请求和数据库故障。

后续接入 ISR 时需明确所有文本与目录的发布、撤稿、删除、slug 和配置变更失效规则，再考虑静态产物与离线降级。
