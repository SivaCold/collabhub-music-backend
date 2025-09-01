package handlers

import (
    "collabhub-music-backend/internal/services"
    "collabhub-music-backend/internal/utils"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
    projectService *services.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
    return &ProjectHandler{
        projectService: projectService,
    }
}

// GetProjects retrieves all projects for the authenticated user
// @Summary Get user projects
// @Description Get all projects where user is owner or collaborator
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponse{data=[]models.Project}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projects, err := h.projectService.GetUserProjects(parsedUserID)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Projects retrieved successfully", projects)
}

// CreateProject creates a new project
// @Summary Create project
// @Description Create a new music project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body services.CreateProjectRequest true "Project data"
// @Success 201 {object} utils.SuccessResponse{data=models.Project}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    var req services.CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    project, err := h.projectService.CreateProject(parsedUserID, &req)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "Project created successfully", project)
}

// GetProject retrieves a specific project
// @Summary Get project
// @Description Get project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} utils.SuccessResponse{data=models.Project}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
        return
    }

    project, err := h.projectService.GetProject(parsedUserID, projectID)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Project retrieved successfully", project)
}

// UpdateProject updates an existing project
// @Summary Update project
// @Description Update project information
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param request body services.UpdateProjectRequest true "Project update data"
// @Success 200 {object} utils.SuccessResponse{data=models.Project}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
        return
    }

    var req services.UpdateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    project, err := h.projectService.UpdateProject(parsedUserID, projectID, &req)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Project updated successfully", project)
}

// DeleteProject deletes a project
// @Summary Delete project
// @Description Delete a project (only owner)
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
        return
    }

    err = h.projectService.DeleteProject(parsedUserID, projectID)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Project deleted successfully", nil)
}

// AddCollaborator adds a collaborator to a project
// @Summary Add collaborator
// @Description Add a user as collaborator to a project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param request body AddCollaboratorRequest true "Collaborator data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /projects/{id}/collaborators [post]
func (h *ProjectHandler) AddCollaborator(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
        return
    }

    var req AddCollaboratorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    collaboratorID, err := uuid.Parse(req.UserID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid collaborator ID", err)
        return
    }

    err = h.projectService.AddCollaborator(parsedUserID, projectID, collaboratorID, req.Role)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Collaborator added successfully", nil)
}

// RemoveCollaborator removes a collaborator from a project
// @Summary Remove collaborator
// @Description Remove a collaborator from a project
// @Tags projects
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Project ID"
// @Param userId path string true "User ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /projects/{id}/collaborators/{userId} [delete]
func (h *ProjectHandler) RemoveCollaborator(c *gin.Context) {
    userID := c.GetString("user_id")
    parsedUserID, err := uuid.Parse(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }

    projectID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid project ID", err)
        return
    }

    collaboratorID, err := uuid.Parse(c.Param("userId"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid collaborator ID", err)
        return
    }

    err = h.projectService.RemoveCollaborator(parsedUserID, projectID, collaboratorID)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Collaborator removed successfully", nil)
}

// Request structs
type AddCollaboratorRequest struct {
    UserID string `json:"user_id" binding:"required"`
    Role   string `json:"role" binding:"required,oneof=admin collaborator viewer"`
}