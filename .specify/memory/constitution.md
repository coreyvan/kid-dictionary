<!--
Sync Impact Report
==================
Version change: 1.0.0 → 1.1.0 (MINOR - new guidance added)

Modified principles: None

Added sections:
- Development Workflow > Documentation Maintenance (new subsection)

Removed sections: None

Templates requiring updates:
- .specify/templates/plan-template.md ✅ (no changes needed - Constitution Check compatible)
- .specify/templates/spec-template.md ✅ (no changes needed)
- .specify/templates/tasks-template.md ✅ (Polish phase already includes documentation updates)

Follow-up TODOs: None
-->

# Kid Dictionary Constitution

## Core Principles

### I. Protobuf-First API Design

All API contracts MUST be defined in Protocol Buffer files before implementation begins.

- Proto files are the source of truth for service definitions
- Use Connect-Go to generate handlers and clients
- Separate request/response messages for each RPC (enables independent evolution)
- Use enums for bounded value sets (age brackets, message roles, content tiers)
- Version packages explicitly (e.g., `kiddictionary.v1`)
- Breaking changes require MAJOR version bump and migration path

**Rationale**: Protobuf-first design enforces explicit contracts between frontend and backend, enables type-safe client generation, and makes API evolution deliberate rather than accidental.

### II. Test-Alongside Development

Tests MUST accompany implementation for all business logic and service boundaries.

- Unit tests for business logic (services, domain operations)
- Integration tests for repository implementations against real PostgreSQL
- Contract tests for API boundaries using generated Connect clients
- Tests may be written during implementation, not required to fail first
- Mock LLM providers in tests to avoid cost and flakiness

**Rationale**: Tests validate correctness and prevent regressions. Allowing tests alongside (rather than strictly before) implementation balances quality assurance with development velocity.

### III. Simplicity & YAGNI

Every abstraction, pattern, and architectural decision MUST justify its existence.

- Start with the simplest solution that works
- Add complexity only when concrete problems demand it
- Avoid speculative features, premature optimization, and unnecessary indirection
- Three similar lines of code are better than a premature abstraction
- Delete unused code entirely; no backward-compatibility hacks for internal code

**Rationale**: Over-engineering obscures intent, increases maintenance burden, and slows iteration. Kid Dictionary's value is in clear explanations, not architectural sophistication.

### IV. Observability

All production code MUST be debuggable through structured logging and metrics.

- Use slog for structured logging with consistent fields (request_id, user_id, method, duration)
- Log at appropriate levels: error for failures, info for requests, debug for tracing
- Include request IDs for correlation across service boundaries
- Instrument key operations: request latency, LLM call duration, error rates
- Health endpoints MUST verify database and core dependencies

**Rationale**: When things go wrong in production, observability is the difference between quick diagnosis and prolonged outages. Structured logs and metrics make systems transparent.

### V. Content Safety

All generated content MUST be appropriate for the specified age bracket.

- Enforce age bracket classification (Little Ones 0-5, Growing Minds 5-10, Pre-Teens 10+)
- Apply content tier handling (Normal, Sensitive, Contextual, Redirect)
- Sensitive topics require soft guidance prefixes
- Harmful content requests MUST be politely declined with alternative suggestions
- System prompts are configuration, not hardcoded strings

**Rationale**: The core product promise is helping adults explain concepts to children safely. Content appropriateness is non-negotiable for trust and child safety.

## Additional Constraints

### Technology Stack

- **Language**: Go (current LTS version)
- **RPC Framework**: Connect-Go with Chi router
- **Database**: PostgreSQL with pgx driver
- **LLM Integration**: Provider-agnostic interface (OpenAI, Anthropic, etc.)
- **Configuration**: Environment variables following 12-factor principles
- **Containerization**: Multi-stage Docker builds with distroless/alpine runtime

### Security Requirements

- Never commit secrets to version control
- Use environment variables or secrets management for credentials
- Validate all input at service boundaries
- Map domain errors to appropriate Connect status codes
- JWT authentication with short-lived access tokens

## Development Workflow

### Code Organization

- `cmd/` - Application entrypoints
- `internal/` - Private application code (config, transport, domain)
- `pkg/` - Public libraries (env helpers, utilities)
- `proto/` - Protocol Buffer definitions
- `gen/` - Generated code from buf

### Documentation Maintenance

The README.md MUST be kept current with each feature implementation.

- **When to update**: Every feature that changes or adds items relevant to developers who clone the repository
- **What to include**:
  - New task commands added to Taskfile.yml
  - New environment variables or configuration options
  - Changes to development workflow (e.g., new prerequisites, setup steps)
  - New API endpoints or capabilities
  - Changes to architecture or package structure
- **Review checkpoint**: README updates SHOULD be included in the same PR as the feature implementation
- **Scope**: Focus on information needed to develop, run, and understand the codebase—not exhaustive feature documentation

**Rationale**: The README is the first point of contact for developers. Outdated documentation wastes time, causes confusion, and erodes trust in the codebase. Keeping it current as part of feature work prevents documentation debt.

### Quality Gates

- All PRs require passing tests before merge
- Buf lint and breaking change detection in CI
- golangci-lint for Go code style enforcement
- Generated code MUST be committed to avoid build-time dependencies on buf

### Review Standards

- Verify compliance with constitution principles
- Complexity MUST be justified against Principle III (Simplicity)
- Content-handling changes require extra scrutiny for Principle V (Safety)

## Governance

This constitution supersedes all other development practices. Amendments require:

1. Documented rationale for the change
2. Impact assessment on existing code and templates
3. Migration plan for affected areas
4. Version bump following semantic versioning:
   - MAJOR: Principle removal or incompatible redefinition
   - MINOR: New principle or material expansion
   - PATCH: Clarifications, typos, non-semantic refinements

All pull requests and code reviews MUST verify compliance with these principles. Non-compliance requires explicit justification documented in the PR.

Use CLAUDE.md for runtime development guidance specific to tooling and commands.

**Version**: 1.1.0 | **Ratified**: 2025-12-04 | **Last Amended**: 2025-12-06
