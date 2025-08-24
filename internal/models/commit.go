package models

import (
    "time"
    "github.com/google/uuid"
)

type Commit struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    BranchID    uuid.UUID  `json:"branch_id" db:"branch_id"`
    ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
    Message     string     `json:"message" db:"message"`
    Description *string    `json:"description,omitempty" db:"description"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
    AuthorID    uuid.UUID  `json:"author_id" db:"author_id"`
    Hash        string     `json:"hash" db:"hash"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

type CommitFile struct {
    ID       uuid.UUID `json:"id" db:"id"`
    CommitID uuid.UUID `json:"commit_id" db:"commit_id"`
    FileID   uuid.UUID `json:"file_id" db:"file_id"`
    Action   string    `json:"action" db:"action"` // added, modified, deleted
}