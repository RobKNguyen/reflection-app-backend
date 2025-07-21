// internal/models/user.go
package models

import (
    "errors"
    "time"
)

var ErrInvalidRequest = errors.New("invalid request")

type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    FirstName    string    `json:"first_name" db:"first_name"`
    LastName     string    `json:"last_name" db:"last_name"`
    PasswordHash string    `json:"-" db:"password_hash"` // Don't include in JSON responses
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}