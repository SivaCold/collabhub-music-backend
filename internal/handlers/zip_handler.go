package handlers

import (
    "fmt"
    "net/http"
    "path/filepath"
    "strconv"

    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/services"
    "collabhub-music-backend/pkg/utils"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// ZipHandler handles ZIP file operations
type ZipHandler struct {
    zipService *services.ZipService
}

// NewZipHandler creates a new ZIP handler
func NewZipHandler(zipService *services.ZipService) *ZipHandler {
    return &ZipHandler{
        zipService: zipService,
    }
}

// UploadZip godoc
// @Summary Upload and validate ZIP file
// @Description Upload a ZIP file and validate its contents for audio files
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "ZIP file to upload"
// @Success 200 {object} utils.APIResponse{data=models.ZipValidationResult} "ZIP file validated successfully"
// @Failure 400 {object} utils.APIError "Bad request - invalid file"
// @Failure 413 {object} utils.APIError "File too large"
// @Failure 422 {object} utils.APIError "Validation failed"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/zip/upload [post]
func (h *ZipHandler) UploadZip(c *gin.Context) {
    // Get uploaded file
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(
            "No file uploaded",
        ))
        return
    }

    // Validate file type
    if filepath.Ext(file.Filename) != ".zip" {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("File must be a ZIP archive"))
        return
    }

    // Check file size (500MB limit)
    maxSize := int64(500 * 1024 * 1024) // 500MB
    if file.Size > maxSize {
        c.JSON(http.StatusRequestEntityTooLarge, utils.ErrorResponse("File size exceeds 500MB limit"))
        return
    }

    // Generate unique filename
    fileID := uuid.New()
    filename := fmt.Sprintf("%s_%s", fileID.String(), file.Filename)
    uploadPath := filepath.Join("uploads", "zips", filename)

    // Save uploaded file
    if err := c.SaveUploadedFile(file, uploadPath); err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save uploaded file"))
        return
    }

    // Validate ZIP contents
    validation, err := h.zipService.ValidateZip(uploadPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to validate ZIP file"))
        return
    }

    if !validation.IsValid {
        c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse(validation.Error))
        return
    }

    // Add file path to response
    response := struct {
        *models.ZipValidationResult
        FileID   string `json:"file_id"`
        FilePath string `json:"file_path"`
    }{
        ZipValidationResult: validation,
        FileID:             fileID.String(),
        FilePath:           uploadPath,
    }

    c.JSON(http.StatusOK, utils.SuccessResponse(response))
}

// ValidateZip godoc
// @Summary Validate ZIP file contents
// @Description Analyze ZIP file and return information about its contents
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_id path string true "File ID from upload response"
// @Success 200 {object} utils.APIResponse{data=models.ZipValidationResult} "ZIP validation result"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 404 {object} utils.APIError "File not found"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/zip/{file_id}/validate [get]
func (h *ZipHandler) ValidateZip(c *gin.Context) {
    fileID := c.Param("file_id")
    if fileID == "" {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("File ID is required"))
        return
    }

    // Find ZIP file by ID (simplified - in real app, store file info in database)
    zipPath := filepath.Join("uploads", "zips", fileID+"_*.zip")
    
    // In a real implementation, you would query the database for the file path
    // For now, we'll assume the file exists at the expected location
    matches, err := filepath.Glob(zipPath)
    if err != nil || len(matches) == 0 {
        c.JSON(http.StatusNotFound, utils.ErrorResponse("ZIP file not found"))
        return
    }

    validation, err := h.zipService.ValidateZip(matches[0])
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to validate ZIP file"))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse(validation))
}

// ExtractZip godoc
// @Summary Extract ZIP file
// @Description Extract ZIP file contents to project directory
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_id path string true "File ID from upload response"
// @Param project_id query string false "Project ID (if not provided, generates new UUID)"
// @Success 200 {object} utils.APIResponse{data=models.ZipExtractionResult} "ZIP extracted successfully"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 404 {object} utils.APIError "File not found"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/zip/{file_id}/extract [post]
func (h *ZipHandler) ExtractZip(c *gin.Context) {
    fileID := c.Param("file_id")
    if fileID == "" {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("File ID is required"))
        return
    }

    // Get project ID or generate new one
    var projectID uuid.UUID
    projectIDStr := c.Query("project_id")
    if projectIDStr != "" {
        parsedID, err := uuid.Parse(projectIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid project ID format"))
            return
        }
        projectID = parsedID
    } else {
        projectID = uuid.New()
    }

    // Find ZIP file
    zipPath := filepath.Join("uploads", "zips", fileID+"_*.zip")
    matches, err := filepath.Glob(zipPath)
    if err != nil || len(matches) == 0 {
        c.JSON(http.StatusNotFound, utils.ErrorResponse("ZIP file not found"))
        return
    }

    // Extract ZIP
    result, err := h.zipService.ExtractZip(matches[0], projectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to extract ZIP file"))
        return
    }

    if !result.Success {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse(result.Error))
        return
    }

    response := struct {
        *models.ZipExtractionResult
        ProjectID string `json:"project_id"`
    }{
        ZipExtractionResult: result,
        ProjectID:          projectID.String(),
    }

    c.JSON(http.StatusOK, utils.SuccessResponse(response))
}

