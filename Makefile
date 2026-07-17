# ==============================================================================
# 1. LOAD ENVIRONMENT VARIABLES
# ==============================================================================
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ==============================================================================
# 2. VARIABLES & DEFAULTS
# ==============================================================================
APP_NAME=Go-CollabSpace
MIGRATION_DIR=migrations

DB_DSN := host=$(DB_HOST) \
          port=$(DB_PORT) \
          user=$(DB_USER) \
          password=$(DB_PASSWORD) \
          dbname=$(DB_NAME) \
          sslmode=disable

# ==============================================================================
# 3. COMMANDS
# ==============================================================================

.PHONY: help
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Development:"
	@echo "  run              Start the API server"
	@echo "  build            Build binary to bin/ folder"
	@echo "  tidy             Tidy go modules"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt              Format code with goimports"
	@echo "  fmt-check        Check if code is formatted"
	@echo "  lint             Run golangci-lint"
	@echo "  test             Run tests"
	@echo "  test-coverage    Run tests with coverage"
	@echo "  pre-commit       Run all checks (fmt + lint + test)"
	@echo ""
	@echo "Database Migrations (Goose):"
	@echo "  db-create name=  Create new migration (e.g., make db-create name=add_users)"
	@echo "  db-up            Apply all migrations"
	@echo "  db-down          Rollback last migration"
	@echo "  db-reset         Reset database (down-to 0 + up)"
	@echo "  db-status        Show migration status"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up        Start containers"
	@echo "  docker-down      Stop containers"
	@echo "  docker-logs      View logs"
	@echo "  hooks            Install git hooks via lefthook"

# ==============================================================================
# DEVELOPMENT
# ==============================================================================
.PHONY: run
run:
	@echo "Running server..."
	@go run cmd/server/main.go

.PHONY: build
build:
	@echo "🔨 Building binary..."
	@go build -o bin/$(APP_NAME) cmd/server/main.go
	@echo "Binary created: bin/$(APP_NAME)"

.PHONY: tidy
tidy:
	@echo "🧹 Tidying modules..."
	@go mod tidy
	@go mod verify
	@echo "Modules tidied!"

# ==============================================================================
# CODE QUALITY
# ==============================================================================
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code with goimports..."
	@goimports -w -local $(APP_NAME) .
	@echo "Code formatted!"

.PHONY: fmt-check
fmt-check:
	@echo "🔍 Checking code format..."
	@test -z "$$(goimports -l -local $(APP_NAME) . | tee /dev/stderr)" || \
		(echo "❌ Files need formatting. Run 'make fmt'" && exit 1)
	@echo "All files properly formatted!"

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@echo "No linting issues!"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -short ./...

.PHONY: test-coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: pre-commit
pre-commit: fmt lint test
	@echo "All pre-commit checks passed!"

.PHONY: hooks
hooks:
	@echo "🪝 Installing git hooks..."
	@lefthook install
	@echo "Git hooks installed!"
# ==============================================================================
# DATABASE MIGRATIONS (GOOSE)
# ==============================================================================
.PHONY: db-create
db-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: name is required"; \
		echo "Usage: make db-create name=add_users"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@goose -dir $(MIGRATION_DIR) create $(name) sql
	@echo "Migration created!"

.PHONY: db-up
db-up:
	@echo "⬆Running migrations..."
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" up
	@echo "Migrations applied!"

.PHONY: db-down
db-down:
	@echo "⬇Rolling back last migration..."
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down
	@echo "Migration rolled back!"

.PHONY: db-status
db-status:
	@echo "Migration status:"
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" status

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down-to 0
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" up
	@echo "Database reset complete!"

# ==============================================================================
# DOCKER
# ==============================================================================
.PHONY: docker-up
docker-up:
	@echo "Starting containers..."
	@docker-compose up -d
	@echo "Containers started!"

.PHONY: docker-down
docker-down:
	@echo "Stopping containers..."
	@docker-compose down
	@echo "Containers stopped!"

.PHONY: docker-logs
docker-logs:
	@docker-compose logs -f

# ==============================================================================
# TOOLS INSTALLATION
# ==============================================================================
.PHONY: install-tools
install-tools:
	@echo "📦 Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/evilmartians/lefthook@latest
	@echo "Tools installed!"

.PHONY: verify-tools
verify-tools:
	@echo "🔍 Verifying tools..."
	@which goimports || echo "❌ goimports not found"
	@which golangci-lint || echo "❌ golangci-lint not found"
	@which goose || echo "❌ goose not found"
	@echo "Verification complete!"