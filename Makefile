-include .env
export

# ============
# Enviroments
# ============

env-up:
	@docker compose up -d shortener-postgres

env-down:
	@docker compose down shortener-postgres

env-cleanup:
	@read -p "Clear all volume files in the environment? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down shortener-postgres && \
		rm -rf ./out/pgdata && \
		echo "Environment files have been cleared"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi


logs-cleanup:
	@read -p "Clear all log files? Risk of losing logs. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ./out/logs && \
		echo "Logs files have been cleared"; \
	else \
		echo "Logs files cleanup cancelled"; \
	fi

# ============
# Migrations
# ============

GOOSE_DRIVER ?= postgres
GOOSE_MIGRATION_DIR ?= ./migrations
GOOSE_DBSTRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

goose = GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=add_clicks"; exit 1; \
	fi
	@$(goose) -s create $(name) sql

migrate-up:
	@$(goose) up

migrate-down:
	@$(goose) down

migrate-status:
	@$(goose) status

# ============
# Startap app
# ============
	
shortener-run:
	@go mod tidy && \
	go run ./cmd/server/main.go

# ============
# Deploy
# ============

shortener-deploy:
	@docker compose up -d --build shortener

shortener-undeploy:
	@docker compose down shortener

# ============
# Swagger
# ============
