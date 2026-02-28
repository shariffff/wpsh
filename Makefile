.PHONY: all build install install-user test test-cli test-ansible test-full clean fmt lint help \
       molecule-install test-molecule test-molecule-provision test-molecule-website test-molecule-roles \
       lint-yaml lint-ansible docker-build docker-build-all \
       docker-build-linux-amd64 docker-build-linux-arm64 \
       docker-build-darwin-amd64 docker-build-darwin-arm64 docker-build-windows-amd64

# Version from version.txt
VERSION=$(shell cat version.txt)

# Default target
all: build

# Build CLI
build:
	@echo "Building WordMon CLI (v$(VERSION))..."
	@cd cli && make build
	@echo "✓ Build complete: cli/wordmon"

# Install CLI to /usr/local/bin (requires sudo)
install: build
	@echo "Installing WordMon CLI to /usr/local/bin..."
	@cd cli && make install
	@echo ""
	@echo "✓ Installation complete!"
	@echo ""
	@echo "Next step: Run 'wordmon init' to set up your environment"

# Install CLI to ~/bin (no sudo required)
install-user: build
	@echo "Installing WordMon CLI to ~/bin..."
	@cd cli && make install-user
	@echo ""
	@echo "✓ Installation complete!"
	@echo ""
	@echo "Next step: Run 'wordmon init' to set up your environment"

# Run all tests
test: test-cli test-ansible
	@echo "✓ All tests passed"

# Run CLI tests
test-cli:
	@echo "Running CLI tests..."
	@cd cli && make test

# Test Ansible playbooks (syntax check)
test-ansible:
	@echo "Validating Ansible playbooks..."
	@cd ansible && ansible-playbook --syntax-check provision.yml
	@cd ansible && ansible-playbook --syntax-check website.yml
	@cd ansible && ansible-playbook --syntax-check playbooks/domain_management.yml
	@cd ansible && ansible-playbook --syntax-check playbooks/delete_site.yml
	@echo "✓ Ansible syntax validation passed"

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@cd cli && make fmt

# Lint Go code
lint:
	@echo "Linting Go code..."
	@cd cli && make lint

# Docker build (no Go installation required)
docker-build:
	@echo "Building WordMon CLI via Docker (v$(VERSION))..."
	@cd cli && make docker-build VERSION=$(VERSION)
	@echo "✓ Docker build complete: cli/wordmon"

# Docker build for all platforms
docker-build-all:
	@echo "Building WordMon CLI for all platforms (v$(VERSION))..."
	@cd cli && make docker-build-all VERSION=$(VERSION)

# Individual platform builds
docker-build-linux-amd64:
	@cd cli && make docker-build-linux-amd64 VERSION=$(VERSION)

docker-build-linux-arm64:
	@cd cli && make docker-build-linux-arm64 VERSION=$(VERSION)

docker-build-darwin-amd64:
	@cd cli && make docker-build-darwin-amd64 VERSION=$(VERSION)

docker-build-darwin-arm64:
	@cd cli && make docker-build-darwin-arm64 VERSION=$(VERSION)

docker-build-windows-amd64:
	@cd cli && make docker-build-windows-amd64 VERSION=$(VERSION)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@cd cli && make clean
	@rm -f wordmon
	@echo "✓ Clean complete"

# Development: quick run
run: build
	@./cli/wordmon

# ===== Molecule Testing =====

# Install Molecule and dependencies
molecule-install:
	@echo "Installing Molecule and dependencies..."
	pip install molecule molecule-plugins[docker] ansible-lint yamllint
	@echo "✓ Molecule installation complete"

# Run all Molecule tests (roles + integration)
test-molecule: test-molecule-roles test-molecule-provision test-molecule-website
	@echo "✓ All Molecule tests passed"

# Run individual role tests
test-molecule-roles:
	@echo "Running Molecule tests for all roles..."
	@cd ansible/roles/bootstrap && molecule test
	@cd ansible/roles/database && molecule test
	@cd ansible/roles/nginx && molecule test
	@cd ansible/roles/php && molecule test
	@cd ansible/roles/security && molecule test
	@echo "✓ All role tests passed"

# Run provision integration test
test-molecule-provision:
	@echo "Running provision integration test..."
	@cd ansible && molecule test -s provision
	@echo "✓ Provision integration test passed"

# Run website integration test
test-molecule-website:
	@echo "Running website integration test..."
	@cd ansible && molecule test -s website
	@echo "✓ Website integration test passed"

# Lint YAML files
lint-yaml:
	@echo "Linting YAML files..."
	@cd ansible && yamllint .
	@echo "✓ YAML lint passed"

# Lint Ansible files
lint-ansible:
	@echo "Linting Ansible files..."
	@cd ansible && ansible-lint
	@echo "✓ Ansible lint passed"

# Full test suite (CLI + Ansible syntax + Molecule)
test-full: test lint-yaml lint-ansible test-molecule
	@echo "✓ Full test suite passed"

# Show help
help:
	@echo "WordMon Build System"
	@echo "Version: $(VERSION)"
	@echo ""
	@echo "Available targets:"
	@echo "  build                   - Build the CLI binary (requires Go 1.24+)"
	@echo "  install                 - Install CLI to /usr/local/bin (requires sudo)"
	@echo "  install-user            - Install CLI to ~/bin (no sudo required)"
	@echo ""
	@echo "Docker builds (no Go required):"
	@echo "  docker-build            - Build for current platform"
	@echo "  docker-build-all        - Build for all platforms"
	@echo "  docker-build-linux-amd64"
	@echo "  docker-build-linux-arm64"
	@echo "  docker-build-darwin-amd64  (macOS Intel)"
	@echo "  docker-build-darwin-arm64  (macOS Apple Silicon)"
	@echo "  docker-build-windows-amd64"
	@echo ""
	@echo "Development:"
	@echo "  test                    - Run all tests (CLI + Ansible syntax)"
	@echo "  test-cli                - Run CLI tests only"
	@echo "  test-ansible            - Validate Ansible playbook syntax"
	@echo "  fmt                     - Format Go code"
	@echo "  lint                    - Lint Go code"
	@echo "  clean                   - Remove build artifacts"
	@echo "  run                     - Build and run CLI"
	@echo ""
	@echo "Molecule Testing (integration):"
	@echo "  molecule-install        - Install Molecule and dependencies"
	@echo "  test-molecule           - Run all Molecule tests"
	@echo "  test-molecule-roles     - Run individual role tests"
	@echo "  test-molecule-provision - Run provision integration test"
	@echo "  test-molecule-website   - Run website integration test"
	@echo "  lint-yaml               - Lint YAML files"
	@echo "  lint-ansible            - Lint Ansible files"
	@echo "  test-full               - Run full test suite (CLI + lint + Molecule)"
	@echo ""
	@echo "  help                    - Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  make build && ./cli/wordmon init"
