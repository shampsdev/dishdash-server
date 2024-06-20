include .env

dev-compose-up:
	docker compose up database adminer --build -d && docker compose logs -f

compose-down:
	docker compose down

swag:
	swag init -g cmd/server/main.go -o docs

run:
	go run cmd/server/main.go