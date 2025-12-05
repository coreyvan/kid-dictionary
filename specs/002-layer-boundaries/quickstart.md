# Quickstart: Layer Boundary Enforcement

**Feature**: 002-layer-boundaries
**Date**: 2025-12-05

## Overview

This guide explains the layered architecture pattern and how to work with it after the refactoring is complete.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Transport Layer                          │
│  internal/transport/                                        │
│  - HTTP/Connect handlers                                    │
│  - Proto type conversion                                    │
│  - Request validation                                       │
│                          │                                  │
│                          ▼                                  │
├─────────────────────────────────────────────────────────────┤
│                     Service Layer                           │
│  internal/service/{conversation,message}/                   │
│  - Business logic                                           │
│  - Orchestration                                            │
│  - Domain validation                                        │
│                          │                                  │
│                          ▼                                  │
├─────────────────────────────────────────────────────────────┤
│                    Repository Layer                         │
│  internal/repository/{conversation,message}/                │
│  - PostgreSQL persistence                                   │
│  - Database row conversion                                  │
│  - SQL queries                                              │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Package                          │
│  internal/domain/                                           │
│  - Canonical domain models                                  │
│  - Repository interfaces                                    │
│  - Domain errors                                            │
│  (Imported by all layers)                                   │
└─────────────────────────────────────────────────────────────┘
```

## Import Rules

| Layer | Can Import | Cannot Import |
|-------|------------|---------------|
| Transport | `domain`, `service/*` | `repository/*` |
| Service | `domain` only | `transport`, `repository/*` |
| Repository | `domain` | `transport`, `service/*` |
| Cross-cutting (`config`, `llm`, `connect`) | Any | (exempt from layer rules) |

## Adding a New Feature

### 1. Define Domain Types

Add models and interfaces to `internal/domain/`:

```go
// internal/domain/newentity.go
package domain

type NewEntity struct {
    ID        uuid.UUID
    Name      string
    CreatedAt time.Time
}

type NewEntityRepository interface {
    Create(ctx context.Context, e *NewEntity) error
    GetByID(ctx context.Context, id uuid.UUID) (*NewEntity, error)
}
```

### 2. Implement Service Layer

Create business logic in `internal/service/newentity/`:

```go
// internal/service/newentity/service.go
package newentity

import (
    "context"
    "github.com/coreyvan/kid-dictionary/internal/domain"
)

type Service struct {
    repo domain.NewEntityRepository
}

func New(repo domain.NewEntityRepository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string) (*domain.NewEntity, error) {
    // Business logic here
    entity := &domain.NewEntity{Name: name}
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, err
    }
    return entity, nil
}
```

### 3. Implement Repository Layer

Create persistence in `internal/repository/newentity/`:

```go
// internal/repository/newentity/postgres.go
package newentity

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/coreyvan/kid-dictionary/internal/domain"
)

type PostgresRepository struct {
    db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *PostgresRepository {
    return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, e *domain.NewEntity) error {
    // Convert domain to row, execute SQL
    return nil
}
```

### 4. Add Transport Handler

Update `internal/transport/transport.go`:

```go
import (
    newentitysvc "github.com/coreyvan/kid-dictionary/internal/service/newentity"
)

// In server struct
newEntitySvc *newentitysvc.Service

// In handler method
func (s *server) CreateNewEntity(ctx context.Context, req *connect.Request[v1.CreateNewEntityRequest]) (*connect.Response[v1.CreateNewEntityResponse], error) {
    entity, err := s.newEntitySvc.Create(ctx, req.Msg.Name)
    if err != nil {
        return nil, intconnect.MapError(err)
    }

    // Convert domain to proto
    return connect.NewResponse(&v1.CreateNewEntityResponse{
        Entity: &v1.NewEntity{
            Id:   entity.ID.String(),
            Name: entity.Name,
        },
    }), nil
}
```

### 5. Wire Dependencies

Update `internal/config/wireup.go`:

```go
import (
    newentityrepo "github.com/coreyvan/kid-dictionary/internal/repository/newentity"
    newentitysvc "github.com/coreyvan/kid-dictionary/internal/service/newentity"
)

func (w *Wireup) ProvideNewEntityService() *newentitysvc.Service {
    repo := newentityrepo.NewPostgres(w.MustProvideDB())
    return newentitysvc.New(repo)
}
```

## Verifying Layer Compliance

Run the layer check before committing:

```bash
# Quick check (under 10 seconds per SC-003)
task layer-check

# Or manually
go build ./... && echo "Layer boundaries OK"
```

The build will fail if there are import violations because Go's package system prevents circular imports, and the structure ensures proper dependency direction.

## Common Mistakes

### ❌ Importing repository from transport

```go
// internal/transport/transport.go
import "github.com/coreyvan/kid-dictionary/internal/repository/message" // WRONG!
```

### ✅ Import service instead

```go
// internal/transport/transport.go
import "github.com/coreyvan/kid-dictionary/internal/service/message" // Correct
```

### ❌ Importing transport from service

```go
// internal/service/message/service.go
import "github.com/coreyvan/kid-dictionary/internal/transport" // WRONG!
```

### ✅ Service should only import domain

```go
// internal/service/message/service.go
import "github.com/coreyvan/kid-dictionary/internal/domain" // Correct
```

## Testing

Each layer can be tested independently:

- **Service tests**: Mock repository interface from domain package
- **Repository tests**: Use testcontainers for real PostgreSQL
- **Transport tests**: Mock service, test request/response mapping

```go
// internal/service/message/service_test.go
type mockRepo struct {
    messages []*domain.Message
}

func (m *mockRepo) Create(ctx context.Context, msg *domain.Message) error {
    m.messages = append(m.messages, msg)
    return nil
}

func TestService_Create(t *testing.T) {
    repo := &mockRepo{}
    svc := New(repo)
    // Test business logic...
}
```
