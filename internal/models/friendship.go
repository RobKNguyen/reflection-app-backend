package models

import (
    "time"
)

type Friendship struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    FriendID  int       `json:"friend_id"`
    Status    string    `json:"status"` // pending, accepted, rejected
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type FriendshipRequest struct {
    FriendID int `json:"friend_id"`
}

type FriendshipResponse struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    FriendID  int       `json:"friend_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // Additional fields for UI
    FriendUsername string `json:"friend_username,omitempty"`
    FriendName     string `json:"friend_name,omitempty"`
}

type FriendListResponse struct {
    Friends []FriendshipResponse `json:"friends"`
    Count   int                 `json:"count"`
}

type PendingRequestsResponse struct {
    Requests []FriendshipResponse `json:"requests"`
    Count    int                 `json:"count"`
} 