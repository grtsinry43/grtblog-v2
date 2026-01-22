# Repository Guidelines

## Project Structure & Module Organization
- `cmd/api`: application entry point (config load, dependency wiring, Fiber startup).
- `internal/`: core packages (config, database, HTTP handlers/routers, services, domain models, persistence).
- `configs/`: runtime configuration files (e.g., RBAC model).
- `migrations/`: Goose SQL migrations using `NNNN_description.sql` naming.
- `docs/`: generated OpenAPI artifacts (`swagger.json`).
- `storage/`: default SQLite database location for local use.

## Build, Test, and Development Commands
- `go mod tidy`: sync Go module dependencies.
- `APP_PORT=8080 go run ./cmd/api`: run the API locally.
- `make migrate-up|migrate-down|migrate-status|migrate-version`: manage database migrations via Goose.
- `make migrate-create NAME=add_posts_table`: create a new SQL migration.
- `make docs`: regenerate OpenAPI JSON from Swagger annotations.
  - Note: `swag` is sensitive to annotation order; keep `@BasePath /api/v2` in the main comment block (prefer near the end) to ensure it is emitted into `docs/swagger.json`.

## Coding Style & Naming Conventions
- Use standard Go formatting (`gofmt`) and idiomatic Go naming (PascalCase exports, camelCase locals).
- Keep packages cohesive and aligned with the existing layout (e.g., `internal/app/*` for services, `internal/domain/*` for entities/repositories).
- Migrations must follow `NNNN_description.sql` so Goose can order them.

## Additional Agent Requirements
- 0. Guarantee high-quality code; forbid placeholder implementations, fake implementations, or problematic code.
- 1. For major changes, list a modification plan first, including strategy, impacted files, technical approach, and choices with reasons.
- 2. Act as an excellent Go engineer; adhere to Go and Fiber best practices, handle Go legacy issues with modern syntax, and avoid GORM pitfalls using experienced practices.
- 3. Before using any library, check its latest version and documentation; leverage search and network access; do not guess or assume usage.

## Testing Guidelines
- No dedicated test suite is present in this repo yet. When adding tests, place `_test.go` files alongside the code under `internal/` and run `go test ./...`.
- Prefer table-driven tests for handler/service logic where possible.

## Commit & Pull Request Guidelines
- Follow the observed Conventional Commits style (e.g., `feat: ...`, `feat(server): ...`, `fix: ...`).
- PRs should include a clear summary, rationale, and any required config or migration notes.
- If you change API handlers or models, update `docs/swagger.json` via `make docs` and mention it in the PR.

## Security & Configuration Tips
- Runtime behavior is controlled via env vars like `APP_PORT`, `DB_DRIVER`, `DB_DSN`, `AUTH_SECRET`, and `AUTH_DEFAULT_ROLES`.
