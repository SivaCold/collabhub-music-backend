package models

import (
    "time"
    "github.com/google/uuid"
)

type Review struct {
    ID        uuid.UUID  `json:"id" db:"id"`
    ProjectID *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
    TrackID   *uuid.UUID `json:"track_id,omitempty" db:"track_id"`
    UserID    uuid.UUID  `json:"user_id" db:"user_id"`
    Title     string     `json:"title" db:"title"`
    Content   string     `json:"content" db:"content"`
    Rating    *int       `json:"rating,omitempty" db:"rating"` // 1-5 stars
    Status    string     `json:"status" db:"status"` // draft, published
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}