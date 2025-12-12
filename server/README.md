# grtblog server

Go Fiber + GORM API scaffold that powers grtblog. The layout follows idiomatic Go conventions by splitting the executable (`cmd/api`) from reusable packages (`internal/...`).

## Structure

- `cmd/api`: program entry point that loads configuration, initializes dependencies, and starts Fiber.
- `internal/config`: environment-driven configuration helpers.
- `internal/database`: database (GORM) initialization.
- `internal/http`: HTTP handlers and routers.
- `internal/server`: Fiber server wiring.
- `migrations`: SQL migration files (compatible with [Goose](https://github.com/pressly/goose)).
- `storage`: default SQLite location (safe to replace with any driver supported by GORM).

## Getting started

```bash
cd server
go mod tidy             # download Fiber + GORM dependencies
APP_PORT=8080 go run ./cmd/api
```

Environment variables such as `APP_PORT`, `APP_ENV`, `DB_DRIVER`, `DB_DSN`, and `DB_AUTO_MIGRATE` can be set to customize behavior. Default database driver is SQLite for quick local usage. Switch to Postgres by setting `DB_DRIVER=postgres` and providing a valid `DB_DSN`.

## Database migrations (Goose)

Migrations live under `migrations/` and follow the `NNNN_description.sql` naming convention so that they work seamlessly with [Goose](https://github.com/pressly/goose).

### Install Goose CLI

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Useful commands

```bash
# Apply all up migrations (uses the DSN from your environment or override inline)
DB_DSN=postgres://postgres:postgres@localhost:5432/grtblog?sslmode=disable \
  make migrate-up

# Roll the last migration back
DB_DSN=postgres://postgres:postgres@localhost:5432/grtblog?sslmode=disable \
  make migrate-down

# Create a new timestamped SQL migration under migrations/
make migrate-create NAME=add_posts_table
```

The Make targets are thin wrappers around `goose` and rely on it being installed in your `PATH`. Feel free to call the Goose binary directly if you prefer more control.

## API documentation (Swagger + Scalar)

Swagger annotations (via [swaggo/swag](https://github.com/swaggo/swag)) describe the HTTP handlers and generate an OpenAPI schema that fuels the Scalar UI.

1. 安装 `swag` CLI（需要网络环境）：
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
2. 在 `server/` 目录下运行：
   ```bash
   make docs
   ```
   该命令会执行 `swag init -g cmd/api/main.go -o docs --outputTypes json --parseInternal`，输出 `docs/swagger.json`。
3. 启动 API 服务后访问 [http://localhost:8080/docs](http://localhost:8080/docs) 查看 Scalar 渲染的交互式文档，`/docs/openapi.json` 则返回原始 OpenAPI JSON。

若你还未安装 `swag`，仓库自带的 `docs/swagger.json` 仍可直接用于开发，后续更新接口时重新执行 `make docs` 即可刷新文档。

## Authentication & RBAC

- 通过 `AUTH_SECRET`、`AUTH_ISSUER`、`AUTH_ACCESS_TTL`、`AUTH_DEFAULT_ROLES` 配置 JWT 及默认角色。
- 当前 API 前缀为 `/api/v2`。注册：`POST /api/v2/auth/register`（示例实现 SHA 加密 + 默认角色绑定）；登录：`POST /api/v2/auth/login`，签发带角色/权限的 JWT。
- OAuth/OIDC：`auth.Service` 暴露 `RegisterProvider`、`LoginWithProvider` 预留扩展点，可在后续接入外部身份提供方。
- 路由保护：组合 `middleware.RequireAuth`、`RequirePermission`（底层使用 Casbin），即可保护敏感接口。示例中 `/api/v2/website-info` GET 需要 `config:read`，写操作需要 `config:write`，而 `/api/v2/public/website-info` 为无需鉴权的公开读取接口。
- RBAC 模型在 `configs/rbac_model.conf`，Casbin 会从 `role_permission` 数据表加载策略；如需自动刷新请设置 `RBAC_AUTO_RELOAD=true`。
- 风控：登录/注册接口已加 IP 级限流（Fiber `limiter`），并支持 Cloudflare Turnstile 人机校验。Turnstile 开关/凭据可通过数据库表 `sys_config` 动态调整：`turnstile.enabled` (bool)、`turnstile.secret`、`turnstile.siteKey`、`turnstile.verifyURL`、`turnstile.timeoutSeconds`，未配置时回落到环境变量 `TURNSTILE_*` 默认值。
