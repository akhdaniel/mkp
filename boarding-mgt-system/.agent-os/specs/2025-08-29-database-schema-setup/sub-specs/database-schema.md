# Database Schema

This is the database schema implementation for the spec detailed in @.agent-os/specs/2025-08-29-database-schema-setup/spec.md

> Created: 2025-08-29
> Version: 1.0.0

## Schema Changes

### Core Entity Tables

```sql
-- Ferry operators (multi-tenant root entity)
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

-- Ports/terminals
CREATE TABLE ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    coordinates POINT,
    facilities JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vessels
CREATE TABLE vessels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(100) UNIQUE NOT NULL,
    vessel_type VARCHAR(50) NOT NULL, -- 'passenger', 'cargo', 'mixed'
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    deck_count INTEGER DEFAULT 1,
    seat_configuration JSONB NOT NULL, -- JSON structure of decks/sections/seats
    amenities JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Routes between ports
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
    CONSTRAINT different_ports CHECK (departure_port_id != arrival_port_id)
);
```

### User Management Schema

```sql
-- Users (customers, agents, operators)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    date_of_birth DATE,
    nationality VARCHAR(3), -- ISO country code
    user_type VARCHAR(20) NOT NULL, -- 'customer', 'agent', 'operator_admin'
    operator_id UUID REFERENCES operators(id), -- NULL for customers
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User sessions for JWT management
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
```

### Scheduling and Booking Schema

```sql
-- Ferry schedules
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    vessel_id UUID NOT NULL REFERENCES vessels(id) ON DELETE CASCADE,
    departure_date DATE NOT NULL,
    departure_time TIME NOT NULL,
    arrival_time TIME NOT NULL,
    base_price DECIMAL(10,2) NOT NULL CHECK (base_price >= 0),
    total_capacity INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled', -- 'scheduled', 'boarding', 'departed', 'arrived', 'cancelled'
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    version INTEGER DEFAULT 1, -- for optimistic locking
    CONSTRAINT valid_capacity CHECK (available_seats <= total_capacity AND available_seats >= 0)
);

-- Bookings
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_reference VARCHAR(20) UNIQUE NOT NULL,
    schedule_id UUID NOT NULL REFERENCES schedules(id),
    customer_id UUID NOT NULL REFERENCES users(id),
    passenger_count INTEGER NOT NULL CHECK (passenger_count > 0),
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    booking_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'confirmed', 'cancelled', 'refunded'
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'failed', 'refunded'
    booking_channel VARCHAR(20) NOT NULL, -- 'online', 'pos', 'phone', 'agent'
    special_requirements TEXT,
    booking_agent_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Individual tickets (passengers)
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    passenger_name VARCHAR(255) NOT NULL,
    passenger_type VARCHAR(20) DEFAULT 'adult', -- 'adult', 'child', 'infant', 'senior'
    seat_number VARCHAR(20),
    ticket_price DECIMAL(10,2) NOT NULL CHECK (ticket_price >= 0),
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    check_in_status VARCHAR(20) DEFAULT 'not_checked_in', -- 'not_checked_in', 'checked_in', 'boarded'
    check_in_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Payment and Financial Schema

```sql
-- Payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    payment_method VARCHAR(20) NOT NULL, -- 'credit_card', 'cash', 'bank_transfer', 'mobile_money'
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'completed', 'failed', 'cancelled'
    gateway_transaction_id VARCHAR(255),
    gateway_response JSONB,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Refunds
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    payment_id UUID NOT NULL REFERENCES payments(id),
    refund_amount DECIMAL(10,2) NOT NULL CHECK (refund_amount > 0),
    refund_reason VARCHAR(50) NOT NULL,
    refund_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processed', 'failed'
    processed_by UUID REFERENCES users(id),
    gateway_refund_id VARCHAR(255),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Support and Audit Schema

```sql
-- Customer support tickets
CREATE TABLE support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id UUID NOT NULL REFERENCES users(id),
    booking_id UUID REFERENCES bookings(id),
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    priority VARCHAR(10) DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'
    status VARCHAR(20) DEFAULT 'open', -- 'open', 'in_progress', 'resolved', 'closed'
    assigned_agent_id UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Audit log for tracking changes
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(100) NOT NULL,
    record_id UUID NOT NULL,
    action VARCHAR(20) NOT NULL, -- 'INSERT', 'UPDATE', 'DELETE'
    old_values JSONB,
    new_values JSONB,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Indexes and Performance Optimization

```sql
-- Performance indexes
CREATE INDEX idx_schedules_route_date ON schedules(route_id, departure_date);
CREATE INDEX idx_schedules_operator_date ON schedules(operator_id, departure_date);
CREATE INDEX idx_bookings_customer ON bookings(customer_id);
CREATE INDEX idx_bookings_schedule ON bookings(schedule_id);
CREATE INDEX idx_tickets_booking ON tickets(booking_id);
CREATE INDEX idx_tickets_qr_code ON tickets(qr_code);
CREATE INDEX idx_user_sessions_user ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token_hash);
CREATE INDEX idx_payments_booking ON payments(booking_id);
CREATE INDEX idx_support_tickets_customer ON support_tickets(customer_id);
CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);

-- Partial indexes for active records
CREATE INDEX idx_schedules_active ON schedules(departure_date, route_id) WHERE status != 'cancelled';
CREATE INDEX idx_vessels_active ON vessels(operator_id) WHERE is_active = true;
CREATE INDEX idx_bookings_confirmed ON bookings(schedule_id) WHERE booking_status = 'confirmed';
```

## Migrations

Migration files will be created using golang-migrate with sequential numbering:

1. **001_create_operators_and_ports.up.sql** - Core reference tables
2. **002_create_vessels_and_routes.up.sql** - Fleet management
3. **003_create_users_and_sessions.up.sql** - Authentication system
4. **004_create_schedules_and_bookings.up.sql** - Core booking system
5. **005_create_tickets_and_payments.up.sql** - Transaction processing
6. **006_create_support_and_audit.up.sql** - Support and logging
7. **007_create_indexes.up.sql** - Performance optimization

Each migration will include corresponding .down.sql files for rollback capability.