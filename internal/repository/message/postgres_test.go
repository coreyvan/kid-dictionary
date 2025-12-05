//go:build integration

package message_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/repository/message"
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

func createTestConversation(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	id := uuid.New()
	_, err := pool.Exec(ctx, `
		INSERT INTO conversations (id, title, age_bracket, created_at, updated_at)
		VALUES ($1, 'Test Conversation', 1, NOW(), NOW())
	`, id)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM conversations WHERE id = $1`, id)
	})

	return id
}

func TestPostgresRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := message.NewPostgresRepository(pool)
	ctx := context.Background()

	convID := createTestConversation(t, pool)

	msg := &domain.Message{
		ConversationID: convID,
		Role:           domain.RoleUser,
		Content:        "Why is the sky blue?",
	}

	err := repo.Create(ctx, msg)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, msg.ID)
	assert.False(t, msg.CreatedAt.IsZero())
}

func TestPostgresRepository_CreateWithContentTier(t *testing.T) {
	pool := setupTestDB(t)
	repo := message.NewPostgresRepository(pool)
	ctx := context.Background()

	convID := createTestConversation(t, pool)

	msg := &domain.Message{
		ConversationID: convID,
		Role:           domain.RoleAssistant,
		Content:        "The sky is blue because...",
		ContentTier:    domain.ContentTierNormal,
	}

	err := repo.Create(ctx, msg)
	require.NoError(t, err)

	// Verify content tier was saved
	msgs, err := repo.GetByConversationID(ctx, convID)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	assert.Equal(t, domain.ContentTierNormal, msgs[0].ContentTier)
}

func TestPostgresRepository_GetByConversationID(t *testing.T) {
	pool := setupTestDB(t)
	repo := message.NewPostgresRepository(pool)
	ctx := context.Background()

	convID := createTestConversation(t, pool)

	// Create multiple messages
	for i := 0; i < 5; i++ {
		role := domain.RoleUser
		if i%2 == 1 {
			role = domain.RoleAssistant
		}
		msg := &domain.Message{
			ConversationID: convID,
			Role:           role,
			Content:        "Test message",
		}
		err := repo.Create(ctx, msg)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Get all messages
	msgs, err := repo.GetByConversationID(ctx, convID)
	require.NoError(t, err)
	assert.Len(t, msgs, 5)

	// Verify order (ascending by created_at)
	for i := 1; i < len(msgs); i++ {
		assert.True(t, msgs[i].CreatedAt.After(msgs[i-1].CreatedAt) || msgs[i].CreatedAt.Equal(msgs[i-1].CreatedAt))
	}
}

func TestPostgresRepository_GetRecentByConversationID(t *testing.T) {
	pool := setupTestDB(t)
	repo := message.NewPostgresRepository(pool)
	ctx := context.Background()

	convID := createTestConversation(t, pool)

	// Create 10 messages
	for i := 0; i < 10; i++ {
		msg := &domain.Message{
			ConversationID: convID,
			Role:           domain.RoleUser,
			Content:        "Test message",
		}
		err := repo.Create(ctx, msg)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// Get last 5 messages
	msgs, err := repo.GetRecentByConversationID(ctx, convID, 5)
	require.NoError(t, err)
	assert.Len(t, msgs, 5)

	// Verify order (ascending by created_at)
	for i := 1; i < len(msgs); i++ {
		assert.True(t, msgs[i].CreatedAt.After(msgs[i-1].CreatedAt) || msgs[i].CreatedAt.Equal(msgs[i-1].CreatedAt))
	}
}

func TestPostgresRepository_GetByConversationID_Empty(t *testing.T) {
	pool := setupTestDB(t)
	repo := message.NewPostgresRepository(pool)
	ctx := context.Background()

	convID := createTestConversation(t, pool)

	msgs, err := repo.GetByConversationID(ctx, convID)
	require.NoError(t, err)
	assert.Empty(t, msgs)
}
