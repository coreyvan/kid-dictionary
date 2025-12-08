package message

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/llm"
)

const (
	maxContentLength   = 500
	maxContextMessages = 10
	llmTimeout         = 60 * time.Second
	titleTimeout       = 10 * time.Second
	fallbackTitle      = "New Conversation"
)

// Service handles message business logic.
type Service struct {
	logger           *slog.Logger
	messageRepo      domain.MessageRepository
	conversationRepo domain.ConversationRepository
	llmProvider      llm.Provider
	classifier       *Classifier
}

// NewService creates a new message service.
func NewService(messageRepo domain.MessageRepository, conversationRepo domain.ConversationRepository, llmProvider llm.Provider, l *slog.Logger) *Service {
	return &Service{
		logger:           l,
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
// If conversationID is uuid.Nil, a new conversation is auto-created using ageBracket.
// The ageBracket parameter is only used when conversationID is nil; it is ignored otherwise.
func (s *Service) SendMessage(ctx context.Context, conversationID uuid.UUID, content string, ageBracket *domain.AgeBracket) (*domain.SendMessageResult, error) {
	// Validate input
	if content == "" {
		return nil, domain.ErrEmptyContent
	}
	if len(content) > maxContentLength {
		return nil, domain.ErrContentTooLong
	}

	s.logger.Debug("SendMessage", "content", content, "conversationID", conversationID)

	var conv *domain.Conversation
	var autoCreated bool

	// Check if we need to auto-create a conversation
	if conversationID == uuid.Nil {
		// Validate age bracket is provided for new conversations
		if ageBracket == nil || *ageBracket == domain.AgeBracketUnspecified {
			return nil, domain.ErrAgeBracketRequired
		}
		if *ageBracket < domain.AgeBracketLittleOnes || *ageBracket > domain.AgeBracketPreTeens {
			return nil, domain.ErrInvalidAgeBracket
		}

		// Generate title via LLM with fallback
		title := s.generateTitleWithFallback(ctx, content)

		// Create the conversation
		conv = &domain.Conversation{
			Title:      title,
			AgeBracket: *ageBracket,
		}
		if err := s.conversationRepo.Create(ctx, conv); err != nil {
			return nil, err
		}
		conversationID = conv.ID
		autoCreated = true
		s.logger.Info("auto-created conversation", "conversation_id", conv.ID, "title", conv.Title)
	} else {
		// Get existing conversation to determine age bracket
		var err error
		conv, err = s.conversationRepo.GetByID(ctx, conversationID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrConversationNotFound
			}
			return nil, err
		}
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

	result := &domain.SendMessageResult{
		UserMessage:      userMsg,
		AssistantMessage: assistantMsg,
	}

	// Include conversation in result if it was auto-created
	if autoCreated {
		result.Conversation = conv
	}

	return result, nil
}

// generateTitleWithFallback attempts to generate a title via LLM.
// On any error, it logs the failure and returns the fallback title.
func (s *Service) generateTitleWithFallback(ctx context.Context, content string) string {
	start := time.Now()

	titleCtx, cancel := context.WithTimeout(ctx, titleTimeout)
	defer cancel()

	title, err := s.llmProvider.GenerateTitle(titleCtx, content)
	latency := time.Since(start)

	if err != nil {
		s.logger.Warn("title_generation",
			"event", "title_generation",
			"success", false,
			"latency_ms", latency.Milliseconds(),
			"error", err.Error(),
		)
		return fallbackTitle
	}

	s.logger.Info("title_generation",
		"event", "title_generation",
		"success", true,
		"latency_ms", latency.Milliseconds(),
		"title", title,
	)

	return title
}
