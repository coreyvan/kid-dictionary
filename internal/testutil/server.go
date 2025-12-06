package testutil

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/coreyvan/kid-dictionary/internal/config"
	"github.com/coreyvan/kid-dictionary/internal/llm"
)

// NewTestServer creates an httptest.Server with the production handlers
// but with injected test dependencies (database pool and LLM provider).
//
// The server is automatically closed when the test completes.
func NewTestServer(t *testing.T, pool *pgxpool.Pool, mockLLM llm.Provider) *httptest.Server {
	t.Helper()

	cfg := config.Config{
		BindAddr: "127.0.0.1",
		Port:     "0", // Let httptest choose
	}

	logger := slog.Default()
	wireup := config.NewWiring(cfg, *logger).(*config.Wireup)

	// Inject test dependencies
	wireup.Deps.DBPool = pool
	wireup.Deps.LLMProvider = mockLLM

	// Build the server with injected dependencies
	server := wireup.MustProvideServer()

	// Create httptest.Server using the handler
	ts := httptest.NewServer(server.Handler())

	t.Cleanup(func() {
		ts.Close()
	})

	return ts
}
