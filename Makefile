# openribcage - A2A Protocol Client
# Build automation for Go project

.PHONY: help build test lint clean install dev-setup docker run-kagent-test

# Variables
GO_VERSION := 1.21
BINARY_NAME := openribcage
BUILD_DIR := ./bin
CMD_DIR := ./cmd
PKG_DIR := ./pkg

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development setup
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@go version
	@go mod download
	@go mod tidy
	@which golangci-lint || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "Development environment ready!"

# Build targets
build: ## Build the main application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(BINARY_NAME)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build all applications
	@echo "Building all applications..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(BINARY_NAME)
	@go build -o $(BUILD_DIR)/discovery $(CMD_DIR)/discovery
	@echo "All builds complete in $(BUILD_DIR)/"

install: ## Install applications to $GOPATH/bin
	@echo "Installing applications..."
	@go install $(CMD_DIR)/$(BINARY_NAME)
	@go install $(CMD_DIR)/discovery
	@echo "Installation complete"

# Testing targets
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests (requires kagent sandbox)
	@echo "Running integration tests..."
	@go test -v -tags=integration ./test/integration/...

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

# A2A specific targets
run-kagent-test: ## Test A2A client against kagent sandbox
	@echo "Testing A2A client with kagent..."
	@./scripts/test-a2a-client.sh

validate-agentcard: ## Validate AgentCard format
	@echo "Validating AgentCard format..."
	@go run ./tools/agentcard-validator

a2a-compliance: ## Run A2A protocol compliance tests
	@echo "Running A2A compliance tests..."
	@go run ./tools/a2a-tester

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t openribcage:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 openribcage:latest

# Cleanup targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

clean-deps: ## Clean dependency cache
	@echo "Cleaning dependency cache..."
	@go clean -modcache

# Development workflow
dev: clean fmt lint test build ## Complete development workflow

release: clean check build-all ## Prepare release build
	@echo "Release build complete!"
