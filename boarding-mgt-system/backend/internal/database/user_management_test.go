package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserManagementTables(t *testing.T) {
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

	t.Run("Verify users table structure", func(t *testing.T) {
		// Check table exists
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'users'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "users table should exist")

		// Verify columns
		rows, err := db.Pool.Query(ctx, `
			SELECT column_name, data_type, is_nullable
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = 'users'
			ORDER BY ordinal_position
		`)
		assert.NoError(t, err)
		defer rows.Close()

		expectedColumns := map[string]string{
			"id":            "uuid",
			"email":         "character varying",
			"password_hash": "character varying",
			"first_name":    "character varying",
			"last_name":     "character varying",
			"phone":         "character varying",
			"date_of_birth": "date",
			"nationality":   "character varying",
			"user_type":     "character varying",
			"operator_id":   "uuid",
			"is_verified":   "boolean",
			"is_active":     "boolean",
			"last_login_at": "timestamp with time zone",
			"created_at":    "timestamp with time zone",
			"updated_at":    "timestamp with time zone",
		}

		columnCount := 0
		for rows.Next() {
			var columnName, dataType, isNullable string
			err := rows.Scan(&columnName, &dataType, &isNullable)
			assert.NoError(t, err)

			expectedType, exists := expectedColumns[columnName]
			assert.True(t, exists, "Column %s should be expected", columnName)
			assert.Equal(t, expectedType, dataType, "Column %s should have correct type", columnName)
			columnCount++
		}
		assert.Equal(t, len(expectedColumns), columnCount, "Should have all expected columns")

		// Verify unique constraint on email
		var constraintExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'users_email_key'
			)
		`).Scan(&constraintExists)
		assert.NoError(t, err)
		assert.True(t, constraintExists, "Unique constraint on email should exist")

		// Verify check constraint on user_type
		var checkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'valid_user_type'
				AND contype = 'c'
			)
		`).Scan(&checkExists)
		assert.NoError(t, err)
		assert.True(t, checkExists, "Check constraint on user_type should exist")
	})

	t.Run("Verify user_sessions table structure", func(t *testing.T) {
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'user_sessions'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "user_sessions table should exist")

		// Verify foreign key to users
		var fkExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.table_constraints 
				WHERE constraint_type = 'FOREIGN KEY' 
				AND table_name = 'user_sessions'
				AND constraint_name LIKE '%user%'
			)
		`).Scan(&fkExists)
		assert.NoError(t, err)
		assert.True(t, fkExists, "Foreign key to users should exist")

		// Verify INET type for ip_address
		var dataType string
		err = db.Pool.QueryRow(ctx, `
			SELECT udt_name
			FROM information_schema.columns
			WHERE table_schema = 'public' 
			AND table_name = 'user_sessions'
			AND column_name = 'ip_address'
		`).Scan(&dataType)
		assert.NoError(t, err)
		assert.Equal(t, "inet", dataType, "ip_address should be INET type")
	})

	t.Run("Test user registration workflow", func(t *testing.T) {
		// Clean up any existing test users
		_, err := db.Pool.Exec(ctx, "DELETE FROM users WHERE email LIKE 'test%@example.com'")
		require.NoError(t, err)

		// Register a customer
		var customerID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name, 
				phone, user_type, is_verified
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`, "testcustomer@example.com", "$2a$10$hashedpassword", "Test", "Customer",
			"+1234567890", "customer", false).Scan(&customerID)
		assert.NoError(t, err, "Should register customer")
		assert.NotEmpty(t, customerID)

		// Register an agent with operator
		var operatorID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Test Operator", "TOP001", "operator@example.com").Scan(&operatorID)
		require.NoError(t, err)

		var agentID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name,
				user_type, operator_id, is_verified
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`, "testagent@example.com", "$2a$10$hashedpassword", "Test", "Agent",
			"agent", operatorID, true).Scan(&agentID)
		assert.NoError(t, err, "Should register agent")
		assert.NotEmpty(t, agentID)

		// Test unique email constraint
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name, user_type
			) VALUES ($1, $2, $3, $4, $5)
		`, "testcustomer@example.com", "$2a$10$newpassword", "Another", "User", "customer")
		assert.Error(t, err, "Should fail on duplicate email")
		assert.Contains(t, err.Error(), "duplicate key", "Error should mention duplicate key")

		// Test invalid user_type
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name, user_type
			) VALUES ($1, $2, $3, $4, $5)
		`, "testinvalid@example.com", "$2a$10$hashedpassword", "Invalid", "Type", "invalid_type")
		assert.Error(t, err, "Should fail on invalid user_type")
		assert.Contains(t, err.Error(), "valid_user_type", "Error should mention user_type constraint")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", customerID)
		assert.NoError(t, err)
	})

	t.Run("Test session management", func(t *testing.T) {
		// Create a test user
		var userID string
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name, user_type
			) VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "testsession@example.com", "$2a$10$hashedpassword", "Session", "User", "customer").Scan(&userID)
		require.NoError(t, err)

		// Create a session
		var sessionID string
		expiresAt := time.Now().Add(24 * time.Hour)
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO user_sessions (
				user_id, token_hash, refresh_token_hash, 
				expires_at, ip_address, user_agent
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, userID, "tokenHash123", "refreshHash456",
			expiresAt, "192.168.1.1", "Mozilla/5.0").Scan(&sessionID)
		assert.NoError(t, err, "Should create session")
		assert.NotEmpty(t, sessionID)

		// Query active sessions for user
		var activeCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM user_sessions 
			WHERE user_id = $1 AND is_active = true
		`, userID).Scan(&activeCount)
		assert.NoError(t, err)
		assert.Equal(t, 1, activeCount, "Should have one active session")

		// Deactivate session
		_, err = db.Pool.Exec(ctx, `
			UPDATE user_sessions SET is_active = false 
			WHERE id = $1
		`, sessionID)
		assert.NoError(t, err)

		// Test cascade delete
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
		assert.NoError(t, err)

		// Verify session was deleted via cascade
		var sessionCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM user_sessions WHERE user_id = $1
		`, userID).Scan(&sessionCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, sessionCount, "Sessions should be deleted via cascade")
	})

	t.Run("Test multi-tenant user isolation", func(t *testing.T) {
		// Create two operators
		var operator1ID, operator2ID string
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Operator One", "OP001", "op1@example.com").Scan(&operator1ID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Operator Two", "OP002", "op2@example.com").Scan(&operator2ID)
		require.NoError(t, err)

		// Create users for each operator
		var user1ID, user2ID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name,
				user_type, operator_id
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "user1@op1.com", "$2a$10$hash", "User", "One", "agent", operator1ID).Scan(&user1ID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name,
				user_type, operator_id
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "user2@op2.com", "$2a$10$hash", "User", "Two", "agent", operator2ID).Scan(&user2ID)
		require.NoError(t, err)

		// Query users by operator
		var count1, count2 int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM users WHERE operator_id = $1
		`, operator1ID).Scan(&count1)
		assert.NoError(t, err)
		assert.Equal(t, 1, count1, "Operator 1 should have 1 user")

		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM users WHERE operator_id = $2
		`, operator2ID).Scan(&count2)
		assert.NoError(t, err)
		assert.Equal(t, 1, count2, "Operator 2 should have 1 user")

		// Create a customer (no operator)
		var customerID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (
				email, password_hash, first_name, last_name,
				user_type, operator_id
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "customer@public.com", "$2a$10$hash", "Public", "Customer", "customer", nil).Scan(&customerID)
		assert.NoError(t, err, "Should create customer without operator")

		// Verify customer has no operator
		var operatorIDNull *string
		err = db.Pool.QueryRow(ctx, `
			SELECT operator_id FROM users WHERE id = $1
		`, customerID).Scan(&operatorIDNull)
		assert.NoError(t, err)
		assert.Nil(t, operatorIDNull, "Customer should have NULL operator_id")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id IN ($1, $2)", operator1ID, operator2ID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", customerID)
		assert.NoError(t, err)
	})
}