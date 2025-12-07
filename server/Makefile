GOOSE ?= goose
GOOSE_DRIVER ?= postgres
MIGRATIONS_DIR ?= migrations
DB_DSN ?= postgres://postgres:postgres@localhost:5432/grtblog?sslmode=disable

.PHONY: migrate-up migrate-down migrate-status migrate-version migrate-create

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_DSN)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_DSN)" status

migrate-version:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_DSN)" version

migrate-create:
ifndef NAME
	$(error NAME is required, e.g. make migrate-create NAME=add_posts_table)
endif
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(NAME) sql
