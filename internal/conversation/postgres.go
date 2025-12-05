package conversation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a conversation is not found.
var ErrNotFound = errors.New("conversation not found")

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL-backed repository.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create stores a new conversation and populates its ID.
func (r *PostgresRepository) Create(ctx context.Context, conv *Conversation) error {
	if conv.ID == uuid.Nil {
		conv.ID = uuid.New()
	}
	now := time.Now()
	conv.CreatedAt = now
	conv.UpdatedAt = now

	query := `
		INSERT INTO conversations (id, user_id, title, age_bracket, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		conv.ID,
		conv.UserID,
		conv.Title,
		conv.AgeBracket,
		conv.CreatedAt,
		conv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert conversation: %w", err)
	}
	return nil
}

// GetByID retrieves a conversation by its ID.
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Conversation, error) {
	query := `
		SELECT id, user_id, title, age_bracket, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	conv := &Conversation{}
	err := row.Scan(
		&conv.ID,
		&conv.UserID,
		&conv.Title,
		&conv.AgeBracket,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query conversation: %w", err)
	}
	return conv, nil
}

// List returns conversations, optionally filtered by user ID.
func (r *PostgresRepository) List(ctx context.Context, userID *uuid.UUID, limit, offset int) ([]*Conversation, error) {
	var query string
	var args []interface{}

	if userID != nil {
		query = `
			SELECT id, user_id, title, age_bracket, created_at, updated_at
			FROM conversations
			WHERE user_id = $1
			ORDER BY updated_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{*userID, limit, offset}
	} else {
		query = `
			SELECT id, user_id, title, age_bracket, created_at, updated_at
			FROM conversations
			ORDER BY updated_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		conv := &Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.Title,
			&conv.AgeBracket,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversations: %w", err)
	}
	return conversations, nil
}

// Update modifies an existing conversation.
func (r *PostgresRepository) Update(ctx context.Context, conv *Conversation) error {
	conv.UpdatedAt = time.Now()

	query := `
		UPDATE conversations
		SET title = $2, age_bracket = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query,
		conv.ID,
		conv.Title,
		conv.AgeBracket,
		conv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update conversation: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a conversation and all its messages (via CASCADE).
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM conversations WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete conversation: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
