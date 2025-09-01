package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData represents a standard API response structure
type ResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, ResponseData{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	c.JSON(statusCode, ResponseData{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string, validationErrors interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": message,
		"errors":  validationErrors,
	})
}

// ParseJSON parses JSON request body into the given struct
func ParseJSON(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindJSON(v); err != nil {
		return err
	}
	return nil
}
