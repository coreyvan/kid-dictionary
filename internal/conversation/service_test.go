package conversation_test

import (
	"context"
	"testing"
	"time"

	"github.com/coreyvan/kid-dictionary/internal/conversation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepository is a simple mock for testing the service.
type mockRepository struct {
	createFn  func(ctx context.Context, conv *conversation.Conversation) error
	getByIDFn func(ctx context.Context, id uuid.UUID) (*conversation.Conversation, error)
	listFn    func(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*conversation.Conversation, error)
	updateFn  func(ctx context.Context, conv *conversation.Conversation) error
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockRepository) Create(ctx context.Context, conv *conversation.Conversation) error {
	if m.createFn != nil {
		return m.createFn(ctx, conv)
	}
	conv.ID = uuid.New()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*conversation.Conversation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, conversation.ErrNotFound
}

func (m *mockRepository) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*conversation.Conversation, error) {
	if m.listFn != nil {
		return m.listFn(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, conv *conversation.Conversation) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, conv)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return conversation.ErrNotFound
}

func TestService_CreateConversation(t *testing.T) {
	t.Run("success with title", func(t *testing.T) {
		repo := &mockRepository{}
		svc := conversation.NewService(repo)

		conv, err := svc.CreateConversation(context.Background(), "Test Title", conversation.AgeBracketLittleOnes)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, conv.ID)
		assert.Equal(t, "Test Title", conv.Title)
		assert.Equal(t, conversation.AgeBracketLittleOnes, conv.AgeBracket)
	})

	t.Run("success with empty title uses default", func(t *testing.T) {
		repo := &mockRepository{}
		svc := conversation.NewService(repo)

		conv, err := svc.CreateConversation(context.Background(), "", conversation.AgeBracketGrowingMinds)

		require.NoError(t, err)
		assert.Equal(t, "New Conversation", conv.Title)
	})

	t.Run("error with unspecified age bracket", func(t *testing.T) {
		repo := &mockRepository{}
		svc := conversation.NewService(repo)

		_, err := svc.CreateConversation(context.Background(), "Test", conversation.AgeBracketUnspecified)

		assert.ErrorIs(t, err, conversation.ErrInvalidAgeBracket)
	})

	t.Run("error with invalid age bracket", func(t *testing.T) {
		repo := &mockRepository{}
		svc := conversation.NewService(repo)

		_, err := svc.CreateConversation(context.Background(), "Test", conversation.AgeBracket(99))

		assert.ErrorIs(t, err, conversation.ErrInvalidAgeBracket)
	})
}

func TestService_GetConversation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		convID := uuid.New()
		repo := &mockRepository{
			getByIDFn: func(ctx context.Context, id uuid.UUID) (*conversation.Conversation, error) {
				return &conversation.Conversation{
					ID:         convID,
					Title:      "Test",
					AgeBracket: conversation.AgeBracketPreTeens,
				}, nil
			},
		}
		svc := conversation.NewService(repo)

		conv, err := svc.GetConversation(context.Background(), convID)

		require.NoError(t, err)
		assert.Equal(t, convID, conv.ID)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockRepository{}
		svc := conversation.NewService(repo)

		_, err := svc.GetConversation(context.Background(), uuid.New())

		assert.ErrorIs(t, err, conversation.ErrNotFound)
	})
}

func TestService_ListConversations(t *testing.T) {
	t.Run("success with default limit", func(t *testing.T) {
		repo := &mockRepository{
			listFn: func(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*conversation.Conversation, error) {
				assert.Equal(t, 20, limit)
				assert.Equal(t, 0, offset)
				return []*conversation.Conversation{
					{ID: uuid.New(), Title: "Conv 1"},
					{ID: uuid.New(), Title: "Conv 2"},
				}, nil
			},
		}
		svc := conversation.NewService(repo)

		convs, err := svc.ListConversations(context.Background(), nil, 0, 0)

		require.NoError(t, err)
		assert.Len(t, convs, 2)
	})

	t.Run("respects limit cap", func(t *testing.T) {
		repo := &mockRepository{
			listFn: func(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*conversation.Conversation, error) {
				assert.Equal(t, 100, limit) // Should be capped at 100
				return nil, nil
			},
		}
		svc := conversation.NewService(repo)

		_, err := svc.ListConversations(context.Background(), nil, 500, 0)

		require.NoError(t, err)
	})
}

func TestService_DeleteConversation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		convID := uuid.New()
		repo := &mockRepository{
			deleteFn: func(ctx context.Context, id uuid.UUID) error {
				assert.Equal(t, convID, id)
				return nil
			},
		}
		svc := conversation.NewService(repo)

		err := svc.DeleteConversation(context.Background(), convID)

		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockRepository{
			deleteFn: func(ctx context.Context, id uuid.UUID) error {
				return conversation.ErrNotFound
			},
		}
		svc := conversation.NewService(repo)

		err := svc.DeleteConversation(context.Background(), uuid.New())

		assert.ErrorIs(t, err, conversation.ErrNotFound)
	})
}
