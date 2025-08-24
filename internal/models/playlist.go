package models

import (
    "time"
    "github.com/google/uuid"
)

type Playlist struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description *string   `json:"description,omitempty" db:"description"`
    UserID      uuid.UUID `json:"user_id" db:"user_id"`
    IsPublic    bool      `json:"is_public" db:"is_public"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type PlaylistTrack struct {
    ID         uuid.UUID `json:"id" db:"id"`
    PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id"`
    TrackID    uuid.UUID `json:"track_id" db:"track_id"`
    Position   int       `json:"position" db:"position"`
    AddedBy    uuid.UUID `json:"added_by" db:"added_by"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}