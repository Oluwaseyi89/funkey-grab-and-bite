APP_NAME := funkey-bite-api
LAMBDA_NAME := bootstrap
CMD_DIR := ./cmd/api
BIN_DIR := ./bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
LAMBDA_PATH := $(BIN_DIR)/$(LAMBDA_NAME)
ENV_FILE := .env
ENV_EXAMPLE := .env.example
GO ?= go
MIGRATE ?= migrate
MIGRATIONS_DIR := ./migrations

# Load DB_* settings from .env (if present) so migrate targets hit the same database as the API.
-include $(ENV_FILE)
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help env bootstrap deps tidy fmt vet test test-race contract check run dev build build-local build-lambda start clean migrate-up migrate-down migrate-version migrate-force migrate-create

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## Create .env from .env.example if it does not exist
	@if [ -f $(ENV_FILE) ]; then \
		echo "$(ENV_FILE) already exists"; \
	else \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
		echo "Created $(ENV_FILE) from $(ENV_EXAMPLE)"; \
	fi

bootstrap: env deps ## Prepare local development prerequisites

deps: ## Download Go module dependencies
	$(GO) mod download

tidy: ## Tidy Go module files
	$(GO) mod tidy

fmt: ## Format Go source files
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

test: ## Run the test suite
	$(GO) test ./...

test-race: ## Run tests with the race detector
	$(GO) test -race ./...

contract: ## Validate the canonical OpenAPI contract against admin and web consumers
	node ../scripts/check-openapi-contract.mjs all

check: fmt vet test contract ## Run formatting, vetting, tests, and contract validation

run: build-local ## Run the local API binary directly
	$(BIN_PATH)

dev: run ## Alias for local development run

build-local: ## Build the standard API binary for local/container development
	@mkdir -p $(BIN_DIR)
	@echo "Building local development binary..."
	$(GO) build -o $(BIN_PATH) $(CMD_DIR)/main_local.go
	@echo "Local binary built at $(BIN_PATH)"

build-lambda: ## Cross-compile optimized bootstrap binary for AWS Lambda (Linux ARM64)
	@mkdir -p $(BIN_DIR)
	@echo "Cross-compiling for AWS Lambda (Linux/ARM64)..."
	GOOS=linux GOARCH=arm64 $(GO) build -tags lambda -o $(LAMBDA_PATH) $(CMD_DIR)/main_lambda.go
	@echo "Lambda bootstrap binary built at $(LAMBDA_PATH)"

build: build-local build-lambda ## Build both local and Lambda targets simultaneously

start: build-local ## Build and start the compiled local API binary
	$(BIN_PATH)

clean: ## Remove built artifacts
	rm -rf $(BIN_DIR)

migrate-up: ## Apply all pending migrations (the API also does this on startup)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Roll back the most recent migration (N=<count> to roll back more)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down $(or $(N),1)

migrate-version: ## Show the currently applied migration version
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-force: ## Mark version V as applied and clean after fixing a failed migration (V=<version>)
	@test -n "$(V)" || (echo "usage: make migrate-force V=<version>" && exit 1)
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(V)

migrate-create: ## Create a new up/down migration pair (NAME=<snake_case_description>)
	@test -n "$(NAME)" || (echo "usage: make migrate-create NAME=<snake_case_description>" && exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq -digits 6 $(NAME)
