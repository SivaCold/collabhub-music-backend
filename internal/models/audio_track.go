package models

import (
    "github.com/google/uuid"
    "time"
)

type AudioTrack struct {
    ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    ProjectID      uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
    Name           string    `json:"name" gorm:"type:varchar(255);not null"`
    Description    string    `json:"description" gorm:"type:text"`
    FilePath       string    `json:"file_path" gorm:"type:varchar(1000);not null"`
    FileSize       int64     `json:"file_size"`
    DurationSeconds float64   `json:"duration_seconds"`
    Format         string    `json:"format" gorm:"type:varchar(10)"`
    SampleRate     int       `json:"sample_rate"`
    BitDepth       int       `json:"bit_depth"`
    Channels       int       `json:"channels"`
    TrackType      string    `json:"track_type" gorm:"type:varchar(20)"`
    Instrument     string    `json:"instrument" gorm:"type:varchar(50)"`
    VersionNumber  int       `json:"version_number" gorm:"default:1"`
    ParentTrackID  *uuid.UUID `json:"parent_track_id" gorm:"type:uuid"`
    CreatedBy      uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
    CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}