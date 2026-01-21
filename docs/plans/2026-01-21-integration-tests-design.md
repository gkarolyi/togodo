# Integration Test Implementation Plan

## Overview

This document outlines the plan for implementing todo.txt-cli feature parity by making the ported integration tests pass. The project has completed a refactoring to:
- Reduce the API surface of `todotxtlib/` package
- Make `cmd/` package a general-purpose business logic service layer
- Move CLI functionality into `internal/cli/` package

The goal is to run the ported bash integration tests and implement missing functionality step-by-step.

## Current State

- **Test Suite**: ~133 test cases across 20+ test files in `tests/` directory
- **Test Files**: Numbered to match original todo.txt-cli bash tests (t0000, t1000, t1010, etc.)
- **Test Framework**: Buffer-based testing using `TestEnvironment` with injected dependencies
- **Architecture**: `todotxtlib` (data) → `cmd` (business logic) → `internal/cli` (CLI wrappers)

## Guiding Principles

### 1. Minimize Refactoring
- Focus changes on `internal/cli/` command wrappers only
- Avoid touching `todotxtlib/` unless absolutely necessary for core functionality
- Avoid touching `cmd/` business logic unless the feature is completely missing
- Keep the existing architecture intact

### 2. Output Format Over Logic
- Most failures will be output format mismatches, not missing functionality
- Fix by adjusting how commands print results, not how they compute them
- Use `cmd.OutOrStdout()` for all output to enable test capture via Cobra's built-in redirection

### 3. Ship Incrementally
- Commit after each test file passes (or after each logical group of tests)
- Don't batch up multiple test files before committing
- Creates clear audit trail and makes it easy to revert if needed

### 4. Feature Gaps as Last Resort
- If a test needs a feature that doesn't exist at all, mark with `t.Skip("TODO: ...")` first
- Group missing features together and implement in a batch (Phase 2)
- Prefer fixing what exists over adding new functionality in Phase 1

### 5. Unit Testing is Mandatory
- Every function in `cmd/` must have unit tests
- Every method in `todotxtlib/` must have unit tests
- Integration tests verify end-to-end CLI behavior, unit tests verify business logic correctness
- Run all tests before committing (enforced by lefthook)

## Phase 1: Output Format Fixes

Fix output formatting to match todo.txt-cli without changing business logic.

### Test Execution Order

**Skip t0000_config initially:**
- Config tests involve Viper configuration persistence
- More complex than other commands
- Return to these after core todo.txt functionality is solid

**Start with t1000_add_list and proceed sequentially:**
```
1. t1000_add_list      - Core add/list functionality
2. t1010_add_date      - Date handling on add
3. t1040_add_priority  - Priority handling on add
4. t1100_replace       - Replace command
5. t1200_pri           - Set priority
6. t1250_listpri       - List by priority
7. t1310_listcon       - List contexts
8. t1320_listproj      - List projects
9. t1350_listall       - List all including done
10. t1400_prepend      - Prepend to todo
11. t1500_do           - Mark done
12. t1600_append       - Append to todo
13. t1700_depri        - Remove priority
14. t1800_del          - Delete todo
15. t1850_move         - Move/renumber todos
16. t1900_archive      - Archive completed
17. t1910_deduplicate  - Remove duplicates
18. t1950_report       - Generate reports
19. t2000_multiline    - Multiline todo support
```

### Repeatable Workflow (Per Test File)

#### Step 1: Run & Observe
```bash
cd tests
go test -v -run TestFileName
```
- Note which tests fail and why (missing command, wrong output, wrong exit code)
- Identify patterns (all failures in one file usually have same root cause)

#### Step 2: Categorize Failure Type

**Type A: Command doesn't exist**
- Create new command file in `internal/cli/`
- Wire up to existing `cmd/` business logic if it exists
- If business logic missing, mark test as `t.Skip("TODO: Implement X in cmd package")`

**Type B: Output format wrong**
- Fix `fmt.Fprintf(cmd.OutOrStdout(), ...)` calls in command
- Match exact output format from test expectations
- Common fixes: add/remove newlines, adjust spacing, change message text

