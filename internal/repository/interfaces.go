package repository

import (
	"collabhub-music-backend/internal/models"

	"github.com/google/uuid"
)

// Type aliases for backward compatibility
type UserRepository = UserRepositoryInterface
type ProjectRepository = ProjectRepositoryInterface
type OrganizationRepository = OrganizationRepositoryInterface

// UserRepositoryInterface defines methods for user repository
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

// ProjectRepositoryInterface defines methods for project repository
type ProjectRepositoryInterface interface {
	Create(project *models.Project) error
	GetByID(id uuid.UUID) (*models.Project, error)
	GetByUserID(userID uuid.UUID) ([]*models.Project, error)
	Update(project *models.Project) error
	Delete(id uuid.UUID) error
	AddCollaborator(projectCollaborator *models.ProjectCollaborator) error
	RemoveCollaborator(projectID, userID uuid.UUID) error
	GetCollaborators(projectID uuid.UUID) ([]*models.ProjectCollaborator, error)
}

// OrganizationRepositoryInterface defines methods for organization repository
type OrganizationRepositoryInterface interface {
	Create(organization *models.Organization) error
	GetByID(id uuid.UUID) (*models.Organization, error)
	GetByUserID(userID uuid.UUID) ([]*models.Organization, error)
	Update(organization *models.Organization) error
	Delete(id uuid.UUID) error
	AddMember(member *models.OrganizationMember) error
	RemoveMember(organizationID, userID uuid.UUID) error
	GetMembers(organizationID uuid.UUID) ([]*models.OrganizationMember, error)
}

// FileRepositoryInterface defines methods for file repository
type FileRepositoryInterface interface {
	Create(file *models.File) error
	GetByID(id uuid.UUID) (*models.File, error)
	GetByProjectID(projectID uuid.UUID) ([]*models.File, error)
	GetByBranchID(branchID uuid.UUID) ([]*models.File, error)
	Update(file *models.File) error
	Delete(id uuid.UUID) error
	CreateVersion(version *models.FileVersion) error
	GetVersions(fileID uuid.UUID) ([]*models.FileVersion, error)
	CreateAudioMetadata(metadata *models.AudioMetadata) error
	UpdateAudioMetadata(metadata *models.AudioMetadata) error
}

// BranchRepositoryInterface defines methods for branch repository
type BranchRepositoryInterface interface {
	Create(branch *models.Branch) error
	GetByID(id uuid.UUID) (*models.Branch, error)
	GetByProjectID(projectID uuid.UUID) ([]*models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uuid.UUID) error
	SetDefault(branchID uuid.UUID) error
}
