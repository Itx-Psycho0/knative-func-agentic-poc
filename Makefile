# Makefile for Knative Functions Agentic Workflow POC

# Variables
BINARY_NAME=func-agentic
BUILD_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"
COVERAGE_FILE=coverage.out

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

.PHONY: all build clean test test-unit test-integration test-e2e test-coverage lint fmt vet check help install deps

# Default target
all: clean deps fmt vet lint test build

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)Knative Functions Agentic Workflow POC$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'

## build: Build the binary
build:
	@echo "$(COLOR_BLUE)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "$(COLOR_GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

## clean: Remove build artifacts
clean:
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE)
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## deps: Download dependencies
deps:
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

## fmt: Format code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GO) vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## lint: Run golangci-lint
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Lint complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(COLOR_RESET)"; \
	fi

## test: Run all tests
test: test-unit

## test-unit: Run unit tests
test-unit:
	@echo "$(COLOR_BLUE)Running unit tests...$(COLOR_RESET)"
	$(GO) test -v -race -timeout 30s ./...
	@echo "$(COLOR_GREEN)✓ Unit tests complete$(COLOR_RESET)"

## test-integration: Run integration tests
test-integration:
	@echo "$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	$(GO) test -v -race -timeout 5m -tags=integration ./tests/integration/...
	@echo "$(COLOR_GREEN)✓ Integration tests complete$(COLOR_RESET)"

## test-e2e: Run end-to-end tests
test-e2e:
	@echo "$(COLOR_BLUE)Running end-to-end tests...$(COLOR_RESET)"
	$(GO) test -v -timeout 10m -tags=e2e ./tests/e2e/...
	@echo "$(COLOR_GREEN)✓ End-to-end tests complete$(COLOR_RESET)"

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	$(GO) test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

## install: Install the binary
install: build
	@echo "$(COLOR_BLUE)Installing $(BINARY_NAME)...$(COLOR_RESET)"
	$(GO) install ./cmd/main.go
	@echo "$(COLOR_GREEN)✓ Installed$(COLOR_RESET)"

## run: Run the application
run: build
	@echo "$(COLOR_BLUE)Running $(BINARY_NAME)...$(COLOR_RESET)"
	./$(BUILD_DIR)/$(BINARY_NAME)

## docker-build: Build Docker image
docker-build:
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	docker build -t $(BINARY_NAME):latest .
	@echo "$(COLOR_GREEN)✓ Docker image built$(COLOR_RESET)"

## setup-dev: Setup development environment
setup-dev:
	@echo "$(COLOR_BLUE)Setting up development environment...$(COLOR_RESET)"
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(COLOR_GREEN)✓ Development environment ready$(COLOR_RESET)"

## demo: Run demo scenario
demo: build
	@echo "$(COLOR_BLUE)Running demo scenario...$(COLOR_RESET)"
	./scripts/run-demo.sh

## docs: Generate documentation
docs:
	@echo "$(COLOR_BLUE)Generating documentation...$(COLOR_RESET)"
	$(GO) doc -all ./... > docs/API_REFERENCE.md
	@echo "$(COLOR_GREEN)✓ Documentation generated$(COLOR_RESET)"

# Development helpers
.PHONY: watch
watch:
	@echo "$(COLOR_BLUE)Watching for changes...$(COLOR_RESET)"
	@while true; do \
		make build; \
		inotifywait -qre close_write .; \
	done
