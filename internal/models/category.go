// internal/models/category.go
package models

import "time"

type Category struct {
    ID          int       `json:"id" db:"id"`
    UserID      int       `json:"user_id" db:"user_id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    ParentID    *int      `json:"parent_id" db:"parent_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type SubCategory struct {
    ID          int       `json:"id" db:"id"`
    CategoryID  int       `json:"category_id" db:"category_id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}