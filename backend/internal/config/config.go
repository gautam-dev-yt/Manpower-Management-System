package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration, loaded from environment variables.
type Config struct {
	Port      string
	DB        DBConfig
	JWTSecret string
	Upload    UploadConfig
}

// DBConfig holds PostgreSQL connection details.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// UploadConfig holds file upload settings.
type UploadConfig struct {
	Dir     string // Local directory for file uploads
	BaseURL string // URL prefix for serving uploaded files
}

// Load reads configuration from environment variables (with .env fallback).
func Load() (*Config, error) {
	// Load .env file for local development — silently ignored in production
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "manpower_dev"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Upload: UploadConfig{
			Dir:     getEnv("UPLOAD_DIR", "./uploads"),
			BaseURL: "", // Set below after port is known
		},
	}

	// Build the file serving URL from the port
	cfg.Upload.BaseURL = getEnv("UPLOAD_BASE_URL",
		fmt.Sprintf("http://localhost:%s/api/files", cfg.Port),
	)

	// Required fields
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
