package connect

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/coreyvan/kid-dictionary/internal/conversation"
	"github.com/coreyvan/kid-dictionary/internal/llm"
	"github.com/coreyvan/kid-dictionary/internal/message"
)

// User-friendly error messages
const (
	msgConversationNotFound = "The conversation you're looking for doesn't exist. Please start a new conversation."
	msgInvalidAgeBracket    = "Please select a valid age bracket: Little Ones (0-5), Growing Minds (5-10), or Pre-Teens (10+)."
	msgContentRequired      = "Please enter a question or concept you'd like help explaining."
	msgContentTooLong       = "Your message is too long. Please keep it under 500 characters."
	msgServiceUnavailable   = "We're having trouble connecting to our AI service. Please try again in a moment."
	msgTimeout              = "Your request took too long. Please try again with a simpler question."
	msgRateLimited          = "You've sent too many requests. Please wait a moment before trying again."
	msgUnexpected           = "Something unexpected happened. Please try again."
)

// MapError converts domain errors to Connect errors with appropriate status codes and user-friendly messages.
func MapError(err error) *connect.Error {
	if err == nil {
		return nil
	}

	// Context errors (timeout)
	if errors.Is(err, context.DeadlineExceeded) {
		return connect.NewError(connect.CodeDeadlineExceeded, errors.New(msgTimeout))
	}

	// Conversation errors
	if errors.Is(err, conversation.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, errors.New(msgConversationNotFound))
	}
	if errors.Is(err, conversation.ErrInvalidAgeBracket) {
		return connect.NewError(connect.CodeInvalidArgument, errors.New(msgInvalidAgeBracket))
	}
	if errors.Is(err, conversation.ErrTitleRequired) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Message errors
	if errors.Is(err, message.ErrContentRequired) {
		return connect.NewError(connect.CodeInvalidArgument, errors.New(msgContentRequired))
	}
	if errors.Is(err, message.ErrContentTooLong) {
		return connect.NewError(connect.CodeInvalidArgument, errors.New(msgContentTooLong))
	}
	if errors.Is(err, message.ErrConversationNotFound) {
		return connect.NewError(connect.CodeNotFound, errors.New(msgConversationNotFound))
	}

	// LLM errors
	if errors.Is(err, llm.ErrInvalidAPIKey) {
		return connect.NewError(connect.CodeInternal, errors.New(msgServiceUnavailable))
	}
	if errors.Is(err, llm.ErrRateLimited) {
		return connect.NewError(connect.CodeResourceExhausted, errors.New(msgRateLimited))
	}
	if errors.Is(err, llm.ErrProviderUnavailable) {
		return connect.NewError(connect.CodeUnavailable, errors.New(msgServiceUnavailable))
	}
	if errors.Is(err, llm.ErrTimeout) {
		return connect.NewError(connect.CodeDeadlineExceeded, errors.New(msgTimeout))
	}
	if errors.Is(err, llm.ErrUserRateLimited) {
		return connect.NewError(connect.CodeResourceExhausted, errors.New(msgRateLimited))
	}

	// Default to internal error
	return connect.NewError(connect.CodeInternal, errors.New(msgUnexpected))
}
