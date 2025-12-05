# Kid Dictionary

A Go backend service that helps adults explain concepts to children using age-appropriate language.

## Prerequisites

- [Docker Desktop](https://docs.docker.com/desktop/), [Podman](https://podman.io/), or compatible container runtime
- [Taskfile](https://taskfile.dev/installation/) (optional, for convenience commands)

## Quick Start

```bash
# 1. Clone and configure
git clone https://github.com/coreyvan/kid-dictionary.git
cd kid-dictionary
cp .env.example .env
# Edit .env with your OPENAI_API_KEY

# 2. Start the environment
task up
# Or: docker compose up -d

# 3. Verify it's running
task status
curl http://localhost:8080/health
```

## Development Commands

| Command | Description |
|---------|-------------|
| `task up` | Start the complete development environment |
| `task down` | Stop all containers (preserves data) |
| `task reset` | Stop containers and wipe database (clean slate) |
| `task status` | Show container status |
| `task logs` | Follow all container logs |
| `task logs:app` | Follow app container logs only |
| `task db:shell` | Open PostgreSQL shell |

## Environment Variables

Required (set in `.env`):
- `OPENAI_API_KEY` - Your OpenAI API key
- `DB_PASSWORD` - Database password (default: `postgres`)

Optional (have sensible defaults):
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_NAME` - Database connection
- `LISTEN_PORT` - Service port (default: `8080`)
- `LOG_LEVEL` - Log level (default: `info`)

## Architecture

The service uses a three-layer architecture:
- **Transport** - Connect-Go RPC handlers
- **Service** - Business logic
- **Repository** - PostgreSQL persistence

Hot reload is enabled via [air](https://github.com/air-verse/air) - code changes are automatically rebuilt.

## Local Development (without Docker)

```bash
task init    # Copy .env.example to .env
task tidy    # Install Go dependencies
task run     # Run the server locally
```