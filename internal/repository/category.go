// internal/repository/category.go
package repository

import (
    "database/sql"
    "reflection-app/internal/models"
)

type CategoryRepository struct {
    db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategoriesByUser(userID int) ([]models.Category, error) {
    query := `SELECT id, user_id, name, description, parent_id, created_at FROM categories WHERE user_id = $1 ORDER BY name`
    
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var categories []models.Category
    for rows.Next() {
        var cat models.Category
        err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Description, &cat.ParentID, &cat.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, cat)
    }
    
    return categories, nil
}

func (r *CategoryRepository) GetSubCategoriesByCategory(categoryID int) ([]models.SubCategory, error) {
    query := `SELECT id, category_id, name, description, created_at FROM sub_categories WHERE category_id = $1 ORDER BY name`
    
    rows, err := r.db.Query(query, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var subcategories []models.SubCategory
    for rows.Next() {
        var subcat models.SubCategory
        var description sql.NullString  // Handle NULL descriptions
        
        err := rows.Scan(&subcat.ID, &subcat.CategoryID, &subcat.Name, &description, &subcat.CreatedAt)
        if err != nil {
            return nil, err
        }
        
        // Convert NullString to regular string
        if description.Valid {
            subcat.Description = description.String
        } else {
            subcat.Description = ""
        }
        
        subcategories = append(subcategories, subcat)
    }
    
    return subcategories, nil
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
    query := `
        INSERT INTO categories (user_id, name, description, parent_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`
    
    return r.db.QueryRow(
        query,
        category.UserID,
        category.Name,
        category.Description,
        category.ParentID,
    ).Scan(&category.ID, &category.CreatedAt)
}

func (r *CategoryRepository) CreateSubCategory(subcategory *models.SubCategory) error {
    query := `
        INSERT INTO sub_categories (category_id, name, description)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`
    
    return r.db.QueryRow(
        query,
        subcategory.CategoryID,
        subcategory.Name,
        subcategory.Description,
    ).Scan(&subcategory.ID, &subcategory.CreatedAt)
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
    query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3`
    
    result, err := r.db.Exec(query, category.Name, category.Description, category.ID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *CategoryRepository) DeleteCategory(categoryID int) error {
    query := `DELETE FROM categories WHERE id = $1`
    
    result, err := r.db.Exec(query, categoryID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *CategoryRepository) UpdateSubCategory(subcategory *models.SubCategory) error {
    query := `UPDATE sub_categories SET name = $1, description = $2 WHERE id = $3`
    
    result, err := r.db.Exec(query, subcategory.Name, subcategory.Description, subcategory.ID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *CategoryRepository) DeleteSubCategory(subcategoryID int) error {
    query := `DELETE FROM sub_categories WHERE id = $1`
    
    result, err := r.db.Exec(query, subcategoryID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}