ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Run
run: ## Run server
	go run cmd/server/main.go

compose-up: ## Launch full service in docker-compose
	docker compose up --build -d && docker compose logs -f

compose-down: ## Compose down
	docker compose down

##@ Database
db-compose-up: ## Launch database+adminer from docker-compose
	docker compose up database adminer --build -d && docker compose logs -f

DB_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

db-migrate-up: migrate-install ## Migrate up
	$(MIGRATE) -database $(DB_URL) -path migrations up

db-migrate-down: migrate-install ## Migrate down
	$(MIGRATE) -database $(DB_URL) -path migrations down

db-default-data: ## Add necessary default data to db (postgresql-client-16 (psql) needed)
	psql $(DB_URL) -a -f $(PROJECT_DIR)/migrations/data/default.sql

db-reset: ## Reset database (down + up + default-data)
	make db-migrate-down
	make db-migrate-up
	make db-default-data


##@ Tools
GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
golangci-lint-install:
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3)

SWAG = $(shell pwd)/bin/swag
swag-install:
	$(call go-get-tool,$(SWAG),github.com/swaggo/swag/cmd/swag@v1.16.3)

MIGRATE = $(shell pwd)/bin/migrate
migrate-install:
	$(call go-get-tool,$(MIGRATE),-tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1)

lint: golangci-lint-install ## Lint (github.com/golangci/golangci-lint)
	$(GOLANGCI_LINT) run

lint-fix: golangci-lint-install ## Lint fix
	$(GOLANGCI_LINT) run --fix

swag: swag-install ## Generate swag documentation (github.com/swaggo/swag)
	$(SWAG) init -g cmd/server/main.go -o docs

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
