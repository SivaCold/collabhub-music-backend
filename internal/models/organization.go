package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	Website     string    `json:"website"`
	Visibility  string    `json:"visibility" gorm:"default:'public'"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Creator  User                 `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Members  []OrganizationMember `json:"members,omitempty" gorm:"foreignKey:OrganizationID"`
	Projects []Project            `json:"projects,omitempty" gorm:"foreignKey:OrganizationID"`
}

// OrganizationMember represents the relationship between users and organizations
type OrganizationMember struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Role           string     `json:"role" gorm:"default:'member'"` // owner, admin, member
	InvitedAt      time.Time  `json:"invited_at"`
	JoinedAt       *time.Time `json:"joined_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	User         User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
