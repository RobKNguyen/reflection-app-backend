package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "reflection-app/internal/models"
    "reflection-app/internal/service"
)

type FriendshipHandler struct {
    friendshipService *service.FriendshipService
}

func NewFriendshipHandler(friendshipService *service.FriendshipService) *FriendshipHandler {
    return &FriendshipHandler{
        friendshipService: friendshipService,
    }
}

// getUserIDFromQuery gets user ID from query parameters
func (h *FriendshipHandler) getUserIDFromQuery(r *http.Request) (int, error) {
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

// SendFriendRequest sends a friend request
func (h *FriendshipHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    var request models.FriendshipRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if request.FriendID <= 0 {
        http.Error(w, "Invalid friend ID", http.StatusBadRequest)
        return
    }
    
    err = h.friendshipService.SendFriendRequest(userID, request.FriendID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Friend request sent successfully"})
}

// AcceptFriendRequest accepts a friend request
func (h *FriendshipHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    var request models.FriendshipRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if request.FriendID <= 0 {
        http.Error(w, "Invalid friend ID", http.StatusBadRequest)
        return
    }
    
    err = h.friendshipService.AcceptFriendRequest(userID, request.FriendID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Friend request accepted successfully"})
}

// RejectFriendRequest rejects a friend request
func (h *FriendshipHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    var request models.FriendshipRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if request.FriendID <= 0 {
        http.Error(w, "Invalid friend ID", http.StatusBadRequest)
        return
    }
    
    err = h.friendshipService.RejectFriendRequest(userID, request.FriendID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Friend request rejected successfully"})
}

// GetFriendsList gets all friends for a user
func (h *FriendshipHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    friends, err := h.friendshipService.GetFriendsList(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(friends)
}

// GetPendingRequests gets all pending friend requests for a user
func (h *FriendshipHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    requests, err := h.friendshipService.GetPendingRequests(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(requests)
}

// RemoveFriend removes a friendship
func (h *FriendshipHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    userID, err := h.getUserIDFromQuery(r)
    if err != nil {
        http.Error(w, "user_id parameter is required", http.StatusBadRequest)
        return
    }
    
    friendIDStr := r.URL.Query().Get("friend_id")
    if friendIDStr == "" {
        http.Error(w, "Friend ID is required", http.StatusBadRequest)
        return
    }
    
    friendID, err := strconv.Atoi(friendIDStr)
    if err != nil || friendID <= 0 {
        http.Error(w, "Invalid friend ID", http.StatusBadRequest)
        return
    }
    
    err = h.friendshipService.RemoveFriend(userID, friendID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Friend removed successfully"})
} 