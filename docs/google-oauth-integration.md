# Google OAuth Integration Guide

## 1. Overview

This guide explains how to integrate Google OAuth authentication into the existing tokenprovider architecture.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         CURRENT VS GOOGLE AUTH                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

  CURRENT (Email/Password)              GOOGLE OAUTH
  ────────────────────────              ────────────────────────
  ┌────────┐                            ┌────────┐
  │ Client │                            │ Client │
  └───┬────┘                            └───┬────┘
      │ POST /login                         │ GET /auth/google
      │ {email, password}                   │
      ▼                                     ▼
  ┌────────┐                            ┌────────┐
  │ Server │                            │ Google │ ← User logs in with Google
  └───┬────┘                            └───┬────┘
      │ Verify in DB                        │ Returns authorization code
      ▼                                     ▼
  ┌────────┐                            ┌────────┐
  │  JWT   │                            │ Server │ ← Exchange code for token
  │ Token  │                            └───┬────┘
  └────────┘                                │ Get user info from Google
                                            ▼
                                        ┌────────┐
                                        │  JWT   │ ← Generate OUR token
                                        │ Token  │
                                        └────────┘
```

---

## 2. Google OAuth Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          GOOGLE OAUTH 2.0 FLOW                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────┐         ┌──────────┐         ┌──────────┐         ┌──────────┐
│  Client  │         │  Server  │         │  Google  │         │    DB    │
│ (Browser)│         │  (API)   │         │  OAuth   │         │(Postgres)│
└────┬─────┘         └────┬─────┘         └────┬─────┘         └────┬─────┘
     │                    │                    │                    │
     │ 1. GET /auth/google│                    │                    │
     │───────────────────>│                    │                    │
     │                    │                    │                    │
     │ 2. Redirect to Google OAuth URL         │                    │
     │<─ ─ ─ ─ ─ ─ ─ ─ ─ ─│                    │                    │
     │                    │                    │                    │
     │ 3. User logs in with Google account     │                    │
     │────────────────────────────────────────>│                    │
     │                    │                    │                    │
     │ 4. Google redirects with auth code      │                    │
     │<────────────────────────────────────────│                    │
     │    /auth/google/callback?code=xxx       │                    │
     │                    │                    │                    │
     │ 5. Send auth code to server             │                    │
     │───────────────────>│                    │                    │
     │                    │                    │                    │
     │                    │ 6. Exchange code   │                    │
     │                    │    for tokens      │                    │
     │                    │───────────────────>│                    │
     │                    │                    │                    │
     │                    │ 7. Return access   │                    │
     │                    │    token + id_token│                    │
     │                    │<───────────────────│                    │
     │                    │                    │                    │
     │                    │ 8. Get user info   │                    │
     │                    │    (email, name)   │                    │
     │                    │───────────────────>│                    │
     │                    │                    │                    │
     │                    │ 9. Return user     │                    │
     │                    │    profile         │                    │
     │                    │<───────────────────│                    │
     │                    │                    │                    │
     │                    │ 10. Find or Create user                 │
     │                    │────────────────────────────────────────>│
     │                    │                    │                    │
     │                    │ 11. Return user data                    │
     │                    │<────────────────────────────────────────│
     │                    │                    │                    │
     │                    │ 12. Generate JWT   │                    │
     │                    │     (our token)    │                    │
     │                    │                    │                    │
     │ 13. Return JWT token                    │                    │
     │<───────────────────│                    │                    │
     │                    │                    │                    │
```

---

## 3. Architecture Extension

### 3.1 New Provider Interface

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        EXTENDED PROVIDER ARCHITECTURE                           │
└─────────────────────────────────────────────────────────────────────────────────┘

                    ┌─────────────────────────────────┐
                    │      Provider (interface)       │
                    │         (existing)              │
                    │─────────────────────────────────│
                    │ + Generate()                    │
                    │ + Validate()                    │
                    │ + SecretKey()                   │
                    └─────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                    ▼                               ▼
        ┌───────────────────────┐     ┌───────────────────────┐
        │    jwtProvider        │     │   googleProvider      │
        │     (existing)        │     │       (NEW)           │
        │───────────────────────│     │───────────────────────│
        │ - secret: string      │     │ - clientID: string    │
        │ - prefix: string      │     │ - clientSecret: string│
        │───────────────────────│     │ - redirectURL: string │
        │ + Generate()          │     │───────────────────────│
        │ + Validate()          │     │ + GetAuthURL()        │
        │                       │     │ + ExchangeCode()      │
        └───────────────────────┘     │ + GetUserInfo()       │
                                      └───────────────────────┘


                    ┌─────────────────────────────────┐
                    │   OAuthProvider (interface)     │
                    │           (NEW)                 │
                    │─────────────────────────────────│
                    │ + GetAuthURL(state) string      │
                    │ + ExchangeCode(code) (*Token)   │
                    │ + GetUserInfo(token) (*User)    │
                    └─────────────────────────────────┘
