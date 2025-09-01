package utils

// APIResponse represents a successful API response
type APIResponse struct {
    Status  string      `json:"status" example:"success"`
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// APIError represents an error API response
type APIError struct {
    Status  string `json:"status" example:"error"`
    Error   string `json:"error" example:"Something went wrong"`
    Code    int    `json:"code,omitempty" example:"400"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) APIResponse {
    return APIResponse{
        Status: "success",
        Data:   data,
    }
}

// SuccessResponseWithMessage creates a success response with message
func SuccessResponseWithMessage(data interface{}, message string) APIResponse {
    return APIResponse{
        Status:  "success",
        Data:    data,
        Message: message,
    }
}

// ErrorResponse creates an error response
func ErrorResponse(message string) APIError {
    return APIError{
        Status: "error",
        Error:  message,
    }
}

// ErrorResponseWithCode creates an error response with status code
func ErrorResponseWithCode(message string, code int) APIError {
    return APIError{
        Status: "error",
        Error:  message,
        Code:   code,
    }
}