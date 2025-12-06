# Contracts: Mock LLM Acceptance Tests

This feature does not define new API contracts.

**Reason**: The acceptance tests use existing generated Connect-Go clients from:
- `gen/kiddictionary/v1/kiddictionaryv1connect/conversation.connect.go`
- `gen/kiddictionary/v1/kiddictionaryv1connect/message.connect.go`

These clients were generated from the existing Protocol Buffer definitions and are already committed to the repository per the project's quality gates.

## Existing Endpoints Used

| Endpoint | Client Method | Used For |
|----------|---------------|----------|
| `CreateConversation` | `ConversationServiceClient.CreateConversation()` | US2: Create conversation tests |
| `SendMessage` | `MessageServiceClient.SendMessage()` | US1: Message flow tests |
| `GetConversation` | `ConversationServiceClient.GetConversation()` | Verification queries |

See `gen/kiddictionary/v1/*.proto` for the source definitions.
