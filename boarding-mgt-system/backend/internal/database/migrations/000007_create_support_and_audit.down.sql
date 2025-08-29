-- Drop audit triggers
DROP TRIGGER IF EXISTS audit_schedules ON schedules;
DROP TRIGGER IF EXISTS audit_refunds ON refunds;
DROP TRIGGER IF EXISTS audit_payments ON payments;
DROP TRIGGER IF EXISTS audit_bookings ON bookings;
DROP TRIGGER IF EXISTS audit_users ON users;
DROP TRIGGER IF EXISTS audit_operators ON operators;

-- Drop support triggers
DROP TRIGGER IF EXISTS auto_assign_high_priority_tickets ON support_tickets;
DROP TRIGGER IF EXISTS update_support_tickets_updated_at ON support_tickets;

-- Drop functions
DROP FUNCTION IF EXISTS clean_old_audit_logs(INTEGER);
DROP FUNCTION IF EXISTS auto_assign_support_ticket();
DROP FUNCTION IF EXISTS audit_trigger_function();
DROP FUNCTION IF EXISTS generate_ticket_number();

-- Drop tables
DROP TABLE IF EXISTS audit.audit_logs CASCADE;
DROP TABLE IF EXISTS support_messages CASCADE;
DROP TABLE IF EXISTS support_tickets CASCADE;