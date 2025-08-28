# Database Schema

This is the database schema implementation for the spec detailed in @.agent-os/specs/2025-08-28-database-schema-setup/spec.md

> Created: 2025-08-28
> Version: 1.0.0

## Database Structure Overview

### Schema Organization
- **public**: Shared system tables and tenant registry
- **tenant_{id}**: Dynamic schemas for each ferry operator
- **audit**: Centralized audit logging

## Core Tables Design

### Public Schema Tables

```sql
-- Tenant registry for multi-tenant support
CREATE TABLE public.tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  schema_name VARCHAR(63) UNIQUE NOT NULL,
  config JSONB DEFAULT '{}',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- System-wide configuration
CREATE TABLE public.system_config (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key VARCHAR(255) UNIQUE NOT NULL,
  value JSONB NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Supported currencies
CREATE TABLE public.currencies (
  code VARCHAR(3) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  symbol VARCHAR(5),
  decimal_places INTEGER DEFAULT 2
);

-- Countries reference
CREATE TABLE public.countries (
  code VARCHAR(2) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone_code VARCHAR(10)
);
```

### Tenant Schema Tables

```sql
-- Users and authentication
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'operator', 'agent', 'customer', 'support')),
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  phone VARCHAR(50),
  metadata JSONB DEFAULT '{}',
  is_active BOOLEAN DEFAULT true,
  email_verified BOOLEAN DEFAULT false,
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- Authentication tokens
CREATE TABLE auth_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) UNIQUE NOT NULL,
  token_type VARCHAR(20) NOT NULL CHECK (token_type IN ('access', 'refresh', 'reset')),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Operators (ferry companies)
CREATE TABLE operators (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  legal_name VARCHAR(255),
  tax_id VARCHAR(50),
  contact_email VARCHAR(255),
  contact_phone VARCHAR(50),
  address JSONB,
  config JSONB DEFAULT '{}',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- Ports/Terminals
CREATE TABLE ports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(10) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  city VARCHAR(100),
  country_code VARCHAR(2) REFERENCES public.countries(code),
  timezone VARCHAR(50) NOT NULL,
  coordinates JSONB,
  facilities JSONB DEFAULT '[]',
  contact_info JSONB,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- Vessels
CREATE TABLE vessels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operator_id UUID NOT NULL REFERENCES operators(id),
  code VARCHAR(20) NOT NULL,
  name VARCHAR(255) NOT NULL,
  vessel_type VARCHAR(50),
  capacity_passengers INTEGER NOT NULL,
  capacity_vehicles INTEGER DEFAULT 0,
  facilities JSONB DEFAULT '[]',
  seat_configuration JSONB NOT NULL,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  UNIQUE(operator_id, code)
);

-- Routes
CREATE TABLE routes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operator_id UUID NOT NULL REFERENCES operators(id),
  code VARCHAR(20) NOT NULL,
  name VARCHAR(255) NOT NULL,
  origin_port_id UUID NOT NULL REFERENCES ports(id),
  destination_port_id UUID NOT NULL REFERENCES ports(id),
  distance_km DECIMAL(10,2),
  duration_minutes INTEGER,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  UNIQUE(operator_id, code)
);

-- Schedules
CREATE TABLE schedules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  route_id UUID NOT NULL REFERENCES routes(id),
  vessel_id UUID NOT NULL REFERENCES vessels(id),
  departure_time TIMESTAMPTZ NOT NULL,
  arrival_time TIMESTAMPTZ NOT NULL,
  schedule_type VARCHAR(20) CHECK (schedule_type IN ('regular', 'special', 'seasonal')),
  recurrence_rule JSONB,
  valid_from DATE NOT NULL,
  valid_until DATE,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- Schedule instances (actual sailings)
CREATE TABLE schedule_instances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  schedule_id UUID NOT NULL REFERENCES schedules(id),
  vessel_id UUID NOT NULL REFERENCES vessels(id),
  departure_time TIMESTAMPTZ NOT NULL,
  arrival_time TIMESTAMPTZ NOT NULL,
  status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'boarding', 'departed', 'arrived', 'cancelled', 'delayed')),
  available_seats INTEGER,
  booked_seats INTEGER DEFAULT 0,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Seat classes
CREATE TABLE seat_classes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operator_id UUID NOT NULL REFERENCES operators(id),
  code VARCHAR(20) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  amenities JSONB DEFAULT '[]',
  priority_boarding BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(operator_id, code)
);

-- Pricing rules
CREATE TABLE pricing_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  route_id UUID NOT NULL REFERENCES routes(id),
  seat_class_id UUID REFERENCES seat_classes(id),
  passenger_type VARCHAR(20) DEFAULT 'adult' CHECK (passenger_type IN ('adult', 'child', 'infant', 'senior', 'student')),
  base_price DECIMAL(19,4) NOT NULL,
  currency_code VARCHAR(3) NOT NULL REFERENCES public.currencies(code),
  valid_from DATE NOT NULL,
  valid_until DATE,
  conditions JSONB DEFAULT '{}',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Bookings
CREATE TABLE bookings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_reference VARCHAR(20) UNIQUE NOT NULL,
  schedule_instance_id UUID NOT NULL REFERENCES schedule_instances(id),
  user_id UUID REFERENCES users(id),
  customer_email VARCHAR(255) NOT NULL,
  customer_phone VARCHAR(50),
  customer_name VARCHAR(255) NOT NULL,
  total_amount DECIMAL(19,4) NOT NULL,
  currency_code VARCHAR(3) NOT NULL REFERENCES public.currencies(code),
  payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'processing', 'paid', 'failed', 'refunded', 'partial_refund')),
  booking_status VARCHAR(20) DEFAULT 'confirmed' CHECK (booking_status IN ('draft', 'confirmed', 'cancelled', 'completed', 'no_show')),
  booking_source VARCHAR(20) CHECK (booking_source IN ('online', 'pos', 'phone', 'agent')),
  agent_id UUID REFERENCES users(id),
  metadata JSONB DEFAULT '{}',
  expires_at TIMESTAMPTZ,
  confirmed_at TIMESTAMPTZ,
  cancelled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Passengers in bookings
CREATE TABLE booking_passengers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
  passenger_type VARCHAR(20) DEFAULT 'adult' CHECK (passenger_type IN ('adult', 'child', 'infant', 'senior', 'student')),
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  document_type VARCHAR(20),
  document_number VARCHAR(50),
  date_of_birth DATE,
  nationality VARCHAR(2) REFERENCES public.countries(code),
  special_requirements JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tickets (one per passenger)
CREATE TABLE tickets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_number VARCHAR(30) UNIQUE NOT NULL,
  booking_id UUID NOT NULL REFERENCES bookings(id),
  passenger_id UUID NOT NULL REFERENCES booking_passengers(id),
  seat_class_id UUID REFERENCES seat_classes(id),
  seat_number VARCHAR(10),
  price DECIMAL(19,4) NOT NULL,
  qr_code TEXT,
  status VARCHAR(20) DEFAULT 'valid' CHECK (status IN ('valid', 'used', 'cancelled', 'expired')),
  checked_in_at TIMESTAMPTZ,
  boarded_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Payment transactions
CREATE TABLE payment_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id UUID NOT NULL REFERENCES bookings(id),
  transaction_reference VARCHAR(100) UNIQUE NOT NULL,
  amount DECIMAL(19,4) NOT NULL,
  currency_code VARCHAR(3) NOT NULL REFERENCES public.currencies(code),
  payment_method VARCHAR(20) CHECK (payment_method IN ('card', 'cash', 'bank_transfer', 'wallet', 'other')),
  payment_provider VARCHAR(50),
  provider_reference VARCHAR(255),
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'success', 'failed', 'cancelled')),
  metadata JSONB DEFAULT '{}',
  processed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Refunds
CREATE TABLE refunds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id UUID NOT NULL REFERENCES bookings(id),
  payment_transaction_id UUID REFERENCES payment_transactions(id),
  refund_reference VARCHAR(100) UNIQUE NOT NULL,
  amount DECIMAL(19,4) NOT NULL,
  currency_code VARCHAR(3) NOT NULL REFERENCES public.currencies(code),
  reason VARCHAR(255),
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'processing', 'completed', 'failed', 'rejected')),
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMPTZ,
  processed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Helpdesk tickets
CREATE TABLE helpdesk_tickets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_number VARCHAR(20) UNIQUE NOT NULL,
  booking_id UUID REFERENCES bookings(id),
  user_id UUID REFERENCES users(id),
  assigned_to UUID REFERENCES users(id),
  category VARCHAR(50),
  priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
  status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'waiting_customer', 'waiting_internal', 'resolved', 'closed')),
  subject VARCHAR(255) NOT NULL,
  description TEXT,
  resolution TEXT,
  metadata JSONB DEFAULT '{}',
  resolved_at TIMESTAMPTZ,
  closed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Helpdesk conversations
CREATE TABLE helpdesk_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_id UUID NOT NULL REFERENCES helpdesk_tickets(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id),
  sender_type VARCHAR(20) CHECK (sender_type IN ('customer', 'agent', 'system', 'bot')),
  message TEXT NOT NULL,
  attachments JSONB DEFAULT '[]',
  is_internal BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- FAQ categories
CREATE TABLE faq_categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  display_order INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- FAQ items
CREATE TABLE faq_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  category_id UUID NOT NULL REFERENCES faq_categories(id),
  question TEXT NOT NULL,
  answer TEXT NOT NULL,
  tags JSONB DEFAULT '[]',
  views INTEGER DEFAULT 0,
  helpful_count INTEGER DEFAULT 0,
  display_order INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### Audit Schema Tables

```sql
-- Audit log for all data changes
CREATE TABLE audit.audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES public.tenants(id),
  table_name VARCHAR(63) NOT NULL,
  record_id UUID NOT NULL,
  action VARCHAR(20) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
  user_id UUID,
  old_data JSONB,
  new_data JSONB,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create monthly partitions for audit logs
