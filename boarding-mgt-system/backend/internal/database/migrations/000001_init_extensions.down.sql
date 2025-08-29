-- Drop schemas
DROP SCHEMA IF EXISTS audit CASCADE;

-- Note: We don't drop extensions as they might be used by other databases
-- Extensions are database-level, not schema-level