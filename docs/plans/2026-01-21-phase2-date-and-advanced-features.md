# Phase 2: Date Support and Advanced Features Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement Phase 2 features to achieve full todo.txt-cli parity - date support, configuration management, advanced delete/do features, move command, multiline support, and UI enhancements.

**Architecture:** Follow TDD workflow established in Phase 1 - run test, implement minimal fix, verify, commit. Focus on `cmd/` business logic with unit tests, then `internal/cli/` wrappers. Extend `todotxtlib/` only when core functionality is missing. Prioritize features by dependency order and user impact.

**Tech Stack:** Go 1.24.4, Cobra CLI, Viper config, todo.txt format specification

---

## Context

**Phase 1 Status:** ✅ Complete - 40/40 tests passing (100%)

**Phase 2 Scope:** 16 skipped tests covering:
- Date support (creation dates, completion dates)
- Configuration management (read/write/list)
- Advanced command features (multiple item operations)
- Move command (inter-file task movement)
- Multiline task support
- UI enhancements (highlighting, plain mode, short help)

**Working Directory:** `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity`

---

## CRITICAL: Tool Usage Rules

**YOU MUST follow these tool usage rules strictly.**

### Use Specialized Tools (NOT bash)

- **Read tool** - For viewing file contents (NOT cat/head/tail/less)
- **Write tool** - For creating new files (NOT echo >/cat <<EOF)
- **Edit tool** - For modifying existing files (NOT sed/awk/perl)
- **Glob tool** - For finding files by pattern (NOT find/ls)
- **Grep tool** - For searching file contents (NOT grep/rg)

### Use Bash ONLY For

1. **Running tests**: `go test -v -run TestName`
2. **Git operations**: `git add`, `git commit`, `git status`
3. **Go commands**: `go build`, `go vet`

### Use Serena's Semantic Search Tools (When Available)

**Serena provides powerful semantic code navigation tools that should be preferred over basic text search:**

- **mcp__plugin_serena_serena__find_symbol** - Find functions, classes, methods by name pattern
  - Use instead of Grep when looking for specific symbols (e.g., "find ConfigRead function")
  - Supports wildcards and path patterns: `ConfigRead`, `Config*`, `/Config/Read`
  - Returns symbol metadata (location, signature) without reading full files
  - Example: "Find the Add function in cmd package" → `find_symbol` with pattern `Add` and path `cmd/`

- **mcp__plugin_serena_serena__get_symbols_overview** - Get high-level file structure
  - Use FIRST when examining a new file to understand its organization
  - Shows all functions, classes, types without full content
  - Much faster than reading entire file
  - Example: "What's in config.go?" → `get_symbols_overview` on `cmd/config.go`

- **mcp__plugin_serena_serena__find_referencing_symbols** - Find where a symbol is used
  - Use when you need to understand impact of changes
  - Shows all call sites with code snippets
  - Example: "Where is ConfigRead called?" → `find_referencing_symbols` for `ConfigRead`

- **mcp__plugin_serena_serena__search_for_pattern** - Advanced regex search with context
  - Use instead of Grep for complex patterns or when you need context lines
  - Supports file filtering by glob patterns
  - Example: Search for error handling patterns in config code

**When to use what:**
- **Symbol lookup** → Use `find_symbol` (not Grep)
- **File overview** → Use `get_symbols_overview` (not Read full file)
- **Find usages** → Use `find_referencing_symbols` (not Grep)
- **File content** → Use Read tool only after confirming the file exists
- **Text patterns** → Use `search_for_pattern` for complex searches, Grep for simple ones

**Workflow example:**
1. Use `get_symbols_overview` to understand file structure
2. Use `find_symbol` to locate specific functions
3. Use Read tool to view implementation details
4. Use `find_referencing_symbols` to understand usage

### Integration Tests Already Exist

- **DO NOT create new test files**
- **DO NOT create temporary test files in /tmp**
- The integration tests ALREADY EXIST in `/tests/` directory
- Just run them: `go test -v ./tests -run TestName`

**If you use bash for file operations or create temporary files, you are doing it wrong. Stop and use the correct tool.**

---

## Test Files to Process (Priority Order)

### High Priority - Core Features
1. ✅ Config management (t0000_config_test.go) - Foundation for other features
2. ✅ Date support in add (t1010_add_date_test.go, t1040_add_priority_test.go)
3. ✅ Advanced delete (t1800_del_test.go - multiple items, term matching)
4. ✅ Advanced do (t1500_do_test.go - multiple items, flags)

### Medium Priority - Additional Commands
5. ✅ Move command (t1850_move_test.go)
6. ✅ Deduplicate enhancements (t1910_deduplicate_test.go)

