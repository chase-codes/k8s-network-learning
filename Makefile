.PHONY: build run fmt test doctor clean install help

# Default target
help: ## Show this help message
	@echo "NetLab - Interactive Networking Learning Environment"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the NetLab binary
	@echo "🔨 Building NetLab..."
	@go build -o bin/netlab .
	@echo "✅ Build complete! Binary available at bin/netlab"

run: ## Run NetLab in development mode
	@echo "🚀 Starting NetLab..."
	@go run . start

install: build ## Install NetLab binary to /usr/local/bin
	@echo "📦 Installing NetLab..."
	@sudo cp bin/netlab /usr/local/bin/
	@echo "✅ NetLab installed! Run 'netlab start' from anywhere."

fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests completed!"

doctor: build ## Run environment diagnostics
	@echo "🔍 Running environment diagnostics..."
	@./bin/netlab doctor

clean: ## Clean build artifacts
	@echo "🧹 Cleaning up..."
	@rm -rf bin/
	@go clean
	@echo "✅ Cleanup complete!"

setup: ## Run setup script to check dependencies
	@echo "🛠️  Running setup script..."
	@./scripts/setup.sh

deps: ## Download and verify dependencies
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "✅ Dependencies ready!"

# Development helpers
dev-run: ## Run with auto-reload (requires 'air' tool)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ 'air' not found. Install with: go install github.com/cosmtrek/air@latest"; \
	fi

lint: ## Run linter (requires golangci-lint)
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ 'golangci-lint' not found. Install from: https://golangci-lint.run/"; \
	fi

all: clean deps fmt test build ## Run full build pipeline

# Module shortcuts
osi: build ## Run OSI Model module
	@./bin/netlab module 01-osi-model

# Quick start
start: build ## Build and start NetLab
	@./bin/netlab start 