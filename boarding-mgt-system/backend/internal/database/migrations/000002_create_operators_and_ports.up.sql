-- Create operators table (multi-tenant root entity)
CREATE TABLE operators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    address TEXT,
    is_active BOOLEAN DEFAULT true,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on operators for faster lookups
CREATE INDEX idx_operators_code ON operators(code);
CREATE INDEX idx_operators_is_active ON operators(is_active);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_operators_updated_at BEFORE UPDATE ON operators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create ports/terminals table
CREATE TABLE ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    coordinates POINT, -- PostGIS point type for lat/long
    facilities JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on ports
CREATE INDEX idx_ports_code ON ports(code);
CREATE INDEX idx_ports_city_country ON ports(city, country);
CREATE INDEX idx_ports_is_active ON ports(is_active);
-- Spatial index for coordinate-based queries
CREATE INDEX idx_ports_coordinates ON ports USING GIST(coordinates) WHERE coordinates IS NOT NULL;

-- Create trigger for ports updated_at
CREATE TRIGGER update_ports_updated_at BEFORE UPDATE ON ports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE operators IS 'Ferry operators - root entity for multi-tenant system';
COMMENT ON COLUMN operators.code IS 'Unique operator code for identification';
COMMENT ON COLUMN operators.settings IS 'JSON configuration for operator-specific settings';

COMMENT ON TABLE ports IS 'Ports and terminals where ferries operate';
COMMENT ON COLUMN ports.code IS 'IATA-style port code';
COMMENT ON COLUMN ports.coordinates IS 'Geographic coordinates (latitude, longitude)';
COMMENT ON COLUMN ports.facilities IS 'Available facilities at the port (parking, restaurants, etc.)';