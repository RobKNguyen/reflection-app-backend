// Add this to your database package (internal/database/database.go)
package database

import (
    "database/sql"
    "log"
)

// DropAllTables drops all tables using CASCADE to handle foreign key constraints
func DropAllTables(db *sql.DB) error {
    // For PostgreSQL, we can use CASCADE to automatically drop dependent objects
    tables := []string{
        "reflection_reactions",
        "reaction_prompts",
        "friendships",
        "reflection_tracking",
        "actions",
        "reflections", 
        "sub_categories",
        "categories",
        "users",
    }
    
    for _, table := range tables {
        query := `DROP TABLE IF EXISTS ` + table + ` CASCADE`
        if _, err := db.Exec(query); err != nil {
            log.Printf("Error dropping table %s: %v", table, err)
            return err
        }
        log.Printf("Dropped table: %s", table)
    }
    
    return nil
}