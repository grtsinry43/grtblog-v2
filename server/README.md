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
