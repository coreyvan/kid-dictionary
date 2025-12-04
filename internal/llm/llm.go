package llm

import (
	"context"
	"errors"
)

// Domain errors for LLM operations.
var (
	ErrInvalidAPIKey       = errors.New("invalid API key")
	ErrRateLimited         = errors.New("rate limited")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrTimeout             = errors.New("request timeout")
	ErrInputTooLong        = errors.New("input exceeds maximum token limit")
	ErrUserRateLimited     = errors.New("user rate limit exceeded")
)

// AgeBracket represents the target age group for explanations.
type AgeBracket int

const (
	AgeBracketUnspecified  AgeBracket = iota
	AgeBracketLittleOnes              // 0-5 years
	AgeBracketGrowingMinds            // 5-10 years
	AgeBracketPreTeens                // 10+ years
)

// Message represents a single message in a conversation.
type Message struct {
	Role    string // "user" or "assistant"
	Content string
}

// CompletionRequest contains the data needed to generate a completion.
type CompletionRequest struct {
	UserID     string     // Used for per-user rate limiting
	AgeBracket AgeBracket
	Messages   []Message // Conversation history
}

// CompletionResponse contains the result of a completion request.
type CompletionResponse struct {
	Content    string
	TokensUsed int
}

// Provider defines the interface for LLM interactions.
type Provider interface {
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}
