-- Create vessels table
CREATE TABLE vessels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(100) UNIQUE NOT NULL,
    vessel_type VARCHAR(50) NOT NULL, -- 'passenger', 'cargo', 'mixed'
    capacity INTEGER NOT NULL,
    deck_count INTEGER DEFAULT 1,
    seat_configuration JSONB NOT NULL, -- JSON structure of decks/sections/seats
    amenities JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT vessels_capacity_check CHECK (capacity > 0),
    CONSTRAINT valid_vessel_type CHECK (vessel_type IN ('passenger', 'cargo', 'mixed')),
    CONSTRAINT valid_deck_count CHECK (deck_count > 0)
);

-- Create indexes on vessels
CREATE INDEX idx_vessels_operator_id ON vessels(operator_id);
CREATE INDEX idx_vessels_registration_number ON vessels(registration_number);
CREATE INDEX idx_vessels_vessel_type ON vessels(vessel_type);
CREATE INDEX idx_vessels_is_active ON vessels(is_active);

-- Create trigger for vessels updated_at
CREATE TRIGGER update_vessels_updated_at BEFORE UPDATE ON vessels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create routes table
CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    departure_port_id UUID NOT NULL REFERENCES ports(id),
    arrival_port_id UUID NOT NULL REFERENCES ports(id),
    distance_km DECIMAL(8,2),
    estimated_duration INTERVAL NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT different_ports CHECK (departure_port_id != arrival_port_id),
    CONSTRAINT valid_distance CHECK (distance_km IS NULL OR distance_km > 0)
);

-- Create indexes on routes
CREATE INDEX idx_routes_operator_id ON routes(operator_id);
CREATE INDEX idx_routes_departure_port_id ON routes(departure_port_id);
CREATE INDEX idx_routes_arrival_port_id ON routes(arrival_port_id);
CREATE INDEX idx_routes_ports ON routes(departure_port_id, arrival_port_id);
CREATE INDEX idx_routes_is_active ON routes(is_active);

-- Create trigger for routes updated_at
CREATE TRIGGER update_routes_updated_at BEFORE UPDATE ON routes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE vessels IS 'Vessels operated by ferry companies';
COMMENT ON COLUMN vessels.operator_id IS 'Reference to the operator that owns this vessel';
COMMENT ON COLUMN vessels.vessel_type IS 'Type of vessel: passenger, cargo, or mixed';
COMMENT ON COLUMN vessels.capacity IS 'Total passenger capacity of the vessel';
COMMENT ON COLUMN vessels.seat_configuration IS 'JSON structure defining deck layout and seat arrangements';
COMMENT ON COLUMN vessels.amenities IS 'Available amenities on the vessel (cafe, wifi, etc.)';

COMMENT ON TABLE routes IS 'Routes between ports operated by ferry services';
COMMENT ON COLUMN routes.operator_id IS 'Reference to the operator that runs this route';
COMMENT ON COLUMN routes.departure_port_id IS 'Port where the route begins';
COMMENT ON COLUMN routes.arrival_port_id IS 'Port where the route ends';
COMMENT ON COLUMN routes.distance_km IS 'Distance in kilometers between ports';
COMMENT ON COLUMN routes.estimated_duration IS 'Expected travel time for this route';