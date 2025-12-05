# API Contracts: Parent Explanation Chat

**Date**: 2025-12-04
**Branch**: `001-parent-explanation-chat`

## Overview

This document describes the Connect-Go API contracts for the parent explanation chat feature. The proto definitions already exist in `proto/kiddictionary/v1/`.

---

## Services

### ConversationService

Manages conversation lifecycle.

**Proto File**: `proto/kiddictionary/v1/conversation.proto`

| RPC | Description | Request | Response |
|-----|-------------|---------|----------|
| CreateConversation | Start a new conversation | CreateConversationRequest | CreateConversationResponse |
| GetConversation | Retrieve conversation with messages | GetConversationRequest | GetConversationResponse |
| ListConversations | List user's conversations | ListConversationsRequest | ListConversationsResponse |
| DeleteConversation | Remove conversation and messages | DeleteConversationRequest | DeleteConversationResponse |

---

### MessageService

Handles sending messages and receiving AI responses.

**Proto File**: `proto/kiddictionary/v1/message.proto`

| RPC | Description | Request | Response |
|-----|-------------|---------|----------|
| SendMessage | Send user message, get AI response | SendMessageRequest | SendMessageResponse |

---

## Endpoint Details

### CreateConversation

Creates a new conversation with a selected age bracket.

**Request**:
```protobuf
message CreateConversationRequest {
  string title = 1;           // Optional - auto-generated if empty
  AgeBracket age_bracket = 2; // Required - LITTLE_ONES, GROWING_MINDS, or PRE_TEENS
}
```

**Response**:
```protobuf
message CreateConversationResponse {
  Conversation conversation = 1;
}
```

**Errors**:
| Code | Condition |
|------|-----------|
| InvalidArgument | age_bracket is UNSPECIFIED |

**Example**:
```json
// Request
{
  "age_bracket": "AGE_BRACKET_LITTLE_ONES"
}

// Response
{
  "conversation": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "New Conversation",
    "age_bracket": "AGE_BRACKET_LITTLE_ONES",
    "created_at": "2025-12-04T10:00:00Z",
    "updated_at": "2025-12-04T10:00:00Z"
  }
}
```

---

### GetConversation

Retrieves a conversation with its message history.

**Request**:
```protobuf
message GetConversationRequest {
  string id = 1; // Required - conversation UUID
}
```

**Response**:
```protobuf
message GetConversationResponse {
  Conversation conversation = 1;
  repeated Message messages = 2; // Ordered by created_at ascending
}
```

**Errors**:
| Code | Condition |
|------|-----------|
| InvalidArgument | id is empty or malformed |
| NotFound | conversation does not exist |

---

### ListConversations

Returns paginated list of user's conversations.

**Request**:
```protobuf
message ListConversationsRequest {
  int32 page_size = 1;    // Max items per page (default 20, max 100)
  string page_token = 2;  // Token for next page
}
```

**Response**:
```protobuf
message ListConversationsResponse {
  repeated Conversation conversations = 1; // Ordered by updated_at descending
  string next_page_token = 2;              // Empty if no more pages
}
```

---

### DeleteConversation

Removes a conversation and all its messages.

**Request**:
```protobuf
message DeleteConversationRequest {
  string id = 1; // Required - conversation UUID
}
```

**Response**:
```protobuf
message DeleteConversationResponse {} // Empty on success
```

**Errors**:
| Code | Condition |
|------|-----------|
| NotFound | conversation does not exist |

---

### SendMessage

Sends a user message and receives an AI-generated response.

**Request**:
```protobuf
message SendMessageRequest {
  string conversation_id = 1; // Required - conversation UUID
  string content = 2;         // Required - 1-500 characters
}
```

**Response**:
```protobuf
message SendMessageResponse {
  Message user_message = 1;      // The saved user message
  Message assistant_message = 2; // The AI-generated response
}
```

**Errors**:
| Code | Condition |
|------|-----------|
| InvalidArgument | content empty or >500 chars |
| NotFound | conversation does not exist |
| Unavailable | LLM service unavailable |
| ResourceExhausted | rate limit exceeded |

**Example - Normal Topic**:
```json
// Request
{
  "conversation_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Why is the sky blue?"
}

// Response
{
  "user_message": {
    "id": "...",
    "conversation_id": "550e8400-e29b-41d4-a716-446655440000",
    "role": "MESSAGE_ROLE_USER",
    "content": "Why is the sky blue?",
    "created_at": "2025-12-04T10:01:00Z"
  },
  "assistant_message": {
    "id": "...",
    "conversation_id": "550e8400-e29b-41d4-a716-446655440000",
    "role": "MESSAGE_ROLE_ASSISTANT",
    "content": "The sky looks blue because of the sunlight. When sunlight comes down, it bounces off tiny bits of air. Blue light bounces the most, so we see blue everywhere we look up!",
    "created_at": "2025-12-04T10:01:02Z"
  }
}
```

**Example - Sensitive Topic (Tier 2)**:
```json
// Request
{
  "conversation_id": "...",
  "content": "Why do people die?"
}

// Response (assistant_message.content includes soft guidance)
{
  "assistant_message": {
    "content": "This is a sensitive topic. You might want to have this conversation when your child feels safe and comfortable.\n\nWhen something dies, it stops moving and breathing. It can't wake up again. Everything that's alive will die someday - people, animals, and plants. It's okay to feel sad about that."
  }
}
```

**Example - Off-Purpose Request**:
```json
// Request
{
  "conversation_id": "...",
  "content": "Write me a poem about cats"
}

// Response (polite decline)
{
  "assistant_message": {
    "content": "I'm here to help you explain words and concepts to your child in age-appropriate ways. I can't write poems, but I'd be happy to help you explain what poetry is, or answer questions about cats!\n\nWhat would you like to explain to your child?"
  }
}
```

---

## Proto Modifications Required

### Add ContentTier enum to common.proto

```protobuf
// ContentTier classifies topic sensitivity for response handling.
enum ContentTier {
  CONTENT_TIER_UNSPECIFIED = 0;
  CONTENT_TIER_NORMAL = 1;      // Standard topics - direct explanation
  CONTENT_TIER_SENSITIVE = 2;   // Death, divorce, etc. - soft guidance
  CONTENT_TIER_CONTEXTUAL = 3;  // Religion, politics - request framing
  CONTENT_TIER_REDIRECT = 4;    // Harmful/off-purpose - decline
}
```

### Add content_tier to Message

```protobuf
message Message {
  string id = 1;
  string conversation_id = 2;
  MessageRole role = 3;
  string content = 4;
  google.protobuf.Timestamp created_at = 5;
  ContentTier content_tier = 6; // NEW: Classification for assistant messages
}
```

---

## HTTP Paths

All endpoints use Connect protocol POST requests:

| Service | Method | Path |
|---------|--------|------|
| ConversationService | CreateConversation | `/kiddictionary.v1.ConversationService/CreateConversation` |
| ConversationService | GetConversation | `/kiddictionary.v1.ConversationService/GetConversation` |
| ConversationService | ListConversations | `/kiddictionary.v1.ConversationService/ListConversations` |
| ConversationService | DeleteConversation | `/kiddictionary.v1.ConversationService/DeleteConversation` |
| MessageService | SendMessage | `/kiddictionary.v1.MessageService/SendMessage` |
