package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type AppConfig struct {
	Environment string
	Port        int
	Host        string
}

type JWTConfig struct {
	Secret string
	Expiry string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     getEnv("DB_NAME", "ferryflow_dev"),
			User:     getEnv("DB_USER", "ferryflow"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        appPort,
			Host:        getEnv("APP_HOST", "localhost"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
	}, nil
}

func LoadTest() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	dbPort, err := strconv.Atoi(getEnv("TEST_DB_PORT", "5433"))
	if err != nil {
		return nil, fmt.Errorf("invalid TEST_DB_PORT: %w", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     getEnv("TEST_DB_NAME", "ferryflow_test"),
			User:     getEnv("TEST_DB_USER", "ferryflow"),
			Password: getEnv("TEST_DB_PASSWORD", ""),
			SSLMode:  getEnv("TEST_DB_SSL_MODE", "disable"),
		},
		App: AppConfig{
			Environment: "test",
			Port:        8081,
			Host:        "localhost",
		},
		JWT: JWTConfig{
			Secret: "test-secret",
			Expiry: "1h",
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}