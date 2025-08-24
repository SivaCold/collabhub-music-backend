package models

import (
    "time"
    "github.com/google/uuid"
)

type Track struct {
    ID          uuid.UUID `json:"id" db:"id"`
    ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
    Name        string    `json:"name" db:"name"`
    Artist      string    `json:"artist" db:"artist"`
    Duration    int       `json:"duration" db:"duration"` // in seconds
    BPM         *int      `json:"bpm,omitempty" db:"bpm"`
    Key         *string   `json:"key,omitempty" db:"key"`
    Genre       *string   `json:"genre,omitempty" db:"genre"`
    FileID      *uuid.UUID `json:"file_id,omitempty" db:"file_id"`
    LyricsID    *uuid.UUID `json:"lyrics_id,omitempty" db:"lyrics_id"`
    Status      string    `json:"status" db:"status"` // draft, recording, mixing, mastered, released
    CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}