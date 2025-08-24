package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/repository"
)

type OrganizationService struct {
    orgRepo  repository.OrganizationRepository
    userRepo repository.UserRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, userRepo repository.UserRepository) *OrganizationService {
    return &OrganizationService{
        orgRepo:  orgRepo,
        userRepo: userRepo,
    }
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, org *models.Organization) error {
    if org == nil {
        return fmt.Errorf("organization data is required")
    }
    
    if org.Name == "" {
        return fmt.Errorf("organization name is required")
    }
    
    if org.CreatedBy == uuid.Nil {
        return fmt.Errorf("created_by is required")
    }
    
    // Verify that the creator exists
    _, err := s.userRepo.GetUserByID(ctx, org.CreatedBy)
    if err != nil {
        return fmt.Errorf("creator user not found: %w", err)
    }
    
    // Set default values
    if org.ID == uuid.Nil {
        org.ID = uuid.New()
    }
    
    now := time.Now()
    org.CreatedAt = now
    org.UpdatedAt = now
    
    return s.orgRepo.CreateOrganization(ctx, org)
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
    if id == uuid.Nil {
        return nil, fmt.Errorf("organization ID is required")
    }
    
    return s.orgRepo.GetOrganizationByID(ctx, id)
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, org *models.Organization) error {
    if org == nil {
        return fmt.Errorf("organization data is required")
    }
    
    if org.ID == uuid.Nil {
        return fmt.Errorf("organization ID is required")
    }
    
    if org.Name == "" {
        return fmt.Errorf("organization name is required")
    }
    
    // Update timestamp
    org.UpdatedAt = time.Now()
    
    return s.orgRepo.UpdateOrganization(ctx, org)
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
    if id == uuid.Nil {
        return fmt.Errorf("organization ID is required")
    }
    
    return s.orgRepo.DeleteOrganization(ctx, id)
}

func (s *OrganizationService) GetOrganizationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Organization, error) {
    if userID == uuid.Nil {
        return nil, fmt.Errorf("user ID is required")
    }
    
    return s.orgRepo.GetOrganizationsByUserID(ctx, userID)
}

func (s *OrganizationService) ListOrganizations(ctx context.Context, limit, offset int) ([]*models.Organization, error) {
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if limit > 100 {
        limit = 100 // Max limit
    }
    
    if offset < 0 {
        offset = 0
    }
    
    return s.orgRepo.ListOrganizations(ctx, limit, offset)
}

func (s *OrganizationService) GetOrganizationByName(ctx context.Context, name string) (*models.Organization, error) {
    if name == "" {
        return nil, fmt.Errorf("organization name is required")
    }
    
    return s.orgRepo.GetOrganizationByName(ctx, name)
}

func (s *OrganizationService) AddUserToOrganization(ctx context.Context, orgID, userID uuid.UUID) error {
    if orgID == uuid.Nil || userID == uuid.Nil {
        return fmt.Errorf("organization ID and user ID are required")
    }
    
    // Verify organization exists
    _, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
    if err != nil {
        return fmt.Errorf("organization not found: %w", err)
    }
    
    // Verify user exists
    _, err = s.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }
    
    return s.orgRepo.AddUserToOrganization(ctx, orgID, userID)
}

func (s *OrganizationService) RemoveUserFromOrganization(ctx context.Context, orgID, userID uuid.UUID) error {
    if orgID == uuid.Nil || userID == uuid.Nil {
        return fmt.Errorf("organization ID and user ID are required")
    }
    
    return s.orgRepo.RemoveUserFromOrganization(ctx, orgID, userID)
}

func (s *OrganizationService) ValidateUserAccess(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
    if orgID == uuid.Nil || userID == uuid.Nil {
        return false, fmt.Errorf("organization ID and user ID are required")
    }
    
    org, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
    if err != nil {
        return false, err
    }
    
    // Check if user is the creator
    if org.CreatedBy == userID {
        return true, nil
    }
    
    // Check if user is a member of the organization
    userOrgs, err := s.orgRepo.GetOrganizationsByUserID(ctx, userID)
    if err != nil {
        return false, err
    }
    
    for _, userOrg := range userOrgs {
        if userOrg.ID == orgID {
            return true, nil
        }
    }
    
    return false, nil
}

type OrganizationSearchFilter struct {
    UserID *uuid.UUID `json:"user_id,omitempty"`
    Name   string     `json:"name,omitempty"`
    Limit  int        `json:"limit"`
    Offset int        `json:"offset"`
}

func (s *OrganizationService) SearchOrganizations(ctx context.Context, filter *OrganizationSearchFilter) ([]*models.Organization, error) {
    if filter == nil {
        return s.ListOrganizations(ctx, 10, 0)
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
    if filter.UserID != nil {
        return s.GetOrganizationsByUserID(ctx, *filter.UserID)
    }
    
    if filter.Name != "" {
        org, err := s.GetOrganizationByName(ctx, filter.Name)
        if err != nil {
            return []*models.Organization{}, nil
        }
        return []*models.Organization{org}, nil
    }
    
    return s.ListOrganizations(ctx, filter.Limit, filter.Offset)
}

type OrganizationStats struct {
    ID           uuid.UUID `json:"id"`
    Name         string    `json:"name"`
    MemberCount  int       `json:"member_count"`
    ProjectCount int       `json:"project_count"`
    CreatedAt    time.Time `json:"created_at"`
}

func (s *OrganizationService) GetOrganizationStats(ctx context.Context, orgID uuid.UUID) (*OrganizationStats, error) {
    if orgID == uuid.Nil {
        return nil, fmt.Errorf("organization ID is required")
    }
    
    org, err := s.GetOrganizationByID(ctx, orgID)
    if err != nil {
        return nil, err
    }
    
    stats := &OrganizationStats{
        ID:           org.ID,
        Name:         org.Name,
        MemberCount:  0, // Would need to query user-organization relationships
        ProjectCount: 0, // Would need project repository to get actual count
        CreatedAt:    org.CreatedAt,
    }
    
    return stats, nil
}

func (s *OrganizationService) GetOrganizationMembers(ctx context.Context, orgID uuid.UUID) ([]*models.User, error) {
    if orgID == uuid.Nil {
        return nil, fmt.Errorf("organization ID is required")
    }
    
    // Verify organization exists
    _, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
    if err != nil {
        return nil, fmt.Errorf("organization not found: %w", err)
    }
    
    return s.userRepo.GetUsersByOrganizationID(ctx, orgID)
}

func (s *OrganizationService) IsUserMember(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
    if orgID == uuid.Nil || userID == uuid.Nil {
        return false, fmt.Errorf("organization ID and user ID are required")
    }
    
    userOrgs, err := s.orgRepo.GetOrganizationsByUserID(ctx, userID)
    if err != nil {
        return false, err
    }
    
    for _, org := range userOrgs {
        if org.ID == orgID {
            return true, nil
        }
    }
    
    return false, nil
}