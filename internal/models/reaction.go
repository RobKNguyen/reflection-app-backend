package models

import (
    "time"
)

type ReactionType string

const (
    AskMeAboutThis      ReactionType = "ask_me_about_this"
    SimilarExperience    ReactionType = "similar_experience"
    UpdateMe            ReactionType = "update_me"
    AccountabilityBuddy ReactionType = "accountability_buddy"
    DifferentAngle      ReactionType = "different_angle"
    Favorite            ReactionType = "favorite"
)

type ReflectionReaction struct {
    ID           int          `json:"id"`
    ReflectionID int          `json:"reflection_id"`
    UserID       int          `json:"user_id"`
    ReactionType ReactionType `json:"reaction_type"`
    CommentText  string       `json:"comment_text,omitempty"`
    CreatedAt    time.Time    `json:"created_at"`
}

type ReactionRequest struct {
    ReactionType ReactionType `json:"reaction_type"`
    CommentText  string       `json:"comment_text,omitempty"`
}

type ReactionResponse struct {
    ID           int          `json:"id"`
    ReflectionID int          `json:"reflection_id"`
    UserID       int          `json:"user_id"`
    ReactionType ReactionType `json:"reaction_type"`
    CommentText  string       `json:"comment_text,omitempty"`
    CreatedAt    time.Time    `json:"created_at"`
    // Additional fields for UI
    Username string `json:"username,omitempty"`
    UserName string `json:"user_name,omitempty"`
}

type ReactionPrompt struct {
    ID           int          `json:"id"`
    ReactionType ReactionType `json:"reaction_type"`
    PromptText   string       `json:"prompt_text"`
    IsActive     bool         `json:"is_active"`
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`
}

type ReflectionWithReactions struct {
    Reflection
    Reactions []ReactionResponse `json:"reactions"`
    ReactionCounts map[ReactionType]int `json:"reaction_counts"`
} 