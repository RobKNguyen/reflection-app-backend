package repository

import (
    "database/sql"
    "fmt"
    "reflection-app/internal/models"
)

type FriendshipRepository struct {
    db *sql.DB
}

func NewFriendshipRepository(db *sql.DB) *FriendshipRepository {
    return &FriendshipRepository{db: db}
}

// SendFriendRequest sends a friend request from userID to friendID
func (r *FriendshipRepository) SendFriendRequest(userID, friendID int) error {
    // Check if friendship already exists
    var existingID int
    err := r.db.QueryRow(`
        SELECT id FROM friendships 
        WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
    `, userID, friendID).Scan(&existingID)
    
    if err == nil {
        return fmt.Errorf("friendship already exists")
    }
    
    if err != sql.ErrNoRows {
        return err
    }
    
    // Insert new friendship request
    _, err = r.db.Exec(`
        INSERT INTO friendships (user_id, friend_id, status) 
        VALUES ($1, $2, 'pending')
    `, userID, friendID)
    
    return err
}

// AcceptFriendRequest accepts a friend request
func (r *FriendshipRepository) AcceptFriendRequest(userID, friendID int) error {
    result, err := r.db.Exec(`
        UPDATE friendships 
        SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
    `, friendID, userID)
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no pending friend request found")
    }
    
    return nil
}

// RejectFriendRequest rejects a friend request
func (r *FriendshipRepository) RejectFriendRequest(userID, friendID int) error {
    result, err := r.db.Exec(`
        UPDATE friendships 
        SET status = 'rejected', updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'
    `, friendID, userID)
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no pending friend request found")
    }
    
    return nil
}

// GetFriendsList gets all accepted friends for a user
func (r *FriendshipRepository) GetFriendsList(userID int) ([]models.FriendshipResponse, error) {
    rows, err := r.db.Query(`
        SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
               u.username, u.first_name, u.last_name
        FROM friendships f
        JOIN users u ON (f.friend_id = u.id)
        WHERE f.user_id = $1 AND f.status = 'accepted'
        UNION
        SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
               u.username, u.first_name, u.last_name
        FROM friendships f
        JOIN users u ON (f.user_id = u.id)
        WHERE f.friend_id = $1 AND f.status = 'accepted'
        ORDER BY created_at DESC
    `, userID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var friendships []models.FriendshipResponse
    for rows.Next() {
        var f models.FriendshipResponse
        var firstName, lastName string
        err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt,
                        &f.FriendUsername, &firstName, &lastName)
        if err != nil {
            return nil, err
        }
        
        f.FriendName = firstName + " " + lastName
        friendships = append(friendships, f)
    }
    
    return friendships, nil
}

// GetPendingRequests gets all pending friend requests for a user
func (r *FriendshipRepository) GetPendingRequests(userID int) ([]models.FriendshipResponse, error) {
    rows, err := r.db.Query(`
        SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
               u.username, u.first_name, u.last_name
        FROM friendships f
        JOIN users u ON f.user_id = u.id
        WHERE f.friend_id = $1 AND f.status = 'pending'
        ORDER BY f.created_at DESC
    `, userID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var requests []models.FriendshipResponse
    for rows.Next() {
        var f models.FriendshipResponse
        var firstName, lastName string
        err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt,
                        &f.FriendUsername, &firstName, &lastName)
        if err != nil {
            return nil, err
        }
        
        f.FriendName = firstName + " " + lastName
        requests = append(requests, f)
    }
    
    return requests, nil
}

// AreFriends checks if two users are friends
func (r *FriendshipRepository) AreFriends(userID1, userID2 int) (bool, error) {
    var exists bool
    err := r.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM friendships 
            WHERE ((user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))
            AND status = 'accepted'
        )
    `, userID1, userID2).Scan(&exists)
    
    return exists, err
}

// RemoveFriend removes a friendship
func (r *FriendshipRepository) RemoveFriend(userID, friendID int) error {
    result, err := r.db.Exec(`
        DELETE FROM friendships 
        WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
    `, userID, friendID)
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("friendship not found")
    }
    
    return nil
} 