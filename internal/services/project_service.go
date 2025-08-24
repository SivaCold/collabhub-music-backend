package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/repository"
)

type ProjectService struct {
    repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *ProjectService {
    return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project) error {
    if project == nil {
        return fmt.Errorf("project data is required")
    }
    
    if project.Name == "" {
        return fmt.Errorf("project name is required")
    }
    
    if project.CreatedBy == uuid.Nil {
        return fmt.Errorf("created_by is required")
    }
    
    // Set default values
    if project.ID == uuid.Nil {
        project.ID = uuid.New()
    }
    
    now := time.Now()
    project.CreatedAt = now
    project.UpdatedAt = now
    
    return s.repo.CreateProject(ctx, project)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
    if id == uuid.Nil {
        return nil, fmt.Errorf("project ID is required")
    }
    
    return s.repo.GetProjectByID(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *models.Project) error {
    if project == nil {
        return fmt.Errorf("project data is required")
    }
    
    if project.ID == uuid.Nil {
        return fmt.Errorf("project ID is required")
    }
    
    if project.Name == "" {
        return fmt.Errorf("project name is required")
    }
    
    // Update timestamp
    project.UpdatedAt = time.Now()
    
    return s.repo.UpdateProject(ctx, project)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
    if id == uuid.Nil {
        return fmt.Errorf("project ID is required")
    }
    
    return s.repo.DeleteProject(ctx, id)
}

func (s *ProjectService) GetProjectsByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.Project, error) {
    if orgID == uuid.Nil {
        return nil, fmt.Errorf("organization ID is required")
    }
    
    return s.repo.GetProjectsByOrganizationID(ctx, orgID)
}

func (s *ProjectService) GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
    if userID == uuid.Nil {
        return nil, fmt.Errorf("user ID is required")
    }
    
    return s.repo.GetProjectsByUserID(ctx, userID)
}

func (s *ProjectService) ListProjects(ctx context.Context, limit, offset int) ([]*models.Project, error) {
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if limit > 100 {
        limit = 100 // Max limit
    }
    
    if offset < 0 {
        offset = 0
    }
    
    return s.repo.ListProjects(ctx, limit, offset)
}

func (s *ProjectService) SearchProjectsByName(ctx context.Context, name string) ([]*models.Project, error) {
    if name == "" {
        return nil, fmt.Errorf("search name is required")
    }
    
    return s.repo.SearchProjectsByName(ctx, name)
}

func (s *ProjectService) ValidateUserAccess(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
    if projectID == uuid.Nil || userID == uuid.Nil {
        return false, fmt.Errorf("project ID and user ID are required")
    }
    
    project, err := s.repo.GetProjectByID(ctx, projectID)
    if err != nil {
        return false, err
    }
    
    // Check if user is the creator
    if project.CreatedBy == userID {
        return true, nil
    }
    
    // Check if user has access through organization membership
    // This would typically involve checking organization membership
    // For now, we'll just check if the project exists
    return project != nil, nil
}

type ProjectSearchFilter struct {
    OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
    UserID         *uuid.UUID `json:"user_id,omitempty"`
    Name           string     `json:"name,omitempty"`
    Status         string     `json:"status,omitempty"`
    Limit          int        `json:"limit"`
    Offset         int        `json:"offset"`
}

func (s *ProjectService) SearchProjects(ctx context.Context, filter *ProjectSearchFilter) ([]*models.Project, error) {
    if filter == nil {
        return s.ListProjects(ctx, 10, 0)
    }
    
    // Set defaults
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }
    if filter.Offset < 0 {
        filter.Offset = 0
    }
    
    // Apply specific filters
    if filter.OrganizationID != nil {
        return s.GetProjectsByOrganizationID(ctx, *filter.OrganizationID)
    }
    
    if filter.UserID != nil {
        return s.GetProjectsByUserID(ctx, *filter.UserID)
    }
    
    if filter.Name != "" {
        return s.SearchProjectsByName(ctx, filter.Name)
    }
    
    return s.ListProjects(ctx, filter.Limit, filter.Offset)
}

type ProjectStats struct {
    ID             uuid.UUID `json:"id"`
    Name           string    `json:"name"`
    TrackCount     int       `json:"track_count"`
    CollaboratorCount int    `json:"collaborator_count"`
    LastActivity   time.Time `json:"last_activity"`
}

func (s *ProjectService) GetProjectStats(ctx context.Context, projectID uuid.UUID) (*ProjectStats, error) {
    if projectID == uuid.Nil {
        return nil, fmt.Errorf("project ID is required")
    }
    
    project, err := s.GetProjectByID(ctx, projectID)
    if err != nil {
        return nil, err
    }
    
    stats := &ProjectStats{
        ID:             project.ID,
        Name:           project.Name,
        TrackCount:     0, // Would need track repository to get actual count
        CollaboratorCount: 0, // Would need collaboration repository to get actual count
        LastActivity:   project.UpdatedAt,
    }
    
    return stats, nil
}