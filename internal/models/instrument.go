package models

import (
    "time"
    "github.com/google/uuid"
)

type Instrument struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Category    string    `json:"category" db:"category"` // strings, brass, woodwind, percussion, electronic, etc.
    Description *string   `json:"description,omitempty" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TrackInstrument struct {
    ID           uuid.UUID `json:"id" db:"id"`
    TrackID      uuid.UUID `json:"track_id" db:"track_id"`
    InstrumentID uuid.UUID `json:"instrument_id" db:"instrument_id"`
    UserID       uuid.UUID `json:"user_id" db:"user_id"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}