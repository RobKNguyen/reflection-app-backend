package handlers

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
    "os"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// isDevelopment checks if we're in development mode
func isDevelopment() bool {
    return os.Getenv("ENV") != "production"
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if req.Username == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
        return
    }
    
    // Pass pointer to the request
    response, err := h.authService.Register(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Set secure HTTP-only cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    response.Token,
        Path:     "/",
        HttpOnly: true,
        Secure:   !isDevelopment(), // Only secure in production
        SameSite: http.SameSiteStrictMode,
        MaxAge:   86400, // 24 hours
    })
    
    // Set user info in a separate cookie (for frontend access)
    userData, _ := json.Marshal(response.User)
    encodedUserData := base64.StdEncoding.EncodeToString(userData)
    http.SetCookie(w, &http.Cookie{
        Name:     "user_data",
        Value:    encodedUserData,
        Path:     "/",
        HttpOnly: false, // Allow JavaScript access for user info
        Secure:   !isDevelopment(), // Only secure in production
        SameSite: http.SameSiteStrictMode,
        MaxAge:   86400, // 24 hours
    })
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": response.User,
        "message": "Registration successful",
    })
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate required fields
    if req.Username == "" || req.Password == "" {
        http.Error(w, "Username and password are required", http.StatusBadRequest)
        return
    }
    
    // Pass pointer to the request
    response, err := h.authService.Login(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    
    // Set secure HTTP-only cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    response.Token,
        Path:     "/",
        HttpOnly: true,
        Secure:   !isDevelopment(), // Only secure in production
        SameSite: http.SameSiteStrictMode,
        MaxAge:   86400, // 24 hours
    })
    
    // Set user info in a separate cookie (for frontend access)
    userData, _ := json.Marshal(response.User)
    encodedUserData := base64.StdEncoding.EncodeToString(userData)
    http.SetCookie(w, &http.Cookie{
        Name:     "user_data",
        Value:    encodedUserData,
        Path:     "/",
        HttpOnly: false, // Allow JavaScript access for user info
        Secure:   !isDevelopment(), // Only secure in production
        SameSite: http.SameSiteStrictMode,
        MaxAge:   86400, // 24 hours
    })
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": response.User,
        "message": "Login successful",
    })
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    // Clear auth token cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   !isDevelopment(),
        SameSite: http.SameSiteStrictMode,
        MaxAge:   -1, // Delete the cookie
    })
    
    // Clear user data cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "user_data",
        Value:    "",
        Path:     "/",
        HttpOnly: false,
        Secure:   !isDevelopment(),
        SameSite: http.SameSiteStrictMode,
        MaxAge:   -1, // Delete the cookie
    })
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Logout successful",
    })
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    // Get token from cookie
    cookie, err := r.Cookie("auth_token")
    if err != nil {
        http.Error(w, "No authentication token", http.StatusUnauthorized)
        return
    }
    
    // Validate token
    claims, err := h.authService.ValidateToken(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
        return
    }
    
    // Get user from database
    user, err := h.authService.GetUserByID(claims.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}