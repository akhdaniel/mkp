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

func TestSupportAndAuditTables(t *testing.T) {
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

	t.Run("Verify support_tickets table structure", func(t *testing.T) {
		// Check table exists
		var exists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename = 'support_tickets'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "support_tickets table should exist")

		// Verify unique constraint on ticket_number
		var constraintExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'support_tickets_ticket_number_key'
			)
		`).Scan(&constraintExists)
		assert.NoError(t, err)
		assert.True(t, constraintExists, "Unique constraint on ticket_number should exist")

		// Verify check constraints
		var priorityCheck, statusCheck bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'valid_priority' AND contype = 'c'
			)
		`).Scan(&priorityCheck)
		assert.NoError(t, err)
		assert.True(t, priorityCheck, "Check constraint on priority should exist")

		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'valid_ticket_status' AND contype = 'c'
			)
		`).Scan(&statusCheck)
		assert.NoError(t, err)
		assert.True(t, statusCheck, "Check constraint on status should exist")
	})

	t.Run("Verify audit schema and table", func(t *testing.T) {
		// Check audit schema exists
		var schemaExists bool
		err := db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.schemata 
				WHERE schema_name = 'audit'
			)
		`).Scan(&schemaExists)
		assert.NoError(t, err)
		assert.True(t, schemaExists, "audit schema should exist")

		// Check audit_logs table exists
		var tableExists bool
		err = db.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM pg_tables 
				WHERE schemaname = 'audit' 
				AND tablename = 'audit_logs'
			)
		`).Scan(&tableExists)
		assert.NoError(t, err)
		assert.True(t, tableExists, "audit.audit_logs table should exist")
	})

	t.Run("Test support ticket workflow", func(t *testing.T) {
		// Create customer and agent users
		var customerID, agentID, bookingID string
		
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, user_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "support_customer@example.com", "$2a$10$hash", "Support", "Customer", "customer").Scan(&customerID)
		require.NoError(t, err)

		// Create agent with operator
		var operatorID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Support Test Op", "STO001", "support@op.com").Scan(&operatorID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, user_type, operator_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "support_agent@example.com", "$2a$10$hash", "Support", "Agent", "agent", operatorID).Scan(&agentID)
		require.NoError(t, err)

		// Create a support ticket
		var ticketID string
		ticketNumber := fmt.Sprintf("ST%d", time.Now().Unix())
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO support_tickets (
				ticket_number, customer_id, subject, description, priority
			) VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, ticketNumber, customerID, "Booking Issue", "Cannot complete booking", "high").Scan(&ticketID)
		assert.NoError(t, err, "Should create support ticket")

		// Assign agent to ticket
		_, err = db.Pool.Exec(ctx, `
			UPDATE support_tickets 
			SET assigned_agent_id = $1, status = 'in_progress'
			WHERE id = $2
		`, agentID, ticketID)
		assert.NoError(t, err, "Should assign agent to ticket")

		// Resolve ticket
		_, err = db.Pool.Exec(ctx, `
			UPDATE support_tickets 
			SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`, ticketID)
		assert.NoError(t, err, "Should resolve ticket")

		// Verify ticket status
		var status string
		var resolvedAt *time.Time
		err = db.Pool.QueryRow(ctx, `
			SELECT status, resolved_at FROM support_tickets WHERE id = $1
		`, ticketID).Scan(&status, &resolvedAt)
		assert.NoError(t, err)
		assert.Equal(t, "resolved", status)
		assert.NotNil(t, resolvedAt)

		// Test priority escalation
		var escalatedTicketID string
		escalatedNumber := fmt.Sprintf("ST%d-ESC", time.Now().Unix())
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO support_tickets (
				ticket_number, customer_id, subject, description, priority
			) VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, escalatedNumber, customerID, "Urgent Issue", "Payment failed multiple times", "normal").Scan(&escalatedTicketID)
		require.NoError(t, err)

		// Escalate priority
		_, err = db.Pool.Exec(ctx, `
			UPDATE support_tickets 
			SET priority = 'urgent'
			WHERE id = $1
		`, escalatedTicketID)
		assert.NoError(t, err, "Should escalate priority")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", customerID)
		assert.NoError(t, err)
	})

	t.Run("Test audit logging", func(t *testing.T) {
		// Create a test operator for auditing
		var operatorID string
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Audit Test Op", "ATO001", "audit@test.com").Scan(&operatorID)
		require.NoError(t, err)

		// Check if audit trigger captured the insert
		var auditCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs 
			WHERE table_name = 'operators' 
			AND record_id = $1
			AND action = 'INSERT'
		`, operatorID).Scan(&auditCount)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, auditCount, 1, "Should have audit log for insert")

		// Update the operator
		_, err = db.Pool.Exec(ctx, `
			UPDATE operators 
			SET contact_email = 'newemail@test.com'
			WHERE id = $1
		`, operatorID)
		assert.NoError(t, err)

		// Check if audit trigger captured the update
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs 
			WHERE table_name = 'operators' 
			AND record_id = $1
			AND action = 'UPDATE'
		`, operatorID).Scan(&auditCount)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, auditCount, 1, "Should have audit log for update")

		// Verify audit log contains old and new values
		var oldValues, newValues string
		err = db.Pool.QueryRow(ctx, `
			SELECT old_values::text, new_values::text 
			FROM audit.audit_logs 
			WHERE table_name = 'operators' 
			AND record_id = $1
			AND action = 'UPDATE'
			ORDER BY changed_at DESC
			LIMIT 1
		`, operatorID).Scan(&oldValues, &newValues)
		assert.NoError(t, err)
		assert.Contains(t, oldValues, "audit@test.com", "Old values should contain original email")
		assert.Contains(t, newValues, "newemail@test.com", "New values should contain updated email")

		// Delete the operator
		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", operatorID)
		assert.NoError(t, err)

		// Check if audit trigger captured the delete
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs 
			WHERE table_name = 'operators' 
			AND record_id = $1
			AND action = 'DELETE'
		`, operatorID).Scan(&auditCount)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, auditCount, 1, "Should have audit log for delete")
	})

	t.Run("Test support ticket messages", func(t *testing.T) {
		// Create customer and support ticket
		var customerID, ticketID string
		
		err := db.Pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, user_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "msg_customer@example.com", "$2a$10$hash", "Message", "Customer", "customer").Scan(&customerID)
		require.NoError(t, err)

		ticketNumber := fmt.Sprintf("MSG%d", time.Now().Unix())
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO support_tickets (
				ticket_number, customer_id, subject, description, priority
			) VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, ticketNumber, customerID, "Test Message", "Testing message system", "normal").Scan(&ticketID)
		require.NoError(t, err)

		// Add messages to ticket
		var message1ID, message2ID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO support_messages (
				ticket_id, sender_id, message, is_internal
			) VALUES ($1, $2, $3, $4)
			RETURNING id
		`, ticketID, customerID, "Initial customer message", false).Scan(&message1ID)
		assert.NoError(t, err, "Should create customer message")

		// Create agent for internal message
		var agentID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, user_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, "msg_agent@example.com", "$2a$10$hash", "Message", "Agent", "agent").Scan(&agentID)
		require.NoError(t, err)

		err = db.Pool.QueryRow(ctx, `
			INSERT INTO support_messages (
				ticket_id, sender_id, message, is_internal
			) VALUES ($1, $2, $3, $4)
			RETURNING id
		`, ticketID, agentID, "Internal note for team", true).Scan(&message2ID)
		assert.NoError(t, err, "Should create internal message")

		// Verify message count
		var publicCount, internalCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT 
				COUNT(*) FILTER (WHERE is_internal = false) as public,
				COUNT(*) FILTER (WHERE is_internal = true) as internal
			FROM support_messages 
			WHERE ticket_id = $1
		`, ticketID).Scan(&publicCount, &internalCount)
		assert.NoError(t, err)
		assert.Equal(t, 1, publicCount, "Should have 1 public message")
		assert.Equal(t, 1, internalCount, "Should have 1 internal message")

		// Cleanup
		_, err = db.Pool.Exec(ctx, "DELETE FROM support_tickets WHERE id = $1", ticketID)
		assert.NoError(t, err)
		_, err = db.Pool.Exec(ctx, "DELETE FROM users WHERE id IN ($1, $2)", customerID, agentID)
		assert.NoError(t, err)
	})

	t.Run("Test audit log retention", func(t *testing.T) {
		// Get current audit log count
		var beforeCount int
		err := db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs
		`).Scan(&beforeCount)
		assert.NoError(t, err)

		// Create and delete a test record to generate audit logs
		var testID string
		err = db.Pool.QueryRow(ctx, `
			INSERT INTO operators (name, code, contact_email)
			VALUES ('Retention Test', 'RET001', 'ret@test.com')
			RETURNING id
		`).Scan(&testID)
		require.NoError(t, err)

		_, err = db.Pool.Exec(ctx, "DELETE FROM operators WHERE id = $1", testID)
		assert.NoError(t, err)

		// Verify audit logs were created
		var afterCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs
		`).Scan(&afterCount)
		assert.NoError(t, err)
		assert.Greater(t, afterCount, beforeCount, "Should have more audit logs after operations")

		// Test that we can query audit logs by date range
		var recentCount int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM audit.audit_logs
			WHERE changed_at >= CURRENT_TIMESTAMP - INTERVAL '1 minute'
		`).Scan(&recentCount)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, recentCount, 2, "Should have recent audit logs")
	})
}