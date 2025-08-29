-- Drop functions
DROP FUNCTION IF EXISTS update_last_login(UUID);
DROP FUNCTION IF EXISTS deactivate_user_sessions(UUID);
DROP FUNCTION IF EXISTS clean_expired_sessions();

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;