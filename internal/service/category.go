// internal/service/category.go
package service

import (
    "errors"
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type CategoryService struct {
    repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
    return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategoriesByUser(userID int) ([]models.Category, error) {
    return s.repo.GetCategoriesByUser(userID)
}

func (s *CategoryService) GetSubCategoriesByCategory(categoryID int) ([]models.SubCategory, error) {
    return s.repo.GetSubCategoriesByCategory(categoryID)
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
    if category.Name == "" {
        return errors.New("category name is required")
    }
    
    if category.UserID <= 0 {
        return errors.New("user ID is required")
    }
    
    return s.repo.CreateCategory(category)
}

func (s *CategoryService) CreateSubCategory(subcategory *models.SubCategory) error {
    if subcategory.Name == "" {
        return errors.New("subcategory name is required")
    }
    
    if subcategory.CategoryID <= 0 {
        return errors.New("category ID is required")
    }
    
    return s.repo.CreateSubCategory(subcategory)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
    if category.Name == "" {
        return errors.New("category name is required")
    }
    
    return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(categoryID int) error {
    return s.repo.DeleteCategory(categoryID)
}

func (s *CategoryService) UpdateSubCategory(subcategory *models.SubCategory) error {
    if subcategory.Name == "" {
        return errors.New("subcategory name is required")
    }
    
    return s.repo.UpdateSubCategory(subcategory)
}

func (s *CategoryService) DeleteSubCategory(subcategoryID int) error {
    return s.repo.DeleteSubCategory(subcategoryID)
}