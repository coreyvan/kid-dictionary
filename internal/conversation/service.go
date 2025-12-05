package conversation

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service errors.
var (
	ErrInvalidAgeBracket = errors.New("invalid age bracket")
	ErrTitleRequired     = errors.New("title is required")
)

// Service handles conversation business logic.
type Service struct {
	repo Repository
}

// NewService creates a new conversation service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateConversation creates a new conversation with the given title and age bracket.
func (s *Service) CreateConversation(ctx context.Context, title string, ageBracket AgeBracket) (*Conversation, error) {
	if title == "" {
		title = "New Conversation"
	}

	if ageBracket == AgeBracketUnspecified {
		return nil, ErrInvalidAgeBracket
	}

	if ageBracket < AgeBracketLittleOnes || ageBracket > AgeBracketPreTeens {
		return nil, ErrInvalidAgeBracket
	}

	conv := &Conversation{
		Title:      title,
		AgeBracket: ageBracket,
	}

	if err := s.repo.Create(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// GetConversation retrieves a conversation by ID.
func (s *Service) GetConversation(ctx context.Context, id uuid.UUID) (*Conversation, error) {
	conv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// UpdateConversation updates a conversation's title and/or age bracket.
func (s *Service) UpdateConversation(ctx context.Context, id uuid.UUID, title *string, ageBracket *AgeBracket) (*Conversation, error) {
	conv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		conv.Title = *title
	}

	if ageBracket != nil {
		if *ageBracket == AgeBracketUnspecified {
			return nil, ErrInvalidAgeBracket
		}
		if *ageBracket < AgeBracketLittleOnes || *ageBracket > AgeBracketPreTeens {
			return nil, ErrInvalidAgeBracket
		}
		conv.AgeBracket = *ageBracket
	}

	if err := s.repo.Update(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// ListConversations returns a paginated list of conversations.
func (s *Service) ListConversations(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*Conversation, error) {
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
