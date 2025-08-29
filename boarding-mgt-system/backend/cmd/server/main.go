package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/api"
	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/ferryflow/boarding-mgt-system/internal/database"
	_ "github.com/ferryflow/boarding-mgt-system/docs" // Swagger docs
)

// @title FerryFlow Boarding Management API
// @version 1.0
// @description Complete API for ferry boarding management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ferryflow.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	ctx := context.Background()
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	if err := database.RunMigrations(ctx, databaseURL); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Check database extensions
	if err := db.CheckExtensions(ctx); err != nil {
		log.Printf("Warning: Failed to check extensions: %v", err)
	}

	// Create API server
	server := api.NewServer(cfg, db)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      server.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s:%d", cfg.App.Host, cfg.App.Port)
		log.Printf("Swagger documentation available at http://%s:%d/swagger/index.html", cfg.App.Host, cfg.App.Port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}