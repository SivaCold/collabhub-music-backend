package handlers

import (
    
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/services"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/middleware"
)

type ProjectHandler struct {
    projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
    return &ProjectHandler{projectService: projectService}
}

// CreateProject handles the creation of a new music project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
    // Get authenticated user
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID, err := uuid.Parse(currentUserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var project models.Project
    if err := c.ShouldBindJSON(&project); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Validation
    if project.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Project name is required"})
        return
    }

    // Set creator
    project.CreatedBy = userID

    if err := h.projectService.CreateProject(c.Request.Context(), &project); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project: " + err.Error()})
        return
    }

    // Return project without sensitive data
    response := gin.H{
        "id":              project.ID,
        "name":            project.Name,
        "description":     project.Description,
        "organization_id": project.OrganizationID,
        "created_by":      project.CreatedBy,
        "created_at":      project.CreatedAt,
        "updated_at":      project.UpdatedAt,
    }

    c.JSON(http.StatusCreated, response)
}

// GetProject handles retrieving a music project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
    idParam := c.Param("id")
    projectID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }

    project, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }

    // TODO: Add access control check here
    // For now, return the project
    c.JSON(http.StatusOK, project)
}

// UpdateProject handles updating an existing music project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
    idParam := c.Param("id")
    projectID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }

    // Get authenticated user
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID, err := uuid.Parse(currentUserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Check if project exists and user has access
    existingProject, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }

    // Check if user is the creator (basic access control)
    if existingProject.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to update this project"})
        return
    }

    var updateData models.Project
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Set ID for update
    updateData.ID = projectID

    if err := h.projectService.UpdateProject(c.Request.Context(), &updateData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project: " + err.Error()})
        return
    }

    // Get updated project
    updatedProject, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated project"})
        return
    }

    c.JSON(http.StatusOK, updatedProject)
}

// DeleteProject handles deleting a music project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
    idParam := c.Param("id")
    projectID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }

    // Get authenticated user
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID, err := uuid.Parse(currentUserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Check if project exists and user has access
    existingProject, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }

    // Check if user is the creator (basic access control)
    if existingProject.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to delete this project"})
        return
    }

    if err := h.projectService.DeleteProject(c.Request.Context(), projectID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ListProjects handles retrieving projects with pagination
func (h *ProjectHandler) ListProjects(c *gin.Context) {
    // Parse pagination parameters
    limit := 10
    offset := 0

    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }

    if offsetStr := c.Query("offset"); offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }

    projects, err := h.projectService.ListProjects(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "projects": projects,
        "limit":    limit,
        "offset":   offset,
        "count":    len(projects),
    })
}

// GetUserProjects handles retrieving projects for the current user
func (h *ProjectHandler) GetUserProjects(c *gin.Context) {
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID, err := uuid.Parse(currentUserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    projects, err := h.projectService.GetProjectsByUserID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user projects: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "projects": projects,
        "count":    len(projects),
    })
}

// SearchProjects handles project search
func (h *ProjectHandler) SearchProjects(c *gin.Context) {
    searchName := c.Query("name")
    if searchName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search name parameter is required"})
        return
    }

    projects, err := h.projectService.SearchProjectsByName(c.Request.Context(), searchName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search projects: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "projects": projects,
        "count":    len(projects),
        "query":    searchName,
    })
}

// GetProjectStats handles retrieving project statistics
func (h *ProjectHandler) GetProjectStats(c *gin.Context) {
    idParam := c.Param("id")
    projectID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
        return
    }

    stats, err := h.projectService.GetProjectStats(c.Request.Context(), projectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project stats: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, stats)
}