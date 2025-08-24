package models

import (
    "github.com/google/uuid"
    "time"
)

type User struct {
    ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
    KeycloakUserID  string      `json:"keycloak_user_id" gorm:"unique;not null"`
    Username        string      `json:"username" gorm:"unique;not null"`
    FirstName       string      `json:"first_name"`
    LastName        string      `json:"last_name"`
    KeycloakID      string      `json:"keycloak_id" gorm:"unique;not null"`
    Email           string      `json:"email" gorm:"unique;not null"`
    DisplayName     string      `json:"display_name"`
    Bio             string      `json:"bio"`
    AvatarURL       string      `json:"avatar_url"`
    Location        string      `json:"location"`
    Website         string      `json:"website"`
    MusicalGenres   []string    `json:"musical_genres" gorm:"type:text[]"`
    Instruments     []string    `json:"instruments" gorm:"type:text[]"`
    ExperienceLevel string      `json:"experience_level" gorm:"type:varchar(20)"`
    IsVerified      bool        `json:"is_verified" gorm:"default:false"`
    IsActive        bool        `json:"is_active" gorm:"default:true"`
    CreatedAt       time.Time   `json:"created_at" gorm:"default:current_timestamp"`
    UpdatedAt       time.Time   `json:"updated_at" gorm:"default:current_timestamp"`
}