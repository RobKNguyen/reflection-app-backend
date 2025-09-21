package test

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "time"
)

type TestHandler struct{}

func NewTestHandler() *TestHandler {
    return &TestHandler{}
}

func (h *TestHandler) SQLTest(w http.ResponseWriter, r *http.Request) {
    log.Printf("SQL test endpoint called - Method: %s, URL: %s, Headers: %v", r.Method, r.URL.String(), r.Header)
    
    // Log the request body if it exists
    if r.Body != nil {
        body, err := io.ReadAll(r.Body)
        if err == nil && len(body) > 0 {
            log.Printf("Request body: %s", string(body))
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    response := map[string]interface{}{
        "success":   true,
        "message":   "Hello from Render API!",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "method":    r.Method,
        "url":       r.URL.String(),
        "source":    "reflection-app",
    }
    
    json.NewEncoder(w).Encode(response)
}

func (h *TestHandler) Echo(w http.ResponseWriter, r *http.Request) {
    log.Printf("Echo endpoint called - Method: %s", r.Method)
    
    var requestData map[string]interface{}
    
    if r.Body != nil {
        body, err := io.ReadAll(r.Body)
        if err == nil && len(body) > 0 {
            json.Unmarshal(body, &requestData)
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    response := map[string]interface{}{
        "success":     true,
        "message":     "Echo response",
        "timestamp":   time.Now().UTC().Format(time.RFC3339),
        "method":      r.Method,
        "received":    requestData,
    }
    
    json.NewEncoder(w).Encode(response)
}

// Create handles POST requests, parses JSON body and returns created resource info
func (h *TestHandler) Create(w http.ResponseWriter, r *http.Request) {
    log.Printf("Create endpoint called - Method: %s", r.Method)

    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var payload map[string]interface{}

    if r.Body == nil {
        http.Error(w, "empty body", http.StatusBadRequest)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("error reading body: %v", err)
        http.Error(w, "unable to read body", http.StatusBadRequest)
        return
    }

    if len(body) == 0 {
        http.Error(w, "empty body", http.StatusBadRequest)
        return
    }

    if err := json.Unmarshal(body, &payload); err != nil {
        log.Printf("invalid json: %v", err)
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    // minimal validation: require at least one key
    if len(payload) == 0 {
        http.Error(w, "payload must contain at least one field", http.StatusBadRequest)
        return
    }

    // generate a pseudo id using timestamp
    id := time.Now().UTC().UnixNano()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    response := map[string]interface{}{
        "success":   true,
        "message":   "Resource created",
        "id":        id,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "received":  payload,
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("error encoding response: %v", err)
    }
}