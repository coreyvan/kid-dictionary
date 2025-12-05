//go:build integration

package conversation_test

import (
	"context"
	"os"
	"testing"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/repository/conversation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestPostgresRepository_CreateAndGet(t *testing.T) {
	pool := setupTestDB(t)
	repo := conversation.NewPostgresRepository(pool)
	ctx := context.Background()

	// Create a conversation
	conv := &domain.Conversation{
		Title:      "Test Conversation",
		AgeBracket: domain.AgeBracketLittleOnes,
	}

	err := repo.Create(ctx, conv)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, conv.ID)
	assert.False(t, conv.CreatedAt.IsZero())

	// Get the conversation
	got, err := repo.GetByID(ctx, conv.ID)
	require.NoError(t, err)
	assert.Equal(t, conv.ID, got.ID)
	assert.Equal(t, conv.Title, got.Title)
	assert.Equal(t, conv.AgeBracket, got.AgeBracket)

	// Cleanup
	err = repo.Delete(ctx, conv.ID)
	require.NoError(t, err)
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := conversation.NewPostgresRepository(pool)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPostgresRepository_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := conversation.NewPostgresRepository(pool)
	ctx := context.Background()

	// Create multiple conversations
	var ids []uuid.UUID
	for i := 0; i < 3; i++ {
		conv := &domain.Conversation{
			Title:      "Test Conversation",
			AgeBracket: domain.AgeBracket(i + 1),
		}
		err := repo.Create(ctx, conv)
		require.NoError(t, err)
		ids = append(ids, conv.ID)
	}

	// List conversations
	convs, err := repo.List(ctx, nil, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(convs), 3)

	// Cleanup
	for _, id := range ids {
		err := repo.Delete(ctx, id)
		require.NoError(t, err)
	}
}

func TestPostgresRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	repo := conversation.NewPostgresRepository(pool)
	ctx := context.Background()

	// Create a conversation
	conv := &domain.Conversation{
		Title:      "Original Title",
		AgeBracket: domain.AgeBracketLittleOnes,
	}
	err := repo.Create(ctx, conv)
	require.NoError(t, err)

	// Update the conversation
	conv.Title = "Updated Title"
	conv.AgeBracket = domain.AgeBracketPreTeens
	err = repo.Update(ctx, conv)
	require.NoError(t, err)

	// Verify the update
	got, err := repo.GetByID(ctx, conv.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", got.Title)
	assert.Equal(t, domain.AgeBracketPreTeens, got.AgeBracket)

	// Cleanup
	err = repo.Delete(ctx, conv.ID)
	require.NoError(t, err)
}

func TestPostgresRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := conversation.NewPostgresRepository(pool)
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
