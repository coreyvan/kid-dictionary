# Research: Local Containerized Development Environment

**Branch**: `003-local-container-env` | **Date**: 2025-12-05

## Migration Tool Selection

**Decision**: golang-migrate

**Rationale**:
- Official Docker image (`migrate/migrate`) optimized for one-shot containers
- Most widely adopted in the Go ecosystem with extensive database driver support
- Simple, focused on SQL migrations only - aligns with Principle III (Simplicity)
- Well-documented Docker usage patterns
- Natural one-shot behavior: exits with code 0 after successful completion

**Alternatives Considered**:
- **goose**: Better for embedding migrations in application, but less mature Docker image support
- **sql-migrate**: Less commonly used in modern Go projects, fewer community resources

**Migration File Compatibility**: The existing migration files use `-- +migrate Up/Down` syntax which is sql-migrate format. golang-migrate uses a different naming convention (`000001_name.up.sql`/`000001_name.down.sql`).

**Action Required**: Either:
1. Convert existing migrations to golang-migrate format (recommended for long-term)
2. Use a wrapper script that runs sql-migrate in a container

**Updated Decision**: Use **sql-migrate** since migrations already use that format. Run via `rubenv/sql-migrate` Docker image or build a simple migration container.

## Database Readiness Pattern

**Decision**: PostgreSQL healthcheck with `pg_isready` + `start_period`

**Rationale**:
- Native PostgreSQL tool, no additional dependencies
- `start_period` gives PostgreSQL initialization time before failures count
- Combined with query validation ensures actual connectivity

**Configuration**:
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
  interval: 5s
  timeout: 5s
  retries: 5
  start_period: 30s
```

**Alternatives Considered**:
- **wait-for-it.sh**: External script, adds complexity
- **Custom connection loop**: More code to maintain
- **Direct query check**: Overkill for development environment

## Service Dependency Chain

**Decision**: Use Docker Compose `depends_on` with conditions

**Rationale**:
- `service_healthy` for database: ensures PostgreSQL is ready
- `service_completed_successfully` for migrations: ensures schema is applied before app starts
- Native Docker Compose feature, no external tools needed

**Dependency Chain**:
```
db (healthy) → migrate (completed successfully) → app
```

**Configuration**:
```yaml
app:
  depends_on:
    db:
      condition: service_healthy
    migrate:
      condition: service_completed_successfully
```

**Alternatives Considered**:
- **Application-level retry**: Adds complexity to app code
- **docker-compose-wait**: External tool, adds dependencies

## Hot Reload with Air

**Decision**: Multi-stage Dockerfile with Air for development

**Rationale**:
- `air-verse/air` is actively maintained and well-documented
- Multi-stage build keeps production image separate (Principle III)
- Volume mounts for source enable hot reload without rebuild

**Configuration**:
- Bind mount source code: `./:/app`
- Named volumes for caches: `go_modules:/go/pkg/mod`, `go_cache:/root/.cache/go-build`
- Polling enabled for cross-platform compatibility

**Performance Optimizations**:
- Named volumes for Go modules and build cache (faster than bind mounts)
- Exclude tmp/vendor/node_modules from watch
- 1000ms delay before rebuild to batch rapid changes

**Alternatives Considered**:
- **CompileDaemon**: Less actively maintained
- **realize**: Archived project
- **nodemon + go build**: Extra dependency, not Go-native

## Environment Variable Strategy

**Decision**: Combine `env_file` with explicit `environment` interpolation

**Rationale**:
- `.env` file for base configuration (can be committed via `.env.example`)
- Host environment passthrough for sensitive values (never committed)
- Explicit variable listing makes dependencies clear
- Compatible with direnv (`dotenv .env`)

**Configuration Pattern**:
```yaml
env_file:
  - .env
environment:
  # Override DB_HOST for container networking
  - DB_HOST=db
  - DB_PORT=5432
  # Pass through from host (sensitive)
  - DB_PASSWORD=${DB_PASSWORD}
  - OPENAI_API_KEY=${OPENAI_API_KEY}
```

**direnv Integration**:
```bash
# .envrc
dotenv .env
# Additional exports if needed
```

**Alternatives Considered**:
- **Pure env_file**: Doesn't allow container-specific overrides
- **Pure environment block**: Verbose, harder to maintain
- **Docker secrets**: Overkill for local development

## Container Architecture Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Network                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐    ┌───────────┐    ┌────────────────────┐   │
│  │    db    │───▶│  migrate  │───▶│        app         │   │
│  │ postgres │    │ sql-migrate│    │ go + air (hot)     │   │
│  └──────────┘    └───────────┘    └────────────────────┘   │
│       │               │                    │                │
│       │               │                    │                │
│  healthcheck     exits 0/1           hot reload             │
│  pg_isready      on complete         volume mount           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
        │                                     │
        ▼                                     ▼
   postgres_data                         ./:/app (source)
   (named volume)                        go_modules (cache)
```

## Files to Create

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Main orchestration with db, migrate, app services |
| `Dockerfile` | Multi-stage: dev (air) and prod (optimized binary) |
| `.air.toml` | Hot reload configuration |
| `.env.example` | Template with all required variables |
| `Taskfile.yml` (update) | Add `up`, `down`, `reset` commands |

## Environment Variables Required

| Variable | Sensitive | Default | Description |
|----------|-----------|---------|-------------|
| `DB_HOST` | No | `db` | Database hostname (container name) |
| `DB_PORT` | No | `5432` | Database port |
| `DB_USER` | No | `postgres` | Database user |
| `DB_PASSWORD` | **Yes** | - | Database password |
| `DB_NAME` | No | `kid_dictionary` | Database name |
| `BIND_ADDR` | No | `0.0.0.0` | Service bind address |
| `LISTEN_PORT` | No | `8080` | Service listen port |
| `PRETTY_LOG` | No | `true` | Human-readable logs |
| `LOG_LEVEL` | No | `info` | Log level |
| `OPENAI_API_KEY` | **Yes** | - | OpenAI API key |
| `OPENAI_MODEL` | No | `gpt-4o-mini` | OpenAI model |
