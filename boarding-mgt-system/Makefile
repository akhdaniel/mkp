.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: up
up: ## Start all services with Docker Compose
	docker-compose up -d

.PHONY: down
down: ## Stop all services
	docker-compose down

.PHONY: clean
clean: ## Stop services and remove volumes
	docker-compose down -v

.PHONY: build
build: ## Build all Docker images
	docker-compose build

.PHONY: rebuild
rebuild: ## Rebuild all Docker images without cache
	docker-compose build --no-cache

.PHONY: logs
logs: ## Show logs from all services
	docker-compose logs -f

.PHONY: logs-backend
logs-backend: ## Show backend logs
	docker-compose logs -f backend

.PHONY: logs-frontend
logs-frontend: ## Show frontend logs
	docker-compose logs -f frontend

.PHONY: logs-db
logs-db: ## Show database logs
	docker-compose logs -f postgres

.PHONY: migrate-up
migrate-up: ## Run database migrations
	cd backend && go run cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	cd backend && go run cmd/migrate/main.go down -steps 1

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	cd backend && go run cmd/migrate/main.go create $(NAME)

.PHONY: seed
seed: ## Seed database with demo data
	cd backend && go run cmd/seed/main.go

.PHONY: test
test: ## Run all tests
	cd backend && go test ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out

.PHONY: dev-backend
dev-backend: ## Run backend in development mode
	cd backend && go run cmd/server/main.go

.PHONY: dev-frontend
dev-frontend: ## Run frontend in development mode
	cd frontend && npm run dev

.PHONY: dev
dev: ## Run both frontend and backend in development mode
	make -j2 dev-backend dev-frontend

.PHONY: shell-backend
shell-backend: ## Open shell in backend container
	docker-compose exec backend /bin/sh

.PHONY: shell-db
shell-db: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U ferryflow -d ferryflow_dev

.PHONY: backup-db
backup-db: ## Backup database
	docker-compose exec postgres pg_dump -U ferryflow ferryflow_dev > backup_$(shell date +%Y%m%d_%H%M%S).sql

.PHONY: status
status: ## Show status of all services
	docker-compose ps

.PHONY: restart
restart: ## Restart all services
	docker-compose restart

.PHONY: restart-backend
restart-backend: ## Restart backend service
	docker-compose restart backend

.PHONY: restart-frontend
restart-frontend: ## Restart frontend service
	docker-compose restart frontend