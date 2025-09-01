package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Branch represents a project branch for version control
type Branch struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ProjectID   uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    ParentBranch string   `json:"parent_branch"`
    IsDefault   bool      `json:"is_default" gorm:"default:false"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Project   Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
    Creator   User    `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
    Files     []File  `json:"files,omitempty" gorm:"foreignKey:BranchID"`
}

// BeforeCreate hook to set ID
func (b *Branch) BeforeCreate(tx *gorm.DB) error {
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    return nil
}