```

### 3.2 Directory Structure

```
component/
├── tokenprovider/
│   ├── provider.go              # Existing - JWT interfaces
│   └── jwt/
│       └── jwt.go               # Existing - JWT implementation
│
└── oauthprovider/               # NEW FOLDER
    ├── provider.go              # OAuth interface definitions
    └── google/
        └── google.go            # Google OAuth implementation
```

---

## 4. Implementation Details

### 4.1 OAuth Provider Interface

```go
// component/oauthprovider/provider.go

package oauthprovider

import "errors"

// OAuthUserInfo contains user information from OAuth provider
type OAuthUserInfo struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
}

// OAuthToken contains tokens from OAuth provider
type OAuthToken struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    IDToken      string `json:"id_token,omitempty"`
}

// OAuthProvider interface for OAuth implementations
type OAuthProvider interface {
    // GetAuthURL returns the URL to redirect user for authentication
    GetAuthURL(state string) string

    // ExchangeCode exchanges authorization code for tokens
    ExchangeCode(code string) (*OAuthToken, error)

    // GetUserInfo retrieves user information using access token
    GetUserInfo(accessToken string) (*OAuthUserInfo, error)

    // GetProviderName returns the provider name (e.g., "google", "facebook")
    GetProviderName() string
}

// Common errors
var (
    ErrInvalidCode     = errors.New("invalid authorization code")
    ErrInvalidToken    = errors.New("invalid access token")
    ErrUserInfoFailed  = errors.New("failed to get user info")
    ErrProviderError   = errors.New("oauth provider error")
)
```

### 4.2 Google Provider Implementation

```go
// component/oauthprovider/google/google.go

package google

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "todolist/component/oauthprovider"
)

