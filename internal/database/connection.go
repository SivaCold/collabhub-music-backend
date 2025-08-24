package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "collabhub-music-backend/internal/config"
    _ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the database
func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
    // Validate configuration
    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("invalid database configuration: %w", err)
    }

    dsn := buildDSN(cfg)
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %w", err)
    }

    // Configure connection pool with defaults if not set
    configureConnectionPool(db, cfg)

    // Test the connection with retry logic
    if err := testConnection(db); err != nil {
        db.Close()
        return nil, fmt.Errorf("error testing database connection: %w", err)
    }

    DB = db
    log.Printf("Successfully connected to database: %s on %s:%s", cfg.DBName, cfg.Host, cfg.Port)
    return db, nil
}

// Close closes the database connection
func Close(db *sql.DB) error {
    if db == nil {
        return nil
    }
    
    if err := db.Close(); err != nil {
        return fmt.Errorf("error closing database connection: %w", err)
    }
    
    log.Println("Database connection closed successfully")
    return nil
}

// GetDB returns the global database instance
func GetDB() *sql.DB {
    return DB
}

// validateConfig validates the database configuration
func validateConfig(cfg config.DatabaseConfig) error {
    if cfg.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if cfg.Port == "" {
        return fmt.Errorf("database port is required")
    }
    if cfg.User == "" {
        return fmt.Errorf("database user is required")
    }
    if cfg.DBName == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}

// buildDSN constructs the database connection string
func buildDSN(cfg config.DatabaseConfig) string {
    sslMode := cfg.SSLMode
    if sslMode == "" {
        sslMode = "require" // Default to secure connection
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
        cfg.Host,
        cfg.Port,
        cfg.User,
        cfg.DBName,
        sslMode,
    )

    // Add password if provided
    if cfg.Password != "" {
        dsn += fmt.Sprintf(" password=%s", cfg.Password)
    }

    // Add timezone if provided
    if cfg.Timezone != "" {
        dsn += fmt.Sprintf(" timezone=%s", cfg.Timezone)
    }

    return dsn
}

// configureConnectionPool sets up the database connection pool
func configureConnectionPool(db *sql.DB, cfg config.DatabaseConfig) {
    maxOpenConns := cfg.MaxOpenConns
    if maxOpenConns <= 0 {
        maxOpenConns = 25 // Default maximum open connections
    }

    maxIdleConns := cfg.MaxIdleConns
    if maxIdleConns <= 0 {
        maxIdleConns = 5 // Default maximum idle connections
    }

    connMaxLifetime := cfg.ConnMaxLifetime
    if connMaxLifetime <= 0 {
        connMaxLifetime = 300 // Default 5 minutes
    }

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

    log.Printf("Database connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%ds",
        maxOpenConns, maxIdleConns, connMaxLifetime)
}

// testConnection tests the database connection with retry logic
func testConnection(db *sql.DB) error {
    maxRetries := 3
    retryDelay := 2 * time.Second

    for i := 0; i < maxRetries; i++ {
        if err := db.Ping(); err != nil {
            if i == maxRetries-1 {
                return err
            }
            log.Printf("Database connection test failed (attempt %d/%d): %v. Retrying in %v...",
                i+1, maxRetries, err, retryDelay)
            time.Sleep(retryDelay)
            continue
        }
        
        if i > 0 {
            log.Printf("Database connection successful after %d attempts", i+1)
        }
        return nil
    }
    
    return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

// Ping checks if the database connection is alive
func Ping() error {
    if DB == nil {
        return fmt.Errorf("database connection is not initialized")
    }
    return DB.Ping()
}

// Stats returns database connection statistics
func Stats() sql.DBStats {
    if DB == nil {
        return sql.DBStats{}
    }
    return DB.Stats()
}