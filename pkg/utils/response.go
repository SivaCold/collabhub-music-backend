package utils

import (
    "github.com/gin-gonic/gin"
)

// Response represents a standard API response structure
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

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