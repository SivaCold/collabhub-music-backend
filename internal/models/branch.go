package models

import (
    "time"
    "github.com/google/uuid"
)

type Branch struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
    Name        string     `json:"name" db:"name"`
    Description *string    `json:"description,omitempty" db:"description"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
    IsDefault   bool       `json:"is_default" db:"is_default"`
    Status      string     `json:"status" db:"status"` // active, merged, archived
    CreatedBy   uuid.UUID  `json:"created_by" db:"created_by"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}