# Quickstart: Git Hygiene Hooks

**Feature**: 005-git-hygiene-hooks
**Date**: 2025-12-06

## Prerequisites

1. Go toolchain installed (1.25+)
2. golangci-lint installed (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
3. Docker running (for acceptance tests)
4. Database running (`task up`)

## Setup

### Install Lefthook and Configure Hooks

```bash
# Install lefthook and set up hooks
task init:hooks
```

This will:
1. Install Lefthook via `go install`
2. Run `lefthook install` to set up git hooks

## Usage Scenarios

### Scenario 1: Normal Commit with Clean Code

```bash
# Make changes to Go files
vim internal/service/message/service.go

# Stage changes
git add internal/service/message/service.go

# Commit - pre-commit hook runs golangci-lint on staged files
git commit -m "feat: add validation"
# Output:
# ╭──────────────────────────────────────╮
# │ 🥊 lefthook v2.x.x                   │
# ╰──────────────────────────────────────╯
# ┃  lint go files ━━━━━━━━━━━━━━━━━━━━━ ✓
#
# [branch abc1234] feat: add validation
#  1 file changed, 10 insertions(+)
```

### Scenario 2: Commit Blocked by Linting Errors

```bash
# Stage a file with formatting issues
git add internal/transport/transport.go

# Attempt commit
git commit -m "wip"
# Output:
# ╭──────────────────────────────────────╮
# │ 🥊 lefthook v2.x.x                   │
# ╰──────────────────────────────────────╯
# ┃  lint go files ━━━━━━━━━━━━━━━━━━━━━ ✖
#
# internal/transport/transport.go:42:1: File is not `gofmt`-ed (gofmt)
#
# ❌ Commit blocked. Fix linting errors and try again.
```

### Scenario 3: Push with All Validations Passing

```bash
# Push to remote - pre-push hook runs full validation
git push origin feature-branch
# Output:
# ╭──────────────────────────────────────╮
# │ 🥊 lefthook v2.x.x                   │
# ╰──────────────────────────────────────╯
# ┃  regenerate code ━━━━━━━━━━━━━━━━━━━ ✓
# ┃  tidy modules ━━━━━━━━━━━━━━━━━━━━━━ ✓
# ┃  unit tests ━━━━━━━━━━━━━━━━━━━━━━━━ ✓
# ┃  acceptance tests ━━━━━━━━━━━━━━━━━━ ✓
# ┃  build binary ━━━━━━━━━━━━━━━━━━━━━━ ✓
#
# To github.com:user/kid-dictionary.git
#    abc1234..def5678  feature-branch -> feature-branch
```

### Scenario 4: Push Blocked by Failing Tests

```bash
git push origin feature-branch
# Output:
# ╭──────────────────────────────────────╮
# │ 🥊 lefthook v2.x.x                   │
# ╰──────────────────────────────────────╯
# ┃  regenerate code ━━━━━━━━━━━━━━━━━━━ ✓
# ┃  tidy modules ━━━━━━━━━━━━━━━━━━━━━━ ✓
# ┃  unit tests ━━━━━━━━━━━━━━━━━━━━━━━━ ✖
#
# --- FAIL: TestSendMessage (0.05s)
#     message_test.go:42: expected "hello" got "goodbye"
# FAIL
#
# ❌ Push blocked. Fix failing tests and try again.
```

### Scenario 5: Push Blocked by Uncommitted Generated Files

```bash
git push origin feature-branch
# Output:
# ╭──────────────────────────────────────╮
# │ 🥊 lefthook v2.x.x                   │
# ╰──────────────────────────────────────╯
# ┃  regenerate code ━━━━━━━━━━━━━━━━━━━ ✖
#
# ❌ Working tree has uncommitted changes after regeneration:
#  M gen/kiddictionary/v1/message.pb.go
#  M gen/kiddictionary/v1/kiddictionaryv1connect/message.connect.go
#
# Please commit the generated files and try again.
```

### Scenario 6: Bypass Hooks for Emergency

```bash
# Emergency commit without linting
git commit --no-verify -m "hotfix: critical bug"

# Emergency push without validation
git push --no-verify origin main
```

## Taskfile Commands

```bash
# Run linter on entire codebase
task lint

# Run linter on staged files only (used by pre-commit)
task lint:staged

# Run unit tests
task test:unit

# Run acceptance tests (requires database)
task test:acceptance

# Set up hooks (run once after clone)
task init:hooks
```

## Troubleshooting

### "golangci-lint: command not found"

Install golangci-lint:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### "lefthook: command not found"

Install lefthook:
```bash
go install github.com/evilmartians/lefthook/v2@latest
```

### Acceptance tests fail with "database connection refused"

Start the database:
```bash
task up
```

### Hooks not running after git clone

Set up hooks:
```bash
task init:hooks
```
