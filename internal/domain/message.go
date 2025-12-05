package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Role indicates who sent a message in a conversation.
type Role int

const (
	RoleUnspecified Role = 0
	RoleUser        Role = 1
	RoleAssistant   Role = 2
)

// ContentTier classifies topic sensitivity for response handling.
type ContentTier int

const (
	ContentTierUnspecified ContentTier = 0
	ContentTierNormal      ContentTier = 1 // Standard topics - direct explanation
	ContentTierSensitive   ContentTier = 2 // Death, divorce, etc. - soft guidance prefix
	ContentTierContextual  ContentTier = 3 // Religion, politics - request framing preference
	ContentTierRedirect    ContentTier = 4 // Harmful or off-purpose - polite decline
)

// Message represents a single message in a conversation.
type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	Role           Role
	Content        string
	ContentTier    ContentTier // Only set for assistant messages
	CreatedAt      time.Time
}

// SendMessageResult contains both the user message and assistant response.
type SendMessageResult struct {
	UserMessage      *Message
	AssistantMessage *Message
}

// MessageRepository defines the interface for message persistence.
type MessageRepository interface {
	// Create stores a new message and returns it with generated ID.
	Create(ctx context.Context, msg *Message) error

	// GetByConversationID retrieves all messages for a conversation.
	// Results are ordered by created_at ascending.
	GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*Message, error)

	// GetRecentByConversationID retrieves the last N messages for context.
	// Results are ordered by created_at ascending.
	GetRecentByConversationID(ctx context.Context, conversationID uuid.UUID, limit int) ([]*Message, error)
}
