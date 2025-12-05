package domain

import "errors"

// Domain errors shared across layers.
var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidAgeBracket is returned when an age bracket is unspecified or out of range.
	ErrInvalidAgeBracket = errors.New("invalid age bracket")

	// ErrTitleRequired is returned when a title is empty when required.
	ErrTitleRequired = errors.New("title is required")

	// ErrConversationNotFound is returned when a conversation does not exist.
	ErrConversationNotFound = errors.New("conversation not found")

	// ErrEmptyContent is returned when message content is empty.
	ErrEmptyContent = errors.New("content is required")

	// ErrContentTooLong is returned when message content exceeds maximum length.
	ErrContentTooLong = errors.New("content exceeds maximum length")
)
