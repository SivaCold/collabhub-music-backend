# ===========================================
# CollabHub Music Backend Development Script (PowerShell)
# ===========================================

param(
    [Parameter(Position=0)]
    [ValidateSet('setup', 'start', 'stop', 'restart', 'logs', 'build', 'test', 'clean', 'help')]
    [string]$Command = 'help',
    
    [Parameter(Position=1)]
    [string]$Service = ''
)

# Configuration
$ProjectName = "CollabHub Music Backend"
$DockerComposeFile = "docker-compose.yml"
$EnvFile = ".env"
$EnvExampleFile = ".env.example"

# Function to print colored output
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

# Function to check if command exists
function Test-Command {
    param([string]$CommandName)
    return $null -ne (Get-Command $CommandName -ErrorAction SilentlyContinue)
}

# Function to check prerequisites
function Test-Prerequisites {
    Write-Status "Checking prerequisites..."
    
    $missingDeps = @()
    
    if (-not (Test-Command "docker")) {
        $missingDeps += "docker"
    }
    
    if (-not (Test-Command "docker-compose")) {
        $missingDeps += "docker-compose"
    }
    
    if (-not (Test-Command "go")) {
        $missingDeps += "go"
    }
    
    if ($missingDeps.Count -gt 0) {
        Write-Error "Missing dependencies: $($missingDeps -join ', ')"
        Write-Error "Please install the missing dependencies and try again."
        exit 1
    }
    
    Write-Success "All prerequisites are installed!"
}

# Function to setup environment
function Initialize-Environment {
    Write-Status "Setting up environment..."
    
    # Copy .env.example to .env if it doesn't exist
    if (-not (Test-Path $EnvFile)) {
        if (Test-Path $EnvExampleFile) {
            Copy-Item $EnvExampleFile $EnvFile
            Write-Success "Created $EnvFile from $EnvExampleFile"
        } else {
            Write-Error "$EnvExampleFile not found!"
            exit 1
        }
    } else {
        Write-Warning "$EnvFile already exists, skipping copy"
    }
    
    # Create necessary directories
    $directories = @("certs", "storage", "logs", "storage\audio", "storage\images", "storage\temp", "storage\backups")
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
    }
    
    Write-Success "Environment setup completed!"
}

# Function to generate TLS certificates
function New-Certificates {
    Write-Status "Generating TLS certificates..."
    
    if ((Test-Path "certs\server.crt") -and (Test-Path "certs\server.key")) {
        Write-Warning "Certificates already exist, skipping generation"
        return
    }
    
    # Check if OpenSSL is available
    if (Test-Command "openssl") {
        # Generate self-signed certificate for development using OpenSSL
        $opensslCmd = @(
            "req", "-x509", "-newkey", "rsa:4096", "-keyout", "certs\server.key",
            "-out", "certs\server.crt", "-days", "365", "-nodes",
            "-subj", "/C=US/ST=State/L=City/O=CollabHub/CN=localhost"
        )
        & openssl @opensslCmd
    } else {
        # Use PowerShell to generate certificate (Windows only)
        Write-Status "OpenSSL not found, using PowerShell certificate generation..."
        
        $cert = New-SelfSignedCertificate -DnsName "localhost", "127.0.0.1" -CertStoreLocation "cert:\CurrentUser\My" -KeyExportPolicy Exportable -KeySpec Signature -KeyLength 2048 -KeyAlgorithm RSA -HashAlgorithm SHA256
        
        # Export certificate
        Export-Certificate -Cert $cert -FilePath "certs\server.crt"
        
        # Export private key
        $password = ConvertTo-SecureString -String "temp" -Force -AsPlainText
        Export-PfxCertificate -Cert $cert -FilePath "certs\server.pfx" -Password $password
        
        # Convert PFX to PEM format (requires OpenSSL)
        if (Test-Command "openssl") {
            & openssl pkcs12 -in "certs\server.pfx" -out "certs\server.key" -nodes -nocerts -passin pass:temp
            Remove-Item "certs\server.pfx"
        } else {
            Write-Warning "OpenSSL not available. PFX certificate created instead of PEM format."
        }
        
        # Remove certificate from store
        Remove-Item -Path "cert:\CurrentUser\My\$($cert.Thumbprint)"
    }
    
    Write-Success "TLS certificates generated successfully!"
}

# Function to build the application
function Build-App {
    Write-Status "Building the application..."
    
    & go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to tidy Go modules"
        exit 1
    }
    
    & go build -o collabhub-backend.exe cmd\server\main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build application"
        exit 1
    }
    
    Write-Success "Application built successfully!"
}

