# Research: Git Hygiene Hooks

**Feature**: 005-git-hygiene-hooks
**Date**: 2025-12-06

## Decision 1: Hook Manager Tool

**Decision**: Use Lefthook

**Rationale**:
- Written in Go, making it fast and lightweight with no runtime dependencies
- Language-agnostic - works well for Go projects without requiring Node.js
- Supports parallel execution of hooks out of the box
- Native `{staged_files}` placeholder for filtering on staged files only
- Simple YAML configuration stored in repository
- Active maintenance by Evil Martians
- Single binary installation via `go install`

**Alternatives Considered**:
- **Husky**: Requires Node.js runtime, designed for JavaScript ecosystem, sequential execution only
- **pre-commit (Python)**: Requires Python runtime, more complex setup, community hooks focused on Python ecosystem
- **Manual scripts in .git/hooks**: Not version-controlled by default, harder to distribute to team

**Sources**:
- [Lefthook GitHub](https://github.com/evilmartians/lefthook)
- [Lefthook vs Husky comparison](https://medium.com/@yuseferi/best-tool-for-managing-git-hooks-in-golang-lefthook-vs-husky-892ff8a810b6)
- [Edopedia comparison 2025](https://www.edopedia.com/blog/lefthook-vs-husky/)

## Decision 2: Go Linting Tool

**Decision**: Use golangci-lint

**Rationale**:
- Already installed in the development environment (v1.64.5)
- Aggregates multiple linters (gofmt, goimports, govet, staticcheck, etc.) in one tool
- Fast due to parallelization and caching
- Supports filtering on specific files via command-line arguments
- Constitution mentions "golangci-lint for Go code style enforcement" in Quality Gates
- Industry standard for Go projects

**Alternatives Considered**:
- **gofmt + govet separately**: Less comprehensive, would require multiple commands
- **staticcheck alone**: Good but golangci-lint includes it plus more
- **revive**: Good alternative but golangci-lint is more comprehensive

## Decision 3: Pre-Commit Strategy

**Decision**: Run golangci-lint on staged Go files only

**Rationale**:
- Fast feedback (SC-001: under 10 seconds for typical commits)
- Only validates code being committed, not entire codebase
- Uses Lefthook's `{staged_files}` placeholder for file filtering
- Fails fast on first error to minimize wait time

**Configuration approach**:
```yaml
pre-commit:
  jobs:
    - name: lint go files
      glob: "*.go"
      run: golangci-lint run {staged_files} --new-from-rev=HEAD
```

## Decision 4: Pre-Push Strategy

**Decision**: Sequential validation with early exit on failure

**Rationale**:
- Order matters: regenerate → tidy → test → build
- Regeneration must happen before tests to ensure generated code is current
- go mod tidy must run before tests to ensure dependencies are correct
- Tests validate correctness before building
- Build is final gate to ensure binary compiles

**Execution order**:
1. `task generate` - Regenerate protobuf code
2. `task generate:mocks` - Regenerate mock implementations
3. `go mod tidy` - Tidy dependencies (fail if changes produced)
4. `go test ./...` - Run unit tests
5. `task test:acceptance` - Run acceptance tests (requires database)
6. `task build` - Build the binary

## Decision 5: Dirty Working Tree Detection

**Decision**: Fail pre-push if regeneration or go mod tidy produces uncommitted changes

**Rationale**:
- Generated code should be committed (per Constitution: "Generated code MUST be committed")
- go.mod/go.sum changes should be committed before push
- Prevents pushing code that would fail CI due to missing generated files
- Clear error message tells developer what to commit

**Implementation**:
```bash
# After regeneration/tidy
if ! git diff --quiet; then
  echo "ERROR: Working tree has uncommitted changes after regeneration"
  git status --short
  exit 1
fi
```

## Decision 6: Lefthook Installation Method

**Decision**: Install via `go install` with Taskfile setup task

**Rationale**:
- Go developers already have Go toolchain
- Single command installation: `go install github.com/evilmartians/lefthook/v2@latest`
- Taskfile task `init:hooks` will install lefthook and run `lefthook install`
- Consistent with existing Taskfile patterns in the project

## Technical Notes

### Existing Taskfile Tasks to Leverage
- `task generate` - Protobuf generation
- `task generate:mocks` - Mock generation
- `task build` - Binary build
- `task test:acceptance` - Acceptance tests

### Tasks to Add
- `task lint` - Run golangci-lint on entire codebase
- `task lint:staged` - Run golangci-lint on staged files only
- `task init:hooks` - Install lefthook and set up hooks
- `task test:unit` - Run unit tests only (currently missing)

### File Locations
- `lefthook.yml` - Hook configuration (repository root)
- `.golangci.yml` - Linter configuration (repository root, to be created)