type googleProvider struct {
    clientID     string
    clientSecret string
    redirectURL  string
    scopes       []string
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(clientID, clientSecret, redirectURL string) *googleProvider {
    return &googleProvider{
        clientID:     clientID,
        clientSecret: clientSecret,
        redirectURL:  redirectURL,
        scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
    }
}

func (g *googleProvider) GetProviderName() string {
    return "google"
}

// GetAuthURL generates Google OAuth URL
func (g *googleProvider) GetAuthURL(state string) string {
    baseURL := "https://accounts.google.com/o/oauth2/v2/auth"

    params := url.Values{}
    params.Add("client_id", g.clientID)
    params.Add("redirect_uri", g.redirectURL)
    params.Add("response_type", "code")
    params.Add("scope", strings.Join(g.scopes, " "))
    params.Add("state", state)                    // CSRF protection
    params.Add("access_type", "offline")          // Get refresh token
    params.Add("prompt", "consent")               // Force consent screen

    return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCode exchanges authorization code for tokens
func (g *googleProvider) ExchangeCode(code string) (*oauthprovider.OAuthToken, error) {
    tokenURL := "https://oauth2.googleapis.com/token"

    data := url.Values{}
    data.Set("code", code)
    data.Set("client_id", g.clientID)
    data.Set("client_secret", g.clientSecret)
    data.Set("redirect_uri", g.redirectURL)
    data.Set("grant_type", "authorization_code")

    resp, err := http.PostForm(tokenURL, data)
    if err != nil {
        return nil, oauthprovider.ErrProviderError
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, oauthprovider.ErrInvalidCode
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var token oauthprovider.OAuthToken
    if err := json.Unmarshal(body, &token); err != nil {
        return nil, err
    }

    return &token, nil
}

// GetUserInfo retrieves user profile from Google
func (g *googleProvider) GetUserInfo(accessToken string) (*oauthprovider.OAuthUserInfo, error) {
    userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

    req, err := http.NewRequest("GET", userInfoURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+accessToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, oauthprovider.ErrProviderError
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, oauthprovider.ErrInvalidToken
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var userInfo oauthprovider.OAuthUserInfo
    if err := json.Unmarshal(body, &userInfo); err != nil {
        return nil, err
    }

    return &userInfo, nil
}
```

---

## 5. Handler Implementation

### 5.1 Google Auth Handlers

```go
// module/user/transport/gin/google_auth_handler.go

package ginuser

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "todolist/common"
    "todolist/component/oauthprovider"
    "todolist/component/tokenprovider"
    "todolist/module/user/biz"
    "todolist/module/user/storage"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// generateState creates a random state for CSRF protection
func generateState() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// GoogleLogin redirects user to Google OAuth
func GoogleLogin(googleProvider oauthprovider.OAuthProvider) gin.HandlerFunc {
    return func(c *gin.Context) {
        state := generateState()

        // Store state in cookie for verification
        c.SetCookie("oauth_state", state, 600, "/", "", false, true)

        // Redirect to Google
        authURL := googleProvider.GetAuthURL(state)
        c.Redirect(http.StatusTemporaryRedirect, authURL)
    }
}

// GoogleCallback handles Google OAuth callback
func GoogleCallback(
    db *gorm.DB,
    googleProvider oauthprovider.OAuthProvider,
    tokenProvider tokenprovider.Provider,
) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Verify state (CSRF protection)
        state := c.Query("state")
        savedState, err := c.Cookie("oauth_state")
        if err != nil || state != savedState {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
            return
        }

        // 2. Get authorization code
        code := c.Query("code")
        if code == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
            return
        }

        // 3. Exchange code for tokens
        oauthToken, err := googleProvider.ExchangeCode(code)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange code"})
            return
        }

        // 4. Get user info from Google
        userInfo, err := googleProvider.GetUserInfo(oauthToken.AccessToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user info"})
            return
        }

        // 5. Find or create user in database
        store := storage.NewSQLStore(db)
        business := biz.NewGoogleLoginBusiness(store, tokenProvider)

        token, err := business.LoginWithGoogle(c.Request.Context(), userInfo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // 6. Return our JWT token
        c.JSON(http.StatusOK, common.SimpleSuccessResponse(token))
    }
}
```

### 5.2 Google Login Business Logic

```go
// module/user/biz/google_login.go

package biz

import (
    "context"
    "todolist/common"
    "todolist/component/oauthprovider"
    "todolist/component/tokenprovider"
    "todolist/module/user/model"
)

type GoogleLoginStorage interface {
    FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error)
    CreateUser(ctx context.Context, data *model.UserCreate) error
}

type googleLoginBusiness struct {
    store         GoogleLoginStorage
    tokenProvider tokenprovider.Provider
    expiry        int
}

func NewGoogleLoginBusiness(
    store GoogleLoginStorage,
    tokenProvider tokenprovider.Provider,
) *googleLoginBusiness {
    return &googleLoginBusiness{
        store:         store,
        tokenProvider: tokenProvider,
        expiry:        60 * 60 * 24 * 30, // 30 days
    }
}

func (b *googleLoginBusiness) LoginWithGoogle(
    ctx context.Context,
    googleUser *oauthprovider.OAuthUserInfo,
) (tokenprovider.Token, error) {

    // 1. Try to find existing user by email
    user, err := b.store.FindUser(ctx, map[string]interface{}{
        "email": googleUser.Email,
    })

    // 2. If user not found, create new user
    if err == common.RecordNotFound {
        newUser := &model.UserCreate{
            Email:     googleUser.Email,
            FirstName: googleUser.GivenName,
            LastName:  googleUser.FamilyName,
            Password:  "",  // No password for OAuth users
            Salt:      "",
            // You might want to add: GoogleID, Avatar, etc.
        }

        if err := b.store.CreateUser(ctx, newUser); err != nil {
            return nil, common.ErrCanNotCreateEntity(model.EntityName, err)
        }

        // Fetch the created user
        user, err = b.store.FindUser(ctx, map[string]interface{}{
            "email": googleUser.Email,
        })
        if err != nil {
            return nil, err
        }
    } else if err != nil {
        return nil, common.ErrDB(err)
    }

    // 3. Generate our JWT token (using existing tokenprovider)
    payload := &common.TokenPayload{
        UId:   user.Id,
        URole: user.Role.String(),
    }

    token, err := b.tokenProvider.Generate(payload, b.expiry)
    if err != nil {
        return nil, common.ErrInternal(err)
    }

    return token, nil
}
```

---

## 6. Route Configuration

### 6.1 Update main.go

```go
// main.go

