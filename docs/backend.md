# Kid Dictionary: Connect-Go Implementation Guide

A conceptual guide for building a Connect-based Go backend with Chi routing, PostgreSQL persistence, and LLM integration.

---

## What is Connect?

Connect is a family of libraries for building browser and gRPC-compatible RPC APIs. Unlike pure gRPC, Connect services work natively over HTTP/1.1 and HTTP/2 without a proxy, making them ideal for web frontends.

Key characteristics:

- **Protocol flexibility**: Supports Connect protocol (HTTP + JSON or Protobuf), gRPC, and gRPC-Web from a single service definition
- **Browser-native**: Works with fetch API directly, no gRPC-Web proxy needed
- **Standard HTTP handlers**: Generates `http.Handler` implementations that integrate with any Go router
- **Generated clients**: Produces idiomatic TypeScript and Go clients from proto definitions

For Kid Dictionary, this means your Vue frontend can call your Go backend using generated TypeScript clients with full type safety, and later you can add streaming without changing your transport layer.

---

## Protobuf-First API Design

### Philosophy

The proto files become your source of truth for the API contract. Design them as a product artifact, not an implementation detail. Think about how clients will consume the API, what data they need, and what operations make sense from a user's perspective.

### Message Design Principles

**Separate request/response messages**: Every RPC should have its own request and response message types, even if they seem redundant now. This allows independent evolution without breaking changes.

**Use wrapper messages for optional primitives**: Protobuf 3 doesn't distinguish between "not set" and "zero value" for primitives. For optional fields where the distinction matters, use message wrappers or the `optional` keyword.

**Prefer enums over strings for known value sets**: Age brackets, message roles, and similar bounded sets should be enums. This gives you compile-time checking and makes the valid values explicit in the schema.

**Include metadata fields for extensibility**: Timestamps, IDs, and version fields in responses make clients more resilient. Consider what information clients might need for caching, optimistic updates, or debugging.

### Service Organization

Group RPCs by resource and capability. For Kid Dictionary, you might organize around:

- **Authentication**: Registration, login, token refresh
- **Conversations**: CRUD operations for conversation sessions
- **Messages**: Sending messages and receiving AI responses
- **User preferences**: Managing user settings

Whether these are separate services or methods on a single service depends on how you want to structure your server implementation and whether different services might scale independently.

### Versioning Strategy

Include version in the package path (e.g., `kiddictionary.v1`). This allows you to run v1 and v2 simultaneously during migrations. Plan for the proto package structure to outlive any particular implementation.

---

## Buf Toolchain

### Purpose

Buf replaces the traditional protoc workflow with a more ergonomic toolchain. It handles linting, breaking change detection, and code generation with a consistent configuration approach.

### Configuration Concepts

**buf.yaml**: Defines your module, lint rules, and breaking change policy. The lint rules enforce proto style conventions, and breaking change detection prevents accidental API breakage.

**buf.gen.yaml**: Specifies which plugins generate what code and where. For Connect-Go, you'll generate Go types, Connect service interfaces, and optionally TypeScript clients.

### Generation Workflow

The typical flow is: edit protos → run buf generate → generated code appears in your output directory. The generated Go code includes:

- Message types with getters and proto marshal/unmarshal
- Connect handler interfaces your server implements
- Connect client constructors for calling other services

For TypeScript, you get equivalent generated clients that your Vue app imports directly.

### Dependency Management

Buf has a schema registry (BSR) for sharing proto dependencies. For standard types like timestamps and durations, you can import from `buf.build/googleapis/googleapis` or similar. This avoids vendoring common proto definitions.

---

## Connect + Chi Integration

### How It Fits Together

Connect generates `http.Handler` implementations from your service definitions. Chi is an `http.Handler`-based router. The integration is natural: mount Connect handlers on Chi paths.

Chi handles:
- Route organization and grouping
- Middleware chain (logging, recovery, CORS, etc.)
- Request context propagation
- Any non-RPC routes (health checks, metrics, static files)

