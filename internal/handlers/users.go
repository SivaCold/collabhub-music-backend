package handlers

import (
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "collabhub-music-backend/internal/services"
    "collabhub-music-backend/internal/models"
    "collabhub-music-backend/internal/middleware"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}


// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user account with Keycloak integration
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserRegistrationRequest true "User registration data"
// @Success 201 {object} models.APIResponse{data=models.User} "User created successfully"
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 409 {object} models.APIError "User already exists"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Basic validation
    if user.Username == "" || user.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username and email are required"})
        return
    }

    if user.FirstName == "" || user.LastName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "First name and last name are required"})
        return
    }

    if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
        return
    }

    // Don't return sensitive data
    response := gin.H{
        "id":         user.ID,
        "username":   user.Username,
        "email":      user.Email,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "created_at": user.CreatedAt,
    }

    c.JSON(http.StatusCreated, response)
}

// UpdateUserProfile godoc
// @Summary Update current user profile
// @Description Update the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.UserUpdateRequest true "User update data"
// @Success 200 {object} models.APIResponse{data=models.User} "User updated successfully"
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 401 {object} models.APIError "Unauthorized"
// @Failure 404 {object} models.APIError "User not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /users/me [put]
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
    userIDParam := c.Param("id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    // Check if current user is updating their own profile
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    if currentUserID != userID.String() {
        c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update another user's profile"})
        return
    }

    var updateData models.User
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }

    // Set the ID for update
    updateData.ID = userID

    if err := h.userService.UpdateUser(c.Request.Context(), &updateData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
        return
    }

    // Get updated user data
    updatedUser, err := h.userService.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user"})
        return
    }

    // Don't return sensitive data
    response := gin.H{
        "id":         updatedUser.ID,
        "username":   updatedUser.Username,
        "email":      updatedUser.Email,
        "first_name": updatedUser.FirstName,
        "last_name":  updatedUser.LastName,
        "updated_at": updatedUser.UpdatedAt,
    }

    c.JSON(http.StatusOK, response)
}

// GetUserProfile handles retrieving a user profile
func (h *UserHandler) GetUserProfile(c *gin.Context) {
    userIDParam := c.Param("id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    user, err := h.userService.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Return public profile data only
    response := gin.H{
        "id":         user.ID,
        "username":   user.Username,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "created_at": user.CreatedAt,
    }

    c.JSON(http.StatusOK, response)
}

// GetCurrentUser handles retrieving the current authenticated user
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
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

    user, err := h.userService.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Return full profile data for current user
    response := gin.H{
        "id":         user.ID,
        "username":   user.Username,
        "email":      user.Email,
        "first_name": user.FirstName,
        "last_name":  user.LastName,
        "created_at": user.CreatedAt,
        "updated_at": user.UpdatedAt,
    }

    c.JSON(http.StatusOK, response)
}

// ListUsers handles retrieving a paginated list of users
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Parse query parameters
    limit := 10
    offset := 0

    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := parseInt(limitStr); err == nil && l > 0 {
            limit = l
        }
    }

    if offsetStr := c.Query("offset"); offsetStr != "" {
        if o, err := parseInt(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }

    users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    // Return public data only
    var response []gin.H
    for _, user := range users {
        response = append(response, gin.H{
            "id":         user.ID,
            "username":   user.Username,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "created_at": user.CreatedAt,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "users":  response,
        "limit":  limit,
        "offset": offset,
        "count":  len(response),
    })
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
    userIDParam := c.Param("id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    // Check if current user is deleting their own account
    currentUserID, exists := middleware.GetCurrentUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    if currentUserID != userID.String() {
        c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete another user's account"})
        return
    }

    if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Helper function to parse integers
func parseInt(s string) (int, error) {
    var result int
    if _, err := fmt.Sscanf(s, "%d", &result); err != nil {
        return 0, err
    }
    return result, nil
}