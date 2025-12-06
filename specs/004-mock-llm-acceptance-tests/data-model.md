# Data Model: Mock LLM Acceptance Tests

**Branch**: `004-mock-llm-acceptance-tests` | **Date**: 2025-12-05

## Overview

This feature is test infrastructure focused and does not introduce new domain entities. It adds test doubles and test utilities that operate on existing domain models (Conversation, Message).

## Test Infrastructure Entities

### MockProvider

A test double implementing the `llm.Provider` interface.

| Field | Type | Description |
|-------|------|-------------|
| DefaultResponse | string | Response returned when no age-bracket-specific response configured |
| Responses | map[AgeBracket]string | Age-bracket-specific responses for targeted testing |
| CallCount | int | Number of Complete() calls (for verification) |

**Implements**: `llm.Provider` interface from `internal/llm/llm.go`

**Behavior**:
- Returns configured response for matching age bracket
- Falls back to DefaultResponse if no match
- Increments CallCount on each call
- Always returns `TokensUsed: 10` (deterministic for testing)

### TestServer

A test HTTP server wrapping the production handlers.

| Component | Type | Description |
|-----------|------|-------------|
| Server | *httptest.Server | stdlib test server with auto-allocated port |
| URL | string | Base URL for client connections |
| MockLLM | *MockProvider | Injected mock for LLM operations |
| DB | *sql.DB | Isolated test database from pgtestdb |

**Lifecycle**:
1. Create isolated database via pgtestdb
2. Initialize wireup with mock dependencies
3. Start httptest.Server
4. Return URL for client connections
5. Cleanup handled by t.Cleanup()

## Existing Entities Used (No Changes)

### Conversation

Existing domain entity at `internal/domain/conversation.go`. Used in tests via generated Connect client.

| Field | Used In Tests |
|-------|---------------|
| ID | CreateConversation response validation |
| Title | Request/response verification |
| AgeBracket | Mock response routing |

### Message

Existing domain entity at `internal/domain/message.go`. Used in tests via SendMessage endpoint.

| Field | Used In Tests |
|-------|---------------|
| ID | Response validation |
| Content | Mock response verification |
| Role | USER vs ASSISTANT distinction |
| ContentTier | Optional field verification |

## Relationship Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Test Execution                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────┐         ┌────────────────────────────┐  │
│  │   Test Case    │         │       pgtestdb             │  │
│  │                │────────▶│  (isolated DB per test)    │  │
│  └────────────────┘         └────────────────────────────┘  │
│           │                              │                   │
│           │                              ▼                   │
│           │                 ┌────────────────────────────┐  │
│           │                 │   WireupDeps (injected)    │  │
│           │                 │  - DBPool: test DB         │  │
│           │                 │  - LLMProvider: MockProvider│  │
│           │                 └────────────────────────────┘  │
│           │                              │                   │
│           ▼                              ▼                   │
│  ┌────────────────┐         ┌────────────────────────────┐  │
│  │ Generated      │────────▶│    httptest.Server         │  │
│  │ Connect Client │         │  (production handlers)     │  │
│  └────────────────┘         └────────────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## No Schema Changes

This feature adds no database tables or migrations. All persistence uses existing Conversation and Message tables via the existing repositories.
