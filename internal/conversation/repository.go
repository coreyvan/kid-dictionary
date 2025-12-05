package conversation

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AgeBracket represents the target age group for explanations.
type AgeBracket int

const (
	AgeBracketUnspecified  AgeBracket = 0
	AgeBracketLittleOnes   AgeBracket = 1 // 0-5 years
	AgeBracketGrowingMinds AgeBracket = 2 // 5-10 years
	AgeBracketPreTeens     AgeBracket = 3 // 10+ years
)

// Conversation represents a chat session.
type Conversation struct {
	ID         uuid.UUID
	UserID     *uuid.UUID // nullable for anonymous MVP
	Title      string
	AgeBracket AgeBracket
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Repository defines the interface for conversation persistence.
type Repository interface {
	// Create stores a new conversation and returns it with generated ID.
	Create(ctx context.Context, conv *Conversation) error

	// GetByID retrieves a conversation by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Conversation, error)

	// List returns conversations, optionally filtered by user ID.
	// Results are ordered by updated_at descending.
	List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*Conversation, error)

	// Update modifies an existing conversation.
	Update(ctx context.Context, conv *Conversation) error

	// Delete removes a conversation and all its messages.
	Delete(ctx context.Context, id uuid.UUID) error
}
