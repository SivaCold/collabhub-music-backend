#!/bin/bash

# ===========================================
# CollabHub Music Backend Development Script
# ===========================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="CollabHub Music Backend"
DOCKER_COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env"
ENV_EXAMPLE_FILE=".env.example"

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    local missing_deps=()
    
    if ! command_exists docker; then
        missing_deps+=("docker")
    fi
    
    if ! command_exists docker-compose; then
        missing_deps+=("docker-compose")
    fi
    
    if ! command_exists go; then
        missing_deps+=("go")
    fi
    
    if ! command_exists openssl; then
        missing_deps+=("openssl")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        print_error "Please install the missing dependencies and try again."
        exit 1
    fi
    
    print_success "All prerequisites are installed!"
}

# Function to setup environment
setup_environment() {
    print_status "Setting up environment..."
    
    # Copy .env.example to .env if it doesn't exist
    if [ ! -f "$ENV_FILE" ]; then
        if [ -f "$ENV_EXAMPLE_FILE" ]; then
            cp "$ENV_EXAMPLE_FILE" "$ENV_FILE"
            print_success "Created $ENV_FILE from $ENV_EXAMPLE_FILE"
        else
            print_error "$ENV_EXAMPLE_FILE not found!"
            exit 1
        fi
    else
        print_warning "$ENV_FILE already exists, skipping copy"
    fi
    
    # Create necessary directories
    mkdir -p certs storage logs
    mkdir -p storage/{audio,images,temp,backups}
    
    print_success "Environment setup completed!"
}

# Function to generate TLS certificates
generate_certificates() {
    print_status "Generating TLS certificates..."
    
    if [ -f "certs/server.crt" ] && [ -f "certs/server.key" ]; then
        print_warning "Certificates already exist, skipping generation"
        return
    fi
    
    # Generate self-signed certificate for development
    openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \
        -out certs/server.crt -days 365 -nodes \
        -subj "/C=US/ST=State/L=City/O=CollabHub/CN=localhost" \
        -addext "subjectAltName=DNS:localhost,IP:127.0.0.1,IP:0.0.0.0"
    
    # Set proper permissions
    chmod 600 certs/server.key
    chmod 644 certs/server.crt
    
    print_success "TLS certificates generated successfully!"
}

# Function to build the application
build_app() {
    print_status "Building the application..."
    
    go mod tidy
    go build -o collabhub-backend cmd/server/main.go
    
    print_success "Application built successfully!"
}

# Function to run database migrations
run_migrations() {
    print_status "Running database migrations..."
    
    # Wait for database to be ready
    print_status "Waiting for database to be ready..."
    until docker-compose exec postgres pg_isready -U collabhub_user -d collabhub_music; do
        sleep 2
    done
    
    # Run migrations (this will be handled by the application)
    print_success "Database is ready for migrations!"
}

# Function to setup Keycloak
setup_keycloak() {
    print_status "Setting up Keycloak..."
    
    # Wait for Keycloak to be ready
    print_status "Waiting for Keycloak to be ready..."
    timeout=300 # 5 minutes timeout
    elapsed=0
    
    while ! curl -s http://localhost:8080/realms/master >/dev/null; do
        if [ $elapsed -ge $timeout ]; then
            print_error "Timeout waiting for Keycloak to start"
            exit 1
        fi
        sleep 5
        elapsed=$((elapsed + 5))
        print_status "Still waiting for Keycloak... ($elapsed/$timeout seconds)"
    done
    
    print_success "Keycloak is ready!"
    print_status "Please configure Keycloak realm and client manually at: http://localhost:8080/admin"
    print_status "Default admin credentials: admin / admin123"
}

# Function to start services
start_services() {
    print_status "Starting services with Docker Compose..."
    
    # Start infrastructure services first
    docker-compose up -d postgres redis
    
    # Wait a bit for services to initialize
    sleep 10
    
    # Start Keycloak
    docker-compose up -d keycloak
    
    # Setup Keycloak
    setup_keycloak
    
    # Start the backend
    docker-compose up -d backend
    
    print_success "All services started successfully!"
    print_status "Services status:"
    docker-compose ps
    
    echo ""
    print_success "🚀 $PROJECT_NAME is now running!"
    echo ""
    echo "📋 Service URLs:"
    echo "   • Backend API: https://localhost:8443"
    echo "   • API Documentation: https://localhost:8443/swagger/index.html"
    echo "   • Health Check: https://localhost:8443/health"
    echo "   • Keycloak Admin: http://localhost:8080/admin (admin/admin123)"
    echo "   • Database: localhost:5432 (collabhub_user/collabhub_password123)"
    echo ""
    echo "📱 React Native Development:"
    echo "   • iOS Simulator: https://localhost:8443/api/v1"
    echo "   • Android Emulator: https://10.0.2.2:8443/api/v1"
    echo "   • Physical Device: https://YOUR_IP:8443/api/v1"
    echo ""
}

# Function to stop services
stop_services() {
    print_status "Stopping services..."
    
    docker-compose down
    
    print_success "All services stopped!"
}

# Function to view logs
view_logs() {
    local service=${1:-""}
    
    if [ -n "$service" ]; then
        print_status "Viewing logs for $service..."
        docker-compose logs -f "$service"
    else
        print_status "Viewing logs for all services..."
        docker-compose logs -f
    fi
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    
    go test ./...
    
    print_success "All tests completed!"
}

# Function to clean up
cleanup() {
    print_status "Cleaning up..."
    
    # Stop and remove containers
    docker-compose down -v
    
    # Remove built binary
    [ -f "collabhub-backend" ] && rm collabhub-backend
    
    # Clean Go module cache
    go clean -modcache
    
    print_success "Cleanup completed!"
}

# Function to show help
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  setup       Setup environment and generate certificates"
    echo "  start       Start all services"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  logs        View logs for all services"
    echo "  logs <svc>  View logs for specific service"
    echo "  build       Build the application"
    echo "  test        Run tests"
    echo "  clean       Clean up containers and volumes"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 setup           # First time setup"
    echo "  $0 start           # Start all services"
    echo "  $0 logs backend    # View backend logs"
    echo "  $0 clean           # Clean everything"
}

# Main script logic
main() {
    local command=${1:-"help"}
    
    echo ""
    print_status "🎵 $PROJECT_NAME Development Script"
    echo ""
    
    case "$command" in
        "setup")
            check_prerequisites
            setup_environment
            generate_certificates
            print_success "Setup completed! Run '$0 start' to start services."
            ;;
        "start")
            check_prerequisites
            setup_environment
            generate_certificates
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            sleep 2
            start_services
            ;;
        "logs")
            view_logs "$2"
            ;;
        "build")
            build_app
            ;;
        "test")
            run_tests
            ;;
        "clean")
            cleanup
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# Run main function with all arguments
main "$@"
