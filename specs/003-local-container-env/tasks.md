# Tasks: Local Containerized Development Environment

**Input**: Design documents from `/specs/003-local-container-env/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, quickstart.md

**Tests**: No automated tests requested for this infrastructure feature. Manual validation via quickstart.md scenarios.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This feature adds Docker infrastructure files at repository root:
- `docker-compose.yml` - Main orchestration
- `Dockerfile` - Multi-stage build
- `.air.toml` - Hot reload config
- `Taskfile.yml` - Updated with container commands

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create base Docker configuration and update environment template

- [X] T001 [P] Update .env.example with database connection variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
- [X] T002 [P] Create .air.toml configuration for Go hot reload in repository root
- [X] T003 Create Dockerfile with multi-stage build (dev target with air, prod target with optimized binary)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core Docker Compose infrastructure that MUST be complete before user story commands work

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Create docker-compose.yml with PostgreSQL service (db) including healthcheck configuration
- [X] T005 Add migration service to docker-compose.yml using sql-migrate with depends_on db:service_healthy
- [X] T006 Add app service to docker-compose.yml with air hot reload, volume mounts, and depends_on migrate:service_completed_successfully
- [X] T007 Configure named volumes in docker-compose.yml (postgres_data, go_modules, go_cache)
- [X] T008 Configure environment variable passthrough in docker-compose.yml (DB_*, OPENAI_*, app config)

**Checkpoint**: Docker Compose file complete with all three services properly configured

---

## Phase 3: User Story 1 - Start Development Environment (Priority: P1) 🎯 MVP

**Goal**: Enable developers to start the complete environment with a single command

**Independent Test**: Run `task up`, verify database accessible on port 5432, service accessible on port 8080, migrations applied

### Implementation for User Story 1

- [X] T009 [US1] Add `up` task to Taskfile.yml that runs docker compose up -d
- [X] T010 [US1] Add `status` task to Taskfile.yml that runs docker compose ps
- [X] T011 [US1] Add `logs` task to Taskfile.yml that runs docker compose logs -f
- [X] T012 [US1] Add `logs:app` task to Taskfile.yml that runs docker compose logs -f app
- [X] T013 [US1] Update internal/config/wireup.go to read database connection from DB_* environment variables
- [X] T014 [US1] Validate start workflow: run task up, verify db healthcheck passes, migrations complete, app starts with hot reload

**Checkpoint**: `task up` starts complete environment, service responds to requests, hot reload works

---

## Phase 4: User Story 2 - Stop Development Environment (Priority: P2)

**Goal**: Enable developers to stop the environment cleanly with a single command

**Independent Test**: Run `task up`, then `task down`, verify all containers stopped and ports freed

### Implementation for User Story 2

- [X] T015 [US2] Add `down` task to Taskfile.yml that runs docker compose down
- [X] T016 [US2] Validate stop workflow: run task down after task up, verify containers stopped, ports released

**Checkpoint**: `task down` cleanly stops all containers, data preserved in volumes

---

## Phase 5: User Story 3 - Reset Development Environment (Priority: P3)

**Goal**: Enable developers to reset the database to a clean migrated state

**Independent Test**: Run `task up`, insert test data, run `task reset`, verify database is clean with fresh migrations

### Implementation for User Story 3

- [X] T017 [US3] Add `reset` task to Taskfile.yml that runs docker compose down -v && docker compose up -d
- [X] T018 [US3] Add `db:shell` task to Taskfile.yml that runs docker compose exec db psql -U postgres -d kid_dictionary
- [X] T019 [US3] Validate reset workflow: insert data, run task reset, verify database is empty but migrated

**Checkpoint**: `task reset` wipes all data and reapplies migrations

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, edge cases, and developer experience improvements

- [X] T020 [P] Update README.md with Docker development instructions
- [X] T021 [P] Add .dockerignore file to exclude unnecessary files from build context
- [X] T022 Validate quickstart.md scenarios end-to-end (fresh clone, setup, start, stop, reset)
- [X] T023 Test idempotent behavior: run task up twice, run task down twice, verify no errors
- [ ] T024a Test error handling: start with port 5432 already bound, verify clear error message displayed
- [ ] T024b Test error handling: start with port 8080 already bound, verify clear error message displayed
- [ ] T024c Test error handling: stop Docker daemon, run task up, verify "Docker not running" or similar message
- [ ] T024d Test error handling: corrupt migration file, run task up, verify migration error displayed and service does not start

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can proceed sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Logically requires US1 for testing
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Logically requires US1 for testing

### Within Each Phase

- Tasks marked [P] can run in parallel
- Sequential tasks within a story depend on previous tasks
- Validate checkpoint before moving to next phase

### Parallel Opportunities

- T001 and T002 can run in parallel (different files)
- T020 and T021 can run in parallel (different files)

---

## Parallel Example: Phase 1 Setup

```bash
# Launch setup tasks in parallel:
Task: "Update .env.example with database connection variables"
Task: "Create .air.toml configuration for Go hot reload"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T008)
3. Complete Phase 3: User Story 1 (T009-T014)
4. **STOP and VALIDATE**: Test `task up` works end-to-end
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Docker infrastructure ready
2. Add User Story 1 → `task up` works → MVP!
3. Add User Story 2 → `task down` works
4. Add User Story 3 → `task reset` works
5. Add Polish → Documentation and edge cases complete

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- T013 (wireup.go changes) may require updating existing code to use environment variables for database connection
