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
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been cleared"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi


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
