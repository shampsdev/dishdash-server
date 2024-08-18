include .env

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

db-compose-up: # Launch database+adminer from docker-compose
	docker compose up database adminer --build -d && docker compose logs -f

compose-up: # Launch full service in docker-compose
	docker compose up --build -d && docker compose logs -f

compose-down: # Compose down
	docker compose down

swag: # Generate swag documentation (github.com/swaggo/swag)
	swag init -g cmd/server/main.go -o docs

run: # Run server
	go run cmd/server/main.go

lint: # Lint (github.com/golangci/golangci-lint)
	golangci-lint run

lint-fix: # Lint fix
	golangci-lint run --fix

test-e2e: # Test e2e
	go test -v -race e2e/e2e_test.go

update-golden-e2e: # Update golden
	go test -v -race e2e/e2e_test.go -update-golden