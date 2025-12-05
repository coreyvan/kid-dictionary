# Tasks: Parent Explanation Chat

**Input**: Design documents from `/specs/001-parent-explanation-chat/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Tests are included per constitution principle "II. Test-Alongside Development".

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go project: `internal/`, `cmd/`, `proto/`, `gen/`
- Existing infrastructure: `internal/config/`, `internal/transport/`, `internal/llm/`
- New packages: `internal/conversation/`, `internal/message/`, `internal/connect/`

---

## Phase 1: Setup

**Purpose**: Proto modifications and code generation

- [x] T001 Add ContentTier enum to proto/kiddictionary/v1/common.proto
- [x] T002 Add content_tier field to Message in proto/kiddictionary/v1/common.proto
- [x] T003 Run buf generate to update gen/kiddictionary/v1/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Database schema and shared infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Create database migration for conversations table in migrations/001_create_conversations.sql
- [x] T005 Create database migration for messages table in migrations/002_create_messages.sql
- [x] T006 [P] Create conversation repository interface in internal/conversation/repository.go
- [x] T007 [P] Create message repository interface in internal/message/repository.go
- [x] T008 Implement PostgreSQL conversation repository in internal/conversation/postgres.go
- [x] T009 Implement PostgreSQL message repository in internal/message/postgres.go
- [x] T010 [P] Add repository providers to internal/config/wireup.go
- [x] T011 [P] Write integration test for conversation repository in internal/conversation/postgres_test.go
- [x] T012 [P] Write integration test for message repository in internal/message/postgres_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Ask a Question (Priority: P1)

**Goal**: A parent can ask a concept question and receive an age-appropriate explanation

**Independent Test**: Enter a concept like "why is the sky blue?" with age bracket "Little Ones" and verify the response uses simple, concrete words in 2-3 sentences

### Implementation for User Story 1

- [x] T013 [US1] Create conversation service with CreateConversation in internal/conversation/service.go
- [x] T014 [US1] Create message service with SendMessage in internal/message/service.go
- [x] T015 [US1] Extend system prompts with classification instructions in internal/llm/prompts.go
- [x] T016 [US1] Implement Connect ConversationService handler (CreateConversation only) in internal/connect/conversation.go
- [x] T017 [US1] Implement Connect MessageService handler (SendMessage) in internal/connect/message.go
- [x] T018 [US1] Wire Connect handlers into Chi router in internal/transport/transport.go
- [x] T019 [US1] Add input validation (500 char limit, non-empty) in internal/message/service.go
- [x] T020 [US1] Add error mapping to Connect status codes in internal/connect/errors.go
- [x] T021 [P] [US1] Write unit test for conversation service in internal/conversation/service_test.go
- [x] T022 [P] [US1] Write unit test for message service in internal/message/service_test.go

**Checkpoint**: User Story 1 complete - parents can ask questions and receive age-appropriate explanations

---

## Phase 4: User Story 2 - Select Age Bracket (Priority: P2)

**Goal**: Parents can select different age brackets and receive appropriately tailored responses

**Independent Test**: Ask the same question ("why do people die?") with each age bracket and verify response complexity varies appropriately

### Implementation for User Story 2

- [x] T023 [US2] Add age bracket validation to conversation service in internal/conversation/service.go
- [x] T024 [US2] Ensure prompts.go correctly returns different prompts per bracket in internal/llm/prompts.go
- [x] T025 [US2] Add UpdateConversation RPC to allow bracket changes in proto/kiddictionary/v1/conversation.proto
- [x] T026 [US2] Implement UpdateConversation in conversation service in internal/conversation/service.go
- [x] T027 [US2] Implement UpdateConversation Connect handler in internal/connect/conversation.go
- [x] T028 [P] [US2] Write test verifying different responses per age bracket in internal/message/service_test.go

**Checkpoint**: User Story 2 complete - parents can select and switch age brackets

---

## Phase 5: User Story 3 - Handle Sensitive Topics (Priority: P3)

**Goal**: System provides appropriate handling for sensitive, contextual, and harmful content

**Independent Test**: Enter sensitive topics ("divorce", "death") and verify response includes soft guidance prefix; enter off-purpose request ("write me a poem") and verify polite decline

### Implementation for User Story 3

- [x] T029 [US3] Create content classifier with tier detection in internal/message/classifier.go
- [x] T030 [US3] Add off-purpose request detection to classifier in internal/message/classifier.go
- [x] T031 [US3] Implement soft guidance prefix for Tier 2 (Sensitive) in internal/message/service.go
- [x] T032 [US3] Implement framing preference request for Tier 3 (Contextual) in internal/message/service.go
- [x] T033 [US3] Implement polite decline for Tier 4 (Redirect/Off-purpose) in internal/message/service.go
- [x] T034 [US3] Store content_tier in message records via repository in internal/message/service.go
- [x] T035 [P] [US3] Write unit test for content classifier in internal/message/classifier_test.go
- [x] T036 [P] [US3] Write test for sensitive topic handling in internal/message/service_test.go
- [x] T037 [P] [US3] Write test for off-purpose request handling in internal/message/service_test.go

**Checkpoint**: User Story 3 complete - system handles sensitive topics and off-purpose requests appropriately

---

## Phase 6: User Story 4 - Conversation History (Priority: P4)

**Goal**: Parents can view conversation history, continue conversations, and ask follow-up questions with context

**Independent Test**: Ask multiple questions in a conversation, retrieve conversation with GetConversation, verify all messages appear in order; ask follow-up question and verify AI considers context

### Implementation for User Story 4

- [x] T038 [US4] Implement GetConversation in conversation service in internal/conversation/service.go
- [x] T039 [US4] Implement ListConversations in conversation service in internal/conversation/service.go
- [x] T040 [US4] Implement DeleteConversation in conversation service in internal/conversation/service.go
- [x] T041 [US4] Add GetConversation Connect handler in internal/connect/conversation.go
- [x] T042 [US4] Add ListConversations Connect handler in internal/connect/conversation.go
- [x] T043 [US4] Add DeleteConversation Connect handler in internal/connect/conversation.go
- [x] T044 [US4] Fetch conversation history for LLM context in internal/message/service.go
- [x] T045 [US4] Limit context to last 10 messages for token management in internal/message/service.go
- [x] T046 [P] [US4] Write test for conversation retrieval in internal/conversation/service_test.go
- [x] T047 [P] [US4] Write test for follow-up question context in internal/message/service_test.go

**Checkpoint**: User Story 4 complete - parents can access conversation history and continue conversations

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Error handling, observability, and cleanup

- [x] T048 [P] Add structured logging with request_id, duration to Connect handlers in internal/connect/
- [x] T049 [P] Add LLM call duration metrics in internal/llm/openai.go
- [x] T050 Add timeout handling (10s) for LLM calls in internal/message/service.go
- [x] T051 Add user-friendly error messages for all error cases in internal/connect/errors.go
- [x] T052 [P] Run quickstart.md validation checklist in specs/001-parent-explanation-chat/quickstart.md
- [x] T053 Update CLAUDE.md with any new patterns discovered during implementation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Phase 2 completion
  - US1 (Phase 3): No story dependencies - MVP
  - US2 (Phase 4): Builds on US1 infrastructure
  - US3 (Phase 5): Builds on US1 infrastructure
  - US4 (Phase 6): Builds on US1 infrastructure
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories - **MVP**
- **User Story 2 (P2)**: Can start after Phase 2, minimal dependency on US1 (uses same service structure)
- **User Story 3 (P3)**: Can start after Phase 2, extends US1 message service
- **User Story 4 (P4)**: Can start after Phase 2, extends US1 conversation service

### Within Each User Story

- Models/repositories before services
- Services before Connect handlers
- Core implementation before tests
- Tests can run in parallel with each other

### Parallel Opportunities

- Phase 1: T001-T002 sequential (same file), T003 depends on both
- Phase 2: T006-T007 parallel, T008-T009 sequential per repo, T010-T012 parallel
- Phase 3: T013-T020 mostly sequential, T021-T022 parallel
- Phase 4: T023-T027 sequential, T028 parallel after T027
- Phase 5: T029-T034 mostly sequential, T035-T037 parallel
- Phase 6: T038-T045 mostly sequential, T046-T047 parallel
- Phase 7: T048-T049 parallel, T050-T053 can run in any order

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (proto changes)
2. Complete Phase 2: Foundational (database, repositories)
3. Complete Phase 3: User Story 1 (ask questions, get explanations)
4. **STOP and VALIDATE**: Test with curl commands per quickstart.md
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (**MVP!**)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Polish → Final validation → Production ready

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Existing LLM provider infrastructure (internal/llm/) is reused - no new abstractions needed
