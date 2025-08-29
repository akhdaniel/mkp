-- Create tickets table (individual passengers)
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    passenger_name VARCHAR(255) NOT NULL,
    passenger_type VARCHAR(20) DEFAULT 'adult',
    seat_number VARCHAR(20),
    ticket_price DECIMAL(10,2) NOT NULL,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    check_in_status VARCHAR(20) DEFAULT 'not_checked_in',
    check_in_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tickets_ticket_price_check CHECK (ticket_price >= 0),
    CONSTRAINT valid_passenger_type CHECK (passenger_type IN ('adult', 'child', 'infant', 'senior')),
    CONSTRAINT valid_check_in_status CHECK (check_in_status IN ('not_checked_in', 'checked_in', 'boarded'))
);

-- Create indexes on tickets
CREATE INDEX idx_tickets_booking_id ON tickets(booking_id);
CREATE INDEX idx_tickets_qr_code ON tickets(qr_code);
CREATE INDEX idx_tickets_passenger_type ON tickets(passenger_type);
CREATE INDEX idx_tickets_check_in_status ON tickets(check_in_status);
CREATE INDEX idx_tickets_seat_number ON tickets(seat_number) WHERE seat_number IS NOT NULL;

-- Create trigger for tickets updated_at
CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    payment_method VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_status VARCHAR(20) DEFAULT 'pending',
    gateway_transaction_id VARCHAR(255),
    gateway_response JSONB,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT payments_amount_check CHECK (amount > 0),
    CONSTRAINT valid_payment_method CHECK (payment_method IN ('credit_card', 'debit_card', 'cash', 'bank_transfer', 'mobile_money', 'paypal')),
    CONSTRAINT valid_payment_status CHECK (payment_status IN ('pending', 'completed', 'failed', 'cancelled')),
    CONSTRAINT valid_currency CHECK (currency ~ '^[A-Z]{3}$')
);

-- Create indexes on payments
CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_payment_status ON payments(payment_status);
CREATE INDEX idx_payments_payment_method ON payments(payment_method);
CREATE INDEX idx_payments_gateway_transaction_id ON payments(gateway_transaction_id) WHERE gateway_transaction_id IS NOT NULL;
CREATE INDEX idx_payments_processed_at ON payments(processed_at DESC) WHERE processed_at IS NOT NULL;

-- Create trigger for payments updated_at
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create refunds table
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    payment_id UUID NOT NULL REFERENCES payments(id),
    refund_amount DECIMAL(10,2) NOT NULL,
    refund_reason VARCHAR(50) NOT NULL,
    refund_status VARCHAR(20) DEFAULT 'pending',
    processed_by UUID REFERENCES users(id),
    gateway_refund_id VARCHAR(255),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT refunds_refund_amount_check CHECK (refund_amount > 0),
    CONSTRAINT valid_refund_status CHECK (refund_status IN ('pending', 'processed', 'failed')),
    CONSTRAINT valid_refund_reason CHECK (refund_reason IN ('cancellation', 'schedule_change', 'no_show', 'service_issue', 'duplicate_payment', 'other'))
);

-- Create indexes on refunds
CREATE INDEX idx_refunds_booking_id ON refunds(booking_id);
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_refund_status ON refunds(refund_status);
CREATE INDEX idx_refunds_processed_by ON refunds(processed_by) WHERE processed_by IS NOT NULL;
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);

-- Create trigger for refunds updated_at
CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to generate QR code
CREATE OR REPLACE FUNCTION generate_qr_code(p_booking_id UUID, p_ticket_number INTEGER)
RETURNS VARCHAR(255) AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR(255);
    i INTEGER;
