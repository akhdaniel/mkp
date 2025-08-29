-- Drop triggers
DROP TRIGGER IF EXISTS update_routes_updated_at ON routes;
DROP TRIGGER IF EXISTS update_vessels_updated_at ON vessels;

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS routes CASCADE;
DROP TABLE IF EXISTS vessels CASCADE;