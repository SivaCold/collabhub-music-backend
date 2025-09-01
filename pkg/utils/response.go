package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response structure for React Native consumption
type Response struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message" example:"Operation completed successfully"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Success bool                   `json:"success" example:"false"`
	Message string                 `json:"message" example:"Validation failed"`
	Error   string                 `json:"error" example:"VALIDATION_ERROR"`
	Details []ValidationError      `json:"details,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
	Code    string `json:"code" example:"INVALID_FORMAT"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int  `json:"page" example:"1"`
	Limit      int  `json:"limit" example:"10"`
	Total      int  `json:"total" example:"100"`
	TotalPages int  `json:"total_pages" example:"10"`
	HasMore    bool `json:"has_more" example:"true"`
}

// SuccessResponse creates a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// SuccessResponseWithPagination creates a successful response with pagination
func SuccessResponseWithPagination(c *gin.Context, statusCode int, message string, data interface{}, pagination *Pagination) {
	response := Response{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
	c.JSON(statusCode, response)
}

// ErrorResponseJSON creates an error response
func ErrorResponseJSON(c *gin.Context, statusCode int, message string, errorCode string, details []ValidationError) {
	response := ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorCode,
		Details: details,
	}
	c.JSON(statusCode, response)
}

// BadRequestResponse creates a 400 Bad Request response
func BadRequestResponse(c *gin.Context, message string, details []ValidationError) {
	ErrorResponseJSON(c, http.StatusBadRequest, message, "BAD_REQUEST", details)
}

// UnauthorizedResponse creates a 401 Unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusUnauthorized, message, "UNAUTHORIZED", nil)
}

// ForbiddenResponse creates a 403 Forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusForbidden, message, "FORBIDDEN", nil)
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusNotFound, message, "NOT_FOUND", nil)
}

// InternalErrorResponse creates a 500 Internal Server Error response
func InternalErrorResponse(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusInternalServerError, message, "INTERNAL_ERROR", nil)
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(c *gin.Context, message string, details []ValidationError) {
	ErrorResponseJSON(c, http.StatusUnprocessableEntity, message, "VALIDATION_ERROR", details)
}

// ConflictResponse creates a 409 Conflict response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponseJSON(c, http.StatusConflict, message, "CONFLICT", nil)
}

// CreatePagination creates pagination metadata
func CreatePagination(page, limit, total int) *Pagination {
	totalPages := (total + limit - 1) / limit
	hasMore := page < totalPages

	return &Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}

// ParsePaginationParams extracts pagination parameters from query string
func ParsePaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if p := c.Query("page"); p != "" {
		if parsed, err := parsePositiveInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := parsePositiveInt(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}

// parsePositiveInt parses a string to a positive integer
func parsePositiveInt(s string) (int, error) {
	var result int
	for _, digit := range s {
		if digit < '0' || digit > '9' {
			return 0, gin.Error{}
		}
		result = result*10 + int(digit-'0')
	}
	return result, nil
}

// GetUserIDFromContext extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}

// GetUserFromContext extracts the user object from the Gin context
func GetUserFromContext(c *gin.Context) (interface{}, bool) {
	user, exists := c.Get("user")
	return user, exists
}

// Legacy function for backward compatibility
// SendResponse sends a JSON response to the client
func SendResponse(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response to the client
func SendError(c *gin.Context, statusCode int, message string) {
	SendResponse(c, statusCode, false, message, nil)
}
