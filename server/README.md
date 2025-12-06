# grtblog server

Go Fiber + GORM API scaffold that powers grtblog. The layout follows idiomatic Go conventions by splitting the executable (`cmd/api`) from reusable packages (`internal/...`).

## Structure

- `cmd/api`: program entry point that loads configuration, initializes dependencies, and starts Fiber.
- `internal/config`: environment-driven configuration helpers.
- `internal/database`: database (GORM) initialization.
- `internal/http`: HTTP handlers and routers.
- `internal/server`: Fiber server wiring.
- `storage`: default SQLite location (safe to replace with any driver supported by GORM).

## Getting started

```bash
cd server
go mod tidy             # download Fiber + GORM dependencies
APP_PORT=8080 go run ./cmd/api
```

Environment variables such as `APP_PORT`, `APP_ENV`, `DB_DRIVER`, `DB_DSN`, and `DB_AUTO_MIGRATE` can be set to customize behavior. Default database driver is SQLite for quick local usage. Switch to Postgres by setting `DB_DRIVER=postgres` and providing a valid `DB_DSN`.
