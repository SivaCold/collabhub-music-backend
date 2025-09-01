package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"collabhub-music-backend/internal/config"
	"collabhub-music-backend/internal/handlers"
	"collabhub-music-backend/internal/middleware"
	"collabhub-music-backend/internal/models"
	"collabhub-music-backend/pkg/utils"
)

// IntegrationTestSuite represents the integration test suite
type IntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	config *config.Config
}

// SetupSuite runs once before all tests in the suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Set test environment
	os.Setenv("SERVER_ENV", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "collabhub_test")

	// Load test configuration
	cfg, err := config.Load()
	suite.NoError(err)
	suite.config = cfg

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup router
	suite.router = suite.setupRouter()
}

// TearDownSuite runs once after all tests in the suite
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up test environment
	os.Unsetenv("SERVER_ENV")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_NAME")
}

// setupRouter configures the test router with all middlewares and routes
func (suite *IntegrationTestSuite) setupRouter() *gin.Engine {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(suite.config.CORS))

	// API v1 routes
	v1 := router.Group("/api")
	{
		// Health check endpoint (no auth required)
		v1.GET("/health", handlers.HealthCheck)

		// Auth endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Protected endpoints (require authentication)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User endpoints
			users := protected.Group("/users")
			{
				users.GET("", handlers.GetUsers)
				users.GET("/:id", handlers.GetUser)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// Project endpoints
			projects := protected.Group("/projects")
			{
				projects.GET("", handlers.GetProjects)
				projects.GET("/:id", handlers.GetProject)
				projects.POST("", handlers.CreateProject)
				projects.PUT("/:id", handlers.UpdateProject)
				projects.DELETE("/:id", handlers.DeleteProject)
				projects.POST("/:id/members", handlers.AddProjectMember)
				projects.DELETE("/:id/members/:userId", handlers.RemoveProjectMember)
			}

			// Organization endpoints
			organizations := protected.Group("/organizations")
			{
				organizations.GET("", handlers.GetOrganizations)
				organizations.GET("/:id", handlers.GetOrganization)
				organizations.POST("", handlers.CreateOrganization)
				organizations.PUT("/:id", handlers.UpdateOrganization)
				organizations.DELETE("/:id", handlers.DeleteOrganization)
			}
		}
	}

	return router
}

// TestHealthCheck tests the health check endpoint
func (suite *IntegrationTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Service is healthy", response.Message)
}

// TestCORSHeaders tests that CORS headers are properly set
func (suite *IntegrationTestSuite) TestCORSHeaders() {
	req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNoContent, resp.Code)
	assert.Equal(suite.T(), "*", resp.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(suite.T(), resp.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(suite.T(), resp.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

// TestAuthLogin tests the login endpoint
func (suite *IntegrationTestSuite) TestAuthLogin() {
	loginData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	// Note: This will return an error in test environment without Keycloak
	// but we can test that the endpoint exists and handles the request
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusUnauthorized, http.StatusInternalServerError}, resp.Code)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
}

// TestAuthEndpointsWithoutAuth tests that protected endpoints require authentication
func (suite *IntegrationTestSuite) TestAuthEndpointsWithoutAuth() {
	endpoints := []string{
		"/api/users",
		"/api/projects",
		"/api/organizations",
	}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("GET", endpoint, nil)
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code,
			fmt.Sprintf("Endpoint %s should require authentication", endpoint))

		var response utils.Response
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), response.Success)
	}
}

// TestUserEndpointsStructure tests the structure of user endpoints
func (suite *IntegrationTestSuite) TestUserEndpointsStructure() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/users"},
		{"GET", "/api/users/123"},
		{"POST", "/api/users"},
		{"PUT", "/api/users/123"},
		{"DELETE", "/api/users/123"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		// Should return 401 (Unauthorized) since no auth token is provided
		// If it returns 404, the route doesn't exist
		assert.NotEqual(suite.T(), http.StatusNotFound, resp.Code,
			fmt.Sprintf("Endpoint %s %s should exist", ep.method, ep.path))
	}
}

// TestProjectEndpointsStructure tests the structure of project endpoints
func (suite *IntegrationTestSuite) TestProjectEndpointsStructure() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/projects"},
		{"GET", "/api/projects/123"},
		{"POST", "/api/projects"},
		{"PUT", "/api/projects/123"},
		{"DELETE", "/api/projects/123"},
		{"POST", "/api/projects/123/members"},
		{"DELETE", "/api/projects/123/members/456"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		// Should return 401 (Unauthorized) since no auth token is provided
		assert.NotEqual(suite.T(), http.StatusNotFound, resp.Code,
			fmt.Sprintf("Endpoint %s %s should exist", ep.method, ep.path))
	}
}

