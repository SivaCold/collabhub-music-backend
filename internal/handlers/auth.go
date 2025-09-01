package handlers

import "github.com/gin-gonic/gin"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
    c.JSON(200, map[string]interface{}{
        "success": true,
        "message": "login endpoint",
    })
}

func (h *AuthHandler) Register(c *gin.Context) {
    c.JSON(200, map[string]interface{}{
        "success": true,
        "message": "register endpoint",
    })
}

func (h *AuthHandler) Logout(c *gin.Context) {
    c.JSON(200, map[string]interface{}{
        "success": true,
        "message": "logout endpoint",
    })
}