**Type C: Exit code wrong**
- Return appropriate error from command's `RunE`
- Cobra translates errors to exit code 1 automatically

**Type D: Missing core feature**
- Skip test with `t.Skip("TODO: ...")`
- Document what's needed in business logic
- Batch these together for Phase 2 implementation

#### Step 3: Fix & Iterate
- Make minimal change to fix one test
- Re-run test file
- Repeat until all tests pass (or are explicitly skipped)

#### Step 4: Commit & Move On
```bash
git add .
git commit -m "test: make t1000_add_list tests pass"
```
- Move to next test file
- Repeat workflow

### Key Implementation Detail: Remove Presenter Abstraction

The `Presenter` abstraction is unnecessary because:
- Cobra commands have `OutOrStdout()` and `ErrOrStderr()` methods
- Tests already redirect via `env.rootCmd.SetOut(env.output)`
- Original todo.txt-cli uses plain text output only

**Changes needed:**
1. Remove `presenter` parameter from all command constructors
2. Update `internal/cli/root.go` to not inject presenter
3. Replace all `fmt.Printf(...)` with `fmt.Fprintf(cmd.OutOrStdout(), ...)`
4. Delete unused presenter files (unless used by TUI): `presenter.go`, `formatter.go`, `output.go`, `themes.go`

## Phase 2: Missing Business Logic

After Phase 1, you'll have a collection of skipped tests. This phase implements the missing features.

### Workflow for Missing Features

#### Step 1: Audit & Categorize Skipped Tests

```bash
cd tests
go test -v ./... 2>&1 | grep "TODO:"
```

**Categorize by location:**
- **Missing cmd/ functions**: Business logic doesn't exist (e.g., `cmd.Archive()`, `cmd.Deduplicate()`)
- **Missing todotxtlib/ features**: Core data operations don't exist (e.g., done.txt support, completion dates)
- **Missing CLI commands**: Need new command file in `internal/cli/` (e.g., `replace`, `prepend`, `append`)

#### Step 2: Prioritize by Dependency & Complexity

**Example grouping:**
```
Group A: Independent additions (can do in any order)
  - cmd.Prepend() - adds text to beginning of todo
  - cmd.Append() - adds text to end of todo
  - cmd.Replace() - replaces entire todo text

Group B: Depends on done.txt support
  - cmd.Archive() - moves completed todos to done.txt
  - ReadDoneFile() in tests
  - todotxtlib support for second file

Group C: Depends on date handling
  - Completion dates when marking done
  - Creation dates on add
  - Date parsing/formatting in todotxtlib
```

#### Step 3: Implement One Feature at a Time

**A. Check existing tests first**
```bash
# Before touching cmd/ or todotxtlib/
go test ./cmd -v
go test ./todotxtlib -v
```
All existing tests must pass before starting.

**B. Implement business logic in `cmd/` with unit tests**

```go
// cmd/prepend.go
func Prepend(repo todotxtlib.TodoRepository, index int, text string) (PrependResult, error) {
    // Implementation
}
```

```go
// cmd/prepend_test.go
func TestPrepend(t *testing.T) {
    // Create test repository
    // Test successful prepend
    // Test error cases (invalid index, etc.)
    // Verify repository state
}
```

**Verify unit tests pass:**
```bash
go test ./cmd -v -run TestPrepend
```

**C. If todotxtlib/ changes needed, unit test first**

Only modify `todotxtlib/` when:
- New core data operations needed (e.g., done.txt file support)
- New Todo fields required (e.g., completion date, creation date)
- New repository methods needed by multiple commands

**Principle: Push down sparingly**
- Try to implement in `cmd/` using existing todotxtlib operations first
- Only add to todotxtlib when the operation is truly fundamental to the data model
- Example: `Prepend()` might just be `Get()` + modify text + `Update()` - stays in cmd/

```go
// todotxtlib/repository.go - only if needed
func (r *FileRepository) SomeNewCoreMethod(...) error {
    // Implementation
}
```