# Function to start services
function Start-Services {
    Write-Status "Starting services with Docker Compose..."
    
    # Start infrastructure services first
    & docker-compose up -d postgres redis
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to start infrastructure services"
        exit 1
    }
    
    # Wait a bit for services to initialize
    Write-Status "Waiting for infrastructure services to initialize..."
    Start-Sleep -Seconds 10
    
    # Start Keycloak
    & docker-compose up -d keycloak
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to start Keycloak"
        exit 1
    }
    
    # Wait for Keycloak to be ready
    Write-Status "Waiting for Keycloak to be ready..."
    $timeout = 300 # 5 minutes
    $elapsed = 0
    
    do {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/realms/master" -TimeoutSec 5 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                break
            }
        } catch {
            # Continue waiting
        }
        
        if ($elapsed -ge $timeout) {
            Write-Error "Timeout waiting for Keycloak to start"
            exit 1
        }
        
        Start-Sleep -Seconds 5
        $elapsed += 5
        Write-Status "Still waiting for Keycloak... ($elapsed/$timeout seconds)"
    } while ($true)
    
    Write-Success "Keycloak is ready!"
    
    # Start the backend
    & docker-compose up -d backend
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to start backend"
        exit 1
    }
    
    Write-Success "All services started successfully!"
    Write-Status "Services status:"
    & docker-compose ps
    
    Write-Host ""
    Write-Success "🚀 $ProjectName is now running!"
    Write-Host ""
    Write-Host "📋 Service URLs:" -ForegroundColor Cyan
    Write-Host "   • Backend API: https://localhost:8443" -ForegroundColor Gray
    Write-Host "   • API Documentation: https://localhost:8443/swagger/index.html" -ForegroundColor Gray
    Write-Host "   • Health Check: https://localhost:8443/health" -ForegroundColor Gray
    Write-Host "   • Keycloak Admin: http://localhost:8080/admin (admin/admin123)" -ForegroundColor Gray
    Write-Host "   • Database: localhost:5432 (collabhub_user/collabhub_password123)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "📱 React Native Development:" -ForegroundColor Cyan
    Write-Host "   • iOS Simulator: https://localhost:8443/api/v1" -ForegroundColor Gray
    Write-Host "   • Android Emulator: https://10.0.2.2:8443/api/v1" -ForegroundColor Gray
    Write-Host "   • Physical Device: https://YOUR_IP:8443/api/v1" -ForegroundColor Gray
    Write-Host ""
}

# Function to stop services
function Stop-Services {
    Write-Status "Stopping services..."
    
    & docker-compose down
    
    Write-Success "All services stopped!"
}

# Function to view logs
function Show-Logs {
    param([string]$ServiceName)
    
    if ($ServiceName) {
        Write-Status "Viewing logs for $ServiceName..."
        & docker-compose logs -f $ServiceName
    } else {
        Write-Status "Viewing logs for all services..."
        & docker-compose logs -f
    }
}

# Function to run tests
function Invoke-Tests {
    Write-Status "Running tests..."
    
    & go test ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Tests failed"
        exit 1
    }
    
    Write-Success "All tests completed!"
}

# Function to clean up
function Remove-All {
    Write-Status "Cleaning up..."
    
    # Stop and remove containers
    & docker-compose down -v
    
    # Remove built binary
    if (Test-Path "collabhub-backend.exe") {
        Remove-Item "collabhub-backend.exe"
    }
    
    # Clean Go module cache
    & go clean -modcache
    
    Write-Success "Cleanup completed!"
}

# Function to show help
function Show-Help {
    Write-Host "Usage: .\dev.ps1 [COMMAND]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  setup       Setup environment and generate certificates"
    Write-Host "  start       Start all services"
    Write-Host "  stop        Stop all services"
    Write-Host "  restart     Restart all services"
    Write-Host "  logs        View logs for all services"
    Write-Host "  logs <svc>  View logs for specific service"
    Write-Host "  build       Build the application"
    Write-Host "  test        Run tests"
    Write-Host "  clean       Clean up containers and volumes"
    Write-Host "  help        Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\dev.ps1 setup           # First time setup"
    Write-Host "  .\dev.ps1 start           # Start all services"
    Write-Host "  .\dev.ps1 logs backend    # View backend logs"
    Write-Host "  .\dev.ps1 clean           # Clean everything"
}

# Main script logic
function Main {
    Write-Host ""
    Write-Status "🎵 $ProjectName Development Script"
    Write-Host ""
    
    switch ($Command) {
        'setup' {
            Test-Prerequisites
            Initialize-Environment
            New-Certificates
            Write-Success "Setup completed! Run '.\dev.ps1 start' to start services."
        }
        'start' {
            Test-Prerequisites
            Initialize-Environment
            New-Certificates
            Start-Services
        }
        'stop' {
            Stop-Services
        }
        'restart' {
            Stop-Services
            Start-Sleep -Seconds 2
            Start-Services
        }
        'logs' {
            Show-Logs -ServiceName $Service
        }
        'build' {
            Build-App
        }
        'test' {
            Invoke-Tests
        }
        'clean' {
            Remove-All
        }
        default {
            Show-Help
        }
    }
}

# Run main function
Main
