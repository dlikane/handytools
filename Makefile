.DEFAULT_GOAL := help

# Variables
BINARY_NAME := img
CMD_PATH := ./cmd/img

.PHONY: help
help: ## Display this help message
	@echo Available commands:
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the img binary
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(CMD_PATH)

.PHONY: run
run: build ## Run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

.PHONY: deploy
deploy: ## Install img globally
	@echo "Installing $(BINARY_NAME) globally..."
	go install $(CMD_PATH)
	@echo "Deployment complete! Run '$(BINARY_NAME)' from anywhere."

.PHONY: clean
clean: ## Remove built binary
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: lint
lint: ## Run linter (golangci-lint)
	@echo "Running linter..."
	golangci-lint run

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test ./...

.PHONY: install
install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod tidy

.PHONY: all
all: fmt lint test build ## Run format, lint, test, and build
