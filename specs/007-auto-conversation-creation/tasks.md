# Tasks: Auto Conversation Creation

**Input**: Design documents from `/specs/007-auto-conversation-creation/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are included per Constitution Principle II (Test-Alongside Development).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Based on plan.md structure:
- **Proto definitions**: `proto/kiddictionary/v1/`
- **Generated code**: `gen/`
- **LLM layer**: `internal/llm/`
- **Service layer**: `internal/service/message/`
- **Transport layer**: `internal/transport/`
- **Acceptance tests**: `tests/acceptance/`

---

## Phase 1: Setup (Proto Changes)

**Purpose**: Update API contract before implementation

- [X] T001 Modify SendMessageRequest in proto/kiddictionary/v1/message.proto to add optional age_bracket field
- [X] T002 Modify SendMessageResponse in proto/kiddictionary/v1/message.proto to add optional conversation field
- [X] T003 Add import for conversation.proto in proto/kiddictionary/v1/message.proto if not present
- [X] T004 Run buf generate to regenerate Connect-Go code in gen/

---

## Phase 2: Foundational (LLM Interface Extension)

**Purpose**: Extend LLM provider interface with title generation - MUST complete before user story implementation

**CRITICAL**: User story implementation depends on GenerateTitle being available

- [X] T005 Add GenerateTitle method signature to Provider interface in internal/llm/llm.go
- [X] T006 Add title generation prompt constant in internal/llm/prompts.go
- [X] T007 Implement GenerateTitle for OpenAI provider in internal/llm/openai.go
- [X] T008 Run go generate ./internal/llm/... to regenerate mock with GenerateTitle method

**Checkpoint**: LLM interface ready - user story implementation can now begin

---

## Phase 3: User Story 1 - First Message Creates Conversation (Priority: P1)

**Goal**: Auto-create conversation when SendMessage is called without conversation_id, generating a 3-4 word title via LLM

**Independent Test**: Send a message without conversation_id, verify conversation is created with auto-generated title, and response includes conversation ID

### Tests for User Story 1

- [X] T009 [P] [US1] Add unit test TestSendMessage_AutoCreatesConversation in internal/service/message/service_test.go
- [X] T010 [P] [US1] Add unit test TestSendMessage_TitleGenerationFallback in internal/service/message/service_test.go
- [X] T011 [P] [US1] Add unit test TestSendMessage_RequiresAgeBracketWhenNoConversationID in internal/service/message/service_test.go

### Implementation for User Story 1

- [X] T012 [US1] Add SendMessageWithAutoCreateResult type to internal/domain/message.go
- [X] T013 [US1] Modify SendMessage method signature in internal/service/message/service.go to accept optional ageBracket parameter
- [X] T014 [US1] Implement auto-creation logic in SendMessage: check if conversationID is nil/zero, validate ageBracket, call GenerateTitle in internal/service/message/service.go
- [X] T015 [US1] Implement fallback title logic when GenerateTitle fails in internal/service/message/service.go
- [X] T016 [US1] Add structured logging for title generation events (success/failure, latency, conversation_id) in internal/service/message/service.go
- [X] T017 [US1] Create conversation via conversationRepo.Create after title generation in internal/service/message/service.go
- [X] T018 [US1] Update SendMessage transport handler to extract age_bracket from request in internal/transport/message.go
- [X] T019 [US1] Update SendMessage transport handler to include conversation in response when auto-created in internal/transport/message.go

### Acceptance Test for User Story 1

- [X] T020 [US1] Add acceptance test TestSendMessage_AutoCreate_E2E in tests/acceptance/message_test.go

**Checkpoint**: User Story 1 complete - first message auto-creates conversation with LLM-generated title

---

## Phase 4: User Story 2 - Subsequent Messages Use Returned Conversation ID (Priority: P2)

**Goal**: Ensure follow-up messages work correctly with the conversation ID returned from auto-creation

**Independent Test**: Auto-create a conversation, capture the returned ID, send follow-up message with that ID, verify context is maintained

### Tests for User Story 2

- [X] T021 [P] [US2] Add unit test TestSendMessage_FollowUpWithAutoCreatedConversation in internal/service/message/service_test.go

### Implementation for User Story 2

- [X] T022 [US2] Verify existing SendMessage flow handles conversation_id correctly after auto-creation in internal/service/message/service.go (may be no-op if existing logic already works)

### Acceptance Test for User Story 2

- [X] T023 [US2] Add acceptance test TestSendMessage_FollowUp_AfterAutoCreate in tests/acceptance/message_test.go

**Checkpoint**: User Story 2 complete - follow-up messages work correctly with auto-created conversations

---

## Phase 5: User Story 3 - Backwards Compatibility (Priority: P3)

**Goal**: Ensure existing clients that explicitly create conversations continue to work unchanged

**Independent Test**: Create conversation via CreateConversation, send message with that ID, verify no auto-creation occurs

### Tests for User Story 3

- [X] T024 [P] [US3] Add unit test TestSendMessage_ExistingConversation_NoAutoCreate in internal/service/message/service_test.go
- [X] T025 [P] [US3] Add unit test TestSendMessage_ExistingConversation_IgnoresAgeBracket in internal/service/message/service_test.go

### Implementation for User Story 3

- [X] T026 [US3] Verify SendMessage skips title generation when valid conversation_id provided in internal/service/message/service.go (may be no-op if logic correct)

### Acceptance Test for User Story 3

- [X] T027 [US3] Add acceptance test TestSendMessage_ExistingFlow_Unchanged in tests/acceptance/message_test.go

**Checkpoint**: User Story 3 complete - backwards compatibility verified

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [X] T028 [P] Run all unit tests: go test ./internal/...
- [X] T029 [P] Run all acceptance tests: go test ./tests/acceptance/...
- [X] T030 [P] Run linter: golangci-lint run
- [X] T031 Run buf lint to validate proto changes
- [X] T032 Verify quickstart.md checklist items are complete
- [X] T033 Update CLAUDE.md if any new technologies or patterns introduced (no new technologies)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - proto changes first
- **Foundational (Phase 2)**: Depends on Phase 1 (needs generated proto types) - BLOCKS all user stories
- **User Stories (Phases 3-5)**: All depend on Phase 2 (need GenerateTitle method)
  - User stories can proceed sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Core auto-creation feature
- **User Story 2 (P2)**: Can start after US1 - Verifies follow-up flow works with auto-created conversations
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of US1/US2, verifies backwards compat

### Within Each User Story

- Tests written before/alongside implementation
- Service layer before transport layer
- Core implementation before logging/validation refinements

### Parallel Opportunities

- T001, T002, T003 are sequential (same file)
- T009, T010, T011 can run in parallel (different test functions)
- T021 is independent
- T024, T025 can run in parallel (different test functions)
- T028, T029, T030 can run in parallel (different tools)

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all US1 tests in parallel:
Task: "Add unit test TestSendMessage_AutoCreatesConversation in internal/service/message/service_test.go"
Task: "Add unit test TestSendMessage_TitleGenerationFallback in internal/service/message/service_test.go"
Task: "Add unit test TestSendMessage_RequiresAgeBracketWhenNoConversationID in internal/service/message/service_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (proto changes + buf generate)
2. Complete Phase 2: Foundational (LLM interface extension)
3. Complete Phase 3: User Story 1 (auto-creation + title generation)
4. **STOP and VALIDATE**: Test auto-creation flow end-to-end
5. Deploy/demo if ready - users can now skip conversation creation step

### Incremental Delivery

1. Complete Setup + Foundational → LLM interface ready
2. Add User Story 1 → Test independently → Core feature working
3. Add User Story 2 → Verify follow-ups work → Full conversation flow
4. Add User Story 3 → Verify backwards compat → Safe for existing clients
5. Each story adds confidence without breaking previous functionality

---

## Notes

- [P] tasks = different files or independent operations
- [Story] label maps task to specific user story for traceability
- Constitution Principle II requires tests alongside implementation
- FR-010 requires title generation logging (covered in T016)
- Fallback title "New Conversation" per spec edge case
- Proto changes are backwards compatible (new optional fields)
