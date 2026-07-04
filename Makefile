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

# Detect package manager (apt / apk / dnf / brew)
APT   := $(shell command -v apt-get  2>/dev/null)
APK   := $(shell command -v apk      2>/dev/null)
DNF   := $(shell command -v dnf      2>/dev/null)
BREW  := $(shell command -v brew     2>/dev/null)

# Use sudo for system package installs on Linux; brew on macOS doesn't need it.
SUDO := $(if $(BREW),,sudo)

.PHONY: install-deps-unix
install-deps-unix: ## Install system + Go deps on Linux/macOS (Linux: run as root or with sudo make)
	@echo "==> Installing system packages..."
ifdef APT
	$(SUDO) apt-get update -qq
	$(SUDO) apt-get install -y --no-install-recommends \
		make git curl wget ca-certificates gnupg \
		chromium-browser
else ifdef APK
	$(SUDO) apk add --no-cache make git curl wget ca-certificates chromium
else ifdef DNF
	$(SUDO) dnf install -y make git curl wget ca-certificates chromium
else ifdef BREW
	brew install make git curl wget
	brew install --cask google-chrome
else
	$(error No supported package manager found: expected apt-get, apk, dnf, or brew)
endif
	@echo "==> Installing Go (if not present)..."
	@if ! command -v go >/dev/null 2>&1 && [ ! -f /usr/local/go/bin/go ]; then \
		curl -fsSL https://golang.org/dl/go1.22.4.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
		$(SUDO) tar -C /usr/local -xzf /tmp/go.tar.gz && \
		echo 'export PATH=$$PATH:/usr/local/go/bin' | $(SUDO) tee /etc/profile.d/go.sh > /dev/null && \
		echo "Go installed — open a new shell or run: source /etc/profile.d/go.sh"; \
	else \
		echo "Go already installed"; \
	fi
	@echo "==> Installing golangci-lint..."
	@PATH="$${PATH}:/usr/local/go/bin" GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "==> Running go mod tidy..."
	PATH="$${PATH}:/usr/local/go/bin" go mod tidy
	@echo "==> Done. Run 'make build' to compile."

.PHONY: all
all: fmt lint test build ## Run format, lint, test, and build
