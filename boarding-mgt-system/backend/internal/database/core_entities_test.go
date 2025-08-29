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

func TestCoreEntityTables(t *testing.T) {
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

	t.Run("Verify operators table structure", func(t *testing.T) {
		// Check table exists
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'operators'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "operators table should exist")

		// Verify columns
		rows, err := db.Pool.Query(ctx, `
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = 'operators'
			ORDER BY ordinal_position
		`)
		assert.NoError(t, err)
		defer rows.Close()

		expectedColumns := map[string]string{
			"id":            "uuid",
			"name":          "character varying",
			"code":          "character varying",
			"contact_email": "character varying",
			"contact_phone": "character varying",
			"address":       "text",
			"is_active":     "boolean",
			"settings":      "jsonb",
			"created_at":    "timestamp with time zone",
			"updated_at":    "timestamp with time zone",
		}

		columnCount := 0
		for rows.Next() {
			var columnName, dataType, isNullable string
			var columnDefault *string
			err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)
			assert.NoError(t, err)

			expectedType, exists := expectedColumns[columnName]
			assert.True(t, exists, "Column %s should be expected", columnName)
			assert.Equal(t, expectedType, dataType, "Column %s should have correct type", columnName)
			columnCount++
		}
		assert.Equal(t, len(expectedColumns), columnCount, "Should have all expected columns")

		// Verify unique constraint on code
		var constraintExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'operators_code_key'
			)
		`).Scan(&constraintExists)
		assert.NoError(t, err)
		assert.True(t, constraintExists, "Unique constraint on code should exist")
	})

	t.Run("Verify ports table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'ports'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "ports table should exist")

		// Verify POINT column for coordinates
		var dataType string
		err = db.Pool.QueryRow(ctx, `
			SELECT udt_name
			FROM information_schema.columns
			WHERE table_schema = 'public' 
			AND table_name = 'ports'
			AND column_name = 'coordinates'
		`).Scan(&dataType)
		assert.NoError(t, err)
		assert.Equal(t, "point", dataType, "coordinates should be POINT type")
	})

	t.Run("Verify vessels table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'vessels'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "vessels table should exist")

		// Verify foreign key to operators
		var fkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.table_constraints 
				WHERE constraint_type = 'FOREIGN KEY' 
				AND table_name = 'vessels'
				AND constraint_name LIKE '%operator%'
			)
		`).Scan(&fkExists)
		assert.NoError(t, err)
		assert.True(t, fkExists, "Foreign key to operators should exist")

		// Verify check constraint on capacity
		var checkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'vessels_capacity_check'
				AND contype = 'c'
			)
		`).Scan(&checkExists)
		assert.NoError(t, err)
		assert.True(t, checkExists, "Check constraint on capacity should exist")
	})

	t.Run("Verify routes table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'routes'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "routes table should exist")

		// Verify check constraint for different ports
		var checkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'different_ports'
				AND contype = 'c'
			)
		`).Scan(&checkExists)
		assert.NoError(t, err)
		assert.True(t, checkExists, "Check constraint for different ports should exist")

		// Verify foreign keys
		var fkCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM information_schema.table_constraints 
			WHERE constraint_type = 'FOREIGN KEY' 
			AND table_name = 'routes'
		`).Scan(&fkCount)
		assert.NoError(t, err)
		assert.Equal(t, 3, fkCount, "Routes should have 3 foreign keys (operator, departure_port, arrival_port)")
	})

	t.Run("Test inserting and querying operators", func(t *testing.T) {
		// Insert test operator
		var operatorID string
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email, contact_phone, address)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "Test Ferry Co", "TFC001", "test@ferry.com", "+1234567890", "123 Ferry St").Scan(&operatorID)
		assert.NoError(t, err, "Should insert operator")
		assert.NotEmpty(t, operatorID, "Should return operator ID")

		// Query operator
		var name, code string
		err = db.Pool.QueryRow(ctx, `
			SELECT name, code FROM operators WHERE id = $1
		`, operatorID).Scan(&name, &code)
		assert.NoError(t, err)
		assert.Equal(t, "Test Ferry Co", name)
		assert.Equal(t, "TFC001", code)

		// Test unique constraint
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
		`, "Another Ferry", "TFC001", "another@ferry.com")
		assert.Error(t, err, "Should fail on duplicate code")
		assert.Contains(t, err.Error(), "duplicate key", "Error should mention duplicate key")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
	})

	t.Run("Test inserting and querying complete route", func(t *testing.T) {
		// Insert operator
		var operatorID string
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Route Test Ferry", "RTF001", "route@ferry.com").Scan(&operatorID)
		require.NoError(t, err)

		// Insert ports
		var port1ID, port2ID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO ports (name, code, city, country, timezone)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "Port Alpha", "PA001", "City A", "Country A", "UTC").Scan(&port1ID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO ports (name, code, city, country, timezone, coordinates)
			VALUES ($1, $2, $3, $4, $5, POINT($6, $7))
			RETURNING id
		`, "Port Beta", "PB001", "City B", "Country B", "UTC", 40.7128, -74.0060).Scan(&port2ID)
		require.NoError(t, err)

		// Insert vessel
		var vesselID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO vessels (operator_id, name, registration_number, vessel_type, capacity, seat_configuration)
			VALUES ($1, $2, $3, $4, $5, $6::jsonb)
			RETURNING id
		`, operatorID, "Test Vessel", "TV001", "passenger", 100, `{"decks": 2}`).Scan(&vesselID)
		require.NoError(t, err)

		// Insert route
		var routeID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, distance_km, estimated_duration)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, operatorID, "Alpha to Beta", port1ID, port2ID, 150.5, "2 hours").Scan(&routeID)
		assert.NoError(t, err, "Should create route")
		assert.NotEmpty(t, routeID)

		// Test constraint - same port for departure and arrival
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, estimated_duration)
			VALUES ($1, $2, $3, $4, $5)
		`, operatorID, "Invalid Route", port1ID, port1ID, "1 hour")
		assert.Error(t, err, "Should fail when departure and arrival ports are the same")
		assert.Contains(t, err.Error(), "different_ports", "Error should mention the constraint")

		// Cleanup (cascade delete from operator)
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
		
		// Verify cascade delete worked
		var count int
		err = db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM vessels WHERE operator_id = $1", operatorID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Vessels should be deleted via cascade")

		err = db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM routes WHERE operator_id = $1", operatorID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Routes should be deleted via cascade")

		// Cleanup ports
		_, err = db.Pool.Exec(ctx, "DELETE FROM ports WHERE id IN ($1, $2)", port1ID, port2ID)
		assert.NoError(t, err)
	})
}