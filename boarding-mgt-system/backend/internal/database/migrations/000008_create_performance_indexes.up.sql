-- Performance optimization indexes for common queries

-- ============================================
-- User Authentication & Session Management
-- ============================================

-- Fast user lookup by email for authentication
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_active 
    ON users(email) 
    WHERE is_active = true;

-- Fast session validation
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_sessions_active_tokens 
    ON user_sessions(token_hash, expires_at) 
    WHERE is_active = true;

-- ============================================
-- Booking Search & Availability Queries
-- ============================================

-- Fast schedule search by date range and route
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_search 
    ON schedules(route_id, departure_date, departure_time) 
    WHERE status = 'scheduled';

-- Fast availability check
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_available 
    ON schedules(departure_date, available_seats) 
    WHERE status = 'scheduled' AND available_seats > 0;

-- Customer booking history (covering index)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_customer_complete 
    ON bookings(customer_id, created_at DESC) 
    INCLUDE (booking_reference, booking_status, total_amount);

-- Active bookings by schedule
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_schedule_active 
    ON bookings(schedule_id, booking_status) 
    WHERE booking_status IN ('pending', 'confirmed');

-- ============================================
-- Operational Queries
-- ============================================

-- Operator's active schedules
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_operator_active 
    ON schedules(operator_id, departure_date, status) 
    WHERE status != 'cancelled';

-- Vessel utilization
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_vessel_dates 
    ON schedules(vessel_id, departure_date, departure_time) 
    WHERE status IN ('scheduled', 'boarding', 'departed');

-- Route performance analysis
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_route_analysis 
    ON bookings b
    USING btree (
        (SELECT s.route_id FROM schedules s WHERE s.id = b.schedule_id),
        created_at DESC
    );

-- ============================================
-- Financial & Reporting Queries
-- ============================================

-- Payment reconciliation
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_date_status 
    ON payments(processed_at::date, payment_status) 
    WHERE payment_status = 'completed';

-- Daily revenue by operator
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_operator_revenue 
    ON bookings b
    USING btree (
        (SELECT s.operator_id FROM schedules s WHERE s.id = b.schedule_id),
        created_at::date,
        total_amount
    )
    WHERE booking_status = 'confirmed';

-- Pending refunds
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_refunds_pending 
    ON refunds(created_at, refund_status) 
    WHERE refund_status = 'pending';

-- ============================================
-- Support & Customer Service
-- ============================================

-- Open tickets by priority
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_support_priority_queue 
    ON support_tickets(priority, created_at) 
    WHERE status IN ('open', 'in_progress');

-- Agent workload
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_support_agent_workload 
    ON support_tickets(assigned_agent_id, status) 
    WHERE assigned_agent_id IS NOT NULL AND status != 'closed';

-- ============================================
-- Audit & Compliance
-- ============================================

-- Audit log date partitioning helper
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_date_partition 
    ON audit.audit_logs(changed_at::date, table_name);

-- User activity tracking
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_user_activity 
    ON audit.audit_logs(changed_by, changed_at DESC) 
    WHERE changed_by IS NOT NULL;

-- ============================================
-- Text Search Indexes
-- ============================================

-- Full text search on support tickets
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_support_tickets_search 
    ON support_tickets 
    USING gin(to_tsvector('english', subject || ' ' || description));

-- Port search by name or code
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_ports_search 
    ON ports 
    USING gin(to_tsvector('english', name || ' ' || code || ' ' || city));

-- ============================================
-- Composite & Partial Indexes for Complex Queries
-- ============================================

-- Schedule capacity planning
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_capacity_planning 
    ON schedules(operator_id, departure_date, vessel_id, available_seats) 
    WHERE status = 'scheduled' AND departure_date >= CURRENT_DATE;

-- Upcoming departures
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_schedules_upcoming 
    ON schedules(departure_date, departure_time, status) 
    WHERE departure_date >= CURRENT_DATE AND status = 'scheduled';

-- Recent bookings for fraud detection
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bookings_recent_by_customer 
    ON bookings(customer_id, created_at DESC) 
    WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours';

-- ============================================
-- Statistics Update for Query Planner
-- ============================================

-- Update statistics on frequently queried tables
ANALYZE operators;
ANALYZE users;
ANALYZE schedules;
ANALYZE bookings;
ANALYZE tickets;
ANALYZE payments;
ANALYZE support_tickets;

-- ============================================
-- Add Comments for Index Documentation
-- ============================================

COMMENT ON INDEX idx_users_email_active IS 'Optimizes user authentication queries';
COMMENT ON INDEX idx_schedules_search IS 'Optimizes schedule search by route and date';
COMMENT ON INDEX idx_bookings_customer_complete IS 'Covering index for customer booking history';
COMMENT ON INDEX idx_schedules_capacity_planning IS 'Optimizes capacity planning queries';
COMMENT ON INDEX idx_support_priority_queue IS 'Optimizes support ticket queue management';
COMMENT ON INDEX idx_audit_date_partition IS 'Helps with audit log partitioning by date';