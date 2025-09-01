#!/bin/bash

# ===========================================
# CollabHub Music Backend Setup Script
# ===========================================

set -e

echo "🎵 Setting up CollabHub Music Backend..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        echo "Visit https://golang.org/doc/install for installation instructions."
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d ' ' -f 3 | sed 's/go//')
    print_success "Go version $GO_VERSION found"
}

# Check if Docker is installed
check_docker() {
    if command -v docker &> /dev/null; then
        print_success "Docker found"
        DOCKER_VERSION=$(docker --version | cut -d ' ' -f 3 | sed 's/,//')
        print_status "Docker version: $DOCKER_VERSION"
    else
        print_warning "Docker not found. Docker is optional but recommended for development."
    fi
    
    if command -v docker-compose &> /dev/null; then
        print_success "Docker Compose found"
    elif command -v docker compose &> /dev/null; then
        print_success "Docker Compose (v2) found"
    else
        print_warning "Docker Compose not found."
    fi
}

# Initialize Go module if not exists
init_go_module() {
    if [ ! -f "go.mod" ]; then
        print_status "Initializing Go module..."
        go mod init collabhub-music-backend
        print_success "Go module initialized"
    else
        print_status "Go module already exists"
    fi
}

# Install Go dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    
    # Core dependencies
    go get github.com/gin-gonic/gin@latest
    go get github.com/joho/godotenv@latest
    go get gorm.io/gorm@latest
    go get gorm.io/driver/postgres@latest
    
    # Authentication and JWT
    go get github.com/golang-jwt/jwt/v5@latest
    go get github.com/Nerzal/gocloak/v13@latest
    
    # Validation
    go get github.com/go-playground/validator/v10@latest
    
    # Logging
    go get github.com/sirupsen/logrus@latest
    
    # Testing
    go get github.com/stretchr/testify@latest
    
    # API documentation
    go get github.com/swaggo/swag/cmd/swag@latest
    go get github.com/swaggo/gin-swagger@latest
    go get github.com/swaggo/files@latest
    
    # Utilities
    go get github.com/google/uuid@latest
    go get golang.org/x/crypto@latest
    
    # Database migrations (optional)
    go get -u github.com/golang-migrate/migrate/v4@latest
    go get -u github.com/golang-migrate/migrate/v4/database/postgres@latest
    go get -u github.com/golang-migrate/migrate/v4/source/file@latest
    
    print_success "Dependencies installed successfully"
}

# Install development tools
install_dev_tools() {
    print_status "Installing development tools..."
    
    # Hot reload tool
    go install github.com/cosmtrek/air@latest
    
    # Swagger documentation generator
    go install github.com/swaggo/swag/cmd/swag@latest
    
    # Code formatting and linting
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # Security scanning
    go install github.com/securecodewarrior/sast-scan@latest
    
    print_success "Development tools installed"
}

# Update go.mod and go.sum
update_dependencies() {
    print_status "Updating dependencies..."
    go mod tidy
    go mod download
    print_success "Dependencies updated"
}

# Create environment file
create_env_file() {
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            print_status "Creating .env file from template..."
            cp .env.example .env
            print_success ".env file created"
            print_warning "Please update .env file with your configuration"
        else
            print_warning ".env.example not found. Creating basic .env file..."
            cat > .env << EOF
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=collabhub_music
DB_USER=collabhub
DB_PASSWORD=your_password_here

# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8180
KEYCLOAK_REALM=collabhub-music
KEYCLOAK_CLIENT_ID=collabhub-backend
KEYCLOAK_CLIENT_SECRET=your_client_secret_here
EOF
            print_success "Basic .env file created"
            print_warning "Please update .env file with your actual configuration"
        fi
    else
        print_status ".env file already exists"
    fi
}

# Generate SSL certificates for development
generate_certs() {
    if [ ! -d "certs" ]; then
        mkdir -p certs
    fi
    
    if [ ! -f "certs/server.crt" ] || [ ! -f "certs/server.key" ]; then
        print_status "Generating SSL certificates for development..."
        
        # Generate private key
        openssl genrsa -out certs/server.key 2048
        
        # Generate certificate
        openssl req -new -x509 -key certs/server.key -out certs/server.crt -days 365 -subj "/CN=localhost"
        
        print_success "SSL certificates generated"
    else
        print_status "SSL certificates already exist"
    fi
}

# Create necessary directories
create_directories() {
    print_status "Creating project directories..."
    
    directories=(
        "logs"
        "storage"
        "tmp"
        "docs"
        "scripts"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            print_status "Created directory: $dir"
        fi
    done
    
    # Create .gitkeep files to ensure empty directories are tracked
    for dir in "${directories[@]}"; do
        if [ ! -f "$dir/.gitkeep" ]; then
            touch "$dir/.gitkeep"
        fi
    done
    
    print_success "Project directories created"
}

# Generate API documentation
generate_docs() {
    if command -v swag &> /dev/null; then
        print_status "Generating API documentation..."
        swag init -g cmd/server/main.go -o docs/
        print_success "API documentation generated"
    else
        print_warning "Swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
    fi
}

# Run tests to verify setup
run_tests() {
    print_status "Running tests to verify setup..."
    
    if go test ./... -v; then
        print_success "Tests passed successfully"
    else
        print_warning "Some tests failed. This might be expected if database is not configured."
    fi
}

# Build the application
build_app() {
    print_status "Building application..."
    
    if go build -o bin/collabhub-music-backend cmd/server/main.go; then
        print_success "Application built successfully"
        print_status "Binary created: bin/collabhub-music-backend"
    else
        print_error "Build failed"
        exit 1
    fi
}

# Main setup function
main() {
    echo "🚀 Starting CollabHub Music Backend setup..."
    echo
    
    check_go
    check_docker
    init_go_module
    install_dependencies
    install_dev_tools
    update_dependencies
    create_env_file
    create_directories
    generate_certs
    generate_docs
    build_app
    
    echo
    print_success "🎉 Setup completed successfully!"
    echo
    echo "Next steps:"
    echo "1. Update .env file with your configuration"
    echo "2. Set up your PostgreSQL database"
    echo "3. Set up your Keycloak server"
    echo "4. Run 'air' or './dev.sh' to start development server"
    echo "5. Visit http://localhost:8080/docs/index.html for API documentation"
    echo
    echo "Available commands:"
    echo "  make help          - Show all available commands"
    echo "  air               - Start development server with hot reload"
    echo "  ./dev.sh          - Run development environment (Linux/macOS)"
    echo "  ./dev.ps1         - Run development environment (Windows)"
    echo "  docker-compose up - Start with Docker"
    echo
}

# Run main function
main "$@"
