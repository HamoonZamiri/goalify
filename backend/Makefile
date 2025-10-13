# Database connection - adjust as needed
# Load environment variables from .env file
export
-include .env
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR ?= ./internal/db/migrations

# Build commands
build:
	@go build -o bin/goalify ./cmd/main.go

run: build
	@./bin/goalify

# Development commands
dev: up
	@echo "Waiting for Docker to start..."
	@sleep 2
	@go tool github.com/air-verse/air

jqdev: up
	@echo "Waiting for Docker to start..."
	@sleep 2
	@echo "Starting air server with pretty json printing"
	@go tool github.com/air-verse/air --build.cmd "go build -o bin/goalify ./cmd/main.go" --build.bin "./bin/goalify"  | jq -C -R 'fromjson?'

# Migration commands
migrate-up:
	@echo "Running migrations up..."
	@go tool github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	@echo "Rolling back one migration..."
	@go tool github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	@echo "Checking migration status..."
	@go tool github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@echo "Creating migration: $(NAME)"
	@go tool github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Code generation (if using sqlc)
generate:
	@echo "Generating code..."
	@go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Linting
lint:
	@echo "Running linter..."
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run

lint-fix:
	@echo "Running linter with fixes..."
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run --fix

# Testing commands
unit:
	@echo "Running unit tests..."
	@go test -count=1 ./pkg/... ./internal/db/... ./internal/middleware/

store:
	@echo "Running store/repository tests..."
	@go test -count=1 -v ./internal/users/stores/...
	@go test -count=1 -v ./internal/goals/stores/...
	@go test -count=1 -v ./internal/loot/stores/...

inte:
	@echo "Running integration tests..."
	@go test -count=1 -v ./tests/...

test: unit store inte

testall:
	@go test -count=1 -v ./...

# Docker commands
up:
	@docker compose up -d

down: 
	@docker compose down

start: up
	@echo "Waiting for Docker to start..."
	@sleep 2
	@make run

# Clean up
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

list-tools:
	@echo "Listing installed Go tools"
	@go tool list


# Help command
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Start development server with hot reload"
	@echo "  jqdev         - Start development server with pretty JSON logs"
	@echo "  migrate-up    - Run all pending migrations"
	@echo "  migrate-down  - Rollback one migration"
	@echo "  migrate-status- Check migration status"
	@echo "  migrate-create NAME=name - Create new migration"
	@echo "  generate      - Generate code (sqlc)"
	@echo "  lint          - Run linter"
	@echo "  lint-fix      - Run linter with auto-fixes"
	@echo "  test          - Run all tests (unit + store + integration)"
	@echo "  unit          - Run unit tests only"
	@echo "  store         - Run store/repository tests only"
	@echo "  inte          - Run integration tests only (tests/ directory)"
	@echo "  testall       - Run all tests including edge cases"
	@echo "  up            - Start Docker services"
	@echo "  down          - Stop Docker services"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help"
	@echo "  list-tools    - List installed go tools

.PHONY: build run dev jqdev migrate-up migrate-down migrate-status migrate-create generate lint lint-fix unit store inte test testall up down start clean help
