package repository

import (
    "database/sql"
    "time"
    "reflection-app/internal/models"
    "log"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (username, email, first_name, last_name, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
    now := time.Now()
    err := r.db.QueryRow(query, 
        user.Username, 
        user.Email, 
        user.FirstName, 
        user.LastName, 
        user.PasswordHash, 
        now, 
        now,
    ).Scan(&user.ID)
    
    if err != nil {
        return err
    }
    
    user.CreatedAt = now
    user.UpdatedAt = now
    return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users WHERE id = $1
    `
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.FirstName, 
        &user.LastName, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users WHERE username = $1
    `
    err := r.db.QueryRow(query, username).Scan(
        &user.ID, &user.Username, &user.Email, &user.FirstName, 
        &user.LastName, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users WHERE email = $1
    `
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Username, &user.Email, &user.FirstName, 
        &user.LastName, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
    query := `
        UPDATE users 
        SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5
        WHERE id = $6
    `
    _, err := r.db.Exec(query, 
        user.Username, 
        user.Email, 
        user.FirstName, 
        user.LastName, 
        time.Now(),
        user.ID,
    )
    return err
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}

// SearchByUsername searches for users by username (case-insensitive)
func (r *UserRepository) SearchByUsername(query string) ([]models.User, error) {
    log.Printf("UserRepository.SearchByUsername called with query: '%s'", query)
    
    sqlQuery := `
        SELECT id, username, email, first_name, last_name, created_at, updated_at
        FROM users 
        WHERE username ILIKE $1
        ORDER BY username ASC
        LIMIT 10`
    
    log.Printf("Executing SQL query: %s with parameter: '%%%s%%'", sqlQuery, query)
    
    rows, err := r.db.Query(sqlQuery, "%"+query+"%")
    if err != nil {
        log.Printf("Database query error: %v", err)
        return nil, err
    }
    defer rows.Close()
    
    var users []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.FirstName,
            &user.LastName,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            log.Printf("Row scan error: %v", err)
            return nil, err
        }
        users = append(users, user)
    }
    
    log.Printf("Found %d users in database", len(users))
    return users, nil
}

// GetAll retrieves all users (optional, for admin purposes)
func (r *UserRepository) GetAll() ([]*models.User, error) {
    query := `
        SELECT id, username, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users ORDER BY created_at DESC
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*models.User
    for rows.Next() {
        user := &models.User{}
        err := rows.Scan(
            &user.ID, &user.Username, &user.Email, &user.FirstName,
            &user.LastName, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}