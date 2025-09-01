package services

import (
	"fmt"

	"collabhub-music-backend/internal/models"
	"collabhub-music-backend/internal/repository"
	"github.com/google/uuid"
)

// OrganizationService provides organization-related business logic
type OrganizationServiceInterface struct {
	orgRepo  repository.OrganizationRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

// NewOrganizationService creates a new instance of OrganizationService
func NewOrganizationService(orgRepo repository.OrganizationRepositoryInterface, userRepo repository.UserRepositoryInterface) *OrganizationServiceInterface {
	return &OrganizationServiceInterface{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// CreateOrganization creates a new organization
func (s *OrganizationServiceInterface) CreateOrganization(org *models.Organization) error {
	return s.orgRepo.Create(org)
}

// GetOrganizationByID retrieves an organization by ID
func (s *OrganizationServiceInterface) GetOrganizationByID(id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.GetByID(id)
}

// UpdateOrganization updates an organization
func (s *OrganizationServiceInterface) UpdateOrganization(org *models.Organization) error {
	return s.orgRepo.Update(org)
}

// DeleteOrganization deletes an organization
func (s *OrganizationServiceInterface) DeleteOrganization(id uuid.UUID) error {
	return s.orgRepo.Delete(id)
}

// GetOrganizationsByUserID gets organizations for a user
func (s *OrganizationServiceInterface) GetOrganizationsByUserID(userID uuid.UUID) ([]*models.Organization, error) {
	return s.orgRepo.GetByUserID(userID)
}

// AddMember adds a member to an organization
func (s *OrganizationServiceInterface) AddMember(member *models.OrganizationMember) error {
	return s.orgRepo.AddMember(member)
}

// RemoveMember removes a member from an organization
func (s *OrganizationServiceInterface) RemoveMember(organizationID, userID uuid.UUID) error {
	return s.orgRepo.RemoveMember(organizationID, userID)
}

// GetMembers gets all members of an organization
func (s *OrganizationServiceInterface) GetMembers(organizationID uuid.UUID) ([]*models.OrganizationMember, error) {
	return s.orgRepo.GetMembers(organizationID)
}

// ListOrganizations lists all organizations (placeholder)
func (s *OrganizationServiceInterface) ListOrganizations() ([]*models.Organization, error) {
	return nil, fmt.Errorf("not implemented")
}
