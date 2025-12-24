# Makefile

# Auto-load .env if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Variables (fallback defaults)
API_DB_DSN ?=

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  help              Print this help message"
	@echo "  run               Run the cmd/api application"
	@echo "  migrations-new    Create a new database migration (NAME=...)"
	@echo "  migrations-down   Roll back the last database migration"
	@echo "  migrations-up     Apply all up database migrations"
	@echo "  tidy              Format Go files and tidy dependencies"

# Internal confirmation target
.PHONY: confirm
confirm:
	@printf "Are you sure? [y/N] " && read ans && [ "$${ans:-N}" = "y" ]

# Run application
.PHONY: run
run:
	go run ./cmd/api

# Create new migration
.PHONY: migrations-new
migrations-new:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required"; \
		echo "Usage: make migrations-new NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration file for $(NAME)..."
	migrate create -ext sql -dir ./internal/database/migrations -seq $(NAME)

# Roll back last migration
.PHONY: migrations-down
migrations-down: confirm
	@echo "Rolling back last migration..."
	migrate -database "$(API_DB_DSN)" -path ./internal/database/migrations down 1

# Apply migrations
.PHONY: migrations-up
migrations-up: confirm
	@echo "Running up migrations..."
	migrate -database "$(API_DB_DSN)" -path ./internal/database/migrations up

# Go maintenance
.PHONY: tidy
tidy:
	@echo "Formatting .go files..."
	go fmt ./...
	@echo "Tidying module dependencies..."
	go mod tidy
	@echo "Verifying dependencies..."
	go mod verify
