package models

import (
    "time"
    "github.com/google/uuid"
)

type Tag struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Color       *string   `json:"color,omitempty" db:"color"`
    Description *string   `json:"description,omitempty" db:"description"`
    CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TrackTag struct {
    ID        uuid.UUID `json:"id" db:"id"`
    TrackID   uuid.UUID `json:"track_id" db:"track_id"`
    TagID     uuid.UUID `json:"tag_id" db:"tag_id"`
    CreatedBy uuid.UUID `json:"created_by" db:"created_by"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}