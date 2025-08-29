package database

import (
	"context"
	"testing"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	cfg, err := config.LoadTest()
	require.NoError(t, err, "Failed to load test config")

	ctx := context.Background()

	t.Run("Connect to database", func(t *testing.T) {
		db, err := New(&cfg.Database)
		require.NoError(t, err, "Failed to connect to database")
		defer db.Close()

		err = db.Ping(ctx)
		assert.NoError(t, err, "Failed to ping database")
	})

	t.Run("Check connection pool", func(t *testing.T) {
		db, err := New(&cfg.Database)
		require.NoError(t, err, "Failed to connect to database")
		defer db.Close()

		stats := db.Pool.Stat()
		assert.GreaterOrEqual(t, stats.MaxConns, int32(5), "Max connections should be at least 5")
		assert.GreaterOrEqual(t, stats.IdleConns, int32(0), "Should have idle connections")
	})

	t.Run("Connection timeout", func(t *testing.T) {
		invalidCfg := &config.DatabaseConfig{
			Host:     "invalid-host",
			Port:     5432,
			Name:     "test",
			User:     "test",
			Password: "test",
			SSLMode:  "disable",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := New(invalidCfg)
		assert.Error(t, err, "Should fail with invalid host")
	})
}

func TestExtensions(t *testing.T) {
	cfg, err := config.LoadTest()
	require.NoError(t, err, "Failed to load test config")

	db, err := New(&cfg.Database)
	require.NoError(t, err, "Failed to connect to database")
	defer db.Close()

	ctx := context.Background()

	t.Run("Check required extensions", func(t *testing.T) {
		err := db.CheckExtensions(ctx)
		assert.NoError(t, err, "Failed to check/create extensions")
	})

	t.Run("Verify pgcrypto functions", func(t *testing.T) {
		var result string
		query := `SELECT encode(digest('test', 'sha256'), 'hex')`
		err := db.Pool.QueryRow(ctx, query).Scan(&result)
		assert.NoError(t, err, "pgcrypto should be available")
		assert.NotEmpty(t, result, "Should return hash")
	})

	t.Run("Verify UUID generation", func(t *testing.T) {
		var uuid string
		query := `SELECT gen_random_uuid()::text`
		err := db.Pool.QueryRow(ctx, query).Scan(&uuid)
		assert.NoError(t, err, "UUID generation should work")
		assert.Len(t, uuid, 36, "UUID should be 36 characters")
	})

	t.Run("Verify pg_trgm functions", func(t *testing.T) {
		var similarity float64
		query := `SELECT similarity('test', 'text')`
		err := db.Pool.QueryRow(ctx, query).Scan(&similarity)
		assert.NoError(t, err, "pg_trgm should be available")
		assert.Greater(t, similarity, 0.0, "Should return similarity score")
	})
}