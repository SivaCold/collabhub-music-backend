package models

import (
    "github.com/google/uuid"
    "time"
)

type Project struct {
    ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Name            string    `json:"name" gorm:"not null"`
    Slug            string    `json:"slug" gorm:"not null;unique"`
    Description     string    `json:"description"`
    OrganizationID  uuid.UUID `json:"organization_id" gorm:"type:uuid"`
    CreatedBy       uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
    Visibility      string    `json:"visibility" gorm:"default:'public'"`
    Genre           string    `json:"genre"`
    Tempo           int       `json:"tempo"` // BPM
    KeySignature    string    `json:"key_signature"`
    TimeSignature    string    `json:"time_signature"`
    Status          string    `json:"status" gorm:"default:'active'"`
    License         string    `json:"license"`
    StarsCount      int       `json:"stars_count" gorm:"default:0"`
    ForksCount      int       `json:"forks_count" gorm:"default:0"`
    CollaboratorsCount int    `json:"collaborators_count" gorm:"default:0"`
    CreatedAt       time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt       time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}