Connect handles:
- Protocol negotiation (Connect, gRPC, gRPC-Web)
- Request deserialization and response serialization
- Service method dispatch
- Streaming lifecycle (when you add it later)

### Middleware Considerations

Chi middleware runs before Connect handlers, so authentication, logging, rate limiting, and similar concerns work naturally. The request context flows through, so you can attach user identity in middleware and read it in your service implementation.

Connect also has its own interceptor concept for cross-cutting concerns that are RPC-aware (like per-method authorization or request validation). Use Chi middleware for HTTP-level concerns and Connect interceptors for RPC-level concerns.

### Path Structure

Connect handlers mount at paths derived from the proto service name. Plan your URL structure knowing that API paths will look like `/package.Service/Method`. You can mount under a prefix like `/api` if you want to separate from other routes.

---

## Service Implementation Patterns

### Handler vs. Business Logic Separation

The Connect service implementation should be a thin translation layer. It receives proto request messages, calls into your business logic layer, and returns proto response messages. This keeps the generated proto types at the edges and lets your core domain use whatever internal representations make sense.

Benefits of this separation:
- Business logic remains testable without proto dependencies
- You can refactor internal types without touching the API
- Multiple transports (Connect, CLI, background jobs) can share business logic

### Dependency Injection

Your service implementation needs access to repositories, LLM providers, configuration, and other dependencies. Constructor injection keeps dependencies explicit and makes testing straightforward. The service struct holds interfaces, not concrete types.

### Error Handling

Connect uses a status code model similar to gRPC. Map your domain errors to appropriate Connect error codes:

- Not found → `NotFound`
- Validation failure → `InvalidArgument`
- Authentication required → `Unauthenticated`
- Permission denied → `PermissionDenied`
- Rate limited → `ResourceExhausted`
- Upstream failure (LLM down) → `Unavailable`

Connect errors can include details and metadata. Use this for structured error information clients can programmatically handle, like field-level validation errors.

### Request Validation

Validate requests at the service boundary before passing to business logic. Options include:

- **protovalidate**: Define validation rules in proto files using CEL expressions. Validation runs automatically via interceptor.
- **Manual validation**: Check fields in your service implementation and return `InvalidArgument` errors.

protovalidate is powerful but adds another concept to learn. For MVP, manual validation is fine.

---

## Authentication and Authorization

### Token-Based Authentication

JWT remains a solid choice for stateless authentication. The flow:

1. User registers or logs in via dedicated auth RPCs
2. Server returns access token (short-lived) and optionally refresh token (longer-lived)
3. Client includes access token in subsequent requests
4. Middleware validates token and attaches user identity to context

For Connect specifically, the token typically goes in the `Authorization` header, which works identically to REST.

### Middleware Placement

Authentication middleware runs at the Chi level before requests reach Connect handlers. This lets you:

- Allow certain paths (login, register, health) to skip auth
- Attach user identity to context for downstream use
- Return 401/403 without invoking the service layer

### Per-Method Authorization

Some RPCs need authorization beyond "is authenticated." For example, deleting a conversation requires ownership verification. This logic belongs in the service layer where you have access to both the user identity and the resource.

Connect interceptors can also enforce authorization rules, which is useful for coarse-grained checks like "this whole service requires admin role."

---

## Repository Layer

### Interface Design

Define repository interfaces in terms of your domain concepts, not database concepts. The interface should express what the application needs, not how the database stores it.

Repository methods should:
- Accept context for cancellation and timeout propagation
- Return domain types, not database row types
- Return domain errors (not found, conflict) not database errors (no rows, unique violation)
- Handle transactions internally or accept a transaction context

### PostgreSQL with pgx

pgx is the modern choice for Postgres in Go. Key concepts:

**Connection pooling**: Use `pgxpool` for concurrent access. Configure pool size based on expected load and Postgres connection limits.

