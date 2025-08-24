package services

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "sync"
    "time"
    
    "github.com/go-resty/resty/v2"
)

type KeycloakService struct {
    baseURL      string
    realm        string
    clientID     string
    clientSecret string
    adminToken   string
    tokenExpiry  time.Time
    mutex        sync.RWMutex
    client       *resty.Client
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    TokenType    string `json:"token_type"`
    Scope        string `json:"scope,omitempty"`
}

type KeycloakUser struct {
    ID          string                 `json:"id,omitempty"`
    Username    string                 `json:"username"`
    Email       string                 `json:"email"`
    FirstName   string                 `json:"firstName"`
    LastName    string                 `json:"lastName"`
    Enabled     bool                   `json:"enabled"`
    Attributes  map[string][]string    `json:"attributes,omitempty"`
    Credentials []KeycloakCredential   `json:"credentials,omitempty"`
}

type KeycloakCredential struct {
    Type      string `json:"type"`
    Value     string `json:"value"`
    Temporary bool   `json:"temporary"`
}

func NewKeycloakService(baseURL, realm, clientID, clientSecret string) *KeycloakService {
    client := resty.New()
    client.SetTimeout(10 * time.Second)
    client.SetRetryCount(3)
    
    return &KeycloakService{
        baseURL:      strings.TrimSuffix(baseURL, "/"),
        realm:        realm,
        clientID:     clientID,
        clientSecret: clientSecret,
        client:       client,
    }
}

// getAdminToken obtient un token d'administration pour les opérations admin
func (k *KeycloakService) getAdminToken(ctx context.Context) (string, error) {
    k.mutex.RLock()
    if k.adminToken != "" && time.Now().Before(k.tokenExpiry.Add(-30*time.Second)) {
        token := k.adminToken
        k.mutex.RUnlock()
        return token, nil
    }
    k.mutex.RUnlock()

    k.mutex.Lock()
    defer k.mutex.Unlock()

    // Double-check après avoir acquis le verrou
    if k.adminToken != "" && time.Now().Before(k.tokenExpiry.Add(-30*time.Second)) {
        return k.adminToken, nil
    }

    tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid_connect/token", k.baseURL, k.realm)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Content-Type", "application/x-www-form-urlencoded").
        SetFormData(map[string]string{
            "grant_type":    "client_credentials",
            "client_id":     k.clientID,
            "client_secret": k.clientSecret,
        }).
        Post(tokenURL)

    if err != nil {
        return "", fmt.Errorf("failed to get admin token: %w", err)
    }

    if resp.StatusCode() != http.StatusOK {
        return "", fmt.Errorf("failed to get admin token: status %d, body: %s", resp.StatusCode(), resp.String())
    }

    var tokenResp TokenResponse
    if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
        return "", fmt.Errorf("failed to parse token response: %w", err)
    }

    k.adminToken = tokenResp.AccessToken
    k.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

    return k.adminToken, nil
}

func (k *KeycloakService) GetUserInfo(ctx context.Context, token string) (*KeycloakUser, error) {
    if token == "" {
        return nil, fmt.Errorf("token is required")
    }

    userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid_connect/userinfo", k.baseURL, k.realm)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+token).
        Get(userInfoURL)

    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }

    switch resp.StatusCode() {
    case http.StatusOK:
        // Continue processing
    case http.StatusUnauthorized:
        return nil, fmt.Errorf("invalid or expired token")
    case http.StatusForbidden:
        return nil, fmt.Errorf("insufficient permissions")
    default:
        return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode(), resp.String())
    }

    var userInfo KeycloakUser
    if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
        return nil, fmt.Errorf("failed to parse user info: %w", err)
    }

    return &userInfo, nil
}

func (k *KeycloakService) CreateUser(ctx context.Context, user *KeycloakUser) (string, error) {
    if user == nil {
        return "", fmt.Errorf("user data is required")
    }

    if user.Username == "" || user.Email == "" {
        return "", fmt.Errorf("username and email are required")
    }

    adminToken, err := k.getAdminToken(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to get admin token: %w", err)
    }

    createUserURL := fmt.Sprintf("%s/admin/realms/%s/users", k.baseURL, k.realm)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+adminToken).
        SetHeader("Content-Type", "application/json").
        SetBody(user).
        Post(createUserURL)

    if err != nil {
        return "", fmt.Errorf("failed to create user: %w", err)
    }

    switch resp.StatusCode() {
    case http.StatusCreated:
        // Extraire l'ID utilisateur de la réponse Location header
        location := resp.Header().Get("Location")
        if location != "" {
            parts := strings.Split(location, "/")
            if len(parts) > 0 {
                return parts[len(parts)-1], nil
            }
        }
        return "", fmt.Errorf("user created but ID not found in response")
    case http.StatusConflict:
        return "", fmt.Errorf("user with username or email already exists")
    case http.StatusBadRequest:
        return "", fmt.Errorf("invalid user data: %s", resp.String())
    case http.StatusUnauthorized:
        return "", fmt.Errorf("unauthorized: invalid admin token")
    case http.StatusForbidden:
        return "", fmt.Errorf("insufficient permissions to create user")
    default:
        return "", fmt.Errorf("failed to create user: status %d, body: %s", resp.StatusCode(), resp.String())
    }
}

