# Technical Specification

This is the technical specification for the spec detailed in @.agent-os/specs/2025-08-28-database-schema-setup/spec.md

> Created: 2025-08-28
> Version: 1.0.0

## Technical Requirements

- PostgreSQL 17+ database server with full SQL support
- Support for UUID primary keys for distributed system compatibility
- Timestamp columns with timezone support for all temporal data
- JSON/JSONB columns for flexible metadata storage
- Full-text search capability for customer service features
- Proper indexing strategy for high-performance queries
- Foreign key constraints with appropriate cascade rules
- Check constraints for data validation at database level
- Trigger functions for audit logging and data consistency
- Database-level encryption for sensitive data fields

## Database Design Principles

### Multi-Tenant Architecture

**Option A:** Shared schema with tenant_id columns
- Pros: Easier maintenance, shared resources, simpler deployment
- Cons: Risk of data leakage, complex row-level security

**Option B:** Schema per tenant (Selected)
- Pros: Complete data isolation, easier compliance, independent scaling
- Cons: More complex deployment, higher resource usage

**Rationale:** Ferry operators require complete data isolation for compliance and security. Schema-per-tenant provides the strongest isolation while allowing shared infrastructure.

### Primary Key Strategy

**Option A:** Auto-incrementing integers
- Pros: Smaller storage, faster joins, simpler
- Cons: Predictable IDs, issues with distributed systems

**Option B:** UUIDs (Selected)
- Pros: Globally unique, better for distributed systems, no enumeration attacks
- Cons: Larger storage, slightly slower joins

**Rationale:** UUIDs provide better security and prepare the system for future distributed architecture needs.

## Schema Organization

### Public Schema
- Shared configuration tables
- System-wide lookups (countries, currencies, timezones)
- Tenant registry

### Tenant Schemas (Dynamic)
- All operator-specific tables
- Isolated data per ferry operator
- Independent migrations possible

### Audit Schema
- Centralized audit logs across all tenants
- System event tracking
- Compliance reporting data

## Key Technical Decisions

### Soft Deletes
- Implement deleted_at timestamps instead of hard deletes
- Maintain data history for audit compliance
- Use partial indexes to exclude deleted records

### Temporal Data
- Use timestamptz for all timestamps
- Store schedules in UTC, display in local timezone
- Implement validity periods for prices and schedules

### Financial Precision
- Use DECIMAL(19,4) for all monetary values
- Store amounts in smallest currency unit (cents)
- Track currency per transaction

## Performance Considerations

### Indexing Strategy
- B-tree indexes on foreign keys and frequently queried columns
- Partial indexes for soft-deleted records
- GiST indexes for scheduling overlap queries
- GIN indexes for full-text search on help content

### Partitioning Plan
- Partition bookings table by created_at (monthly)
- Partition audit_logs by created_at (monthly)
- Archive old partitions after 2 years

## Security Measures

### Data Encryption
- Encrypt PII columns using pgcrypto
- Separate encryption keys per tenant
- Transparent encryption for application layer

### Access Control
- Row Level Security (RLS) policies per tenant
- Database roles matching application roles
- Separate read-only role for reporting

## External Dependencies

- **pgcrypto extension** - For encryption functions
- **Justification:** Required for PII protection and compliance

- **uuid-ossp extension** - For UUID generation
- **Justification:** Needed for distributed-safe primary keys

- **pg_trgm extension** - For fuzzy text search
- **Justification:** Enables better search in helpdesk and FAQs

- **btree_gist extension** - For exclusion constraints
- **Justification:** Prevents scheduling conflicts and double-booking