.PHONY: build install clean test docker-build docker-build-all \
       docker-build-linux-amd64 docker-build-linux-arm64 \
       docker-build-darwin-amd64 docker-build-darwin-arm64 \
       docker-build-windows-amd64

# Binary name
BINARY_NAME=wordmon

# Build information
VERSION?=dev
COMMIT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
LDFLAGS=-ldflags "-X github.com/wordmon/cli/cmd.Version=$(VERSION) -X github.com/wordmon/cli/cmd.CommitSHA=$(COMMIT_SHA) -X github.com/wordmon/cli/cmd.BuildDate=$(BUILD_DATE)"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete. Run 'wordmon --help' to get started."

# Install to current user's bin directory (no sudo required)
install-user: build
	@echo "Installing $(BINARY_NAME) to ~/bin..."
	@mkdir -p ~/bin
	@cp $(BINARY_NAME) ~/bin/
	@echo "Installation complete."
	@echo "Make sure ~/bin is in your PATH."
	@echo "Add this to your ~/.bashrc or ~/.zshrc if not already present:"
	@echo '  export PATH="$$HOME/bin:$$PATH"'

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

# Run the binary
run: build
	@./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed. Run: brew install golangci-lint"

# Docker build - builds for current platform
docker-build:
	@echo "Building $(BINARY_NAME) using Docker..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--output type=local,dest=. \
		--target binary \
		.
	@echo "Build complete: ./$(BINARY_NAME)"

# Docker build for Linux (amd64)
docker-build-linux-amd64:
	@echo "Building $(BINARY_NAME) for Linux (amd64)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GO_OS=linux \
		--build-arg GO_ARCH=amd64 \
		--output type=local,dest=./dist/linux-amd64 \
		--target binary \
		.
	@echo "Build complete: ./dist/linux-amd64/$(BINARY_NAME)"

# Docker build for Linux (arm64)
docker-build-linux-arm64:
	@echo "Building $(BINARY_NAME) for Linux (arm64)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GO_OS=linux \
		--build-arg GO_ARCH=arm64 \
		--output type=local,dest=./dist/linux-arm64 \
		--target binary \
		.
	@echo "Build complete: ./dist/linux-arm64/$(BINARY_NAME)"

# Docker build for macOS (amd64 - Intel)
docker-build-darwin-amd64:
	@echo "Building $(BINARY_NAME) for macOS (amd64)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GO_OS=darwin \
		--build-arg GO_ARCH=amd64 \
		--output type=local,dest=./dist/darwin-amd64 \
		--target binary \
		.
	@echo "Build complete: ./dist/darwin-amd64/$(BINARY_NAME)"

# Docker build for macOS (arm64 - Apple Silicon)
docker-build-darwin-arm64:
	@echo "Building $(BINARY_NAME) for macOS (arm64)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GO_OS=darwin \
		--build-arg GO_ARCH=arm64 \
		--output type=local,dest=./dist/darwin-arm64 \
		--target binary \
		.
	@echo "Build complete: ./dist/darwin-arm64/$(BINARY_NAME)"

# Docker build for Windows (amd64)
docker-build-windows-amd64:
	@echo "Building $(BINARY_NAME) for Windows (amd64)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT_SHA=$(COMMIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GO_OS=windows \
		--build-arg GO_ARCH=amd64 \
		--output type=local,dest=./dist/windows-amd64 \
		--target binary \
		.
	@mv ./dist/windows-amd64/wordmon ./dist/windows-amd64/wordmon.exe 2>/dev/null || true
	@echo "Build complete: ./dist/windows-amd64/$(BINARY_NAME).exe"

# Docker build for all platforms
docker-build-all: docker-build-linux-amd64 docker-build-linux-arm64 docker-build-darwin-amd64 docker-build-darwin-arm64 docker-build-windows-amd64
	@echo ""
	@echo "All platform builds complete:"
	@echo "  ./dist/linux-amd64/wordmon"
	@echo "  ./dist/linux-arm64/wordmon"
	@echo "  ./dist/darwin-amd64/wordmon"
	@echo "  ./dist/darwin-arm64/wordmon"
	@echo "  ./dist/windows-amd64/wordmon.exe"

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary (requires Go 1.24+)"
	@echo "  install            - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user       - Install to ~/bin (no sudo)"
	@echo "  test               - Run tests"
	@echo "  clean              - Remove build artifacts"
	@echo "  run                - Build and run the binary"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo ""
	@echo "Docker builds (no Go installation required):"
	@echo "  docker-build              - Build for current platform"
	@echo "  docker-build-linux-amd64  - Build for Linux (amd64)"
	@echo "  docker-build-linux-arm64  - Build for Linux (arm64)"
	@echo "  docker-build-darwin-amd64 - Build for macOS Intel"
	@echo "  docker-build-darwin-arm64 - Build for macOS Apple Silicon"
	@echo "  docker-build-windows-amd64- Build for Windows (amd64)"
	@echo "  docker-build-all          - Build for all platforms"
