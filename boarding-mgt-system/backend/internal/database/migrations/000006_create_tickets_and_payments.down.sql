-- Drop triggers
DROP TRIGGER IF EXISTS validate_refund ON refunds;
DROP TRIGGER IF EXISTS update_booking_on_payment ON payments;
DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;

-- Drop functions
DROP FUNCTION IF EXISTS check_in_ticket(VARCHAR);
DROP FUNCTION IF EXISTS validate_refund_amount();
DROP FUNCTION IF EXISTS update_booking_payment_status();
DROP FUNCTION IF EXISTS generate_qr_code(UUID, INTEGER);

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS refunds CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;