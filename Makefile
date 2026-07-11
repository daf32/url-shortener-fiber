include .env
export

export PROJECT_ROOT=$(shell pwd)

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
		docker compose down shortener-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been cleared"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

logs-cleanup:
	@read -p "Clear all log files? Risk of losing logs. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Logs files have been cleared"; \
	else \
		echo "Logs files cleanup cancelled"; \
	fi

# ============
# Migrations
# ============

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "The required seq parameter is missing. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "The required action parameter is missing. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

# ============
# Startap app
# ============
	
shortener-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/server/main.go

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

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/server/main.go \
		-o docs \
		--parseInternal \
		--parseDependency
