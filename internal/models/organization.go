package models

import (
    "github.com/google/uuid"
    "time"
)

type Organization struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Name        string    `json:"name" gorm:"not null"`
    Slug        string    `json:"slug" gorm:"unique;not null"`
    Description string    `json:"description"`
    AvatarURL   string    `json:"avatar_url"`
    Website     string    `json:"website"`
    Visibility  string    `json:"visibility" gorm:"default:'public'"`
    CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid"`
    CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}