#!/bin/bash

# Test script for database migrations
# This script tests forward and backward migrations to ensure they work correctly

set -e  # Exit on error

echo "========================================="
echo "Testing Database Migrations"
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Change to backend directory
cd "$(dirname "$0")/.."

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.22+ first."
    exit 1
fi

print_info "Starting migration tests..."

# 1. Check current version
print_info "Checking current migration version..."
go run cmd/migrate/main.go version || true

# 2. Rollback all migrations first (clean slate)
print_info "Rolling back all existing migrations..."
go run cmd/migrate/main.go down || true

# 3. Run migrations up one by one
print_info "Testing migration 001 (extensions)..."
go run cmd/migrate/main.go up -steps 1
if [ $? -eq 0 ]; then
    print_status "Migration 001 applied successfully"
else
    print_error "Migration 001 failed"
    exit 1
fi

print_info "Testing migration 002 (operators and ports)..."
go run cmd/migrate/main.go up -steps 1
if [ $? -eq 0 ]; then
    print_status "Migration 002 applied successfully"
else
    print_error "Migration 002 failed"
    exit 1
fi

print_info "Testing migration 003 (vessels and routes)..."
go run cmd/migrate/main.go up -steps 1
if [ $? -eq 0 ]; then
    print_status "Migration 003 applied successfully"
else
    print_error "Migration 003 failed"
    exit 1
fi

# 4. Check final version
print_info "Verifying all migrations are applied..."
VERSION=$(go run cmd/migrate/main.go version 2>&1)
echo "Current version: $VERSION"

# 5. Test rollback of migration 003
print_info "Testing rollback of migration 003..."
go run cmd/migrate/main.go down -steps 1
if [ $? -eq 0 ]; then
    print_status "Migration 003 rolled back successfully"
else
    print_error "Migration 003 rollback failed"
    exit 1
fi

# 6. Re-apply migration 003
print_info "Re-applying migration 003..."
go run cmd/migrate/main.go up -steps 1
if [ $? -eq 0 ]; then
    print_status "Migration 003 re-applied successfully"
else
    print_error "Migration 003 re-apply failed"
    exit 1
fi

# 7. Test full rollback
print_info "Testing full rollback..."
go run cmd/migrate/main.go down
if [ $? -eq 0 ]; then
    print_status "All migrations rolled back successfully"
else
    print_error "Full rollback failed"
    exit 1
fi

# 8. Test full migration up
print_info "Testing full migration up..."
go run cmd/migrate/main.go up
if [ $? -eq 0 ]; then
    print_status "All migrations applied successfully"
else
    print_error "Full migration up failed"
    exit 1
fi

echo ""
echo "========================================="
echo -e "${GREEN}All migration tests passed!${NC}"
echo "========================================="

# Run the Go tests for core entities
print_info "Running core entity tests..."
go test -v ./internal/database -run TestCoreEntityTables

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================="
    echo -e "${GREEN}All tests completed successfully!${NC}"
    echo "========================================="
else
    echo ""
    echo "========================================="
    echo -e "${RED}Some tests failed!${NC}"
    echo "========================================="
    exit 1
fi