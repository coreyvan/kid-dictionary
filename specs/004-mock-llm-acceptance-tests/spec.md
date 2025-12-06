# Feature Specification: Mock LLM Acceptance Tests

**Feature Branch**: `004-mock-llm-acceptance-tests`
**Created**: 2025-12-05
**Status**: Draft
**Input**: User description: "Create a mock implementation of the interface we're using to retrieve responses from LLM's. The goal is to create an acceptance test where we can inject a mock LLM interface, make a call (with a client generated from the openAPI spec) to the running service and ensure that we get happy path successes back."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Send Message Happy Path (Priority: P1)

A developer runs acceptance tests to verify that when a user sends a message to an existing conversation, the service returns both the user's message and an assistant response without requiring a real LLM provider.

**Why this priority**: This is the core functionality - validating the primary message flow end-to-end. Without this working, no other acceptance tests have value.

**Independent Test**: Can be fully tested by creating a conversation, sending a message via the generated API client, and verifying the response contains both user and assistant messages with expected content.

**Acceptance Scenarios**:

1. **Given** a running service with mock LLM provider injected, **When** a test creates a conversation and sends a message via the API client, **Then** the response contains the user's message and a predictable assistant response.
2. **Given** a running service with mock LLM provider configured to return specific content, **When** a test sends a message, **Then** the assistant response matches the configured mock content.

---

### User Story 2 - Conversation Creation Happy Path (Priority: P2)

A developer runs acceptance tests to verify that creating a new conversation via the API returns a valid conversation with expected properties.

**Why this priority**: Conversations are a prerequisite to sending messages, but this can be tested independently to isolate failures.

**Independent Test**: Can be fully tested by calling the create conversation endpoint and verifying the response contains a valid ID and expected properties.

**Acceptance Scenarios**:

1. **Given** a running service with test database, **When** a test creates a conversation with a specified age bracket, **Then** the response contains a valid conversation ID and the specified age bracket.

---

### User Story 3 - Mock LLM Response Customization (Priority: P3)

A developer configures the mock LLM to return specific responses for different test scenarios, enabling targeted testing of edge cases.

**Why this priority**: Enables more sophisticated testing once basic happy path is validated.

**Independent Test**: Can be fully tested by configuring mock responses and verifying the service returns those exact responses.

**Acceptance Scenarios**:

1. **Given** a mock LLM configured with a specific response for age bracket "Little Ones", **When** a test sends a message to a conversation with that age bracket, **Then** the assistant response matches the configured mock response.

---

### Edge Cases

- What happens when the mock LLM is not configured? (Should return a sensible default response)
- What happens when a message is sent to a non-existent conversation? (Should return appropriate error)
- Each test receives a fresh isolated database via pgtestdb, eliminating cross-test state contamination

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a mock implementation of the LLM Provider interface that can be injected at test time
- **FR-002**: System MUST allow acceptance tests to start an in-process HTTP test server using the same wireup pattern as the production entrypoint, with mock dependencies injected via WireupDeps
- **FR-003**: System MUST generate an API client from the service's protocol buffer definitions for use in tests
- **FR-004**: Acceptance tests MUST be able to create conversations via the generated client
- **FR-005**: Acceptance tests MUST be able to send messages via the generated client and receive responses
- **FR-006**: The mock LLM provider MUST return deterministic, predictable responses for test assertions
- **FR-007**: Acceptance tests MUST run against an actual PostgreSQL database using pgtestdb for per-test database isolation
- **FR-008**: System MUST allow configuring the mock LLM's response content for different test scenarios
- **FR-009**: Each acceptance test MUST receive a fresh, isolated database instance to prevent cross-test contamination

### Key Entities

- **Mock LLM Provider**: A test double implementing the LLM Provider interface that returns configurable, predictable responses
- **Test Server**: An in-process HTTP test server instantiated using the production wireup pattern with mock dependencies injected, mirroring deployed service initialization
- **Generated Client**: Protocol buffer-generated client code used to make requests to the test server

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Acceptance tests complete successfully in under 30 seconds when run locally
- **SC-002**: Developers can run acceptance tests without requiring external API keys or network access to LLM providers
- **SC-003**: Tests provide clear failure messages when assertions fail, identifying which response field was unexpected
- **SC-004**: New team members can run acceptance tests within 5 minutes of cloning the repository (using existing Docker setup)

## Assumptions

- The existing `llm.Provider` interface at `internal/llm/llm.go` will be used as the contract for the mock implementation
- The generated Connect-Go client from the protocol buffers will be used for making test requests
- Tests will use the existing Docker Compose setup for database access (PostgreSQL)
- The wireup pattern in `internal/config/wireup.go` already supports dependency injection via the `WireupDeps` struct, allowing mock injection

## Clarifications

### Session 2025-12-05

- Q: How should test database state be managed between test runs? → A: Fresh database per test using pgtestdb (github.com/peterldowns/pgtestdb) to spin up isolated PostgreSQL databases on the shared server
- Q: How should the test server be instantiated for acceptance tests? → A: In-process HTTP test server using the same wireup pattern as the main entrypoint, with mocks injected via WireupDeps
