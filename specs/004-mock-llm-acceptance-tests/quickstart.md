# Quickstart: Mock LLM Acceptance Tests

**Branch**: `004-mock-llm-acceptance-tests` | **Date**: 2025-12-05

## Prerequisites

- Docker Desktop, Podman, or compatible container runtime (for PostgreSQL)
- Go 1.25.1+
- Existing environment setup (`task up` has been run at least once)

## Running Acceptance Tests

### 1. Start Database (if not already running)

```bash
# Ensure PostgreSQL is running
task up

# Verify database is accessible
docker compose ps | grep db
```

### 2. Run Acceptance Tests

```bash
# Run all acceptance tests
go test ./tests/acceptance/... -v

# Run specific test
go test ./tests/acceptance/... -v -run TestSendMessageHappyPath

# Run with race detection (recommended for CI)
go test ./tests/acceptance/... -v -race
```

### 3. Verify Results

Expected output for successful run:
```
=== RUN   TestConversationCreation
--- PASS: TestConversationCreation (0.05s)
=== RUN   TestSendMessageHappyPath
--- PASS: TestSendMessageHappyPath (0.08s)
=== RUN   TestMockResponseCustomization
--- PASS: TestMockResponseCustomization (0.06s)
PASS
ok      github.com/coreyvan/kid-dictionary/tests/acceptance    0.230s
```

## Writing New Acceptance Tests

### Basic Test Structure

```go
func TestMyFeature(t *testing.T) {
    // 1. Get isolated test database
    db := testutil.NewTestDB(t)

    // 2. Create mock LLM with expected response
    mockLLM := &llm.MockProvider{
        DefaultResponse: "This is a test response for kids!",
    }

    // 3. Start test server with mocks injected
    server := testutil.NewTestServer(t, db, mockLLM)
    defer server.Close()

    // 4. Create client pointing to test server
    client := kiddictionaryv1connect.NewConversationServiceClient(
        http.DefaultClient,
        server.URL,
    )

    // 5. Make requests and assert
    ctx := context.Background()
    resp, err := client.CreateConversation(ctx, connect.NewRequest(&v1.CreateConversationRequest{
        Title:      "Test",
        AgeBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
    }))

    require.NoError(t, err)
    assert.NotEmpty(t, resp.Msg.Conversation.Id)
}
```

### Customizing Mock Responses

```go
// Per-age-bracket responses
mockLLM := &llm.MockProvider{
    DefaultResponse: "Generic response",
    Responses: map[llm.AgeBracket]string{
        llm.AgeBracketLittleOnes:   "Simple words for little kids!",
        llm.AgeBracketGrowingMinds: "A bit more detail for bigger kids.",
        llm.AgeBracketPreTeens:     "More nuanced explanation.",
    },
}
```

### Verifying Mock Was Called

```go
mockLLM := &llm.MockProvider{DefaultResponse: "Response"}
// ... run test ...
assert.Equal(t, 1, mockLLM.CallCount, "LLM should be called exactly once")
```

## Troubleshooting

### "connection refused" errors

Database not running. Run `task up` first.

### "database does not exist" errors

pgtestdb creates databases on-demand. Ensure:
1. PostgreSQL is running and healthy: `docker compose ps`
2. Postgres user has CREATE DATABASE permission (default docker user does)

### Tests are slow (>30s)

Check if template database needs recreation:
```bash
# Connect to postgres and check templates
docker compose exec db psql -U postgres -c "\l" | grep pgtestdb
```

If many stale templates exist, reset:
```bash
task reset
```

### Mock not being used

Verify wireup injection:
1. Check that `NewTestServer` properly sets `WireupDeps.LLMProvider`
2. Ensure test is using the server URL from `NewTestServer`, not localhost:8080

## CI Integration

Add to CI workflow:

```yaml
- name: Run Acceptance Tests
  run: |
    docker compose up -d db
    sleep 5  # Wait for DB to be ready
    go test ./tests/acceptance/... -v -race
  env:
    DB_HOST: localhost
    DB_PORT: 5432
    DB_USER: postgres
    DB_PASSWORD: postgres
    DB_NAME: kid_dictionary
```

## Available Test Utilities

| Function | Package | Description |
|----------|---------|-------------|
| `NewTestDB(t)` | `testutil` | Creates isolated database via pgtestdb |
| `NewTestServer(t, db, mockLLM)` | `testutil` | Starts httptest.Server with injected mocks |
| `MockProvider` | `llm` | Mock LLM provider with configurable responses |
