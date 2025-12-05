# Contracts: Layer Boundary Enforcement

**Feature**: 002-layer-boundaries

## No API Changes

This feature is an internal refactoring that enforces layer boundaries. There are no changes to:

- Proto definitions (`proto/kiddictionary/v1/`)
- Generated Connect handlers
- OpenAPI specifications
- External API behavior

The existing proto contracts remain unchanged. All changes are to internal Go package structure.

## Internal Contracts

The refactoring introduces internal contracts through repository interfaces in `internal/domain/`:

- `domain.ConversationRepository` - interface for conversation persistence
- `domain.MessageRepository` - interface for message persistence

These interfaces define the contract between the service and repository layers.
