package models

import (
    "time"
    "github.com/google/uuid"
)

type Lyrics struct {
    ID        uuid.UUID  `json:"id" db:"id"`
    TrackID   uuid.UUID  `json:"track_id" db:"track_id"`
    Content   string     `json:"content" db:"content"`
    Language  string     `json:"language" db:"language"`
    Version   int        `json:"version" db:"version"`
    ParentID  *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
    CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}