.PHONY: help dev-infra dev-infra-down dev-be dev-fe dev-scheduler dev-worker deploy deploy-down clean

help:
	@echo "Available commands:"
	@echo "  dev-infra       - Start infrastructure services (PostgreSQL, Redis, etc.)"
	@echo "  dev-infra-down  - Stop infrastructure services"
	@echo "  dev-be          - Start backend server locally"
	@echo "  dev-fe          - Start frontend server locally"
	@echo "  dev-scheduler   - Start scheduler server locally"
	@echo "  dev-worker      - Start worker server locally"
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

# Cleanup
clean:
	docker-compose -f docker-compose.dev.yml down -v
	docker-compose -f docker-compose.prod.yml down -v
	docker system prune -f