package conversation

import (
	"context"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/google/uuid"
)

// Service handles conversation business logic.
type Service struct {
	repo domain.ConversationRepository
}

// NewService creates a new conversation service.
func NewService(repo domain.ConversationRepository) *Service {
	return &Service{repo: repo}
}

// CreateConversation creates a new conversation with the given title and age bracket.
func (s *Service) CreateConversation(ctx context.Context, title string, ageBracket domain.AgeBracket) (*domain.Conversation, error) {
	if title == "" {
		title = "New Conversation"
	}

	if ageBracket == domain.AgeBracketUnspecified {
		return nil, domain.ErrInvalidAgeBracket
	}

	if ageBracket < domain.AgeBracketLittleOnes || ageBracket > domain.AgeBracketPreTeens {
		return nil, domain.ErrInvalidAgeBracket
	}

	conv := &domain.Conversation{
		Title:      title,
		AgeBracket: ageBracket,
	}

	if err := s.repo.Create(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// GetConversation retrieves a conversation by ID.
func (s *Service) GetConversation(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	conv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// UpdateConversation updates a conversation's title and/or age bracket.
func (s *Service) UpdateConversation(ctx context.Context, id uuid.UUID, title *string, ageBracket *domain.AgeBracket) (*domain.Conversation, error) {
	conv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		conv.Title = *title
	}

	if ageBracket != nil {
		if *ageBracket == domain.AgeBracketUnspecified {
			return nil, domain.ErrInvalidAgeBracket
		}
		if *ageBracket < domain.AgeBracketLittleOnes || *ageBracket > domain.AgeBracketPreTeens {
			return nil, domain.ErrInvalidAgeBracket
		}
		conv.AgeBracket = *ageBracket
	}

	if err := s.repo.Update(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// ListConversations returns a paginated list of conversations.
func (s *Service) ListConversations(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*domain.Conversation, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(ctx, userID, limit, offset)
}

// DeleteConversation removes a conversation and all its messages.
func (s *Service) DeleteConversation(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetMessages retrieves all messages for a conversation through the message repository.
// This is a convenience method for the transport layer.
func (s *Service) GetMessages(ctx context.Context, msgRepo domain.MessageRepository, conversationID uuid.UUID) ([]*domain.Message, error) {
	return msgRepo.GetByConversationID(ctx, conversationID)
}
