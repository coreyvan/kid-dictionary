# Implementation Plan: Mock LLM Acceptance Tests

**Branch**: `004-mock-llm-acceptance-tests` | **Date**: 2025-12-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-mock-llm-acceptance-tests/spec.md`

## Summary

Create a mock implementation of the `llm.Provider` interface and acceptance test infrastructure that allows developers to run end-to-end tests against an in-process HTTP server without requiring real LLM API calls. Tests use pgtestdb for per-test database isolation and the existing wireup pattern for dependency injection.

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**: Connect-Go (RPC), Chi (routing), pgtestdb (test DB isolation), testify (assertions)
**Storage**: PostgreSQL with pgx v5 driver, pgtestdb for test isolation
**Testing**: go test with acceptance tests in separate package
**Target Platform**: Developer workstations with Docker (macOS/Linux)
**Project Type**: Single backend service (existing)
**Performance Goals**: Acceptance tests complete in <30 seconds
**Constraints**: No external API keys required for tests, offline-capable after Docker images cached
**Scale/Scope**: ~3-5 acceptance test cases covering happy paths

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | PASS | Using existing generated Connect-Go clients from proto definitions |
| II. Test-Alongside Development | PASS | This feature explicitly implements contract tests for API boundaries |
| III. Simplicity & YAGNI | PASS | Mock provider is minimal; reuses existing wireup pattern |
| IV. Observability | N/A | Test infrastructure, not production code |
| V. Content Safety | N/A | Tests use mock responses, no real content generation |
| Technology Stack | PASS | Go, Connect-Go, PostgreSQL - all aligned |
| Security Requirements | PASS | No secrets required for tests; database credentials are local Docker |

**Gate Status**: PASS - No violations requiring justification

## Project Structure

### Documentation (this feature)

```text
specs/004-mock-llm-acceptance-tests/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal - test infrastructure)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (empty - using existing protos)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# New files for this feature
internal/llm/mock.go                    # Mock LLM provider implementation
internal/testutil/server.go             # Test server setup helper using wireup
internal/testutil/database.go           # pgtestdb integration helper
tests/acceptance/acceptance_test.go     # Acceptance tests using generated clients

# Existing files (no changes needed)
internal/llm/llm.go                     # Provider interface (already exists)
internal/config/wireup.go               # Dependency injection (already supports mocks)
gen/kiddictionary/v1/...                # Generated Connect-Go clients (already exist)
```

**Structure Decision**: Single project structure preserved. Test infrastructure added under `internal/testutil/` following Go conventions. Acceptance tests in `tests/acceptance/` to separate from unit tests.

## Complexity Tracking

> No violations to justify - Constitution Check passed.
