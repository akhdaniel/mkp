package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	cfg, err := config.LoadTest()
	require.NoError(t, err, "Failed to load test config")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	t.Run("Run migrations up", func(t *testing.T) {
		migrator, err := NewMigrator(databaseURL)
		require.NoError(t, err, "Failed to create migrator")
		defer migrator.Close()

		err = migrator.Up()
		assert.NoError(t, err, "Failed to run migrations")

		version, dirty, err := migrator.Version()
		assert.NoError(t, err, "Failed to get version")
		assert.False(t, dirty, "Migration should not be dirty")
		assert.Equal(t, uint(1), version, "Should be at version 1")
	})

	t.Run("Verify extensions created", func(t *testing.T) {
		db, err := New(&cfg.Database)
		require.NoError(t, err, "Failed to connect to database")
		defer db.Close()

		ctx := context.Background()

		// Check extensions exist
		extensions := []string{"pgcrypto", "uuid-ossp", "pg_trgm", "btree_gist"}
		for _, ext := range extensions {
			var exists bool
			query := `SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = $1)`
			err := db.Pool.QueryRow(ctx, query, ext).Scan(&exists)
			assert.NoError(t, err, "Failed to check extension %s", ext)
			assert.True(t, exists, "Extension %s should exist", ext)
		}
	})

	t.Run("Verify schemas created", func(t *testing.T) {
		db, err := New(&cfg.Database)
		require.NoError(t, err, "Failed to connect to database")
		defer db.Close()

		ctx := context.Background()

		// Check audit schema exists
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'audit')`
		err = db.Pool.QueryRow(ctx, query).Scan(&exists)
		assert.NoError(t, err, "Failed to check audit schema")
		assert.True(t, exists, "Audit schema should exist")
	})
}