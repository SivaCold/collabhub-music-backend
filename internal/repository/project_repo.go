package repository

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
)

// projectRepo implements the ProjectRepository interface
type projectRepo struct {
    db *sql.DB
}

// NewProjectRepo creates a new instance of projectRepo
func NewProjectRepo(db *sql.DB) ProjectRepository {
    return &projectRepo{db: db}
}

// CreateProject adds a new project to the database
func (r *projectRepo) CreateProject(ctx context.Context, project *models.Project) error {
    query := `
        INSERT INTO projects (id, name, description, organization_id, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    _, err := r.db.ExecContext(ctx, query, project.ID, project.Name, project.Description,
        project.OrganizationID, project.CreatedBy, project.CreatedAt, project.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to create project: %w", err)
    }
    return nil
}

// GetProjectByID retrieves a project by its ID
func (r *projectRepo) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
    query := `
        SELECT id, name, description, organization_id, created_by, created_at, updated_at
        FROM projects
        WHERE id = $1 AND deleted_at IS NULL
    `
    row := r.db.QueryRowContext(ctx, query, id)
    
    var project models.Project
    err := row.Scan(&project.ID, &project.Name, &project.Description, &project.OrganizationID,
        &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("project not found")
        }
        return nil, fmt.Errorf("failed to get project: %w", err)
    }
    return &project, nil
}

// UpdateProject modifies an existing project
func (r *projectRepo) UpdateProject(ctx context.Context, project *models.Project) error {
    query := `
        UPDATE projects
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, project.Name, project.Description,
        project.UpdatedAt, project.ID)
    if err != nil {
        return fmt.Errorf("failed to update project: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("project not found or already deleted")
    }
    
    return nil
}

// DeleteProject removes a project by its ID (soft delete)
func (r *projectRepo) DeleteProject(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE projects
        SET deleted_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete project: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("project not found or already deleted")
    }
    
    return nil
}

// GetProjectsByOrganizationID retrieves projects for a specific organization
func (r *projectRepo) GetProjectsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.Project, error) {
    query := `
        SELECT id, name, description, organization_id, created_by, created_at, updated_at
        FROM projects
        WHERE organization_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, orgID)
    if err != nil {
        return nil, fmt.Errorf("failed to get projects by organization ID: %w", err)
    }
    defer rows.Close()

    var projects []*models.Project
    for rows.Next() {
        var project models.Project
        err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OrganizationID,
            &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan project: %w", err)
        }
        projects = append(projects, &project)
    }

    return projects, nil
}

// GetProjectsByUserID retrieves projects for a specific user
func (r *projectRepo) GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
    query := `
        SELECT DISTINCT p.id, p.name, p.description, p.organization_id, p.created_by, p.created_at, p.updated_at
        FROM projects p
        LEFT JOIN project_members pm ON p.id = pm.project_id
        WHERE (p.created_by = $1 OR pm.user_id = $1) AND p.deleted_at IS NULL
        AND (pm.deleted_at IS NULL OR pm.deleted_at IS NULL)
        ORDER BY p.created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get projects by user ID: %w", err)
    }
    defer rows.Close()

    var projects []*models.Project
    for rows.Next() {
        var project models.Project
        err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OrganizationID,
            &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan project: %w", err)
        }
        projects = append(projects, &project)
    }

    return projects, nil
}

// ListProjects retrieves projects with pagination
func (r *projectRepo) ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, error) {
    query := `
        SELECT id, name, description, organization_id, created_by, created_at, updated_at
        FROM projects
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list projects: %w", err)
    }
    defer rows.Close()

    var projects []*models.Project
    for rows.Next() {
        var project models.Project
        err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OrganizationID,
            &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan project: %w", err)
        }
        projects = append(projects, &project)
    }

    return projects, nil
}

// SearchProjectsByName searches projects by name
func (r *projectRepo) SearchProjectsByName(ctx context.Context, name string) ([]*models.Project, error) {
    query := `
        SELECT id, name, description, organization_id, created_by, created_at, updated_at
        FROM projects
        WHERE LOWER(name) LIKE LOWER($1) AND deleted_at IS NULL
        ORDER BY created_at DESC
    `
    searchPattern := "%" + name + "%"
    rows, err := r.db.QueryContext(ctx, query, searchPattern)
    if err != nil {
        return nil, fmt.Errorf("failed to search projects by name: %w", err)
    }
    defer rows.Close()

    var projects []*models.Project
    for rows.Next() {
        var project models.Project
        err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.OrganizationID,
            &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan project: %w", err)
        }
        projects = append(projects, &project)
    }

    return projects, nil
}