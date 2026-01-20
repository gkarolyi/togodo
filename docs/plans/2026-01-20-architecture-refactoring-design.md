# Architecture Refactoring Design

**Date:** 2026-01-20
**Status:** Approved
**Author:** Design session with user

## Problem Statement

The current architecture doesn't support the integration tests effectively. The tests expect CLI output that matches todo.txt-cli (like "TODO: 1 added"), but the cmd package is also used by the TUI, which doesn't need or want this output.

Additionally:
- The service layer provides minimal value and adds unnecessary indirection
- The TodoRepository interface has 20+ methods, many overly specific
- No clear separation between business logic and presentation concerns

## Goals

1. **Enable CLI/TUI to share business logic** while presenting output differently
2. **Simplify the todotxtlib API** to be a clean, reusable library
3. **Achieve feature parity** with todo.txt-cli through comprehensive integration tests
4. **Maintain clean separation** between data, business logic, and presentation layers

## Architecture Overview

### Three-Layer Architecture

```
┌─────────────────────────────────────────────────┐
│         Presentation Layer                      │
│  ┌──────────────┐        ┌──────────────┐      │
│  │ internal/cli │        │ internal/tui │      │
│  │ (Cobra)      │        │ (Bubble Tea) │      │
│  └──────┬───────┘        └──────┬───────┘      │
│         │                       │               │
│         └───────────┬───────────┘               │
└─────────────────────┼─────────────────────────┘
                      │
┌─────────────────────┼─────────────────────────┐
│         Business Logic Layer                   │
│              ┌──────▼───────┐                  │
│              │  cmd package │                  │
│              │  (functions) │                  │
│              └──────┬───────┘                  │
└─────────────────────┼─────────────────────────┘
                      │
┌─────────────────────┼─────────────────────────┐
│         Data Layer                             │
│              ┌──────▼───────┐                  │
│              │  todotxtlib  │                  │
│              │ (repository) │                  │
│              └──────────────┘                  │
└─────────────────────────────────────────────────┘
```

**Layer 1: todotxtlib (Data & CRUD)**
- Repository pattern with CRUD operations on todo.txt
- Filtering, sorting, searching primitives
- No orchestration or business logic
- Pure data layer usable by both CLI and external libraries

**Layer 2: cmd (Business Logic)**
- Business logic functions (not Cobra commands)
- Orchestrates repository operations for specific commands
- Handles validation, multi-step operations, sort/save coordination
- Returns typed result structs specific to each command
- Example: `cmd.Add(repo TodoRepository, args []string) (AddResult, error)`

