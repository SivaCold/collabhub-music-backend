package services

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/repository"
)

type UserService struct {
    userRepo        repository.UserRepository
    keycloakService *KeycloakService
}

func NewUserService(userRepo repository.UserRepository, keycloakService *KeycloakService) *UserService {
    return &UserService{
        userRepo:        userRepo,
        keycloakService: keycloakService,
    }
}

// CreateUser creates a user in both Keycloak and local database
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
    if user == nil {
        return fmt.Errorf("user data is required")
    }

    if user.Email == "" || user.Username == "" {
        return fmt.Errorf("email and username are required")
    }

    // Create user in Keycloak first
    keycloakUser := &KeycloakUser{
        Username:  user.Username,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Enabled:   true,
    }

    keycloakID, err := s.keycloakService.CreateUser(ctx, keycloakUser)
    if err != nil {
        return fmt.Errorf("failed to create user in Keycloak: %w", err)
    }

    // Set user properties
    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    user.KeycloakID = keycloakID
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now

    // Create user in local database
    if err := s.userRepo.CreateUser(ctx, user); err != nil {
        // Rollback: delete from Keycloak if local creation fails
        _ = s.keycloakService.DeleteUser(ctx, keycloakID)
        return fmt.Errorf("failed to create user in database: %w", err)
    }

    return nil
}

// GetUserByID retrieves a user by ID from local database
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    if id == uuid.Nil {
        return nil, fmt.Errorf("user ID is required")
    }

    return s.userRepo.GetUserByID(ctx, id)
}

// GetUserByKeycloakID retrieves a user by their Keycloak ID
func (s *UserService) GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error) {
    if keycloakID == "" {
        return nil, fmt.Errorf("Keycloak ID is required")
    }

    return s.userRepo.GetUserByKeycloakID(ctx, keycloakID)
}

// GetUserByEmail retrieves a user by email from local database
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    if email == "" {
        return nil, fmt.Errorf("email is required")
    }

    return s.userRepo.GetUserByEmail(ctx, email)
}

// UpdateUser updates user in both Keycloak and local database
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
    if user == nil {
        return fmt.Errorf("user data is required")
    }

    if user.ID == uuid.Nil {
        return fmt.Errorf("user ID is required")
    }

    // Get current user data
    currentUser, err := s.userRepo.GetUserByID(ctx, user.ID)
    if err != nil {
        return fmt.Errorf("failed to get current user: %w", err)
    }

    // Update in Keycloak
    keycloakUser := &KeycloakUser{
        ID:        currentUser.KeycloakID,
        Username:  user.Username,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Enabled:   true,
    }

    if err := s.keycloakService.UpdateUser(ctx, currentUser.KeycloakID, keycloakUser); err != nil {
        return fmt.Errorf("failed to update user in Keycloak: %w", err)
    }

    // Update timestamp
    user.UpdatedAt = time.Now()

    // Update in local database
    if err := s.userRepo.UpdateUser(ctx, user); err != nil {
        return fmt.Errorf("failed to update user in database: %w", err)
    }

    return nil
}

// DeleteUser deletes user from both Keycloak and local database
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
    if id == uuid.Nil {
        return fmt.Errorf("user ID is required")
    }

    // Get user to get Keycloak ID
    user, err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }

    // Delete from Keycloak
    if err := s.keycloakService.DeleteUser(ctx, user.KeycloakID); err != nil {
        return fmt.Errorf("failed to delete user from Keycloak: %w", err)
    }

    // Delete from local database (soft delete)
    if err := s.userRepo.DeleteUser(ctx, id); err != nil {
        return fmt.Errorf("failed to delete user from database: %w", err)
    }

    return nil
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if limit > 100 {
        limit = 100 // Max limit
    }
    
    if offset < 0 {
        offset = 0
    }

    return s.userRepo.ListUsers(ctx, limit, offset)
}

// GetUsersByOrganizationID retrieves users belonging to an organization
func (s *UserService) GetUsersByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]*models.User, error) {
    if orgID == uuid.Nil {
        return nil, fmt.Errorf("organization ID is required")
    }

    return s.userRepo.GetUsersByOrganizationID(ctx, orgID)
}

// SyncUserFromKeycloak synchronizes user data from Keycloak token
func (s *UserService) SyncUserFromKeycloak(ctx context.Context, token string) (*models.User, error) {
    if token == "" {
        return nil, fmt.Errorf("token is required")
    }

    // Get user info from Keycloak
    keycloakUser, err := s.keycloakService.GetUserInfo(ctx, token)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info from Keycloak: %w", err)
    }

    // Try to find existing user by Keycloak ID
    existingUser, err := s.userRepo.GetUserByKeycloakID(ctx, keycloakUser.ID)
    if err == nil {
        // Update existing user with latest Keycloak data
        existingUser.Username = keycloakUser.Username
        existingUser.Email = keycloakUser.Email
        existingUser.FirstName = keycloakUser.FirstName
        existingUser.LastName = keycloakUser.LastName
        existingUser.UpdatedAt = time.Now()

        if err := s.userRepo.UpdateUser(ctx, existingUser); err != nil {
            return nil, fmt.Errorf("failed to update existing user: %w", err)
        }

        return existingUser, nil
    }

    // Create new user if not found
    newUser := &models.User{
        ID:          uuid.New(),
        KeycloakID:  keycloakUser.ID,
        Username:    keycloakUser.Username,
        Email:       keycloakUser.Email,
        FirstName:   keycloakUser.FirstName,
        LastName:    keycloakUser.LastName,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
        return nil, fmt.Errorf("failed to create new user: %w", err)
    }

    return newUser, nil
}

// ValidateUserExists checks if a user exists and is active
func (s *UserService) ValidateUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
    if userID == uuid.Nil {
        return false, fmt.Errorf("user ID is required")
    }

    user, err := s.userRepo.GetUserByID(ctx, userID)
    if err != nil {
        return false, nil // User doesn't exist
    }

    return user != nil, nil
}