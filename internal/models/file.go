package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// File represents a file in a project
type File struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ProjectID    uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
    BranchID     uuid.UUID `json:"branch_id" gorm:"type:uuid"`
    Name         string    `json:"name" gorm:"not null"`
    OriginalName string    `json:"original_name"`
    Path         string    `json:"path"`
    FileType     string    `json:"file_type"` // audio, image, video, document, code, other
    MimeType     string    `json:"mime_type"`
    Size         int64     `json:"size"`
    Checksum     string    `json:"checksum"`
    StoragePath  string    `json:"storage_path"` // Path on disk
    IsPublic     bool      `json:"is_public" gorm:"default:false"`
    UploadedBy   uuid.UUID `json:"uploaded_by" gorm:"type:uuid;not null"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Project       Project        `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
    Branch        Branch         `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
    Uploader      User           `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
    Versions      []FileVersion  `json:"versions,omitempty" gorm:"foreignKey:FileID"`
    AudioMetadata *AudioMetadata `json:"audio_metadata,omitempty" gorm:"foreignKey:FileID"`
}

// FileVersion represents different versions of a file
type FileVersion struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FileID      uuid.UUID `json:"file_id" gorm:"type:uuid;not null"`
    Version     int       `json:"version" gorm:"not null"`
    StoragePath string    `json:"storage_path"`
    Size        int64     `json:"size"`
    Checksum    string    `json:"checksum"`
    Comment     string    `json:"comment"`
    CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
    CreatedAt   time.Time `json:"created_at"`

    // Relationships
    File    File `json:"file,omitempty" gorm:"foreignKey:FileID"`
    Creator User `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

// AudioMetadata represents metadata for audio files
type AudioMetadata struct {
    ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FileID   uuid.UUID `json:"file_id" gorm:"type:uuid;not null;uniqueIndex"`
    Title    string    `json:"title"`
    Artist   string    `json:"artist"`
    Album    string    `json:"album"`
    Genre    string    `json:"genre"`
    Year     int       `json:"year"`
    Track    int       `json:"track"`
    Duration float64   `json:"duration"` // in seconds
    BitRate  int       `json:"bit_rate"`
    SampleRate int     `json:"sample_rate"`
    Channels int       `json:"channels"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // Relationships
    File File `json:"file,omitempty" gorm:"foreignKey:FileID"`
}

// BeforeCreate hooks
func (f *File) BeforeCreate(tx *gorm.DB) error {
    if f.ID == uuid.Nil {
        f.ID = uuid.New()
    }
    return nil
}

func (fv *FileVersion) BeforeCreate(tx *gorm.DB) error {
    if fv.ID == uuid.Nil {
        fv.ID = uuid.New()
    }
    return nil
}

func (am *AudioMetadata) BeforeCreate(tx *gorm.DB) error {
    if am.ID == uuid.Nil {
        am.ID = uuid.New()
    }
    return nil
}