**Query modes**: pgx supports both direct queries and prepared statements. For dynamic queries, use `pgx.NamedArgs` or `pgx.StrictNamedArgs` to avoid SQL injection while keeping queries readable.

**Scanning**: pgx can scan into structs using `pgx.RowToStructByName`. Define database DTOs that mirror your table structure, then map to domain types.

**Transactions**: Use `pool.BeginTx` for operations that need atomicity. Pass the transaction to repository methods or use a unit-of-work pattern.

### Migrations

Use a migration tool (golang-migrate, goose, or similar) to version your schema. Migrations should:

- Be idempotent where possible
- Have corresponding down migrations for development
- Live in version control alongside application code
- Run automatically on deployment or as a separate step

### Connection String Management

For containerized deployment, the connection string comes from environment variables. In production, you'll likely use IAM authentication or secrets management rather than static passwords.

---

## LLM Provider Layer

### Abstraction Design

Define an interface for LLM interactions that hides provider-specific details. The interface should express what your application needs: take a conversation context and age bracket, return a completion.

This abstraction enables:
- Swapping providers without changing business logic
- Implementing caching as a wrapper
- Adding fallback providers for resilience
- Mocking in tests

### Provider Implementation

Each LLM provider (OpenAI, Anthropic, etc.) gets its own implementation. The implementation handles:

- API authentication
- Request formatting (system prompts, message structure)
- Response parsing
- Error mapping (rate limits, overload, invalid requests)
- Timeout configuration

### System Prompt Management

System prompts vary by age bracket and potentially by topic sensitivity tier. Store these as configuration or embedded files rather than hardcoding in the provider. This makes iteration easier and keeps prompts version-controlled.

### Cost Control

Several strategies reduce LLM costs:

**Semantic caching**: Before calling the LLM, check if a similar question (same age bracket, high embedding similarity) was recently answered. Serve cached responses for near-duplicates.

**Context limiting**: Implement the hard conversation length limit discussed earlier. Truncate or summarize older messages to control input token count.

**Model selection**: Use cheaper models (GPT-4o-mini, Claude Haiku) for most requests. Reserve expensive models for complex queries if quality differs noticeably.

**Rate limiting**: Per-user limits prevent abuse and make costs predictable.

---

## Streaming (Future)

Connect supports server streaming, client streaming, and bidirectional streaming natively. For Kid Dictionary, server streaming (responses that arrive incrementally) improves the chat UX.

When you're ready to add streaming:

1. Define the RPC as `stream` in the proto response
2. Implement the handler to yield response chunks
3. Wrap the LLM provider to expose streaming
4. Update the frontend to consume the stream

The Connect TypeScript client handles streaming naturally with async iterators. The Chi integration works the same—streaming handlers are still `http.Handler` implementations.

---

## Testing Strategy

### Unit Testing

Test business logic in isolation by mocking repository and LLM provider interfaces. Focus on:

- Service method behavior for various inputs
- Error handling and edge cases
- Authorization logic

Use testify/suite for organized test setup and teardown.

### Integration Testing

Test repository implementations against a real Postgres instance. Use testcontainers-go or a docker-compose test environment to spin up Postgres automatically.

Test the full service implementation with real repositories but mocked LLM providers. This validates the wiring without incurring LLM costs.

### End-to-End Testing

Use the generated Connect client to call your running server. This validates proto compatibility and the full request/response cycle.

### Contract Testing

Buf's breaking change detection serves as a form of contract testing for the API. Run it in CI to catch unintentional API breakage.

---

## Configuration

### Environment-Based

Use environment variables for all configuration. This aligns with 12-factor app principles and works naturally with container orchestration.

Configuration categories:
- **Server**: Port, host, timeouts
- **Database**: Connection string, pool settings
- **LLM**: Provider API keys, model selection, timeout
- **Auth**: JWT signing key, token expiration
- **Feature flags**: Enable/disable streaming, caching, etc.

### Secrets Management