func (k *KeycloakService) GetUser(ctx context.Context, userID string) (*KeycloakUser, error) {
    if userID == "" {
        return nil, fmt.Errorf("user ID is required")
    }

    adminToken, err := k.getAdminToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get admin token: %w", err)
    }

    getUserURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", k.baseURL, k.realm, userID)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+adminToken).
        Get(getUserURL)

    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    switch resp.StatusCode() {
    case http.StatusOK:
        // Continue processing
    case http.StatusNotFound:
        return nil, fmt.Errorf("user not found")
    case http.StatusUnauthorized:
        return nil, fmt.Errorf("unauthorized: invalid admin token")
    default:
        return nil, fmt.Errorf("failed to get user: status %d, body: %s", resp.StatusCode(), resp.String())
    }

    var user KeycloakUser
    if err := json.Unmarshal(resp.Body(), &user); err != nil {
        return nil, fmt.Errorf("failed to parse user data: %w", err)
    }

    return &user, nil
}

func (k *KeycloakService) UpdateUser(ctx context.Context, userID string, user *KeycloakUser) error {
    if userID == "" {
        return fmt.Errorf("user ID is required")
    }
    if user == nil {
        return fmt.Errorf("user data is required")
    }

    adminToken, err := k.getAdminToken(ctx)
    if err != nil {
        return fmt.Errorf("failed to get admin token: %w", err)
    }

    updateUserURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", k.baseURL, k.realm, userID)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+adminToken).
        SetHeader("Content-Type", "application/json").
        SetBody(user).
        Put(updateUserURL)

    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    switch resp.StatusCode() {
    case http.StatusNoContent:
        return nil
    case http.StatusNotFound:
        return fmt.Errorf("user not found")
    case http.StatusBadRequest:
        return fmt.Errorf("invalid user data: %s", resp.String())
    case http.StatusUnauthorized:
        return fmt.Errorf("unauthorized: invalid admin token")
    default:
        return fmt.Errorf("failed to update user: status %d, body: %s", resp.StatusCode(), resp.String())
    }
}

func (k *KeycloakService) DeleteUser(ctx context.Context, userID string) error {
    if userID == "" {
        return fmt.Errorf("user ID is required")
    }

    adminToken, err := k.getAdminToken(ctx)
    if err != nil {
        return fmt.Errorf("failed to get admin token: %w", err)
    }

    deleteUserURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", k.baseURL, k.realm, userID)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Authorization", "Bearer "+adminToken).
        Delete(deleteUserURL)

    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    switch resp.StatusCode() {
    case http.StatusNoContent:
        return nil
    case http.StatusNotFound:
        return fmt.Errorf("user not found")
    case http.StatusUnauthorized:
        return fmt.Errorf("unauthorized: invalid admin token")
    default:
        return fmt.Errorf("failed to delete user: status %d, body: %s", resp.StatusCode(), resp.String())
    }
}

func (k *KeycloakService) ValidateToken(ctx context.Context, token string) (bool, error) {
    if token == "" {
        return false, fmt.Errorf("token is required")
    }

    introspectURL := fmt.Sprintf("%s/realms/%s/protocol/openid_connect/token/introspect", k.baseURL, k.realm)
    
    resp, err := k.client.R().
        SetContext(ctx).
        SetHeader("Content-Type", "application/x-www-form-urlencoded").
        SetFormData(map[string]string{
            "token":         token,
            "client_id":     k.clientID,
            "client_secret": k.clientSecret,
        }).
        Post(introspectURL)

    if err != nil {
        return false, fmt.Errorf("failed to validate token: %w", err)
    }

    if resp.StatusCode() != http.StatusOK {
        return false, fmt.Errorf("failed to validate token: status %d", resp.StatusCode())
    }

    var introspectionResp map[string]interface{}
    if err := json.Unmarshal(resp.Body(), &introspectionResp); err != nil {
        return false, fmt.Errorf("failed to parse introspection response: %w", err)
    }

    active, ok := introspectionResp["active"].(bool)
    if !ok {
        return false, fmt.Errorf("invalid introspection response format")
    }

    return active, nil
}