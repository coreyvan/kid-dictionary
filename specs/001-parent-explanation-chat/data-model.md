# Data Model: Parent Explanation Chat

**Date**: 2025-12-04
**Branch**: `001-parent-explanation-chat`

## Overview

This document defines the data model for the parent explanation chat feature, derived from the feature specification entities and proto definitions.

---

## Entities

### Conversation

Represents a chat session between a parent and the system.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, auto-generated | Unique identifier |
| user_id | UUID | FK (nullable for MVP) | Owner of the conversation |
| title | string | NOT NULL, max 200 chars | Display title (auto-generated from first question) |
| age_bracket | AgeBracket (enum) | NOT NULL | Target age group for explanations |
| created_at | timestamp | NOT NULL, default NOW() | Creation time |
| updated_at | timestamp | NOT NULL, default NOW() | Last modification time |

**Business Rules**:
- Title is auto-generated from the first message content (truncated to 50 chars + "...")
- Age bracket can be changed mid-conversation (affects subsequent responses only)
- Deletion cascades to all associated messages

**Proto Mapping**: `kiddictionary.v1.Conversation`

---

### Message

A single exchange in a conversation (user question or assistant response).

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, auto-generated | Unique identifier |
| conversation_id | UUID | FK, NOT NULL | Parent conversation |
| role | MessageRole (enum) | NOT NULL | USER or ASSISTANT |
| content | string | NOT NULL, max 5000 chars | Message text |
| content_tier | ContentTier (enum) | nullable | Classification (assistant messages only) |
| created_at | timestamp | NOT NULL, default NOW() | Creation time |

**Business Rules**:
- User messages: max 500 chars (validated at input)
- Assistant messages: max 5000 chars (LLM output limit)
- Content tier is set on assistant messages for audit purposes
- Messages are ordered by created_at within a conversation

**Proto Mapping**: `kiddictionary.v1.Message`

---

## Enumerations

### AgeBracket

| Value | Code | Description |
|-------|------|-------------|
| UNSPECIFIED | 0 | Not set (invalid for API calls) |
| LITTLE_ONES | 1 | Ages 0-5: simple, concrete words |
| GROWING_MINDS | 2 | Ages 5-10: metaphors, cause/effect |
| PRE_TEENS | 3 | Ages 10+: nuance, multiple perspectives |

**Proto Mapping**: `kiddictionary.v1.AgeBracket`

---

### MessageRole

| Value | Code | Description |
|-------|------|-------------|
| UNSPECIFIED | 0 | Not set (invalid) |
| USER | 1 | Parent's question |
| ASSISTANT | 2 | System's response |

**Proto Mapping**: `kiddictionary.v1.MessageRole`

---

### ContentTier

| Value | Code | Description | Handling |
|-------|------|-------------|----------|
| UNSPECIFIED | 0 | Not classified | N/A |
| NORMAL | 1 | Standard topics | Direct explanation |
| SENSITIVE | 2 | Death, divorce, illness | Soft guidance prefix |
| CONTEXTUAL | 3 | Religion, politics | Request framing preference |
| REDIRECT | 4 | Harmful or off-purpose | Polite decline |

**Note**: Not currently in proto; add to `common.proto` during implementation.

---

## Relationships

```
┌─────────────────┐
│     User        │ (future - out of scope for MVP)
│─────────────────│
│ id              │
│ email           │
│ default_bracket │
└────────┬────────┘
         │ 1
         │
         │ owns (nullable for MVP)
         │
         │ *
┌────────┴────────┐
│  Conversation   │
│─────────────────│
│ id              │
│ user_id (FK)    │
│ title           │
│ age_bracket     │
│ created_at      │
│ updated_at      │
└────────┬────────┘
         │ 1
         │
         │ contains
         │
         │ *
┌────────┴────────┐
│    Message      │
│─────────────────│
│ id              │
│ conversation_id │
│ role            │
│ content         │
│ content_tier    │
│ created_at      │
└─────────────────┘
```

---

## State Transitions

### Conversation Lifecycle

```
[Created] ──────────────────────────────────────► [Deleted]
    │                                                  ▲
    │                                                  │
    ▼                                                  │
[Active] ─────► messages added ─────► [Active]        │
    │                                                  │
    └──────────────── delete request ─────────────────┘
```

**States**:
- **Created**: New conversation with age bracket selected
- **Active**: Has one or more messages; can receive more
- **Deleted**: Removed along with all messages (hard delete)

### Message Processing Flow

```
[User Input]
    │
    ▼
[Validate] ──── invalid ───► [Return Error]
    │
    │ valid
    ▼
[Classify] ──── off-purpose ───► [Return Decline + Guidance]
    │
    │ on-purpose
    ▼
[Check Tier]
    │
    ├── Tier 1 (Normal) ───► [Generate Explanation]
    │
    ├── Tier 2 (Sensitive) ───► [Generate with Soft Guidance]
    │
    ├── Tier 3 (Contextual) ───► [Request Framing Preference]
    │
    └── Tier 4 (Redirect) ───► [Return Decline + Alternative]
```

---

## Indexes

| Table | Index | Columns | Purpose |
|-------|-------|---------|---------|
| conversations | idx_conversations_user_id | user_id | List user's conversations |
| conversations | idx_conversations_created_at | created_at | Order by recency |
| messages | idx_messages_conversation_id | conversation_id | Fetch conversation messages |
| messages | idx_messages_created_at | conversation_id, created_at | Order messages in conversation |

---

## Data Retention

Per spec assumptions:
- Conversations persist for at least 30 days
- Implementation: Scheduled job or TTL-based cleanup
- Soft delete not required for MVP

---

## Validation Rules

### Conversation
- `age_bracket` must be 1, 2, or 3 (not UNSPECIFIED)
- `title` auto-generated, max 200 chars

### Message (User Input)
- `content` required, 1-500 characters after trimming whitespace
- `conversation_id` must exist and be accessible

### Message (Assistant Output)
- `content` max 5000 characters
- `content_tier` set based on classification
