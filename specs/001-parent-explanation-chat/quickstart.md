# Quickstart: Parent Explanation Chat

**Date**: 2025-12-04
**Branch**: `001-parent-explanation-chat`

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Buf CLI (for proto generation)
- OpenAI API key

## Setup

### 1. Clone and Initialize

```bash
git clone <repo>
cd kid-dictionary
git checkout 001-parent-explanation-chat
task init  # Copies .env.example to .env
```

### 2. Configure Environment

Edit `.env` with your settings:

```bash
# Server
BIND_ADDR=0.0.0.0
LISTEN_PORT=8080

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/kiddictionary?sslmode=disable

# LLM Provider
OPENAI_API_KEY=sk-...

# Logging
PRETTY_LOG=true
LOG_LEVEL=debug
```

### 3. Start PostgreSQL

```bash
# Using Docker
docker run -d \
  --name kiddictionary-db \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=pass \
  -e POSTGRES_DB=kiddictionary \
  -p 5432:5432 \
  postgres:14

# Or use your local PostgreSQL
createdb kiddictionary
```

### 4. Run Migrations

```bash
# Migrations will be in migrations/ directory
# Use golang-migrate or goose to apply
migrate -path migrations -database "$DATABASE_URL" up
```

### 5. Generate Proto Code

```bash
buf generate
```

### 6. Start the Server

```bash
go run cmd/server/main.go
```

Server starts at `http://localhost:8080`

## Usage Examples

### Create a Conversation

```bash
curl -X POST http://localhost:8080/kiddictionary.v1.ConversationService/CreateConversation \
  -H "Content-Type: application/json" \
  -d '{
    "age_bracket": "AGE_BRACKET_LITTLE_ONES"
  }'
```

Response:
```json
{
  "conversation": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "New Conversation",
    "ageBracket": "AGE_BRACKET_LITTLE_ONES",
    "createdAt": "2025-12-04T10:00:00Z"
  }
}
```

### Send a Message

```bash
curl -X POST http://localhost:8080/kiddictionary.v1.MessageService/SendMessage \
  -H "Content-Type: application/json" \
  -d '{
    "conversation_id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Why is the sky blue?"
  }'
```

Response:
```json
{
  "userMessage": {
    "id": "...",
    "role": "MESSAGE_ROLE_USER",
    "content": "Why is the sky blue?"
  },
  "assistantMessage": {
    "id": "...",
    "role": "MESSAGE_ROLE_ASSISTANT",
    "content": "The sky looks blue because of the sunlight..."
  }
}
```

### List Conversations

```bash
curl -X POST http://localhost:8080/kiddictionary.v1.ConversationService/ListConversations \
  -H "Content-Type: application/json" \
  -d '{
    "page_size": 10
  }'
```

### Get Conversation with Messages

```bash
curl -X POST http://localhost:8080/kiddictionary.v1.ConversationService/GetConversation \
  -H "Content-Type: application/json" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

## Testing

### Run Unit Tests

```bash
go test ./...
```

### Run Integration Tests

Requires PostgreSQL and OpenAI API key:

```bash
go test ./internal/llm -tags=integration
go test ./internal/conversation -tags=integration
```

## Validation Checklist

After implementation, verify:

- [ ] Create conversation with each age bracket
- [ ] Send message and receive age-appropriate response
- [ ] Verify Little Ones (0-5) uses simple 2-3 sentence explanations
- [ ] Verify Growing Minds (5-10) uses metaphors and 3-5 sentences
- [ ] Verify Pre-Teens (10+) acknowledges complexity
- [ ] Test sensitive topic (e.g., "Why do people die?") includes guidance prefix
- [ ] Test off-purpose request (e.g., "Write me a poem") is politely declined
- [ ] Verify conversation history is persisted and retrievable
- [ ] Test follow-up questions maintain context
- [ ] Verify response time is under 10 seconds
- [ ] Test error handling when LLM is unavailable
- [ ] Verify 500 character input limit is enforced

## Troubleshooting

### "invalid API key" error
- Check `OPENAI_API_KEY` in `.env`
- Ensure API key has sufficient credits

### Database connection errors
- Verify `DATABASE_URL` format
- Check PostgreSQL is running
- Ensure migrations have been applied

### Proto generation errors
- Run `buf dep update` to update dependencies
- Ensure `buf.yaml` and `buf.gen.yaml` are configured

### Slow responses
- Check OpenAI API status
- Consider reducing conversation context length
- Review rate limiting configuration
