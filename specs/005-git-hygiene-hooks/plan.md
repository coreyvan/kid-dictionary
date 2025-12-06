# Implementation Plan: Git Hygiene Hooks

**Branch**: `005-git-hygiene-hooks` | **Date**: 2025-12-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-git-hygiene-hooks/spec.md`

## Summary

Implement automated code quality checks via git hooks using Lefthook. Pre-commit hooks will run golangci-lint on staged Go files only for fast feedback (<10 seconds). Pre-push hooks will run comprehensive validation: code regeneration, dependency tidying, unit tests, acceptance tests, and binary build. All failures block the git operation with clear error messages.

## Technical Context

**Language/Version**: Shell scripts (bash), Go 1.25+ (for lefthook installation)
**Primary Dependencies**: Lefthook (hook manager), golangci-lint (Go linter)
**Storage**: N/A (configuration files only)
**Testing**: Manual testing of hook behavior, existing test suite validation
**Target Platform**: macOS/Linux (shell scripts)
**Project Type**: Developer tooling / build infrastructure
**Performance Goals**: Pre-commit <10 seconds, pre-push within test suite time + 30s overhead
**Constraints**: Must work with existing Taskfile tasks, no Node.js dependency
**Scale/Scope**: Single developer workflow, all Go files in repository

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Protobuf-First API Design | N/A | No API changes in this feature |
| II. Test-Alongside Development | PASS | Validates existing tests run before push |
| III. Simplicity & YAGNI | PASS | Single tool (Lefthook), minimal config, no over-engineering |
| IV. Observability | PASS | Clear output for each validation step per FR-012 |
| V. Content Safety | N/A | No content generation in this feature |

**Quality Gates alignment**:
- "golangci-lint for Go code style enforcement" - Implemented via pre-commit hook
- "Generated code MUST be committed" - Enforced via pre-push regeneration check
- "All PRs require passing tests before merge" - Pre-push ensures tests pass locally first

## Project Structure

### Documentation (this feature)

```text
specs/005-git-hygiene-hooks/
├── plan.md              # This file
├── research.md          # Phase 0 output (complete)
├── data-model.md        # Phase 1 output (N/A - no data model)
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# New files for this feature
lefthook.yml             # Hook configuration (repository root)
.golangci.yml            # Linter configuration (repository root)

# Modified files
Taskfile.yml             # Add lint, lint:staged, init:hooks, test:unit tasks
```

**Structure Decision**: This feature adds configuration files at repository root only. No changes to existing Go source structure. Lefthook manages hooks via its YAML config rather than raw .git/hooks scripts.

## Complexity Tracking

No constitution violations. The implementation follows Principle III (Simplicity):
- Single tool choice (Lefthook) over complex multi-tool setup
- YAML configuration over custom scripts
- Leverages existing Taskfile tasks where possible
