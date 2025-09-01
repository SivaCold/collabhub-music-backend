package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Keycloak KeycloakConfig
	Storage  StorageConfig
	CORS     CORSConfig
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Host        string
	Port        string
	GinMode     string
	SSLEnabled  bool
	SSLCertPath string
	SSLKeyPath  string
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	Timezone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// KeycloakConfig contains Keycloak integration configuration
type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	UploadPath   string
	MaxFileSize  string
	AllowedTypes []string
}

// CORSConfig contains CORS configuration for frontend integration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// Load loads configuration from environment variables and files
func Load() *Config {
	// Load from environment file based on GO_ENV
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Try to load .env file
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Could not load %s file: %v", envFile, err)
		// Try loading default .env file
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Could not load default .env file: %v", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:        getEnv("SERVER_HOST", "localhost"),
			Port:        getEnv("SERVER_PORT", "8444"),
			GinMode:     getEnv("GIN_MODE", "debug"),
			SSLEnabled:  getBoolEnv("SSL_ENABLED", false),
			SSLCertPath: getEnv("SSL_CERT_PATH", "./certs/server.crt"),
			SSLKeyPath:  getEnv("SSL_KEY_PATH", "./certs/server.key"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "dev_app_user"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "collabhub_music_dev"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			Timezone:        getEnv("DB_TIMEZONE", "UTC"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getIntEnv("DB_CONN_MAX_LIFETIME", 300),
		},
		Keycloak: KeycloakConfig{
			URL:          getEnv("KEYCLOAK_URL", "http://localhost:8080"),
			Realm:        getEnv("KEYCLOAK_REALM", "collabhub"),
			ClientID:     getEnv("KEYCLOAK_CLIENT_ID", "collabhub-backend"),
			ClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		},
		Storage: StorageConfig{
			UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:  getEnv("MAX_FILE_SIZE", "100MB"),
			AllowedTypes: []string{"audio/*", "image/*", "application/pdf"},
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				"http://localhost:8081",  // React Native Metro
				"http://localhost:3000",  // React Web
				"https://localhost:3000", // React Web HTTPS
				"exp://localhost:19000",  // Expo
				"exp://192.168.*:19000",  // Expo LAN
				"*",                      // Allow all origins in development
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		},
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		log.Printf("Configuration validation warning: %v", err)
	}

	return cfg
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if cfg.Keycloak.URL == "" {
		return fmt.Errorf("keycloak URL is required")
	}

	return nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
