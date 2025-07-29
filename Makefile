.PHONY: build clean version help

# Default target
all: build

# Build with version information
build:
	@./build.sh

# Quick build without version injection
quick:
	@echo "Building aurvt (quick build)..."
	@go build -o aurvt

publish: build
	@echo "Publishing aurvt to GitHub Releases..."
	@git tag v$(VERSION)
	@git push origin --tags
	@gh release create v$(VERSION) aurvt

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f aurvt

# Show version information
version:
	@./aurvt version

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
	@echo "  build     - Build with version information (default)"
	@echo "  quick     - Quick build without version injection"
	@echo "  clean     - Remove build artifacts"
	@echo "  version   - Show version information"
	@echo "  install   - Install to /usr/local/bin"
	@echo "  uninstall - Remove from /usr/local/bin"
	@echo "  test      - Run tests"
	@echo "  help      - Show this help"