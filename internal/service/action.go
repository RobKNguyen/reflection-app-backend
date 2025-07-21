package service

import (
    "errors"
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type ActionService struct {
    repo *repository.ActionRepository
}

func NewActionService(repo *repository.ActionRepository) *ActionService {
    return &ActionService{repo: repo}
}

func (s *ActionService) CreateAction(action *models.ActionItem) error {
    if action.Action == "" {
        return errors.New("action text is required")
    }
    
    if action.ReflectionID <= 0 {
        return errors.New("reflection ID is required")
    }
    
    return s.repo.Create(action)
}

func (s *ActionService) GetActionsByReflection(reflectionID int) ([]models.ActionItem, error) {
    if reflectionID <= 0 {
        return nil, errors.New("invalid reflection ID")
    }
    
    return s.repo.GetByReflectionID(reflectionID)
}

func (s *ActionService) GetAction(id int) (*models.ActionItem, error) {
    if id <= 0 {
        return nil, errors.New("invalid action ID")
    }
    
    return s.repo.GetByID(id)
}

func (s *ActionService) UpdateAction(action *models.ActionItem) error {
    if action.ID <= 0 {
        return errors.New("invalid action ID")
    }
    
    if action.Action == "" {
        return errors.New("action text is required")
    }
    
    return s.repo.Update(action)
}

func (s *ActionService) CompleteAction(id int) error {
    action, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    
    action.Status = "Done"
    
    return s.repo.Update(action)
}

func (s *ActionService) UpdateActionStatus(id int, status string) error {
    if id <= 0 {
        return errors.New("invalid action ID")
    }
    
    if status != "Done" && status != "Pending" {
        return errors.New("invalid status: must be 'Done' or 'Pending'")
    }
    
    action, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    
    action.Status = status
    
    return s.repo.Update(action)
}

func (s *ActionService) DeleteAction(id int) error {
    if id <= 0 {
        return errors.New("invalid action ID")
    }
    
    return s.repo.Delete(id)
}