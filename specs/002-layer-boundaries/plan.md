# Implementation Plan: Layer Boundary Enforcement

**Branch**: `002-layer-boundaries` | **Date**: 2025-12-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-layer-boundaries/spec.md`

## Summary

Refactor the codebase to enforce strict three-layer architecture (transport → service → repository) where:
- Transport layer only imports service layer packages
- Service layer owns domain models and defines repository interfaces; has zero imports from transport or repository packages
- Repository layer implements service-defined interfaces and translates between persistence types and service domain models

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**: Connect-Go, Chi router, pgx v5
**Storage**: PostgreSQL
**Testing**: go test with testify
**Target Platform**: Linux server (Docker)
**Project Type**: single (Go backend service)
**Performance Goals**: N/A (architectural refactoring, no performance changes)
**Constraints**: Must maintain API compatibility; no breaking changes to proto definitions
**Scale/Scope**: ~20 Go files across 4 internal packages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | ✅ Pass | No API changes; proto definitions remain unchanged |
| II. Test-Alongside Development | ✅ Pass | Existing tests will be updated to match new package structure |
| III. Simplicity & YAGNI | ⚠️ Review | Layer separation adds packages but reduces coupling. Justified for maintainability |
| IV. Observability | ✅ Pass | No changes to logging or metrics |
| V. Content Safety | ✅ Pass | No changes to content handling |

**Complexity Justification for Principle III**:
The layer separation introduces additional packages (service, repository, transport per domain), but this is the minimum complexity required to enforce proper dependency direction. The alternative (current mixed packages) leads to tight coupling and makes testing harder. This aligns with "Simplicity" by making each layer independently testable.

## Project Structure

### Documentation (this feature)

```text
specs/002-layer-boundaries/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (empty - no API changes)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

**Current Structure (with violations)**:
```text
internal/
├── config/              # Cross-cutting ✓
├── connect/             # Transport helper ✓
├── conversation/        # MIXED: service + repository types
│   ├── repository.go    # Domain model + interface
│   ├── service.go       # Business logic
│   ├── postgres.go      # Repository impl
│   └── *_test.go
├── llm/                 # Cross-cutting ✓
├── message/             # MIXED: service + repository types
│   ├── repository.go    # Domain model + interface
│   ├── service.go       # Business logic
│   ├── classifier.go    # Business logic
│   ├── postgres.go      # Repository impl
│   └── *_test.go
└── transport/           # Transport layer ✓ (but imports repository)
    └── transport.go
```

**Target Structure (layer-compliant)**:
```text
internal/
├── config/              # Cross-cutting (unchanged)
├── connect/             # Transport helper (unchanged)
├── llm/                 # Cross-cutting (unchanged)
├── domain/              # NEW: Shared domain models & service interfaces
│   ├── conversation.go  # Conversation model, AgeBracket, Repository interface
│   ├── message.go       # Message model, ContentTier, Role, Repository interface
│   └── errors.go        # Domain errors
├── service/             # NEW: Business logic services
│   ├── conversation/
│   │   ├── service.go   # ConversationService
│   │   └── service_test.go
│   └── message/
│       ├── service.go   # MessageService
│       ├── classifier.go
│       └── *_test.go
├── repository/          # NEW: Repository implementations
│   ├── conversation/
│   │   ├── postgres.go
│   │   └── postgres_test.go
│   └── message/
│       ├── postgres.go
│       └── postgres_test.go
└── transport/           # Transport layer (updated imports)
    └── transport.go
```

**Structure Decision**: The target structure separates concerns into distinct packages:
- `internal/domain/` - Shared domain models owned by the service layer conceptually, but in a separate package for import clarity
- `internal/service/` - Business logic, imports only domain
- `internal/repository/` - Persistence, imports domain to implement interfaces
- `internal/transport/` - HTTP/RPC handling, imports domain and service

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| New `domain/` package | Central location for shared types | Keeping types in service packages creates circular import risk |
| Separate service subpackages | Clear boundaries per domain entity | Single service package would grow unwieldy |
| Separate repository subpackages | Mirrors service structure | Consistency with service layer |

## Current Violations Found

1. **Transport imports repository directly**:
   - `transport/transport.go:38` - `messageRepo message.Repository` field
   - `transport/transport.go:158` - Direct call to `s.messageRepo.GetByConversationID()`

2. **Mixed packages**:
   - `conversation/` package contains: domain model, repository interface, service, postgres impl
   - `message/` package contains: domain model, repository interface, service, classifier, postgres impl

3. **Transport uses domain types from "mixed" packages**:
   - Creates tight coupling between transport and repository implementation
