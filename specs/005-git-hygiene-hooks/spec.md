# Feature Specification: Git Hygiene Hooks

**Feature Branch**: `005-git-hygiene-hooks`
**Created**: 2025-12-06
**Status**: Draft
**Input**: User description: "the repository should run hygiene checks on every git commit and git push. git commit hooks should run linters. git push hooks should regenerate files, go mod tidy, run unit and acceptance tests, and build the app binary. any failures to tool calls should fail the commit or push hook. during the pre commit hook, where it makes sense, we should only run on staged files (like linting)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Pre-Commit Linting on Staged Files (Priority: P1)

As a developer, I want my code to be automatically linted when I commit so that code quality issues are caught early before they enter the codebase. The linting should only run on the files I'm actually committing (staged files) to keep the process fast.

**Why this priority**: This is the first line of defense against code quality issues. Running on every commit ensures problems are caught at the earliest possible moment, before code is shared or pushed.

**Independent Test**: Can be fully tested by staging a file with linting errors, attempting to commit, and verifying the commit is blocked with clear error messages.

**Acceptance Scenarios**:

1. **Given** a developer has staged Go files with formatting issues, **When** they run `git commit`, **Then** the commit is blocked and specific formatting errors are displayed
2. **Given** a developer has staged Go files with no linting issues, **When** they run `git commit`, **Then** the commit proceeds successfully
3. **Given** a developer has staged only non-Go files (markdown, config), **When** they run `git commit`, **Then** the commit proceeds without running Go-specific linters
4. **Given** a developer has unstaged Go files with linting errors but staged files are clean, **When** they run `git commit`, **Then** the commit proceeds (only staged files are checked)

---

### User Story 2 - Pre-Push Full Validation (Priority: P2)

As a developer, I want comprehensive validation to run before my code is pushed to the remote repository so that I don't push broken code that could affect other team members or CI pipelines.

**Why this priority**: This is the final local quality gate before code reaches the shared repository. It catches issues that might not be visible in individual commits but appear when the full codebase is considered.

**Independent Test**: Can be fully tested by making changes that break tests or the build, attempting to push, and verifying the push is blocked.

**Acceptance Scenarios**:

1. **Given** a developer has commits ready to push, **When** they run `git push`, **Then** the system regenerates code, tidies dependencies, runs tests, and builds before allowing the push
2. **Given** generated code is out of sync with source files, **When** a developer runs `git push`, **Then** the push is blocked and the developer is informed which files need regeneration
3. **Given** unit tests fail, **When** a developer runs `git push`, **Then** the push is blocked with clear test failure output
4. **Given** acceptance tests fail, **When** a developer runs `git push`, **Then** the push is blocked with clear test failure output
5. **Given** the application fails to build, **When** a developer runs `git push`, **Then** the push is blocked with build error output
6. **Given** all validations pass, **When** a developer runs `git push`, **Then** the push proceeds to the remote

---

### User Story 3 - Bypass Mechanism for Exceptional Cases (Priority: P3)

As a developer, I need a way to bypass hooks in exceptional circumstances (emergency fixes, documentation-only changes, WIP commits) while being aware I'm skipping safety checks.

**Why this priority**: While hooks provide valuable safety nets, there are legitimate cases where they need to be bypassed. This should be possible but discouraged for normal workflow.

**Independent Test**: Can be fully tested by using the standard git bypass flags and verifying hooks are skipped.

**Acceptance Scenarios**:

1. **Given** a developer needs to make an emergency commit, **When** they use `git commit --no-verify`, **Then** pre-commit hooks are skipped
2. **Given** a developer needs to push without validation, **When** they use `git push --no-verify`, **Then** pre-push hooks are skipped

---

### Edge Cases

- What happens when the database is not running during pre-push acceptance tests?
  - The hook should fail with a clear message indicating the database is required
- What happens when linters are not installed on the developer's machine?
  - The hook should fail with instructions on how to install required tools
- How does the system handle commits during a rebase or merge conflict resolution?
  - Hooks should run normally; developers can use --no-verify if needed during complex operations
- What happens when go mod tidy produces changes?
  - The hook should fail and inform the developer to commit the go.mod/go.sum changes first
- What happens when generated files are modified?
  - The regeneration step should restore them; if they differ, the hook fails

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST run linting checks on staged Go files during pre-commit
- **FR-002**: System MUST only lint files that are staged for commit (not all files in the repository)
- **FR-003**: System MUST block commits when any staged file fails linting
- **FR-004**: System MUST display clear error messages indicating which files failed and why
- **FR-005**: System MUST regenerate code (protobuf, mocks) during pre-push
- **FR-006**: System MUST run `go mod tidy` during pre-push and fail if it produces changes
- **FR-007**: System MUST run unit tests during pre-push
- **FR-008**: System MUST run acceptance tests during pre-push
- **FR-009**: System MUST build the application binary during pre-push
- **FR-010**: System MUST block the push if any pre-push validation step fails
- **FR-011**: System MUST allow bypass via standard git flags (--no-verify)
- **FR-012**: System MUST provide clear, actionable output for each validation step
- **FR-013**: System MUST exit with non-zero status code on any failure
- **FR-014**: System MUST use a hook manager tool for hook installation and lifecycle management
- **FR-015**: Hook configuration MUST be version-controlled so all developers use the same hooks

### Key Entities

- **Pre-Commit Hook**: Script that runs before a commit is finalized, focusing on staged files
- **Pre-Push Hook**: Script that runs before commits are pushed to a remote, validating the full codebase
- **Staged Files**: Files added to the git index that will be included in the next commit

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Pre-commit hook completes in under 10 seconds for typical commits (1-10 files)
- **SC-002**: Developers receive feedback on commit issues within 5 seconds of running `git commit`
- **SC-003**: Pre-push hook completes within the time required to run the full test suite plus 30 seconds overhead
- **SC-004**: 100% of commits with linting errors are blocked before reaching the repository
- **SC-005**: 100% of pushes with failing tests are blocked before reaching the remote
- **SC-006**: Error messages clearly identify the failing step and provide resolution guidance

## Clarifications

### Session 2025-12-06

- Q: How should git hooks be installed/distributed to developer machines? → A: Use a hook manager tool (lefthook, husky, pre-commit framework)

## Assumptions

- Developers have Go toolchain installed and properly configured
- Developers have Docker running for acceptance tests (which require PostgreSQL)
- The existing Taskfile tasks (generate, generate:mocks, test:acceptance, build) are functional
- A hook manager tool will be used for hook installation and management
- Shell scripts are sufficient (no need for cross-platform solutions beyond macOS/Linux)

## Out of Scope

- Windows-specific hook implementations
- GUI/IDE-specific integrations
- Automatic fixing of linting issues (only reporting)
- Custom hook configuration per developer
- CI/CD pipeline integration (this feature is for local development only)
