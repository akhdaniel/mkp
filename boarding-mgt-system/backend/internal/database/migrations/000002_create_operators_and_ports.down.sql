-- Drop triggers
DROP TRIGGER IF EXISTS update_ports_updated_at ON ports;
DROP TRIGGER IF EXISTS update_operators_updated_at ON operators;

-- Drop function (only if no other tables use it)
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of creation)
DROP TABLE IF EXISTS ports CASCADE;
DROP TABLE IF EXISTS operators CASCADE;