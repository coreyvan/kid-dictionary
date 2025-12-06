# Tasks: Git Hygiene Hooks

**Input**: Design documents from `/specs/005-git-hygiene-hooks/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Tests**: Not explicitly requested in the feature specification. Manual testing via quickstart.md scenarios.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and tool configuration

- [X] T001 Create golangci-lint configuration file at .golangci.yml
- [X] T002 [P] Create lefthook configuration file at lefthook.yml (skeleton structure only)
- [X] T003 Add `init:hooks` task to Taskfile.yml for lefthook installation and setup

---

## Phase 2: Foundational (Taskfile Tasks)

**Purpose**: Add Taskfile tasks that hooks will depend on

**CRITICAL**: These tasks MUST be complete before hook implementation

- [X] T004 Add `lint` task to Taskfile.yml to run golangci-lint on entire codebase
- [X] T005 [P] Add `lint:staged` task to Taskfile.yml to run golangci-lint on staged files only
- [X] T006 [P] Add `test:unit` task to Taskfile.yml to run unit tests (excluding acceptance tests)

**Checkpoint**: All supporting Taskfile tasks ready - hook implementation can begin

---

## Phase 3: User Story 1 - Pre-Commit Linting on Staged Files (Priority: P1)

**Goal**: Automatically lint staged Go files when committing to catch code quality issues early

**Independent Test**: Stage a file with linting errors, attempt to commit, verify commit is blocked with clear error messages

### Implementation for User Story 1

- [X] T007 [US1] Configure pre-commit hook in lefthook.yml to run golangci-lint on staged Go files
- [X] T008 [US1] Verify pre-commit hook blocks commits with linting errors (manual test per quickstart.md Scenario 2)
- [X] T009 [US1] Verify pre-commit hook allows clean commits (manual test per quickstart.md Scenario 1)
- [X] T010 [US1] Verify pre-commit hook skips non-Go files (manual test per quickstart.md acceptance scenario 3)

**Checkpoint**: User Story 1 complete - pre-commit linting functional

---

## Phase 4: User Story 2 - Pre-Push Full Validation (Priority: P2)

**Goal**: Comprehensive validation before pushing to ensure broken code doesn't reach the remote repository

**Independent Test**: Make changes that break tests or build, attempt to push, verify push is blocked

### Implementation for User Story 2

- [X] T011 [US2] Configure pre-push hook in lefthook.yml for code regeneration (task generate, task generate:mocks)
- [X] T012 [US2] Add dirty working tree detection to pre-push hook to fail if regeneration produces uncommitted changes
- [X] T013 [US2] Configure pre-push hook to run go mod tidy and fail if changes are produced
- [X] T014 [US2] Configure pre-push hook to run unit tests (task test:unit or go test ./...)
- [X] T015 [US2] Configure pre-push hook to run acceptance tests (task test:acceptance)
- [X] T016 [US2] Configure pre-push hook to build the binary (task build)
- [X] T017 [US2] Verify pre-push hook blocks on failing tests (manual test per quickstart.md Scenario 4)
- [X] T018 [US2] Verify pre-push hook blocks on uncommitted generated files (manual test per quickstart.md Scenario 5)
- [X] T019 [US2] Verify pre-push hook allows push when all validations pass (manual test per quickstart.md Scenario 3)

**Checkpoint**: User Story 2 complete - pre-push validation functional

---

## Phase 5: User Story 3 - Bypass Mechanism for Exceptional Cases (Priority: P3)

**Goal**: Allow developers to bypass hooks in exceptional circumstances using standard git flags

**Independent Test**: Use git commit/push with --no-verify flag and verify hooks are skipped

### Implementation for User Story 3

- [X] T020 [US3] Verify git commit --no-verify bypasses pre-commit hooks (manual test per quickstart.md Scenario 6)
- [X] T021 [US3] Verify git push --no-verify bypasses pre-push hooks (manual test per quickstart.md Scenario 6)

**Checkpoint**: User Story 3 complete - bypass mechanism verified (built-in git functionality)

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [X] T022 Run full quickstart.md validation (all 6 scenarios)
- [X] T023 Update CLAUDE.md Recent Changes section with feature summary and technologies

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on T001 (.golangci.yml) for lint tasks - BLOCKS hook implementation
- **User Story 1 (Phase 3)**: Depends on T004, T005 (lint tasks)
- **User Story 2 (Phase 4)**: Depends on T006 (test:unit task), can start after Phase 2
- **User Story 3 (Phase 5)**: Built-in git functionality, can start after US1 and US2 are testable
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - No dependencies on US1
- **User Story 3 (P3)**: Can start after US1 and US2 are implemented - verifies bypass of both hook types

### Within Each User Story

- Configuration before manual testing
- Verify each hook behavior matches quickstart.md scenarios
- Story complete before moving to next priority

### Parallel Opportunities

- T001 and T002 can run in parallel (different config files)
- T004, T005, T006 can run in parallel (different Taskfile tasks)
- US1 and US2 can be worked on in parallel after Foundational phase (separate hooks)

---

## Parallel Example: Foundational Phase

```bash
# Launch all Taskfile task additions together (after T001 complete):
Task: "Add lint task to Taskfile.yml"
Task: "Add lint:staged task to Taskfile.yml"
Task: "Add test:unit task to Taskfile.yml"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T006)
3. Complete Phase 3: User Story 1 (T007-T010)
4. **STOP and VALIDATE**: Test pre-commit linting independently
5. Commit and demo if ready

### Incremental Delivery

1. Complete Setup + Foundational -> Foundation ready
2. Add User Story 1 -> Test independently -> Commit (MVP: pre-commit linting!)
3. Add User Story 2 -> Test independently -> Commit (full pre-push validation)
4. Add User Story 3 -> Test independently -> Commit (bypass mechanism verified)
5. Each story adds value without breaking previous stories

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Manual testing follows quickstart.md scenarios exactly
- Commit after each phase or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