CREATE TABLE audit.audit_logs_2025_01 PARTITION OF audit.audit_logs
  FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
```

## Indexes

```sql
-- User indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;

-- Booking indexes  
CREATE INDEX idx_bookings_reference ON bookings(booking_reference);
CREATE INDEX idx_bookings_schedule ON bookings(schedule_instance_id);
CREATE INDEX idx_bookings_user ON bookings(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_bookings_status ON bookings(booking_status) WHERE booking_status != 'completed';
CREATE INDEX idx_bookings_created ON bookings(created_at DESC);

-- Schedule indexes
CREATE INDEX idx_schedules_route ON schedules(route_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_schedule_instances_departure ON schedule_instances(departure_time) WHERE status = 'scheduled';

-- Ticket indexes
CREATE INDEX idx_tickets_booking ON tickets(booking_id);
CREATE INDEX idx_tickets_number ON tickets(ticket_number);

-- Helpdesk indexes
CREATE INDEX idx_helpdesk_tickets_user ON helpdesk_tickets(user_id);
CREATE INDEX idx_helpdesk_tickets_status ON helpdesk_tickets(status) WHERE status NOT IN ('resolved', 'closed');

-- Full-text search indexes
CREATE INDEX idx_faq_search ON faq_items USING gin(to_tsvector('english', question || ' ' || answer));
```

## Constraints and Triggers

```sql
-- Ensure departure is before arrival
ALTER TABLE schedules ADD CONSTRAINT chk_schedule_times 
  CHECK (departure_time < arrival_time);

ALTER TABLE schedule_instances ADD CONSTRAINT chk_instance_times
  CHECK (departure_time < arrival_time);

-- Prevent overbooking
ALTER TABLE schedule_instances ADD CONSTRAINT chk_booking_capacity
  CHECK (booked_seats <= available_seats);

-- Ensure valid price ranges
ALTER TABLE pricing_rules ADD CONSTRAINT chk_price_positive
  CHECK (base_price > 0);

-- Create update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply update trigger to all tables with updated_at
CREATE TRIGGER update_users_timestamp BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_bookings_timestamp BEFORE UPDATE ON bookings
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Add triggers for other tables...

-- Audit trigger function
CREATE OR REPLACE FUNCTION audit.log_changes()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'DELETE' THEN
    INSERT INTO audit.audit_logs (tenant_id, table_name, record_id, action, old_data, new_data)
    VALUES (
      current_setting('app.current_tenant')::UUID,
      TG_TABLE_NAME,
      OLD.id,
      TG_OP,
      row_to_json(OLD),
      NULL
    );
    RETURN OLD;
  ELSE
    INSERT INTO audit.audit_logs (tenant_id, table_name, record_id, action, old_data, new_data)
    VALUES (
      current_setting('app.current_tenant')::UUID,
      TG_TABLE_NAME,
      NEW.id,
      TG_OP,
      CASE WHEN TG_OP = 'UPDATE' THEN row_to_json(OLD) ELSE NULL END,
      row_to_json(NEW)
    );
    RETURN NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;
```

## Migration Strategy

1. Create database and enable required extensions
2. Create schemas (public, audit)
3. Create public schema tables
4. Create template tenant schema
5. Create audit schema and partitions
6. Apply indexes and constraints
7. Create trigger functions and apply triggers
8. Insert initial reference data

## Rollback Plan

Each migration should have a corresponding rollback script to undo changes if needed. Store migrations in numbered files with up/down SQL scripts.