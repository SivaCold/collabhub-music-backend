package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User represents a user in the system
type User struct {
    ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    KeycloakID string    `json:"keycloak_id" gorm:"uniqueIndex"`
    Email      string    `json:"email" gorm:"uniqueIndex;not null"`
    Username   string    `json:"username" gorm:"uniqueIndex;not null"`
    FirstName  string    `json:"first_name"`
    LastName   string    `json:"last_name"`
    Avatar     string    `json:"avatar"`
    Password   string    `json:"-" gorm:"not null"` // Never include in JSON
    IsActive   bool      `json:"is_active" gorm:"default:true"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    OwnedProjects []Project `json:"owned_projects,omitempty" gorm:"foreignKey:OwnerID"`
    Collaborations []ProjectCollaborator `json:"collaborations,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook to set ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
    if u.FirstName != "" && u.LastName != "" {
        return u.FirstName + " " + u.LastName
    }
    if u.FirstName != "" {
        return u.FirstName
    }
    if u.LastName != "" {
        return u.LastName
    }
    return u.Username
}