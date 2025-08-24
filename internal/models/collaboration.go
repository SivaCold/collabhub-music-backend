package models

import (
    "time"
    "github.com/google/uuid"
)

type Collaboration struct {
    ID          uuid.UUID `json:"id" db:"id"`
    ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
    UserID      uuid.UUID `json:"user_id" db:"user_id"`
    Role        string    `json:"role" db:"role"` // owner, admin, collaborator, viewer
    Permissions []string  `json:"permissions" db:"permissions"`
    InvitedBy   uuid.UUID `json:"invited_by" db:"invited_by"`
    Status      string    `json:"status" db:"status"` // pending, accepted, declined
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}