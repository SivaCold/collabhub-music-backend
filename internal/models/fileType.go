package models

import (
	"time"

	"github.com/google/uuid"
)

// FileType represents the type of file
type FileType string

const (
	FileTypeAudio    FileType = "audio"
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeCode     FileType = "code"
	FileTypeOther    FileType = "other"
)

// ProjectFile represents a file in a project
type ProjectFile struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	BranchID    uuid.UUID `json:"branch_id" db:"branch_id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Path        string    `json:"path" db:"path" binding:"required"`
	Type        FileType  `json:"type" db:"type"`
	Size        int64     `json:"size" db:"size"`
	MimeType    string    `json:"mime_type" db:"mime_type"`
	StoragePath string    `json:"storage_path" db:"storage_path"`
	Checksum    string    `json:"checksum" db:"checksum"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Audio-specific metadata
	AudioMetadata *AudioMetadata `json:"audio_metadata,omitempty"`

	// For API responses - file content is not stored in DB
	Content []byte `json:"content,omitempty" db:"-"`
	URL     string `json:"url,omitempty" db:"-"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	BranchID  uuid.UUID `json:"branch_id" binding:"required"`
	Path      string    `json:"path" binding:"required"`
}

// FileTreeNode represents a node in file tree structure
type FileTreeNode struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Type     string         `json:"type"` // "file" or "folder"
	Path     string         `json:"path"`
	Size     *int64         `json:"size,omitempty"`
	Children []FileTreeNode `json:"children,omitempty"`
	File     *ProjectFile   `json:"file,omitempty"`
}

// TableName returns the database table name for ProjectFile
func (ProjectFile) TableName() string {
	return "project_files"
}
