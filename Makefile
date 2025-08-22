.PHONY: help dev-infra dev-infra-down dev-be dev-fe dev-scheduler dev-worker deploy deploy-down clean build-all mod-tidy test health-check

help:
	@echo "Available commands:"
	@echo "  dev-infra       - Start infrastructure services (PostgreSQL, Redis, etc.)"
	@echo "  dev-infra-down  - Stop infrastructure services"
	@echo "  dev-be          - Start backend server locally"
	@echo "  dev-fe          - Start frontend server locally"
	@echo "  dev-scheduler   - Start scheduler server locally"
	@echo "  dev-worker      - Start worker server locally"
	@echo "  build-all       - Build all Go services"
	@echo "  mod-tidy        - Tidy all Go modules"
	@echo "  test            - Run all tests"
	@echo "  health-check    - Check if all services are healthy"
	@echo "  deploy          - Deploy all services with Docker"
	@echo "  deploy-down     - Stop all deployed services"
	@echo "  clean           - Clean up Docker containers and volumes"

# Local development - infrastructure only
dev-infra:
	docker-compose -f docker-compose.dev.yml up -d

dev-infra-down:
	docker-compose -f docker-compose.dev.yml down

# Local development - servers
dev-be:
	cd weave-be && go run main.go

dev-fe:
	cd weave-fe && npm start

dev-scheduler:
	cd weave-scheduler && go run main.go

dev-worker:
	cd weave-worker && go run main.go

# Production deployment
deploy:
	docker-compose -f docker-compose.prod.yml up -d --build

deploy-down:
	docker-compose -f docker-compose.prod.yml down

# Build targets
build-all: build-module build-be build-scheduler build-worker

build-module:
	cd weave-module && go mod tidy

build-be:
	cd weave-be && go mod tidy && go build -o bin/weave-be .

build-scheduler:
	cd weave-scheduler && go mod tidy && go build -o bin/weave-scheduler .

build-worker:
	cd weave-worker && go mod tidy && go build -o bin/weave-worker .

# Module management
mod-tidy:
	cd weave-module && go mod tidy
	cd weave-be && go mod tidy
	cd weave-scheduler && go mod tidy
	cd weave-worker && go mod tidy

# Testing
test:
	cd weave-module && go test ./...
	cd weave-be && go test ./...
	cd weave-scheduler && go test ./...
	cd weave-worker && go test ./...
	cd weave-fe && npm test --passWithNoTests

# Health check
health-check:
	@echo "Checking infrastructure services..."
	@docker-compose -f docker-compose.dev.yml ps
	@echo "Checking Go modules..."
	@cd weave-module && go mod verify
	@cd weave-be && go mod verify
	@cd weave-scheduler && go mod verify
	@cd weave-worker && go mod verify
	@echo "Checking frontend dependencies..."
	@cd weave-fe && npm list --depth=0 || true

# Cleanup
clean:
	docker-compose -f docker-compose.dev.yml down -v
	docker-compose -f docker-compose.prod.yml down -v
	docker system prune -f
	cd weave-be && rm -rf bin/
	cd weave-scheduler && rm -rf bin/
	cd weave-worker && rm -rf bin/
	cd weave-fe && rm -rf build/ node_modules/.cache/