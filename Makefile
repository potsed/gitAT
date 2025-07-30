# GitAT Makefile
# Build and manage the GitAT Go application

# Variables
BINARY_NAME=git-@
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.1.0")
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building GitAT..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gitat

# Build for current platform
.PHONY: build-local
build-local:
	@echo "Building GitAT for local platform..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gitat

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building GitAT for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/gitat
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/gitat
	
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/gitat
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/gitat
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/gitat
	GOOS=windows GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/gitat

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install the binary and Git extension
.PHONY: install
install: build-local
	@echo "Installing GitAT..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "GitAT installed to /usr/local/bin/$(BINARY_NAME)"
	@echo "Git extension 'git @' is now available"

# Uninstall the binary and Git extension
.PHONY: uninstall
uninstall:
	@echo "Uninstalling GitAT..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "GitAT uninstalled"

# Run the application
.PHONY: run
run: build-local
	@echo "Running GitAT..."
	./$(BINARY_NAME)

# Show help
.PHONY: help
help:
	@echo "GitAT Makefile"
	@echo "=============="
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  build-local  - Build for current platform"
	@echo "  build-all    - Build for all platforms (Linux, macOS, Windows)"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  install      - Install binary and Git extension to /usr/local/bin"
	@echo "  uninstall    - Remove binary and Git extension from /usr/local/bin"
	@echo "  run          - Build and run the application"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_HASH)"
	@echo "Build Date: $(BUILD_DATE)"

# Development helpers
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development environment ready!"

# Create release
.PHONY: release
release: build-all
	@echo "Creating release..."
	@mkdir -p release
	@cd $(BUILD_DIR) && tar -czf ../release/gitat-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/gitat-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf ../release/gitat-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/gitat-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip ../release/gitat-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@cd $(BUILD_DIR) && zip ../release/gitat-$(VERSION)-windows-arm64.zip $(BINARY_NAME)-windows-arm64.exe
	@echo "Release packages created in release/ directory" 