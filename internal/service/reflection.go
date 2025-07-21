package service

import (
    "errors"
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
    "time"
)

type ReflectionService struct {
    repo *repository.ReflectionRepository
}

func NewReflectionService(repo *repository.ReflectionRepository) *ReflectionService {
    return &ReflectionService{repo: repo}
}

func (s *ReflectionService) CreateReflection(reflection *models.Reflection) error {
    // Business logic validation
    if reflection.ReflectionText == "" {
        return errors.New("reflection text is required")
    }
    
    if reflection.CategoryID <= 0 {
        return errors.New("category is required")
    }
    
    return s.repo.Create(reflection)
}

func (s *ReflectionService) GetUserReflections(userID int) ([]models.Reflection, error) {
    return s.repo.GetByUserID(userID)
}

func (s *ReflectionService) GetReflection(id int) (*models.Reflection, error) {
    if id <= 0 {
        return nil, errors.New("invalid reflection ID")
    }
    
    return s.repo.GetByID(id)
}

func (s *ReflectionService) UpdateReflection(reflection *models.Reflection) error {
    if reflection.ID <= 0 {
        return errors.New("invalid reflection ID")
    }
    
    if reflection.ReflectionText == "" {
        return errors.New("reflection text is required")
    }
    
    if reflection.CategoryID <= 0 {
        return errors.New("category is required")
    }
    
    return s.repo.Update(reflection)
}

func (s *ReflectionService) DeleteReflection(id int) error {
    if id <= 0 {
        return errors.New("invalid reflection ID")
    }
    
    return s.repo.Delete(id)
}

func (s *ReflectionService) GetReflectionsByCategory(userID, categoryID int) ([]models.Reflection, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    if categoryID <= 0 {
        return nil, errors.New("invalid category ID")
    }
    
    return s.repo.GetByCategory(userID, categoryID)
}

func (s *ReflectionService) TrackReflection(reflectionID int) error {
    if reflectionID <= 0 {
        return errors.New("invalid reflection ID")
    }
    
    return s.repo.TrackReflection(reflectionID)
}

func (s *ReflectionService) GetReflectionTrackingAnalytics(userID int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    if startDate.After(endDate) {
        return nil, errors.New("start date cannot be after end date")
    }
    
    return s.repo.GetReflectionTrackingAnalytics(userID, startDate, endDate)
}

func (s *ReflectionService) GetReflectionTrackingByCategory(userID int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    if startDate.After(endDate) {
        return nil, errors.New("start date cannot be after end date")
    }
    
    return s.repo.GetReflectionTrackingByCategory(userID, startDate, endDate)
}

// GetFriendsFeed gets public reflections from friends for the social feed
func (s *ReflectionService) GetFriendsFeed(userID int, limit, offset int) ([]models.Reflection, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if offset < 0 {
        offset = 0
    }
    
    return s.repo.GetFriendsFeed(userID, limit, offset)
}

// GetFriendReflections gets all public reflections from a specific friend
func (s *ReflectionService) GetFriendReflections(userID, friendID int, limit, offset int) ([]models.Reflection, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    if friendID <= 0 {
        return nil, errors.New("invalid friend ID")
    }
    
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if offset < 0 {
        offset = 0
    }
    
    return s.repo.GetFriendReflections(userID, friendID, limit, offset)
}

// GetFriendReflectionsByUsername gets public reflections from a friend by username
func (s *ReflectionService) GetFriendReflectionsByUsername(currentUserID int, friendUsername string, limit, offset int) ([]models.Reflection, error) {
    if currentUserID <= 0 {
        return []models.Reflection{}, errors.New("invalid current user ID")
    }
    
    if friendUsername == "" {
        return []models.Reflection{}, errors.New("friend username is required")
    }
    
    if limit <= 0 {
        limit = 10 // Default limit
    }
    
    if offset < 0 {
        offset = 0
    }
    
    reflections, err := s.repo.GetFriendReflectionsByUsername(currentUserID, friendUsername, limit, offset)
    if err != nil {
        return []models.Reflection{}, err
    }
    
    // Ensure we always return an array, even if empty
    if reflections == nil {
        return []models.Reflection{}, nil
    }
    
    return reflections, nil
}