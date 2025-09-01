// filepath: collabhub-music-backend/internal/models/project.go
package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Project represents a music project
type Project struct {
    ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name           string    `json:"name" gorm:"not null"`
    Description    string    `json:"description"`
    OwnerID        uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
    OrganizationID *uuid.UUID `json:"organization_id,omitempty" gorm:"type:uuid"`
    CreatedBy      uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
    IsPublic       bool      `json:"is_public" gorm:"default:false"`
    CurrentBranch  string    `json:"current_branch" gorm:"default:'main'"`
    Settings       ProjectSettings `json:"settings" gorm:"type:jsonb"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Owner         User                   `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
    Creator       User                   `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
    Organization  *Organization          `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    Collaborators []ProjectCollaborator  `json:"collaborators,omitempty" gorm:"foreignKey:ProjectID"`
    Branches      []Branch               `json:"branches,omitempty" gorm:"foreignKey:ProjectID"`
    Files         []File                 `json:"files,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectSettings holds project-specific settings
type ProjectSettings struct {
    SampleRate    int    `json:"sample_rate"`
    BitDepth      int    `json:"bit_depth"`
    Tempo         int    `json:"tempo"`
    TimeSignature string `json:"time_signature"`
    Key           string `json:"key"`
}

// ProjectCollaborator represents the relationship between users and projects
type ProjectCollaborator struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
    UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
    Role      string    `json:"role" gorm:"default:'collaborator'"` // owner, admin, collaborator, viewer
    InvitedAt time.Time `json:"invited_at"`
    JoinedAt  *time.Time `json:"joined_at"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // Relationships
    Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
    User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook to set ID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}

// BeforeCreate hook for ProjectCollaborator
func (pc *ProjectCollaborator) BeforeCreate(tx *gorm.DB) error {
    if pc.ID == uuid.Nil {
        pc.ID = uuid.New()
    }
    return nil
}