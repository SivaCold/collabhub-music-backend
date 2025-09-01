# ===========================================
# CollabHub Music Backend Makefile
# ===========================================

# Variables
BINARY_NAME=collabhub-backend
BINARY_PATH=./$(BINARY_NAME)
MAIN_PATH=./cmd/server/main.go
PKG_PATH=./...
DOCKER_IMAGE=collabhub-backend
DOCKER_TAG=latest

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-w -s"
BUILD_FLAGS=-v $(LDFLAGS)

# Environment
ENV_FILE=.env
ENV_EXAMPLE=.env.example

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test coverage deps help run dev docker docker-build docker-run setup certificates

# Default target
all: deps build

## Build the application
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Build completed successfully!$(NC)"

## Build for multiple platforms
build-cross:
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "$(GREEN)Cross-platform build completed!$(NC)"

## Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning up...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
	rm -rf ./dist
	@echo "$(GREEN)Clean completed!$(NC)"

## Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies installed!$(NC)"

## Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...

## Run tests with coverage
coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## Run tests with race detection
test-race:
	@echo "$(BLUE)Running tests with race detection...$(NC)"
	$(GOTEST) -race -short ./...

## Run benchmarks
bench:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

## Run the application in development mode
run: build
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

## Run with hot reload (requires air)
dev:
	@echo "$(BLUE)Starting development server with hot reload...$(NC)"
	@if ! command -v air > /dev/null; then \
		echo "$(YELLOW)Installing air for hot reload...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

## Setup development environment
setup: setup-env certificates deps
	@echo "$(GREEN)Development environment setup completed!$(NC)"

## Setup environment file
setup-env:
	@echo "$(BLUE)Setting up environment file...$(NC)"
	@if [ ! -f $(ENV_FILE) ]; then \
		if [ -f $(ENV_EXAMPLE) ]; then \
			cp $(ENV_EXAMPLE) $(ENV_FILE); \
			echo "$(GREEN)Created $(ENV_FILE) from $(ENV_EXAMPLE)$(NC)"; \
		else \
			echo "$(RED)$(ENV_EXAMPLE) not found!$(NC)"; \
			exit 1; \
		fi \
	else \
		echo "$(YELLOW)$(ENV_FILE) already exists$(NC)"; \
	fi

## Generate TLS certificates for development
certificates:
	@echo "$(BLUE)Generating TLS certificates...$(NC)"
	@mkdir -p certs
	@if [ ! -f certs/server.crt ] || [ ! -f certs/server.key ]; then \
		openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \
			-out certs/server.crt -days 365 -nodes \
			-subj "/C=US/ST=State/L=City/O=CollabHub/CN=localhost"; \
		chmod 600 certs/server.key; \
		chmod 644 certs/server.crt; \
		echo "$(GREEN)TLS certificates generated!$(NC)"; \
	else \
		echo "$(YELLOW)TLS certificates already exist$(NC)"; \
	fi

## Create necessary directories
dirs:
	@echo "$(BLUE)Creating necessary directories...$(NC)"
	@mkdir -p certs storage logs
	@mkdir -p storage/{audio,images,temp,backups}
	@echo "$(GREEN)Directories created!$(NC)"

## Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

## Lint code
lint:
	@echo "$(BLUE)Linting code...$(NC)"
	@if ! command -v golangci-lint > /dev/null; then \
		echo "$(YELLOW)Installing golangci-lint...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	golangci-lint run

## Check for security issues
security:
	@echo "$(BLUE)Checking for security issues...$(NC)"
	@if ! command -v gosec > /dev/null; then \
		echo "$(YELLOW)Installing gosec...$(NC)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...

## Generate API documentation
docs:
	@echo "$(BLUE)Generating API documentation...$(NC)"
	@if ! command -v swag > /dev/null; then \
		echo "$(YELLOW)Installing swag...$(NC)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g cmd/server/main.go -o ./docs
	@echo "$(GREEN)API documentation generated!$(NC)"

## Build Docker image
docker-build:
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

## Run with Docker Compose
docker-up:
	@echo "$(BLUE)Starting services with Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"
	@echo "$(BLUE)Service URLs:$(NC)"
	@echo "  • Backend API: https://localhost:8443"
	@echo "  • Swagger Docs: https://localhost:8443/swagger/index.html"
	@echo "  • Keycloak Admin: http://localhost:8080/admin"

## Stop Docker Compose services
docker-down:
	@echo "$(BLUE)Stopping services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped!$(NC)"

## View Docker Compose logs
docker-logs:
	docker-compose logs -f

## Database operations
db-migrate:
	@echo "$(BLUE)Running database migrations...$(NC)"
	@# Migrations are run automatically by the application
	@echo "$(GREEN)Database migrations completed!$(NC)"

## Reset database
db-reset:
	@echo "$(BLUE)Resetting database...$(NC)"
	docker-compose down
	docker volume rm collabhub-postgres-data || true
	docker-compose up -d postgres
	@echo "$(GREEN)Database reset completed!$(NC)"

## Database shell
db-shell:
	@echo "$(BLUE)Opening database shell...$(NC)"
	docker-compose exec postgres psql -U collabhub_user -d collabhub_music

## Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	@echo "$(GREEN)Development tools installed!$(NC)"

## Check application health
health:
	@echo "$(BLUE)Checking application health...$(NC)"
	@curl -k -s https://localhost:8443/health | jq . || echo "$(RED)Health check failed or jq not installed$(NC)"

## Load test endpoints
load-test:
	@echo "$(BLUE)Running basic load test...$(NC)"
	@if ! command -v hey > /dev/null; then \
		echo "$(YELLOW)Installing hey for load testing...$(NC)"; \
		go install github.com/rakyll/hey@latest; \
	fi
	hey -n 1000 -c 10 -k https://localhost:8443/health

## Show application version
version:
	@echo "CollabHub Music Backend v1.0.0"

## Deploy to production (placeholder)
deploy:
	@echo "$(BLUE)Deploying to production...$(NC)"
	@echo "$(YELLOW)Production deployment should be handled by CI/CD pipeline$(NC)"

## Show help
help:
	@echo "$(BLUE)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BLUE)Examples:$(NC)"
	@echo "  make setup          # First time setup"
	@echo "  make dev            # Start development server with hot reload"
	@echo "  make docker-up      # Start all services with Docker"
	@echo "  make test           # Run tests"
	@echo "  make docs           # Generate API documentation"
