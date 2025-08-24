package repository

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
)

// organizationRepo implements the OrganizationRepository interface
type organizationRepo struct {
    db *sql.DB
}

// NewOrganizationRepo creates a new instance of organizationRepo
func NewOrganizationRepo(db *sql.DB) OrganizationRepository {
    return &organizationRepo{db: db}
}

// CreateOrganization adds a new organization to the database
func (r *organizationRepo) CreateOrganization(ctx context.Context, org *models.Organization) error {
    query := `
        INSERT INTO organizations (id, name, description, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query, org.ID, org.Name, org.Description, 
        org.CreatedBy, org.CreatedAt, org.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to create organization: %w", err)
    }
    return nil
}

// GetOrganizationByID retrieves an organization by its ID
func (r *organizationRepo) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
    query := `
        SELECT id, name, description, created_by, created_at, updated_at
        FROM organizations
        WHERE id = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, id)
    
    var org models.Organization
    err := row.Scan(&org.ID, &org.Name, &org.Description, &org.CreatedBy, 
        &org.CreatedAt, &org.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("organization not found")
        }
        return nil, fmt.Errorf("failed to get organization: %w", err)
    }
    return &org, nil
}

// UpdateOrganization modifies an existing organization
func (r *organizationRepo) UpdateOrganization(ctx context.Context, org *models.Organization) error {
    query := `
        UPDATE organizations
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, org.Name, org.Description, 
        org.UpdatedAt, org.ID)
    if err != nil {
        return fmt.Errorf("failed to update organization: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("organization not found or already deleted")
    }
    
    return nil
}

// DeleteOrganization removes an organization by its ID (soft delete)
func (r *organizationRepo) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE organizations
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete organization: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("organization not found or already deleted")
    }
    
    return nil
}

// GetOrganizationsByUserID retrieves organizations for a specific user
func (r *organizationRepo) GetOrganizationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
    query := `
        SELECT o.id, o.name, o.description, o.created_by, o.created_at, o.updated_at
        FROM organizations o
        JOIN organization_members om ON o.id = om.organization_id
        WHERE om.user_id = $1 AND o.deleted_at IS NULL AND om.deleted_at IS NULL
        ORDER BY o.created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get organizations by user ID: %w", err)
    }
    defer rows.Close()

    var organizations []*models.Organization
    for rows.Next() {
        var org models.Organization
        err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.CreatedBy,
            &org.CreatedAt, &org.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan organization: %w", err)
        }
        organizations = append(organizations, &org)
    }

    return organizations, nil
}

// ListOrganizations retrieves organizations with pagination
func (r *organizationRepo) ListOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error) {
    query := `
        SELECT id, name, description, created_by, created_at, updated_at
        FROM organizations
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list organizations: %w", err)
    }
    defer rows.Close()

    var organizations []*models.Organization
    for rows.Next() {
        var org models.Organization
        err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.CreatedBy,
            &org.CreatedAt, &org.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan organization: %w", err)
        }
        organizations = append(organizations, &org)
    }

    return organizations, nil
}

// GetOrganizationByName retrieves an organization by its name
func (r *organizationRepo) GetOrganizationByName(ctx context.Context, name string) (*models.Organization, error) {
    query := `
        SELECT id, name, description, created_by, created_at, updated_at
        FROM organizations
        WHERE name = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, name)
    
    var org models.Organization
    err := row.Scan(&org.ID, &org.Name, &org.Description, &org.CreatedBy,
        &org.CreatedAt, &org.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("organization not found")
        }
        return nil, fmt.Errorf("failed to get organization by name: %w", err)
    }
    return &org, nil
}

// AddUserToOrganization adds a user to an organization
func (r *organizationRepo) AddUserToOrganization(ctx context.Context, orgID, userID uuid.UUID) error {
    query := `
        INSERT INTO organization_members (id, organization_id, user_id, created_at)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (organization_id, user_id) DO NOTHING
    `
    memberID := uuid.New()
    _, err := r.db.ExecContext(ctx, query, memberID, orgID, userID)
    if err != nil {
        return fmt.Errorf("failed to add user to organization: %w", err)
    }
    return nil
}

// RemoveUserFromOrganization removes a user from an organization
func (r *organizationRepo) RemoveUserFromOrganization(ctx context.Context, orgID, userID uuid.UUID) error {
    query := `
        UPDATE organization_members
        SET deleted_at = NOW()
        WHERE organization_id = $1 AND user_id = $2 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, orgID, userID)
    if err != nil {
        return fmt.Errorf("failed to remove user from organization: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found in organization or already removed")
    }
    
    return nil
}