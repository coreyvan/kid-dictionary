package message_test

import (
	"context"
	"errors"
	"log/slog"
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
	createFn  func(ctx context.Context, conv *domain.Conversation) error
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error)
}

func (m *mockConversationRepo) Create(ctx context.Context, conv *domain.Conversation) error {
	if m.createFn != nil {
		return m.createFn(ctx, conv)
	}
	conv.ID = uuid.New()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
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
	completeFn      func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error)
	generateTitleFn func(ctx context.Context, content string) (string, error)
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

func (m *mockLLMProvider) GenerateTitle(ctx context.Context, content string) (string, error) {
	if m.generateTitleFn != nil {
		return m.generateTitleFn(ctx, content)
	}
	return "Generated Title", nil
}

// Helper to get pointer to age bracket
func ageBracketPtr(ab domain.AgeBracket) *domain.AgeBracket {
	return &ab
}

// ============================================================================
// User Story 1 Tests: Auto-Create Conversation
// ============================================================================

func TestSendMessage_AutoCreatesConversation(t *testing.T) {
	msgRepo := &mockMessageRepo{}
	var capturedConv *domain.Conversation
	convRepo := &mockConversationRepo{
		createFn: func(ctx context.Context, conv *domain.Conversation) error {
			capturedConv = conv
			conv.ID = uuid.New()
			conv.CreatedAt = time.Now()
			conv.UpdatedAt = time.Now()
			return nil
		},
	}
	llmProvider := &mockLLMProvider{
		generateTitleFn: func(ctx context.Context, content string) (string, error) {
			return "Sky Being Blue", nil
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), uuid.Nil, "Why is the sky blue?", ageBracketPtr(domain.AgeBracketLittleOnes))

	require.NoError(t, err)
	require.NotNil(t, result.Conversation)
	assert.Equal(t, "Sky Being Blue", capturedConv.Title)
	assert.Equal(t, domain.AgeBracketLittleOnes, capturedConv.AgeBracket)
	assert.NotEqual(t, uuid.Nil, result.Conversation.ID)
}

func TestSendMessage_TitleGenerationFallback(t *testing.T) {
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{}
	llmProvider := &mockLLMProvider{
		generateTitleFn: func(ctx context.Context, content string) (string, error) {
			return "", errors.New("LLM unavailable")
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), uuid.Nil, "Why is the sky blue?", ageBracketPtr(domain.AgeBracketLittleOnes))

	require.NoError(t, err)
	require.NotNil(t, result.Conversation)
	assert.Equal(t, "New Conversation", result.Conversation.Title)
}

func TestSendMessage_RequiresAgeBracketWhenNoConversationID(t *testing.T) {
	svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, &mockLLMProvider{}, slog.Default())

	// Test with nil age bracket
	_, err := svc.SendMessage(context.Background(), uuid.Nil, "Why is the sky blue?", nil)
	assert.ErrorIs(t, err, domain.ErrAgeBracketRequired)

	// Test with unspecified age bracket
	_, err = svc.SendMessage(context.Background(), uuid.Nil, "Why is the sky blue?", ageBracketPtr(domain.AgeBracketUnspecified))
	assert.ErrorIs(t, err, domain.ErrAgeBracketRequired)
}

// ============================================================================
// User Story 2 Tests: Follow-up Messages
// ============================================================================

func TestSendMessage_FollowUpWithAutoCreatedConversation(t *testing.T) {
	autoCreatedConvID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		createFn: func(ctx context.Context, conv *domain.Conversation) error {
			conv.ID = autoCreatedConvID
			conv.CreatedAt = time.Now()
			conv.UpdatedAt = time.Now()
			return nil
		},
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			if id == autoCreatedConvID {
				return &domain.Conversation{
					ID:         autoCreatedConvID,
					Title:      "Sky Being Blue",
					AgeBracket: domain.AgeBracketLittleOnes,
				}, nil
			}
			return nil, domain.ErrNotFound
		},
	}
	llmProvider := &mockLLMProvider{
		generateTitleFn: func(ctx context.Context, content string) (string, error) {
			return "Sky Being Blue", nil
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())

	// First message auto-creates conversation
	result1, err := svc.SendMessage(context.Background(), uuid.Nil, "Why is the sky blue?", ageBracketPtr(domain.AgeBracketLittleOnes))
	require.NoError(t, err)
	require.NotNil(t, result1.Conversation)

	// Follow-up message uses the returned conversation ID
	result2, err := svc.SendMessage(context.Background(), result1.Conversation.ID, "Tell me more", nil)
	require.NoError(t, err)
	assert.Nil(t, result2.Conversation) // Should not return conversation for existing conv
	assert.Equal(t, result1.Conversation.ID, result2.UserMessage.ConversationID)
}

