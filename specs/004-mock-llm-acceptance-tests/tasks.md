# Tasks: Mock LLM Acceptance Tests

**Input**: Design documents from `/specs/004-mock-llm-acceptance-tests/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: This feature is itself test infrastructure. Acceptance tests ARE the deliverables.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Based on plan.md structure:
- Mock provider: `internal/llm/mock.go`
- Test utilities: `internal/testutil/`
- Acceptance tests: `tests/acceptance/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add dependencies and create project structure for acceptance tests

- [X] T001 Add pgtestdb dependencies: `go get github.com/peterldowns/pgtestdb github.com/peterldowns/pgtestdb/migrators/sqlmigrator`
- [X] T002 [P] Create tests/acceptance/ directory structure
- [X] T003 [P] Create internal/testutil/ directory structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core test infrastructure that MUST be complete before ANY user story tests can be written

**⚠️ CRITICAL**: No acceptance test implementation can begin until this phase is complete

- [X] T004 Create MockProvider struct implementing llm.Provider interface in internal/llm/mock_provider.go using moq code generation
- [X] T005 Implement MockProvider.Complete() method in internal/llm/mock_provider.go returning configured or default response
- [X] T006 Add Handler() method to transport.Server interface in internal/transport/transport.go
- [X] T007 Implement Handler() method on server struct in internal/transport/transport.go returning http.Handler
- [X] T008 Create NewTestDB helper function in internal/testutil/database.go using pgtestdb with sqlmigrator
- [X] T009 Create NewTestServer helper function in internal/testutil/server.go using wireup pattern with mock injection
- [X] T010 Create pgx pool wrapper in internal/testutil/database.go to convert sql.DB to pgxpool.Pool interface

**Checkpoint**: Foundation ready - all test infrastructure components available for acceptance tests

---

## Phase 3: User Story 1 - Send Message Happy Path (Priority: P1) 🎯 MVP

**Goal**: Verify end-to-end message flow works with mock LLM provider

**Independent Test**: Run `go test ./tests/acceptance/... -run TestSendMessageHappyPath -v`

### Implementation for User Story 1

- [X] T011 [US1] Create TestSendMessageHappyPath function in tests/acceptance/acceptance_test.go
- [X] T012 [US1] In TestSendMessageHappyPath: setup test database using testutil.NewTestDB(t)
- [X] T013 [US1] In TestSendMessageHappyPath: create MockProvider with DefaultResponse set
- [X] T014 [US1] In TestSendMessageHappyPath: start test server using testutil.NewTestServer(t, db, mockLLM)
- [X] T015 [US1] In TestSendMessageHappyPath: create conversation using ConversationServiceClient
- [X] T016 [US1] In TestSendMessageHappyPath: send message using MessageServiceClient
- [X] T017 [US1] In TestSendMessageHappyPath: assert response contains user message and mock assistant response
- [X] T018 [US1] Validate TestSendMessageHappyPath runs successfully: `go test ./tests/acceptance/... -run TestSendMessageHappyPath -v`

**Checkpoint**: User Story 1 complete - can send messages and receive mock responses end-to-end

---

## Phase 4: User Story 2 - Conversation Creation Happy Path (Priority: P2)

**Goal**: Verify conversation creation works independently

**Independent Test**: Run `go test ./tests/acceptance/... -run TestConversationCreation -v`

### Implementation for User Story 2

- [X] T019 [US2] Create TestConversationCreation function in tests/acceptance/acceptance_test.go
- [X] T020 [US2] In TestConversationCreation: setup test database and mock LLM
- [X] T021 [US2] In TestConversationCreation: start test server
- [X] T022 [US2] In TestConversationCreation: create conversation with specific age bracket via ConversationServiceClient
- [X] T023 [US2] In TestConversationCreation: assert response contains valid ID and correct age bracket
- [X] T024 [US2] Validate TestConversationCreation runs successfully: `go test ./tests/acceptance/... -run TestConversationCreation -v`

**Checkpoint**: User Story 2 complete - conversation creation validated independently

