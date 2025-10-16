# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library implementing the Functional Options Pattern using generics. It provides a flexible way to configure objects through functional options while supporting defaults, validation, and struct-based initialization.

**Core Concepts:**
- `with.Func[O]` - Generic function type for option functions
- `Defaulted` interface - Implement `SetDefaults()` to provide default values
- `Validated` interface - Implement `Validate()` to validate options
- `DefaultThenAddWith()` - Apply defaults then functional options
- `AddWith()` - Apply functional options to existing options struct

## Development Commands

### Build
```bash
go build -v ./...
```

### Test
```bash
# Run all tests (main package has no tests, only samples have tests)
go test -v ./...

# Run tests including samples
go test -v -tags samples ./...

# Run tests with coverage
go test -v -cover ./...
go test -v -cover -tags samples ./...
```

### Lint
```bash
# Run golangci-lint
golangci-lint run
```

### Pre-commit Hooks
This repository uses pre-commit hooks configured in `.pre-commit-config.yaml`:
- `go-fmt` - Format Go code
- `golangci-lint` - Lint Go code
- `go-unit-tests` - Run unit tests
- `go-mod-tidy` - Tidy go.mod automatically
- `commitlint` - Enforce conventional commit messages

Install hooks: `pre-commit install`

### Commit Messages
Commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) format (enforced by commitlint):
- `feat: add new feature`
- `fix: resolve bug`
- `chore: update dependencies`
- `docs: update documentation`

## Architecture

### Main Package (`with.go`)

The core library consists of:
- **Interfaces:**
  - `Defaulted` - Optional interface for types that have defaults
  - `Validated` - Optional interface for types that need validation
- **Types:**
  - `Func[O any]` - Generic function type that returns an error
- **Functions:**
  - `DefaultThenAddWith[O any]()` - Calls SetDefaults if available, then applies options
  - `AddWith[O any]()` - Applies options and validates if Validated is implemented
  - `Nop[O any]()` - Returns a no-op option function

Error messages automatically include the function name of failed `With...` functions for better debugging.

### Sample Code (`samples/server/`)

The sample demonstrates the pattern with a server configuration that requires `Host`, `Port`, and `Timeout`. Tests are tagged with `//go:build samples` to separate them from the main library tests.

**Constructor Patterns:**
1. `NewServer(withOptions ...)` - Starts with defaults, applies options
2. `NewServerFromOptions(options, withOptions ...)` - Starts with provided struct, applies additional options

## Key Implementation Details

- Options structs should export fields (capitalize) if callers need to construct them directly
- `With...` functions can return errors for per-option validation
- The `Validate()` method performs whole-options validation after all options are applied
- The library uses reflection to extract function names for error messages
