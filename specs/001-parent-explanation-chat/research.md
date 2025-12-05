# Research: Parent Explanation Chat

**Date**: 2025-12-04
**Branch**: `001-parent-explanation-chat`

## Overview

This document captures research findings and decisions for implementing the parent explanation chat feature. The codebase already has foundational infrastructure (LLM provider, protos, transport), so research focuses on new components.

---

## 1. Content Classification Approach

**Decision**: Use LLM-based classification with structured prompts

**Rationale**:
- Content tier classification (Normal, Sensitive, Contextual, Redirect) and off-purpose detection require semantic understanding
- Rule-based keyword matching is brittle and easily bypassed
- The same LLM call can both classify and generate, reducing latency
- Existing LLM provider interface supports this pattern

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Keyword blocklists | Too rigid; misses context ("bank" could be financial or riverbank) |
| Separate classification model | Adds latency; over-engineering for MVP |
| Pre-trained moderation API | Extra cost; doesn't handle off-purpose detection |

**Implementation Approach**:
- Extend system prompt to include classification instructions
- LLM returns structured response with tier classification and explanation
- Parse response to extract tier before delivering to user

---

## 2. Off-Purpose Request Detection

**Decision**: Include purpose-constraint instructions in system prompt

**Rationale**:
- FR-012 requires detecting off-purpose requests (not about explaining concepts to kids)
- LLM can distinguish "explain gravity" from "write me a poem" contextually
- Clarification allows meta-guidance ("how do I talk to my child about death")

**Implementation Approach**:
- System prompt explicitly states app purpose
- Instructions to classify requests as on-purpose or off-purpose
- Off-purpose responses use a standard template: explain app purpose, prompt for concept question

**On-Purpose Examples** (from clarifications):
- "What is photosynthesis?" → Direct concept explanation
- "How do I explain divorce to my 5-year-old?" → Meta-guidance allowed
- "Why do people die?" → Sensitive but on-purpose

**Off-Purpose Examples**:
- "Write me a poem" → Politely decline
- "Help with my taxes" → Politely decline
- "Tell me a joke" → Politely decline

---

## 3. Sensitive Topic Handling

**Decision**: Use content tier system from spec with soft guidance prefixes

**Rationale**:
- Tier 2 (Sensitive): death, divorce, illness, reproduction - need gentle framing
- Tier 3 (Contextual): religion, politics - need user's framing preference
- Tier 4 (Redirect): harmful content - decline and suggest alternative

**Implementation Approach**:
- Classification happens in initial LLM call
- Tier 2: Prepend soft guidance to response (e.g., "This is a topic that children may find difficult...")
- Tier 3: Return a follow-up question asking for framing preference before generating
- Tier 4: Return decline message without calling LLM for generation

**Soft Guidance Prefix Template** (Tier 2):
```
This is a sensitive topic. You might want to have this conversation when your child feels safe and comfortable. Here's a way to explain it:
```

---

## 4. Conversation Context Management

**Decision**: Store messages in PostgreSQL; pass recent history to LLM

**Rationale**:
- FR-010 requires maintaining context for follow-up questions
- FR-011 requires persistent history across sessions
- Existing pgx infrastructure in codebase

**Implementation Approach**:
- Store all messages (user and assistant) in `messages` table
- On SendMessage, retrieve last N messages (e.g., 10) for context
- Pass to LLM as conversation history
- Context window management: truncate older messages if token limit approached

**Data Retention**:
- Per assumptions: 30-day retention
- Implement via scheduled cleanup or soft delete with TTL

---

## 5. Response Time Optimization

**Decision**: Synchronous LLM call with 10-second timeout

**Rationale**:
- SC-001 requires <10 second response time
- Streaming would improve UX but adds complexity (constitution Principle III)
- OpenAI typical response time is 2-5 seconds for short completions

**Implementation Approach**:
- Set 10-second timeout on LLM provider calls
- Return user-friendly error if timeout exceeded (FR-009)
- Future: Consider streaming as separate feature

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|------------------|
| Streaming responses | Adds frontend complexity; MVP can use synchronous |
| Background job + polling | Over-engineering; latency already acceptable |
| Caching similar questions | Semantic matching is complex; defer to future |

---

## 6. Input Validation

**Decision**: Validate at service boundary before LLM call

**Rationale**:
- FR-008: 500 character limit
- Constitution security requirements: validate all input at service boundaries

**Validations**:
- Message content: max 500 characters, non-empty after trimming
- Age bracket: must be valid enum value
- Conversation ID: must be valid UUID, must belong to user (when auth added)

---

## 7. Database Schema

**Decision**: Use existing proto-defined structure; add content_tier to messages

**Rationale**:
- Protos define Conversation and Message types
- Need to persist content tier for audit/analytics
- User table exists but auth is out of scope for MVP (anonymous usage)

**Schema Additions**:
```sql
-- Conversations table (maps to proto Conversation)
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,  -- nullable for anonymous MVP
    title TEXT NOT NULL,
    age_bracket INTEGER NOT NULL,  -- enum value
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Messages table (maps to proto Message)
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role INTEGER NOT NULL,  -- enum: user=1, assistant=2
    content TEXT NOT NULL,
    content_tier INTEGER,  -- enum: normal=1, sensitive=2, contextual=3, redirect=4
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_conversations_created_at ON conversations(created_at);
```

---

## 8. Error Handling Strategy

**Decision**: Map domain errors to Connect status codes

**Rationale**:
- Constitution security requirements specify Connect error mapping
- FR-009 requires user-friendly error messages

**Error Mapping**:
| Domain Error | Connect Code | User Message |
|--------------|--------------|--------------|
| Input too long | InvalidArgument | "Your question is too long. Please keep it under 500 characters." |
| LLM unavailable | Unavailable | "We're having trouble right now. Please try again in a moment." |
| Rate limited | ResourceExhausted | "You're sending messages too quickly. Please wait a moment." |
| Conversation not found | NotFound | "Conversation not found." |
| Invalid age bracket | InvalidArgument | "Please select a valid age group." |

---

## Summary of Decisions

| Area | Decision |
|------|----------|
| Content classification | LLM-based with structured prompts |
| Off-purpose detection | System prompt constraints |
| Sensitive topics | Content tier system with soft guidance |
| Context management | PostgreSQL storage, pass recent history to LLM |
| Response time | Synchronous with 10s timeout |
| Input validation | Service boundary validation |
| Database | PostgreSQL with conversations + messages tables |
| Error handling | Domain errors mapped to Connect codes |

---

## Unresolved Items

None - all technical decisions resolved. Ready for Phase 1 design.
