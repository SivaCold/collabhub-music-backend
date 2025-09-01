package main

import (
    "log"
    "os"

    "collabhub-music-backend/internal/handlers"
    "collabhub-music-backend/internal/services"

    "github.com/gin-gonic/gin"
)

func main() {
    // Create upload directories
    uploadPath := "uploads"
    zipUploadPath := "uploads/zips"
    extractPath := "uploads/extracted"
    
    os.MkdirAll(zipUploadPath, 0755)
    os.MkdirAll(extractPath, 0755)

    // Create Gin router
    r := gin.Default()
    
    // Set max form size (500MB for file uploads)
    r.MaxMultipartMemory = 500 << 20 // 500MB

    // Create services
    zipService := services.NewZipService(uploadPath, extractPath)

    // Create handlers
    authHandler := handlers.NewAuthHandler()
    zipHandler := handlers.NewZipHandler(zipService)

    // Setup routes
    api := r.Group("/api/v1")
    {
        // Authentication routes
        auth := api.Group("/auth")
        {
            auth.POST("/login", authHandler.Login)
            auth.POST("/register", authHandler.Register)
            auth.POST("/logout", authHandler.Logout)
        }

        // File upload and ZIP handling routes
        files := api.Group("/files")
        {
            // ZIP file operations
            zip := files.Group("/zip")
            {
                zip.POST("/upload", zipHandler.UploadZip)
                zip.GET("/:file_id/validate", zipHandler.ValidateZip)
                zip.GET("/:file_id/info", zipHandler.GetZipInfo)
                zip.POST("/:file_id/extract", zipHandler.ExtractZip)
                zip.POST("/:file_id/project", zipHandler.CreateProjectFromZip)
            }

            // Project file operations
            projects := files.Group("/projects")
            {
                projects.GET("/:project_id/files", zipHandler.ListExtractedFiles)
                projects.DELETE("/:project_id/cleanup", zipHandler.CleanupProject)
            }
        }

        // Health check
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status":  "ok",
                "message": "CollabHub Music Backend is running",
                "version": "1.0.0",
            })
        })
    }

    log.Println("Starting server on :8081")
    log.Println("Upload directory:", uploadPath)
    log.Println("Extract directory:", extractPath)
    
    if err := r.Run(":8081"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}