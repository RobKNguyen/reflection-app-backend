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