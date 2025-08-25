package database

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
)

// RunMigrations executes the complete SQL schema
func RunMigrations(db *sql.DB) error {
    // Check if migrations have already been run
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'").Scan(&count)
    if err != nil {
        return fmt.Errorf("failed to check if migrations exist: %w", err)
    }

    if count > 0 {
        log.Println("Database schema already exists, skipping migrations")
        return nil
    }

    // Read and execute the init.sql file
    sqlFile := "internal/database/migrations/init.sql"
    sqlContent, err := ioutil.ReadFile(sqlFile)
    if err != nil {
        return fmt.Errorf("failed to read init.sql: %w", err)
    }

    // Execute the complete SQL script
    _, err = db.Exec(string(sqlContent))
    if err != nil {
        return fmt.Errorf("failed to execute init.sql: %w", err)
    }

    log.Println("Database schema created successfully from init.sql")
    return nil
}