func main() {
    db := db.ConnectDB()

    // Existing JWT provider
    systemSecret := os.Getenv("SYSTEM_SECRET")
    tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)

    // NEW: Google OAuth provider
    googleProvider := google.NewGoogleProvider(
        os.Getenv("GOOGLE_CLIENT_ID"),
        os.Getenv("GOOGLE_CLIENT_SECRET"),
        os.Getenv("GOOGLE_REDIRECT_URL"), // e.g., "http://localhost:3001/v1/auth/google/callback"
    )

    r := gin.Default()

    v1 := r.Group("v1")
    {
        // Existing routes
        v1.POST("/register", ginuser.Register(db))
        v1.POST("/login", ginuser.Login(db, tokenProvider))

        // NEW: Google OAuth routes
        auth := v1.Group("/auth")
        {
            auth.GET("/google", ginuser.GoogleLogin(googleProvider))
            auth.GET("/google/callback", ginuser.GoogleCallback(db, googleProvider, tokenProvider))
        }

        // ... other routes
    }

    r.Run(":3001")
}
```

### 6.2 Environment Variables

```env
# .env

# Existing
DB_USER=postgres
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=todolist
DB_SSLMODE=disable
SYSTEM_SECRET=your-jwt-secret-key

# NEW: Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3001/v1/auth/google/callback
```

---

## 7. Google Cloud Console Setup

### Step 1: Create Project
```
1. Go to https://console.cloud.google.com/
2. Create a new project or select existing one
3. Enable "Google+ API" or "Google People API"
```

### Step 2: Create OAuth Credentials
```
1. Go to "APIs & Services" > "Credentials"
2. Click "Create Credentials" > "OAuth client ID"
3. Select "Web application"
4. Add Authorized redirect URIs:
   - http://localhost:3001/v1/auth/google/callback (development)
   - https://yourdomain.com/v1/auth/google/callback (production)
5. Copy Client ID and Client Secret
```

### Step 3: Configure Consent Screen
```
1. Go to "OAuth consent screen"
2. Select "External" (for testing)
3. Fill in App name, User support email
4. Add scopes:
   - .../auth/userinfo.email
   - .../auth/userinfo.profile