**Layer 3: Presentation (CLI & TUI)**
- **internal/cli/**: Cobra command wrappers that format cmd results for CLI (matching todo.txt-cli)
- **internal/tui/**: Bubble Tea interface that calls cmd functions and displays in UI
- Both call identical cmd business logic
- Each formats output for its context

### What Changes

**Removed:**
- `todotxtlib/service.go` - Business logic moves to cmd package
- Overly specific repository methods (AddContext, RemoveContext, etc.)

**Added:**
- `cmd/` functions become pure business logic (not just Cobra wrappers)
- `internal/cli/` for Cobra command wrappers
- Result types for each command

**Modified:**
- `todotxtlib/repository.go` - Simplified interface (20+ → 12 methods)
- `internal/tui/` - Calls cmd functions instead of repository directly

## Detailed Design

### 1. Simplified TodoRepository Interface

**Current problems:**
- 20+ methods, many overly specific (AddContext, RemoveContext, AddProject, RemoveProject)
- No `Get(index)` method - forces awkward ListAll() + array access
- Inconsistent Save()/WriteToString() pattern
- Search is just a thin wrapper around Filter

**Proposed interface (12 methods):**

```go
type TodoRepository interface {
    // Core CRUD
    Add(text string) (Todo, error)
    Get(index int) (Todo, error)              // NEW - essential for modifications
    Update(index int, todo Todo) (Todo, error)
    Remove(index int) (Todo, error)

    // Common operations (convenience for TUI and common patterns)
    ToggleDone(index int) (Todo, error)
    SetPriority(index int, priority string) (Todo, error)

    // Queries
    ListAll() ([]Todo, error)
    Filter(filter Filter) ([]Todo, error)
    ListContexts() ([]string, error)
    ListProjects() ([]string, error)

    // Sorting & Persistence
    Sort(sort Sort)  // nil = default sort
    Save() error
}
```

**Key changes:**

1. **Added `Get(index)`** - Essential for commands that modify individual todos (append, prepend, replace)
2. **Removed Search()** - Use `Filter(Filter{Text: query})` instead
3. **Removed WriteToString()** - Use Save() with BufferWriter instead (consistent with writer abstraction)
4. **Removed SortDefault()** - Use `Sort(nil)` for default sort (idiomatic Go pattern)
5. **Removed micro-methods** - SetContexts, AddContext, RemoveContext, etc. Use Get() + Update() instead
6. **Kept convenience methods** - ToggleDone and SetPriority remain (used by TUI, very common operations)

**Rationale:**
- Get + Update pattern is more flexible than many specific methods
- Writer abstraction already handles both file and buffer cases
- Reduces interface size while maintaining all functionality
- More idiomatic Go (smaller interfaces, nil as default)

### 2. cmd Package - Business Logic Functions

The cmd package contains pure business logic functions that orchestrate repository operations.

**Structure:**
```
cmd/
  add.go       - Add(repo, args) (AddResult, error)
  list.go      - List(repo, filter) (ListResult, error)
  do.go        - Do(repo, indices) (DoResult, error)
  delete.go    - Delete(repo, index) (DeleteResult, error)
  archive.go   - Archive(repo, doneRepo) (ArchiveResult, error)
  pri.go       - SetPriority(repo, indices, priority) (PriResult, error)
  depri.go     - RemovePriority(repo, indices) (DepriResult, error)
  append.go    - Append(repo, index, text) (AppendResult, error)
  prepend.go   - Prepend(repo, index, text) (PrependResult, error)
  replace.go   - Replace(repo, index, text) (ReplaceResult, error)
  // ... etc for remaining commands
```

**Key patterns:**

1. **Dependencies as parameters** (no globals, easy testing)
2. **Each command has its own result type** (specific to needs)
3. **Functions orchestrate operations** (validate → operate → sort → save)

**Example implementation:**

```go
// cmd/add.go
type AddResult struct {
    Todo       Todo
    LineNumber int
}

func Add(repo TodoRepository, args []string) (AddResult, error) {
    // Validate
    if len(args) == 0 {
        return AddResult{}, fmt.Errorf("task text required")
    }

    // Business logic: join args into single task
    text := strings.Join(args, " ")

    // Add to repository
    todo, err := repo.Add(text)
    if err != nil {
        return AddResult{}, fmt.Errorf("failed to add todo: %w", err)
    }

    // Sort and save
    repo.Sort(nil)  // default sort
    if err := repo.Save(); err != nil {
        return AddResult{}, fmt.Errorf("failed to save: %w", err)
    }

    // Calculate line number after sort
    lineNum, err := findLineNumber(repo, todo)
    if err != nil {
        return AddResult{}, fmt.Errorf("failed to find line number: %w", err)
    }

    return AddResult{Todo: todo, LineNumber: lineNum}, nil
}
```

**Example result types:**

```go
// Each command defines what it needs to return
type AddResult struct {
    Todo       Todo
    LineNumber int
}

type ListResult struct {
    Todos       []Todo
    TotalCount  int
    ShownCount  int
}

type ArchiveResult struct {
    ArchivedTodos []Todo
}

type DeleteResult struct {
    DeletedTodo Todo
    LineNumber  int
}
```

**Benefits:**
- Pure functions, easy to test
- No Cobra dependencies in business logic
- Both CLI and TUI can use identical logic
- Result types provide all data needed for presentation

### 3. internal/cli Package - CLI Presentation Layer

The internal/cli package provides Cobra command wrappers that format cmd results to match todo.txt-cli output exactly.

**Structure:**
```
internal/cli/
  add.go        - NewAddCmd(repo) *cobra.Command
  list.go       - NewListCmd(repo) *cobra.Command
  do.go         - NewDoCmd(repo) *cobra.Command
  delete.go     - NewDeleteCmd(repo) *cobra.Command
  // ... etc
  presenter.go  - Todo formatting (already exists)
  formatter.go  - Output utilities (already exists)
```

**Pattern for each command:**

```go
// internal/cli/add.go
func NewAddCmd(repo todotxtlib.TodoRepository) *cobra.Command {
    return &cobra.Command{
        Use:   "add [TASK]",
        Short: "Add a new todo item",
        Args:  cobra.MinimumNArgs(1),
        Aliases: []string{"a"},
        RunE: func(cmd *cobra.Command, args []string) error {
            // Call business logic
            result, err := togocmd.Add(repo, args)
            if err != nil {
                return err
            }

            // Format output to match todo.txt-cli exactly
            fmt.Printf("%d %s\n", result.LineNumber, result.Todo.String())
            fmt.Printf("TODO: %d added.\n", result.LineNumber)
            return nil
        },
    }
}
```

**Responsibilities:**
- Cobra command setup (flags, usage, aliases, validation)
- Call cmd business logic functions
- Format results to match todo.txt-cli output exactly
- Handle CLI-specific concerns (exit codes, error messages)
- Use existing presenter/formatter utilities

**Note on pattern:** We're using factory functions (`NewAddCmd()`) rather than package-level variables (standard Cobra pattern) to enable dependency injection. This makes testing easier and keeps dependencies explicit, even though it's slightly non-standard for Cobra.

### 4. TUI Integration

The TUI calls cmd business logic functions and refreshes display from repository queries.

**Current approach (direct repository calls):**
```go
// internal/tui/update.go (current)
m.repository.Add(m.input.Value())
m.repository.ToggleDone(i)
m.repository.SetPriority(i, priority)
```

**New approach (cmd functions + refresh):**
```go
// internal/tui/update.go (proposed)
_, err := togocmd.Add(m.repository, []string{m.input.Value()})
if err != nil {
    // handle error in UI
}

// Refresh display based on current view
if m.filtering {
    m.choices, _ = m.repository.Filter(Filter{Text: m.filter})
} else {
    m.choices, _ = m.repository.ListAll()
}
```

**Key changes:**
1. Call cmd functions instead of repository directly
2. Refresh display with ListAll()/Filter() after operations
3. cmd functions handle sort/save coordination
4. TUI just displays current state

**Benefits:**
- TUI and CLI use identical business logic
- Bug fixes benefit both interfaces
- Simpler TUI code (no orchestration logic)
- Always shows consistent state after operations

## Migration Strategy

This is a significant refactoring. The strategy keeps things working while migrating incrementally.

### Step 1: Simplify todotxtlib (Foundation)

**Goal:** Clean up the data layer foundation

**Tasks:**
- Simplify TodoRepository interface (remove 8+ methods, add Get())
- Update FileRepository implementation
- Remove service.go entirely
- Update writer usage (remove WriteToString, use Save with BufferWriter)
- Update Sort to accept nil for default

**Impact:** Breaking changes to cmd package
**Validation:** `go test ./todotxtlib`

**Files modified:**
- `todotxtlib/repository.go`
- `todotxtlib/service.go` (deleted)

### Step 2: Create cmd business logic functions (Core)

**Goal:** Establish business logic layer

**Tasks:**
- Create result types for each command
- Write cmd functions for existing commands (add, list, do, pri, tidy)
- Add unit tests for each function
- Keep existing Cobra commands temporarily (parallel implementation)

**Impact:** No breaking changes yet (cmd grows without changing existing code)
**Validation:** Unit tests for each cmd function

**Files created:**
- `cmd/add.go` - Add function + AddResult type
- `cmd/list.go` - List function + ListResult type
- `cmd/do.go` - Do function + DoResult type
- `cmd/pri.go` - SetPriority function + PriResult type
- `cmd/tidy.go` - Tidy function + TidyResult type
- `cmd/*_test.go` - Unit tests for each

### Step 3: Migrate CLI to use cmd functions (Presentation)

**Goal:** Separate Cobra wrappers from business logic

**Tasks:**
- Create internal/cli/ directory
- Move/rewrite Cobra wrappers to call cmd functions
- Format output to match todo.txt-cli exactly
- Wire up in cmd/root.go
- Remove old cmd implementations

**Impact:** CLI commands work through new path, output matches tests
**Validation:** Integration tests in `/tests` start passing

**Files created:**
- `internal/cli/add.go` - NewAddCmd wrapper
- `internal/cli/list.go` - NewListCmd wrapper
- `internal/cli/do.go` - NewDoCmd wrapper
- `internal/cli/pri.go` - NewPriCmd wrapper
- `internal/cli/tidy.go` - NewTidyCmd wrapper

**Files modified:**
- `cmd/root.go` - Register commands from internal/cli

### Step 4: Migrate TUI to use cmd functions (Presentation)

**Goal:** TUI uses same business logic as CLI

**Tasks:**
- Update internal/tui/update.go to call cmd functions
- Add refresh logic (ListAll/Filter after operations)
- Manual testing of TUI functionality

**Impact:** TUI behavior may change slightly (more consistent with CLI)
**Validation:** Manual TUI testing, verify all operations work

**Files modified:**
- `internal/tui/update.go`
- `internal/tui/model.go` (if needed)

### Step 5: Implement missing commands (Feature parity)

**Goal:** Add remaining commands following established pattern

**Tasks:**
- For each missing command:
  1. Create cmd function with result type
  2. Create internal/cli wrapper
  3. Run integration tests
  4. Fix until tests pass
- Commands to implement: del, archive, depri, append, prepend, replace, listpri, listcon, listproj, listall, deduplicate, report, move

**Impact:** Feature parity with todo.txt-cli achieved
**Validation:** Integration tests pass for each command

**Files created:**
- `cmd/delete.go`, `internal/cli/delete.go`
- `cmd/archive.go`, `internal/cli/archive.go`
- (... and so on for each command)

## Testing Strategy

### Unit Tests
- Each cmd function gets comprehensive unit tests
- Use BufferWriter for repository in tests
- Test edge cases, validation, error handling

### Integration Tests
- Existing tests in `/tests` validate CLI behavior
- Tests guide implementation (TDD approach)
- Must match todo.txt-cli output exactly

### TUI Tests
- Manual testing for now (TUI testing is complex)
- Verify all operations work after refactoring
- Consider automated tests later if needed

## Open Questions & Future Considerations

These will be addressed as we implement, guided by the integration tests:

1. **done.txt support** - Archive command needs a second repository for done.txt
2. **Line number tracking** - Tests expect line numbers matching displayed list (especially after filtering)
3. **Error handling patterns** - Consistent error messages across commands
4. **Date handling** - Auto-date support for add command (configurable)
5. **Sorting behavior** - Match todo.txt-cli sorting exactly

The integration tests will tell us exactly what's needed for each of these.

## Success Criteria

### Technical
- [ ] TodoRepository has ≤12 methods
- [ ] No service layer
- [ ] cmd functions are pure (no Cobra dependencies)
- [ ] CLI and TUI share business logic
- [ ] All todotxtlib tests pass
- [ ] All cmd unit tests pass

### Feature Parity
- [ ] All critical issues fixed (from FEATURE_GAP_REPORT.md)
- [ ] 90%+ integration tests passing
- [ ] CLI output matches todo.txt-cli exactly
- [ ] TUI functionality preserved

### Code Quality
- [ ] Clean separation of concerns
- [ ] todotxtlib usable as standalone library
- [ ] Easy to add new commands (clear pattern)
- [ ] Good test coverage

## Summary

This refactoring establishes a clean three-layer architecture:

1. **todotxtlib** - Simple, focused repository for todo.txt CRUD (12 methods)
2. **cmd** - Pure business logic functions with typed results
3. **CLI/TUI** - Presentation layers that share business logic but format differently

The migration strategy is incremental and safe, with tests validating each step. Once complete, achieving feature parity becomes straightforward: follow the established pattern for each new command, let the integration tests guide implementation.

The architecture makes todotxtlib a useful standalone library while enabling the togodo CLI to achieve complete parity with todo.txt-cli through comprehensive integration testing.