### Lower Priority - UI/UX Enhancements
7. ✅ Listall enhancements (t1350_listall_test.go - highlighting)
8. ✅ Priority plain flag (t1200_pri_test.go)
9. ✅ Short help (t2100_help_test.go)
10. ✅ Multiline support (t2000_multiline_test.go)

---

## Task Group 1: Configuration Management Foundation

**Goal:** Implement config read/write/list to enable feature configuration

---

### Task 1: Implement Config Read Command

**Files:**
- Test: `tests/t0000_config_test.go` (already exists)
- Create: `cmd/config.go`
- Create: `cmd/config_test.go`
- Modify: `internal/cli/config.go` (may already exist, add read subcommand)

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestConfigRead`
Expected: SKIP or FAIL - "config" command may exist but missing read functionality

**Step 2: Examine existing config command**

Use `mcp__plugin_serena_serena__get_symbols_overview` on `internal/cli/config.go` to see what functions exist.
Then use `mcp__plugin_serena_serena__find_symbol` with pattern `NewConfigCmd` to locate the command definition.
Finally, use Read tool on the file to see full implementation details.

**Step 3: Create business logic for config read**

Create: `cmd/config.go`

```go
package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

// ConfigReadResult contains the result of reading a config value
type ConfigReadResult struct {
	Key   string
	Value interface{}
	Found bool
}

// ConfigRead reads a configuration value by key
func ConfigRead(key string) (ConfigReadResult, error) {
	if !viper.IsSet(key) {
		return ConfigReadResult{
			Key:   key,
			Found: false,
		}, fmt.Errorf("configuration key '%s' not found", key)
	}

	return ConfigReadResult{
		Key:   key,
		Value: viper.Get(key),
		Found: true,
	}, nil
}
```

**Step 4: Create unit test for config read**

Create: `cmd/config_test.go`

```go
package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestConfigRead(t *testing.T) {
	// Setup test config
	viper.Set("test_key", "test_value")
	defer viper.Set("test_key", nil)

	result, err := ConfigRead("test_key")
	if err != nil {
		t.Fatalf("ConfigRead failed: %v", err)
	}

	if !result.Found {
		t.Error("Expected key to be found")
	}

	if result.Value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", result.Value)
	}
}

func TestConfigReadNotFound(t *testing.T) {
	result, err := ConfigRead("nonexistent_key")
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}

	if result.Found {
		t.Error("Expected key not to be found")
	}
}
```

**Step 5: Run unit test**

Run: `go test ./cmd -v -run TestConfigRead`
Expected: PASS

**Step 6: Update CLI wrapper to support read**

Use Read tool on `internal/cli/config.go`, then Edit tool to add read subcommand if needed.

Expected CLI behavior:
```bash
togodo config <key>
# Output: <value>
```

**Step 7: Run integration test**

Run: `go test -v ./tests -run TestConfigRead`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/config.go cmd/config_test.go internal/cli/config.go
git commit -m "feat: implement config read command

- Add ConfigRead() business logic
- Add unit tests for config read
- Update CLI wrapper to support config <key>
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2: Implement Config Write Command

**Files:**
- Test: `tests/t0000_config_test.go` (already exists)
- Modify: `cmd/config.go`
- Modify: `cmd/config_test.go`
- Modify: `internal/cli/config.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestConfigWrite`
Expected: SKIP or FAIL

**Step 2: Add config write business logic**

Modify: `cmd/config.go`

Add after ConfigRead:

```go
// ConfigWriteResult contains the result of writing a config value
type ConfigWriteResult struct {
	Key      string
	OldValue interface{}
	NewValue string
	Created  bool
}

// ConfigWrite writes a configuration value by key
func ConfigWrite(key string, value string) (ConfigWriteResult, error) {
	oldValue := viper.Get(key)
	created := !viper.IsSet(key)

	viper.Set(key, value)

	// Save config to file
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if err := viper.SafeWriteConfig(); err != nil {
			return ConfigWriteResult{}, fmt.Errorf("failed to write config: %w", err)
		}
	}

	return ConfigWriteResult{
		Key:      key,
		OldValue: oldValue,
		NewValue: value,
		Created:  created,
	}, nil
}
```

**Step 3: Add unit test for config write**

Modify: `cmd/config_test.go`

Add:

```go
func TestConfigWrite(t *testing.T) {
	result, err := ConfigWrite("test_key", "new_value")
	if err != nil {
		t.Fatalf("ConfigWrite failed: %v", err)
	}

	if result.NewValue != "new_value" {
		t.Errorf("Expected 'new_value', got '%s'", result.NewValue)
	}

	// Verify it was actually set
	if viper.GetString("test_key") != "new_value" {
		t.Error("Config value was not persisted")
	}
}

