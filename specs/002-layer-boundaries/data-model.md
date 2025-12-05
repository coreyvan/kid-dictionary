# Data Model: Layer Boundary Enforcement

**Feature**: 002-layer-boundaries
**Date**: 2025-12-05

## Overview

This document defines the domain models and their placement in the new layered architecture. The key principle is that domain models are owned by the service layer (conceptually) but placed in a shared `domain` package to enable proper import direction.

## Layer-Specific Model Ownership

### Domain Layer (`internal/domain/`)

Contains canonical domain models and repository interfaces. All layers can import this package.

### Transport Layer

Owns API-specific types (proto-generated types from `gen/kiddictionary/v1`). Transport handles:
- `v1.Conversation` ↔ `domain.Conversation`
- `v1.Message` ↔ `domain.Message`
- `v1.AgeBracket` ↔ `domain.AgeBracket`

### Repository Layer

Owns persistence-specific types (database row structs). Repository handles:
- `conversationRow` ↔ `domain.Conversation`
- `messageRow` ↔ `domain.Message`

## Domain Models

### Conversation

**Location**: `internal/domain/conversation.go`

| Field | Type | Description |
|-------|------|-------------|
| ID | `uuid.UUID` | Unique identifier |
| UserID | `*uuid.UUID` | Optional user association (nil for anonymous) |
| Title | `string` | Conversation title |
| AgeBracket | `AgeBracket` | Target age group for explanations |
| CreatedAt | `time.Time` | Creation timestamp |
| UpdatedAt | `time.Time` | Last modification timestamp |

**Enums**:
```
AgeBracket:
  - Unspecified (0)
  - LittleOnes (1)   // 0-5 years
  - GrowingMinds (2) // 5-10 years
  - PreTeens (3)     // 10+ years
```

**Validation Rules**:
- AgeBracket must be 1-3 (not Unspecified)
- Title defaults to "New Conversation" if empty

---

### Message

**Location**: `internal/domain/message.go`

| Field | Type | Description |
|-------|------|-------------|
| ID | `uuid.UUID` | Unique identifier |
| ConversationID | `uuid.UUID` | Parent conversation reference |
| Role | `Role` | Who sent the message |
| Content | `string` | Message text content |
| ContentTier | `ContentTier` | Topic sensitivity classification |
| CreatedAt | `time.Time` | Creation timestamp |

**Enums**:
```
Role:
  - Unspecified (0)
  - User (1)
  - Assistant (2)

ContentTier:
  - Unspecified (0)
  - Normal (1)      // Standard topics - direct explanation
  - Sensitive (2)   // Death, divorce, etc. - soft guidance prefix
  - Contextual (3)  // Religion, politics - request framing preference
  - Redirect (4)    // Harmful or off-purpose - polite decline
```

**Validation Rules**:
- Role must be User or Assistant (not Unspecified)
- ContentTier is only set for Assistant messages
- Content cannot be empty for User messages

---

### SendMessageResult

**Location**: `internal/domain/message.go`

| Field | Type | Description |
|-------|------|-------------|
| UserMessage | `*Message` | The stored user message |
| AssistantMessage | `*Message` | The generated assistant response |

## Repository Interfaces

### ConversationRepository

**Location**: `internal/domain/conversation.go`

```go
type ConversationRepository interface {
    Create(ctx context.Context, conv *Conversation) error
    GetByID(ctx context.Context, id uuid.UUID) (*Conversation, error)
    List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*Conversation, error)
    Update(ctx context.Context, conv *Conversation) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### MessageRepository

**Location**: `internal/domain/message.go`

```go
type MessageRepository interface {
    Create(ctx context.Context, msg *Message) error
    GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*Message, error)
    GetRecentByConversationID(ctx context.Context, conversationID uuid.UUID, limit int) ([]*Message, error)
}
```

## Domain Errors

**Location**: `internal/domain/errors.go`

| Error | Description |
|-------|-------------|
| `ErrNotFound` | Entity not found |
| `ErrInvalidAgeBracket` | Age bracket is unspecified or out of range |
| `ErrTitleRequired` | Title is empty when required |
| `ErrConversationNotFound` | Conversation does not exist |
| `ErrEmptyContent` | Message content is empty |

## Relationships

```
Conversation 1 ──────< Message
     │                    │
     └── AgeBracket       └── Role, ContentTier
```

- A Conversation has many Messages (ordered by CreatedAt ascending)
- Deleting a Conversation cascades to delete all its Messages

## State Transitions

No complex state machines. Entities are created, updated, and deleted.

## Data Volume Assumptions

- Conversations: Hundreds per user
- Messages: Tens to hundreds per conversation
- No pagination concerns at current scale
