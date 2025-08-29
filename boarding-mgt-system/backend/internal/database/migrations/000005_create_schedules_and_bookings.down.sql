-- Drop triggers
DROP TRIGGER IF EXISTS manage_schedule_availability ON bookings;
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP TRIGGER IF EXISTS update_schedules_updated_at ON schedules;

-- Drop functions
DROP FUNCTION IF EXISTS update_schedule_availability();
DROP FUNCTION IF EXISTS generate_booking_reference();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS schedules CASCADE;