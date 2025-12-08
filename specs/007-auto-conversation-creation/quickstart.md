# Quickstart: Auto Conversation Creation

**Branch**: `007-auto-conversation-creation`

## What This Feature Does

Simplifies the user experience by allowing messages to be sent without first creating a conversation. When a message is sent without a `conversation_id`, the system:

1. Generates a 3-4 word title from the message content using the LLM
2. Creates a new conversation with that title
3. Processes the message normally
4. Returns the conversation details along with the response

## Key Files to Modify

### Proto Definition
- `proto/kiddictionary/v1/message.proto` - Add `age_bracket` to request, `conversation` to response

### LLM Layer
- `internal/llm/llm.go` - Add `GenerateTitle` to Provider interface
- `internal/llm/openai.go` - Implement `GenerateTitle` for OpenAI
- `internal/llm/prompts.go` - Add title generation prompt
- `internal/llm/provider_mock.go` - Regenerate mock with new method

### Service Layer
- `internal/service/message/service.go` - Add auto-creation logic to `SendMessage`

### Transport Layer
- `internal/transport/message.go` - Handle new request/response fields

## Implementation Order

1. **Proto changes** → `buf generate`
2. **LLM interface** → Add `GenerateTitle` method
3. **OpenAI implementation** → Implement title generation
4. **Mock regeneration** → `go generate ./internal/llm/...`
5. **Service changes** → Auto-create logic
6. **Transport changes** → Wire up new fields
7. **Tests** → Unit and acceptance tests

## Testing the Feature

### Unit Test (Service)
```go
func TestSendMessage_AutoCreatesConversation(t *testing.T) {
    // Setup mock LLM that returns a title
    mockLLM := &MockProvider{
        GenerateTitleFunc: func(ctx context.Context, content string) (string, error) {
            return "Test Title Here", nil
        },
        // ... Complete func
    }

    // Call SendMessage without conversation_id
    result, err := svc.SendMessage(ctx, uuid.Nil, "Little Ones", "Why is the sky blue?")

    // Assert conversation was created
    assert.NotNil(t, result.Conversation)
    assert.Equal(t, "Test Title Here", result.Conversation.Title)
}
```

### Acceptance Test
```go
func TestSendMessage_AutoCreate_E2E(t *testing.T) {
    // Send request without conversation_id
    resp, err := client.SendMessage(ctx, connect.NewRequest(&v1.SendMessageRequest{
        Content:    "Why do leaves change color?",
        AgeBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
    }))

    // Verify response includes conversation
    require.NoError(t, err)
    require.NotNil(t, resp.Msg.Conversation)
    require.NotEmpty(t, resp.Msg.Conversation.Id)
    require.NotEmpty(t, resp.Msg.Conversation.Title)
}
```

## Common Pitfalls

1. **Forgetting to regenerate protos** after modifying message.proto
2. **Not updating the mock** after adding `GenerateTitle` to the interface
3. **Missing age_bracket validation** when conversation_id is empty
4. **Not logging title generation events** per FR-010

## Verification Checklist

- [X] Proto changes compile with `buf generate`
- [X] `GenerateTitle` added to Provider interface
- [X] OpenAI implementation generates 3-4 word titles
- [X] Mock regenerated with `go generate`
- [X] Service auto-creates conversation when conversation_id empty
- [X] Service validates age_bracket when conversation_id empty
- [X] Service uses fallback title on LLM failure
- [X] Title generation logged with success/failure and latency
- [X] Response includes conversation when auto-created
- [X] Existing flow (with conversation_id) unchanged
- [X] Unit tests pass
- [X] Acceptance tests pass
