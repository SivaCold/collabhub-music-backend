package models

import (
    "time"
    "github.com/google/uuid"
)

type Genre struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description *string   `json:"description,omitempty" db:"description"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"` // for sub-genres
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}