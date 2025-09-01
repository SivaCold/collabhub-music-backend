package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidateJSON middleware validates JSON request body against struct tags
func ValidateJSON(v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(v); err != nil {
			var errors []ValidationError

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, validationError := range validationErrors {
					errors = append(errors, ValidationError{
						Field:   validationError.Field(),
						Message: getValidationMessage(validationError),
						Value:   validationError.Value(),
					})
				}
			}

			c.JSON(400, gin.H{
				"success": false,
				"message": "Validation failed",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		// Store validated data in context
		c.Set("validatedData", v)
		c.Next()
	}
}

// ValidateQuery validates query parameters
func ValidateQuery(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindQuery(v); err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}

// ValidateParams validates path parameters
func ValidateParams(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindUri(v); err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Value is too short (minimum " + fe.Param() + " characters)"
	case "max":
		return "Value is too long (maximum " + fe.Param() + " characters)"
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	case "uuid4":
		return "Must be a valid UUID"
	case "oneof":
		return "Value must be one of: " + fe.Param()
	default:
		return "Invalid value"
	}
}

// ValidationMiddleware returns a middleware that validates request data
func ValidationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add validation helper to context
		c.Set("validator", validate)
		c.Next()
	})
}
