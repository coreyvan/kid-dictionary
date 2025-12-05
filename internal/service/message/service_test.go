package message_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/llm"
	"github.com/coreyvan/kid-dictionary/internal/service/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMessageRepo implements domain.MessageRepository for testing.
type mockMessageRepo struct {
	createFn                    func(ctx context.Context, msg *domain.Message) error
	getByConversationIDFn       func(ctx context.Context, conversationID uuid.UUID) ([]*domain.Message, error)
	getRecentByConversationIDFn func(ctx context.Context, conversationID uuid.UUID, limit int) ([]*domain.Message, error)
}

func (m *mockMessageRepo) Create(ctx context.Context, msg *domain.Message) error {
	if m.createFn != nil {
		return m.createFn(ctx, msg)
	}
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()
	return nil
}

func (m *mockMessageRepo) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*domain.Message, error) {
	if m.getByConversationIDFn != nil {
		return m.getByConversationIDFn(ctx, conversationID)
	}
	return nil, nil
}

func (m *mockMessageRepo) GetRecentByConversationID(ctx context.Context, conversationID uuid.UUID, limit int) ([]*domain.Message, error) {
	if m.getRecentByConversationIDFn != nil {
		return m.getRecentByConversationIDFn(ctx, conversationID, limit)
	}
	return nil, nil
}

// mockConversationRepo implements domain.ConversationRepository for testing.
type mockConversationRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error)
}

func (m *mockConversationRepo) Create(ctx context.Context, conv *domain.Conversation) error {
	return nil
}

func (m *mockConversationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return &domain.Conversation{
		ID:         id,
		AgeBracket: domain.AgeBracketLittleOnes,
	}, nil
}

func (m *mockConversationRepo) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*domain.Conversation, error) {
	return nil, nil
}

func (m *mockConversationRepo) Update(ctx context.Context, conv *domain.Conversation) error {
	return nil
}

func (m *mockConversationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// mockLLMProvider implements llm.Provider for testing.
type mockLLMProvider struct {
	completeFn func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error)
}

func (m *mockLLMProvider) Complete(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
	if m.completeFn != nil {
		return m.completeFn(ctx, req)
	}
	return llm.CompletionResponse{
		Content:    "The sky is blue because of light!",
		TokensUsed: 10,
	}, nil
}

func TestService_SendMessage_FollowUpContext(t *testing.T) {
	convID := uuid.New()

	// Create history with previous messages
	previousMessages := []*domain.Message{
		{ID: uuid.New(), ConversationID: convID, Role: domain.RoleUser, Content: "Why is the sky blue?"},
		{ID: uuid.New(), ConversationID: convID, Role: domain.RoleAssistant, Content: "Light scatters!"},
	}

	msgRepo := &mockMessageRepo{
		getRecentByConversationIDFn: func(ctx context.Context, conversationID uuid.UUID, limit int) ([]*domain.Message, error) {
			// Should include the new message plus history
			return append(previousMessages, &domain.Message{
				ID:             uuid.New(),
				ConversationID: convID,
				Role:           domain.RoleUser,
				Content:        "Tell me more about that",
			}), nil
		},
	}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         convID,
				AgeBracket: domain.AgeBracketGrowingMinds,
			}, nil
		},
	}

	var capturedMessages []llm.Message
	llmProvider := &mockLLMProvider{
		completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			capturedMessages = req.Messages
			return llm.CompletionResponse{
				Content:    "More details about light scattering...",
				TokensUsed: 20,
			}, nil
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider)
	result, err := svc.SendMessage(context.Background(), convID, "Tell me more about that")

	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify context was passed to LLM (should have history + new message)
	assert.GreaterOrEqual(t, len(capturedMessages), 2, "LLM should receive conversation history")
}

func TestService_SendMessage_AgeBracketVariation(t *testing.T) {
	testCases := []struct {
		name           string
		ageBracket     domain.AgeBracket
		expectedLLMAge llm.AgeBracket
	}{
		{
			name:           "little ones",
			ageBracket:     domain.AgeBracketLittleOnes,
			expectedLLMAge: llm.AgeBracketLittleOnes,
		},
		{
			name:           "growing minds",
			ageBracket:     domain.AgeBracketGrowingMinds,
			expectedLLMAge: llm.AgeBracketGrowingMinds,
		},
		{
			name:           "pre-teens",
			ageBracket:     domain.AgeBracketPreTeens,
			expectedLLMAge: llm.AgeBracketPreTeens,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			convID := uuid.New()
			msgRepo := &mockMessageRepo{}
			convRepo := &mockConversationRepo{
				getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
					return &domain.Conversation{
						ID:         convID,
						AgeBracket: tc.ageBracket,
					}, nil
				},
			}

			var capturedAgeBracket llm.AgeBracket
			llmProvider := &mockLLMProvider{
				completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
					capturedAgeBracket = req.AgeBracket
					return llm.CompletionResponse{
						Content:    "Response for " + tc.name,
						TokensUsed: 10,
					}, nil
				},
			}

			svc := message.NewService(msgRepo, convRepo, llmProvider)
			_, err := svc.SendMessage(context.Background(), convID, "Why is the sky blue?")

			require.NoError(t, err)
			assert.Equal(t, tc.expectedLLMAge, capturedAgeBracket, "LLM should receive correct age bracket")
		})
	}
}