```go
// todotxtlib/repository_test.go
func TestSomeNewCoreMethod(t *testing.T) {
    // Test thoroughly
}
```

**Verify all todotxtlib tests still pass:**
```bash
go test ./todotxtlib -v
```

**D. Wire up CLI command in `internal/cli/`**

```go
// internal/cli/prepend.go
func NewPrependCmd(repo todotxtlib.TodoRepository) *cobra.Command {
    return &cobra.Command{
        Use:   "prepend ITEM# TEXT",
        Short: "Prepend text to a todo item",
        Args:  cobra.MinimumNArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Parse index
            // Call cmd.Prepend()
            // Format output using cmd.OutOrStdout()
            return nil
        },
    }
}
```

**E. Register in root.go**
```go
rootCmd.AddCommand(NewPrependCmd(repo))
```

**F. Unskip integration test and verify**
```bash
# Remove t.Skip() from test
go test -v -run TestPrepend ./tests
```

**G. Run full test suite before committing**
```bash
go test ./cmd -v
go test ./todotxtlib -v
go test ./tests -v -run TestPrepend
```

**H. Commit with both unit and integration tests**
```bash
git add .
git commit -m "feat: implement prepend command

- Add cmd.Prepend() with unit tests
- Add cli.NewPrependCmd()
- Integration tests passing"
```

## Testing Strategy

### Test Fixtures & Data

**Fixtures are ported from bash tests:**
- Tests use `env.WriteTodoFile(content)` to set up initial state
- Content matches exactly what the original bash tests used
- No need to create separate fixture files

**Buffer-based testing:**
- `BufferReader` / `BufferWriter` already implemented
- Test helpers create buffer-backed repositories
- No filesystem I/O in tests (fast, isolated)
- Same data semantics as real files

**Example pattern:**
```go
func TestSomething(t *testing.T) {
    env := SetupTestEnv(t)

    // Set up fixture data (same as bash test)
    env.WriteTodoFile("Buy milk\nCall dentist\n(A) Important task")

    // Run command
    output, code := env.RunCommand("list")

    // Assert output matches todo.txt-cli behavior
}
```

### Pre-commit Automation

Lefthook automatically runs on every commit:
- `go vet ./...` - Package consistency checks
- `gofmt -w {files}` - Code formatting
- `go test ./...` - Full test suite (unit + integration)
- `go build` - Build verification

**Trust the automation** - if lefthook passes, all tests passed.

### Unit Testing Patterns

**For cmd/ functions:**
```go
func TestAdd(t *testing.T) {
    // Use test helpers from cmd/test_helpers.go
    repo := NewTestRepository(t)

    result, err := Add(repo, []string{"test task"})

    // Assert no error
    // Assert result contains expected todo
    // Assert repository contains the todo
}
```

**For todotxtlib/ methods:**
```go
func TestRepositoryAdd(t *testing.T) {
    // Use buffer-based repository
    buffer := &bytes.Buffer{}
    repo, _ := todotxtlib.NewFileRepository(
        todotxtlib.NewBufferReader(buffer),
        todotxtlib.NewBufferWriter(buffer),
    )

    todo, err := repo.Add("test task")

    // Assert behavior
}
```

## Success Criteria

**Phase 1 Complete:**
- All integration tests either pass or are explicitly skipped with documented reason
- No changes to `todotxtlib/` or `cmd/` unless absolutely necessary
- All output format matches todo.txt-cli exactly
- Committed incrementally (one test file or logical group per commit)

**Phase 2 Complete:**
- All skipped tests are unskipped and passing
- All new business logic has unit tests
- All changes to `todotxtlib/` have unit tests
- Full test suite passes (`go test ./...`)
- Feature parity with todo.txt-cli achieved

## Next Steps

1. Remove Presenter abstraction from commands
2. Start with t1000_add_list test file
3. Follow repeatable workflow through remaining test files
4. Document skipped tests for Phase 2
5. Implement missing features with unit tests
6. Achieve full feature parity
