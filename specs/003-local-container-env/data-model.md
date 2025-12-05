# Data Model: Local Containerized Development Environment

**Branch**: `003-local-container-env` | **Date**: 2025-12-05

## Overview

This feature is infrastructure-focused and does not introduce new domain entities. The data model section documents the container entities and their relationships for orchestration purposes.

## Container Entities

### Database Container (db)

| Attribute | Type | Description |
|-----------|------|-------------|
| image | string | `postgres:16-alpine` |
| volume | named | `postgres_data` for persistent storage |
| ports | exposed | `${DB_PORT}:5432` |
| healthcheck | pg_isready | Validates database is accepting connections |

**State Transitions**:
- Starting → Healthy (after pg_isready succeeds)
- Healthy → Stopped (on `docker compose down`)

### Migration Container (migrate)

| Attribute | Type | Description |
|-----------|------|-------------|
| image | string | Custom or `rubenv/sql-migrate` |
| volume | bind | `./migrations:/migrations` |
| depends_on | service | db (condition: service_healthy) |

**State Transitions**:
- Waiting → Running (when db is healthy)
- Running → Completed (exit 0) or Failed (exit 1)

### Service Container (app)

| Attribute | Type | Description |
|-----------|------|-------------|
| build | Dockerfile | Multi-stage, target: dev |
| volumes | bind + named | Source code + Go caches |
| ports | exposed | `${LISTEN_PORT}:8080` |
| depends_on | services | db (healthy), migrate (completed) |

**State Transitions**:
- Waiting → Building → Running (with hot reload)
- Running → Rebuilding → Running (on source change)
- Running → Stopped (on `docker compose down`)

## Volume Entities

| Volume | Type | Purpose | Persistence |
|--------|------|---------|-------------|
| `postgres_data` | named | Database files | Across restarts |
| `go_modules` | named | Go module cache | Across restarts |
| `go_cache` | named | Go build cache | Across restarts |
| `./:/app` | bind | Source code | Host filesystem |
| `./migrations:/migrations` | bind | Migration SQL files | Host filesystem |

## Network Entity

| Network | Driver | Purpose |
|---------|--------|---------|
| `app-network` | bridge | Inter-container communication |

**Service DNS**:
- `db` → Database container
- `app` → Service container

## Environment Variable Entity

See [research.md](./research.md#environment-variables-required) for complete variable specification.

## Relationship Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                       Host Machine                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  .env                    ./                ./migrations/     │
│    │                      │                     │            │
│    │ (passthrough)        │ (bind mount)        │ (bind)     │
│    ▼                      ▼                     ▼            │
│  ┌─────────────────────────────────────────────────────────┐│
│  │                   Docker Network                         ││
│  │                                                          ││
│  │  ┌────────────┐         ┌────────────┐                  ││
│  │  │     db     │◀───────▶│  postgres_ │                  ││
│  │  │  :5432     │         │    data    │                  ││
│  │  └─────┬──────┘         └────────────┘                  ││
│  │        │ healthy                                         ││
│  │        ▼                                                 ││
│  │  ┌────────────┐                                         ││
│  │  │  migrate   │                                         ││
│  │  │ (one-shot) │                                         ││
│  │  └─────┬──────┘                                         ││
│  │        │ completed                                       ││
│  │        ▼                                                 ││
│  │  ┌────────────┐         ┌────────────┐                  ││
│  │  │    app     │◀───────▶│ go_modules │                  ││
│  │  │  :8080     │         │ go_cache   │                  ││
│  │  └────────────┘         └────────────┘                  ││
│  │                                                          ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```
