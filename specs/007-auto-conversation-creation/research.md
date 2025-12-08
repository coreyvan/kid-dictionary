# Research: Auto Conversation Creation

**Date**: 2025-12-07
**Branch**: `007-auto-conversation-creation`

## Research Topics

### 1. Title Generation Approach

**Decision**: Add a dedicated title generation method to the LLM provider interface.

**Rationale**:
- The existing `Provider` interface has only `Complete()` for explanation generation
- Title generation has different requirements: shorter output, no age bracket context, simpler prompt
- Adding a separate `GenerateTitle(ctx, content string) (string, error)` method keeps concerns separated
- Can use the same underlying OpenAI client with different parameters (lower max_tokens, different system prompt)

**Alternatives Considered**:
1. **Reuse Complete() with special parameters** - Rejected: overloads the method signature, age bracket param is meaningless for titles
2. **Separate title generation service** - Rejected: over-engineering for a simple string transformation
3. **String truncation without LLM** - Rejected: doesn't capture semantic meaning of the question

### 2. Title Generation Prompt Design

**Decision**: Use a minimal system prompt that instructs the LLM to summarize the question in 3-4 words.

**Rationale**:
- Simple, focused prompt reduces latency and token usage
- Explicit word count constraint ensures predictable output
- No age bracket context needed for title generation

**Prompt Design**:
```
System: Summarize the following question in 3-4 words. Output only the title, no punctuation, no quotes.
User: {message_content}
```

**Alternatives Considered**:
1. **Complex prompt with examples** - Rejected: adds tokens, not necessary for simple summarization
2. **Using message content directly as title** - Rejected: often too long, not descriptive

### 3. Error Handling for Title Generation

**Decision**: On title generation failure, use fallback title "New Conversation" and continue with message processing.

**Rationale**:
- Title is a UX enhancement, not critical functionality
- User should still get their explanation even if title generation fails
- Fallback is the same as current behavior when CreateConversation is called without a title
- Log the failure for debugging/monitoring (FR-010)

**Alternatives Considered**:
1. **Fail the entire request** - Rejected: poor UX, title is secondary to explanation
2. **Use first N characters of message** - Rejected: may be ugly/truncated mid-word
3. **Retry title generation** - Rejected: adds latency, not worth it for non-critical feature

### 4. Proto Message Modifications

**Decision**: Modify SendMessageRequest and SendMessageResponse in message.proto.

**Rationale**:
- Proto-first design (Constitution Principle I)
- Backwards compatible: existing fields remain, new fields are optional
- `conversation_id` becomes optional (empty string = auto-create)
- `age_bracket` added as optional field (required when conversation_id is empty)
- Response includes optional `Conversation` message when auto-created

**Proto Changes**:
```protobuf
message SendMessageRequest {
  string conversation_id = 1;  // Optional: if empty, creates new conversation
  string content = 2;
  AgeBracket age_bracket = 3;  // Required if conversation_id is empty
}

message SendMessageResponse {
  Message user_message = 1;
  Message assistant_message = 2;
  Conversation conversation = 3;  // Present when conversation was auto-created
}
```

### 5. Service Layer Changes

**Decision**: Modify `message.Service.SendMessage()` to handle auto-creation logic.

**Rationale**:
- Message service already depends on conversation repository
- Keeps the logic close to where it's used
- Avoids circular dependencies between services

**Flow**:
1. Check if conversation_id is empty
2. If empty:
   a. Validate age_bracket is provided
   b. Generate title via LLM (with fallback on failure)
   c. Create conversation via repository
   d. Log title generation event
3. Continue with existing SendMessage logic
4. Include conversation in response if auto-created

### 6. LLM Provider Interface Extension

**Decision**: Add `GenerateTitle(ctx context.Context, content string) (string, error)` to Provider interface.

**Rationale**:
- Clear separation of concerns
- Allows different implementations (mock for testing)
- Minimal interface extension

**Implementation Notes**:
- OpenAI implementation uses same client, different prompt
- Lower max_tokens (20 is sufficient for 3-4 words)
- Mock implementation returns predictable title for tests

### 7. Logging Requirements

**Decision**: Add structured logging for title generation events per FR-010.

**Log Fields**:
- `event`: "title_generation"
- `success`: bool
- `latency_ms`: int64
- `conversation_id`: string (after creation)
- `error`: string (on failure)

**Rationale**:
- Enables monitoring of title generation performance
- Helps identify LLM issues
- Follows Constitution Principle IV (Observability)

## Dependencies

| Dependency | Impact | Notes |
|------------|--------|-------|
| OpenAI Go SDK | Low | Already in use, no new dependency |
| message.proto | Medium | Breaking change for clients not providing age_bracket when auto-creating |
| LLM Provider interface | Medium | Adding new method requires updating mock |

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Title generation latency spikes | Medium | Low | Fallback title, timeout |
| LLM rate limiting | Low | Low | Same rate limiter as explanation calls |
| Proto breaking changes | Low | Medium | age_bracket only required when conversation_id empty |
