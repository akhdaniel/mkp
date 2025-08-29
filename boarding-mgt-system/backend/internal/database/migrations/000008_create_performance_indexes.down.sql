-- Drop all performance optimization indexes

-- User Authentication & Session Management
DROP INDEX IF EXISTS idx_users_email_active;
DROP INDEX IF EXISTS idx_user_sessions_active_tokens;

-- Booking Search & Availability
DROP INDEX IF EXISTS idx_schedules_search;
DROP INDEX IF EXISTS idx_schedules_available;
DROP INDEX IF EXISTS idx_bookings_customer_complete;
DROP INDEX IF EXISTS idx_bookings_schedule_active;

-- Operational Queries
DROP INDEX IF EXISTS idx_schedules_operator_active;
DROP INDEX IF EXISTS idx_schedules_vessel_dates;
DROP INDEX IF EXISTS idx_bookings_route_analysis;

-- Financial & Reporting
DROP INDEX IF EXISTS idx_payments_date_status;
DROP INDEX IF EXISTS idx_bookings_operator_revenue;
DROP INDEX IF EXISTS idx_refunds_pending;

-- Support & Customer Service
DROP INDEX IF EXISTS idx_support_priority_queue;
DROP INDEX IF EXISTS idx_support_agent_workload;

-- Audit & Compliance
DROP INDEX IF EXISTS idx_audit_date_partition;
DROP INDEX IF EXISTS idx_audit_user_activity;

-- Text Search
DROP INDEX IF EXISTS idx_support_tickets_search;
DROP INDEX IF EXISTS idx_ports_search;

-- Composite & Partial Indexes
DROP INDEX IF EXISTS idx_schedules_capacity_planning;
DROP INDEX IF EXISTS idx_schedules_upcoming;
DROP INDEX IF EXISTS idx_bookings_recent_by_customer;