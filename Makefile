.DEFAULT_GOAL := help

# Variables
BINARY_NAME := img
CMD_PATH := ./cmd/img

#
# Windows batch rename setup
#
DEPLOY_DIR := C:\Users\User\go\bin
REG_FILE := internal/batchrename/batchrename.reg
SRC_BAT := internal/batchrename/batchrename.bat
SRC_VBS := internal/batchrename/batchrename.vbs

.PHONY: deploy-win
deploy-win: deploy ## Take care of Windows deployment
	@if not exist "$(DEPLOY_DIR)" mkdir "$(DEPLOY_DIR)"
	@powershell -Command "if (Test-Path '$(SRC_BAT)') { Copy-Item -Force '$(SRC_BAT)' '$(DEPLOY_DIR)\batchrename.bat' } else { Write-Host 'Warning: $(SRC_BAT) not found' }"
	@powershell -Command "if (Test-Path '$(SRC_VBS)') { Copy-Item -Force '$(SRC_VBS)' '$(DEPLOY_DIR)\batchrename.vbs' } else { Write-Host 'Warning: $(SRC_VBS) not found' }"
	@powershell -Command "if (Test-Path '$(REG_FILE)') { Start-Process reg -ArgumentList 'import $(REG_FILE)' -NoNewWindow -Wait } else { Write-Host 'Warning: $(REG_FILE) not found' }"
	@echo "Batch rename helper files installed"

.PHONY: deploy
deploy:  ## Install and deploy batch rename
	@echo "Installing $(BINARY_NAME) globally..."
	go install $(CMD_PATH)

# Helper commands
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

.PHONY: install-deps
install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod tidy

.PHONY: all
all: fmt lint test build ## Run format, lint, test, and build
