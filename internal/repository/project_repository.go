package repository

import (
	"collabhub-music-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// projectRepository implements the ProjectRepositoryInterface
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new instance of projectRepository
func NewProjectRepository(db *gorm.DB) ProjectRepositoryInterface {
	return &projectRepository{db: db}
}

// Create adds a new project to the database
func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Owner").Preload("Creator").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByUserID retrieves projects by user ID
func (r *projectRepository) GetByUserID(userID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Preload("Owner").Where("owner_id = ? OR created_by = ?", userID, userID).Find(&projects).Error
	return projects, err
}

// Update updates a project in the database
func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete deletes a project from the database
func (r *projectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// AddCollaborator adds a collaborator to a project
func (r *projectRepository) AddCollaborator(projectCollaborator *models.ProjectCollaborator) error {
	return r.db.Create(projectCollaborator).Error
}

// RemoveCollaborator removes a collaborator from a project
func (r *projectRepository) RemoveCollaborator(projectID, userID uuid.UUID) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectCollaborator{}).Error
}

// GetCollaborators gets all collaborators for a project
func (r *projectRepository) GetCollaborators(projectID uuid.UUID) ([]*models.ProjectCollaborator, error) {
	var collaborators []*models.ProjectCollaborator
	err := r.db.Preload("User").Where("project_id = ?", projectID).Find(&collaborators).Error
	return collaborators, err
}
