package middleware

import (
    "crypto/rsa"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "math/big"
    "net/http"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type JWKSet struct {
    Keys []JWK `json:"keys"`
}

type JWK struct {
    Kty string `json:"kty"`
    Kid string `json:"kid"`
    Use string `json:"use"`
    N   string `json:"n"`
    E   string `json:"e"`
}

type KeycloakClaims struct {
    jwt.RegisteredClaims
    PreferredUsername string `json:"preferred_username"`
    Email            string `json:"email"`
    Name             string `json:"name"`
    RealmAccess      struct {
        Roles []string `json:"roles"`
    } `json:"realm_access"`
    ResourceAccess map[string]struct {
        Roles []string `json:"roles"`
    } `json:"resource_access"`
}

type JWTMiddleware struct {
    keycloakURL string
    realm       string
    publicKeys  map[string]*rsa.PublicKey
    lastUpdate  time.Time
    mutex       sync.RWMutex
}

func NewJWTMiddleware(keycloakURL, realm string) *JWTMiddleware {
    return &JWTMiddleware{
        keycloakURL: keycloakURL,
        realm:       realm,
        publicKeys:  make(map[string]*rsa.PublicKey),
    }
}

func (j *JWTMiddleware) ValidateJWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        // Parse token without verification first to get kid
        token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            kid, ok := token.Header["kid"].(string)
            if !ok {
                return nil, fmt.Errorf("kid header missing")
            }

            // Get public key for this kid
            publicKey, err := j.getPublicKey(kid)
            if err != nil {
                return nil, err
            }

            return publicKey, nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            c.Abort()
            return
        }

        if !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(*KeycloakClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // Validate token expiration
        if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            c.Abort()
            return
        }

        // Store user info in context
        c.Set("user_id", claims.Subject)
        c.Set("username", claims.PreferredUsername)
        c.Set("email", claims.Email)
        c.Set("name", claims.Name)
        c.Set("roles", claims.RealmAccess.Roles)

        c.Next()
    }
}

func (j *JWTMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
    j.mutex.RLock()
    // Check if we have the key cached and it's not too old
    if key, exists := j.publicKeys[kid]; exists && time.Since(j.lastUpdate) < time.Hour {
        j.mutex.RUnlock()
        return key, nil
    }
    j.mutex.RUnlock()

    j.mutex.Lock()
    defer j.mutex.Unlock()

    // Double-check after acquiring write lock
    if key, exists := j.publicKeys[kid]; exists && time.Since(j.lastUpdate) < time.Hour {
        return key, nil
    }

    // Fetch JWKs from Keycloak
    jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid_connect/certs", j.keycloakURL, j.realm)
    
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp, err := client.Get(jwksURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch JWKs: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch JWKs: status %d", resp.StatusCode)
    }

    var jwkSet JWKSet
    if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
        return nil, fmt.Errorf("failed to decode JWKs: %w", err)
    }

    // Find the key with matching kid
    for _, key := range jwkSet.Keys {
        if key.Kid == kid && key.Kty == "RSA" {
            publicKey, err := j.parseRSAPublicKey(key)
            if err != nil {
                continue
            }
            
            j.publicKeys[kid] = publicKey
            j.lastUpdate = time.Now()
            return publicKey, nil
        }
    }

    return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func (j *JWTMiddleware) parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
    nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
    if err != nil {
        return nil, fmt.Errorf("failed to decode N: %w", err)
    }

    eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
    if err != nil {
        return nil, fmt.Errorf("failed to decode E: %w", err)
    }

    n := big.NewInt(0).SetBytes(nBytes)
    e := big.NewInt(0).SetBytes(eBytes)

    if !e.IsInt64() {
        return nil, fmt.Errorf("exponent too large")
    }

    return &rsa.PublicKey{
        N: n,
        E: int(e.Int64()),
    }, nil
}