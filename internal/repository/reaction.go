package repository

import (
    "database/sql"
    "fmt"
    "reflection-app/internal/models"
)

type ReactionRepository struct {
    db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
    return &ReactionRepository{db: db}
}

// AddReaction adds a reaction to a reflection
func (r *ReactionRepository) AddReaction(reflectionID, userID int, reactionType models.ReactionType, commentText string) error {
    // Check if user already reacted with this type
    var existingID int
    err := r.db.QueryRow(`
        SELECT id FROM reflection_reactions 
        WHERE reflection_id = $1 AND user_id = $2 AND reaction_type = $3
    `, reflectionID, userID, reactionType).Scan(&existingID)
    
    if err == nil {
        return fmt.Errorf("user already reacted with this type")
    }
    
    if err != sql.ErrNoRows {
        return err
    }
    
    // Insert new reaction
    _, err = r.db.Exec(`
        INSERT INTO reflection_reactions (reflection_id, user_id, reaction_type, comment_text) 
        VALUES ($1, $2, $3, $4)
    `, reflectionID, userID, reactionType, commentText)
    
    return err
}

// RemoveReaction removes a reaction from a reflection
func (r *ReactionRepository) RemoveReaction(reflectionID, userID int, reactionType models.ReactionType) error {
    result, err := r.db.Exec(`
        DELETE FROM reflection_reactions 
        WHERE reflection_id = $1 AND user_id = $2 AND reaction_type = $3
    `, reflectionID, userID, reactionType)
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("reaction not found")
    }
    
    return nil
}

// GetReactionsForReflection gets all reactions for a reflection
func (r *ReactionRepository) GetReactionsForReflection(reflectionID int) ([]models.ReactionResponse, error) {
    rows, err := r.db.Query(`
        SELECT rr.id, rr.reflection_id, rr.user_id, rr.reaction_type, rr.comment_text, rr.created_at,
               u.username, u.first_name, u.last_name
        FROM reflection_reactions rr
        JOIN users u ON rr.user_id = u.id
        WHERE rr.reflection_id = $1
        ORDER BY rr.created_at ASC
    `, reflectionID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var reactions []models.ReactionResponse
    for rows.Next() {
        var reaction models.ReactionResponse
        var firstName, lastName string
        err := rows.Scan(&reaction.ID, &reaction.ReflectionID, &reaction.UserID, &reaction.ReactionType,
                        &reaction.CommentText, &reaction.CreatedAt, &reaction.Username, &firstName, &lastName)
        if err != nil {
            return nil, err
        }
        
        reaction.UserName = firstName + " " + lastName
        reactions = append(reactions, reaction)
    }
    
    return reactions, nil
}

// GetReactionCountsForReflection gets reaction counts by type for a reflection
func (r *ReactionRepository) GetReactionCountsForReflection(reflectionID int) (map[models.ReactionType]int, error) {
    rows, err := r.db.Query(`
        SELECT reaction_type, COUNT(*) as count
        FROM reflection_reactions
        WHERE reflection_id = $1
        GROUP BY reaction_type
    `, reflectionID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    counts := make(map[models.ReactionType]int)
    for rows.Next() {
        var reactionType models.ReactionType
        var count int
        err := rows.Scan(&reactionType, &count)
        if err != nil {
            return nil, err
        }
        counts[reactionType] = count
    }
    
    return counts, nil
}

// GetUserReactionForReflection gets a user's reaction for a specific reflection
func (r *ReactionRepository) GetUserReactionForReflection(reflectionID, userID int) (*models.ReactionResponse, error) {
    var reaction models.ReactionResponse
    var firstName, lastName string
    err := r.db.QueryRow(`
        SELECT rr.id, rr.reflection_id, rr.user_id, rr.reaction_type, rr.comment_text, rr.created_at,
               u.username, u.first_name, u.last_name
        FROM reflection_reactions rr
        JOIN users u ON rr.user_id = u.id
        WHERE rr.reflection_id = $1 AND rr.user_id = $2
    `, reflectionID, userID).Scan(&reaction.ID, &reaction.ReflectionID, &reaction.UserID, &reaction.ReactionType,
                                   &reaction.CommentText, &reaction.CreatedAt, &reaction.Username, &firstName, &lastName)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    if err != nil {
        return nil, err
    }
    
    reaction.UserName = firstName + " " + lastName
    return &reaction, nil
}

// GetReactionPrompts gets all active reaction prompts
func (r *ReactionRepository) GetReactionPrompts() ([]models.ReactionPrompt, error) {
    rows, err := r.db.Query(`
        SELECT id, reaction_type, prompt_text, is_active, created_at, updated_at
        FROM reaction_prompts
        WHERE is_active = true
        ORDER BY reaction_type
    `)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var prompts []models.ReactionPrompt
    for rows.Next() {
        var prompt models.ReactionPrompt
        err := rows.Scan(&prompt.ID, &prompt.ReactionType, &prompt.PromptText, &prompt.IsActive,
                        &prompt.CreatedAt, &prompt.UpdatedAt)
        if err != nil {
            return nil, err
        }
        prompts = append(prompts, prompt)
    }
    
    return prompts, nil
} 