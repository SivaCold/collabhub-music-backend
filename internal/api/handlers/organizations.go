package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/services"
    "collabhub-music-backend/internal/middleware"
)

type OrganizationHandler struct {
    service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
    return &OrganizationHandler{service: service}
}

// CreateOrganization handles the creation of a new organization
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
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

    var org models.Organization
    if err := c.ShouldBindJSON(&org); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Validation
    if org.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Organization name is required"})
        return
    }

    // Set creator
    org.CreatedBy = userID

    if err := h.service.CreateOrganization(c.Request.Context(), &org); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization: " + err.Error()})
        return
    }

    // Return organization data
    response := gin.H{
        "id":          org.ID,
        "name":        org.Name,
        "description": org.Description,
        "created_by":  org.CreatedBy,
        "created_at":  org.CreatedAt,
        "updated_at":  org.UpdatedAt,
    }

    c.JSON(http.StatusCreated, response)
}

// GetOrganization handles retrieving an organization by ID
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
    idParam := c.Param("id")
    orgID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
        return
    }

    org, err := h.service.GetOrganizationByID(c.Request.Context(), orgID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
        return
    }

    c.JSON(http.StatusOK, org)
}

// UpdateOrganization handles updating an existing organization
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
    idParam := c.Param("id")
    orgID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
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

    // Check if organization exists and user has access
    existingOrg, err := h.service.GetOrganizationByID(c.Request.Context(), orgID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
        return
    }

    // Check if user is the creator (basic access control)
    if existingOrg.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to update this organization"})
        return
    }

    var updateData models.Organization
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Set ID for update
    updateData.ID = orgID

    if err := h.service.UpdateOrganization(c.Request.Context(), &updateData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization: " + err.Error()})
        return
    }

    // Get updated organization
    updatedOrg, err := h.service.GetOrganizationByID(c.Request.Context(), orgID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated organization"})
        return
    }

    c.JSON(http.StatusOK, updatedOrg)
}

// DeleteOrganization handles deleting an organization
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
    idParam := c.Param("id")
    orgID, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
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

    // Check if organization exists and user has access
    existingOrg, err := h.service.GetOrganizationByID(c.Request.Context(), orgID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
        return
    }

    // Check if user is the creator (basic access control)
    if existingOrg.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to delete this organization"})
        return
    }

    if err := h.service.DeleteOrganization(c.Request.Context(), orgID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

// ListOrganizations handles retrieving organizations with pagination
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
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

    organizations, err := h.service.ListOrganizations(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organizations: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "organizations": organizations,
        "limit":         limit,
        "offset":        offset,
        "count":         len(organizations),
    })
}

// GetUserOrganizations handles retrieving organizations for the current user
func (h *OrganizationHandler) GetUserOrganizations(c *gin.Context) {
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

    organizations, err := h.service.GetOrganizationsByUserID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user organizations: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "organizations": organizations,
        "count":         len(organizations),
    })
}

// AddUserToOrganization handles adding a user to an organization
func (h *OrganizationHandler) AddUserToOrganization(c *gin.Context) {
    orgIDParam := c.Param("id")
    orgID, err := uuid.Parse(orgIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
        return
    }

    var requestData struct {
        UserID string `json:"user_id" binding:"required"`
    }

    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    userID, err := uuid.Parse(requestData.UserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    if err := h.service.AddUserToOrganization(c.Request.Context(), orgID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to organization: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User added to organization successfully"})
}

// RemoveUserFromOrganization handles removing a user from an organization
func (h *OrganizationHandler) RemoveUserFromOrganization(c *gin.Context) {
    orgIDParam := c.Param("id")
    orgID, err := uuid.Parse(orgIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
        return
    }

    userIDParam := c.Param("user_id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    if err := h.service.RemoveUserFromOrganization(c.Request.Context(), orgID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from organization: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User removed from organization successfully"})
}

// GetOrganizationByName handles retrieving an organization by name
func (h *OrganizationHandler) GetOrganizationByName(c *gin.Context) {
    name := c.Query("name")
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Organization name parameter is required"})
        return
    }

    org, err := h.service.GetOrganizationByName(c.Request.Context(), name)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
        return
    }

    c.JSON(http.StatusOK, org)
}