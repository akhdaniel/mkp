package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/ferryflow/boarding-mgt-system/internal/database"
)

func main() {
	var (
		direction string
		steps     int
		version   int
		force     bool
	)

	flag.StringVar(&direction, "direction", "", "Migration direction: up, down, or force")
	flag.IntVar(&steps, "steps", 0, "Number of migrations to run (default: all)")
	flag.IntVar(&version, "version", 0, "Migrate to specific version")
	flag.BoolVar(&force, "force", false, "Force set version without running migrations")
	flag.Parse()

	// Also accept direction as first argument for convenience
	if direction == "" && flag.NArg() > 0 {
		direction = flag.Arg(0)
	}

	if direction == "" {
		fmt.Println("Usage: migrate [up|down|version|force] [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -steps int    Number of migrations to run")
		fmt.Println("  -version int  Migrate to specific version")
		fmt.Println("  -force        Force set version without running migrations")
		fmt.Println("\nExamples:")
		fmt.Println("  migrate up                  # Run all pending migrations")
		fmt.Println("  migrate down -steps 1       # Rollback one migration")
		fmt.Println("  migrate version            # Show current version")
		fmt.Println("  migrate force -version 3   # Force set to version 3")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build database URL
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Create migrator
	migrator, err := database.NewMigrator(databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	switch direction {
	case "up":
		fmt.Println("Running migrations...")
		if steps > 0 {
			for i := 0; i < steps; i++ {
				if err := migrator.Steps(1); err != nil {
					log.Fatalf("Failed to run migration: %v", err)
				}
			}
			fmt.Printf("Successfully ran %d migration(s)\n", steps)
		} else {
			if err := migrator.Up(); err != nil {
				log.Fatalf("Failed to run migrations: %v", err)
			}
			fmt.Println("All migrations completed successfully")
		}

	case "down":
		if steps > 0 {
			fmt.Printf("Rolling back %d migration(s)...\n", steps)
			for i := 0; i < steps; i++ {
				if err := migrator.Steps(-1); err != nil {
					log.Fatalf("Failed to rollback migration: %v", err)
				}
			}
			fmt.Printf("Successfully rolled back %d migration(s)\n", steps)
		} else {
			fmt.Println("Rolling back all migrations...")
			if err := migrator.Down(); err != nil {
				log.Fatalf("Failed to rollback migrations: %v", err)
			}
			fmt.Println("All migrations rolled back successfully")
		}

	case "version":
		v, dirty, err := migrator.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current version: %d", v)
		if dirty {
			fmt.Print(" (dirty)")
		}
		fmt.Println()

	case "force":
		if version == 0 {
			log.Fatal("Version number required for force command")
		}
		fmt.Printf("Forcing version to %d...\n", version)
		if err := migrator.Force(version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		fmt.Printf("Version forced to %d\n", version)

	case "create":
		// Get migration name from next argument
		if flag.NArg() < 2 {
			log.Fatal("Migration name required for create command")
		}
		name := flag.Arg(1)
		if err := createMigration(name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

	default:
		log.Fatalf("Unknown direction: %s", direction)
	}

	// Verify final state
	if direction != "version" {
		v, dirty, err := migrator.Version()
		if err != nil {
			fmt.Printf("Warning: Could not verify final version: %v\n", err)
		} else {
			fmt.Printf("Current version: %d", v)
			if dirty {
				fmt.Print(" (dirty)")
			}
			fmt.Println()
		}
	}
}

func createMigration(name string) error {
	timestamp := time.Now().Unix()
	upFile := fmt.Sprintf("backend/internal/database/migrations/%06d_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("backend/internal/database/migrations/%06d_%s.down.sql", timestamp, name)

	// Create up migration file
	if err := os.WriteFile(upFile, []byte("-- Add your UP migration here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}

	// Create down migration file
	if err := os.WriteFile(downFile, []byte("-- Add your DOWN migration here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)
	return nil
}