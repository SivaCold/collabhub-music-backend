package config

import (
    "os"
    "strconv"
)

type Config struct {
    Database   DatabaseConfig
    Keycloak   KeycloakConfig
    Server     ServerConfig
    Storage    StorageConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime int
	Timezone string
}

type KeycloakConfig struct {
    URL      string
    Realm    string
    ClientID string
	ClientSecret string
}

type ServerConfig struct {
    Port        string
    GinMode     string
    SSLCertPath string
    SSLKeyPath  string
}

type StorageConfig struct {
    UploadPath  string
    MaxFileSize int64
}

func Load() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "collabhub_music"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        Keycloak: KeycloakConfig{
            URL:      getEnv("KEYCLOAK_URL", "http://localhost:8080"),
            Realm:    getEnv("KEYCLOAK_REALM", "collabhub"),
            ClientID: getEnv("KEYCLOAK_CLIENT_ID", "collabhub-backend"),
        },
        Server: ServerConfig{
            Port:        getEnv("PORT", "8443"),
            GinMode:     getEnv("GIN_MODE", "debug"),
            SSLCertPath: getEnv("SSL_CERT_PATH", "./certs/server.crt"),
            SSLKeyPath:  getEnv("SSL_KEY_PATH", "./certs/server.key"),
        },
        Storage: StorageConfig{
            UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
            MaxFileSize: parseFileSize(getEnv("MAX_FILE_SIZE", "100MB")),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func parseFileSize(size string) int64 {
    // Simple parser for sizes like "100MB"
    if len(size) < 3 {
        return 100 * 1024 * 1024 // Default 100MB
    }
    
    unit := size[len(size)-2:]
    valueStr := size[:len(size)-2]
    value, err := strconv.ParseInt(valueStr, 10, 64)
    if err != nil {
        return 100 * 1024 * 1024
    }
    
    switch unit {
    case "KB":
        return value * 1024
    case "MB":
        return value * 1024 * 1024
    case "GB":
        return value * 1024 * 1024 * 1024
    default:
        return value
    }
}