func TestConfigWriteUpdate(t *testing.T) {
	viper.Set("existing_key", "old_value")
	defer viper.Set("existing_key", nil)

	result, err := ConfigWrite("existing_key", "updated_value")
	if err != nil {
		t.Fatalf("ConfigWrite failed: %v", err)
	}

	if result.Created {
		t.Error("Expected update, not creation")
	}

	if result.OldValue != "old_value" {
		t.Errorf("Expected old value 'old_value', got '%v'", result.OldValue)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestConfigWrite`
Expected: PASS

**Step 5: Update CLI wrapper to support write**

Use Edit tool on `internal/cli/config.go` to add write functionality.

Expected CLI behavior:
```bash
togodo config <key> <value>
# Output: Configuration updated
```

**Step 6: Run integration test**

Run: `go test -v ./tests -run TestConfigWrite`
Expected: PASS

**Step 7: Commit**

```bash
git add cmd/config.go cmd/config_test.go internal/cli/config.go
git commit -m "feat: implement config write command

- Add ConfigWrite() business logic
- Add unit tests for config write
- Update CLI wrapper to support config <key> <value>
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3: Implement Config List Command

**Files:**
- Test: `tests/t0000_config_test.go` (already exists)
- Modify: `cmd/config.go`
- Modify: `cmd/config_test.go`
- Modify: `internal/cli/config.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestConfigList`
Expected: SKIP or FAIL

**Step 2: Add config list business logic**

Modify: `cmd/config.go`

Add:

```go
// ConfigListResult contains all configuration settings
type ConfigListResult struct {
	Settings map[string]interface{}
}

// ConfigList lists all configuration settings
func ConfigList() (ConfigListResult, error) {
	settings := viper.AllSettings()

	return ConfigListResult{
		Settings: settings,
	}, nil
}
```

**Step 3: Add unit test for config list**

Modify: `cmd/config_test.go`

Add:

```go
func TestConfigList(t *testing.T) {
	// Set some test config values
	viper.Set("key1", "value1")
	viper.Set("key2", "value2")
	defer func() {
		viper.Set("key1", nil)
		viper.Set("key2", nil)
	}()

	result, err := ConfigList()
	if err != nil {
		t.Fatalf("ConfigList failed: %v", err)
	}

	if len(result.Settings) == 0 {
		t.Error("Expected settings to be returned")
	}

	if result.Settings["key1"] != "value1" {
		t.Error("Expected key1 to be in settings")
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestConfigList`
Expected: PASS

**Step 5: Update CLI wrapper to support list**

Use Edit tool on `internal/cli/config.go` to add list functionality.

Expected CLI behavior:
```bash
togodo config
# Output: Lists all key=value pairs
```

**Step 6: Run integration test**

Run: `go test -v ./tests -run TestConfigList`
Expected: PASS

**Step 7: Commit**

```bash
git add cmd/config.go cmd/config_test.go internal/cli/config.go
git commit -m "feat: implement config list command

- Add ConfigList() business logic
- Add unit tests for config list
- Update CLI wrapper to support config (no args)
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 2: Date Support

**Goal:** Add creation and completion date support to match todo.txt format

---

### Task 4: Implement Creation Date in Add Command

**Files:**
- Test: `tests/t1010_add_date_test.go` (already exists)
- Modify: `todotxtlib/todo.go` (add date parsing)
- Modify: `todotxtlib/parser.go` or create if needed
- Modify: `cmd/add.go` (add date handling)
- Modify: `cmd/add_test.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestAddWithDate`
Expected: SKIP or FAIL

**Step 2: Read test to understand date format**

Use Read tool on `tests/t1010_add_date_test.go` to see expected format.

Expected format: `YYYY-MM-DD task text` or auto-added date

**Step 3: Add date fields to Todo struct**

Use `mcp__plugin_serena_serena__find_symbol` with pattern `Todo` and path `todotxtlib/todo.go` to locate the struct.
Use Read tool to view the full struct definition.
Then use Edit tool to add:

```go
type Todo struct {
	// ... existing fields ...
	CreationDate   string // YYYY-MM-DD format
	CompletionDate string // YYYY-MM-DD format (only when Done is true)
}
```

**Step 4: Update todo parser to extract dates**

Use `mcp__plugin_serena_serena__find_symbol` with pattern `Parse*` and path `todotxtlib/` to find existing parser functions.
If parser doesn't exist, create it in `todotxtlib/`:

```go
// ParseTodo parses a todo.txt format line
// Format: x (priority) YYYY-MM-DD YYYY-MM-DD text @context +project
// Where first date is completion date, second is creation date
func ParseTodo(text string) Todo {
	// Extract dates according to todo.txt spec
	// If done: x (A) 2026-01-15 2026-01-10 task text
	// If not done: (A) 2026-01-10 task text
}
```

**Step 5: Update Add command to support dates**

Use `mcp__plugin_serena_serena__find_symbol` with pattern `Add` and path `cmd/add.go` to locate the Add function.
Use `mcp__plugin_serena_serena__find_referencing_symbols` to see where Add is called from CLI layer.
Then modify: `cmd/add.go`

Check if `--date` flag or config enables auto-dating, add creation date if so.

**Step 6: Add unit test for date parsing**

Modify: `cmd/add_test.go`

Add tests for:
- Auto-adding creation date
- Preserving manually specified creation date
- Format validation

**Step 7: Run unit test**

Run: `go test ./cmd -v -run TestAdd`
Expected: PASS

**Step 8: Run integration test**

Run: `go test -v ./tests -run TestAddWithDate`
Expected: PASS

**Step 9: Commit**

```bash
git add todotxtlib/todo.go todotxtlib/parser.go cmd/add.go cmd/add_test.go
git commit -m "feat: add creation date support to add command

- Add CreationDate field to Todo struct
- Update parser to extract dates from todo.txt format
- Add auto-dating support to add command
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 5: Implement Priority with Date in Add Command

**Files:**
- Test: `tests/t1040_add_priority_test.go` (already exists)
- Modify: `cmd/add.go`
- Modify: `cmd/add_test.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestAddWithPriority`
Expected: SKIP or FAIL

**Step 2: Read test to understand format**

Use Read tool on `tests/t1040_add_priority_test.go`.

Expected format: `(A) YYYY-MM-DD task text` (priority before date)

**Step 3: Update Add to parse priority in input**

Modify: `cmd/add.go`

Allow adding tasks with priority prefix in the input text:
- `togodo add "(A) task"` → creates `(A) 2026-01-21 task`

**Step 4: Update parser to handle priority + date**

Ensure parser correctly handles: `(A) YYYY-MM-DD task text`

**Step 5: Add unit test**

Modify: `cmd/add_test.go`

Test priority + date combinations.

**Step 6: Run unit test**

Run: `go test ./cmd -v -run TestAdd`
Expected: PASS

**Step 7: Run integration test**

Run: `go test -v ./tests -run TestAddWithPriority`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/add.go cmd/add_test.go todotxtlib/parser.go
git commit -m "feat: support priority with date in add command

- Parse priority from input text
- Maintain proper order: (priority) date text
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 6: Implement Completion Date in Do Command

**Files:**
- Test: `tests/t1500_do_test.go` (check for date requirements)
- Modify: `cmd/do.go`
- Modify: `cmd/do_test.go`

**Step 1: Check if test requires completion dates**

Run: `go test -v ./tests -run TestDoAlreadyDone`
Look for date-related assertions.

**Step 2: Update Do command to add completion date**

Modify: `cmd/do.go`

When marking task done, add completion date:
- Format: `x 2026-01-21 2026-01-10 task text`
- First date: completion, second date: creation (if exists)

**Step 3: Add unit test**

Modify: `cmd/do_test.go`

Test completion date is added correctly.

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestDo`
Expected: PASS

**Step 5: Run integration test**

Run: `go test -v ./tests -run TestDoAlreadyDone`
Expected: PASS if date support required, SKIP if not

**Step 6: Commit if changes needed**

```bash
git add cmd/do.go cmd/do_test.go
git commit -m "feat: add completion date when marking tasks done

- Add current date as completion date
- Maintain creation date if present
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 3: Advanced Delete Features

**Goal:** Support multiple item deletion and term-based deletion

---

### Task 7: Implement Multiple Item Delete

**Files:**
- Test: `tests/t1800_del_test.go` (TestDelMultiple)
- Modify: `cmd/del.go`
- Modify: `cmd/del_test.go`
- Modify: `internal/cli/del.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestDelMultiple`
Expected: SKIP or FAIL

**Step 2: Read test to understand syntax**

Use Read tool on test file.

Expected syntax: `togodo del 1 3 5` (space-separated line numbers)

**Step 3: Update Del business logic for multiple items**

Modify: `cmd/del.go`

```go
// DelMultipleResult contains results of deleting multiple todos
type DelMultipleResult struct {
	DeletedTodos []todotxtlib.Todo
	LineNumbers  []int
	Count        int
}

// DelMultiple deletes multiple todos by line numbers
func DelMultiple(repo todotxtlib.TodoRepository, lineNumbers []int) (DelMultipleResult, error) {
	var deleted []todotxtlib.Todo

	// Sort line numbers descending to avoid index shifting
	sort.Sort(sort.Reverse(sort.IntSlice(lineNumbers)))

	for _, lineNum := range lineNumbers {
		index, err := repo.FindIndexByLineNumber(lineNum)
		if err != nil {
			return DelMultipleResult{}, fmt.Errorf("line %d: %w", lineNum, err)
		}

		todo, err := repo.Remove(index)
		if err != nil {
			return DelMultipleResult{}, fmt.Errorf("line %d: %w", lineNum, err)
		}

		deleted = append(deleted, todo)
	}

	if err := repo.Save(); err != nil {
		return DelMultipleResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelMultipleResult{
		DeletedTodos: deleted,
		LineNumbers:  lineNumbers,
		Count:        len(deleted),
	}, nil
}
```

**Step 4: Add unit test**

Modify: `cmd/del_test.go`

```go
func TestDelMultiple(t *testing.T) {
	repo := NewTestRepository(t)

	// Add test todos
	repo.Add("task 1")
	repo.Add("task 2")
	repo.Add("task 3")

	// Delete multiple (1 and 3)
	result, err := DelMultiple(repo, []int{1, 3})
	if err != nil {
		t.Fatalf("DelMultiple failed: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("Expected 2 deletions, got %d", result.Count)
	}

	// Verify only task 2 remains
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo remaining, got %d", len(todos))
	}

	if todos[0].Text != "task 2" {
		t.Errorf("Expected 'task 2', got '%s'", todos[0].Text)
	}
}
```

**Step 5: Run unit test**

Run: `go test ./cmd -v -run TestDelMultiple`
Expected: PASS

**Step 6: Update CLI wrapper to accept multiple args**

Modify: `internal/cli/del.go`

Change `Args: cobra.ExactArgs(1)` to `Args: cobra.MinimumNArgs(1)`

Parse all args as line numbers, call DelMultiple if more than one.

**Step 7: Run integration test**

Run: `go test -v ./tests -run TestDelMultiple`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/del.go cmd/del_test.go internal/cli/del.go
git commit -m "feat: support deleting multiple items at once

- Add DelMultiple() for batch deletion
- Update CLI to accept multiple line numbers
- Sort descending to avoid index issues
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 8: Implement Delete by Term

**Files:**
- Test: `tests/t1800_del_test.go` (TestDelWithTerm)
- Modify: `cmd/del.go`
- Modify: `cmd/del_test.go`
- Modify: `internal/cli/del.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestDelWithTerm`
Expected: SKIP or FAIL

**Step 2: Read test to understand syntax**

Use Read tool on test file.

Expected syntax: `togodo del "term"` (deletes all tasks containing term)

**Step 3: Add delete by term business logic**

Modify: `cmd/del.go`

```go
// DelByTermResult contains results of deleting by search term
type DelByTermResult struct {
	DeletedTodos []todotxtlib.Todo
	Term         string
	Count        int
}

// DelByTerm deletes all todos matching a search term
func DelByTerm(repo todotxtlib.TodoRepository, term string) (DelByTermResult, error) {
	allTodos, err := repo.ListAll()
	if err != nil {
		return DelByTermResult{}, err
	}

	var deleted []todotxtlib.Todo
	var indices []int

	// Find matching todos
	for i, todo := range allTodos {
		if strings.Contains(strings.ToLower(todo.Text), strings.ToLower(term)) {
			indices = append(indices, i)
			deleted = append(deleted, todo)
		}
	}

	if len(indices) == 0 {
		return DelByTermResult{
			Term:  term,
			Count: 0,
		}, fmt.Errorf("no tasks found matching '%s'", term)
	}

	// Delete in reverse order
	for i := len(indices) - 1; i >= 0; i-- {
		repo.Remove(indices[i])
	}

	if err := repo.Save(); err != nil {
		return DelByTermResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelByTermResult{
		DeletedTodos: deleted,
		Term:         term,
		Count:        len(deleted),
	}, nil
}
```

**Step 4: Add unit test**

Modify: `cmd/del_test.go`

```go
func TestDelByTerm(t *testing.T) {
	repo := NewTestRepository(t)

	repo.Add("task with roses")
	repo.Add("task with daisies")
	repo.Add("another roses task")

	result, err := DelByTerm(repo, "roses")
	if err != nil {
		t.Fatalf("DelByTerm failed: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("Expected 2 deletions, got %d", result.Count)
	}

	// Verify only daisies task remains
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
}
```

**Step 5: Run unit test**

Run: `go test ./cmd -v -run TestDelByTerm`
Expected: PASS

**Step 6: Update CLI wrapper to detect term vs line number**

Modify: `internal/cli/del.go`

If arg is not a number, treat as term search.

**Step 7: Run integration test**

Run: `go test -v ./tests -run TestDelWithTerm`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/del.go cmd/del_test.go internal/cli/del.go
git commit -m "feat: support deleting by search term

- Add DelByTerm() for term-based deletion
- Case-insensitive matching
- Delete all matching tasks
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 4: Advanced Do Features

**Goal:** Support multiple items and flags for do command

---

### Task 9: Implement Multiple Item Do with Comma Separation

**Files:**
- Test: `tests/t1500_do_test.go` (TestDoMultipleWithComma)
- Modify: `cmd/do.go`
- Modify: `cmd/do_test.go`
- Modify: `internal/cli/do.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestDoMultipleWithComma`
Expected: SKIP or FAIL

**Step 2: Read test to understand syntax**

Expected syntax: `togodo do 1,3,5` (comma-separated line numbers)

**Step 3: Update CLI wrapper to parse comma-separated list**

Modify: `internal/cli/do.go`

Split args by comma, parse each as line number.

**Step 4: Business logic already supports multiple items**

Verify `cmd.Do(repo, []int{...})` already accepts array.

**Step 5: Add unit test if needed**

Modify: `cmd/do_test.go` if additional testing needed.

**Step 6: Run integration test**

Run: `go test -v ./tests -run TestDoMultipleWithComma`
Expected: PASS

**Step 7: Commit**

```bash
git add internal/cli/do.go
git commit -m "feat: support comma-separated line numbers in do command

- Parse '1,3,5' syntax
- Call existing Do() with array
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 10: Implement Do Command Flags

**Files:**
- Test: `tests/t1500_do_test.go` (TestDoWithFlags)
- Modify: `internal/cli/do.go`

**Step 1: Run test to see what flags are needed**

Run: `go test -v ./tests -run TestDoWithFlags`
Expected: SKIP or FAIL

**Step 2: Read test to understand flags**

Use Read tool on test file.

Common flags: `--auto-archive`, `--no-date`

**Step 3: Add flags to CLI command**

Modify: `internal/cli/do.go`

Add flags:
```go
doCmd.Flags().Bool("auto-archive", false, "Automatically archive completed tasks")
doCmd.Flags().Bool("no-date", false, "Don't add completion date")
```

**Step 4: Update command logic to use flags**

Check flag values, modify behavior accordingly.

**Step 5: Run integration test**

Run: `go test -v ./tests -run TestDoWithFlags`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/cli/do.go
git commit -m "feat: add flags to do command

- Add --auto-archive flag
- Add --no-date flag
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 5: Move Command

**Goal:** Implement move command for inter-file task movement

---

### Task 11: Implement Move Command

**Files:**
- Test: `tests/t1850_move_test.go` (TestBasicMove, TestMoveUsage)
- Create: `cmd/move.go`
- Create: `cmd/move_test.go`
- Create: `internal/cli/move.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestBasicMove`
Expected: SKIP or FAIL

**Step 2: Read test to understand move semantics**

Use Read tool on test file.

Expected: Move task from todo.txt to another file (or vice versa).

**Step 3: Create move business logic**

Create: `cmd/move.go`

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// MoveResult contains the result of moving a todo
type MoveResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
	SourceFile string
	DestFile   string
}

// Move moves a todo from one file to another
func Move(sourceRepo, destRepo todotxtlib.TodoRepository, lineNumber int) (MoveResult, error) {
	// Find index by line number
	index, err := sourceRepo.FindIndexByLineNumber(lineNumber)
	if err != nil {
		return MoveResult{}, fmt.Errorf("failed to find task: %w", err)
	}

	// Get the todo
	todo, err := sourceRepo.Get(index)
	if err != nil {
		return MoveResult{}, fmt.Errorf("failed to get task: %w", err)
	}

	// Add to destination
	if _, err := destRepo.Add(todo.Text); err != nil {
		return MoveResult{}, fmt.Errorf("failed to add to destination: %w", err)
	}

	// Remove from source
	if _, err := sourceRepo.Remove(index); err != nil {
		return MoveResult{}, fmt.Errorf("failed to remove from source: %w", err)
	}

	// Save both
	if err := sourceRepo.Save(); err != nil {
		return MoveResult{}, fmt.Errorf("failed to save source: %w", err)
	}

	if err := destRepo.Save(); err != nil {
		return MoveResult{}, fmt.Errorf("failed to save destination: %w", err)
	}

	return MoveResult{
		Todo:       todo,
		LineNumber: lineNumber,
	}, nil
}
```

**Step 4: Create unit test**

Create: `cmd/move_test.go`

```go
package cmd

import (
	"testing"
)

func TestMove(t *testing.T) {
	sourceRepo := NewTestRepository(t)
	destRepo := NewTestRepository(t)

	sourceRepo.Add("task to move")
	sourceRepo.Add("task to keep")

	result, err := Move(sourceRepo, destRepo, 1)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	if result.Todo.Text != "task to move" {
		t.Errorf("Expected 'task to move', got '%s'", result.Todo.Text)
	}

	// Verify removed from source
	sourceTodos, _ := sourceRepo.ListAll()
	if len(sourceTodos) != 1 {
		t.Errorf("Expected 1 task in source, got %d", len(sourceTodos))
	}

	// Verify added to destination
	destTodos, _ := destRepo.ListAll()
	if len(destTodos) != 1 {
		t.Errorf("Expected 1 task in destination, got %d", len(destTodos))
	}
}
```

**Step 5: Run unit test**

Run: `go test ./cmd -v -run TestMove`
Expected: PASS

**Step 6: Create CLI wrapper**

Create: `internal/cli/move.go`

```go
package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewMoveCmd creates a Cobra command for moving todos
func NewMoveCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "move ITEM# DEST_FILE",
		Short: "Move a todo item to another file",
		Long: `Moves a todo item from current file to another file.

# move task 1 to done.txt
togodo move 1 done.txt
`,
		Args: cobra.ExactArgs(2),
		Aliases: []string{"mv"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			destFile := args[1]

			// Create destination repository
			// TODO: Create repo from dest file path
			destRepo, err := todotxtlib.NewRepository(destFile)
			if err != nil {
				return fmt.Errorf("failed to open destination: %w", err)
			}

			// Call business logic
			result, err := cmd.Move(repo, destRepo, lineNum)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: Moved to %s\n", destFile)
			return nil
		},
	}
}
```

**Step 7: Register command**

Modify: `internal/cli/root.go`

Add: `rootCmd.AddCommand(NewMoveCmd(repo))`

**Step 8: Run integration test**

Run: `go test -v ./tests -run TestBasicMove`
Expected: PASS

**Step 9: Commit**

```bash
git add cmd/move.go cmd/move_test.go internal/cli/move.go internal/cli/root.go
git commit -m "feat: implement move command

- Add Move() business logic
- Support moving tasks between files
- Add unit and integration tests
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 6: UI Enhancements

**Goal:** Add highlighting, plain mode, and other UI features

---

### Task 12: Implement Listall with Highlighting

**Files:**
- Test: `tests/t1350_listall_test.go` (TestListallHighlighting)
- Modify: `internal/cli/listall.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestListallHighlighting`
Expected: SKIP or FAIL

**Step 2: Read test to understand highlighting requirements**

Use Read tool on test file.

Expected: Color-coding for priorities, contexts, projects.

**Step 3: Add color support to listall**

Modify: `internal/cli/listall.go`

Use lipgloss for coloring if it's available in the project.

**Step 4: Run integration test**

Run: `go test -v ./tests -run TestListallHighlighting`
Expected: PASS

**Step 5: Commit if changes needed**

```bash
git add internal/cli/listall.go
git commit -m "feat: add highlighting to listall output

- Color-code priorities
- Highlight contexts and projects
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 13: Implement Priority Plain Flag

**Files:**
- Test: `tests/t1200_pri_test.go` (TestPriorityWithPlainFlag)
- Modify: `internal/cli/pri.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestPriorityWithPlainFlag`
Expected: SKIP or FAIL

**Step 2: Add --plain flag to pri command**

Modify: `internal/cli/pri.go`

Add: `priCmd.Flags().Bool("plain", false, "Plain output without colors")`

**Step 3: Update output logic to respect flag**

Check flag and disable colors if set.

**Step 4: Run integration test**

Run: `go test -v ./tests -run TestPriorityWithPlainFlag`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/pri.go
git commit -m "feat: add --plain flag to pri command

- Disable color output when flag set
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 14: Implement Short Help

**Files:**
- Test: `tests/t2100_help_test.go` (TestShortHelp)
- Modify: `internal/cli/root.go` or help command

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestShortHelp`
Expected: SKIP or FAIL

**Step 2: Read test to understand format**

Use Read tool on test file.

Expected: `togodo -h` shows short help (just command names).

**Step 3: Customize help output**

Modify help template to show short format for `-h` flag.

**Step 4: Run integration test**

Run: `go test -v ./tests -run TestShortHelp`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/root.go
git commit -m "feat: implement short help output

- Show concise help for -h flag
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 7: Multiline Support

**Goal:** Support multi-line task text

---

### Task 15: Implement Multiline Task Support

**Files:**
- Test: `tests/t2000_multiline_test.go` (TestMultilineAdd, TestMultilineHandling)
- Modify: `todotxtlib/parser.go`
- Modify: `cmd/add.go`
- Modify: `internal/cli/list.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestMultilineAdd`
Expected: SKIP or FAIL

**Step 2: Read test to understand format**

Use Read tool on test file.

Expected format: Tasks with embedded newlines or special markers.

**Step 3: Update parser to handle multiline**

Modify: `todotxtlib/parser.go`

Parse multiline task format (may use escaped newlines `\n` or other markers).

**Step 4: Update add command to accept multiline input**

Modify: `cmd/add.go`

Support creating multiline tasks.

**Step 5: Update list to display multiline correctly**

Modify: `internal/cli/list.go`

Format multiline tasks for display.

**Step 6: Run integration test**

Run: `go test -v ./tests -run TestMultilineAdd`
Expected: PASS

**Step 7: Commit**

```bash
git add todotxtlib/parser.go cmd/add.go internal/cli/list.go
git commit -m "feat: implement multiline task support

- Parse multiline task format
- Support creating multiline tasks
- Display multiline tasks correctly
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task Group 8: Additional Enhancements

**Goal:** Complete remaining test coverage

---

### Task 16: Implement Deduplicate with Priority

**Files:**
- Test: `tests/t1910_deduplicate_test.go` (TestDeduplicateWithPriority)
- Modify: `cmd/deduplicate.go`
- Modify: `cmd/deduplicate_test.go`

**Step 1: Run test to see what's expected**

Run: `go test -v ./tests -run TestDeduplicateWithPriority`
Expected: SKIP or FAIL

**Step 2: Read test to understand requirements**

Tasks with same text but different priorities should be treated as unique.

**Step 3: Update deduplicate logic**

Modify: `cmd/deduplicate.go`

Compare both text and priority when detecting duplicates.

**Step 4: Add unit test**

Modify: `cmd/deduplicate_test.go`

**Step 5: Run tests**

Run: `go test ./cmd -v -run TestDeduplicate`
Run: `go test -v ./tests -run TestDeduplicateWithPriority`
Expected: PASS

**Step 6: Commit**

```bash
git add cmd/deduplicate.go cmd/deduplicate_test.go
git commit -m "feat: treat priority as part of uniqueness in deduplicate

- Same text with different priorities are unique
- Integration test passing

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 17: Run Full Test Suite and Document Results

**Files:**
- Update: `docs/test-execution-summary.md`
- Update: `docs/implementation-complete-summary.md`

**Step 1: Run full test suite**

Run: `go test ./... -v 2>&1 | tee phase2-test-results.txt`

**Step 2: Count results**

Run: `grep -c "PASS" phase2-test-results.txt`
Run: `grep -c "SKIP" phase2-test-results.txt`
Run: `grep -c "FAIL" phase2-test-results.txt`

**Step 3: Update documentation**

Use Edit tool on documentation files to reflect Phase 2 completion.

**Step 4: Commit documentation**

```bash
git add docs/
git commit -m "docs: update for Phase 2 completion

- All Phase 2 tests passing
- Full todo.txt-cli parity achieved
- Updated test execution summary

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Success Criteria

- **All 56 tests passing** (40 Phase 1 + 16 Phase 2)
- 0% skip rate
- Full todo.txt-cli parity
- Each feature has:
  - Unit tests in `cmd/`
  - CLI wrapper in `internal/cli/`
  - Integration test passing
- Frequent commits (one per feature)
- All tests pass with `go test ./...`

## Notes

- Follow Phase 1 patterns (service layer, repository pattern, TDD)
- **Prefer Serena's semantic tools** for code navigation (`find_symbol`, `get_symbols_overview`, `find_referencing_symbols`) over basic text search
- Use `cmd.OutOrStdout()` for all output
- Match exact todo.txt-cli format and behavior
- Keep commits small and focused
- Test each feature thoroughly before moving to next

---

## Estimated Effort

**Total:** 50-70 hours

**Breakdown:**
- Config management: 4-6 hours
- Date support: 12-16 hours
- Advanced delete/do: 8-12 hours
- Move command: 6-8 hours
- UI enhancements: 8-12 hours
- Multiline support: 8-12 hours
- Documentation: 4-6 hours
