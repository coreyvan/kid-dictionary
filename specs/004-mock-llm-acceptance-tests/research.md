# Research: Mock LLM Acceptance Tests

**Branch**: `004-mock-llm-acceptance-tests` | **Date**: 2025-12-05

## pgtestdb Integration

**Decision**: Use pgtestdb with sqlmigrator for per-test database isolation

**Rationale**:
- Each test gets a completely isolated database in ~20ms (template database cloning)
- Works with existing sql-migrate format migrations (`-- +migrate Up/Down`)
- Concurrency-safe via `migrate.MigrationSet` (no global state)
- Uses existing Docker Compose PostgreSQL server

**Configuration**:
```go
import (
    "github.com/peterldowns/pgtestdb"
    "github.com/peterldowns/pgtestdb/migrators/sqlmigrator"
    "github.com/rubenv/sql-migrate"
)

func NewTestDB(t *testing.T) *sql.DB {
    t.Helper()

    migrator := sqlmigrator.New(
        &migrate.FileMigrationSource{Dir: "../../migrations"},
        nil, // default MigrationSet
    )

    return pgtestdb.New(t, pgtestdb.Config{
        DriverName: "pgx",
        Host:       "localhost",
        Port:       "5432",
        User:       "postgres",
        Password:   "postgres",
        Options:    "sslmode=disable",
    }, migrator)
}
```

**Alternatives Considered**:
- **Transaction rollback**: Simpler but doesn't test actual commit behavior
- **Truncate tables**: Requires ordering for foreign keys, doesn't reset sequences
- **Docker container per test**: Too slow (~seconds vs ~20ms)

## Mock LLM Provider Pattern

**Decision**: Simple struct implementing `llm.Provider` interface with configurable responses

**Rationale**:
- Interface already defined at `internal/llm/llm.go:48`
- Wireup pattern at `internal/config/wireup.go` already supports injection via `WireupDeps.LLMProvider`
- No need for complex mocking frameworks - simple struct with map of responses

**Implementation Pattern**:
```go
type MockProvider struct {
    DefaultResponse string
    Responses       map[llm.AgeBracket]string
    CallCount       int
}

func (m *MockProvider) Complete(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
    m.CallCount++
    content := m.DefaultResponse
    if resp, ok := m.Responses[req.AgeBracket]; ok {
        content = resp
    }
    return llm.CompletionResponse{Content: content, TokensUsed: 10}, nil
}
```

**Alternatives Considered**:
- **testify/mock**: Overkill for simple interface; adds dependency
- **gomock**: Code generation overhead; simple struct sufficient
- **httptest.Server mocking OpenAI**: Tests OpenAI client, not our business logic

## Test Server Setup Pattern

**Decision**: Create httptest.Server using same wireup as production, with mocks injected

**Rationale**:
- Constitution Principle III (Simplicity): Reuse existing wireup pattern
- Tests exercise same code paths as production
- httptest.Server handles port allocation automatically
- Matches clarification decision to "use wireup the same way we're using in the main entrypoint"

**Implementation Pattern**:
```go
func NewTestServer(t *testing.T, db *sql.DB, mockLLM llm.Provider) *httptest.Server {
    t.Helper()

    pool := WrapDBAsPool(db) // Convert sql.DB to pgx pool interface

    cfg := config.Config{
        BindAddr: "127.0.0.1",
        Port:     "0", // Let httptest choose
    }

    wireup := config.NewWiring(cfg, slog.Default())
    wireup.(*config.Wireup).Deps = config.WireupDeps{
        DBPool:      pool,
        LLMProvider: mockLLM,
    }

    server := wireup.MustProvideServer()
    return httptest.NewServer(server.Handler())
}
```

**Challenge**: `transport.Server` doesn't expose handler directly. Need to either:
1. Add `Handler()` method to Server interface (preferred - minimal change)
2. Create server and extract handler before Listen()

**Alternatives Considered**:
- **Docker container with mock**: Constitution Principle III violation; adds complexity
- **Separate test binary**: Harder to inject mocks; slower test iteration

## Generated Client Usage

**Decision**: Use existing Connect-Go generated clients from `gen/kiddictionary/v1/kiddictionaryv1connect/`

**Rationale**:
- Already generated and committed per Constitution quality gates
- Provides type-safe request/response handling
- Matches FR-003 requirement

**Usage Pattern**:
```go
client := kiddictionaryv1connect.NewConversationServiceClient(
    http.DefaultClient,
    testServer.URL,
)

resp, err := client.CreateConversation(ctx, connect.NewRequest(&v1.CreateConversationRequest{
    Title:      "Test Conversation",
    AgeBracket: v1.AgeBracket_AGE_BRACKET_LITTLE_ONES,
}))
```

## Test File Organization

**Decision**: Acceptance tests in `tests/acceptance/` package

**Rationale**:
- Separates acceptance tests from unit tests
- Clear naming convention
- Can run independently: `go test ./tests/acceptance/...`

**File Structure**:
```
tests/
└── acceptance/
    ├── acceptance_test.go    # Main acceptance tests
    └── testutil_test.go      # Test helpers (NewTestDB, NewTestServer)
```

**Alternatives Considered**:
- **`internal/acceptance/`**: Acceptance tests shouldn't be internal
- **`*_acceptance_test.go` suffix**: Less discoverable; harder to run separately

## Dependencies to Add

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/peterldowns/pgtestdb` | latest | Per-test database isolation |
| `github.com/peterldowns/pgtestdb/migrators/sqlmigrator` | latest | sql-migrate integration |
| `github.com/rubenv/sql-migrate` | latest | Migration execution (indirect, already used) |

## Open Issue: Server Handler Access

The current `transport.Server` interface doesn't expose the HTTP handler. Two options:

1. **Add Handler() method** (recommended):
   ```go
   type Server interface {
       Listen(addr string) error
       Ready() chan struct{}
       Shutdown(ctx context.Context) error
       Addr() (string, error)
       Handler() http.Handler  // New
   }
   ```

2. **Build handler without Listen()**:
   Create handler in wireup, wrap in httptest.Server without calling Listen()

Decision deferred to implementation - both are simple. Option 1 is cleaner.
