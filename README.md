# CollabHub Music Backend

CollabHub Music is a collaborative platform for musicians, allowing users to manage music projects, collaborate with others, and share their work. This backend service is built using Go and the Gin framework, providing a robust, secure, and efficient REST API with HTTPS support and Keycloak authentication.

## Features

### Core Features
- **User Management**: Complete user lifecycle with Keycloak integration
- **Project Management**: Create, manage, and collaborate on music projects
- **Organization Management**: Group users and projects under organizations
- **Authentication & Authorization**: JWT-based authentication with Keycloak
- **HTTPS Support**: Secure communication with TLS 1.2+ encryption
- **Database Migrations**: Automatic schema management and versioning
- **Graceful Shutdown**: Proper resource cleanup and connection management

### API Features
- RESTful API design with JSON responses
- Comprehensive error handling and validation
- CORS support for cross-origin requests
- Request logging and monitoring
- Health check endpoints
- Pagination support for list endpoints

## Architecture

### Project Structure
```
collabhub-music-backend/
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and migrations
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/          # Custom middleware (Auth, CORS, Logging)
│   ├── models/              # Data models
│   ├── repository/          # Data access layer
│   └── services/            # Business logic layer
├── docker/                  # Docker configuration files
├── docs/                    # API documentation
└── scripts/                 # Utility scripts
```

### Technology Stack
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15+
- **Authentication**: Keycloak
- **Containerization**: Docker & Docker Compose
- **Security**: HTTPS/TLS, JWT tokens

## Getting Started

### Prerequisites

- **Go**: 1.21 or later
- **PostgreSQL**: 15 or later
- **Docker & Docker Compose**: Latest versions
- **Keycloak**: 20.0 or later
- **SSL Certificates**: For HTTPS (production)

### Environment Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/collabhub-music-backend.git
   cd collabhub-music-backend
   ```

2. **Create environment files**:

   **Development (`.env.development`)**:
   ```env
   # Server Configuration
   GO_ENV=development
   SERVER_HOST=localhost
   SERVER_PORT=8443
   GIN_MODE=debug
   SSL_ENABLED=false
   SSL_CERT_PATH=
   SSL_KEY_PATH=

   # Database Configuration
   DATABASE_HOST=localhost
   DATABASE_PORT=5432
   DATABASE_USER=collabhub
   DATABASE_PASSWORD=your_password
   DATABASE_NAME=collabhub_music_dev
   DATABASE_SSL_MODE=disable
   DATABASE_TIMEZONE=UTC
   DATABASE_MAX_OPEN_CONNS=25
   DATABASE_MAX_IDLE_CONNS=5
   DATABASE_CONN_MAX_LIFETIME=300

   # Keycloak Configuration
   KEYCLOAK_URL=http://localhost:8080
   KEYCLOAK_REALM=collabhub
   KEYCLOAK_CLIENT_ID=collabhub-backend
   KEYCLOAK_CLIENT_SECRET=your_client_secret
   ```

   **Production (`.env.production`)**:
   ```env
   # Server Configuration
   GO_ENV=production
   SERVER_HOST=your-domain.com
   SERVER_PORT=443
   GIN_MODE=release
   SSL_ENABLED=true
   SSL_CERT_PATH=/path/to/cert.pem
   SSL_KEY_PATH=/path/to/key.pem

   # Database Configuration (use secure values)
   DATABASE_HOST=your-db-host
   DATABASE_PORT=5432
   DATABASE_USER=collabhub_prod
   DATABASE_PASSWORD=secure_password
   DATABASE_NAME=collabhub_music_prod
   DATABASE_SSL_MODE=require
   DATABASE_TIMEZONE=UTC
   DATABASE_MAX_OPEN_CONNS=50
   DATABASE_MAX_IDLE_CONNS=10
   DATABASE_CONN_MAX_LIFETIME=600

   # Keycloak Configuration
   KEYCLOAK_URL=https://auth.your-domain.com
   KEYCLOAK_REALM=collabhub-prod
   KEYCLOAK_CLIENT_ID=collabhub-backend
   KEYCLOAK_CLIENT_SECRET=production_secret
   ```

### Development Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Start PostgreSQL**:
   ```bash
   docker-compose -f docker/docker-compose.dev.yml up -d postgres
   ```

3. **Start Keycloak**:
   ```bash
   docker-compose -f docker/docker-compose.dev.yml up -d keycloak
   ```

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

### Production Deployment

1. **Build Docker image**:
   ```bash
   docker build -t collabhub-music-backend:latest .
   ```

2. **Deploy with Docker Compose**:
   ```bash
   docker-compose -f docker/docker-compose.prod.yml up -d
   ```

## API Documentation

### Base URL
- **Development**: `http://localhost:8443/api/v1`
- **Production**: `https://your-domain.com/api/v1`

### Authentication
All protected endpoints require a Bearer token:
```bash
Authorization: Bearer <jwt-token>
```

### Core Endpoints

#### Public Endpoints
- `GET /health` - Service health check
- `POST /users/register` - User registration

