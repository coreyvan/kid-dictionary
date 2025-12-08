# Feature Specification: Auto Conversation Creation

**Feature Branch**: `007-auto-conversation-creation`
**Created**: 2025-12-07
**Status**: Draft
**Input**: User description: "Alter the conversation creation/message send flow so that a new conversation is automatically created when no conversation ID is provided. The first message triggers LLM-based title generation (3-4 words summarizing the context), and subsequent messages use the conversation ID returned in the response."

## Clarifications

### Session 2025-12-07

- Q: Should title and explanation generation happen sequentially or in parallel? → A: Sequential - generate title first, create conversation, then generate explanation.
- Q: Should title generation have separate observability (logging/metrics)? → A: Log title generation as a distinct event with success/failure and latency.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - First Message Creates Conversation (Priority: P1)

A parent opens the app and immediately types a question without having to create a conversation first. The system automatically creates a conversation, generates a meaningful title from their question, and returns the AI response along with the new conversation ID for subsequent messages.

**Why this priority**: This is the core UX improvement - eliminating the friction of manual conversation creation before asking the first question. It directly addresses the user's request and provides immediate value.

**Independent Test**: Can be fully tested by sending a message without a conversation ID and verifying that a conversation is created with an auto-generated title, and the response includes both the AI explanation and the new conversation ID.

**Acceptance Scenarios**:

1. **Given** a parent has selected an age bracket, **When** they send a message without providing a conversation ID, **Then** the system creates a new conversation with an auto-generated title and returns the AI response along with the new conversation ID.

2. **Given** a parent sends the first message "Why do leaves change color?", **When** the system processes the request, **Then** the conversation is created with a short, contextual title like "Leaves Changing Color" (3-4 words).

3. **Given** the title generation request is sent to the LLM, **When** the response is received, **Then** the title is exactly 3-4 words that capture the essence of the user's question.

4. **Given** a parent sends a first message, **When** the system creates the conversation, **Then** the response includes the new conversation ID that the client can use for subsequent messages.

---

### User Story 2 - Subsequent Messages Use Returned Conversation ID (Priority: P2)

After the first message creates a conversation, the parent continues the conversation by sending follow-up messages using the conversation ID returned from the first response. The flow remains unchanged from the existing behavior.

**Why this priority**: This ensures continuity of the conversation experience. Without this, the auto-creation feature would break the follow-up question flow.

**Independent Test**: Can be tested by first sending a message without conversation ID, capturing the returned conversation ID, and then sending a follow-up message with that ID to verify context is maintained.

**Acceptance Scenarios**:

1. **Given** a parent has received a response with a conversation ID from their first message, **When** they send a follow-up message with that conversation ID, **Then** the system processes it as a continuation of the same conversation with full context.

2. **Given** a conversation was auto-created, **When** the parent asks a follow-up question, **Then** the AI response considers the previous messages in the conversation for context.

---

### User Story 3 - Backwards Compatibility (Priority: P3)

Existing clients that explicitly create conversations before sending messages continue to work without modification. The auto-creation only triggers when no conversation ID is provided.

**Why this priority**: Ensures existing integrations and client implementations are not broken by this change.

**Independent Test**: Can be tested by using the existing flow (CreateConversation followed by SendMessage with conversation ID) and verifying it works identically to before.

**Acceptance Scenarios**:

1. **Given** a client creates a conversation using CreateConversation, **When** they send a message with that conversation ID, **Then** the system processes it exactly as before (no auto-creation occurs).

2. **Given** a client provides a valid conversation ID, **When** SendMessage is called, **Then** the title generation step is skipped and the message is added to the existing conversation.

---

### Edge Cases

- What happens when the LLM title generation fails? The system uses a fallback title of "New Conversation" and continues processing the message normally.
- What happens when the first message is empty or only whitespace? The system returns an InvalidArgument error before attempting conversation creation.
- What happens when the first message exceeds 500 characters? Standard validation applies - return InvalidArgument error.
- What happens when both conversation_id and age_bracket are missing? The system returns an InvalidArgument error indicating age_bracket is required for new conversations.
- What happens if the generated title exceeds 4 words? The system truncates to the first 4 words.
- What happens if the generated title is fewer than 3 words? The title is accepted as-is (minimum is a guideline, not a strict requirement).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept SendMessage requests without a conversation_id field.
- **FR-002**: When conversation_id is not provided, system MUST create a new conversation automatically before processing the message.
- **FR-003**: System MUST require age_bracket in SendMessage requests when conversation_id is not provided.
- **FR-004**: System MUST generate a 3-4 word title for auto-created conversations using an LLM call with the first message content.
- **FR-005**: System MUST return the newly created conversation_id in the SendMessage response when auto-creation occurs.
- **FR-006**: System MUST use a fallback title of "New Conversation" when LLM title generation fails.
- **FR-007**: System MUST continue to support explicit conversation_id in SendMessage requests (backwards compatibility).
- **FR-008**: When a valid conversation_id is provided, system MUST NOT trigger title generation and MUST add the message to the existing conversation.
- **FR-009**: System MUST include the conversation object (with id and title) in the SendMessage response when a new conversation is created.
- **FR-010**: System MUST log title generation events as distinct log entries including: success/failure status, latency, and conversation ID.

### Key Entities

- **Conversation** (existing): Extended to support auto-creation via SendMessage. Title field populated by LLM-generated content rather than client-provided value.
- **SendMessageRequest** (modified): conversation_id becomes optional; age_bracket becomes conditionally required (required when conversation_id is absent).
- **SendMessageResponse** (modified): Includes optional conversation field to return the newly created conversation details.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Parents can start asking questions within 5 seconds of opening the app (reduced from the current flow requiring conversation creation first).
- **SC-002**: Auto-generated titles accurately reflect the question content in 90% of cases (based on user feedback or quality review).
- **SC-003**: The end-to-end latency for first message (including title generation) remains under 15 seconds.
- **SC-004**: 100% of existing client integrations continue to function without modification.
- **SC-005**: Title generation failures are gracefully handled with fallback in 100% of error cases.

## Assumptions

- The LLM provider used for explanations can also be used for title generation with a simple prompt.
- Title generation adds minimal latency (estimated 1-3 seconds) that is acceptable for the improved UX.
- Clients will update to use the returned conversation_id for subsequent messages (this is a client-side responsibility).
- The same age_bracket validation rules apply whether conversation is created explicitly or auto-created.
- Title generation and explanation generation MUST be performed sequentially: title first, then conversation creation, then explanation generation. This ensures data integrity and simplifies error handling.