func TestService_SendMessage_SensitiveTopic(t *testing.T) {
	convID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         convID,
				AgeBracket: domain.AgeBracketLittleOnes,
			}, nil
		},
	}
	llmProvider := &mockLLMProvider{
		completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			return llm.CompletionResponse{
				Content:    "Death is when someone's body stops working.",
				TokensUsed: 10,
			}, nil
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider)
	result, err := svc.SendMessage(context.Background(), convID, "How do I explain death to my child?")

	require.NoError(t, err)
	assert.Equal(t, domain.ContentTierSensitive, result.AssistantMessage.ContentTier)
	assert.Contains(t, result.AssistantMessage.Content, "This is a sensitive topic")
	assert.Contains(t, result.AssistantMessage.Content, "Death is when someone's body stops working.")
}

func TestService_SendMessage_ContextualTopic(t *testing.T) {
	convID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         convID,
				AgeBracket: domain.AgeBracketGrowingMinds,
			}, nil
		},
	}
	llmProvider := &mockLLMProvider{}

	svc := message.NewService(msgRepo, convRepo, llmProvider)
	result, err := svc.SendMessage(context.Background(), convID, "What is religion?")

	require.NoError(t, err)
	assert.Equal(t, domain.ContentTierContextual, result.AssistantMessage.ContentTier)
	assert.Contains(t, result.AssistantMessage.Content, "frame this explanation")
}

func TestService_SendMessage_OffPurposeRequest(t *testing.T) {
	convID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         convID,
				AgeBracket: domain.AgeBracketPreTeens,
			}, nil
		},
	}
	llmProvider := &mockLLMProvider{}

	svc := message.NewService(msgRepo, convRepo, llmProvider)
	result, err := svc.SendMessage(context.Background(), convID, "Write me a poem about cats")

	require.NoError(t, err)
	assert.Equal(t, domain.ContentTierRedirect, result.AssistantMessage.ContentTier)
	assert.Contains(t, result.AssistantMessage.Content, "I'm here to help parents explain")
}

func TestService_SendMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		convID := uuid.New()
		msgRepo := &mockMessageRepo{}
		convRepo := &mockConversationRepo{
			getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
				return &domain.Conversation{
					ID:         convID,
					AgeBracket: domain.AgeBracketLittleOnes,
				}, nil
			},
		}
		llmProvider := &mockLLMProvider{
			completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
				assert.Equal(t, llm.AgeBracketLittleOnes, req.AgeBracket)
				return llm.CompletionResponse{
					Content:    "The sky is blue because of light!",
					TokensUsed: 10,
				}, nil
			},
		}

		svc := message.NewService(msgRepo, convRepo, llmProvider)
		result, err := svc.SendMessage(context.Background(), convID, "Why is the sky blue?")

		require.NoError(t, err)
		assert.NotNil(t, result.UserMessage)
		assert.NotNil(t, result.AssistantMessage)
		assert.Equal(t, "Why is the sky blue?", result.UserMessage.Content)
		assert.Equal(t, domain.RoleUser, result.UserMessage.Role)
		assert.Equal(t, "The sky is blue because of light!", result.AssistantMessage.Content)
		assert.Equal(t, domain.RoleAssistant, result.AssistantMessage.Role)
	})

	t.Run("error with empty content", func(t *testing.T) {
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, &mockLLMProvider{})

		_, err := svc.SendMessage(context.Background(), uuid.New(), "")

		assert.ErrorIs(t, err, domain.ErrEmptyContent)
	})

	t.Run("error with content too long", func(t *testing.T) {
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, &mockLLMProvider{})

		longContent := make([]byte, 501)
		for i := range longContent {
			longContent[i] = 'a'
		}

		_, err := svc.SendMessage(context.Background(), uuid.New(), string(longContent))

		assert.ErrorIs(t, err, domain.ErrContentTooLong)
	})

	t.Run("error when conversation not found", func(t *testing.T) {
		convRepo := &mockConversationRepo{
			getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
				return nil, domain.ErrNotFound
			},
		}
		svc := message.NewService(&mockMessageRepo{}, convRepo, &mockLLMProvider{})

		_, err := svc.SendMessage(context.Background(), uuid.New(), "Hello")

		assert.ErrorIs(t, err, domain.ErrConversationNotFound)
	})

	t.Run("error when LLM fails", func(t *testing.T) {
		llmErr := errors.New("LLM error")
		llmProvider := &mockLLMProvider{
			completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
				return llm.CompletionResponse{}, llmErr
			},
		}
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, llmProvider)

		_, err := svc.SendMessage(context.Background(), uuid.New(), "Hello")

		assert.Error(t, err)
	})
}