// ============================================================================
// User Story 3 Tests: Backwards Compatibility
// ============================================================================

func TestSendMessage_ExistingConversation_NoAutoCreate(t *testing.T) {
	existingConvID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         existingConvID,
				Title:      "Existing Conversation",
				AgeBracket: domain.AgeBracketGrowingMinds,
			}, nil
		},
	}
	llmProvider := &mockLLMProvider{}

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), existingConvID, "Why is the sky blue?", nil)

	require.NoError(t, err)
	assert.Nil(t, result.Conversation) // No conversation returned for existing conv
	assert.NotNil(t, result.UserMessage)
	assert.NotNil(t, result.AssistantMessage)
}

func TestSendMessage_ExistingConversation_IgnoresAgeBracket(t *testing.T) {
	existingConvID := uuid.New()
	msgRepo := &mockMessageRepo{}
	convRepo := &mockConversationRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
			return &domain.Conversation{
				ID:         existingConvID,
				Title:      "Existing Conversation",
				AgeBracket: domain.AgeBracketGrowingMinds, // Conversation has GrowingMinds
			}, nil
		},
	}

	var capturedAgeBracket llm.AgeBracket
	llmProvider := &mockLLMProvider{
		completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			capturedAgeBracket = req.AgeBracket
			return llm.CompletionResponse{Content: "Response", TokensUsed: 10}, nil
		},
	}

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())

	// Pass PreTeens age bracket but conversation is GrowingMinds
	_, err := svc.SendMessage(context.Background(), existingConvID, "Why is the sky blue?", ageBracketPtr(domain.AgeBracketPreTeens))

	require.NoError(t, err)
	// Should use conversation's age bracket (GrowingMinds), not the passed one (PreTeens)
	assert.Equal(t, llm.AgeBracketGrowingMinds, capturedAgeBracket)
}

// ============================================================================
// Existing Tests (Updated for new signature)
// ============================================================================

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

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), convID, "Tell me more about that", nil)

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

			svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
			_, err := svc.SendMessage(context.Background(), convID, "Why is the sky blue?", nil)

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

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), convID, "How do I explain death to my child?", nil)

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

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), convID, "What is religion?", nil)

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

	svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
	result, err := svc.SendMessage(context.Background(), convID, "Write me a poem about cats", nil)

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

		svc := message.NewService(msgRepo, convRepo, llmProvider, slog.Default())
		result, err := svc.SendMessage(context.Background(), convID, "Why is the sky blue?", nil)

		require.NoError(t, err)
		assert.NotNil(t, result.UserMessage)
		assert.NotNil(t, result.AssistantMessage)
		assert.Equal(t, "Why is the sky blue?", result.UserMessage.Content)
		assert.Equal(t, domain.RoleUser, result.UserMessage.Role)
		assert.Equal(t, "The sky is blue because of light!", result.AssistantMessage.Content)
		assert.Equal(t, domain.RoleAssistant, result.AssistantMessage.Role)
	})

	t.Run("error with empty content", func(t *testing.T) {
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, &mockLLMProvider{}, slog.Default())

		_, err := svc.SendMessage(context.Background(), uuid.New(), "", nil)

		assert.ErrorIs(t, err, domain.ErrEmptyContent)
	})

	t.Run("error with content too long", func(t *testing.T) {
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, &mockLLMProvider{}, slog.Default())

		longContent := make([]byte, 501)
		for i := range longContent {
			longContent[i] = 'a'
		}

		_, err := svc.SendMessage(context.Background(), uuid.New(), string(longContent), nil)

		assert.ErrorIs(t, err, domain.ErrContentTooLong)
	})

	t.Run("error when conversation not found", func(t *testing.T) {
		convRepo := &mockConversationRepo{
			getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
				return nil, domain.ErrNotFound
			},
		}
		svc := message.NewService(&mockMessageRepo{}, convRepo, &mockLLMProvider{}, slog.Default())

		_, err := svc.SendMessage(context.Background(), uuid.New(), "Hello", nil)

		assert.ErrorIs(t, err, domain.ErrConversationNotFound)
	})

	t.Run("error when LLM fails", func(t *testing.T) {
		llmErr := errors.New("LLM error")
		llmProvider := &mockLLMProvider{
			completeFn: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
				return llm.CompletionResponse{}, llmErr
			},
		}
		svc := message.NewService(&mockMessageRepo{}, &mockConversationRepo{}, llmProvider, slog.Default())

		_, err := svc.SendMessage(context.Background(), uuid.New(), "Hello", nil)

		assert.Error(t, err)
	})
}
