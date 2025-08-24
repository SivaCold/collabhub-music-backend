package repository

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
)

// userRepo implements the UserRepository interface
type userRepo struct {
    db *sql.DB
}

// NewUserRepo creates a new instance of userRepo
func NewUserRepo(db *sql.DB) UserRepository {
    return &userRepo{db: db}
}

// CreateUser adds a new user to the database
func (r *userRepo) CreateUser(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (id, keycloak_id, username, email, first_name, last_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
    _, err := r.db.ExecContext(ctx, query, user.ID, user.KeycloakID, user.Username, 
        user.Email, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

// GetUserByID retrieves a user by their ID
func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    query := `
        SELECT id, keycloak_id, username, email, first_name, last_name, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, id)
    
    var user models.User
    err := row.Scan(&user.ID, &user.KeycloakID, &user.Username, &user.Email,
        &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}

// GetUserByKeycloakID retrieves a user by their Keycloak ID
func (r *userRepo) GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error) {
    query := `
        SELECT id, keycloak_id, username, email, first_name, last_name, created_at, updated_at
        FROM users
        WHERE keycloak_id = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, keycloakID)
    
    var user models.User
    err := row.Scan(&user.ID, &user.KeycloakID, &user.Username, &user.Email,
        &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user by keycloak ID: %w", err)
    }
    return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, keycloak_id, username, email, first_name, last_name, created_at, updated_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, email)
    
    var user models.User
    err := row.Scan(&user.ID, &user.KeycloakID, &user.Username, &user.Email,
        &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    return &user, nil
}

// UpdateUser modifies an existing user
func (r *userRepo) UpdateUser(ctx context.Context, user *models.User) error {
    query := `
        UPDATE users
        SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5
        WHERE id = $6 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, user.Username, user.Email,
        user.FirstName, user.LastName, user.UpdatedAt, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found or already deleted")
    }
    
    return nil
}

// DeleteUser removes a user by their ID (soft delete)
func (r *userRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE users
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found or already deleted")
    }
    
    return nil
}

// ListUsers retrieves users with pagination
func (r *userRepo) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
    query := `
        SELECT id, keycloak_id, username, email, first_name, last_name, created_at, updated_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.KeycloakID, &user.Username, &user.Email,
            &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }

    return users, nil
}

// GetUsersByOrganizationID retrieves users for a specific organization
func (r *userRepo) GetUsersByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.User, error) {
    query := `
        SELECT u.id, u.keycloak_id, u.username, u.email, u.first_name, u.last_name, u.created_at, u.updated_at
        FROM users u
        JOIN organization_members om ON u.id = om.user_id
        WHERE om.organization_id = $1 AND u.deleted_at IS NULL AND om.deleted_at IS NULL
        ORDER BY u.created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, orgID)
    if err != nil {
        return nil, fmt.Errorf("failed to get users by organization ID: %w", err)
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.KeycloakID, &user.Username, &user.Email,
            &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }

    return users, nil
}