package constants

const (
    // Server constants
    ServerPort       = ":8080"
    ServerReadTimeout = 10 // seconds
    ServerWriteTimeout = 10 // seconds

    // Database constants
    DBDriver   = "postgres"
    DBHost     = "localhost"
    DBPort     = "5432"
    DBUser     = "your_db_user"
    DBPassword = "your_db_password"
    DBName     = "collabhub_music"

    // Keycloak constants
    KeycloakURL       = "http://localhost:8080/auth"
    KeycloakRealm     = "your_realm"
    KeycloakClientID  = "your_client_id"
    KeycloakClientSecret = "your_client_secret"

    // JWT constants
    JWTSecret         = "your_jwt_secret"
    JWTExpirationTime = 3600 // seconds

    // CORS constants
    AllowedOrigins    = "*"
)