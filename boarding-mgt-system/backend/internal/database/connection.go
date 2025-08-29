package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.DatabaseConfig) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *DB) CheckExtensions(ctx context.Context) error {
	requiredExtensions := []string{
		"pgcrypto",
		"uuid-ossp", 
		"pg_trgm",
		"btree_gist",
	}

	for _, ext := range requiredExtensions {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = $1)`
		err := db.Pool.QueryRow(ctx, query, ext).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check extension %s: %w", ext, err)
		}
		if !exists {
			// Try to create the extension
			if err := db.CreateExtension(ctx, ext); err != nil {
				return fmt.Errorf("extension %s is not available: %w", ext, err)
			}
		}
	}

	return nil
}

func (db *DB) CreateExtension(ctx context.Context, name string) error {
	query := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", name)
	_, err := db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create extension %s: %w", name, err)
	}
	return nil
}

func (db *DB) CreateDatabase(ctx context.Context, name string) error {
	// Connect to postgres database to create new database
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "postgres",
		User:     "ferryflow",
		Password: "ferryflow_dev_2024",
		SSLMode:  "disable",
	}

	conn, err := pgx.Connect(ctx, cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer conn.Close(ctx)

	// Check if database exists
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		query := fmt.Sprintf("CREATE DATABASE %s OWNER ferryflow", name)
		_, err = conn.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}