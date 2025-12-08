# Implementation Plan: Auto Conversation Creation

**Branch**: `007-auto-conversation-creation` | **Date**: 2025-12-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-auto-conversation-creation/spec.md`

## Summary

Modify the SendMessage flow to automatically create a conversation when no conversation_id is provided. The system generates a 3-4 word title using an LLM call before creating the conversation, then processes the message normally. This eliminates the friction of requiring explicit conversation creation before asking the first question.

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**: Connect-Go (RPC), Chi (routing), OpenAI Go SDK (LLM)
**Storage**: PostgreSQL with pgx v5 driver
**Testing**: Go testing with testify, pgtestdb for integration tests
**Target Platform**: Linux server (Docker containers)
**Project Type**: Web API backend
**Performance Goals**: First message response under 15 seconds (including title generation)
**Constraints**: Title generation adds 1-3 seconds latency; sequential execution required
**Scale/Scope**: Existing message/conversation flow, single API endpoint modification

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Design Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | PASS | Will modify message.proto to add optional age_bracket and optional conversation in response |
| II. Test-Alongside Development | PASS | Tests will accompany implementation; mock LLM for title generation tests |
| III. Simplicity & YAGNI | PASS | Minimal change to existing service; title generation uses existing LLM provider interface |
| IV. Observability | PASS | FR-010 requires logging title generation events with success/failure, latency, conversation ID |
| V. Content Safety | N/A | Title generation doesn't involve content safety classification |
| Documentation Maintenance | PASS | README updates will be included if new env vars or commands added |
| Quality Gates | PASS | Existing CI/lint gates apply |

**Gate Result**: PASS - No violations detected.

### Post-Design Check (Phase 1 Complete)

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | PASS | Proto changes documented in contracts/api-changes.md; backwards compatible |
| II. Test-Alongside Development | PASS | Test approach defined in quickstart.md with unit and acceptance test examples |
| III. Simplicity & YAGNI | PASS | Single method addition to Provider interface; no new abstractions |
| IV. Observability | PASS | Logging fields specified in research.md (event, success, latency_ms, conversation_id, error) |
| V. Content Safety | N/A | Title generation is semantic summarization, not content classification |
| Documentation Maintenance | PASS | No new env vars or commands; existing README sufficient |
| Quality Gates | PASS | Proto regeneration, mock regeneration included in implementation steps |

**Post-Design Gate Result**: PASS - Design complies with all constitution principles.

## Project Structure

### Documentation (this feature)

```text
specs/007-auto-conversation-creation/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# Existing structure (Go backend)
cmd/server/              # Application entrypoint
internal/
├── config/              # Dependency wiring
├── domain/              # Shared models, repository interfaces
├── llm/                 # LLM provider interface and implementations
├── repository/          # PostgreSQL implementations
│   ├── conversation/
│   └── message/
├── service/             # Business logic
│   ├── conversation/    # Conversation service
│   └── message/         # Message service (PRIMARY CHANGES HERE)
└── transport/           # Connect-Go HTTP handlers
proto/kiddictionary/v1/  # Proto definitions (CHANGES HERE)
gen/                     # Generated Connect code

tests/
├── acceptance/          # End-to-end tests
└── (unit tests co-located with packages)
```

**Structure Decision**: Existing Go backend structure. Changes primarily in `internal/service/message/` (business logic), `proto/kiddictionary/v1/message.proto` (API contract), and `internal/llm/` (title generation function).

## Complexity Tracking

> **No violations to justify** - Implementation uses existing patterns and infrastructure.
