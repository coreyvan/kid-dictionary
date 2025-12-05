# Research: Layer Boundary Enforcement

**Feature**: 002-layer-boundaries
**Date**: 2025-12-05

## Research Questions

### 1. Go Package Organization for Layered Architecture

**Question**: What is the best practice for organizing Go packages to enforce layer boundaries without circular imports?

**Decision**: Use a `domain` package for shared types with separate `service` and `repository` package hierarchies.

**Rationale**:
- Go's package system makes it difficult to have service and repository in the same package while keeping them separate
- A central `domain` package containing models and interfaces allows:
  - Service packages to define business logic using domain types
  - Repository packages to implement domain interfaces
  - Transport to use domain types and service packages
- This follows the "Dependency Rule" from Clean Architecture where dependencies point inward

**Alternatives Considered**:
1. **Keep types in service packages, repositories import service**
   - Rejected: Creates awkward coupling where repository "knows about" service
2. **Interface definitions in separate `ports` package**
   - Rejected: Over-engineering for this project size; adds unnecessary indirection
3. **Single `internal` package with file-based separation**
   - Rejected: No compile-time enforcement; relies on convention only

### 2. Transport-to-Service Communication Pattern

**Question**: Should transport layer call service methods directly or use an intermediary?

**Decision**: Direct method calls on service structs. Transport injects service dependencies.

**Rationale**:
- Go idiom favors simple, direct calls over complex patterns
- Services are already designed with clean interfaces
- No need for message passing, command pattern, or mediators at this scale
- Constitution Principle III (Simplicity) dictates minimal indirection

**Alternatives Considered**:
1. **Command/Query pattern (CQRS)**
   - Rejected: Over-engineering; no separate read/write models needed
2. **Message bus / Event-driven**
   - Rejected: Adds complexity without benefit for synchronous HTTP

### 3. Repository Interface Location

**Question**: Where should repository interfaces be defined given the constraint that service layer cannot import repository?

**Decision**: Repository interfaces are defined in the `domain` package alongside the domain models.

**Rationale**:
- Keeps related types together (Conversation model + ConversationRepository interface)
- Service layer imports domain to get both models and interfaces
- Repository layer imports domain to implement interfaces
- No circular dependencies possible

**Alternatives Considered**:
1. **Interfaces in service packages**
   - Rejected: Would require repository to import service (wrong direction)
2. **Interfaces in repository packages**
   - Rejected: Would require service to import repository (violates spec requirement)

### 4. Static Analysis for Enforcement

**Question**: What tools can enforce layer boundary violations in Go?

**Decision**: Use `go-arch-lint` or custom script with `go list -json` for CI enforcement.

**Rationale**:
- Go's type system doesn't natively prevent cross-layer imports
- Static analysis can flag violations before code review
- Can be integrated into CI pipeline for SC-004 (violations detectable before merge)

**Alternatives Considered**:
1. **Manual code review only**
   - Rejected: Error-prone, doesn't meet SC-003 (<10 seconds verification)
2. **Build tags to separate layers**
   - Rejected: Cumbersome; doesn't prevent mistakes during development

### 5. Type Translation Strategy

**Question**: How should transport and repository layers translate to/from domain types?

**Decision**: Explicit converter functions in each layer, no shared converter package.

**Rationale**:
- Transport layer owns `ToProto` / `FromProto` functions for API types
- Repository layer owns `ToEntity` / `FromEntity` functions for DB rows
- Each layer is responsible for its own translation
- Keeps domain types "clean" without serialization/persistence concerns

**Alternatives Considered**:
1. **Shared converter package**
   - Rejected: Creates a dependency that all layers share; muddles boundaries
2. **Methods on domain types (`conv.ToProto()`)**
   - Rejected: Pollutes domain types with layer-specific concerns
3. **Goverter / automatic generation**
   - Rejected: Adds tooling complexity; manual converters are explicit and testable

## Key Findings

1. **Package structure matters**: Go's import system is the enforcement mechanism. Proper package organization prevents violations at compile time.

2. **Domain package is central**: A dedicated domain package containing models and interfaces is the cleanest solution for Go projects.

3. **Simplicity over patterns**: Advanced patterns (CQRS, mediator, event sourcing) are unnecessary for this project size. Direct calls with clean interfaces suffice.

4. **Static analysis fills the gap**: Since Go can't enforce "import X but not Y from X", external tooling is needed for complete enforcement.

## Dependencies Identified

None new. The refactoring uses existing Go standard library and project dependencies.

## Risks Identified

1. **Migration complexity**: Moving files across packages requires updating all import statements
   - Mitigation: Systematic approach, file-by-file with tests running between changes

2. **Test refactoring**: Tests may need significant updates for new import paths
   - Mitigation: Update tests alongside source files

3. **Wireup changes**: `config/wireup.go` will need updates for new package paths
   - Mitigation: Update wireup last, after all packages are in place
