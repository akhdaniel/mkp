# Tests Specification

This is the tests coverage details for the spec detailed in @.agent-os/specs/2025-08-28-database-schema-setup/spec.md

> Created: 2025-08-28
> Version: 1.0.0

## Test Coverage

### Schema Creation Tests

**Database Initialization**
- Verify database can be created with required extensions
- Ensure all required PostgreSQL extensions are available
- Validate schema creation for public, audit schemas
- Test dynamic tenant schema creation

**Table Creation**
- Verify all tables are created with correct columns
- Ensure all data types match specifications
- Validate default values are properly set
- Test partition creation for audit logs

### Constraint Tests

**Primary Key Constraints**
- Verify UUID generation for all primary keys
- Test uniqueness of primary keys
- Ensure primary keys cannot be null

**Foreign Key Constraints**
- Test referential integrity for all relationships
- Verify cascade rules work correctly
- Test restrict and nullify behaviors
- Ensure orphaned records cannot exist

**Check Constraints**
- Test enum value constraints (status fields)
- Verify time constraint (departure < arrival)
- Test positive amount constraints
- Validate capacity constraints (booked <= available)

**Unique Constraints**
- Test unique constraints on code fields
- Verify composite unique constraints
- Test unique constraints with soft deletes

### Data Integrity Tests

**Soft Delete Functionality**
- Verify soft delete timestamps are set correctly
- Test that soft deleted records are excluded from queries
- Ensure foreign key references handle soft deletes
- Test undelete functionality

**Audit Logging**
- Verify audit logs are created for inserts
- Test audit logs capture updates with old/new data
- Ensure deletes are logged properly
- Test audit log partitioning by month

**Timestamp Management**
- Test created_at is set on insert
- Verify updated_at changes on update
- Test timezone handling for all timestamps
- Ensure timestamp triggers work correctly

### Transaction Tests

**Booking Transactions**
- Test atomic booking creation with passengers and tickets
- Verify rollback on payment failure
- Test concurrent booking prevention
- Ensure seat availability updates are atomic

**Payment Processing**
- Test payment status transitions
- Verify refund amount validation
- Test partial refund scenarios
- Ensure payment audit trail integrity

### Multi-Tenant Tests

**Tenant Isolation**
- Test data isolation between tenant schemas
- Verify cross-tenant queries are prevented
- Test tenant-specific table creation
- Ensure shared reference data is accessible

**Schema Management**
- Test dynamic schema creation for new tenants
- Verify schema naming conventions
- Test schema deletion and cleanup
- Ensure migrations apply to all tenant schemas

### Performance Tests

**Index Effectiveness**
- Test query performance with indexes
- Verify index usage in query plans
- Test partial index functionality
- Measure full-text search performance

**Partitioning**
- Test partition pruning for date ranges
- Verify partition creation automation
- Test archival of old partitions
- Measure query performance on partitioned tables

### Security Tests

**Encryption**
- Test PII field encryption/decryption
- Verify encrypted data storage format
- Test key rotation procedures
- Ensure encryption doesn't break queries

**Access Control**
- Test row-level security policies
- Verify role-based permissions
- Test database user privileges
- Ensure application roles map to database roles

### Edge Case Tests

**Boundary Conditions**
- Test maximum capacity bookings
- Verify handling of past dates
- Test extreme price values
- Handle null and empty values

**Concurrency**
- Test simultaneous bookings for same seat
- Verify optimistic locking mechanisms
- Test deadlock prevention
- Ensure transaction isolation levels

### Integration Tests

**Application Integration**
- Test ORM mapping to database schema
- Verify connection pooling behavior
- Test transaction management from application
- Ensure proper error handling

**Migration Tests**
- Test forward migration execution
- Verify rollback functionality
- Test idempotent migrations
- Ensure data preservation during migrations

## Mocking Requirements

### External Services

**UUID Generation**
- Mock UUID generation for predictable testing
- Strategy: Use deterministic UUID generator in tests

**Timestamp Generation**
- Mock CURRENT_TIMESTAMP for time-based tests
- Strategy: Use fixed timestamps in test transactions

**Encryption Keys**
- Mock encryption key service for testing
- Strategy: Use test-only encryption keys

### Test Data

**Reference Data**
- Pre-populate currencies, countries for tests
- Strategy: Use SQL fixtures for reference data

**Tenant Data**
- Create test tenants with known schemas
- Strategy: Setup/teardown test tenant schemas

**Sample Bookings**
- Generate realistic booking scenarios
- Strategy: Use factory patterns for test data

## Test Execution Strategy

1. **Unit Tests**: Run in isolated test database
2. **Integration Tests**: Use Docker PostgreSQL instance
3. **Performance Tests**: Run against production-like data volume
4. **Security Tests**: Run with production-like permissions
5. **Cleanup**: Ensure all test data is removed after runs