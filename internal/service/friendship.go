package service

import (
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type FriendshipService struct {
    friendshipRepo *repository.FriendshipRepository
}

func NewFriendshipService(friendshipRepo *repository.FriendshipRepository) *FriendshipService {
    return &FriendshipService{
        friendshipRepo: friendshipRepo,
    }
}

// SendFriendRequest sends a friend request
func (s *FriendshipService) SendFriendRequest(userID, friendID int) error {
    if userID == friendID {
        return models.ErrInvalidRequest
    }
    
    return s.friendshipRepo.SendFriendRequest(userID, friendID)
}

// AcceptFriendRequest accepts a friend request
func (s *FriendshipService) AcceptFriendRequest(userID, friendID int) error {
    return s.friendshipRepo.AcceptFriendRequest(userID, friendID)
}

// RejectFriendRequest rejects a friend request
func (s *FriendshipService) RejectFriendRequest(userID, friendID int) error {
    return s.friendshipRepo.RejectFriendRequest(userID, friendID)
}

// GetFriendsList gets all friends for a user
func (s *FriendshipService) GetFriendsList(userID int) (*models.FriendListResponse, error) {
    friends, err := s.friendshipRepo.GetFriendsList(userID)
    if err != nil {
        return nil, err
    }
    
    return &models.FriendListResponse{
        Friends: friends,
        Count:   len(friends),
    }, nil
}

// GetPendingRequests gets all pending friend requests for a user
func (s *FriendshipService) GetPendingRequests(userID int) (*models.PendingRequestsResponse, error) {
    requests, err := s.friendshipRepo.GetPendingRequests(userID)
    if err != nil {
        return nil, err
    }
    
    return &models.PendingRequestsResponse{
        Requests: requests,
        Count:    len(requests),
    }, nil
}

// AreFriends checks if two users are friends
func (s *FriendshipService) AreFriends(userID1, userID2 int) (bool, error) {
    return s.friendshipRepo.AreFriends(userID1, userID2)
}

// RemoveFriend removes a friendship
func (s *FriendshipService) RemoveFriend(userID, friendID int) error {
    return s.friendshipRepo.RemoveFriend(userID, friendID)
} 