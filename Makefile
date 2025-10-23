# Makefile for VKarmane API

.PHONY: help build up down logs clean test migrate swagger

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build Docker images"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - Show logs from all services"
	@echo "  clean     - Remove containers and volumes"
	@echo "  test      - Run tests"
	@echo "  migrate   - Apply database migrations"
	@echo "  dev       - Start development environment"
	@echo "  swagger   - Generate Swagger documentation"

# Build Docker images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Clean up containers and volumes
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run tests
test:
	go test ./...

# Apply migrations (for local development)
migrate:
	psql -h localhost -p 5432 -U vkarmane -d vkarmane -f migrations/001_create_tables.sql

# Development environment
dev: build up
	@echo "Development environment started!"
	@echo "API: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "Use 'make logs' to see logs"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	export PATH=$$PATH:$$(go env GOPATH)/bin && swag init -g cmd/api/main.go -o docs/
	@echo "Swagger documentation generated in docs/"

# Production deployment
deploy: build up
	@echo "Production deployment completed!"