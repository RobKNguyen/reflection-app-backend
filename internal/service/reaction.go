package service

import (
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type ReactionService struct {
    reactionRepo *repository.ReactionRepository
}

func NewReactionService(reactionRepo *repository.ReactionRepository) *ReactionService {
    return &ReactionService{
        reactionRepo: reactionRepo,
    }
}

// AddReaction adds a reaction to a reflection
func (s *ReactionService) AddReaction(reflectionID, userID int, reactionType models.ReactionType, commentText string) error {
    // Validate reaction type
    if !isValidReactionType(reactionType) {
        return models.ErrInvalidRequest
    }
    
    // Validate comment length
    if len(commentText) > 100 {
        return models.ErrInvalidRequest
    }
    
    return s.reactionRepo.AddReaction(reflectionID, userID, reactionType, commentText)
}

// RemoveReaction removes a reaction from a reflection
func (s *ReactionService) RemoveReaction(reflectionID, userID int, reactionType models.ReactionType) error {
    return s.reactionRepo.RemoveReaction(reflectionID, userID, reactionType)
}

// GetReactionsForReflection gets all reactions for a reflection
func (s *ReactionService) GetReactionsForReflection(reflectionID int) ([]models.ReactionResponse, error) {
    return s.reactionRepo.GetReactionsForReflection(reflectionID)
}

// GetReactionCountsForReflection gets reaction counts by type for a reflection
func (s *ReactionService) GetReactionCountsForReflection(reflectionID int) (map[models.ReactionType]int, error) {
    return s.reactionRepo.GetReactionCountsForReflection(reflectionID)
}

// GetUserReactionForReflection gets a user's reaction for a specific reflection
func (s *ReactionService) GetUserReactionForReflection(reflectionID, userID int) (*models.ReactionResponse, error) {
    return s.reactionRepo.GetUserReactionForReflection(reflectionID, userID)
}

// GetReactionPrompts gets all active reaction prompts
func (s *ReactionService) GetReactionPrompts() ([]models.ReactionPrompt, error) {
    return s.reactionRepo.GetReactionPrompts()
}

// isValidReactionType validates if a reaction type is valid
func isValidReactionType(reactionType models.ReactionType) bool {
    validTypes := []models.ReactionType{
        models.AskMeAboutThis,
        models.SimilarExperience,
        models.UpdateMe,
        models.AccountabilityBuddy,
        models.DifferentAngle,
        models.Favorite,
    }
    
    for _, validType := range validTypes {
        if reactionType == validType {
            return true
        }
    }
    
    return false
} 