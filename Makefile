.PHONY: help install install-docker start-docker stop-docker check-docker build build-docker run run-docker run-docker-bg stop-docker-container clean clean-docker dev test migrate-db setup-db

# Default target
help:
	@echo "EchoHub Community Forum - Makefile Commands"
	@echo "=========================================="
	@echo ""
	@echo "Setup & Installation:"
	@echo "  make install              - Install Go dependencies"
	@echo "  make install-docker       - Install Docker (rootless)"
	@echo "  make setup-db             - Setup database (run migrations)"
	@echo ""
	@echo "Docker Management:"
	@echo "  make start-docker         - Start Docker daemon (rootless)"
	@echo "  make stop-docker          - Stop Docker daemon"
	@echo "  make check-docker         - Check if Docker is running"
	@echo ""
	@echo "Build:"
	@echo "  make build                - Build Go application"
	@echo "  make build-docker         - Build Docker image"
	@echo ""
	@echo "Run:"
	@echo "  make run                  - Run application locally"
	@echo "  make run-docker           - Run Docker container (foreground)"
	@echo "  make run-docker-bg        - Run Docker container (background)"
	@echo "  make stop-docker-container - Stop running Docker container"
	@echo ""
	@echo "Development:"
	@echo "  make dev                  - Run in development mode"
	@echo "  make test                 - Run tests"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make clean-docker         - Clean Docker images and containers"
	@echo ""

# Setup environment variables for rootless Docker
export PATH := $(HOME)/bin:$(PATH)
export DOCKER_HOST := unix://$(XDG_RUNTIME_DIR)/docker.sock

# Install Go dependencies
install:
	@echo "📦 Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Install Docker (rootless)
install-docker:
	@echo "📦 Installing Docker (rootless)..."
	@if [ -f ./install-docker-rootless.zsh ]; then \
		zsh ./install-docker-rootless.zsh; \
	else \
		echo "❌ install-docker-rootless.zsh not found!"; \
		exit 1; \
	fi

# Start Docker daemon
start-docker:
	@echo "🚀 Starting Docker daemon..."
	@if [ -f ./start-docker.sh ]; then \
		./start-docker.sh; \
	else \
		echo "⚠️  start-docker.sh not found, trying direct command..."; \
		nohup dockerd-rootless.sh > ~/docker-rootless.log 2>&1 & \
		echo "✅ Docker daemon started (log: ~/docker-rootless.log)"; \
	fi
	@sleep 2
	@$(MAKE) check-docker

# Stop Docker daemon
stop-docker:
	@echo "🛑 Stopping Docker daemon..."
	@pkill -f dockerd-rootless || echo "⚠️  Docker daemon not running"
	@echo "✅ Docker daemon stopped"

# Check if Docker is running
check-docker:
	@echo "🔍 Checking Docker status..."
	@docker ps > /dev/null 2>&1 && echo "✅ Docker is running" || (echo "❌ Docker is not running. Run 'make start-docker'" && exit 1)

# Build Go application
build:
	@echo "🔨 Building Go application..."
	go build -o bin/echohubApp ./cmd/api/main.go
	@echo "✅ Build complete! Binary: ./bin/echohubApp"

# Build Docker image
build-docker: check-docker
	@echo "🔨 Building Docker image..."
	docker build -t echohub-community-app .
	@echo "✅ Docker image built: echohub-community-app"

# Run application locally
run: build
	@echo "🚀 Starting application..."
	./bin/echohubApp

# Run Docker container (foreground)
run-docker: build-docker
	@echo "🚀 Starting Docker container..."
	@docker stop echohub-community 2>/dev/null || true
	@docker rm echohub-community 2>/dev/null || true
	docker run -p 8080:8080 --name echohub-community echohub-community-app

# Run Docker container (background)
run-docker-bg: build-docker
	@echo "🚀 Starting Docker container in background..."
	@docker stop echohub-community 2>/dev/null || true
	@docker rm echohub-community 2>/dev/null || true
	@docker run -d -p 8080:8080 --name echohub-community echohub-community-app
	@sleep 2
	@docker ps | grep echohub-community || echo "⚠️  Container may not have started"
	@echo "✅ Container running! Check with: docker ps"
	@echo "📋 View logs: docker logs echohub-community"
	@echo "🌐 Application: http://localhost:8080"

# Stop Docker container
stop-docker-container:
	@echo "🛑 Stopping Docker container..."
	@docker stop echohub-community 2>/dev/null && echo "✅ Container stopped" || echo "⚠️  Container not running"
	@docker rm echohub-community 2>/dev/null || true

# Development mode (with auto-reload if you have air/realize installed)
dev: install
	@echo "🔧 Starting in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	elif command -v realize > /dev/null; then \
		realize start; \
	else \
		echo "⚠️  No auto-reload tool found. Install 'air' or 'realize' for auto-reload."; \
		echo "Running normally..."; \
		$(MAKE) run; \
	fi

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Setup database (run migrations)
setup-db:
	@echo "🗄️  Setting up database..."
	@if [ -f ./internal/database/migration/tables.sql ]; then \
		echo "⚠️  Database migrations should be run manually or via application initialization"; \
		echo "Database file: ./internal/database/echohub-community.db"; \
	else \
		echo "❌ Migration files not found!"; \
	fi

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f *.log
	@echo "✅ Clean complete!"

# Clean Docker images and containers
clean-docker: stop-docker-container
	@echo "🧹 Cleaning Docker resources..."
	@docker rmi echohub-community-app 2>/dev/null && echo "✅ Docker image removed" || echo "⚠️  Image not found"
	@echo "✅ Docker cleanup complete!"

# Full clean (everything)
clean-all: clean clean-docker
	@echo "✅ Full cleanup complete!"
