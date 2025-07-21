package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type ActionHandler struct {
    service *service.ActionService
}

func NewActionHandler(service *service.ActionService) *ActionHandler {
    return &ActionHandler{service: service}
}

func (h *ActionHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
    var action models.ActionItem
    
    // Read the body first for debugging
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Received action data: %s\n", string(body))
    
    if err := json.Unmarshal(body, &action); err != nil {
        fmt.Printf("JSON decode error: %v\n", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Parsed action: %+v\n", action)
    
    if err := h.service.CreateAction(&action); err != nil {
        fmt.Printf("Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(action)
}

func (h *ActionHandler) GetActionsByReflection(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reflectionID, err := strconv.Atoi(vars["reflectionId"])
    if err != nil {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    actions, err := h.service.GetActionsByReflection(reflectionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(actions)
}

func (h *ActionHandler) GetAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid action ID", http.StatusBadRequest)
        return
    }
    
    action, err := h.service.GetAction(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(action)
}

func (h *ActionHandler) UpdateAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid action ID", http.StatusBadRequest)
        return
    }
    
    var action models.ActionItem
    if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    action.ID = id
    if err := h.service.UpdateAction(&action); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(action)
}

func (h *ActionHandler) CompleteAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid action ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.CompleteAction(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *ActionHandler) UpdateActionStatus(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("UpdateActionStatus called with method: %s\n", r.Method)
    fmt.Printf("URL: %s\n", r.URL.String())
    
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        fmt.Printf("Invalid action ID: %v\n", err)
        http.Error(w, "Invalid action ID", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Action ID: %d\n", id)
    
    var request struct {
        Status string `json:"status"`
    }
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("Failed to read body: %v\n", err)
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Request body: %s\n", string(body))
    
    if err := json.Unmarshal(body, &request); err != nil {
        fmt.Printf("JSON decode error: %v\n", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("Status: %s\n", request.Status)
    
    if err := h.service.UpdateActionStatus(id, request.Status); err != nil {
        fmt.Printf("Service error: %v\n", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Get the updated action to return
    action, err := h.service.GetAction(id)
    if err != nil {
        fmt.Printf("Get action error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Printf("Successfully updated action: %+v\n", action)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(action)
}

func (h *ActionHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid action ID", http.StatusBadRequest)
        return
    }
    
    if err := h.service.DeleteAction(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}