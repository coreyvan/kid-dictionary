# Kid Dictionary

A Go backend service that helps adults explain concepts to children using age-appropriate language. Uses Connect-Go (gRPC-compatible RPC) with PostgreSQL for persistence and OpenAI for generating explanations.

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)
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

# 3. Set up git hooks (recommended)
task init:hooks

# 4. Verify it's running
task status
curl http://localhost:8080/health
```

## Development Commands

### Docker Environment

| Command | Description |
|---------|-------------|
| `task up` | Start complete dev environment (database, migrations, app with hot reload) |
| `task down` | Stop all containers (preserves data) |
| `task reset` | Stop containers and wipe database (clean slate) |
| `task status` | Show container status |
| `task logs` | Follow all container logs |
| `task logs:app` | Follow app container logs only |
| `task db:shell` | Open PostgreSQL shell |

### Local Development

| Command | Description |
|---------|-------------|
| `task init` | Copy .env.example to .env |
| `task tidy` | Tidy Go module dependencies |
| `task build` | Build the server binary |
| `task run` | Run the server locally |

### Code Generation

| Command | Description |
|---------|-------------|
| `task generate` | Generate code from protobuf definitions |
| `task generate:mocks` | Generate mock implementations using moq |
| `task lint:proto` | Lint protobuf definitions |

### Testing

| Command | Description |
|---------|-------------|
| `task test:unit` | Run unit tests (excludes acceptance tests) |
| `task test:acceptance` | Run acceptance tests (requires database running) |

### Code Quality

| Command | Description |
|---------|-------------|
| `task lint` | Run golangci-lint on entire codebase |
| `task lint:staged` | Run golangci-lint on staged files only |
| `task layer-check` | Verify layer boundary compliance |
| `task init:hooks` | Install Lefthook and set up git hooks |

## Git Hooks

This project uses [Lefthook](https://github.com/evilmartians/lefthook) for git hooks. Install with:

```bash
task init:hooks
```

### Pre-commit Hook
- Runs golangci-lint on staged Go files
- Fast feedback (~3 seconds)

### Pre-push Hook
Runs full validation before pushing:
1. Regenerates protobuf code and mocks
2. Checks for uncommitted generated files
3. Runs `go mod tidy` and checks for changes
4. Runs unit tests
5. Runs acceptance tests
6. Builds the binary

To bypass hooks in emergencies: `git commit --no-verify` or `git push --no-verify`

## Environment Variables

### Required

| Variable | Description |
|----------|-------------|
| `OPENAI_API_KEY` | Your OpenAI API key |

### Optional (with defaults)

| Variable | Default | Description |
|----------|---------|-------------|
| `BIND_ADDR` | `0.0.0.0` | Server bind address |
| `LISTEN_PORT` | `8080` | Server port |
| `PRETTY_LOG` | `false` | Human-readable log output |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `OPENAI_MODEL` | `gpt-4o-mini` | OpenAI model to use |

### Database Connection

Either set `DATABASE_URL` or individual components:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | - | Full PostgreSQL connection URL |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `kid_dictionary` | Database name |

## Architecture

The service uses a three-layer architecture with strict dependency rules:

```
Transport → Service → Repository
    ↓          ↓          ↓
          Domain (shared by all)
```

| Layer | Can Import | Cannot Import |
|-------|------------|---------------|
| Transport | domain, service | repository |
| Service | domain only | transport, repository |
| Repository | domain | transport, service |

Run `task layer-check` to verify compliance.

### Package Structure

- `cmd/server/` - Application entrypoint
- `internal/transport/` - Connect-Go HTTP handlers
- `internal/service/` - Business logic (conversation, message)
- `internal/repository/` - PostgreSQL persistence
- `internal/domain/` - Shared models and interfaces
- `internal/config/` - Configuration and dependency wiring
- `internal/llm/` - LLM provider integration
- `gen/` - Generated protobuf code
- `tests/acceptance/` - End-to-end acceptance tests

## API

The service exposes Connect-Go (gRPC-compatible) endpoints:

### ConversationService
- `CreateConversation` - Create a new conversation with age bracket
- `GetConversation` - Get conversation with all messages
- `ListConversations` - Paginated list of conversations
- `UpdateConversation` - Update title or age bracket
- `DeleteConversation` - Delete conversation and messages

### MessageService
- `SendMessage` - Send a message and receive AI-generated response

### Age Brackets
- **Little Ones (0-5)** - Simple, concrete words, 2-3 sentences
- **Growing Minds (5-10)** - Simple metaphors, 3-5 sentences
- **Pre-Teens (10+)** - Nuance, multiple perspectives

### Health Check
```bash
curl http://localhost:8080/health
```

## Running Without Docker

```bash
# 1. Start PostgreSQL locally (or use existing instance)
# 2. Configure environment
task init
# Edit .env with your database credentials

# 3. Run migrations
# (migrations are in migrations/ directory)

# 4. Start the server
task run
```
