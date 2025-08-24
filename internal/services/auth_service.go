package services

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "time"
)

type AuthService struct {
    secretKey string
}

func NewAuthService(secretKey string) *AuthService {
    return &AuthService{secretKey: secretKey}
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorUnverifiable)
        }
        return []byte(s.secretKey), nil
    })
}

func (s *AuthService) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.Request.Header.Get("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        token, err := s.ValidateToken(tokenString)
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", token.Claims.(jwt.MapClaims)["user_id"])
        c.Next()
    }
}