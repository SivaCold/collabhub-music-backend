package services

import (
	"collabhub-music-backend/internal/models"
	"collabhub-music-backend/internal/repository"
	"github.com/google/uuid"
)

// ProjectService provides project-related business logic
type ProjectServiceInterface struct {
	projectRepo repository.ProjectRepositoryInterface
	userRepo    repository.UserRepositoryInterface
}

// NewProjectService creates a new instance of ProjectService
func NewProjectService(projectRepo repository.ProjectRepositoryInterface, userRepo repository.UserRepositoryInterface) *ProjectServiceInterface {
	return &ProjectServiceInterface{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateProject creates a new project
func (s *ProjectServiceInterface) CreateProject(project *models.Project) error {
	return s.projectRepo.Create(project)
}

// GetProjectByID retrieves a project by ID
func (s *ProjectServiceInterface) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

// GetProjectsByUserID retrieves projects by user ID
func (s *ProjectServiceInterface) GetProjectsByUserID(userID uuid.UUID) ([]*models.Project, error) {
	return s.projectRepo.GetByUserID(userID)
}

// UpdateProject updates a project
func (s *ProjectServiceInterface) UpdateProject(project *models.Project) error {
	return s.projectRepo.Update(project)
}

// DeleteProject deletes a project
func (s *ProjectServiceInterface) DeleteProject(id uuid.UUID) error {
	return s.projectRepo.Delete(id)
}

// AddCollaborator adds a collaborator to a project
func (s *ProjectServiceInterface) AddCollaborator(collaborator *models.ProjectCollaborator) error {
	return s.projectRepo.AddCollaborator(collaborator)
}

// RemoveCollaborator removes a collaborator from a project
func (s *ProjectServiceInterface) RemoveCollaborator(projectID, userID uuid.UUID) error {
	return s.projectRepo.RemoveCollaborator(projectID, userID)
}

// GetCollaborators gets all collaborators for a project
func (s *ProjectServiceInterface) GetCollaborators(projectID uuid.UUID) ([]*models.ProjectCollaborator, error) {
	return s.projectRepo.GetCollaborators(projectID)
}
