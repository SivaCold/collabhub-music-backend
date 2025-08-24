package models

import (
    "time"
    "github.com/google/uuid"
)

type Album struct {
    ID          uuid.UUID `json:"id" db:"id"`
    ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
    Title       string    `json:"title" db:"title"`
    Description *string   `json:"description,omitempty" db:"description"`
    CoverArt    *string   `json:"cover_art,omitempty" db:"cover_art"`
    ReleaseDate *time.Time `json:"release_date,omitempty" db:"release_date"`
    Status      string    `json:"status" db:"status"` // draft, recording, mixing, mastered, released
    CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type AlbumTrack struct {
    ID       uuid.UUID `json:"id" db:"id"`
    AlbumID  uuid.UUID `json:"album_id" db:"album_id"`
    TrackID  uuid.UUID `json:"track_id" db:"track_id"`
    Position int       `json:"position" db:"position"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}