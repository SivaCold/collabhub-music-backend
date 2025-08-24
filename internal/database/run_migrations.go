package database

import (
    "database/sql"
    "fmt"
    "log"
)

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
    if err := createMigrationsTable(db); err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }

    migrations := []Migration{
        {
            ID:   1,
            Name: "create_users_table",
            SQL:  createUsersTableSQL,
        },
        {
            ID:   2,
            Name: "create_organizations_table",
            SQL:  createOrganizationsTableSQL,
        },
        {
            ID:   3,
            Name: "create_projects_table",
            SQL:  createProjectsTableSQL,
        },
        {
            ID:   4,
            Name: "create_user_organizations_table",
            SQL:  createUserOrganizationsTableSQL,
        },
        {
            ID:   5,
            Name: "add_indexes",
            SQL:  addIndexesSQL,
        },
    }

    for _, migration := range migrations {
        if err := runMigration(db, migration); err != nil {
            return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
        }
    }

    log.Println("All migrations completed successfully")
    return nil
}

type Migration struct {
    ID   int
    Name string
    SQL  string
}

func createMigrationsTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS migrations (
        id INTEGER PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
    
    _, err := db.Exec(query)
    return err
}

func runMigration(db *sql.DB, migration Migration) error {
    // Check if migration already executed
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE id = $1", migration.ID).Scan(&count)
    if err != nil {
        return err
    }

    if count > 0 {
        log.Printf("Migration %s already executed, skipping", migration.Name)
        return nil
    }

    // Execute migration
    _, err = db.Exec(migration.SQL)
    if err != nil {
        return err
    }

    // Record migration
    _, err = db.Exec("INSERT INTO migrations (id, name) VALUES ($1, $2)", migration.ID, migration.Name)
    if err != nil {
        return err
    }

    log.Printf("Migration %s executed successfully", migration.Name)
    return nil
}

const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keycloak_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);`

const createOrganizationsTableSQL = `
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);`

const createProjectsTableSQL = `
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    organization_id UUID REFERENCES organizations(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);`

const createUserOrganizationsTableSQL = `
CREATE TABLE IF NOT EXISTS user_organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, organization_id)
);`

const addIndexesSQL = `
CREATE INDEX IF NOT EXISTS idx_users_keycloak_id ON users(keycloak_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_projects_organization_id ON projects(organization_id);
CREATE INDEX IF NOT EXISTS idx_organizations_created_by ON organizations(created_by);
CREATE INDEX IF NOT EXISTS idx_user_organizations_user_id ON user_organizations(user_id);
CREATE INDEX IF NOT EXISTS idx_user_organizations_organization_id ON user_organizations(organization_id);
`