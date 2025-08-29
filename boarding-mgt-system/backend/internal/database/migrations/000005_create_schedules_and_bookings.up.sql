-- Create schedules table
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    vessel_id UUID NOT NULL REFERENCES vessels(id) ON DELETE CASCADE,
    departure_date DATE NOT NULL,
    departure_time TIME NOT NULL,
    arrival_time TIME NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    total_capacity INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled',
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    version INTEGER DEFAULT 1, -- for optimistic locking
    CONSTRAINT schedules_base_price_check CHECK (base_price >= 0),
    CONSTRAINT valid_capacity CHECK (available_seats <= total_capacity AND available_seats >= 0),
    CONSTRAINT valid_schedule_status CHECK (status IN ('scheduled', 'boarding', 'departed', 'arrived', 'cancelled'))
);

-- Create indexes on schedules
CREATE INDEX idx_schedules_operator_id ON schedules(operator_id);
CREATE INDEX idx_schedules_route_id ON schedules(route_id);
CREATE INDEX idx_schedules_vessel_id ON schedules(vessel_id);
CREATE INDEX idx_schedules_departure_date ON schedules(departure_date);
CREATE INDEX idx_schedules_departure_datetime ON schedules(departure_date, departure_time);
CREATE INDEX idx_schedules_status ON schedules(status) WHERE status != 'arrived';
-- Composite index for finding available schedules
CREATE INDEX idx_schedules_availability ON schedules(route_id, departure_date, status) 
    WHERE status = 'scheduled' AND available_seats > 0;

-- Create trigger for schedules updated_at
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create bookings table
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_reference VARCHAR(20) UNIQUE NOT NULL,
    schedule_id UUID NOT NULL REFERENCES schedules(id),
    customer_id UUID NOT NULL REFERENCES users(id),
    passenger_count INTEGER NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    booking_status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'pending',
    booking_channel VARCHAR(20) NOT NULL,
    special_requirements TEXT,
    booking_agent_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT bookings_passenger_count_check CHECK (passenger_count > 0),
    CONSTRAINT bookings_total_amount_check CHECK (total_amount >= 0),
    CONSTRAINT valid_booking_status CHECK (booking_status IN ('pending', 'confirmed', 'cancelled', 'refunded')),
    CONSTRAINT valid_payment_status CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded')),
    CONSTRAINT valid_booking_channel CHECK (booking_channel IN ('online', 'pos', 'phone', 'agent'))
);

-- Create indexes on bookings
CREATE INDEX idx_bookings_booking_reference ON bookings(booking_reference);
CREATE INDEX idx_bookings_schedule_id ON bookings(schedule_id);
CREATE INDEX idx_bookings_customer_id ON bookings(customer_id);
CREATE INDEX idx_bookings_booking_status ON bookings(booking_status);
CREATE INDEX idx_bookings_payment_status ON bookings(payment_status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at DESC);
CREATE INDEX idx_bookings_agent_id ON bookings(booking_agent_id) WHERE booking_agent_id IS NOT NULL;
-- Composite index for customer booking history
CREATE INDEX idx_bookings_customer_history ON bookings(customer_id, created_at DESC);

-- Create trigger for bookings updated_at
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to generate booking reference
CREATE OR REPLACE FUNCTION generate_booking_reference()
RETURNS VARCHAR(20) AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR(20);
    i INTEGER;
BEGIN
    -- Format: BK + timestamp (6 chars) + random (4 chars)
    result := 'BK' || to_char(CURRENT_TIMESTAMP, 'YYMMDD');
    
    -- Add random characters
    FOR i IN 1..4 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create function to update available seats after booking
CREATE OR REPLACE FUNCTION update_schedule_availability()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Decrease available seats
        UPDATE schedules 
        SET available_seats = available_seats - NEW.passenger_count,
            version = version + 1
        WHERE id = NEW.schedule_id
        AND available_seats >= NEW.passenger_count;
        
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Insufficient seats available for booking';
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- Increase available seats when booking is cancelled
        IF OLD.booking_status = 'confirmed' THEN
            UPDATE schedules 
            SET available_seats = available_seats + OLD.passenger_count,
                version = version + 1
            WHERE id = OLD.schedule_id;
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Handle booking status changes
        IF OLD.booking_status != 'cancelled' AND NEW.booking_status = 'cancelled' THEN
            -- Booking cancelled, return seats
            UPDATE schedules 
            SET available_seats = available_seats + NEW.passenger_count,
                version = version + 1
            WHERE id = NEW.schedule_id;
        ELSIF OLD.booking_status = 'cancelled' AND NEW.booking_status = 'confirmed' THEN
            -- Booking restored, decrease seats
            UPDATE schedules 
            SET available_seats = available_seats - NEW.passenger_count,
                version = version + 1
            WHERE id = NEW.schedule_id
            AND available_seats >= NEW.passenger_count;
            
            IF NOT FOUND THEN
                RAISE EXCEPTION 'Insufficient seats available for booking restoration';
            END IF;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for automatic seat management
CREATE TRIGGER manage_schedule_availability
    AFTER INSERT OR UPDATE OR DELETE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_schedule_availability();

-- Add comments for documentation
COMMENT ON TABLE schedules IS 'Ferry schedules with departure times and capacity';
COMMENT ON COLUMN schedules.version IS 'Version number for optimistic locking to prevent concurrent booking conflicts';
COMMENT ON COLUMN schedules.available_seats IS 'Current number of seats available for booking';
COMMENT ON COLUMN schedules.status IS 'Current status of the schedule: scheduled, boarding, departed, arrived, or cancelled';

COMMENT ON TABLE bookings IS 'Customer bookings for ferry schedules';
COMMENT ON COLUMN bookings.booking_reference IS 'Unique reference number for the booking';
COMMENT ON COLUMN bookings.booking_channel IS 'Channel through which booking was made: online, pos, phone, or agent';
COMMENT ON COLUMN bookings.booking_agent_id IS 'Reference to agent who made the booking (NULL for customer self-service)';

COMMENT ON FUNCTION generate_booking_reference() IS 'Generates unique booking reference with format BK + date + random';
COMMENT ON FUNCTION update_schedule_availability() IS 'Automatically manages seat availability when bookings are created, updated, or cancelled';