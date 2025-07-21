package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type ReactionHandler struct {
    reactionService *service.ReactionService
}

func NewReactionHandler(reactionService *service.ReactionService) *ReactionHandler {
    return &ReactionHandler{
        reactionService: reactionService,
    }
}

// getUserIDFromQuery gets user ID from query parameters
func (h *ReactionHandler) getUserIDFromQuery(r *http.Request) (int, error) {
    userIDStr := r.URL.Query().Get("user_id")
    if userIDStr == "" {
        return 0, strconv.ErrSyntax
    }
    
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        return 0, err
    }
    
    return userID, nil
}

// AddReaction adds a reaction to a reflection
func (h *ReactionHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    var request models.ReactionRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    reflectionIDStr := r.URL.Query().Get("reflection_id")
    if reflectionIDStr == "" {
        http.Error(w, "Reflection ID is required", http.StatusBadRequest)
        return
    }
    
    reflectionID, err := strconv.Atoi(reflectionIDStr)
    if err != nil || reflectionID <= 0 {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    err = h.reactionService.AddReaction(reflectionID, userID, request.ReactionType, request.CommentText)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added successfully"})
}

// RemoveReaction removes a reaction from a reflection
func (h *ReactionHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    reflectionIDStr := r.URL.Query().Get("reflection_id")
    if reflectionIDStr == "" {
        http.Error(w, "Reflection ID is required", http.StatusBadRequest)
        return
    }
    
    reflectionID, err := strconv.Atoi(reflectionIDStr)
    if err != nil || reflectionID <= 0 {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    reactionTypeStr := r.URL.Query().Get("reaction_type")
    if reactionTypeStr == "" {
        http.Error(w, "Reaction type is required", http.StatusBadRequest)
        return
    }
    
    reactionType := models.ReactionType(reactionTypeStr)
    
    err = h.reactionService.RemoveReaction(reflectionID, userID, reactionType)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Reaction removed successfully"})
}

// GetReactionsForReflection gets all reactions for a reflection
func (h *ReactionHandler) GetReactionsForReflection(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    reflectionIDStr := r.URL.Query().Get("reflection_id")
    if reflectionIDStr == "" {
        http.Error(w, "Reflection ID is required", http.StatusBadRequest)
        return
    }
    
    reflectionID, err := strconv.Atoi(reflectionIDStr)
    if err != nil || reflectionID <= 0 {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    reactions, err := h.reactionService.GetReactionsForReflection(reflectionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reactions)
}

// GetReactionCountsForReflection gets reaction counts by type for a reflection
func (h *ReactionHandler) GetReactionCountsForReflection(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    reflectionIDStr := r.URL.Query().Get("reflection_id")
    if reflectionIDStr == "" {
        http.Error(w, "Reflection ID is required", http.StatusBadRequest)
        return
    }
    
    reflectionID, err := strconv.Atoi(reflectionIDStr)
    if err != nil || reflectionID <= 0 {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    counts, err := h.reactionService.GetReactionCountsForReflection(reflectionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(counts)
}

// GetUserReactionForReflection gets a user's reaction for a specific reflection
func (h *ReactionHandler) GetUserReactionForReflection(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    reflectionIDStr := r.URL.Query().Get("reflection_id")
    if reflectionIDStr == "" {
        http.Error(w, "Reflection ID is required", http.StatusBadRequest)
        return
    }
    
    reflectionID, err := strconv.Atoi(reflectionIDStr)
    if err != nil || reflectionID <= 0 {
        http.Error(w, "Invalid reflection ID", http.StatusBadRequest)
        return
    }
    
    reaction, err := h.reactionService.GetUserReactionForReflection(reflectionID, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reaction)
}

// GetReactionPrompts gets all active reaction prompts
func (h *ReactionHandler) GetReactionPrompts(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    prompts, err := h.reactionService.GetReactionPrompts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(prompts)
} 