package repository

import (
	"collabhub-music-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// organizationRepository implements the OrganizationRepositoryInterface
type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new instance of organizationRepository
func NewOrganizationRepository(db *gorm.DB) OrganizationRepositoryInterface {
	return &organizationRepository{db: db}
}

// Create adds a new organization to the database
func (r *organizationRepository) Create(organization *models.Organization) error {
	return r.db.Create(organization).Error
}

// GetByID retrieves an organization by ID
func (r *organizationRepository) GetByID(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Preload("Creator").First(&organization, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

// GetByUserID retrieves organizations by user ID
func (r *organizationRepository) GetByUserID(userID uuid.UUID) ([]*models.Organization, error) {
	var organizations []*models.Organization
	err := r.db.Joins("JOIN organization_members ON organization_members.organization_id = organizations.id").
		Where("organization_members.user_id = ?", userID).
		Find(&organizations).Error
	return organizations, err
}

// Update updates an organization in the database
func (r *organizationRepository) Update(organization *models.Organization) error {
	return r.db.Save(organization).Error
}

// Delete deletes an organization from the database
func (r *organizationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Organization{}, id).Error
}

// AddMember adds a member to an organization
func (r *organizationRepository) AddMember(member *models.OrganizationMember) error {
	return r.db.Create(member).Error
}

// RemoveMember removes a member from an organization
func (r *organizationRepository) RemoveMember(organizationID, userID uuid.UUID) error {
	return r.db.Where("organization_id = ? AND user_id = ?", organizationID, userID).Delete(&models.OrganizationMember{}).Error
}

// GetMembers gets all members for an organization
func (r *organizationRepository) GetMembers(organizationID uuid.UUID) ([]*models.OrganizationMember, error) {
	var members []*models.OrganizationMember
	err := r.db.Preload("User").Where("organization_id = ?", organizationID).Find(&members).Error
	return members, err
}
