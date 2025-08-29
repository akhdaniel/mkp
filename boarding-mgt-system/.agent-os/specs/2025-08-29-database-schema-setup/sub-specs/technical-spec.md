# Technical Specification

This is the technical specification for the spec detailed in @.agent-os/specs/2025-08-29-database-schema-setup/spec.md

> Created: 2025-08-29
> Version: 1.0.0

## Technical Requirements

- PostgreSQL 17+ database with UUID primary keys for all entities
- Multi-tenant architecture with proper data isolation using operator_id foreign keys
- ACID compliance for all financial transactions (bookings, payments, refunds)
- Support for concurrent seat allocation with row-level locking
- Audit trail implementation using triggers or application-level logging
- Optimized indexes for high-frequency queries (availability searches, booking lookups)
- Foreign key constraints with proper cascade rules for data integrity
- JSON/JSONB fields for flexible configuration data (seat maps, vessel specifications)
- Timestamp tracking (created_at, updated_at) on all entities
- Soft delete capability using deleted_at timestamps where appropriate
- Support for database migrations with rollback capability
- Performance optimization for booking search queries across date ranges
- Integration with Go GORM or sqlx for ORM functionality

## Approach

**Selected Approach:** Normalized relational design with strategic denormalization

The database will use a primarily normalized approach to ensure data integrity while strategically denormalizing certain frequently-accessed data for performance. Core entities will be properly normalized with clear relationships, while read-heavy operations like availability searches will be optimized through materialized views or computed columns.

**Multi-tenancy Strategy:** Shared database with operator-level isolation
- Single PostgreSQL instance with operator_id as a tenant identifier
- Row-level security (RLS) policies to enforce data isolation
- Shared reference data (ports, routes) with operator-specific associations

**Concurrency Handling:**
- Optimistic locking for seat allocation using version columns
- Database-level constraints to prevent double-booking
- Transaction isolation levels appropriate for financial operations

**Performance Strategy:**
- Composite indexes on frequently queried columns (operator_id + date ranges)
- Partial indexes for active records (WHERE deleted_at IS NULL)
- Database connection pooling configuration
- Query optimization through EXPLAIN ANALYZE monitoring

## External Dependencies

- **golang-migrate/migrate** - Database migration management
  - **Justification:** Industry standard for Go applications, supports up/down migrations with version tracking

- **google/uuid** - UUID generation for primary keys
  - **Justification:** Provides better distributed system compatibility and avoids sequential ID enumeration attacks

- **PostgreSQL JSONB** - Flexible schema storage for configurations
  - **Justification:** Native PostgreSQL feature for semi-structured data like seat maps and vessel specifications

- **GORM or sqlx** - Go database ORM/toolkit
  - **Justification:** GORM for rapid development with relationships, sqlx for performance-critical queries