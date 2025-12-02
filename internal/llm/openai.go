package llm

import (
	"context"
	"errors"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

const (
	defaultModel = "gpt-5-nano-2025-08-07"
)

var _ Provider = (*OpenAIProvider)(nil)

// OpenAIProvider implements the Provider interface using OpenAI's API.
type OpenAIProvider struct {
	client *openai.Client
	model  shared.ChatModel
}

// OpenAIOption configures the OpenAI provider.
type OpenAIOption func(*OpenAIProvider)

// WithModel sets the model to use for completions.
func WithModel(model string) OpenAIOption {
	return func(p *OpenAIProvider) {
		p.model = model
	}
}

// NewOpenAIProvider creates a new OpenAI provider with the given API key.
func NewOpenAIProvider(apiKey string, opts ...OpenAIOption) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))

	p := &OpenAIProvider{
		client: &client,
		model:  defaultModel,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Complete generates a completion for the given request.
func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	systemPrompt := SystemPromptForAgeBracket(req.AgeBracket)

	// Build messages: system prompt + conversation history
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages)+1)
	messages = append(messages, openai.SystemMessage(systemPrompt))

	for _, msg := range req.Messages {
		switch msg.Role {
		case "user":
			messages = append(messages, openai.UserMessage(msg.Content))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    p.model,
		Messages: messages,
	})
	if err != nil {
		return CompletionResponse{}, mapOpenAIError(err)
	}

	if len(resp.Choices) == 0 {
		return CompletionResponse{}, errors.New("no completion choices returned")
	}

	return CompletionResponse{
		Content:    resp.Choices[0].Message.Content,
		TokensUsed: int(resp.Usage.TotalTokens),
	}, nil
}

// mapOpenAIError converts OpenAI API errors to domain errors.
func mapOpenAIError(err error) error {
	var apiErr *openai.Error
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized:
			return ErrInvalidAPIKey
		case http.StatusTooManyRequests:
			return ErrRateLimited
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			return ErrProviderUnavailable
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}

	return err
}
