APP_NAME := funkey-bite-api
CMD_DIR := ./cmd/api
BIN_DIR := ./bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)
ENV_FILE := .env
ENV_EXAMPLE := .env.example
GO ?= go

.PHONY: help env bootstrap deps tidy fmt vet test test-race check run dev build start clean

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

check: fmt vet test ## Run formatting, vetting, and tests

run: ## Run the API with go run
	$(GO) run $(CMD_DIR)

dev: run ## Alias for local development run

build: ## Build the API binary into ./bin
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_PATH) $(CMD_DIR)

start: build ## Build and start the compiled API binary
	$(BIN_PATH)

clean: ## Remove built artifacts
	rm -rf $(BIN_DIR)
