// @title CollabHub Music API
// @version 1.0.0
// @description A collaborative platform for musicians with GitLab-like versioning system for music projects
// @termsOfService https://collabhub-music.com/terms

// @contact.name CollabHub Music Support
// @contact.url https://collabhub-music.com/support
// @contact.email support@collabhub-music.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8444
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


package main

import (
    "context"
    "crypto/tls"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "collabhub-music-backend/internal/config"
    "collabhub-music-backend/internal/database"
    "collabhub-music-backend/internal/handlers"
    "collabhub-music-backend/internal/middleware"
    "collabhub-music-backend/internal/repository"
    "collabhub-music-backend/internal/services"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
    // Load environment variables
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development"
    }
    
    envFile := ".env." + env
    if err := godotenv.Load(envFile); err != nil {
        log.Printf("Warning: Could not load %s file: %v", envFile, err)
        // Try loading default .env file
        if err := godotenv.Load(); err != nil {
            log.Printf("Warning: Could not load default .env file: %v", err)
        }
    }

    // Load configuration
    cfg := config.Load()

    // Set Gin mode
    gin.SetMode(cfg.Server.GinMode)

    // Initialize database connection
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Run migrations
    if err := database.RunMigrations(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepo(db)
    projectRepo := repository.NewProjectRepo(db)
    organizationRepo := repository.NewOrganizationRepo(db)

    // Initialize services
    keycloakService := services.NewKeycloakService(
        cfg.Keycloak.URL,
        cfg.Keycloak.Realm,
        cfg.Keycloak.ClientID,
        cfg.Keycloak.ClientSecret,
    )
    userService := services.NewUserService(userRepo, keycloakService)
    projectService := services.NewProjectService(projectRepo)
    organizationService := services.NewOrganizationService(organizationRepo, userRepo)

    // Initialize middleware
    jwtMiddleware := middleware.NewJWTMiddleware(cfg.Keycloak.URL, cfg.Keycloak.Realm)
    authMiddleware := middleware.NewAuthMiddleware(jwtMiddleware, keycloakService, userService)

    // Initialize handlers
    userHandler := handlers.NewUserHandler(userService)
    projectHandler := handlers.NewProjectHandler(projectService)
    organizationHandler := handlers.NewOrganizationHandler(organizationService)

    // Setup router
    r := gin.New()
    
    // Add global middleware
    r.Use(middleware.LoggerMiddleware())
    r.Use(gin.Recovery())
    r.Use(middleware.CORSMiddleware())

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Public routes
    public := r.Group("/api/v1")
    {
        public.GET("/health", healthCheck)
        public.POST("/users/register", userHandler.RegisterUser)
    }

    // Protected routes
    protected := r.Group("/api/v1")
    protected.Use(authMiddleware.RequireAuth())
    {
        // Users
        users := protected.Group("/users")
        {
            users.GET("/me", userHandler.GetCurrentUser)
            users.PUT("/me", userHandler.UpdateUserProfile)
            users.DELETE("/me", userHandler.DeleteUser)
            users.GET("/:id", userHandler.GetUserProfile)
            users.GET("", userHandler.ListUsers)
        }

        // Projects
        projects := protected.Group("/projects")
        {
            projects.GET("", projectHandler.ListProjects)
            projects.POST("", projectHandler.CreateProject)
            projects.GET("/user", projectHandler.GetUserProjects)
            projects.GET("/search", projectHandler.SearchProjects)
            projects.GET("/:id", projectHandler.GetProject)
            projects.PUT("/:id", projectHandler.UpdateProject)
            projects.DELETE("/:id", projectHandler.DeleteProject)
            projects.GET("/:id/stats", projectHandler.GetProjectStats)
        }
        
        // Organizations
        organizations := protected.Group("/organizations")
        {
            organizations.GET("", organizationHandler.ListOrganizations)
            organizations.POST("", organizationHandler.CreateOrganization)
            organizations.GET("/user", organizationHandler.GetUserOrganizations)
            organizations.GET("/search", organizationHandler.GetOrganizationByName)
            organizations.GET("/:id", organizationHandler.GetOrganization)
            organizations.PUT("/:id", organizationHandler.UpdateOrganization)
            organizations.DELETE("/:id", organizationHandler.DeleteOrganization)
            organizations.POST("/:id/users", organizationHandler.AddUserToOrganization)
            organizations.DELETE("/:id/users/:user_id", organizationHandler.RemoveUserFromOrganization)
        }
    }

    // Create HTTPS server
    server := &http.Server{
        Addr:    ":" + cfg.Server.Port,
        Handler: r,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
            CurvePreferences: []tls.CurveID{
                tls.CurveP521,
                tls.CurveP384,
                tls.CurveP256,
            },
            PreferServerCipherSuites: true,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
            },
        },
        ReadTimeout:    15 * time.Second,
        WriteTimeout:   15 * time.Second,
        IdleTimeout:    60 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1 MB
    }

    // Start server in a goroutine
    go func() {
        log.Printf("Starting CollabHub Music server on https://localhost:%s", cfg.Server.Port)
        
        // Check if SSL certificates exist
        if _, err := os.Stat(cfg.Server.SSLCertPath); os.IsNotExist(err) {
            log.Fatal("SSL certificate file not found:", cfg.Server.SSLCertPath)
        }
        if _, err := os.Stat(cfg.Server.SSLKeyPath); os.IsNotExist(err) {
            log.Fatal("SSL key file not found:", cfg.Server.SSLKeyPath)
        }

        if err := server.ListenAndServeTLS(cfg.Server.SSLCertPath, cfg.Server.SSLKeyPath); err != nil && err != http.ErrServerClosed {
            log.Fatal("Failed to start HTTPS server:", err)
        }
    }()

    // Also start HTTP server for redirects to HTTPS (optional)
    httpServer := &http.Server{
        Addr: ":80",
        Handler: http.HandlerFunc(redirectToHTTPS),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
    }

    go func() {
        log.Printf("Starting HTTP redirect server on port 80")
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP redirect server error: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Give outstanding requests 30 seconds to complete
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Shutdown HTTPS server
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    // Shutdown HTTP redirect server
    if err := httpServer.Shutdown(ctx); err != nil {
        log.Printf("HTTP redirect server forced to shutdown: %v", err)
    }

    log.Println("Server exited")
}

// healthCheck endpoint
func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "ok",
        "service":   "CollabHub Music",
        "version":   "1.0.0",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "uptime":    time.Since(startTime).String(),
    })
}

// HTTP to HTTPS redirect handler
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
    target := "https://" + r.Host + r.URL.Path
    if len(r.URL.RawQuery) > 0 {
        target += "?" + r.URL.RawQuery
    }
    http.Redirect(w, r, target, http.StatusMovedPermanently)
}

var startTime = time.Now()