#### User Management
- `GET /users/me` - Get current user profile
- `PUT /users/me` - Update current user profile
- `DELETE /users/me` - Delete current user account
- `GET /users/{id}` - Get user profile by ID
- `GET /users` - List users (paginated)

#### Project Management
- `GET /projects` - List all projects (paginated)
- `POST /projects` - Create new project
- `GET /projects/user` - Get current user's projects
- `GET /projects/search` - Search projects by name
- `GET /projects/{id}` - Get project details
- `PUT /projects/{id}` - Update project
- `DELETE /projects/{id}` - Delete project
- `GET /projects/{id}/stats` - Get project statistics

#### Organization Management
- `GET /organizations` - List all organizations (paginated)
- `POST /organizations` - Create new organization
- `GET /organizations/user` - Get current user's organizations
- `GET /organizations/search` - Search organizations by name
- `GET /organizations/{id}` - Get organization details
- `PUT /organizations/{id}` - Update organization
- `DELETE /organizations/{id}` - Delete organization
- `POST /organizations/{id}/users` - Add user to organization
- `DELETE /organizations/{id}/users/{user_id}` - Remove user from organization

### Response Format

**Success Response**:
```json
{
  "data": {...},
  "status": "success"
}
```

**Error Response**:
```json
{
  "error": "Error message",
  "status": "error",
  "code": 400
}
```

**Paginated Response**:
```json
{
  "data": [...],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 100,
    "has_more": true
  }
}
```

## Database

### Schema Overview
- **users**: User profiles and authentication data
- **organizations**: Organization information and settings
- **projects**: Music project details and metadata
- **user_organizations**: Many-to-many relationship with roles
- **migrations**: Database schema version tracking

### Migrations
Database migrations are automatically executed on startup:
```bash
# Manual migration execution
go run cmd/migrate/main.go
```

## Security Features

### HTTPS/TLS
- **TLS 1.2+** minimum version
- **Strong cipher suites** configuration
- **HTTP to HTTPS** automatic redirect
- **HSTS headers** for security

### Authentication
- **JWT tokens** with Keycloak validation
- **Token expiration** and refresh handling
- **Role-based access** control
- **CORS protection** with configurable origins

### Database Security
- **Prepared statements** to prevent SQL injection
- **Connection pooling** with limits
- **Soft delete** for data retention
- **Audit timestamps** on all records

## Monitoring & Logging

### Health Checks
```bash
curl -k https://localhost:8443/api/v1/health
```

### Logging
- **Structured logging** with fields
- **Request/Response** logging
- **Error tracking** with stack traces
- **Performance metrics** (latency, throughput)

### Metrics
- Database connection pool statistics
- HTTP request metrics
- Error rates and response times
- System resource usage

## Development

### Code Structure
- **Clean Architecture** with separated layers
- **Interface-based** design for testability
- **Context propagation** throughout the stack
- **Error wrapping** for better debugging

### Testing
```bash
# Run unit tests
go test ./internal/...

# Run integration tests
go test -tags=integration ./tests/...

# Run with coverage
go test -cover ./internal/...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Security scan
gosec ./...
```

## Docker Support

### Development
```yaml
# docker/docker-compose.dev.yml
services:
  app:
    build: .
    environment:
      - GO_ENV=development
    volumes:
      - .:/app
    ports:
      - "8443:8443"
```

### Production
```yaml
# docker/docker-compose.prod.yml
services:
  app:
    image: collabhub-music-backend:latest
    environment:
      - GO_ENV=production
    volumes:
      - /path/to/certs:/certs:ro
```

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run tests: `go test ./...`
5. Commit changes: `git commit -m 'Add amazing feature'`
6. Push to branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Standards
- Follow Go conventions and best practices
- Write comprehensive tests for new features
- Update documentation for API changes
- Use semantic commit messages
- Ensure HTTPS/security compliance

### Issue Reporting
- Use GitHub Issues for bug reports
- Include reproduction steps and environment details
- Label issues appropriately (bug, enhancement, etc.)

## Deployment

### Prerequisites
- SSL certificates for HTTPS
- Keycloak server configured and running
- PostgreSQL database accessible
- Domain name and DNS configured

### Steps
1. Configure environment variables
2. Set up SSL certificates
3. Deploy using Docker Compose
4. Configure reverse proxy (nginx/traefik)
5. Set up monitoring and logging
6. Configure backup strategies

## Troubleshooting

### Common Issues
- **SSL Certificate errors**: Verify cert paths and permissions
- **Database connection**: Check connection string and firewall
- **Keycloak integration**: Verify client configuration and secrets
- **CORS issues**: Configure allowed origins properly

### Logs Location
- **Application logs**: `/var/log/collabhub/app.log`
- **Access logs**: `/var/log/collabhub/access.log`
- **Error logs**: `/var/log/collabhub/error.log`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **Gin Framework** - High-performance HTTP web framework
- **Keycloak** - Identity and access management solution
- **PostgreSQL** - Advanced open-source relational database
- **Go Community** - For excellent tools and libraries
- **Docker** - Containerization platform for consistent deployments

---

**Support**: For questions and support, please open an issue on GitHub or contact the development team.

**Version**: 1.0.0  
**Last Updated**: December 2024