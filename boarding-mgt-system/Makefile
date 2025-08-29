.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: docker-up
docker-up: ## Start Docker containers
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5

.PHONY: docker-down
docker-down: ## Stop Docker containers
	docker-compose down

.PHONY: docker-clean
docker-clean: ## Stop and remove Docker containers and volumes
	docker-compose down -v

.PHONY: db-create
db-create: ## Create databases
	docker-compose up -d
	@sleep 5
	@echo "Databases created in Docker containers"

.PHONY: deps
deps: ## Download Go dependencies
	cd backend && go mod download
	cd backend && go mod tidy

.PHONY: test
test: ## Run tests
	cd backend && go test -v -cover ./internal/database/...

.PHONY: test-db
test-db: docker-up ## Run database tests with Docker
	cd backend && go test -v -cover ./internal/database/...

.PHONY: migrate-up
migrate-up: ## Run database migrations
	cd backend && go run cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	cd backend && go run cmd/migrate/main.go down

.PHONY: setup
setup: docker-up deps ## Full setup: start Docker and install dependencies
	@echo "Setup complete!"
	@echo "Run 'make test-db' to verify everything is working"

.PHONY: clean
clean: docker-clean ## Clean everything
	@echo "Cleanup complete!"