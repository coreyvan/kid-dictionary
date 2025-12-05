# Feature Specification: Layer Boundary Enforcement

**Feature Branch**: `002-layer-boundaries`
**Created**: 2025-12-05
**Status**: Draft
**Input**: User description: "Create a refactoring spec to make sure we have proper layer boundaries. The transport should only access the service layer. The service layer should be the only thing that uses the repository layer."

## Clarifications

### Session 2025-12-05

- Q: Where should the canonical domain models (the "source of truth" types) live? → A: Service layer owns domain models (transport & repository translate)
- User guidance: Data models at each layer must be distinct; service layer must not import any types from repository or transport layers; outer layers handle translation

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Developer Adds New Transport Handler (Priority: P1)

A developer needs to add a new API endpoint. When implementing the transport handler, they should only be able to access service layer interfaces. Any attempt to directly import or use repository layer code should be prevented or flagged.

**Why this priority**: This is the primary use case for layer boundaries - ensuring new code follows the correct architecture from the start. Preventing violations at development time is more valuable than detecting them later.

**Independent Test**: Can be tested by attempting to add a new transport handler that incorrectly imports a repository package - the violation should be detected.

**Acceptance Scenarios**:

1. **Given** a transport handler is being developed, **When** the developer imports a service interface, **Then** the code compiles and layer rules are satisfied
2. **Given** a transport handler is being developed, **When** the developer attempts to import a repository package directly, **Then** the layer violation is detected and reported
3. **Given** a transport handler needs data access, **When** the developer reviews available imports, **Then** only service layer interfaces are accessible/allowed

---

### User Story 2 - Developer Modifies Service Layer (Priority: P2)

A developer is adding or modifying business logic in the service layer. The service layer defines its own repository interfaces and domain models. It should not have any knowledge of transport-layer or persistence-layer implementation details.

**Why this priority**: Service layer is the core business logic - ensuring it has zero dependencies on outer layers maintains testability and separation of concerns.

**Independent Test**: Can be tested by verifying a service has no imports from transport or repository packages, and can be instantiated with mock implementations of its own interfaces.

**Acceptance Scenarios**:

1. **Given** a service is being developed, **When** the developer defines repository interfaces within the service package, **Then** the code compiles and layer rules are satisfied
2. **Given** a service is being developed, **When** the developer attempts to import a transport or repository package, **Then** the layer violation is detected
3. **Given** a service needs to be tested, **When** a developer creates unit tests, **Then** the service can be tested with mock implementations of its interfaces without any outer layer dependencies

---

### User Story 3 - Team Reviews Codebase for Architecture Compliance (Priority: P3)

A tech lead or architect wants to verify that the entire codebase follows proper layer boundaries. They should be able to run a check that reports all existing violations across the project.

**Why this priority**: Auditing is important for maintaining architectural integrity over time, but is less critical than preventing new violations.

**Independent Test**: Can be tested by running an architecture audit command that scans all packages and reports violations.

**Acceptance Scenarios**:

1. **Given** the codebase has proper layer boundaries, **When** the architecture check runs, **Then** no violations are reported
2. **Given** the codebase has layer violations, **When** the architecture check runs, **Then** all violations are listed with file locations and violation types
3. **Given** a violation report exists, **When** a developer reviews it, **Then** the report clearly identifies which package imported a forbidden package

---

### Edge Cases

- What happens when a package legitimately needs cross-cutting concerns (e.g., logging, configuration)?
  - Cross-cutting concerns in `pkg/` or `internal/config/` are exempt from layer rules as they serve all layers
- How does the system handle shared types or interfaces that multiple layers need?
  - The service layer owns the canonical domain models. Transport and repository layers import service types and translate to/from their layer-specific types (API DTOs, database entities)
- What happens when the transport layer needs to convert between API types and service types?
  - Conversion logic lives in the transport layer, which can import both API types and service types

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST enforce that transport layer packages can only import service layer packages, not repository layer packages
- **FR-002**: System MUST enforce that service layer packages do NOT import any types from transport or repository layers (service layer is dependency-free from outer layers)
- **FR-003**: System MUST allow all layers to import cross-cutting concern packages (configuration, logging, utilities)
- **FR-004**: System MUST provide a mechanism to detect and report layer boundary violations
- **FR-005**: System MUST clearly define which packages belong to which architectural layer
- **FR-006**: Violation reports MUST identify the violating file, the forbidden import, and which layer rule was broken
- **FR-007**: Service layer MUST own the canonical domain models; transport and repository layers MUST NOT export types that service layer imports
- **FR-008**: Transport layer MUST translate between API-specific types and service layer domain models
- **FR-009**: Repository layer MUST translate between persistence-specific types and service layer domain models
- **FR-010**: Repository layer MUST import service layer packages to implement service-defined interfaces and access domain models

### Layer Definitions

- **Transport Layer**: Packages that handle HTTP/RPC request/response handling (e.g., `internal/transport/`, Connect handlers). Owns API-specific types (DTOs, request/response structs). Translates to/from service domain models.
- **Service Layer**: Packages containing business logic and orchestration (e.g., `internal/*/service.go`). Owns canonical domain models and defines repository interfaces.
- **Repository Layer**: Packages handling data persistence and external data access (e.g., `internal/*/repository.go`, `internal/*/postgres.go`). Owns persistence-specific types (database entities). Implements service-defined interfaces and translates to/from service domain models.
- **Cross-cutting**: Configuration, logging, and utility packages that all layers may use

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of transport layer files have zero direct repository imports
- **SC-002**: 100% of service layer files have zero imports from transport or repository packages
- **SC-003**: Developers can verify layer compliance in under 10 seconds during local development
- **SC-004**: All layer violations are detectable before code is merged to the main branch
- **SC-005**: New developers can understand layer boundaries within 5 minutes of reading documentation

## Assumptions

- The current package structure (`internal/conversation/`, `internal/message/`, etc.) will be refactored to clearly separate service and repository code
- Layer boundary enforcement will be integrated into the development workflow (e.g., CI checks, pre-commit hooks, or compile-time enforcement)
- The team agrees on the three-layer architecture pattern (transport -> service -> repository)
- Cross-cutting concerns like `internal/config/`, `internal/llm/`, and `pkg/` packages are intentionally exempt from strict layer rules

## Out of Scope

- Refactoring the entire codebase to follow layers (this spec defines the boundaries; implementation is separate)
- Defining specific technology for enforcement (static analysis tool, linting rules, etc.)
- Changes to the external API or feature behavior - this is purely an internal architectural concern