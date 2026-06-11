APP_NAME := funkey-bite-api
LAMBDA_NAME := bootstrap
CMD_DIR := ./cmd/api
BIN_DIR := ./bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
LAMBDA_PATH := $(BIN_DIR)/$(LAMBDA_NAME)
ENV_FILE := .env
ENV_EXAMPLE := .env.example
GO ?= go

.PHONY: help env bootstrap deps tidy fmt vet test test-race contract check run dev build build-local build-lambda start clean

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

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