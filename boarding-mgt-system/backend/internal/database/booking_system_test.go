package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingSystemTables(t *testing.T) {
	cfg, err := config.LoadTest()
	require.NoError(t, err, "Failed to load test config")

	db, err := New(&cfg.Database)
	require.NoError(t, err, "Failed to connect to database")
	defer db.Close()

	ctx := context.Background()

	// Ensure migrations are run first
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	migrator, err := NewMigrator(databaseURL)
	require.NoError(t, err, "Failed to create migrator")
	defer migrator.Close()

	// Run migrations up to latest
	err = migrator.Up()
	assert.NoError(t, err, "Failed to run migrations")

	t.Run("Verify schedules table structure", func(t *testing.T) {
		// Check table exists
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'schedules'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "schedules table should exist")

		// Verify version column for optimistic locking
		var columnExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_schema = 'public' 
				AND table_name = 'schedules'
				AND column_name = 'version'
			)
		`).Scan(&columnExists)
		assert.NoError(t, err)
		assert.True(t, columnExists, "version column should exist for optimistic locking")

		// Verify check constraint on capacity
		var checkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'valid_capacity'
				AND contype = 'c'
			)
		`).Scan(&checkExists)
		assert.NoError(t, err)
		assert.True(t, checkExists, "Check constraint on capacity should exist")

		// Verify foreign keys
		var fkCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM information_schema.table_constraints 
			WHERE constraint_type = 'FOREIGN KEY' 
			AND table_name = 'schedules'
		`).Scan(&fkCount)
		assert.NoError(t, err)
		assert.Equal(t, 3, fkCount, "Schedules should have 3 foreign keys")
	})

	t.Run("Verify bookings table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'bookings'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "bookings table should exist")

		// Verify unique constraint on booking_reference
		var constraintExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'bookings_booking_reference_key'
			)
		`).Scan(&constraintExists)
		assert.NoError(t, err)
		assert.True(t, constraintExists, "Unique constraint on booking_reference should exist")

		// Verify check constraint on passenger_count
		var checkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'bookings_passenger_count_check'
				AND contype = 'c'
			)
		`).Scan(&checkExists)
		assert.NoError(t, err)
		assert.True(t, checkExists, "Check constraint on passenger_count should exist")
	})

	t.Run("Verify tickets table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'tickets'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "tickets table should exist")

		// Verify unique constraint on qr_code
		var constraintExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'tickets_qr_code_key'
			)
		`).Scan(&constraintExists)
		assert.NoError(t, err)
		assert.True(t, constraintExists, "Unique constraint on qr_code should exist")

		// Verify cascade delete from bookings
		var fkConstraint string
		err = db.Pool.QueryRow(ctx, `
			SELECT confdeltype 
			FROM pg_constraint 
			WHERE conname LIKE 'tickets_booking_id_fkey'
		`).Scan(&fkConstraint)
		assert.NoError(t, err)
		assert.Equal(t, "c", fkConstraint, "Should have CASCADE delete")
	})

	t.Run("Test complete booking workflow", func(t *testing.T) {
		// Setup: Create operator, ports, vessel, route
		var operatorID, port1ID, port2ID, vesselID, routeID, customerID string

		// Create operator
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Booking Test Ferry", "BTF001", "booking@ferry.com").Scan(&operatorID)
		require.NoError(t, err)

		// Create ports
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO ports (name, code, city, country, timezone)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "Departure Port", "DEP001", "City A", "Country", "UTC").Scan(&port1ID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO ports (name, code, city, country, timezone)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "Arrival Port", "ARR001", "City B", "Country", "UTC").Scan(&port2ID)
		require.NoError(t, err)

		// Create vessel
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO vessels (operator_id, name, registration_number, vessel_type, capacity, seat_configuration)
			VALUES ($1, $2, $3, $4, $5, $6::jsonb)
			RETURNING id
		`, operatorID, "Test Ferry", "TF001", "passenger", 100, `{"decks": 2}`).Scan(&vesselID)
		require.NoError(t, err)

		// Create route
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, estimated_duration)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, operatorID, "Test Route", port1ID, port2ID, "2 hours").Scan(&routeID)
		require.NoError(t, err)

		// Create customer
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, user_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "bookingcustomer@example.com", "$2a$10$hash", "Booking", "Customer", "customer").Scan(&customerID)
		require.NoError(t, err)

		// Create schedule
		var scheduleID string
		departureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO schedules (
				operator_id, route_id, vessel_id, 
				departure_date, departure_time, arrival_time,
				base_price, total_capacity, available_seats
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, operatorID, routeID, vesselID, 
			departureDate, "10:00:00", "12:00:00",
			50.00, 100, 100).Scan(&scheduleID)
		assert.NoError(t, err, "Should create schedule")

		// Create booking
		var bookingID string
		bookingRef := fmt.Sprintf("BK%d", time.Now().Unix())
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO bookings (
				booking_reference, schedule_id, customer_id,
				passenger_count, total_amount, booking_channel
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, bookingRef, scheduleID, customerID, 2, 100.00, "online").Scan(&bookingID)
		assert.NoError(t, err, "Should create booking")

		// Create tickets
		var ticket1ID, ticket2ID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO tickets (
				booking_id, passenger_name, passenger_type,
				seat_number, ticket_price, qr_code
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, bookingID, "John Doe", "adult", "A1", 50.00, "QR001"+bookingRef).Scan(&ticket1ID)
		assert.NoError(t, err, "Should create first ticket")

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO tickets (
				booking_id, passenger_name, passenger_type,
				seat_number, ticket_price, qr_code
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, bookingID, "Jane Doe", "adult", "A2", 50.00, "QR002"+bookingRef).Scan(&ticket2ID)
		assert.NoError(t, err, "Should create second ticket")

		// Verify tickets are linked to booking
		var ticketCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM tickets WHERE booking_id = $1
		`, bookingID).Scan(&ticketCount)
		assert.NoError(t, err)
		assert.Equal(t, 2, ticketCount, "Should have 2 tickets")

		// Test cascade delete
		_, err = db.Pool.Exec(ctx, "DELETE FROM bookings WHERE id = $1", bookingID)
		assert.NoError(t, err)

		// Verify tickets were deleted
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM tickets WHERE booking_id = $1
		`, bookingID).Scan(&ticketCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, ticketCount, "Tickets should be deleted via cascade")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM ports WHERE id IN ($1, $2)", port1ID, port2ID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", customerID)
		assert.NoError(t, err)
	})

	t.Run("Test seat capacity constraints", func(t *testing.T) {
		// Setup minimal data
		var operatorID, routeID, vesselID, scheduleID string
		
		// Create minimal setup (reuse from previous test setup)
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ('Capacity Test', 'CAP001', 'cap@test.com')
			RETURNING id
		`).Scan(&operatorID)
		require.NoError(t, err)

		// Create minimal route and vessel (simplified)
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, estimated_duration)
			SELECT $1, 'Test Route', 
				(SELECT id FROM ports LIMIT 1), 
				(SELECT id FROM ports OFFSET 1 LIMIT 1),
				'1 hour'
			RETURNING id
		`, operatorID).Scan(&routeID)
		
		// If no ports exist, create them
		if err != nil {
			var p1, p2 string
			db.Pool.QueryRow(ctx, `
				INSERT INTO ports (name, code, city, country, timezone)
				VALUES ('Port1', 'P1', 'City1', 'Country', 'UTC')
				RETURNING id
			`).Scan(&p1)
			db.Pool.QueryRow(ctx, `
				INSERT INTO ports (name, code, city, country, timezone)
				VALUES ('Port2', 'P2', 'City2', 'Country', 'UTC')
				RETURNING id
			`).Scan(&p2)
			db.Pool.QueryRow(ctx, `
				INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, estimated_duration)
				VALUES ($1, 'Test Route', $2, $3, '1 hour')
				RETURNING id
			`, operatorID, p1, p2).Scan(&routeID)
		}

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO vessels (operator_id, name, registration_number, vessel_type, capacity, seat_configuration)
			VALUES ($1, 'Small Vessel', 'SV001', 'passenger', 5, '{}')
			RETURNING id
		`, operatorID).Scan(&vesselID)
		require.NoError(t, err)

		// Create schedule with limited capacity
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO schedules (
				operator_id, route_id, vessel_id,
				departure_date, departure_time, arrival_time,
				base_price, total_capacity, available_seats
			) VALUES ($1, $2, $3, CURRENT_DATE + INTERVAL '1 day', '10:00', '11:00', 10.00, 5, 5)
			RETURNING id
		`, operatorID, routeID, vesselID).Scan(&scheduleID)
		assert.NoError(t, err)

		// Test capacity constraint
		_, err = db.Pool.Exec(ctx, `
			UPDATE schedules 
			SET available_seats = 10 
			WHERE id = $1
		`, scheduleID)
		assert.Error(t, err, "Should fail when available_seats > total_capacity")
		assert.Contains(t, err.Error(), "valid_capacity", "Error should mention capacity constraint")

		// Test negative seats
		_, err = db.Pool.Exec(ctx, `
			UPDATE schedules 
			SET available_seats = -1 
			WHERE id = $1
		`, scheduleID)
		assert.Error(t, err, "Should fail with negative available_seats")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
	})

	t.Run("Test optimistic locking with version", func(t *testing.T) {
		// Create a schedule with version tracking
		var scheduleID string
		var operatorID, routeID, vesselID string

		// Minimal setup
		db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ('Lock Test', 'LCK001', 'lock@test.com')
			RETURNING id
		`).Scan(&operatorID)

		// Create or get ports for route
		var p1, p2 string
		err := db.Pool.QueryRow(ctx, `SELECT id FROM ports LIMIT 1`).Scan(&p1)
		if err != nil {
			db.Pool.QueryRow(ctx, `
				INSERT INTO ports (name, code, city, country, timezone)
				VALUES ('LockPort1', 'LP1', 'City', 'Country', 'UTC')
				RETURNING id
			`).Scan(&p1)
		}
		err = db.Pool.QueryRow(ctx, `SELECT id FROM ports WHERE id != $1 LIMIT 1`, p1).Scan(&p2)
		if err != nil {
			db.Pool.QueryRow(ctx, `
				INSERT INTO ports (name, code, city, country, timezone)
				VALUES ('LockPort2', 'LP2', 'City', 'Country', 'UTC')
				RETURNING id
			`).Scan(&p2)
		}

		db.Pool.QueryRow(ctx, `
			INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, estimated_duration)
			VALUES ($1, 'Lock Route', $2, $3, '1 hour')
			RETURNING id
		`, operatorID, p1, p2).Scan(&routeID)

		db.Pool.QueryRow(ctx, `
			INSERT INTO vessels (operator_id, name, registration_number, vessel_type, capacity, seat_configuration)
			VALUES ($1, 'Lock Vessel', 'LV001', 'passenger', 100, '{}')
			RETURNING id
		`, operatorID).Scan(&vesselID)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO schedules (
				operator_id, route_id, vessel_id,
				departure_date, departure_time, arrival_time,
				base_price, total_capacity, available_seats, version
			) VALUES ($1, $2, $3, CURRENT_DATE + INTERVAL '2 days', '10:00', '12:00', 50.00, 100, 100, 1)
			RETURNING id
		`, operatorID, routeID, vesselID).Scan(&scheduleID)
		require.NoError(t, err)

		// Simulate concurrent booking attempts
		var version1, version2 int
		var seats1, seats2 int

		// Read version and seats (simulating two concurrent transactions)
		err = db.Pool.QueryRow(ctx, `
			SELECT version, available_seats FROM schedules WHERE id = $1
		`, scheduleID).Scan(&version1, &seats1)
		require.NoError(t, err)

		version2 = version1
		seats2 = seats1

		// First update succeeds
		tag, err := db.Pool.Exec(ctx, `
			UPDATE schedules 
			SET available_seats = available_seats - 2,
			    version = version + 1
			WHERE id = $1 AND version = $2
		`, scheduleID, version1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), tag.RowsAffected(), "First update should succeed")

		// Second update fails (wrong version)
		tag, err = db.Pool.Exec(ctx, `
			UPDATE schedules 
			SET available_seats = available_seats - 3,
			    version = version + 1
			WHERE id = $1 AND version = $2
		`, scheduleID, version2)
		assert.NoError(t, err) // No error, but no rows affected
		assert.Equal(t, int64(0), tag.RowsAffected(), "Second update should fail due to version mismatch")

		// Verify final state
		var finalSeats, finalVersion int
		err = db.Pool.QueryRow(ctx, `
			SELECT available_seats, version FROM schedules WHERE id = $1
		`, scheduleID).Scan(&finalSeats, &finalVersion)
		assert.NoError(t, err)
		assert.Equal(t, 98, finalSeats, "Should have deducted 2 seats from first transaction")
		assert.Equal(t, 2, finalVersion, "Version should be incremented once")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
	})

	t.Run("Test booking reference uniqueness", func(t *testing.T) {
		// Create minimal setup for booking
		var scheduleID, customerID string

		// Get or create a schedule
		err := db.Pool.QueryRow(ctx, `
			SELECT id FROM schedules LIMIT 1
		`).Scan(&scheduleID)
		if err != nil {
			// Create minimal schedule
			var opID string
			db.Pool.QueryRow(ctx, `
				INSERT INTO operators (name, code, contact_email)
				VALUES ('Ref Test', 'REF001', 'ref@test.com')
				RETURNING id
			`).Scan(&opID)
			
			// Create schedule (simplified - assuming routes/vessels exist or create them)
			db.Pool.QueryRow(ctx, `
				INSERT INTO schedules (
					operator_id, 
					route_id, 
					vessel_id,
					departure_date, departure_time, arrival_time,
					base_price, total_capacity, available_seats
				)
				SELECT 
					$1,
					(SELECT id FROM routes WHERE operator_id = $1 LIMIT 1),
					(SELECT id FROM vessels WHERE operator_id = $1 LIMIT 1),
					CURRENT_DATE + INTERVAL '3 days', '10:00', '12:00',
					25.00, 50, 50
				RETURNING id
			`, opID).Scan(&scheduleID)
		}

		// Get or create customer
		err = db.Pool.QueryRow(ctx, `
			SELECT id FROM users WHERE user_type = 'customer' LIMIT 1
		`).Scan(&customerID)
		if err != nil {
			db.Pool.QueryRow(ctx, `
				INSERT INTO users (email, password_hash, first_name, last_name, user_type)
				VALUES ('reftest@example.com', '$2a$10$hash', 'Ref', 'Test', 'customer')
				RETURNING id
			`).Scan(&customerID)
		}

		// Create first booking
		uniqueRef := fmt.Sprintf("REF%d", time.Now().UnixNano())
		var booking1ID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO bookings (
				booking_reference, schedule_id, customer_id,
				passenger_count, total_amount, booking_channel
			) VALUES ($1, $2, $3, 1, 25.00, 'online')
			RETURNING id
		`, uniqueRef, scheduleID, customerID).Scan(&booking1ID)
		assert.NoError(t, err)

		// Try to create another booking with same reference
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO bookings (
				booking_reference, schedule_id, customer_id,
				passenger_count, total_amount, booking_channel
			) VALUES ($1, $2, $3, 1, 25.00, 'online')
		`, uniqueRef, scheduleID, customerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key", "Should fail on duplicate booking reference")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM bookings WHERE id = $1", booking1ID)
		assert.NoError(t, err)
	})
}