---

## Phase 5: User Story 3 - Mock LLM Response Customization (Priority: P3)

**Goal**: Verify mock LLM can return age-bracket-specific responses

**Independent Test**: Run `go test ./tests/acceptance/... -run TestMockResponseCustomization -v`

### Implementation for User Story 3

- [X] T025 [US3] Create TestMockResponseCustomization function in tests/acceptance/acceptance_test.go
- [X] T026 [US3] In TestMockResponseCustomization: setup test database
- [X] T027 [US3] In TestMockResponseCustomization: create MockProvider with Responses map containing age-bracket-specific responses
- [X] T028 [US3] In TestMockResponseCustomization: start test server
- [X] T029 [US3] In TestMockResponseCustomization: create conversation with AgeBracketLittleOnes
- [X] T030 [US3] In TestMockResponseCustomization: send message and assert response matches configured age bracket response
- [X] T031 [US3] In TestMockResponseCustomization: verify MockProvider.CallCount equals 1
- [X] T032 [US3] Validate TestMockResponseCustomization runs successfully: `go test ./tests/acceptance/... -run TestMockResponseCustomization -v`

**Checkpoint**: User Story 3 complete - mock response customization validated

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [X] T033 [P] Run all acceptance tests: `go test ./tests/acceptance/... -v -race`
- [X] T034 [P] Verify tests complete in under 30 seconds (SC-001) - Completed in ~16 seconds
- [X] T035 [P] Verify tests run without external API keys (SC-002)
- [X] T036 Add Taskfile task for acceptance tests in Taskfile.yml: `task test:acceptance`
- [X] T037 Validate quickstart.md scenarios work end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can proceed sequentially in priority order (P1 → P2 → P3)
  - Each story is independently testable after Foundational
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on Foundational - No dependencies on other stories
- **User Story 2 (P2)**: Depends only on Foundational - Independently testable
- **User Story 3 (P3)**: Depends only on Foundational - Independently testable

### Within Each User Story

- Setup test infrastructure (db, mock, server)
- Make API calls via generated clients
- Assert responses match expectations
- Validate test runs successfully

### Parallel Opportunities

**Phase 1 (Setup)**:
```bash
# Can run in parallel:
Task T002: "Create tests/acceptance/ directory structure"
Task T003: "Create internal/testutil/ directory structure"
```

**Phase 2 (Foundational)**:
- T004 and T005 (mock provider) can parallel with T006/T007 (handler) and T008/T009/T010 (test utils)

**Phase 6 (Polish)**:
```bash
# Can run in parallel:
Task T033: "Run all acceptance tests"
Task T034: "Verify tests complete in under 30 seconds"
Task T035: "Verify tests run without external API keys"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (add dependencies, create directories)
2. Complete Phase 2: Foundational (mock provider, test utilities)
3. Complete Phase 3: User Story 1 (TestSendMessageHappyPath)
4. **STOP and VALIDATE**: Run `go test ./tests/acceptance/... -run TestSendMessageHappyPath -v`
5. MVP complete - end-to-end message flow validated

### Incremental Delivery

1. Setup + Foundational → Test infrastructure ready
2. Add User Story 1 → Test message flow → Validate
3. Add User Story 2 → Test conversation creation → Validate
4. Add User Story 3 → Test mock customization → Validate
5. Polish → Add task runner, validate performance

### Parallel Team Strategy

With multiple developers:
1. Dev A: Complete Setup + Foundational (required first)
2. Once Foundational is done:
   - Dev A: User Story 1
   - Dev B: User Story 2
   - Dev C: User Story 3
3. All devs: Polish phase

---

## Notes

- This feature creates TEST infrastructure, not production code
- All acceptance tests should be independently runnable
- MockProvider is generated using moq (`go generate ./internal/llm/...` or `task generate:mocks`)
- pgtestdb provides ~20ms database isolation per test
- Existing wireup pattern enables clean dependency injection
- T006/T007 adds Handler() to expose http.Handler for httptest.Server