For local development, environment variables or a `.env` file work fine. For production, use your platform's secrets management:
- AWS Secrets Manager / SSM Parameter Store
- GCP Secret Manager
- Kubernetes Secrets
- HashiCorp Vault

Never commit secrets to version control.

---

## Containerization

### Dockerfile Patterns

Use multi-stage builds to keep images small:

1. **Builder stage**: Full Go toolchain, build the binary
2. **Runtime stage**: Minimal base (distroless or alpine), copy binary only

The final image should contain only the compiled binary and any runtime dependencies (CA certificates, timezone data).

### Health Checks

Expose a health endpoint (can be a regular Chi route, doesn't need to be an RPC) that verifies:
- Server is accepting connections
- Database connection is healthy
- LLM provider is reachable (optional, may not want to block on this)

Container orchestrators use this for readiness and liveness probes.

### Graceful Shutdown

Handle SIGTERM to drain in-flight requests before exiting. Go's `http.Server.Shutdown` does this, but you also need to close database pools and flush any buffered data.

---

## Horizontal Scaling Considerations

### Statelessness

The server should hold no request-scoped state. All state lives in:
- PostgreSQL (durable)
- Redis (ephemeral/cache)
- LLM provider (external)

This allows running multiple server instances behind a load balancer.

### Connection Pooling

Each server instance maintains its own connection pool to Postgres. Size pools appropriately so that max_instances × pool_size ≤ Postgres max_connections, with headroom for migrations and admin access.

### Caching Layer

When you add semantic caching, use Redis or a similar shared cache rather than in-memory caching. In-memory caches aren't shared across instances and complicate scaling.

### Session Affinity

Not required for this architecture. All instances can handle any request. Avoid designs that require sticky sessions.

---

## Deployment Topology

### MVP

Single instance is fine. Deploy to:
- Railway, Render, or Fly.io for simplicity
- A single VPS (DigitalOcean, Hetzner) for cost

Use managed Postgres (Supabase, Neon, or provider-included) to avoid ops burden.

### Growth

Add horizontal scaling when single-instance limits are reached:
- Load balancer in front of multiple instances
- Auto-scaling based on CPU or request latency
- Connection pooler (PgBouncer) if connection limits become an issue

### Production

- CDN in front for static assets and API caching
- Separate read replicas if read volume justifies it
- Multi-region if latency matters for your users
- Proper observability (metrics, traces, logs)

---

## Observability

### Logging

Use structured logging (slog is in stdlib now) with consistent fields:
- Request ID (propagate through context)
- User ID when authenticated
- RPC method name
- Duration and status

Log at appropriate levels: errors for failures, info for requests, debug for detailed tracing.

### Metrics

Instrument key operations:
- Request rate, latency, and error rate by method
- LLM call duration and token usage
- Database query latency
- Cache hit rate (when implemented)

Prometheus format is widely supported. Connect and Chi both have middleware for automatic request metrics.

### Tracing

Distributed tracing helps debug latency issues across service boundaries. OpenTelemetry is the standard. For MVP, logging with request IDs is sufficient; add tracing when complexity warrants.

---

## Development Workflow

### Local Environment

docker-compose for Postgres (and later Redis). Run the Go server directly with `go run` or use air/gow for hot reload.

Generate proto code before building. A Makefile target like `make generate` keeps this consistent.

### CI Pipeline

1. Lint protos with buf
2. Check for breaking changes with buf
3. Generate code
4. Run Go linter (golangci-lint)
5. Run tests
6. Build container image
7. Push to registry

### Deployment

Automate deployment from main branch or tags. Start simple (push triggers deploy) and add approval gates as needed.

---

## Summary

This architecture gives you:

- **Strong API contracts** via protobuf and Connect
- **Familiar Go patterns** with Chi routing and clean layering
- **Horizontal scalability** through stateless design
- **Cost control** via caching and model selection
- **Future streaming** built into the transport layer
- **Generated clients** for your Vue frontend

Start with the proto definitions—they're the foundation everything else builds on. Get a basic RPC working end-to-end, then iterate on features.