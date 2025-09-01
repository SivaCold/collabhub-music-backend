package services

import (
    "archive/zip"
    "fmt"
    "io"
    "mime"
    "os"
    "path/filepath"
    "strings"

    "collabhub-music-backend/internal/models"
    "github.com/google/uuid"
)

// ZipService handles ZIP file operations
type ZipService struct {
    uploadPath string
    extractPath string
}

// NewZipService creates a new ZIP service
func NewZipService(uploadPath, extractPath string) *ZipService {
    // Ensure directories exist
    os.MkdirAll(uploadPath, 0755)
    os.MkdirAll(extractPath, 0755)
    
    return &ZipService{
        uploadPath:  uploadPath,
        extractPath: extractPath,
    }
}

// ValidateZip validates a ZIP file and returns information about its contents
func (s *ZipService) ValidateZip(zipPath string) (*models.ZipValidationResult, error) {
    reader, err := zip.OpenReader(zipPath)
    if err != nil {
        return &models.ZipValidationResult{
            IsValid: false,
            Error:   fmt.Sprintf("Failed to open ZIP file: %v", err),
        }, nil
    }
    defer reader.Close()

    result := &models.ZipValidationResult{
        IsValid:          true,
        SupportedFiles:   []string{},
        UnsupportedFiles: []string{},
    }

    audioExtensions := map[string]bool{
        ".mp3":  true,
        ".wav":  true,
        ".flac": true,
        ".aac":  true,
        ".ogg":  true,
        ".m4a":  true,
        ".wma":  true,
    }

    for _, file := range reader.File {
        result.TotalFiles++
        result.TotalSize += int64(file.UncompressedSize64)

        if file.FileInfo().IsDir() {
            result.Folders++
            continue
        }

        ext := strings.ToLower(filepath.Ext(file.Name))
        
        if audioExtensions[ext] {
            result.AudioFiles++
            result.SupportedFiles = append(result.SupportedFiles, file.Name)
        } else if ext != "" { // Skip files without extensions (likely directories)
            result.UnsupportedFiles = append(result.UnsupportedFiles, file.Name)
        }
    }

    // Validation rules
    if result.TotalFiles == 0 {
        result.IsValid = false
        result.Error = "ZIP file is empty"
    } else if result.AudioFiles == 0 {
        result.IsValid = false
        result.Error = "No supported audio files found in ZIP"
    } else if result.TotalSize > 500*1024*1024 { // 500MB limit
        result.IsValid = false
        result.Error = "ZIP file is too large (max 500MB)"
    }

    return result, nil
}

// ExtractZip extracts a ZIP file to the specified directory
func (s *ZipService) ExtractZip(zipPath string, projectID uuid.UUID) (*models.ZipExtractionResult, error) {
    reader, err := zip.OpenReader(zipPath)
    if err != nil {
        return &models.ZipExtractionResult{
            Success: false,
            Error:   fmt.Sprintf("Failed to open ZIP file: %v", err),
        }, err
    }
    defer reader.Close()

    extractPath := filepath.Join(s.extractPath, projectID.String())
    if err := os.MkdirAll(extractPath, 0755); err != nil {
        return &models.ZipExtractionResult{
            Success: false,
            Error:   fmt.Sprintf("Failed to create extraction directory: %v", err),
        }, err
    }

    result := &models.ZipExtractionResult{
        Success:        true,
        ExtractedPath:  extractPath,
        ExtractedFiles: []models.ZipFileInfo{},
        AudioFiles:     []models.ZipFileInfo{},
    }

    audioExtensions := map[string]bool{
        ".mp3":  true,
        ".wav":  true,
        ".flac": true,
        ".aac":  true,
        ".ogg":  true,
        ".m4a":  true,
        ".wma":  true,
    }

    for _, file := range reader.File {
        extractedPath := filepath.Join(extractPath, file.Name)
        
        // Security check: prevent directory traversal
        if !strings.HasPrefix(extractedPath, extractPath) {
            continue
        }

        fileInfo := models.ZipFileInfo{
            Name:        filepath.Base(file.Name),
            Path:        file.Name,
            Size:        int64(file.UncompressedSize64),
            IsDirectory: file.FileInfo().IsDir(),
            ModTime:     file.FileInfo().ModTime(),
        }

        if file.FileInfo().IsDir() {
            if err := os.MkdirAll(extractedPath, file.FileInfo().Mode()); err != nil {
                result.Error = fmt.Sprintf("Failed to create directory: %v", err)
                continue
            }
        } else {
            // Ensure parent directory exists
            if err := os.MkdirAll(filepath.Dir(extractedPath), 0755); err != nil {
                result.Error = fmt.Sprintf("Failed to create parent directory: %v", err)
                continue
            }

            // Extract file
            if err := s.extractFile(file, extractedPath); err != nil {
                result.Error = fmt.Sprintf("Failed to extract file %s: %v", file.Name, err)
                continue
            }

            // Set file info
            ext := strings.ToLower(filepath.Ext(file.Name))
            fileInfo.ContentType = mime.TypeByExtension(ext)
            fileInfo.IsAudioFile = audioExtensions[ext]

            if fileInfo.IsAudioFile {
                result.AudioFiles = append(result.AudioFiles, fileInfo)
            }
        }

        result.ExtractedFiles = append(result.ExtractedFiles, fileInfo)
        result.TotalFiles++
        result.TotalSize += fileInfo.Size
    }

    return result, nil
}

// extractFile extracts a single file from ZIP
func (s *ZipService) extractFile(file *zip.File, destPath string) error {
    reader, err := file.Open()
    if err != nil {
        return err
    }
    defer reader.Close()

    writer, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
    if err != nil {
        return err
    }
    defer writer.Close()

    _, err = io.Copy(writer, reader)
    return err
}

// GetZipInfo returns information about ZIP contents without extracting
func (s *ZipService) GetZipInfo(zipPath string) (*models.ZipValidationResult, error) {
    return s.ValidateZip(zipPath)
}

// CleanupExtractedFiles removes extracted files for a project
func (s *ZipService) CleanupExtractedFiles(projectID uuid.UUID) error {
    extractPath := filepath.Join(s.extractPath, projectID.String())
    return os.RemoveAll(extractPath)
}

// ListExtractedFiles lists all files in an extracted project directory
func (s *ZipService) ListExtractedFiles(projectID uuid.UUID) ([]models.ZipFileInfo, error) {
    extractPath := filepath.Join(s.extractPath, projectID.String())
    
    var files []models.ZipFileInfo
    
    err := filepath.Walk(extractPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Get relative path from extract directory
        relPath, err := filepath.Rel(extractPath, path)
        if err != nil {
            return err
        }

        // Skip root directory
        if relPath == "." {
            return nil
        }

        fileInfo := models.ZipFileInfo{
            Name:        info.Name(),
            Path:        relPath,
            Size:        info.Size(),
            IsDirectory: info.IsDir(),
            ModTime:     info.ModTime(),
        }

        if !info.IsDir() {
            ext := strings.ToLower(filepath.Ext(info.Name()))
            fileInfo.ContentType = mime.TypeByExtension(ext)
            
            audioExtensions := map[string]bool{
                ".mp3": true, ".wav": true, ".flac": true,
                ".aac": true, ".ogg": true, ".m4a": true, ".wma": true,
            }
            fileInfo.IsAudioFile = audioExtensions[ext]
        }

        files = append(files, fileInfo)
        return nil
    })

    return files, err
}