// ListExtractedFiles godoc
// @Summary List extracted files
// @Description List all files in an extracted project directory
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param audio_only query boolean false "Return only audio files"
// @Success 200 {object} utils.APIResponse{data=[]models.ZipFileInfo} "List of extracted files"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 404 {object} utils.APIError "Project not found"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/projects/{project_id}/files [get]
func (h *ZipHandler) ListExtractedFiles(c *gin.Context) {
    projectIDStr := c.Param("project_id")
    projectID, err := uuid.Parse(projectIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid project ID format"))
        return
    }

    // Get audio_only parameter
    audioOnly, _ := strconv.ParseBool(c.Query("audio_only"))

    files, err := h.zipService.ListExtractedFiles(projectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to list extracted files"))
        return
    }

    // Filter for audio files only if requested
    if audioOnly {
        var audioFiles []models.ZipFileInfo
        for _, file := range files {
            if file.IsAudioFile {
                audioFiles = append(audioFiles, file)
            }
        }
        files = audioFiles
    }

    response := struct {
        ProjectID  string                `json:"project_id"`
        Files      []models.ZipFileInfo  `json:"files"`
        TotalFiles int                   `json:"total_files"`
        AudioFiles int                   `json:"audio_files"`
    }{
        ProjectID:  projectID.String(),
        Files:      files,
        TotalFiles: len(files),
    }

    // Count audio files
    for _, file := range files {
        if file.IsAudioFile {
            response.AudioFiles++
        }
    }

    c.JSON(http.StatusOK, utils.SuccessResponse(response))
}

// CreateProjectFromZip godoc
// @Summary Create project from ZIP
// @Description Create a new project by extracting and processing a ZIP file
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_id path string true "File ID from upload response"
// @Param project body models.ProjectFromZipRequest true "Project details"
// @Success 201 {object} utils.APIResponse{data=models.Project} "Project created successfully"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 404 {object} utils.APIError "File not found"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/zip/{file_id}/project [post]
func (h *ZipHandler) CreateProjectFromZip(c *gin.Context) {
    fileID := c.Param("file_id")
    if fileID == "" {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("File ID is required"))
        return
    }

    // Parse request body
    var req models.ProjectFromZipRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request data"))
        return
    }

    // Generate project ID
    projectID := uuid.New()

    // Find and extract ZIP file
    zipPath := filepath.Join("uploads", "zips", fileID+"_*.zip")
    matches, err := filepath.Glob(zipPath)
    if err != nil || len(matches) == 0 {
        c.JSON(http.StatusNotFound, utils.ErrorResponse("ZIP file not found"))
        return
    }

    // Extract ZIP
    extractResult, err := h.zipService.ExtractZip(matches[0], projectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to extract ZIP file"))
        return
    }

    if !extractResult.Success {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse(extractResult.Error))
        return
    }

    // Create project model
    project := models.Project{
        ID:          projectID,
        Name:        req.Name,
        Description: req.Description,
        // Add other project fields as needed
    }

    // In a real implementation, you would:
    // 1. Save the project to database
    // 2. Create file records for extracted files
    // 3. Process audio metadata
    // 4. Create initial project structure

    response := struct {
        models.Project
        ExtractedFiles int                   `json:"extracted_files"`
        AudioFiles     int                   `json:"audio_files"`
        ExtractedPath  string                `json:"extracted_path"`
        Files          []models.ZipFileInfo  `json:"files"`
    }{
        Project:        project,
        ExtractedFiles: extractResult.TotalFiles,
        AudioFiles:     len(extractResult.AudioFiles),
        ExtractedPath:  extractResult.ExtractedPath,
        Files:          extractResult.AudioFiles,
    }

    c.JSON(http.StatusCreated, utils.SuccessResponse(response))
}

// GetZipInfo godoc
// @Summary Get ZIP file information
// @Description Get detailed information about ZIP file contents without extracting
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file_id path string true "File ID from upload response"
// @Success 200 {object} utils.APIResponse{data=models.ZipValidationResult} "ZIP file information"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 404 {object} utils.APIError "File not found"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/zip/{file_id}/info [get]
func (h *ZipHandler) GetZipInfo(c *gin.Context) {
    fileID := c.Param("file_id")
    if fileID == "" {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("File ID is required"))
        return
    }

    // Find ZIP file
    zipPath := filepath.Join("uploads", "zips", fileID+"_*.zip")
    matches, err := filepath.Glob(zipPath)
    if err != nil || len(matches) == 0 {
        c.JSON(http.StatusNotFound, utils.ErrorResponse("ZIP file not found"))
        return
    }

    info, err := h.zipService.GetZipInfo(matches[0])
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get ZIP information"))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse(info))
}

// CleanupProject godoc
// @Summary Cleanup project files
// @Description Remove all extracted files for a project
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Success 200 {object} utils.APIResponse{data=string} "Files cleaned up successfully"
// @Failure 400 {object} utils.APIError "Bad request"
// @Failure 500 {object} utils.APIError "Internal server error"
// @Router /files/projects/{project_id}/cleanup [delete]
func (h *ZipHandler) CleanupProject(c *gin.Context) {
    projectIDStr := c.Param("project_id")
    projectID, err := uuid.Parse(projectIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid project ID format"))
        return
    }

    if err := h.zipService.CleanupExtractedFiles(projectID); err != nil {
        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to cleanup project files"))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Project files cleaned up successfully"))
}