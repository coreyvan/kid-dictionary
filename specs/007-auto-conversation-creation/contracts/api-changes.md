# API Contract Changes: Auto Conversation Creation

**Date**: 2025-12-07
**Branch**: `007-auto-conversation-creation`

## Overview

This document describes the changes to the Connect-Go API contracts for the auto conversation creation feature. These changes are backwards compatible - existing clients that provide `conversation_id` will continue to work unchanged.

---

## Proto File Changes

### message.proto

**File**: `proto/kiddictionary/v1/message.proto`

#### SendMessageRequest (Modified)

```protobuf
message SendMessageRequest {
  // conversation_id is optional. If empty, a new conversation is auto-created.
  // When empty, age_bracket must be provided.
  string conversation_id = 1;

  // content is required. Must be 1-500 characters.
  string content = 2;

  // age_bracket is required when conversation_id is empty.
  // Ignored when conversation_id is provided.
  AgeBracket age_bracket = 3;
}
```

#### SendMessageResponse (Modified)

```protobuf
message SendMessageResponse {
  Message user_message = 1;
  Message assistant_message = 2;

  // conversation is present when a new conversation was auto-created.
  // Clients should use conversation.id for subsequent messages.
  Conversation conversation = 3;
}
```

**Import Required**: Add import for conversation.proto if not already present.

---

## Endpoint Behavior Changes

### SendMessage

**Path**: `/kiddictionary.v1.MessageService/SendMessage`

#### New Flow (Auto-Create)

**Request** (without conversation_id):
```json
{
  "content": "Why is the sky blue?",
  "age_bracket": "AGE_BRACKET_LITTLE_ONES"
}
```

**Response**:
```json
{
  "user_message": {
    "id": "msg-uuid-1",
    "conversation_id": "conv-uuid-new",
    "role": "MESSAGE_ROLE_USER",
    "content": "Why is the sky blue?",
    "created_at": "2025-12-07T10:00:00Z"
  },
  "assistant_message": {
    "id": "msg-uuid-2",
    "conversation_id": "conv-uuid-new",
    "role": "MESSAGE_ROLE_ASSISTANT",
    "content": "The sky looks blue because...",
    "created_at": "2025-12-07T10:00:02Z"
  },
  "conversation": {
    "id": "conv-uuid-new",
    "title": "Sky Being Blue",
    "age_bracket": "AGE_BRACKET_LITTLE_ONES",
    "created_at": "2025-12-07T10:00:00Z",
    "updated_at": "2025-12-07T10:00:00Z"
  }
}
```

#### Existing Flow (Unchanged)

**Request** (with conversation_id):
```json
{
  "conversation_id": "existing-conv-uuid",
  "content": "Tell me more about that"
}
```

**Response**:
```json
{
  "user_message": { ... },
  "assistant_message": { ... }
}
```

Note: `conversation` field is not present when using existing conversation.

---

## Error Cases

### New Validation Errors

| Condition | Connect Code | Message |
|-----------|--------------|---------|
| conversation_id empty, age_bracket unspecified | InvalidArgument | "age_bracket is required when conversation_id is not provided" |
| conversation_id empty, age_bracket invalid | InvalidArgument | "age_bracket must be LITTLE_ONES, GROWING_MINDS, or PRE_TEENS" |

### Existing Errors (Unchanged)

| Condition | Connect Code | Message |
|-----------|--------------|---------|
| content empty | InvalidArgument | "content is required" |
| content > 500 chars | InvalidArgument | "content exceeds maximum length of 500 characters" |
| conversation_id not found | NotFound | "conversation not found" |
| LLM unavailable | Unavailable | "service temporarily unavailable" |

---

## Backwards Compatibility

| Client Behavior | Impact |
|-----------------|--------|
| Always provides conversation_id | No change - works as before |
| Never provides age_bracket | Works when conversation_id provided; fails when conversation_id empty |
| Ignores conversation field in response | No change - field is optional |
| Uses conversation.id for follow-ups | Required for new auto-create flow |

---

## Client Migration Guide

### Before (Two-Step Flow)

```javascript
// Step 1: Create conversation
const conv = await client.createConversation({
  age_bracket: "AGE_BRACKET_LITTLE_ONES"
});

// Step 2: Send message
const response = await client.sendMessage({
  conversation_id: conv.id,
  content: "Why is the sky blue?"
});
```

### After (One-Step Flow)

```javascript
// Step 1: Send message (conversation auto-created)
const response = await client.sendMessage({
  content: "Why is the sky blue?",
  age_bracket: "AGE_BRACKET_LITTLE_ONES"
});

// Use returned conversation ID for follow-ups
const conversationId = response.conversation.id;

// Step 2+: Continue conversation
const followUp = await client.sendMessage({
  conversation_id: conversationId,
  content: "Tell me more about that"
});
```

---

## Title Generation Details

When a conversation is auto-created:

1. The system calls the LLM to generate a 3-4 word title from the message content
2. If title generation succeeds, the title is used
3. If title generation fails (timeout, rate limit, etc.), fallback title "New Conversation" is used
4. Title generation events are logged with success/failure, latency, and conversation ID

**Example Titles**:
- "Why is the sky blue?" → "Sky Being Blue"
- "How do plants grow?" → "Plants Growing Process"
- "What are dinosaurs?" → "Understanding Dinosaurs"
