package service

import (
    "errors"
    "log"
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
    if user.Username == "" {
        return errors.New("username is required")
    }
    
    if user.Email == "" {
        return errors.New("email is required")
    }
    
    return s.repo.Create(user)
}

func (s *UserService) GetUser(id int) (*models.User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    return s.repo.GetByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
    if username == "" {
        return nil, errors.New("username is required")
    }
    
    return s.repo.GetByUsername(username)
}

func (s *UserService) UpdateUser(user *models.User) error {
    if user.ID <= 0 {
        return errors.New("invalid user ID")
    }
    
    if user.Username == "" {
        return errors.New("username is required")
    }
    
    if user.Email == "" {
        return errors.New("email is required")
    }
    
    return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
    if id <= 0 {
        return errors.New("invalid user ID")
    }
    
    return s.repo.Delete(id)
}

// SearchUsers searches for users by username
func (s *UserService) SearchUsers(query string) ([]models.User, error) {
    log.Printf("UserService.SearchUsers called with query: '%s'", query)
    
    if query == "" {
        log.Printf("Search query is empty")
        return nil, errors.New("search query is required")
    }
    
    log.Printf("Calling repository.SearchByUsername with query: '%s'", query)
    users, err := s.repo.SearchByUsername(query)
    if err != nil {
        log.Printf("Repository error: %v", err)
        return nil, err
    }
    
    log.Printf("Repository returned %d users", len(users))
    return users, nil
}