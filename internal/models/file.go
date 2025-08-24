package models

import (
    "time"
    "github.com/google/uuid"
)

type File struct {
    ID          uuid.UUID `json:"id" db:"id"`
    ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
    UserID      uuid.UUID `json:"user_id" db:"user_id"`
    Name        string    `json:"name" db:"name"`
    Path        string    `json:"path" db:"path"`
    Size        int64     `json:"size" db:"size"`
    MimeType    string    `json:"mime_type" db:"mime_type"`
    FileType    string    `json:"file_type" db:"file_type"` // audio, midi, sheet, lyrics, video
    Version     int       `json:"version" db:"version"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
    BranchID    uuid.UUID `json:"branch_id" db:"branch_id"`
    CommitID    *uuid.UUID `json:"commit_id,omitempty" db:"commit_id"`
    Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}