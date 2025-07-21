package repository

import (
    "database/sql"
    "reflection-app/internal/models"
)

type ActionRepository struct {
    db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
    return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(action *models.ActionItem) error {
    query := `
        INSERT INTO actions (reflection_id, action, priority, status, due_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`
    
    return r.db.QueryRow(
        query,
        action.ReflectionID,
        action.Action,
        action.Priority,
        action.Status,
        action.DueDate,
    ).Scan(&action.ID, &action.CreatedAt, &action.UpdatedAt)
}

func (r *ActionRepository) GetByReflectionID(reflectionID int) ([]models.ActionItem, error) {
    query := `
        SELECT id, reflection_id, action, priority, status, due_date, created_at, updated_at
        FROM actions 
        WHERE reflection_id = $1 
        ORDER BY created_at DESC`
    
    rows, err := r.db.Query(query, reflectionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var actions []models.ActionItem
    for rows.Next() {
        var action models.ActionItem
        err := rows.Scan(
            &action.ID,
            &action.ReflectionID,
            &action.Action,
            &action.Priority,
            &action.Status,
            &action.DueDate,
            &action.CreatedAt,
            &action.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        actions = append(actions, action)
    }
    
    return actions, nil
}

func (r *ActionRepository) GetByID(id int) (*models.ActionItem, error) {
    query := `
        SELECT id, reflection_id, action, priority, status, due_date, created_at, updated_at
        FROM actions 
        WHERE id = $1`
    
    action := &models.ActionItem{}
    err := r.db.QueryRow(query, id).Scan(
        &action.ID,
        &action.ReflectionID,
        &action.Action,
        &action.Priority,
        &action.Status,
        &action.DueDate,
        &action.CreatedAt,
        &action.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return action, nil
}

func (r *ActionRepository) Update(action *models.ActionItem) error {
    query := `
        UPDATE actions 
        SET action = $1, priority = $2, status = $3, due_date = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5`
    
    _, err := r.db.Exec(
        query,
        action.Action,
        action.Priority,
        action.Status,
        action.DueDate,
        action.ID,
    )
    
    return err
}

func (r *ActionRepository) Delete(id int) error {
    query := `DELETE FROM actions WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}