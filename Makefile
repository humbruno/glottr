swagger:
	@swag init -g ./cmd/api/main.go && swag fmt

docker:
	@docker compose up -d
