# Spec Tasks

These are the tasks to be completed for the spec detailed in @.agent-os/specs/2025-08-28-database-schema-setup/spec.md

> Created: 2025-08-28
> Status: Ready for Implementation

## Tasks

- [ ] 1. Setup PostgreSQL database and create initial structure
  - [ ] 1.1 Write tests for database connection and extension availability
  - [ ] 1.2 Install PostgreSQL locally or configure Docker container
  - [ ] 1.3 Create database and enable required extensions (pgcrypto, uuid-ossp, pg_trgm, btree_gist)
  - [ ] 1.4 Create database connection configuration and environment variables
  - [ ] 1.5 Create initial migration structure and tooling setup
  - [ ] 1.6 Verify all tests pass

- [ ] 2. Implement public schema and core reference tables
  - [ ] 2.1 Write tests for public schema tables creation
  - [ ] 2.2 Create migration for tenants table with multi-tenant support
  - [ ] 2.3 Create system_config, currencies, and countries tables
  - [ ] 2.4 Insert initial reference data (currencies, countries)
  - [ ] 2.5 Add indexes for reference table lookups
  - [ ] 2.6 Verify all tests pass

- [ ] 3. Build tenant schema structure and user management
  - [ ] 3.1 Write tests for tenant schema creation and user tables
  - [ ] 3.2 Create dynamic tenant schema creation function
  - [ ] 3.3 Implement users and auth_tokens tables with security
  - [ ] 3.4 Create operators and ports management tables
  - [ ] 3.5 Add authentication indexes and constraints
  - [ ] 3.6 Implement password hashing and token management functions
  - [ ] 3.7 Verify all tests pass

- [ ] 4. Create booking system core tables
  - [ ] 4.1 Write tests for booking-related tables and constraints
  - [ ] 4.2 Create vessels, routes, and schedules tables
  - [ ] 4.3 Implement schedule_instances with availability tracking
  - [ ] 4.4 Create bookings, booking_passengers, and tickets tables
  - [ ] 4.5 Implement payment_transactions and refunds tables
  - [ ] 4.6 Add seat_classes and pricing_rules tables
  - [ ] 4.7 Create all booking-related indexes and constraints
  - [ ] 4.8 Verify all tests pass

- [ ] 5. Implement support system and audit functionality
  - [ ] 5.1 Write tests for helpdesk and audit logging
  - [ ] 5.2 Create helpdesk_tickets and helpdesk_messages tables
  - [ ] 5.3 Implement FAQ categories and items tables
  - [ ] 5.4 Create audit schema with partitioned audit_logs table
  - [ ] 5.5 Implement audit trigger functions
  - [ ] 5.6 Apply audit triggers to all relevant tables
  - [ ] 5.7 Create monthly partition management for audit logs
  - [ ] 5.8 Verify all tests pass

## Spec Documentation

- Tasks: @.agent-os/specs/2025-08-28-database-schema-setup/tasks.md
- Technical Specification: @.agent-os/specs/2025-08-28-database-schema-setup/sub-specs/technical-spec.md
- Database Schema: @.agent-os/specs/2025-08-28-database-schema-setup/sub-specs/database-schema.md
- Tests Specification: @.agent-os/specs/2025-08-28-database-schema-setup/sub-specs/tests.md