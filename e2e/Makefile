.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ E2E Environment

E2E_ENV_FILE := e2e.env

E2E_DB_HOST ?= localhost
E2E_DB_PORT ?= 5433
E2E_DB_USER ?= root
E2E_DB_PASSWORD ?= root
E2E_DB_NAME ?= root

E2E_DB_URL="postgres://$(E2E_DB_USER):$(E2E_DB_PASSWORD)@$(E2E_DB_HOST):$(E2E_DB_PORT)/$(E2E_DB_NAME)?sslmode=disable"

e2e: e2e-create-env e2e-test e2e-delete-env ## Run all

e2e-create-env: e2e-export-vars e2e-create-db e2e-create-backend ## Create environment

e2e-export-vars:
	@echo "Creating $(E2E_ENV_FILE)"
	@echo "POSTGRES_HOST=$(E2E_DB_HOST)" > $(E2E_ENV_FILE)
	@echo "POSTGRES_PORT=$(E2E_DB_PORT)" >> $(E2E_ENV_FILE)
	@echo "POSTGRES_USER=$(E2E_DB_USER)" >> $(E2E_ENV_FILE)
	@echo "POSTGRES_PASSWORD=$(E2E_DB_PASSWORD)" >> $(E2E_ENV_FILE)
	@echo "POSTGRES_DB=$(E2E_DB_NAME)" >> $(E2E_ENV_FILE)
	@echo "HTTP_PORT=8001" >> $(E2E_ENV_FILE)
	@echo "DEBUG=True" >> $(E2E_ENV_FILE)

e2e-create-db:
	@echo "Launching e2e database"
	@docker compose -f docker-compose.yaml up e2e-database e2e-adminer -d
	@echo "You can manage database in http://localhost:1001/?pgsql=e2e-database&username=$(E2E_DB_USER)&database=$(E2E_DB_NAME)" with password: $(E2E_DB_PASSWORD)

e2e-create-backend:
	@echo "Launching server"
	@docker compose -f docker-compose.yaml up e2e-backend --build -d
	@echo "Backend is running on 8001 port. Visit http://localhost:8001/api/v1/swagger/index.html"


e2e-delete-env: e2e-delete-backend e2e-delete-db ## Delete environment

e2e-delete-db:
	@echo "Deleting e2e database"
	@-docker compose -f docker-compose.yaml down e2e-database e2e-adminer

e2e-delete-backend:
	@echo "Deleting e2e backend"
	@-docker compose -f docker-compose.yaml down e2e-backend

##@ E2E Tests
E2E_TESTS ?= 1-lobby-join 2-lobby-swipe 3-lobby-vote

e2e-test: ## Launch all test
	@for testdir in $(E2E_TESTS); do \
		$(MAKE) -C "$$testdir" test; \
	done
