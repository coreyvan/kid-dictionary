# Tasks: Layer Boundary Enforcement

**Input**: Design documents from `/specs/002-layer-boundaries/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests will be updated alongside implementation (existing tests in `_test.go` files).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go project**: `internal/` for application code, tests co-located as `*_test.go`
- Based on plan.md target structure:
  - `internal/domain/` - Domain models and interfaces
  - `internal/service/{entity}/` - Business logic
  - `internal/repository/{entity}/` - Persistence implementations
  - `internal/transport/` - HTTP/RPC handlers

---

## Phase 1: Setup (Domain Package Foundation)

**Purpose**: Create the new domain package with shared types and interfaces

- [X] T001 Create directory structure: `mkdir -p internal/domain internal/service/conversation internal/service/message internal/repository/conversation internal/repository/message`
- [X] T002 [P] Create domain errors in `internal/domain/errors.go` (ErrNotFound, ErrInvalidAgeBracket, ErrConversationNotFound)
- [X] T003 [P] Create Conversation model and ConversationRepository interface in `internal/domain/conversation.go`
- [X] T004 [P] Create Message model, enums (Role, ContentTier), and MessageRepository interface in `internal/domain/message.go`

---

## Phase 2: Foundational (Repository Layer Migration)

**Purpose**: Move repository implementations to new package structure - MUST complete before service layer migration

**⚠️ CRITICAL**: Service layer cannot be migrated until repositories are in place (to satisfy interface contracts)

- [X] T005 [P] Migrate conversation postgres implementation to `internal/repository/conversation/postgres.go` (implement domain.ConversationRepository)
- [X] T006 [P] Migrate conversation postgres tests to `internal/repository/conversation/postgres_test.go`
- [X] T007 [P] Migrate message postgres implementation to `internal/repository/message/postgres.go` (implement domain.MessageRepository)
- [X] T008 [P] Migrate message postgres tests to `internal/repository/message/postgres_test.go`
- [X] T009 Verify repositories compile and tests pass: `go test ./internal/repository/...`

**Checkpoint**: Repository layer complete - domain interfaces implemented

---

## Phase 3: User Story 1 - Developer Adds New Transport Handler (Priority: P1) 🎯 MVP

**Goal**: Transport layer only imports service layer packages, not repository layer packages

**Independent Test**: Verify `internal/transport/transport.go` has no imports from `internal/repository/*` packages

### Implementation for User Story 1

- [X] T010 [US1] Migrate conversation service to `internal/service/conversation/service.go` (import only domain package)
- [X] T011 [US1] Migrate conversation service tests to `internal/service/conversation/service_test.go`
- [X] T012 [US1] Migrate message service to `internal/service/message/service.go` (import only domain package)
- [X] T013 [US1] Migrate message classifier to `internal/service/message/classifier.go`
- [X] T014 [US1] Migrate message classifier tests to `internal/service/message/classifier_test.go`
- [X] T015 [US1] Migrate message service tests to `internal/service/message/service_test.go`
- [X] T016 [US1] Update `internal/transport/transport.go` imports: replace `internal/conversation` and `internal/message` with `internal/domain` and `internal/service/*`
- [X] T017 [US1] Remove `messageRepo message.Repository` field from transport server struct - use service layer instead
- [X] T018 [US1] Update `GetConversation` handler to call service method instead of direct repository access
- [X] T019 [US1] Verify transport compiles with no repository imports: `go build ./internal/transport/...`

**Checkpoint**: Transport layer correctly imports only service layer - US1 complete

---

## Phase 4: User Story 2 - Developer Modifies Service Layer (Priority: P2)

**Goal**: Service layer has zero imports from transport or repository packages

**Independent Test**: Verify all files in `internal/service/*` import only `internal/domain` and standard library

### Implementation for User Story 2

- [X] T020 [US2] Audit `internal/service/conversation/service.go` for any transport/repository imports and remove
- [X] T021 [US2] Audit `internal/service/message/service.go` for any transport/repository imports and remove
- [X] T022 [US2] Audit `internal/service/message/classifier.go` for any transport/repository imports and remove
- [X] T023 [US2] Verify service layer compiles with only domain imports: `go build ./internal/service/...`
- [X] T024 [US2] Run service tests with mocks to verify independence: `go test ./internal/service/...`

**Checkpoint**: Service layer is dependency-free from outer layers - US2 complete

---

## Phase 5: User Story 3 - Team Reviews Codebase for Architecture Compliance (Priority: P3)

**Goal**: Provide a mechanism to detect and report layer boundary violations

**Independent Test**: Run layer check command and verify it reports zero violations on compliant code, and correctly flags violations when introduced

### Implementation for User Story 3

- [X] T025 [US3] Create Taskfile task `layer-check` that verifies import rules using `go list -json ./internal/...`
- [X] T026 [US3] Add layer boundary documentation to quickstart.md with import rules table
- [X] T027 [US3] Run full layer check and document results
- [X] T028 [US3] Verify all success criteria met: SC-001 through SC-005

**Checkpoint**: Layer enforcement tooling in place - US3 complete

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup and finalization

- [X] T029 [P] Update `internal/config/wireup.go` with new package import paths
- [X] T030 [P] Update `cmd/server/main.go` if any imports changed
- [X] T031 Delete old mixed packages: `rm -rf internal/conversation internal/message` (after verifying no references remain)
- [X] T032 Run full test suite: `go test ./...`
- [X] T033 Run `go build ./...` to verify clean compilation
- [X] T034 Update CLAUDE.md with layer boundary guidance if needed

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - repository implementations need domain interfaces
- **User Story 1 (Phase 3)**: Depends on Foundational - transport needs service layer, which needs repositories
- **User Story 2 (Phase 4)**: Can run in parallel with US1 audit tasks, but service migration (T010-T015) must complete first
- **User Story 3 (Phase 5)**: Depends on US1 and US2 - needs clean architecture to verify
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Phase 2 (repositories must implement domain interfaces first)
- **User Story 2 (P2)**: Depends on service migration from US1 (T010-T015)
- **User Story 3 (P3)**: Depends on US1 and US2 (needs clean architecture to validate)

### Within Each Phase

- T002, T003, T004 can run in parallel (different files)
- T005, T006, T007, T008 can run in parallel (different packages)
- T029, T030 can run in parallel (different files)

### Parallel Opportunities

**Phase 1 (3 parallel tasks):**
```
T002 (errors.go) || T003 (conversation.go) || T004 (message.go)
```

**Phase 2 (4 parallel tasks):**
```
T005 (conv postgres) || T006 (conv test) || T007 (msg postgres) || T008 (msg test)
```

**Phase 6 (2 parallel tasks):**
```
T029 (wireup.go) || T030 (main.go)
```

---

## Parallel Example: Phase 1

```bash
# Launch all domain type definitions together:
Task: "Create domain errors in internal/domain/errors.go"
Task: "Create Conversation model and ConversationRepository interface in internal/domain/conversation.go"
Task: "Create Message model, enums, and MessageRepository interface in internal/domain/message.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (domain package)
2. Complete Phase 2: Foundational (repository migration)
3. Complete Phase 3: User Story 1 (transport + service migration)
4. **STOP and VALIDATE**: Verify `go build ./...` and `go test ./...` pass
5. Transport can no longer import repository - architecture enforced

### Incremental Delivery

1. Complete Setup + Foundational → Domain and repository layers ready
2. Add User Story 1 → Transport correctly layered → Test and verify
3. Add User Story 2 → Service layer isolated → Test and verify
4. Add User Story 3 → Enforcement tooling → Full validation
5. Polish → Cleanup old packages → Final verification

### Recommended Execution

Single developer, sequential:
1. Phase 1 (parallel within) → Phase 2 (parallel within) → Phase 3 → Phase 4 → Phase 5 → Phase 6

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently verifiable
- Run `go build ./...` after each phase to catch import errors early
- Run `go test ./...` after completing each user story
- Commit after each phase completion for easy rollback
- Old packages (`internal/conversation`, `internal/message`) deleted only in Phase 6 after all migrations verified
