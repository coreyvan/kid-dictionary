package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver with database/sql
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/sqlmigrator"
	migrate "github.com/rubenv/sql-migrate"
)

// NewTestDB creates an isolated test database using pgtestdb.
// Each test gets its own database cloned from a template, providing
// complete isolation in ~20ms.
//
// The database is automatically cleaned up when the test completes.
func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	migrator := sqlmigrator.New(
		&migrate.FileMigrationSource{Dir: "../../migrations"},
		nil, // use default MigrationSet
	)

	conf := pgtestdb.Config{
		DriverName: "pgx",
		Host:       "localhost",
		Port:       "5432",
		User:       "postgres",
		Password:   "postgres",
		Options:    "sslmode=disable",
	}

	// Create the test database - this returns a *sql.DB
	db := pgtestdb.New(t, conf, migrator)

	// Get the database name from the connection
	var dbName string
	err := db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		t.Fatalf("failed to get database name: %v", err)
	}

	// Close the sql.DB connection - we'll use pgxpool instead
	_ = db.Close()

	// Build connection URL for pgxpool
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, dbName,
	)

	pool, err := pgxpool.New(context.Background(), connURL)
	if err != nil {
		t.Fatalf("failed to create pgxpool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
