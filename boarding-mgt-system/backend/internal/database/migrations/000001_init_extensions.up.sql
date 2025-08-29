-- Enable required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS audit;

-- Add comment for documentation
COMMENT ON SCHEMA public IS 'Shared system tables and tenant registry';
COMMENT ON SCHEMA audit IS 'Centralized audit logging across all tenants';