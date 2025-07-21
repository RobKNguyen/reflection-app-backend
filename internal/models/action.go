// Updated internal/models/action.go
package models

import "time"

type ActionItem struct {
    ID           int        `json:"id" db:"id"`
    ReflectionID int        `json:"reflection_id" db:"reflection_id"`
    Action       string     `json:"action" db:"action"`
    Priority     string     `json:"priority" db:"priority"`
    Status       string     `json:"status" db:"status"`
    DueDate      *time.Time `json:"due_date" db:"due_date"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}