package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    
    _ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
    // Get database URL from environment variable (for Render deployment)
    dbURL := os.Getenv("DATABASE_URL")
    
    // Fallback to local development settings
    if dbURL == "" {
        dbURL = "postgres://postgres:password@localhost:5432/reflection_db?sslmode=disable"
    }
    
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Connected to PostgreSQL database")
    return db, nil
}

// Updated internal/database/connection.go - Add this to CreateTables function
// Replace the existing CreateTables function in internal/database/connection.go
// (or wherever your CreateTables function currently exists)

func CreateTables(db *sql.DB) error {
    // Create users table
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        first_name VARCHAR(50) NOT NULL,
        last_name VARCHAR(50) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`
    
    if _, err := db.Exec(userTable); err != nil {
        log.Printf("Error creating users table: %v", err)
        return err
    }
    log.Println("Created users table")

    // Create categories table
    categoryTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        parent_id INTEGER REFERENCES categories(id),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, name)
    )`
    
    if _, err := db.Exec(categoryTable); err != nil {
        log.Printf("Error creating categories table: %v", err)
        return err
    }
    log.Println("Created categories table")

    // Create sub_categories table
    subCategoryTable := `
    CREATE TABLE IF NOT EXISTS sub_categories (
        id SERIAL PRIMARY KEY,
        category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(category_id, name)
    )`
    
    if _, err := db.Exec(subCategoryTable); err != nil {
        log.Printf("Error creating sub_categories table: %v", err)
        return err
    }
    log.Println("Created sub_categories table")

    // Create reflections table
    reflectionTable := `
CREATE TABLE IF NOT EXISTS reflections (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id),
    sub_category_id INTEGER REFERENCES sub_categories(id),
    date DATE DEFAULT CURRENT_DATE,
    reflection_text VARCHAR(500) NOT NULL,
    reflection_detail TEXT NOT NULL,
    tags TEXT[],
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)`
    
    if _, err := db.Exec(reflectionTable); err != nil {
        log.Printf("Error creating reflections table: %v", err)
        return err
    }
    log.Println("Created reflections table")

    // Create actions table
    actionTable := `
    CREATE TABLE IF NOT EXISTS actions (
        id SERIAL PRIMARY KEY,
        reflection_id INTEGER NOT NULL REFERENCES reflections(id) ON DELETE CASCADE,
        action VARCHAR(200) NOT NULL,
        priority VARCHAR(20) DEFAULT 'Medium',
        status VARCHAR(20) DEFAULT 'Pending',
        due_date TIMESTAMP WITH TIME ZONE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`
    
    if _, err := db.Exec(actionTable); err != nil {
        log.Printf("Error creating actions table: %v", err)
        return err
    }
    log.Println("Created actions table")

    // Create reflection_tracking table
    trackingTable := `
    CREATE TABLE IF NOT EXISTS reflection_tracking (
        id SERIAL PRIMARY KEY,
        reflection_id INTEGER NOT NULL REFERENCES reflections(id) ON DELETE CASCADE,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        reflected_date DATE NOT NULL DEFAULT CURRENT_DATE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(reflection_id, user_id, reflected_date)
    )`
    
    if _, err := db.Exec(trackingTable); err != nil {
        log.Printf("Error creating reflection_tracking table: %v", err)
        return err
    }
    log.Println("Created reflection_tracking table")

    // Add visibility field to reflections table
    alterReflectionsTable := `
    ALTER TABLE reflections ADD COLUMN IF NOT EXISTS visibility VARCHAR(20) DEFAULT 'private' CHECK (visibility IN ('private', 'public'))
    `
    
    if _, err := db.Exec(alterReflectionsTable); err != nil {
        log.Printf("Error adding visibility column to reflections table: %v", err)
        return err
    }
    log.Println("Added visibility column to reflections table")

    // Create friendships table
    friendshipsTable := `
    CREATE TABLE IF NOT EXISTS friendships (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected')),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, friend_id)
    )`
    
    if _, err := db.Exec(friendshipsTable); err != nil {
        log.Printf("Error creating friendships table: %v", err)
        return err
    }
    log.Println("Created friendships table")

    // Create reflection_reactions table
    reactionsTable := `
    CREATE TABLE IF NOT EXISTS reflection_reactions (
        id SERIAL PRIMARY KEY,
        reflection_id INTEGER NOT NULL REFERENCES reflections(id) ON DELETE CASCADE,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        reaction_type VARCHAR(50) NOT NULL CHECK (reaction_type IN (
            'ask_me_about_this',
            'similar_experience', 
            'update_me',
            'accountability_buddy',
            'different_angle'
        )),
        comment_text VARCHAR(100),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(reflection_id, user_id, reaction_type)
    )`
    
    if _, err := db.Exec(reactionsTable); err != nil {
        log.Printf("Error creating reflection_reactions table: %v", err)
        return err
    }
    log.Println("Created reflection_reactions table")

    // Create reaction_prompts table for configurable prompts
    promptsTable := `
    CREATE TABLE IF NOT EXISTS reaction_prompts (
        id SERIAL PRIMARY KEY,
        reaction_type VARCHAR(50) NOT NULL UNIQUE,
        prompt_text VARCHAR(200) NOT NULL,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    )`
    
    if _, err := db.Exec(promptsTable); err != nil {
        log.Printf("Error creating reaction_prompts table: %v", err)
        return err
    }
    log.Println("Created reaction_prompts table")

    // Insert default reaction prompts
    insertPrompts := `
    INSERT INTO reaction_prompts (reaction_type, prompt_text) VALUES
        ('ask_me_about_this', 'I''d love to hear more about this. What specifically would you like to discuss?'),
        ('similar_experience', 'I''ve been through something similar. What helped me was...'),
        ('update_me', 'I''m invested in your journey. How did this turn out?'),
        ('accountability_buddy', 'I''m here to help you stay on track. What''s your next step?'),
        ('different_angle', 'Have you considered looking at this from a different perspective?')
    ON CONFLICT (reaction_type) DO NOTHING
    `
    
    if _, err := db.Exec(insertPrompts); err != nil {
        log.Printf("Error inserting default prompts: %v", err)
        return err
    }
    log.Println("Inserted default reaction prompts")

    // Add favorite reaction type to the CHECK constraint
    alterReactionsTable := `
    ALTER TABLE reflection_reactions DROP CONSTRAINT IF EXISTS reflection_reactions_reaction_type_check;
    ALTER TABLE reflection_reactions ADD CONSTRAINT reflection_reactions_reaction_type_check 
    CHECK (reaction_type IN (
        'ask_me_about_this',
        'similar_experience', 
        'update_me',
        'accountability_buddy',
        'different_angle',
        'favorite'
    ))
    `
    
    if _, err := db.Exec(alterReactionsTable); err != nil {
        log.Printf("Error updating reflection_reactions table constraint: %v", err)
        return err
    }
    log.Println("Updated reflection_reactions table to include favorite reaction type")

    return nil
}

