.PHONY: build clean version help dev test-build install uninstall test test-package test-multiple

# Default target
all: build

# Build with version information
build:
	@./build.sh

# Development build with custom version
dev:
	@echo "Building aurvt for development..."
	@./build.sh "0.0.0-dev" "$(shell git rev-parse --short HEAD 2>/dev/null || echo 'development')" "dev"

# Test build with specific version
test-build:
	@echo "Building aurvt for testing..."
	@./build.sh "1.0.0-test" "$(shell git rev-parse --short HEAD 2>/dev/null || echo 'test')" "test"

# Quick build without version injection
quick:
	@echo "Building aurvt (quick build)..."
	@go build -o aurvt

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f aurvt

# Show version information
version:
	@./aurvt version

# Test the binary on a sample package
test-package:
	@echo "Testing aurvt on python-conda-libmamba-solver..."
	@./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver

# Test the binary on multiple packages
test-multiple:
	@echo "Testing aurvt on multiple packages..."
	@./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver
	@echo ""
	@./aurvt ../aur-pkgbuilds/rot8-git
	@echo ""
	@./aurvt ../aur-pkgbuilds/openmohaa-git

# Install to system (requires sudo)
install: build
	@echo "Installing aurvt to /usr/local/bin..."
	@sudo cp aurvt /usr/local/bin/
	@echo "Installation complete!"

# Uninstall from system
uninstall:
	@echo "Removing aurvt from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/aurvt
	@echo "Uninstallation complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build with version information (default)"
	@echo "  dev           - Development build with custom version"
	@echo "  test-build    - Test build with specific version"
	@echo "  quick         - Quick build without version injection"
	@echo "  clean         - Remove build artifacts"
	@echo "  version       - Show version information"
	@echo "  test-package  - Test binary on sample package"
	@echo "  test-multiple - Test binary on multiple packages"
	@echo "  install       - Install to /usr/local/bin"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  test          - Run tests"
	@echo "  help          - Show this help"