BEGIN
    -- Format: QR + booking_id (first 8 chars) + ticket_number + random (6 chars)
    result := 'QR' || substr(p_booking_id::text, 1, 8) || '-' || p_ticket_number::text || '-';
    
    -- Add random characters for uniqueness
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create function to update booking payment status
CREATE OR REPLACE FUNCTION update_booking_payment_status()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Update booking payment status based on payment status
        IF NEW.payment_status = 'completed' THEN
            UPDATE bookings 
            SET payment_status = 'paid',
                booking_status = CASE 
                    WHEN booking_status = 'pending' THEN 'confirmed'
                    ELSE booking_status
                END,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.booking_id;
        ELSIF NEW.payment_status = 'failed' THEN
            UPDATE bookings 
            SET payment_status = 'failed',
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.booking_id;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update booking payment status
CREATE TRIGGER update_booking_on_payment
    AFTER INSERT OR UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_booking_payment_status();

-- Create function to validate refund amount
CREATE OR REPLACE FUNCTION validate_refund_amount()
RETURNS TRIGGER AS $$
DECLARE
    total_paid DECIMAL(10,2);
    total_refunded DECIMAL(10,2);
BEGIN
    -- Get total paid amount for the payment
    SELECT amount INTO total_paid
    FROM payments
    WHERE id = NEW.payment_id;
    
    -- Get total already refunded for this payment
    SELECT COALESCE(SUM(refund_amount), 0) INTO total_refunded
    FROM refunds
    WHERE payment_id = NEW.payment_id
    AND id != NEW.id
    AND refund_status = 'processed';
    
    -- Check if refund amount exceeds paid amount
    IF (total_refunded + NEW.refund_amount) > total_paid THEN
        RAISE EXCEPTION 'Refund amount exceeds paid amount. Paid: %, Already refunded: %, Attempting to refund: %',
            total_paid, total_refunded, NEW.refund_amount;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to validate refund amount
CREATE TRIGGER validate_refund
    BEFORE INSERT OR UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION validate_refund_amount();

-- Create function to check in a ticket
CREATE OR REPLACE FUNCTION check_in_ticket(p_qr_code VARCHAR(255))
RETURNS TABLE(
    ticket_id UUID,
    passenger_name VARCHAR(255),
    seat_number VARCHAR(20),
    schedule_id UUID,
    departure_time TIME,
    status VARCHAR(20)
) AS $$
BEGIN
    -- Update ticket check-in status
    UPDATE tickets t
    SET check_in_status = 'checked_in',
        check_in_time = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    WHERE t.qr_code = p_qr_code
    AND t.check_in_status = 'not_checked_in';
    
    -- Return ticket information
    RETURN QUERY
    SELECT 
        t.id,
        t.passenger_name,
        t.seat_number,
        s.id,
        s.departure_time,
        t.check_in_status
    FROM tickets t
    JOIN bookings b ON b.id = t.booking_id
    JOIN schedules s ON s.id = b.schedule_id
    WHERE t.qr_code = p_qr_code;
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE tickets IS 'Individual passenger tickets within a booking';
COMMENT ON COLUMN tickets.qr_code IS 'Unique QR code for ticket validation and check-in';
COMMENT ON COLUMN tickets.check_in_status IS 'Current check-in status: not_checked_in, checked_in, or boarded';

COMMENT ON TABLE payments IS 'Payment transactions for bookings';
COMMENT ON COLUMN payments.gateway_transaction_id IS 'External payment gateway transaction reference';
COMMENT ON COLUMN payments.gateway_response IS 'Full response from payment gateway stored as JSON';

COMMENT ON TABLE refunds IS 'Refund transactions for cancelled or modified bookings';
COMMENT ON COLUMN refunds.refund_reason IS 'Reason for refund: cancellation, schedule_change, no_show, service_issue, duplicate_payment, or other';

COMMENT ON FUNCTION generate_qr_code(UUID, INTEGER) IS 'Generates unique QR code for tickets';
COMMENT ON FUNCTION update_booking_payment_status() IS 'Automatically updates booking status when payment is processed';
COMMENT ON FUNCTION validate_refund_amount() IS 'Ensures refund amount does not exceed paid amount';
COMMENT ON FUNCTION check_in_ticket(VARCHAR) IS 'Processes ticket check-in using QR code';