package models

import (
    "time"
    "github.com/google/uuid"
)

type Comment struct {
    ID              uuid.UUID  `json:"id" db:"id"`
    ProjectID       *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
    TrackID         *uuid.UUID `json:"track_id,omitempty" db:"track_id"`
    FileID          *uuid.UUID `json:"file_id,omitempty" db:"file_id"`
    ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
    UserID          uuid.UUID  `json:"user_id" db:"user_id"`
    Content         string     `json:"content" db:"content"`
    Timestamp       *int       `json:"timestamp,omitempty" db:"timestamp"` // for audio comments in seconds
    CreatedAt       time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}