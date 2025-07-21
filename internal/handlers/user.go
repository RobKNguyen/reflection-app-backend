package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "log"
    
    "github.com/gorilla/mux"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreateUser(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.GetUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]
    
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }
    
    user, err := h.service.GetUserByUsername(username)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = id
    if err := h.service.UpdateUser(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    err = h.service.DeleteUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// SearchUsers searches for users by username
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
    log.Printf("SearchUsers called - Method: %s, URL: %s", r.Method, r.URL.String())
    
    if r.Method != http.MethodGet {
        log.Printf("Method not allowed: %s", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    query := r.URL.Query().Get("q")
    log.Printf("Search query: '%s'", query)
    
    if query == "" {
        log.Printf("Query parameter 'q' is missing")
        http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
        return
    }
    
    log.Printf("Calling service.SearchUsers with query: '%s'", query)
    users, err := h.service.SearchUsers(query)
    if err != nil {
        log.Printf("Service error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    log.Printf("Found %d users", len(users))
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}