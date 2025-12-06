package acceptance

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1"
	"github.com/coreyvan/kid-dictionary/gen/kiddictionary/v1/kiddictionaryv1connect"
	"github.com/coreyvan/kid-dictionary/internal/llm"
	"github.com/coreyvan/kid-dictionary/internal/testutil"
)

// TestSendMessageHappyPath verifies the end-to-end message flow:
// create a conversation, send a message, and receive a mock LLM response.
func TestSendMessageHappyPath(t *testing.T) {
	// Setup test database
	pool := testutil.NewTestDB(t)

	// Create mock LLM with configured response
	mockResponse := "This is a simple explanation for kids!"
	mockLLM := &llm.ProviderMock{
		CompleteFunc: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			return llm.CompletionResponse{
				Content:    mockResponse,
				TokensUsed: 10,
			}, nil
		},
	}

	// Start test server with mock dependencies
	server := testutil.NewTestServer(t, pool, mockLLM)

	// Create clients
	convClient := kiddictionaryv1connect.NewConversationServiceClient(
		http.DefaultClient,
		server.URL,
	)
	msgClient := kiddictionaryv1connect.NewMessageServiceClient(
		http.DefaultClient,
		server.URL,
	)

	ctx := context.Background()

	// Create a conversation
	createResp, err := convClient.CreateConversation(ctx, connect.NewRequest(&v1.CreateConversationRequest{
		Title:      "Test Conversation",
		AgeBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
	}))
	require.NoError(t, err)
	require.NotNil(t, createResp.Msg.Conversation)
	conversationID := createResp.Msg.Conversation.Id
	assert.NotEmpty(t, conversationID)

	// Send a message
	sendResp, err := msgClient.SendMessage(ctx, connect.NewRequest(&v1.SendMessageRequest{
		ConversationId: conversationID,
		Content:        "What is the sun?",
	}))
	require.NoError(t, err)

	// Verify user message
	require.NotNil(t, sendResp.Msg.UserMessage)
	assert.Equal(t, "What is the sun?", sendResp.Msg.UserMessage.Content)
	assert.Equal(t, v1.MessageRole_MESSAGE_ROLE_USER, sendResp.Msg.UserMessage.Role)

	// Verify assistant message contains mock response
	require.NotNil(t, sendResp.Msg.AssistantMessage)
	assert.Equal(t, mockResponse, sendResp.Msg.AssistantMessage.Content)
	assert.Equal(t, v1.MessageRole_MESSAGE_ROLE_ASSISTANT, sendResp.Msg.AssistantMessage.Role)

	// Verify mock was called exactly once
	assert.Len(t, mockLLM.CompleteCalls(), 1)
}

// TestConversationCreation verifies conversation creation works independently.
func TestConversationCreation(t *testing.T) {
	// Setup test database
	pool := testutil.NewTestDB(t)

	// Create mock LLM (not used in this test but required for server setup)
	mockLLM := &llm.ProviderMock{
		CompleteFunc: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			return llm.CompletionResponse{Content: "unused", TokensUsed: 10}, nil
		},
	}

	// Start test server
	server := testutil.NewTestServer(t, pool, mockLLM)

	// Create client
	client := kiddictionaryv1connect.NewConversationServiceClient(
		http.DefaultClient,
		server.URL,
	)

	ctx := context.Background()

	testCases := []struct {
		name       string
		title      string
		ageBracket v1.AgeBracket
	}{
		{
			name:       "little ones",
			title:      "Questions for Little Ones",
			ageBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
		},
		{
			name:       "growing minds",
			title:      "Questions for Growing Minds",
			ageBracket: v1.AgeBracket_AGE_BRACKET_GROWING_MINDS,
		},
		{
			name:       "pre-teens",
			title:      "Questions for Pre-Teens",
			ageBracket: v1.AgeBracket_AGE_BRACKET_PRE_TEENS,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.CreateConversation(ctx, connect.NewRequest(&v1.CreateConversationRequest{
				Title:      tc.title,
				AgeBracket: tc.ageBracket,
			}))

			require.NoError(t, err)
			require.NotNil(t, resp.Msg.Conversation)

			conv := resp.Msg.Conversation
			assert.NotEmpty(t, conv.Id)
			assert.Equal(t, tc.title, conv.Title)
			assert.Equal(t, tc.ageBracket, conv.AgeBracket)
			assert.NotNil(t, conv.CreatedAt)
			assert.NotNil(t, conv.UpdatedAt)
		})
	}
}

// TestMockResponseCustomization verifies that the mock LLM can return
// age-bracket-specific responses for targeted testing.
func TestMockResponseCustomization(t *testing.T) {
	// Setup test database
	pool := testutil.NewTestDB(t)

	// Define age-bracket-specific responses
	littleOnesResponse := "The sun is a big, warm light in the sky!"
	growingMindsResponse := "The sun is a star that gives us light and heat."
	preTeensResponse := "The sun is a G-type main-sequence star at the center of our solar system."

	// Track which age bracket was requested
	var requestedAgeBracket llm.AgeBracket

	// Create mock LLM with age-bracket-specific responses
	mockLLM := &llm.ProviderMock{
		CompleteFunc: func(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
			requestedAgeBracket = req.AgeBracket
			var content string
			switch req.AgeBracket {
			case llm.AgeBracketLittleOnes:
				content = littleOnesResponse
			case llm.AgeBracketGrowingMinds:
				content = growingMindsResponse
			case llm.AgeBracketPreTeens:
				content = preTeensResponse
			default:
				content = "Default response"
			}
			return llm.CompletionResponse{
				Content:    content,
				TokensUsed: 10,
			}, nil
		},
	}

	// Start test server
	server := testutil.NewTestServer(t, pool, mockLLM)

	// Create clients
	convClient := kiddictionaryv1connect.NewConversationServiceClient(
		http.DefaultClient,
		server.URL,
	)
	msgClient := kiddictionaryv1connect.NewMessageServiceClient(
		http.DefaultClient,
		server.URL,
	)

	ctx := context.Background()

	// Create conversation with AgeBracketLittleOnes
	createResp, err := convClient.CreateConversation(ctx, connect.NewRequest(&v1.CreateConversationRequest{
		Title:      "Test for Little Ones",
		AgeBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
	}))
	require.NoError(t, err)
	conversationID := createResp.Msg.Conversation.Id

	// Send message
	sendResp, err := msgClient.SendMessage(ctx, connect.NewRequest(&v1.SendMessageRequest{
		ConversationId: conversationID,
		Content:        "What is the sun?",
	}))
	require.NoError(t, err)

	// Verify the response matches the configured age bracket response
	assert.Equal(t, littleOnesResponse, sendResp.Msg.AssistantMessage.Content)

	// Verify the correct age bracket was passed to the mock
	assert.Equal(t, llm.AgeBracketLittleOnes, requestedAgeBracket)

	// Verify mock was called exactly once
	assert.Len(t, mockLLM.CompleteCalls(), 1)
}
