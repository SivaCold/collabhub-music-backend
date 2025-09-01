// Package main CollabHub Music Backend API
//
// This is the API documentation for CollabHub Music Backend, a collaborative music platform.
// The API provides endpoints for user management, project collaboration, organization management,
// and authentication using Keycloak.
//
// Terms Of Service: https://collabhub-music.com/terms
// Contact: support@collabhub-music.com
// Version: 1.0.0
//
//	@title			CollabHub Music Backend API
//	@description	A collaborative music platform API for managing users, projects, and organizations
//	@version		1.0.0
//	@host			localhost:8443
//	@BasePath		/
//	@schemes		https http
//
//	@contact.name	CollabHub Music Support
//	@contact.url	https://collabhub-music.com/support
//	@contact.email	support@collabhub-music.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Authorization header using the Bearer scheme. Enter 'Bearer' [space] and then your token in the text input below. Example: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
//
//	@tag.name			authentication
//	@tag.description	Authentication and authorization endpoints
//
//	@tag.name			users
//	@tag.description	User management operations
//
//	@tag.name			projects
//	@tag.description	Music project management operations
//
//	@tag.name			organizations
//	@tag.description	Organization management operations
package docs
