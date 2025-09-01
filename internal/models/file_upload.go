package models

import (
    "time"

    "github.com/google/uuid"
)

// FileUpload represents an uploaded file
type FileUpload struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Filename    string    `json:"filename" gorm:"not null"`
    OriginalName string   `json:"original_name" gorm:"not null"`
    ContentType string    `json:"content_type"`
    Size        int64     `json:"size"`
    Path        string    `json:"path" gorm:"not null"`
    IsExtracted bool      `json:"is_extracted" gorm:"default:false"`
    ProjectID   *uuid.UUID `json:"project_id,omitempty" gorm:"type:uuid"`
    UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// ZipValidationResult represents ZIP file validation result
type ZipValidationResult struct {
    IsValid          bool     `json:"is_valid"`
    Error            string   `json:"error,omitempty"`
    TotalFiles       int      `json:"total_files"`
    AudioFiles       int      `json:"audio_files"`
    Folders          int      `json:"folders"`
    TotalSize        int64    `json:"total_size"`
    SupportedFiles   []string `json:"supported_files"`
    UnsupportedFiles []string `json:"unsupported_files"`
}

// ZipFileInfo represents information about a file in ZIP
type ZipFileInfo struct {
    Name         string    `json:"name"`
    Path         string    `json:"path"`
    Size         int64     `json:"size"`
    IsDirectory  bool      `json:"is_directory"`
    ContentType  string    `json:"content_type"`
    IsAudioFile  bool      `json:"is_audio_file"`
    ModTime      time.Time `json:"mod_time"`
}

// ZipExtractionResult represents ZIP extraction result
type ZipExtractionResult struct {
    Success        bool          `json:"success"`
    ExtractedPath  string        `json:"extracted_path"`
    ExtractedFiles []ZipFileInfo `json:"extracted_files"`
    AudioFiles     []ZipFileInfo `json:"audio_files"`
    TotalFiles     int           `json:"total_files"`
    TotalSize      int64         `json:"total_size"`
    Error          string        `json:"error,omitempty"`
}

// ProjectFromZipRequest represents request to create project from ZIP
type ProjectFromZipRequest struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description,omitempty"`
    Genre       string `json:"genre,omitempty"`
    BPM         int    `json:"bpm,omitempty"`
    Key         string `json:"key,omitempty"`
}