MIGRATIONS_PATH = ./internal/database/migrations
DB_ADDR = postgres://admin:password@localhost:5432/glottr?sslmode=disable

swagger:
	@swag init -g ./cmd/api/main.go && swag fmt

run:
	@go run ./cmd/api

docker:
	@docker compose up -d

migration:
	@goose -v -s -dir $(MIGRATIONS_PATH) create $(NAME) sql

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_ADDR) goose -dir $(MIGRATIONS_PATH) up

down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_ADDR) goose -dir $(MIGRATIONS_PATH) down

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_ADDR) goose -dir $(MIGRATIONS_PATH) reset

seed:
	@go run ./cmd/seed/main.go
