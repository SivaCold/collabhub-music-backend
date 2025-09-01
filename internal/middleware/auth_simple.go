package middleware

import (
	"strings"

	"collabhub-music-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT token authentication for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.UnauthorizedResponse(c, "Authorization header must start with 'Bearer '")
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			utils.UnauthorizedResponse(c, "Token is required")
			c.Abort()
			return
		}

		// TODO: Implement token validation with Keycloak
		// For now, we'll validate that a token exists
		// In a real implementation, you would:
		// 1. Validate the JWT token
		// 2. Check if it's not expired
		// 3. Verify the signature
		// 4. Extract user information from the token

		// Mock validation - replace with real Keycloak integration
		if isValidToken(token) {
			// Set user context (this would come from token claims)
			c.Set("user_id", "mock-user-id")
			c.Set("username", "mock-user")
			c.Set("email", "user@example.com")
			c.Next()
		} else {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
		}
	}
}

// isValidToken is a mock function - replace with real Keycloak token validation
func isValidToken(token string) bool {
	// Mock validation - in production, this should validate with Keycloak
	return len(token) > 10 // Very basic check
}

// OptionalAuthMiddleware provides optional authentication for routes that can work with or without auth
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if token != "" && isValidToken(token) {
				// Set authenticated user context
				c.Set("user_id", "mock-user-id")
				c.Set("username", "mock-user")
				c.Set("email", "user@example.com")
				c.Set("authenticated", true)
			}
		} else {
			// No auth provided, continue as anonymous user
			c.Set("authenticated", false)
		}

		c.Next()
	}
}
