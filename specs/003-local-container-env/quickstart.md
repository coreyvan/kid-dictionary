# Quickstart: Local Containerized Development Environment

**Branch**: `003-local-container-env` | **Date**: 2025-12-05

## Prerequisites

- Docker Desktop, Podman, Colima, or compatible container runtime
- Git
- (Optional) direnv for automatic environment loading

## Getting Started

### 1. Clone and Configure

```bash
# Clone the repository
git clone https://github.com/coreyvan/kid-dictionary.git
cd kid-dictionary

# Copy environment template
cp .env.example .env

# Edit .env with your values (especially OPENAI_API_KEY)
# Or use direnv: echo "dotenv .env" > .envrc && direnv allow
```

### 2. Start the Environment

```bash
# Start all services (database, migrations, app with hot reload)
task up

# Or using docker compose directly:
docker compose up -d
```

### 3. Verify Everything is Running

```bash
# Check service status
task status
# Or: docker compose ps

# View logs
task logs
# Or: docker compose logs -f

# Test the API
curl http://localhost:8080/health
```

### 4. Development Workflow

```bash
# Make code changes - Air will automatically rebuild

# View service logs only
task logs:app
# Or: docker compose logs -f app

# Connect to database directly
task db:shell
# Or: docker compose exec db psql -U postgres -d kid_dictionary
```

### 5. Stop the Environment

```bash
# Stop all services (data preserved)
task down
# Or: docker compose down

# Stop and remove volumes (clean slate)
task reset
# Or: docker compose down -v
```

## Available Commands

| Command | Description |
|---------|-------------|
| `task up` | Start the complete development environment |
| `task down` | Stop all containers (preserves data) |
| `task reset` | Stop containers and wipe database |
| `task status` | Show container status |
| `task logs` | Follow all container logs |
| `task logs:app` | Follow app container logs only |
| `task db:shell` | Open PostgreSQL shell |

## Environment Variables

### Required (Sensitive)

| Variable | Description |
|----------|-------------|
| `DB_PASSWORD` | PostgreSQL password |
| `OPENAI_API_KEY` | OpenAI API key for LLM features |

### Optional (with defaults)

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `db` | Database hostname |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_NAME` | `kid_dictionary` | Database name |
| `LISTEN_PORT` | `8080` | Service port |
| `LOG_LEVEL` | `info` | Log level |
| `PRETTY_LOG` | `true` | Human-readable logs |

## Troubleshooting

### Port Already in Use

```bash
# Find what's using the port
lsof -i :8080

# Stop the offending process or change LISTEN_PORT in .env
```

### Migrations Failed

```bash
# Check migration logs
docker compose logs migrate

# Reset and retry
task reset && task up
```

### Hot Reload Not Working

```bash
# Check air logs
docker compose logs app | grep -i air

# Restart the app container
docker compose restart app
```

### Database Connection Issues

```bash
# Verify database is healthy
docker compose ps db

# Check database logs
docker compose logs db

# Connect directly to verify
docker compose exec db pg_isready -U postgres
```

## Architecture Overview

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│    App      │────▶│  Database   │
│             │     │   :8080     │     │   :5432     │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │  Air (hot   │
                    │   reload)   │
                    └─────────────┘
```

**Startup sequence**:
1. Database starts and becomes healthy
2. Migration container runs and applies schema
3. App container starts with hot reload enabled
