package services

import (
	"collabhub-music-backend/internal/models"
	"collabhub-music-backend/internal/repository"
	"github.com/google/uuid"
)

// UserService provides user-related business logic
type UserServiceInterface struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepositoryInterface) *UserServiceInterface {
	return &UserServiceInterface{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserServiceInterface) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}

// GetUserByID retrieves a user by ID
func (s *UserServiceInterface) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserServiceInterface) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// GetUserByUsername retrieves a user by username
func (s *UserServiceInterface) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

// UpdateUser updates a user
func (s *UserServiceInterface) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user
func (s *UserServiceInterface) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
