include .env

db-compose-up:
	docker compose up database adminer --build -d && docker compose logs -f

compose-up:
	docker compose up --build -d && docker compose logs -f

compose-down:
	docker compose down

swag:
	swag init -g cmd/server/main.go -o docs

run:
	go run cmd/server/main.go

lint:
	golangci-lin run

lint-fix:
	golangci-lin run --fix