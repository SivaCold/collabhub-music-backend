# CollabHub Music Backend API

CollabHub Music is a collaborative music production platform backend built with Go and Gin framework. This API provides endpoints for user management, project collaboration, organization management, and Keycloak authentication integration.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Keycloak server
- Docker and Docker Compose (optional but recommended)

### Environment Setup

1. Copy the environment template:
```bash
cp .env.example .env
```

2. Configure your environment variables in `.env`:
```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=collabhub_music
DB_USER=your_username
DB_PASSWORD=your_password

# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8180
KEYCLOAK_REALM=collabhub-music
KEYCLOAK_CLIENT_ID=collabhub-backend
```

### Development with Air (Hot Reload)

1. Install Air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

2. Run the development server:
```bash
air
```

Or use the development scripts:
- **Linux/macOS**: `./dev.sh`
- **Windows**: `.\dev.ps1`

### Docker Development

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## 📱 React Native Integration

### CORS Configuration

The backend is configured to support React Native development with proper CORS settings:

- **Development Origins**: `http://localhost:*`, `http://127.0.0.1:*`, `http://192.168.*.*:*`
- **Expo Origins**: `http://192.168.*.*:19000`, `http://192.168.*.*:19001`, `http://192.168.*.*:19002`
- **Custom Origins**: Configure via `CORS_ALLOWED_ORIGINS` environment variable

### Authentication Flow

1. **Login Process**:
```javascript
// React Native example
const login = async (username, password) => {
  try {
    const response = await fetch('http://your-api-url/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    
    const data = await response.json();
    
    if (data.success) {
      // Store tokens securely
      await AsyncStorage.setItem('accessToken', data.data.access_token);
      await AsyncStorage.setItem('refreshToken', data.data.refresh_token);
    }
  } catch (error) {
    console.error('Login error:', error);
  }
};
```

2. **Authenticated Requests**:
```javascript
const makeAuthenticatedRequest = async (url, options = {}) => {
  const token = await AsyncStorage.getItem('accessToken');
  
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
};
```

### API Endpoints

#### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/auth/profile` - Get user profile

#### Users
- `GET /api/users` - List users (with pagination)
- `GET /api/users/{id}` - Get user by ID
- `POST /api/users` - Create user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

#### Projects
- `GET /api/projects` - List projects
- `GET /api/projects/{id}` - Get project by ID
- `POST /api/projects` - Create project
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project
- `POST /api/projects/{id}/members` - Add project member
- `DELETE /api/projects/{id}/members/{userId}` - Remove project member

#### Organizations
- `GET /api/organizations` - List organizations
- `GET /api/organizations/{id}` - Get organization by ID
- `POST /api/organizations` - Create organization
- `PUT /api/organizations/{id}` - Update organization
- `DELETE /api/organizations/{id}` - Delete organization

## 📄 API Documentation

### Swagger/OpenAPI

The API is fully documented with OpenAPI 3.0 specifications. Access the interactive documentation:

- **Swagger UI**: `http://localhost:8080/docs/index.html`
- **OpenAPI JSON**: `http://localhost:8080/docs/swagger.json`
- **OpenAPI YAML**: `http://localhost:8080/docs/swagger.yaml`

### Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  },
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

Error responses:
```json
{
  "success": false,
  "message": "Error description",
  "error": "VALIDATION_ERROR",
  "details": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ]
}
```

### Authentication Headers

All protected endpoints require authentication:
```
Authorization: Bearer <access_token>
```

## 🛠️ Development

### Building

```bash
# Build the application
make build

# Build for production
make build-prod

# Cross-compile for different platforms
make build-linux
make build-windows
make build-darwin
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

### Database Management

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Reset database
make db-reset

# Create new migration
make migrate-create NAME=your_migration_name
```

## 🐳 Docker

### Development
```bash
# Start development environment
docker-compose up -d

# View application logs
docker-compose logs -f app

# Access database
docker-compose exec db psql -U collabhub -d collabhub_music
```

### Production
```bash
# Build production image
docker build -t collabhub-backend:latest .

# Run production container
docker run -d \
  --name collabhub-backend \
  -p 8080:8080 \
  --env-file .env \
  collabhub-backend:latest
```

## 🔧 Configuration

### Environment Variables

See `.env.example` for all available configuration options:

- **Server**: Host, port, environment, SSL settings
- **Database**: Connection details, pool settings
- **Keycloak**: Authentication server configuration
- **CORS**: Cross-origin settings for React Native
- **Logging**: Log levels and output settings
- **Storage**: File upload and storage settings

### SSL/TLS Configuration

For HTTPS in production:

```env
SERVER_SSL_ENABLED=true
SERVER_SSL_CERT_FILE=./certs/server.crt
SERVER_SSL_KEY_FILE=./certs/server.key
```

Generate self-signed certificates for development:
```bash
make generate-certs
```

## 📊 Monitoring and Logging

### Health Check

```bash
curl http://localhost:8080/api/health
```

### Metrics

The application exposes metrics for monitoring:
- Request duration
- Request count
- Active connections
- Database pool statistics

### Logging

Structured logging with different levels:
- `DEBUG`: Detailed information for troubleshooting
- `INFO`: General application flow
- `WARN`: Warning messages
- `ERROR`: Error conditions
- `FATAL`: Critical errors that cause application exit

## 🚀 Deployment

### Production Checklist

1. Set environment to production: `SERVER_ENV=production`
2. Configure secure database credentials
3. Set up SSL/TLS certificates
4. Configure CORS for production domains
5. Set up monitoring and logging
6. Configure backup strategies
7. Set resource limits (memory, CPU)

### CI/CD Pipeline

The project includes GitHub Actions workflows for:
- Automated testing
- Security scanning
- Docker image building
- Deployment to staging/production

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Code Style

- Follow Go best practices and idioms
- Use `gofmt` for formatting
- Run `golint` and `go vet` before committing
- Write tests for new functionality
- Update documentation as needed

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [React Native Documentation](https://reactnative.dev/docs/getting-started)

## 📞 Support

For support and questions:
- Create an issue on GitHub
- Contact the development team
- Check the documentation and FAQ sections

---

**Happy coding! 🎵🚀**
