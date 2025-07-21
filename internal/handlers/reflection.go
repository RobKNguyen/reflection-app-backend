package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type ReflectionHandler struct {
    service *service.ReflectionService
}

func NewReflectionHandler(service *service.ReflectionService) *ReflectionHandler {
    return &ReflectionHandler{service: service}
}

func (h *ReflectionHandler) CreateReflection(w http.ResponseWriter, r *http.Request) {
    var reflection models.Reflection
    
    // Log the raw request body for debugging
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Received reflection data: %s\n", string(body))
    
    // Create a new reader for the JSON decoder
    r.Body = io.NopCloser(bytes.NewBuffer(body))
    
    if err := json.NewDecoder(r.Body).Decode(&reflection); err != nil {
        fmt.Printf("JSON decode error: %v\n", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Parsed reflection: %+v\n", reflection)
    fmt.Printf("Handler: Reflection date before service call: %s (Local: %s, UTC: %s)\n", 
        reflection.Date.Time().Format("2006-01-02"),
        reflection.Date.Time().Local().Format("2006-01-02 15:04:05 MST"),
        reflection.Date.Time().UTC().Format("2006-01-02 15:04:05 UTC"))
    
    if err := h.service.CreateReflection(&reflection); err != nil {
        fmt.Printf("Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(reflection)
}

func (h *ReflectionHandler) GetReflections(w http.ResponseWriter, r *http.Request) {
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    reflections, err := h.service.GetUserReflections(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflections)
}

func (h *ReflectionHandler) GetUserReflections(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["userId"]
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    reflections, err := h.service.GetUserReflections(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflections)
}

func (h *ReflectionHandler) GetReflection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    reflection, err := h.service.GetReflection(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflection)
}

func (h *ReflectionHandler) UpdateReflection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    // Log the raw request body for debugging
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("UpdateReflection: Received reflection data: %s\n", string(body))
    
    // Create a new reader for the JSON decoder
    r.Body = io.NopCloser(bytes.NewBuffer(body))
    
    var reflection models.Reflection
    if err := json.NewDecoder(r.Body).Decode(&reflection); err != nil {
        fmt.Printf("UpdateReflection: JSON decode error: %v\n", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("UpdateReflection: Parsed reflection: %+v\n", reflection)
    
    reflection.ID = id
    if err := h.service.UpdateReflection(&reflection); err != nil {
        fmt.Printf("UpdateReflection: Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflection)
}

func (h *ReflectionHandler) DeleteReflection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.DeleteReflection(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *ReflectionHandler) TrackReflection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.TrackReflection(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Reflection tracked successfully"})
}

func (h *ReflectionHandler) GetReflectionTrackingAnalytics(w http.ResponseWriter, r *http.Request) {
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get date range from query parameters
    startDateStr := r.URL.Query().Get("start_date")
    endDateStr := r.URL.Query().Get("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        http.Error(w, "start_date and end_date parameters are required", http.StatusBadRequest)
        return
    }
    
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        http.Error(w, "Invalid start_date format (YYYY-MM-DD)", http.StatusBadRequest)
        return
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        http.Error(w, "Invalid end_date format (YYYY-MM-DD)", http.StatusBadRequest)
        return
    }
    
    analytics, err := h.service.GetReflectionTrackingAnalytics(userID, startDate, endDate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(analytics)
}

func (h *ReflectionHandler) GetReflectionTrackingByCategory(w http.ResponseWriter, r *http.Request) {
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get date range from query parameters
    startDateStr := r.URL.Query().Get("start_date")
    endDateStr := r.URL.Query().Get("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        http.Error(w, "start_date and end_date parameters are required", http.StatusBadRequest)
        return
    }
    
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        http.Error(w, "Invalid start_date format (YYYY-MM-DD)", http.StatusBadRequest)
        return
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        http.Error(w, "Invalid end_date format (YYYY-MM-DD)", http.StatusBadRequest)
        return
    }
    
    analytics, err := h.service.GetReflectionTrackingByCategory(userID, startDate, endDate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(analytics)
}

// GetFriendsFeed gets public reflections from friends for the social feed
func (h *ReflectionHandler) GetFriendsFeed(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get pagination parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    
    limit := 10 // Default limit
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }
    
    offset := 0 // Default offset
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }
    
    fmt.Printf("=== BACKEND FEED DEBUG ===\n")
    fmt.Printf("User ID: %d\n", userID)
    fmt.Printf("Limit: %d, Offset: %d\n", limit, offset)
    
    reflections, err := h.service.GetFriendsFeed(userID, limit, offset)
    if err != nil {
        fmt.Printf("Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Printf("Reflections returned: %d\n", len(reflections))
    if len(reflections) > 0 {
        fmt.Printf("First reflection: %+v\n", reflections[0])
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflections)
}

// GetFriendReflections gets all public reflections from a specific friend
func (h *ReflectionHandler) GetFriendReflections(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get friend ID from query parameter
    friendIDStr := r.URL.Query().Get("friend_id")
    if friendIDStr == "" {
        http.Error(w, "friend_id parameter is required", http.StatusBadRequest)
        return
    }
    
    friendID, err := strconv.Atoi(friendIDStr)
    if err != nil {
        http.Error(w, "Invalid friend ID", http.StatusBadRequest)
        return
    }
    
    // Get pagination parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    
    limit := 10 // Default limit
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }
    
    offset := 0 // Default offset
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }
    
    reflections, err := h.service.GetFriendReflections(userID, friendID, limit, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflections)
}

// GetFriendReflectionsByUsername gets public reflections from a friend by username
func (h *ReflectionHandler) GetFriendReflectionsByUsername(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get current user ID from query parameter
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get friend username from query parameter
    friendUsername := r.URL.Query().Get("friend_username")
    if friendUsername == "" {
        http.Error(w, "friend_username parameter is required", http.StatusBadRequest)
        return
    }
    
    // Get pagination parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    
    limit := 10 // Default limit
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }
    
    offset := 0 // Default offset
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = o
        }
    }
    
    fmt.Printf("=== BACKEND FRIEND BY USERNAME DEBUG ===\n")
    fmt.Printf("Current User ID: %d\n", userID)
    fmt.Printf("Friend Username: %s\n", friendUsername)
    fmt.Printf("Limit: %d, Offset: %d\n", limit, offset)
    
    reflections, err := h.service.GetFriendReflectionsByUsername(userID, friendUsername, limit, offset)
    if err != nil {
        fmt.Printf("Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Printf("Reflections returned: %d\n", len(reflections))
    if len(reflections) > 0 {
        fmt.Printf("First reflection: %+v\n", reflections[0])
    }
    
    // Ensure we always return an array, even if empty
    if reflections == nil {
        reflections = []models.Reflection{}
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reflections)
}