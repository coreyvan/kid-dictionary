# Data Model: Auto Conversation Creation

**Date**: 2025-12-07
**Branch**: `007-auto-conversation-creation`

## Overview

This feature modifies the SendMessage API contract to support automatic conversation creation. No changes to the underlying database schema are required - the feature adds optional fields to existing proto messages and extends the LLM provider interface.

---

## Entities

### Existing Entities (No Changes)

#### Conversation

The existing `domain.Conversation` struct remains unchanged:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| ID | UUID | PK, auto-generated | Unique identifier |
| UserID | *UUID | FK (nullable for MVP) | Owner of the conversation |
| Title | string | NOT NULL, max 200 chars | Display title (now auto-generated via LLM) |
| AgeBracket | AgeBracket (enum) | NOT NULL | Target age group for explanations |
| CreatedAt | timestamp | NOT NULL | Creation time |
| UpdatedAt | timestamp | NOT NULL | Last modification time |

**Note**: Title was previously client-provided or defaulted to "New Conversation". Now it can be LLM-generated from the first message content.

#### Message

The existing `domain.Message` struct remains unchanged.

---

## Modified API Messages

### SendMessageRequest (Proto)

**Current**:
```protobuf
message SendMessageRequest {
  string conversation_id = 1;  // Required
  string content = 2;          // Required
}
```

**Modified**:
```protobuf
message SendMessageRequest {
  string conversation_id = 1;  // Optional: if empty, auto-creates conversation
  string content = 2;          // Required: 1-500 characters
  AgeBracket age_bracket = 3;  // Required when conversation_id is empty
}
```

**Validation Rules**:
- `content`: Required, 1-500 characters after whitespace trimming
- `conversation_id`: Optional
  - If provided: must be valid UUID, conversation must exist
  - If empty: `age_bracket` becomes required
- `age_bracket`: Conditionally required
  - Required when `conversation_id` is empty
  - Must be LITTLE_ONES, GROWING_MINDS, or PRE_TEENS (not UNSPECIFIED)
  - Ignored when `conversation_id` is provided

### SendMessageResponse (Proto)

**Current**:
```protobuf
message SendMessageResponse {
  Message user_message = 1;
  Message assistant_message = 2;
}
```

**Modified**:
```protobuf
message SendMessageResponse {
  Message user_message = 1;
  Message assistant_message = 2;
  Conversation conversation = 3;  // Present when conversation was auto-created
}
```

**Business Rules**:
- `conversation`: Only populated when the request triggered auto-creation
- Contains the newly created conversation with LLM-generated or fallback title
- Clients MUST use `conversation.id` for subsequent messages in the same session

---

## New Domain Types

### SendMessageWithAutoCreateResult

Extended result type for the service layer:

```go
// SendMessageWithAutoCreateResult extends SendMessageResult with optional conversation.
type SendMessageWithAutoCreateResult struct {
    UserMessage      *Message
    AssistantMessage *Message
    Conversation     *Conversation  // Non-nil if auto-created
}
```

---

## Interface Extensions

### LLM Provider Interface

**Current**:
```go
type Provider interface {
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}
```

**Extended**:
```go
type Provider interface {
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
    GenerateTitle(ctx context.Context, content string) (string, error)
}
```

**GenerateTitle Specification**:
- Input: User's message content (1-500 characters)
- Output: 3-4 word summary title (no punctuation, no quotes)
- Errors: Same error types as Complete (ErrRateLimited, ErrProviderUnavailable, ErrTimeout)
- Timeout: 10 seconds (configurable)

---

## State Transitions

### Auto-Create Flow

```
[No Conversation]
       │
       │ SendMessage(content, age_bracket) without conversation_id
       ▼
[Generate Title] ──── failure ───► [Use Fallback "New Conversation"]
       │                                        │
       │ success                                │
       ▼                                        │
[Create Conversation] ◄─────────────────────────┘
       │
       ▼
[Process Message] (existing flow)
       │
       ▼
[Return Response with Conversation]
```

### Existing Flow (Unchanged)

```
[Existing Conversation]
       │
       │ SendMessage(conversation_id, content)
       ▼
[Validate Conversation Exists]
       │
       ▼
[Process Message] (existing flow)
       │
       ▼
[Return Response without Conversation]
```

---

## Error Scenarios

| Condition | Error Code | Error Message |
|-----------|------------|---------------|
| content empty | InvalidArgument | "content is required" |
| content > 500 chars | InvalidArgument | "content exceeds maximum length" |
| conversation_id empty AND age_bracket unspecified | InvalidArgument | "age_bracket is required for new conversations" |
| conversation_id empty AND age_bracket invalid | InvalidArgument | "age_bracket must be LITTLE_ONES, GROWING_MINDS, or PRE_TEENS" |
| conversation_id provided AND not found | NotFound | "conversation not found" |
| LLM unavailable | Unavailable | "service temporarily unavailable" |

---

## Database Impact

**No schema changes required.**

The conversation table already supports:
- Auto-generated UUIDs for ID
- Title field (varchar 200)
- All required fields for auto-creation

The only change is that titles are now LLM-generated rather than client-provided or defaulted.
