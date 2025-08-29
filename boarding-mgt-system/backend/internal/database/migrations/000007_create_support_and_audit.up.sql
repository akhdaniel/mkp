-- Create support tickets table
CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id UUID NOT NULL REFERENCES users(id),
    booking_id UUID REFERENCES bookings(id),
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    priority VARCHAR(10) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'open',
    assigned_agent_id UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_priority CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    CONSTRAINT valid_ticket_status CHECK (status IN ('open', 'in_progress', 'resolved', 'closed'))
);

-- Create indexes on support_tickets
CREATE INDEX idx_support_tickets_ticket_number ON support_tickets(ticket_number);
CREATE INDEX idx_support_tickets_customer_id ON support_tickets(customer_id);
CREATE INDEX idx_support_tickets_booking_id ON support_tickets(booking_id) WHERE booking_id IS NOT NULL;
CREATE INDEX idx_support_tickets_status ON support_tickets(status) WHERE status != 'closed';
CREATE INDEX idx_support_tickets_priority ON support_tickets(priority) WHERE priority IN ('high', 'urgent');
CREATE INDEX idx_support_tickets_assigned_agent ON support_tickets(assigned_agent_id) WHERE assigned_agent_id IS NOT NULL;
CREATE INDEX idx_support_tickets_created_at ON support_tickets(created_at DESC);

-- Create trigger for support_tickets updated_at
CREATE TRIGGER update_support_tickets_updated_at BEFORE UPDATE ON support_tickets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create support messages table
CREATE TABLE support_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    message TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT false, -- Internal notes not visible to customers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on support_messages
CREATE INDEX idx_support_messages_ticket_id ON support_messages(ticket_id);
CREATE INDEX idx_support_messages_sender_id ON support_messages(sender_id);
CREATE INDEX idx_support_messages_created_at ON support_messages(created_at DESC);

-- Create audit logs table in audit schema
CREATE TABLE audit.audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(100) NOT NULL,
    record_id UUID NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_audit_action CHECK (action IN ('INSERT', 'UPDATE', 'DELETE'))
);

-- Create indexes on audit_logs
CREATE INDEX idx_audit_logs_table_record ON audit.audit_logs(table_name, record_id);
CREATE INDEX idx_audit_logs_action ON audit.audit_logs(action);
CREATE INDEX idx_audit_logs_changed_at ON audit.audit_logs(changed_at DESC);
CREATE INDEX idx_audit_logs_changed_by ON audit.audit_logs(changed_by) WHERE changed_by IS NOT NULL;

-- Create function to generate support ticket number
CREATE OR REPLACE FUNCTION generate_ticket_number()
RETURNS VARCHAR(20) AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR(20);
    i INTEGER;
BEGIN
    -- Format: ST + year + month + random (6 chars)
    result := 'ST' || to_char(CURRENT_TIMESTAMP, 'YYMM');
    
    -- Add random characters
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create generic audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
DECLARE
    audit_user UUID;
    old_data JSONB;
    new_data JSONB;
BEGIN
    -- Get the user ID from context (would be set by application)
    -- For now, we'll use NULL if not set
    audit_user := current_setting('app.current_user_id', true)::UUID;
    
    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit.audit_logs (
            table_name, record_id, action, new_values, changed_by
        ) VALUES (
            TG_TABLE_NAME, NEW.id, 'INSERT', to_jsonb(NEW), audit_user
        );
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Only log if there are actual changes
        IF NEW IS DISTINCT FROM OLD THEN
            INSERT INTO audit.audit_logs (
                table_name, record_id, action, old_values, new_values, changed_by
            ) VALUES (
                TG_TABLE_NAME, NEW.id, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW), audit_user
            );
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit.audit_logs (
            table_name, record_id, action, old_values, changed_by
        ) VALUES (
            TG_TABLE_NAME, OLD.id, 'DELETE', to_jsonb(OLD), audit_user
        );
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create audit triggers for important tables
CREATE TRIGGER audit_operators AFTER INSERT OR UPDATE OR DELETE ON operators
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

CREATE TRIGGER audit_users AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

CREATE TRIGGER audit_bookings AFTER INSERT OR UPDATE OR DELETE ON bookings
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

CREATE TRIGGER audit_payments AFTER INSERT OR UPDATE OR DELETE ON payments
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

CREATE TRIGGER audit_refunds AFTER INSERT OR UPDATE OR DELETE ON refunds
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

CREATE TRIGGER audit_schedules AFTER INSERT OR UPDATE OR DELETE ON schedules
    FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();

-- Create function to auto-assign support tickets
CREATE OR REPLACE FUNCTION auto_assign_support_ticket()
RETURNS TRIGGER AS $$
DECLARE
    available_agent UUID;
BEGIN
    -- Only auto-assign if priority is high or urgent and no agent assigned
    IF NEW.priority IN ('high', 'urgent') AND NEW.assigned_agent_id IS NULL THEN
        -- Find agent with least active tickets
        SELECT u.id INTO available_agent
        FROM users u
        LEFT JOIN (
            SELECT assigned_agent_id, COUNT(*) as ticket_count
            FROM support_tickets
            WHERE status IN ('open', 'in_progress')
            AND assigned_agent_id IS NOT NULL
            GROUP BY assigned_agent_id
        ) tc ON tc.assigned_agent_id = u.id
        WHERE u.user_type IN ('agent', 'operator_admin')
        AND u.is_active = true
        ORDER BY COALESCE(tc.ticket_count, 0), random()
        LIMIT 1;
        
        IF available_agent IS NOT NULL THEN
            NEW.assigned_agent_id := available_agent;
            NEW.status := 'in_progress';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for auto-assignment
CREATE TRIGGER auto_assign_high_priority_tickets
    BEFORE INSERT ON support_tickets
    FOR EACH ROW EXECUTE FUNCTION auto_assign_support_ticket();

-- Create function to clean old audit logs
CREATE OR REPLACE FUNCTION clean_old_audit_logs(retention_days INTEGER DEFAULT 90)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM audit.audit_logs
    WHERE changed_at < CURRENT_TIMESTAMP - (retention_days || ' days')::INTERVAL;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE support_tickets IS 'Customer support tickets for issues and inquiries';
COMMENT ON COLUMN support_tickets.priority IS 'Ticket priority: low, normal, high, or urgent';
COMMENT ON COLUMN support_tickets.status IS 'Ticket status: open, in_progress, resolved, or closed';

COMMENT ON TABLE support_messages IS 'Messages and communications within support tickets';
COMMENT ON COLUMN support_messages.is_internal IS 'Internal notes visible only to support staff';

COMMENT ON TABLE audit.audit_logs IS 'Audit trail for all changes to critical tables';
COMMENT ON COLUMN audit.audit_logs.old_values IS 'Previous values before change (UPDATE/DELETE only)';
COMMENT ON COLUMN audit.audit_logs.new_values IS 'New values after change (INSERT/UPDATE only)';

COMMENT ON FUNCTION generate_ticket_number() IS 'Generates unique support ticket number';
COMMENT ON FUNCTION audit_trigger_function() IS 'Generic audit logging trigger for any table';
COMMENT ON FUNCTION auto_assign_support_ticket() IS 'Automatically assigns high priority tickets to available agents';
COMMENT ON FUNCTION clean_old_audit_logs(INTEGER) IS 'Removes audit logs older than specified days';