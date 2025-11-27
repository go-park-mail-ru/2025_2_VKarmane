# Makefile for VKarmane API

MODULE := github.com/go-park-mail-ru/2025_2_VKarmane

# Все пакеты проекта
ALL_PKGS := $(shell go list ./...)

# Пакеты для покрытия без моков, cmd/docs/proto, сервисных entrypoint-ов, HTTP-обвязок и интеграционных адаптеров.
PKGS := $(shell go list ./... \
	| grep -vE '(^|/)(mocks|cmd|docs|proto)(/|$$)' \
	| grep -vE 'internal/app/[^/]+$$' \
	| grep -vE 'internal/app/.+/handlers' \
	| grep -vE 'internal/app/image/repository' \
	| grep -vE 'internal/metrics$$' \
	| grep -vE 'internal/models$$' \
	| grep -vE 'scripts$$')

COVER_MODE := atomic
COVER_RAW := coverage.raw.out
COVER_OUT := coverage.out
COVER_HTML := coverage.html
COVER_THRESHOLD := 60

EXCLUDE_FILES_REGEX := \/mocks\/|\/mock_.*\.go

.PHONY: help build up down logs clean test migrate swagger cover cover-check coverhtml dev deploy mocks seed-users

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - Show logs from all services"
	@echo "  clean     - Remove containers and volumes"
	@echo "  test      - Run tests"
	@echo "  cover     - Run tests with coverage"
	@echo "  coverhtml - Generate HTML coverage report"
	@echo "  migrate   - Apply database migrations"
	@echo "  dev       - Start development environment"
	@echo "  swagger   - Generate Swagger documentation"
	@echo "  mocks     - Generate mocks using gomock"
	@echo "  deploy    - Production deployment"
	@echo "  seed-users - Seed test users with accounts"

# Build Docker images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Show logs from all services
logs:
	docker-compose logs -f

# Remove containers and volumes
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f
	rm -f $(COVER_OUT) $(COVER_HTML) $(COVER_RAW)

# Run tests
test:
	@go test $(PKGS)

# Run tests with coverage
cover:
	@echo "Running tests with coverage..."
	@GOFLAGS= go test -covermode=$(COVER_MODE) -coverprofile=$(COVER_OUT) $(PKGS) || true
	@if [ -f $(COVER_OUT) ]; then \
		go tool cover -func=$(COVER_OUT) | grep total:; \
	else \
		echo "coverage.out not found"; \
	fi

cover-check: cover
	@go run ./scripts/coverage_check.go -profile=$(COVER_OUT) -threshold=$(COVER_THRESHOLD)

# Generate HTML coverage report
coverhtml: cover
	@go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "Wrote $(COVER_HTML)"
	@which xdg-open >/dev/null 2>&1 && xdg-open $(COVER_HTML) || true
	@which open >/dev/null 2>&1 && open $(COVER_HTML) || true

# Apply database migrations
migrate:
	docker-compose exec api go run cmd/api/main.go migrate

# Start development environment
dev: build up
	@echo "Development environment started!"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/"
	@echo "PostgreSQL: localhost:5432"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && \
	if ! command -v swag >/dev/null 2>&1; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi && \
	export PATH=$$PATH:$$(go env GOPATH)/bin && \
	swag init -g cmd/api/main.go -o docs || \
	(docker-compose exec api swag init -g cmd/api/main.go -o docs 2>/dev/null || echo "Please install swag: go install github.com/swaggo/swag/cmd/swag@latest")

# Generate mocks using gomock
mocks:
	@echo "Generating mocks..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && \
	if ! command -v mockgen >/dev/null 2>&1; then \
		echo "Installing mockgen..."; \
		go install go.uber.org/mock/mockgen@latest; \
	fi
	@echo "Running go generate..."
	@go generate ./internal/mocks/...
	@echo "Mocks generated successfully!"

# Production deployment
deploy: build up
	@echo "Production deployment completed!"

# Seed test users with accounts (local)
seed-users:
	@echo "Seeding test users with accounts..."
	@go run scripts/seed_test_users.go
	@echo "Seed completed!"
	
