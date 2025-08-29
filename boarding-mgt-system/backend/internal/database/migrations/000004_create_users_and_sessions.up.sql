-- Create users table (customers, agents, operators)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    date_of_birth DATE,
    nationality VARCHAR(3), -- ISO country code
    user_type VARCHAR(20) NOT NULL,
    operator_id UUID REFERENCES operators(id) ON DELETE CASCADE,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_user_type CHECK (user_type IN ('customer', 'agent', 'operator_admin', 'system_admin')),
    CONSTRAINT valid_nationality CHECK (nationality IS NULL OR LENGTH(nationality) = 3),
    CONSTRAINT operator_required_for_staff CHECK (
        (user_type IN ('agent', 'operator_admin') AND operator_id IS NOT NULL) OR
        (user_type IN ('customer', 'system_admin'))
    )
);

-- Create indexes on users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_user_type ON users(user_type);
CREATE INDEX idx_users_operator_id ON users(operator_id) WHERE operator_id IS NOT NULL;
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
-- Composite index for common queries
CREATE INDEX idx_users_type_operator ON users(user_type, operator_id) WHERE operator_id IS NOT NULL;

-- Create trigger for users updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create user_sessions table for JWT management
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on user_sessions
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token_hash ON user_sessions(token_hash) WHERE is_active = true;
CREATE INDEX idx_user_sessions_refresh_token_hash ON user_sessions(refresh_token_hash) WHERE is_active = true;
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at) WHERE is_active = true;
CREATE INDEX idx_user_sessions_is_active ON user_sessions(is_active);

-- Create function to clean expired sessions (can be called by a scheduled job)
CREATE OR REPLACE FUNCTION clean_expired_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM user_sessions
    WHERE expires_at < CURRENT_TIMESTAMP
    OR (is_active = false AND created_at < CURRENT_TIMESTAMP - INTERVAL '30 days');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create function to deactivate all sessions for a user (useful for password reset)
CREATE OR REPLACE FUNCTION deactivate_user_sessions(p_user_id UUID)
RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER;
BEGIN
    UPDATE user_sessions
    SET is_active = false
    WHERE user_id = p_user_id AND is_active = true;
    
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

-- Create function to update last login timestamp
CREATE OR REPLACE FUNCTION update_last_login(p_user_id UUID)
RETURNS VOID AS $$
BEGIN
    UPDATE users
    SET last_login_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE users IS 'System users including customers, agents, and administrators';
COMMENT ON COLUMN users.email IS 'Unique email address used for authentication';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.user_type IS 'User role: customer, agent, operator_admin, or system_admin';
COMMENT ON COLUMN users.operator_id IS 'Reference to operator for staff users, NULL for customers';
COMMENT ON COLUMN users.is_verified IS 'Whether email has been verified';
COMMENT ON COLUMN users.nationality IS 'ISO 3166-1 alpha-3 country code';

COMMENT ON TABLE user_sessions IS 'Active user sessions for JWT token management';
COMMENT ON COLUMN user_sessions.token_hash IS 'SHA256 hash of the JWT access token';
COMMENT ON COLUMN user_sessions.refresh_token_hash IS 'SHA256 hash of the refresh token';
COMMENT ON COLUMN user_sessions.expires_at IS 'Token expiration timestamp';
COMMENT ON COLUMN user_sessions.ip_address IS 'IP address from which session was created';
COMMENT ON COLUMN user_sessions.user_agent IS 'Browser/client user agent string';

COMMENT ON FUNCTION clean_expired_sessions() IS 'Removes expired and old inactive sessions';
COMMENT ON FUNCTION deactivate_user_sessions(UUID) IS 'Deactivates all active sessions for a user';
COMMENT ON FUNCTION update_last_login(UUID) IS 'Updates the last login timestamp for a user';