package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a message is not found.
var ErrNotFound = errors.New("message not found")

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL-backed repository.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create stores a new message and populates its ID.
func (r *PostgresRepository) Create(ctx context.Context, msg *Message) error {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	msg.CreatedAt = time.Now()

	// Convert ContentTier to nullable integer
	var contentTier *int
	if msg.ContentTier != ContentTierUnspecified {
		tier := int(msg.ContentTier)
		contentTier = &tier
	}

	query := `
		INSERT INTO messages (id, conversation_id, role, content, content_tier, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		msg.ID,
		msg.ConversationID,
		msg.Role,
		msg.Content,
		contentTier,
		msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

// GetByConversationID retrieves all messages for a conversation.
func (r *PostgresRepository) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, role, content, content_tier, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer rows.Close()

	return scanMessages(rows)
}

// GetRecentByConversationID retrieves the last N messages for context.
func (r *PostgresRepository) GetRecentByConversationID(ctx context.Context, conversationID uuid.UUID, limit int) ([]*Message, error) {
	// Use subquery to get last N messages, then order ascending
	query := `
		SELECT id, conversation_id, role, content, content_tier, created_at
		FROM (
			SELECT id, conversation_id, role, content, content_tier, created_at
			FROM messages
			WHERE conversation_id = $1
			ORDER BY created_at DESC
			LIMIT $2
		) sub
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent messages: %w", err)
	}
	defer rows.Close()

	return scanMessages(rows)
}

// scanMessages is a helper to scan message rows.
func scanMessages(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}) ([]*Message, error) {
	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		var contentTier *int
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.Role,
			&msg.Content,
			&contentTier,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		if contentTier != nil {
			msg.ContentTier = ContentTier(*contentTier)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}
	return messages, nil
}