5. Add test users (for development)
```

---

## 8. Complete Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                      COMPLETE GOOGLE OAUTH FLOW                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────┐                                                              ┌─────────┐
│ Browser │                                                              │  Google │
└────┬────┘                                                              └────┬────┘
     │                                                                        │
     │  1. Click "Login with Google"                                          │
     │  ─────────────────────────────►  GET /v1/auth/google                   │
     │                                                                        │
     │                              ┌─────────────────────────────────────┐   │
     │                              │ GoogleLogin Handler:                │   │
     │                              │ 1. Generate random state            │   │
     │                              │ 2. Store state in cookie            │   │
     │                              │ 3. Build Google OAuth URL           │   │
     │                              │ 4. Redirect to Google               │   │
     │                              └─────────────────────────────────────┘   │
     │                                                                        │
     │  2. Redirect to Google OAuth URL                                       │
     │  ◄─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─    │
     │                                                                        │
     │  3. User sees Google login page                                        │
     │  ─────────────────────────────────────────────────────────────────────►│
     │                                                                        │
     │                                          ┌─────────────────────────┐   │
     │                                          │ Google:                 │   │
     │                                          │ 1. Show login form      │   │
     │                                          │ 2. Verify credentials   │   │
     │                                          │ 3. Show consent screen  │   │
     │                                          │ 4. Generate auth code   │   │
     │                                          └─────────────────────────┘   │
     │                                                                        │
     │  4. Redirect with authorization code                                   │
     │  ◄─────────────────────────────────────────────────────────────────────│
     │     /v1/auth/google/callback?code=xxx&state=yyy                        │
     │                                                                        │
     │  5. Browser follows redirect                                           │
     │  ─────────────────────────────►  GET /v1/auth/google/callback          │
     │                                                                        │
     │                              ┌─────────────────────────────────────┐   │
     │                              │ GoogleCallback Handler:             │   │
     │                              │                                     │   │
     │                              │ 1. Verify state (CSRF check)        │   │
     │                              │    - Compare with cookie            │   │
     │                              │                                     │   │
     │                              │ 2. Exchange code for tokens ───────────►│
     │                              │    POST oauth2.googleapis.com/token │   │
     │                              │    ◄─────────────────────────────────────│
     │                              │    {access_token, id_token}         │   │
     │                              │                                     │   │
     │                              │ 3. Get user info ───────────────────────►│
     │                              │    GET googleapis.com/oauth2/v2/    │   │
     │                              │        userinfo                     │   │
     │                              │    ◄─────────────────────────────────────│
     │                              │    {email, name, picture}           │   │
     │                              │                                     │   │
     │                              │ 4. Find or Create user in DB        │   │
     │                              │    - Check if email exists          │   │
     │                              │    - Create new user if not         │   │
     │                              │                                     │   │
     │                              │ 5. Generate JWT (our token)         │   │
     │                              │    - Use existing jwtProvider       │   │
     │                              │    - Same token format as login     │   │
     │                              └─────────────────────────────────────┘   │
     │                                                                        │
     │  6. Return JWT token                                                   │
     │  ◄─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─                                         │
     │     {                                                                  │
     │       "data": {                                                        │
     │         "token": "eyJhbGciOiJIUzI1NiIs...",                            │
     │         "created": "2024-01-05T10:30:00Z",                             │
     │         "expiry": 2592000                                              │
     │       }                                                                │
     │     }                                                                  │
     │                                                                        │
     │  7. Use JWT for authenticated requests                                 │
     │  ─────────────────────────────►  Authorization: Bearer <token>         │
     │                                                                        │
```

---

## 9. Database Schema Update (Optional)

If you want to track OAuth provider information:

```sql
-- Add columns to users table
ALTER TABLE users ADD COLUMN google_id VARCHAR(255);
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(500);
ALTER TABLE users ADD COLUMN auth_provider VARCHAR(50) DEFAULT 'local';

-- Create index for faster lookup
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_auth_provider ON users(auth_provider);
```

Updated User model:
```go
type User struct {
    common.SQLModel
    Email        string    `json:"email" gorm:"column:email;"`
    Password     string    `json:"-" gorm:"column:password;"`
    Salt         string    `json:"-" gorm:"column:salt;"`
    LastName     string    `json:"last_name" gorm:"column:last_name;"`
    FirstName    string    `json:"first_name" gorm:"column:first_name;"`
    Phone        string    `json:"phone" gorm:"column:phone;"`
    Role         *UserRole `json:"role" gorm:"column:role;"`
    Status       int       `json:"status" gorm:"column:status;"`
    // NEW fields for OAuth
    GoogleID     string    `json:"-" gorm:"column:google_id;"`
    AvatarURL    string    `json:"avatar_url" gorm:"column:avatar_url;"`
    AuthProvider string    `json:"auth_provider" gorm:"column:auth_provider;default:local"`
}
```

---

## 10. Security Considerations

| Aspect | Implementation |
|--------|----------------|
| **CSRF Protection** | Use `state` parameter with random value stored in cookie |
| **Token Storage** | Store JWT in httpOnly cookie or secure localStorage |
| **HTTPS** | Always use HTTPS in production |
| **Scope Limitation** | Request only necessary scopes (email, profile) |
| **Token Expiry** | Set reasonable expiry for JWT tokens |
| **Refresh Tokens** | Implement refresh token rotation if needed |

---

## 11. Testing with Postman

### Step 1: Get Google Auth URL
```
GET http://localhost:3001/v1/auth/google

Response: Redirects to Google login page
```

### Step 2: After Google Login
```
Google redirects to:
GET http://localhost:3001/v1/auth/google/callback?code=xxx&state=yyy

Response:
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "created": "2024-01-05T10:30:00Z",
    "expiry": 2592000
  }
}
```

### Step 3: Use Token for API Calls
```
GET http://localhost:3001/v1/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Response:
{
  "data": {
    "id": 1,
    "email": "user@gmail.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user"
  }
}
```
