//go:build integration

package llm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOpenAIComplete(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	provider := NewOpenAIProvider(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	req := CompletionRequest{
		AgeBracket: AgeBracketLittleOnes,
		Messages: []Message{
			{
				Role:    "user",
				Content: "Tell me a short story about a brave little toaster.",
			},
		},
	}

	resp, err := provider.Complete(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
