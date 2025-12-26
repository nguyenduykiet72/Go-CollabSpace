# ==============================================================================
# 1. LOAD ENVIRONMENT VARIABLES
# ==============================================================================
# Check nếu file .env tồn tại thì load vào
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ==============================================================================
# 2. VARIABLES & DEFAULTS
# ==============================================================================
# Tên App
APP_NAME=Go-CollabSpace

# Đường dẫn Migration
MIGRATION_DIR=migrations

# Database DSN: Nếu trong .env không có, dùng giá trị mặc định (Fallback)
# Lưu ý: DB_URL phải khớp với format của Goose/Gorm
DB_DSN := host=$(DB_HOST) \
          port=$(DB_PORT) \
          user=$(DB_USER) \
          password=$(DB_PASSWORD) \
          dbname=$(DB_NAME) \
          sslmode=disable
# ==============================================================================
# 3. COMMANDS
# ==============================================================================

# Help: Hiển thị danh sách lệnh (Mặc định khi gõ 'make')
.PHONY: help
help:
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  run              Start the API server (cmd/server/main.go)"
	@echo "  build            Build binary to bin/ folder"
	@echo "  tidy             Format code and tidy mod"
	@echo ""
	@echo "Database Migrations (Goose):"
	@echo "  db-create name=  Create a new SQL migration file (e.g., make db-create name=add_users)"
	@echo "  db-up            Apply all up migrations"
	@echo "  db-down          Rollback the last migration"
	@echo "  db-reset         Rollback all and re-apply (Fresh DB)"
	@echo "  db-status        Show migration status"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up        Start DB containers (docker-compose up -d)"
	@echo "  docker-down      Stop DB containers"

# --- DEVELOPMENT ---
.PHONY: run
run:
	@echo "Running server..."
	go run cmd/server/main.go

.PHONY: build
build:
	@echo "Building binary..."
	go build -o bin/$(APP_NAME) cmd/server/main.go

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy
	go mod verify

# --- MIGRATIONS (GOOSE) ---
.PHONY: db-create
db-create:
	@echo "Creating migration file..."
	@if [ -z "$(name)" ]; then echo "Error: name is required (make db-create name=something)"; exit 1; fi
	goose -dir $(MIGRATION_DIR) create $(name) sql

.PHONY: db-up
db-up:
	@echo "Migrating UP..."
	goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" up

.PHONY: db-down
db-down:
	@echo "Migrating DOWN (Rollback)..."
	goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down

.PHONY: db-status
db-status:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" status

.PHONY: db-reset
db-reset:
	@echo "Resetting Database..."
	goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down-to 0
	goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" up

# --- DOCKER (Optional - Nếu bạn dùng docker-compose cho DB) ---
.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down