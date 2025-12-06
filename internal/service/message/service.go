package message

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/llm"
)

const (
	maxContentLength   = 500
	maxContextMessages = 10
	llmTimeout         = 60 * time.Second
)

// Service handles message business logic.
type Service struct {
	messageRepo      domain.MessageRepository
	conversationRepo domain.ConversationRepository
	llmProvider      llm.Provider
	classifier       *Classifier
}

// NewService creates a new message service.
func NewService(messageRepo domain.MessageRepository, conversationRepo domain.ConversationRepository, llmProvider llm.Provider) *Service {
	return &Service{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		llmProvider:      llmProvider,
		classifier:       NewClassifier(),
	}
}

// Response messages for different content tiers.
const (
	sensitiveGuidancePrefix = "This is a sensitive topic. Before sharing this explanation with your child, you may want to consider the timing and your child's emotional readiness.\n\n"

	contextualFramingRequest = "This topic can be approached from different perspectives. I'd like to help you explain this in a way that aligns with your family's values. Please let me know how you'd like me to frame this explanation, and I'll provide an age-appropriate response."

	offPurposeDecline = "I'm here to help parents explain concepts and words to their children in age-appropriate ways. I can help with questions like 'How do I explain where babies come from?' or 'What is gravity?'\n\nWhat concept or word would you like help explaining to your child?"
)

// SendMessage sends a user message and generates an AI response.
func (s *Service) SendMessage(ctx context.Context, conversationID uuid.UUID, content string) (*domain.SendMessageResult, error) {
	// Validate input
	if content == "" {
		return nil, domain.ErrEmptyContent
	}
	if len(content) > maxContentLength {
		return nil, domain.ErrContentTooLong
	}

	// Get conversation to determine age bracket
	conv, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrConversationNotFound
		}
		return nil, err
	}

	// Create and save user message
	userMsg := &domain.Message{
		ConversationID: conversationID,
		Role:           domain.RoleUser,
		Content:        content,
	}
	if err := s.messageRepo.Create(ctx, userMsg); err != nil {
		return nil, err
	}

	// Classify the content tier
	contentTier := s.classifier.Classify(content)

	// Handle off-purpose requests (Tier 4 - Redirect)
	if contentTier == domain.ContentTierRedirect {
		assistantMsg := &domain.Message{
			ConversationID: conversationID,
			Role:           domain.RoleAssistant,
			Content:        offPurposeDecline,
			ContentTier:    domain.ContentTierRedirect,
		}
		if err := s.messageRepo.Create(ctx, assistantMsg); err != nil {
			return nil, err
		}
		return &domain.SendMessageResult{
			UserMessage:      userMsg,
			AssistantMessage: assistantMsg,
		}, nil
	}

	// Handle contextual topics (Tier 3) - request framing preference
	if contentTier == domain.ContentTierContextual {
		assistantMsg := &domain.Message{
			ConversationID: conversationID,
			Role:           domain.RoleAssistant,
			Content:        contextualFramingRequest,
			ContentTier:    domain.ContentTierContextual,
		}
		if err := s.messageRepo.Create(ctx, assistantMsg); err != nil {
			return nil, err
		}
		return &domain.SendMessageResult{
			UserMessage:      userMsg,
			AssistantMessage: assistantMsg,
		}, nil
	}

	// Get conversation history for context
	history, err := s.messageRepo.GetRecentByConversationID(ctx, conversationID, maxContextMessages)
	if err != nil {
		return nil, err
	}

	// Build LLM messages from history
	llmMessages := make([]llm.Message, 0, len(history))
	for _, msg := range history {
		role := "user"
		if msg.Role == domain.RoleAssistant {
			role = "assistant"
		}
		llmMessages = append(llmMessages, llm.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Map conversation age bracket to LLM age bracket
	llmAgeBracket := llm.AgeBracket(conv.AgeBracket)

	// Call LLM for response with timeout
	llmCtx, cancel := context.WithTimeout(ctx, llmTimeout)
	defer cancel()

	resp, err := s.llmProvider.Complete(llmCtx, llm.CompletionRequest{
		AgeBracket: llmAgeBracket,
		Messages:   llmMessages,
	})
	if err != nil {
		return nil, err
	}

	// Apply soft guidance prefix for sensitive topics (Tier 2)
	responseContent := resp.Content
	if contentTier == domain.ContentTierSensitive {
		responseContent = sensitiveGuidancePrefix + responseContent
	}

	// Create and save assistant message
	assistantMsg := &domain.Message{
		ConversationID: conversationID,
		Role:           domain.RoleAssistant,
		Content:        responseContent,
		ContentTier:    contentTier,
	}
	if err := s.messageRepo.Create(ctx, assistantMsg); err != nil {
		return nil, err
	}

	return &domain.SendMessageResult{
		UserMessage:      userMsg,
		AssistantMessage: assistantMsg,
	}, nil
}