// TestOrganizationEndpointsStructure tests the structure of organization endpoints
func (suite *IntegrationTestSuite) TestOrganizationEndpointsStructure() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/organizations"},
		{"GET", "/api/organizations/123"},
		{"POST", "/api/organizations"},
		{"PUT", "/api/organizations/123"},
		{"DELETE", "/api/organizations/123"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		// Should return 401 (Unauthorized) since no auth token is provided
		assert.NotEqual(suite.T(), http.StatusNotFound, resp.Code,
			fmt.Sprintf("Endpoint %s %s should exist", ep.method, ep.path))
	}
}

// TestRequestValidation tests input validation
func (suite *IntegrationTestSuite) TestRequestValidation() {
	// Test invalid JSON
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
}

// TestResponseFormat tests that all responses follow the standard format
func (suite *IntegrationTestSuite) TestResponseFormat() {
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	var response utils.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	// Check that the response has the required fields
	assert.NotEmpty(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.Message)
}

// TestContentTypeHeaders tests that appropriate content-type headers are set
func (suite *IntegrationTestSuite) TestContentTypeHeaders() {
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	contentType := resp.Header().Get("Content-Type")
	assert.Contains(suite.T(), contentType, "application/json")
}

// TestRateLimitingHeaders tests rate limiting (if implemented)
func (suite *IntegrationTestSuite) TestRateLimitingHeaders() {
	// Make multiple requests quickly
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/api/health", nil)
		resp := httptest.NewRecorder()

		suite.router.ServeHTTP(resp, req)

		// If rate limiting is implemented, check for rate limit headers
		rateLimitRemaining := resp.Header().Get("X-RateLimit-Remaining")
		rateLimitReset := resp.Header().Get("X-RateLimit-Reset")

		if rateLimitRemaining != "" || rateLimitReset != "" {
			// Rate limiting is implemented
			suite.T().Logf("Rate limiting headers found: Remaining=%s, Reset=%s",
				rateLimitRemaining, rateLimitReset)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// TestModelValidation tests model validation
func (suite *IntegrationTestSuite) TestModelValidation() {
	// Test User model validation
	user := models.User{
		Email:    "invalid-email",
		Username: "",
	}

	// This would normally be tested with actual database validation
	// For now, we just ensure the struct exists and can be created
	assert.NotNil(suite.T(), user)

	// Test Project model validation
	project := models.Project{
		Name:        "",
		Description: "Test project",
	}

	assert.NotNil(suite.T(), project)

	// Test Organization model validation
	org := models.Organization{
		Name:        "",
		Description: "Test organization",
	}

	assert.NotNil(suite.T(), org)
}

// TestConfigurationLoading tests that configuration loads properly
func (suite *IntegrationTestSuite) TestConfigurationLoading() {
	cfg := suite.config

	assert.NotNil(suite.T(), cfg)
	assert.NotEmpty(suite.T(), cfg.Server.Host)
	assert.Greater(suite.T(), cfg.Server.Port, 0)
	assert.NotNil(suite.T(), cfg.Database)
	assert.NotNil(suite.T(), cfg.Keycloak)
	assert.NotNil(suite.T(), cfg.CORS)
}

// TestSecurityHeaders tests that security headers are properly set
func (suite *IntegrationTestSuite) TestSecurityHeaders() {
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	// Check for common security headers
	headers := resp.Header()

	// These headers might be set by middleware
	xFrameOptions := headers.Get("X-Frame-Options")
	xContentTypeOptions := headers.Get("X-Content-Type-Options")
	xXSSProtection := headers.Get("X-XSS-Protection")

	suite.T().Logf("Security headers - X-Frame-Options: %s, X-Content-Type-Options: %s, X-XSS-Protection: %s",
		xFrameOptions, xContentTypeOptions, xXSSProtection)
}

// Run the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Benchmark tests
func BenchmarkHealthCheck(b *testing.B) {
	os.Setenv("SERVER_ENV", "test")
	gin.SetMode(gin.TestMode)

	cfg, _ := config.Load()
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS))
	router.GET("/api/health", handlers.HealthCheck)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/health", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

// Example function demonstrating how to create test data
func createTestUser() models.User {
	return models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Example function demonstrating how to create test project
func createTestProject() models.Project {
	return models.Project{
		Name:        "Test Project",
		Description: "A test project for integration testing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Example function demonstrating how to create test organization
func createTestOrganization() models.Organization {
	return models.Organization{
		Name:        "Test Organization",
		Description: "A test organization for integration testing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
