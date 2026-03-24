# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Up app in docker
up:
	docker compose --profile dev up -d

# Logs prod app in docker
logs:
	docker logs go-with-tools-air-1 -f

# Down app in docker
down:
	docker compose --profile dev down && docker compose --profile prod down

# Up prod app in docker
prod-up:
	docker compose --profile prod up --build -d

# Logs prod app in docker
prod-logs:
	docker logs go-with-tools-web-1 -f

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

DB_USER ?= go_with_tools_user
DB_NAME ?= go_with_tools_db
SCHEMA  := internal/database/schema.sql

db-schema:
	pg_dump -U $(DB_USER) -d $(DB_NAME) --schema-only -f $(SCHEMA)
	sed -i '/^\\restrict\b/d; /^\\unrestrict\b/d' $(SCHEMA)

sqlc: db-schema
	sqlc generate

swag:
	swag init -g main.go -d cmd/api,internal/ -o cmd/api/docs

migrate-in-docker:
	docker compose run --rm migrate

.PHONY: all build run test clean watch up logs down prod-up prod-logs prod-down itest db-schema sqlc swag migrate-in-docker
