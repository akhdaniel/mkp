# Database Migration Framework

## Overview

The database migration framework for FerryFlow is built using `golang-migrate/migrate` library, providing robust schema version management and migration capabilities.

## Architecture

### Components

1. **Migrator** (`migration.go`)
   - Handles all migration operations (up, down, steps, force)
   - Embedded migration files for deployment simplicity
   - Connection pooling and timeout management

2. **Database Connection** (`connection.go`)
   - PostgreSQL connection pool management
   - Extension verification and creation
   - Environment-specific configuration

3. **Migration CLI** (`cmd/migrate/main.go`)
   - Command-line interface for migration operations
   - Supports up, down, version, force, and create commands
   - Environment configuration via .env files

## Setup Instructions

### Prerequisites

1. Install Go 1.22+
2. Install PostgreSQL 17+
3. Install Docker (optional, for containerized databases)

### Initial Setup

1. **Start the databases** (if using Docker):
   ```bash
   make docker-up
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Run migrations**:
   ```bash
   make migrate-up
   ```

## Usage

### Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback one migration
make migrate-down

# Rollback all migrations
make migrate-down-all

# Check current version
make migrate-version

# Create new migration
make migrate-create name=add_users_table
```

### Direct CLI Usage

```bash
cd backend

# Run migrations
go run cmd/migrate/main.go up

# Rollback specific number of steps
go run cmd/migrate/main.go down -steps 2

# Force set version (use with caution)
go run cmd/migrate/main.go force -version 3

# Create new migration files
go run cmd/migrate/main.go create add_new_feature
```

## Migration Files

Migrations are stored in `backend/internal/database/migrations/` with the naming convention:
- `NNNNNN_description.up.sql` - Forward migration
- `NNNNNN_description.down.sql` - Rollback migration

Where `NNNNNN` is a sequential number (or timestamp for created migrations).

### Current Migrations

1. **000001_init_extensions** - Initializes PostgreSQL extensions and schemas
   - Creates: pgcrypto, uuid-ossp, pg_trgm, btree_gist extensions
   - Creates: audit schema for audit logging

## Testing

Run the migration tests:

```bash
make test-db
```

The tests verify:
- Migration framework initialization
- Forward and backward migrations
- Step-based migrations
- Extension creation
- Schema creation

## Environment Configuration

Configure database connections in `.env`:

```env
# Development Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ferryflow_dev
DB_USER=ferryflow
DB_PASSWORD=ferryflow_dev_2024
DB_SSL_MODE=disable

# Test Database
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=ferryflow_test
TEST_DB_USER=ferryflow
TEST_DB_PASSWORD=ferryflow_test_2024
TEST_DB_SSL_MODE=disable
```

## Best Practices

1. **Always test migrations** in development before production
2. **Include rollback migrations** for every forward migration
3. **Use transactions** where possible for atomic changes
4. **Version control** all migration files
5. **Never modify** existing migration files; create new ones instead
6. **Document** complex migrations with comments

## Troubleshooting

### Common Issues

1. **"database does not exist"**
   - Run `make db-create` or create manually
   
2. **"permission denied to create extension"**
   - Ensure user has SUPERUSER privileges or pre-create extensions

3. **"dirty database version"**
   - A migration failed midway
   - Fix the issue and use `force` command to reset version

4. **"no migration files found"**
   - Check that migration files are in correct directory
   - Verify file naming convention

## Next Steps

After setting up the migration framework, proceed with:
1. Creating core entity tables (Task 2)
2. Implementing user management schema (Task 3)
3. Building booking and scheduling system (Task 4)