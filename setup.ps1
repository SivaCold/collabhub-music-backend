# ===========================================
# CollabHub Music Backend Setup Script (PowerShell)
# ===========================================

param(
    [switch]$SkipTests,
    [switch]$SkipBuild,
    [switch]$Help
)

# Function to show help
function Show-Help {
    Write-Host "🎵 CollabHub Music Backend Setup Script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\setup.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -SkipTests    Skip running tests during setup"
    Write-Host "  -SkipBuild    Skip building the application"
    Write-Host "  -Help         Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\setup.ps1                 # Full setup"
    Write-Host "  .\setup.ps1 -SkipTests      # Setup without tests"
    Write-Host "  .\setup.ps1 -Help           # Show help"
}

if ($Help) {
    Show-Help
    exit 0
}

# Colors for output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Check if Go is installed
function Test-Go {
    try {
        $goVersion = go version
        if ($goVersion) {
            $version = ($goVersion -split " ")[2] -replace "go", ""
            Write-Success "Go version $version found"
            return $true
        }
    }
    catch {
        Write-Error "Go is not installed. Please install Go 1.21 or higher."
        Write-Host "Visit https://golang.org/doc/install for installation instructions."
        return $false
    }
}

# Check if Docker is installed
function Test-Docker {
    try {
        $dockerVersion = docker --version
        if ($dockerVersion) {
            Write-Success "Docker found"
            Write-Status "Docker version: $($dockerVersion -split " ")[2] -replace ",", "")"
        }
    }
    catch {
        Write-Warning "Docker not found. Docker is optional but recommended for development."
    }
    
    try {
        $composeVersion = docker-compose --version
        if ($composeVersion) {
            Write-Success "Docker Compose found"
        }
    }
    catch {
        try {
            docker compose version | Out-Null
            Write-Success "Docker Compose (v2) found"
        }
        catch {
            Write-Warning "Docker Compose not found."
        }
    }
}

# Initialize Go module if not exists
function Initialize-GoModule {
    if (-not (Test-Path "go.mod")) {
        Write-Status "Initializing Go module..."
        go mod init collabhub-music-backend
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Go module initialized"
        } else {
            Write-Error "Failed to initialize Go module"
            exit 1
        }
    } else {
        Write-Status "Go module already exists"
    }
}

# Install Go dependencies
function Install-Dependencies {
    Write-Status "Installing Go dependencies..."
    
    $dependencies = @(
        "github.com/gin-gonic/gin@latest",
        "github.com/joho/godotenv@latest",
        "gorm.io/gorm@latest",
        "gorm.io/driver/postgres@latest",
        "github.com/golang-jwt/jwt/v5@latest",
        "github.com/Nerzal/gocloak/v13@latest",
        "github.com/go-playground/validator/v10@latest",
        "github.com/sirupsen/logrus@latest",
        "github.com/stretchr/testify@latest",
        "github.com/swaggo/swag/cmd/swag@latest",
        "github.com/swaggo/gin-swagger@latest",
        "github.com/swaggo/files@latest",
        "github.com/google/uuid@latest",
        "golang.org/x/crypto@latest",
        "github.com/golang-migrate/migrate/v4@latest",
        "github.com/golang-migrate/migrate/v4/database/postgres@latest",
        "github.com/golang-migrate/migrate/v4/source/file@latest"
    )
    
    foreach ($dep in $dependencies) {
        Write-Status "Installing $dep..."
        go get $dep
        if ($LASTEXITCODE -ne 0) {
            Write-Warning "Failed to install $dep"
        }
    }
    
    Write-Success "Dependencies installed successfully"
}

# Install development tools
function Install-DevTools {
    Write-Status "Installing development tools..."
    
    $tools = @(
        "github.com/cosmtrek/air@latest",
        "github.com/swaggo/swag/cmd/swag@latest",
        "golang.org/x/tools/cmd/goimports@latest",
        "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    )
    
    foreach ($tool in $tools) {
        Write-Status "Installing $tool..."
        go install $tool
        if ($LASTEXITCODE -ne 0) {
            Write-Warning "Failed to install $tool"
        }
    }
    
    Write-Success "Development tools installed"
}

# Update dependencies
function Update-Dependencies {
    Write-Status "Updating dependencies..."
    go mod tidy
    if ($LASTEXITCODE -eq 0) {
        go mod download
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Dependencies updated"
        } else {
            Write-Error "Failed to download dependencies"
        }
    } else {
        Write-Error "Failed to tidy dependencies"
    }
}

