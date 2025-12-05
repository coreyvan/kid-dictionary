# Implementation Plan: Parent Explanation Chat

**Branch**: `001-parent-explanation-chat` | **Date**: 2025-12-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-parent-explanation-chat/spec.md`

## Summary

Build the core functionality for parents to receive age-appropriate explanations of concepts for their children. The system accepts questions, classifies them by content tier, generates LLM-powered explanations tailored to three age brackets, and persists conversation history. Off-purpose requests are politely declined with guidance.

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: Connect-Go, Chi router, pgx, OpenAI SDK
**Storage**: PostgreSQL (conversation and message persistence)
**Testing**: Go testing with testify, integration tests against real PostgreSQL
**Target Platform**: Linux server (containerized)
**Project Type**: Single backend service (API-only for this feature; frontend is separate)
**Performance Goals**: <10 second response time for LLM completions, 100 concurrent users
**Constraints**: <10s p95 for explanation generation, 500 character input limit
**Scale/Scope**: 100 concurrent users, 30-day conversation retention

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | ✅ PASS | Proto files exist: `proto/kiddictionary/v1/` defines services, messages, enums |
| II. Test-Alongside Development | ✅ PASS | Tests will accompany implementation; integration test exists for OpenAI |
| III. Simplicity & YAGNI | ✅ PASS | Using existing LLM provider interface, no new abstractions required |
| IV. Observability | ✅ PASS | Existing logging infrastructure; will add request_id, duration logging |
| V. Content Safety | ✅ PASS | Age brackets defined in proto; content tier handling in spec (FR-004 to FR-007, FR-012) |

**Gate Result**: PASS - No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/001-parent-explanation-chat/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto analysis)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
cmd/
└── server/
    └── main.go              # Application entrypoint (exists)

internal/
├── config/
│   ├── config.go            # Configuration loading (exists)
│   ├── wireup.go            # Dependency injection (exists)
│   └── logging.go           # Logging setup (exists)
├── transport/
│   └── transport.go         # HTTP server (exists)
├── llm/
│   ├── llm.go               # Provider interface (exists)
│   ├── openai.go            # OpenAI implementation (exists)
│   ├── prompts.go           # Age bracket prompts (exists)
│   └── ratelimit.go         # Rate limiting (exists)
├── conversation/            # NEW: Conversation domain
│   ├── service.go           # Business logic
│   └── repository.go        # PostgreSQL persistence
├── message/                 # NEW: Message domain
│   ├── service.go           # Business logic + LLM integration
│   └── classifier.go        # Content tier + off-purpose classification
└── connect/                 # NEW: Connect handlers
    ├── conversation.go      # ConversationService implementation
    └── message.go           # MessageService implementation

proto/
└── kiddictionary/v1/
    ├── common.proto         # Shared types (exists)
    ├── conversation.proto   # ConversationService (exists)
    ├── message.proto        # MessageService (exists)
    └── auth.proto           # Auth service (exists, out of scope)

gen/
└── kiddictionary/v1/        # Generated Connect handlers
```

**Structure Decision**: Single backend service using existing Go project structure. New packages for conversation, message, and connect handlers. Protos already exist and define the API contract.

## Complexity Tracking

> No violations - table left empty per constitution.
