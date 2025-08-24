package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "collabhub-music-backend/internal/services"
)

type AuthMiddleware struct {
    jwtMiddleware   *JWTMiddleware
    keycloakService *services.KeycloakService
    userService     *services.UserService
}

func NewAuthMiddleware(jwtMiddleware *JWTMiddleware, keycloakService *services.KeycloakService, userService *services.UserService) *AuthMiddleware {
    return &AuthMiddleware{
        jwtMiddleware:   jwtMiddleware,
        keycloakService: keycloakService,
        userService:     userService,
    }
}

// RequireAuth valide le token Keycloak et synchronise l'utilisateur
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        // Valider le token avec Keycloak
        isValid, err := a.keycloakService.ValidateToken(c.Request.Context(), tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to validate token"})
            c.Abort()
            return
        }

        if !isValid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Synchroniser l'utilisateur depuis Keycloak
        user, err := a.userService.SyncUserFromKeycloak(c.Request.Context(), tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to sync user data"})
            c.Abort()
            return
        }

        // Stocker les informations utilisateur dans le contexte
        c.Set("user_id", user.ID.String())
        c.Set("keycloak_id", user.KeycloakID)
        c.Set("username", user.Username)
        c.Set("email", user.Email)
        c.Set("user", user)

        c.Next()
    }
}

// RequireJWT utilise le middleware JWT pour valider les tokens Keycloak
func (a *AuthMiddleware) RequireJWT() gin.HandlerFunc {
    return a.jwtMiddleware.ValidateJWT()
}

// OptionalAuth permet l'accès avec ou sans token
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // Pas de token, continuer sans authentification
            c.Next()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            // Format de token invalide, continuer sans authentification
            c.Next()
            return
        }

        // Essayer de valider le token
        isValid, err := a.keycloakService.ValidateToken(c.Request.Context(), tokenString)
        if err != nil || !isValid {
            // Token invalide, continuer sans authentification
            c.Next()
            return
        }

        // Token valide, synchroniser l'utilisateur
        user, err := a.userService.SyncUserFromKeycloak(c.Request.Context(), tokenString)
        if err == nil && user != nil {
            c.Set("user_id", user.ID.String())
            c.Set("keycloak_id", user.KeycloakID)
            c.Set("username", user.Username)
            c.Set("email", user.Email)
            c.Set("user", user)
        }

        c.Next()
    }
}

// RequireRole vérifie si l'utilisateur a un rôle spécifique
func (a *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles, exists := c.Get("roles")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No roles found"})
            c.Abort()
            return
        }

        userRoles, ok := roles.([]string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
            c.Abort()
            return
        }

        hasRole := false
        for _, role := range userRoles {
            if role == requiredRole {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// RequireAnyRole vérifie si l'utilisateur a au moins un des rôles requis
func (a *AuthMiddleware) RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles, exists := c.Get("roles")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No roles found"})
            c.Abort()
            return
        }

        userRoles, ok := roles.([]string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
            c.Abort()
            return
        }

        hasAnyRole := false
        for _, userRole := range userRoles {
            for _, requiredRole := range requiredRoles {
                if userRole == requiredRole {
                    hasAnyRole = true
                    break
                }
            }
            if hasAnyRole {
                break
            }
        }

        if !hasAnyRole {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// GetCurrentUser récupère l'utilisateur depuis le contexte
func GetCurrentUser(c *gin.Context) (*services.UserService, bool) {
    user, exists := c.Get("user")
    if !exists {
        return nil, false
    }

    userModel, ok := user.(*services.UserService)
    if !ok {
        return nil, false
    }

    return userModel, true
}

// GetCurrentUserID récupère l'ID de l'utilisateur depuis le contexte
func GetCurrentUserID(c *gin.Context) (string, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return "", false
    }

    userIDStr, ok := userID.(string)
    if !ok {
        return "", false
    }

    return userIDStr, true
}

// IsAuthenticated vérifie si l'utilisateur est authentifié
func IsAuthenticated(c *gin.Context) bool {
    _, exists := c.Get("user_id")
    return exists
}