# Create environment file
function New-EnvFile {
    if (-not (Test-Path ".env")) {
        if (Test-Path ".env.example") {
            Write-Status "Creating .env file from template..."
            Copy-Item ".env.example" ".env"
            Write-Success ".env file created"
            Write-Warning "Please update .env file with your configuration"
        } else {
            Write-Warning ".env.example not found. Creating basic .env file..."
            $envContent = @"
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
"@
            $envContent | Out-File -FilePath ".env" -Encoding UTF8
            Write-Success "Basic .env file created"
            Write-Warning "Please update .env file with your actual configuration"
        }
    } else {
        Write-Status ".env file already exists"
    }
}

# Create necessary directories
function New-Directories {
    Write-Status "Creating project directories..."
    
    $directories = @("logs", "storage", "tmp", "docs", "scripts", "bin")
    
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
            Write-Status "Created directory: $dir"
        }
        
        # Create .gitkeep files
        $gitkeepPath = Join-Path $dir ".gitkeep"
        if (-not (Test-Path $gitkeepPath)) {
            New-Item -ItemType File -Path $gitkeepPath -Force | Out-Null
        }
    }
    
    Write-Success "Project directories created"
}

# Generate SSL certificates for development
function New-Certificates {
    if (-not (Test-Path "certs")) {
        New-Item -ItemType Directory -Path "certs" -Force | Out-Null
    }
    
    if (-not (Test-Path "certs\server.crt") -or -not (Test-Path "certs\server.key")) {
        Write-Status "Generating SSL certificates for development..."
        
        try {
            # Check if OpenSSL is available
            openssl version | Out-Null
            
            # Generate private key
            openssl genrsa -out certs\server.key 2048
            
            # Generate certificate
            openssl req -new -x509 -key certs\server.key -out certs\server.crt -days 365 -subj "/CN=localhost"
            
            Write-Success "SSL certificates generated using OpenSSL"
        }
        catch {
            Write-Warning "OpenSSL not found. Creating placeholder certificate files."
            Write-Warning "For production, please generate proper SSL certificates."
            
            # Create placeholder files
            "# Placeholder certificate file" | Out-File -FilePath "certs\server.crt" -Encoding UTF8
            "# Placeholder key file" | Out-File -FilePath "certs\server.key" -Encoding UTF8
            
            Write-Status "Placeholder certificate files created"
        }
    } else {
        Write-Status "SSL certificates already exist"
    }
}

# Generate API documentation
function New-Documentation {
    try {
        swag version | Out-Null
        Write-Status "Generating API documentation..."
        swag init -g cmd/server/main.go -o docs/
        if ($LASTEXITCODE -eq 0) {
            Write-Success "API documentation generated"
        } else {
            Write-Warning "Failed to generate API documentation"
        }
    }
    catch {
        Write-Warning "Swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
    }
}

# Run tests to verify setup
function Invoke-Tests {
    if (-not $SkipTests) {
        Write-Status "Running tests to verify setup..."
        go test ./... -v
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Tests passed successfully"
        } else {
            Write-Warning "Some tests failed. This might be expected if database is not configured."
        }
    } else {
        Write-Status "Skipping tests as requested"
    }
}

# Build the application
function Build-Application {
    if (-not $SkipBuild) {
        Write-Status "Building application..."
        
        if (-not (Test-Path "bin")) {
            New-Item -ItemType Directory -Path "bin" -Force | Out-Null
        }
        
        go build -o bin\collabhub-music-backend.exe cmd\server\main.go
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Application built successfully"
            Write-Status "Binary created: bin\collabhub-music-backend.exe"
        } else {
            Write-Error "Build failed"
            exit 1
        }
    } else {
        Write-Status "Skipping build as requested"
    }
}

# Main setup function
function main {
    Write-Host "🚀 Starting CollabHub Music Backend setup..." -ForegroundColor Cyan
    Write-Host ""
    
    if (-not (Test-Go)) {
        exit 1
    }
    
    Test-Docker
    Initialize-GoModule
    Install-Dependencies
    Install-DevTools
    Update-Dependencies
    New-EnvFile
    New-Directories
    New-Certificates
    New-Documentation
    Invoke-Tests
    Build-Application
    
    Write-Host ""
    Write-Success "🎉 Setup completed successfully!"
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "1. Update .env file with your configuration"
    Write-Host "2. Set up your PostgreSQL database"
    Write-Host "3. Set up your Keycloak server"
    Write-Host "4. Run 'air' or '.\dev.ps1' to start development server"
    Write-Host "5. Visit http://localhost:8080/docs/index.html for API documentation"
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Yellow
    Write-Host "  .\dev.ps1         - Run development environment"
    Write-Host "  air               - Start development server with hot reload"
    Write-Host "  docker-compose up - Start with Docker"
    Write-Host "  go run cmd\server\main.go - Run directly with Go"
    Write-Host ""
}

# Error handling
trap {
    Write-Error "An error occurred during setup: $($_.Exception.Message)"
    Write-Host "Please check the error message above and try again."
    exit 1
}

# Run main function
main
