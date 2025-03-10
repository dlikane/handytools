# Variables
BINARY_NAME = img
CMD_PATH = ./cmd/img

# Default target (shows available commands)
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the img binary"
	@echo "  run        - Run the application"
	@echo "  deploy     - Install img globally"
	@echo "  clean      - Remove built binary"
	@echo "  fmt        - Format Go code"
	@echo "  lint       - Run linter (golangci-lint)"
	@echo "  test       - Run tests"
	@echo "  install    - Install dependencies"
	@echo "  all        - Run format, lint, test, and build"

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(CMD_PATH)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Install (Deploy) the binary to GOPATH/bin
deploy:
	@echo "Installing $(BINARY_NAME) globally..."
	go install $(CMD_PATH)
	@echo "Deployment complete! Run '$(BINARY_NAME)' from anywhere."

# Clean the build output
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Format the code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Lint the code (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod tidy

# Run all steps: format, lint, test, and build
all: fmt lint test build
