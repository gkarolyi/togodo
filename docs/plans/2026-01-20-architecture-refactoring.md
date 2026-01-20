# Architecture Refactoring Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refactor togodo into a clean three-layer architecture (todotxtlib data layer, cmd business logic, CLI/TUI presentation) to enable feature parity with todo.txt-cli.

**Architecture:** Simplify todotxtlib to a 12-method repository interface, remove the service layer, move business logic to pure functions in cmd package, separate Cobra wrappers to internal/cli, and update TUI to use cmd functions.

**Tech Stack:** Go 1.24.4, Cobra (CLI), Bubble Tea (TUI), existing todotxtlib

---

## Phase 1: Simplify todotxtlib (Foundation)

### Task 1.1: Add Get() method to TodoRepository

**Files:**
- Modify: `todotxtlib/repository.go`

**Step 1: Add Get() to interface**

In `todotxtlib/repository.go`, add the Get method to the TodoRepository interface after the Add method:

```go
type TodoRepository interface {
	Add(todoText string) (Todo, error)
	Get(index int) (Todo, error)  // ADD THIS LINE
	Remove(index int) (Todo, error)
	// ... rest of interface
}
```

**Step 2: Implement Get() in FileRepository**

In `todotxtlib/repository.go`, add the Get implementation after the Add method implementation:

```go
// Get returns a todo at the given index
func (r *FileRepository) Get(index int) (Todo, error) {
	if index < 0 || index >= len(r.todos) {
		return Todo{}, fmt.Errorf("index out of bounds")
	}
	return r.todos[index], nil
}
```

**Step 3: Test Get() method**

Run existing tests to ensure nothing breaks:

```bash
go test ./todotxtlib -v
```

Expected: All existing tests pass

**Step 4: Commit**

```bash
git add todotxtlib/repository.go
git commit -m "feat(todotxtlib): add Get() method to repository

Adds Get(index) method for retrieving individual todos by index.
Required for text modification commands (append, prepend, replace).

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 1.2: Update Sort() to accept nil for default

**Files:**
- Modify: `todotxtlib/repository.go`

**Step 1: Update Sort() implementation**

In `todotxtlib/repository.go`, modify the Sort method to handle nil:

```go
// Sort sorts the todos in the repository according to the specified criteria
// Pass nil to use default sort
func (r *FileRepository) Sort(sort Sort) {
	if sort == nil {
		sort = NewDefaultSort()
	}
	sort.Apply(r.todos)
}
```

**Step 2: Remove SortDefault() from interface**

In `todotxtlib/repository.go`, remove the SortDefault() line from the TodoRepository interface:

```go
type TodoRepository interface {
	// ... other methods
	Sort(sort Sort)
	// REMOVE: SortDefault()
	Save() error
	// ... rest
}
```

**Step 3: Remove SortDefault() implementation**

In `todotxtlib/repository.go`, delete the entire SortDefault method:

```go
// DELETE THIS ENTIRE METHOD:
// func (r *FileRepository) SortDefault() {
// 	sort := NewDefaultSort()
// 	sort.Apply(r.todos)
// }
```

**Step 4: Update service.go calls temporarily**

In `todotxtlib/service.go`, replace all `s.repo.SortDefault()` calls with `s.repo.Sort(nil)`:

```bash
# Find all occurrences (should be 3)
grep -n "SortDefault" todotxtlib/service.go
```

Replace each occurrence manually or use sed:

```bash
sed -i 's/s\.repo\.SortDefault()/s.repo.Sort(nil)/g' todotxtlib/service.go
```

**Step 5: Test Sort changes**

```bash
go test ./todotxtlib -v
```

Expected: All tests pass

**Step 6: Commit**

```bash
git add todotxtlib/repository.go todotxtlib/service.go
git commit -m "refactor(todotxtlib): consolidate Sort methods

Replace SortDefault() with Sort(nil) for idiomatic Go.
nil parameter triggers default sort behavior.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 1.3: Remove Search() method from interface

**Files:**
- Modify: `todotxtlib/repository.go`

**Step 1: Remove Search() from interface**

In `todotxtlib/repository.go`, remove the Search line from the TodoRepository interface:

```go
type TodoRepository interface {
	// ... other methods
	Filter(filter Filter) ([]Todo, error)
	// REMOVE: Search(query string) ([]Todo, error)
	Sort(sort Sort)
	// ... rest
}
```

**Step 2: Keep Search() implementation for now**

Keep the Search() implementation in FileRepository - we'll update callers first, then remove it.

**Step 3: Update service.go SearchTodos**

In `todotxtlib/service.go`, update the SearchTodos method to use Filter:

```go
// SearchTodos searches for todos matching the given query
// Returns matching todos
func (s *DefaultTodoService) SearchTodos(query string) ([]Todo, error) {
	if query == "" {
		return s.repo.ListAll()
	}
	filter := Filter{Text: query}
	return s.repo.Filter(filter)
}
```

**Step 4: Update TUI to use Filter**

In `internal/tui/update.go`, replace the Search call (around line 81 and 93):

Find:
```go
filteredTodos, _ := m.repository.Search(m.filter)
```

Replace with:
```go
filter := todotxtlib.Filter{Text: m.filter}
filteredTodos, _ := m.repository.Filter(filter)
```

**Step 5: Test changes**

```bash
go test ./todotxtlib -v
go test ./internal/tui -v
```

Expected: All tests pass

**Step 6: Commit**

```bash
git add todotxtlib/repository.go todotxtlib/service.go internal/tui/update.go
git commit -m "refactor(todotxtlib): remove Search() from interface

Use Filter(Filter{Text: query}) instead.
Updates service and TUI to use Filter.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 1.4: Remove WriteToString() method

**Files:**
- Modify: `todotxtlib/repository.go`

**Step 1: Check for WriteToString usage**

```bash
grep -r "WriteToString" . --include="*.go" --exclude-dir=vendor
```

Expected: Only used in tests or test helpers

**Step 2: Remove WriteToString() from interface**

In `todotxtlib/repository.go`, remove the WriteToString line from the TodoRepository interface:

```go
type TodoRepository interface {
	// ... other methods
	Save() error
	// REMOVE: WriteToString() (string, error)
}
```

**Step 3: Remove WriteToString() implementation**

In `todotxtlib/repository.go`, delete the entire WriteToString method:

```go
// DELETE THIS ENTIRE METHOD:
// func (r *FileRepository) WriteToString() (string, error) {
// 	var buffer bytes.Buffer
// 	writer := NewBufferWriter(&buffer)
// 	err := writer.Write(r.todos)
// 	if err != nil {
// 		return "", err
// 	}
// 	return buffer.String(), nil
// }
```

**Step 4: Update test helpers if needed**

Check `tests/test_helpers.go` - if it uses WriteToString, update to use Save with a buffer:

```go
// Instead of:
// content, err := repo.WriteToString()

// Use:
var buf bytes.Buffer
writer := todotxtlib.NewBufferWriter(&buf)
tempRepo := todotxtlib.NewFileRepository(reader, writer)
err := tempRepo.Save()
content := buf.String()
```

**Step 5: Test changes**

```bash
go test ./todotxtlib -v
go test ./tests -v
```

Expected: All tests pass

**Step 6: Commit**

```bash
git add todotxtlib/repository.go tests/test_helpers.go
git commit -m "refactor(todotxtlib): remove WriteToString() method

Use Save() with BufferWriter instead for consistency.
Tests updated to use writer abstraction.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 1.5: Remove micro-methods from interface

**Files:**
- Modify: `todotxtlib/repository.go`

**Step 1: Check usage of micro-methods**

```bash
grep -r "SetContexts\|SetProjects\|AddContext\|RemoveContext\|AddProject\|RemoveProject" . --include="*.go" --exclude-dir=vendor
```

Expected: Only in internal/tui if anywhere

**Step 2: Remove methods from interface**

In `todotxtlib/repository.go`, remove these lines from TodoRepository interface:

```go
type TodoRepository interface {
	// ... other methods
	// REMOVE these 6 lines:
	// SetContexts(index int, contexts []string) (Todo, error)
	// SetProjects(index int, projects []string) (Todo, error)
	// AddContext(index int, context string) (Todo, error)
	// AddProject(index int, project string) (Todo, error)
	// RemoveContext(index int, context string) (Todo, error)
	// RemoveProject(index int, project string) (Todo, error)
	Filter(filter Filter) ([]Todo, error)
	// ... rest
}
```

**Step 3: Remove implementations from FileRepository**

In `todotxtlib/repository.go`, delete these 6 methods entirely:
- SetContexts
- SetProjects
- AddContext
- AddProject
- RemoveContext
- RemoveProject

**Step 4: Remove ListTodos() and ListDone() from interface**

In `todotxtlib/repository.go`, remove these lines from the interface:

```go
type TodoRepository interface {
	// ... other methods
	ListAll() ([]Todo, error)
	// REMOVE: ListTodos() ([]Todo, error)
	// REMOVE: ListDone() ([]Todo, error)
	ListProjects() ([]string, error)
	// ... rest
}
```

**Step 5: Remove ListTodos() and ListDone() implementations**

Delete the ListTodos and ListDone methods from FileRepository in `todotxtlib/repository.go`.

**Step 6: Update service.go if it uses these methods**

Check `todotxtlib/service.go` for ListDone usage (in RemoveDoneTodos):

```go
// RemoveDoneTodos - replace ListDone() call with Filter
func (s *DefaultTodoService) RemoveDoneTodos() ([]Todo, error) {
	// Get done todos before removing
	doneFilter := Filter{Done: true}
	doneTodos, err := s.repo.Filter(doneFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list done todos: %w", err)
	}

	// Rest of implementation stays the same
	// ...
}
```

**Step 7: Test changes**

```bash
go test ./todotxtlib -v
```

Expected: All tests pass

**Step 8: Commit**

```bash
git add todotxtlib/repository.go todotxtlib/service.go
git commit -m "refactor(todotxtlib): remove micro-methods from interface

Removes: SetContexts, SetProjects, AddContext, RemoveContext,
AddProject, RemoveProject, ListTodos, ListDone.

Use Get() + Update() for modifications, Filter() for queries.
Reduces interface from 20+ to 12 methods.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 1.6: Remove service layer

**Files:**
- Delete: `todotxtlib/service.go`
- Modify: `cmd/root.go`
- Modify: `cmd/*.go` (all command files)

**Step 1: Identify service usage**

```bash
grep -r "TodoService\|NewTodoService" . --include="*.go" --exclude-dir=vendor
```

Expected: Used in cmd/root.go and passed to commands

**Step 2: Update cmd/root.go to remove service**

In `cmd/root.go`, find the Execute function and remove service creation:

Before:
```go
service := todotxtlib.NewTodoService(repository)
rootCmd.AddCommand(NewAddCmd(service, presenter))
```

After:
```go
rootCmd.AddCommand(NewAddCmd(repository, presenter))
```

Update all command constructor calls to pass repository instead of service.

**Step 3: Update command constructors**

Update each command file to accept TodoRepository instead of TodoService:

In `cmd/add.go`:
```go
// Before:
func NewAddCmd(service todotxtlib.TodoService, presenter *cli.Presenter) *cobra.Command {

// After:
func NewAddCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
```

Update the RunE function to call repository directly with orchestration logic.

**Step 4: Update add command logic**

In `cmd/add.go`, replace service call with direct repository calls + orchestration:

```go
RunE: func(cmd *cobra.Command, args []string) error {
	// Add todos
	addedTodos := make([]todotxtlib.Todo, 0, len(args))
	for _, text := range args {
		todo, err := repo.Add(text)
		if err != nil {
			return fmt.Errorf("failed to add todo: %w", err)
		}
		addedTodos = append(addedTodos, todo)
	}

	// Sort and save
	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// Present
	for _, todo := range addedTodos {
		presenter.Print(todo)
	}
	return nil
},
```

**Step 5: Update remaining commands similarly**

Update `cmd/do.go`, `cmd/pri.go`, `cmd/tidy.go`, `cmd/list.go` to:
- Accept TodoRepository instead of TodoService
- Implement orchestration logic directly

**Step 6: Delete service.go**

```bash
rm todotxtlib/service.go
```

**Step 7: Test everything compiles**

```bash
go build
```

Expected: Successful compilation

**Step 8: Run tests**

```bash
go test ./...
```

Expected: Existing tests may fail (expected - we're changing behavior)

**Step 9: Commit**

```bash
git add todotxtlib/service.go cmd/
git commit -m "refactor: remove service layer

Move orchestration logic directly into cmd package.
Commands now call repository directly.
Service layer provided minimal value and added indirection.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 2: Create cmd Business Logic Functions

### Task 2.1: Create Add business logic function

**Files:**
- Create: `cmd/add_logic.go`
- Create: `cmd/add_logic_test.go`

**Step 1: Write failing test**

Create `cmd/add_logic_test.go`:

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestAdd(t *testing.T) {
	t.Run("adds single task", func(t *testing.T) {
		// Setup
		var buf strings.Builder
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := Add(repo, []string{"test", "task"})

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Todo.Text != "test task" {
			t.Errorf("expected 'test task', got '%s'", result.Todo.Text)
		}
		if result.LineNumber != 1 {
			t.Errorf("expected line number 1, got %d", result.LineNumber)
		}
	})
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./cmd -run TestAdd -v
```

Expected: FAIL with "undefined: Add" or "undefined: AddResult"

**Step 3: Create Add function**

Create `cmd/add_logic.go`:

```go
package cmd

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AddResult contains the result of an Add operation
type AddResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Add adds a new todo task to the repository
func Add(repo todotxtlib.TodoRepository, args []string) (AddResult, error) {
	if len(args) == 0 {
		return AddResult{}, fmt.Errorf("task text required")
	}

	// Join args into single task
	text := strings.Join(args, " ")

	// Add to repository
	todo, err := repo.Add(text)
	if err != nil {
		return AddResult{}, fmt.Errorf("failed to add todo: %w", err)
	}

	// Sort and save
	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return AddResult{}, fmt.Errorf("failed to save: %w", err)
	}

	// Find line number after sort
	allTodos, err := repo.ListAll()
	if err != nil {
		return AddResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	lineNumber := 1
	for i, t := range allTodos {
		if t.Text == todo.Text && t.Priority == todo.Priority {
			lineNumber = i + 1
			break
		}
	}

	return AddResult{Todo: todo, LineNumber: lineNumber}, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./cmd -run TestAdd -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add cmd/add_logic.go cmd/add_logic_test.go
git commit -m "feat(cmd): add Add business logic function

Pure function for adding todos with orchestration logic.
Returns AddResult with todo and line number.
Includes unit tests.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2.2: Create List business logic function

**Files:**
- Create: `cmd/list_logic.go`
- Create: `cmd/list_logic_test.go`

**Step 1: Write failing test**

Create `cmd/list_logic_test.go`:

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestList(t *testing.T) {
	t.Run("lists all tasks", func(t *testing.T) {
		// Setup
		buf := strings.Builder{}
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := List(repo, "")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Todos) != 2 {
			t.Errorf("expected 2 todos, got %d", len(result.Todos))
		}
		if result.TotalCount != 2 {
			t.Errorf("expected total count 2, got %d", result.TotalCount)
		}
		if result.ShownCount != 2 {
			t.Errorf("expected shown count 2, got %d", result.ShownCount)
		}
	})

	t.Run("filters tasks", func(t *testing.T) {
		// Setup
		buf := strings.Builder{}
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := List(repo, "one")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(result.Todos))
		}
		if result.TotalCount != 2 {
			t.Errorf("expected total count 2, got %d", result.TotalCount)
		}
		if result.ShownCount != 1 {
			t.Errorf("expected shown count 1, got %d", result.ShownCount)
		}
	})
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./cmd -run TestList -v
```

Expected: FAIL with "undefined: List"

**Step 3: Create List function**

Create `cmd/list_logic.go`:

```go
package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListResult contains the result of a List operation
type ListResult struct {
	Todos      []todotxtlib.Todo
	TotalCount int
	ShownCount int
}

// List lists todos, optionally filtering by search query
func List(repo todotxtlib.TodoRepository, searchQuery string) (ListResult, error) {
	// Get total count
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListResult{}, err
	}
	totalCount := len(allTodos)

	// Filter if search query provided
	var todos []todotxtlib.Todo
	if searchQuery != "" {
		filter := todotxtlib.Filter{Text: searchQuery}
		todos, err = repo.Filter(filter)
		if err != nil {
			return ListResult{}, err
		}
	} else {
		todos = allTodos
	}

	return ListResult{
		Todos:      todos,
		TotalCount: totalCount,
		ShownCount: len(todos),
	}, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./cmd -run TestList -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add cmd/list_logic.go cmd/list_logic_test.go
git commit -m "feat(cmd): add List business logic function

Lists todos with optional filtering.
Returns ListResult with todos and counts.
Includes unit tests.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2.3: Create Do business logic function

**Files:**
- Create: `cmd/do_logic.go`
- Create: `cmd/do_logic_test.go`

**Step 1: Write failing test**

Create `cmd/do_logic_test.go`:

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestDo(t *testing.T) {
	t.Run("toggles task done", func(t *testing.T) {
		// Setup
		buf := strings.Builder{}
		buf.WriteString("task one\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := Do(repo, []int{0})

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.ToggledTodos) != 1 {
			t.Errorf("expected 1 toggled todo, got %d", len(result.ToggledTodos))
		}
		if !result.ToggledTodos[0].Done {
			t.Error("expected todo to be marked done")
		}
	})
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./cmd -run TestDo -v
```

Expected: FAIL with "undefined: Do"

**Step 3: Create Do function**

Create `cmd/do_logic.go`:

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DoResult contains the result of a Do operation
type DoResult struct {
	ToggledTodos []todotxtlib.Todo
}

// Do toggles the done status of todos at the given indices (0-based)
func Do(repo todotxtlib.TodoRepository, indices []int) (DoResult, error) {
	toggledTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.ToggleDone(index)
		if err != nil {
			return DoResult{}, fmt.Errorf("failed to toggle todo at index %d: %w", index, err)
		}
		toggledTodos = append(toggledTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return DoResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DoResult{ToggledTodos: toggledTodos}, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./cmd -run TestDo -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add cmd/do_logic.go cmd/do_logic_test.go
git commit -m "feat(cmd): add Do business logic function

Toggles done status for tasks.
Returns DoResult with toggled todos.
Includes unit tests.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2.4: Create Pri business logic function

**Files:**
- Create: `cmd/pri_logic.go`
- Create: `cmd/pri_logic_test.go`

**Step 1: Write failing test**

Create `cmd/pri_logic_test.go`:

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestSetPriority(t *testing.T) {
	t.Run("sets priority on task", func(t *testing.T) {
		// Setup
		buf := strings.Builder{}
		buf.WriteString("task one\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := SetPriority(repo, []int{0}, "A")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.UpdatedTodos) != 1 {
			t.Errorf("expected 1 updated todo, got %d", len(result.UpdatedTodos))
		}
		if result.UpdatedTodos[0].Priority != "A" {
			t.Errorf("expected priority A, got %s", result.UpdatedTodos[0].Priority)
		}
	})
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./cmd -run TestSetPriority -v
```

Expected: FAIL with "undefined: SetPriority"

**Step 3: Create SetPriority function**

Create `cmd/pri_logic.go`:

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// PriResult contains the result of a SetPriority operation
type PriResult struct {
	UpdatedTodos []todotxtlib.Todo
}

// SetPriority sets the priority for todos at the given indices (0-based)
func SetPriority(repo todotxtlib.TodoRepository, indices []int, priority string) (PriResult, error) {
	updatedTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.SetPriority(index, priority)
		if err != nil {
			return PriResult{}, fmt.Errorf("failed to set priority at index %d: %w", index, err)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := repo.Save(); err != nil {
		return PriResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return PriResult{UpdatedTodos: updatedTodos}, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./cmd -run TestSetPriority -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add cmd/pri_logic.go cmd/pri_logic_test.go
git commit -m "feat(cmd): add SetPriority business logic function

Sets priority on tasks (doesn't sort to preserve order).
Returns PriResult with updated todos.
Includes unit tests.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 2.5: Create Tidy business logic function

**Files:**
- Create: `cmd/tidy_logic.go`
- Create: `cmd/tidy_logic_test.go`

**Step 1: Write failing test**

Create `cmd/tidy_logic_test.go`:

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestTidy(t *testing.T) {
	t.Run("removes done tasks", func(t *testing.T) {
		// Setup
		buf := strings.Builder{}
		buf.WriteString("task one\nx done task\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := Tidy(repo)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.RemovedTodos) != 1 {
			t.Errorf("expected 1 removed todo, got %d", len(result.RemovedTodos))
		}
		if !result.RemovedTodos[0].Done {
			t.Error("expected removed todo to be done")
		}

		// Verify remaining todos
		allTodos, _ := repo.ListAll()
		if len(allTodos) != 1 {
			t.Errorf("expected 1 remaining todo, got %d", len(allTodos))
		}
	})
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./cmd -run TestTidy -v
```

Expected: FAIL with "undefined: Tidy"

**Step 3: Create Tidy function**

Create `cmd/tidy_logic.go`:

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// TidyResult contains the result of a Tidy operation
type TidyResult struct {
	RemovedTodos []todotxtlib.Todo
}

// Tidy removes all completed todos
func Tidy(repo todotxtlib.TodoRepository) (TidyResult, error) {
	// Get done todos before removing
	doneFilter := todotxtlib.Filter{Done: true}
	doneTodos, err := repo.Filter(doneFilter)
	if err != nil {
		return TidyResult{}, fmt.Errorf("failed to filter done todos: %w", err)
	}

	// Get all todos to iterate
	allTodos, err := repo.ListAll()
	if err != nil {
		return TidyResult{}, fmt.Errorf("failed to list all todos: %w", err)
	}

	// Remove backwards to avoid index shifting
	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			if _, err := repo.Remove(i); err != nil {
				return TidyResult{}, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
			}
		}
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return TidyResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return TidyResult{RemovedTodos: doneTodos}, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./cmd -run TestTidy -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add cmd/tidy_logic.go cmd/tidy_logic_test.go
git commit -m "feat(cmd): add Tidy business logic function

Removes completed todos from the list.
Returns TidyResult with removed todos.
Includes unit tests.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Phase 3: Migrate CLI to use cmd Functions

### Task 3.1: Create internal/cli directory structure

**Files:**
- Create: `internal/cli/add.go`

**Step 1: Create directory**

```bash
mkdir -p internal/cli
```

**Step 2: Create add command wrapper**

Create `internal/cli/add.go`:

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewAddCmd creates a Cobra command for adding todos
func NewAddCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "add [TASK]",
		Short: "Add a new todo item to the list",
		Long: `Adds a new task to the list and prints the newly added task.

# add "Buy milk" to the list
togodo add "Buy milk"

# add multiple words
togodo add Buy milk and eggs
`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"a"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Add(repo, args)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			fmt.Printf("%d %s\n", result.LineNumber, result.Todo.String())
			fmt.Printf("TODO: %d added.\n", result.LineNumber)
			return nil
		},
	}
}
```

**Step 3: Build to verify it compiles**

```bash
go build
```

Expected: Successful compilation

**Step 4: Commit**

```bash
git add internal/cli/add.go
git commit -m "feat(cli): add Add command wrapper

Cobra wrapper that calls cmd.Add and formats output.
Matches todo.txt-cli output format.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.2: Create list command wrapper

**Files:**
- Create: `internal/cli/list.go`

**Step 1: Create list command wrapper**

Create `internal/cli/list.go`:

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListCmd creates a Cobra command for listing todos
func NewListCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "list [FILTER]",
		Short: "List all todo items",
		Long: `Lists all todo items, optionally filtered by search term.

# list all tasks
togodo list

# list tasks containing "milk"
togodo list milk
`,
		Aliases: []string{"l", "ls"},
		RunE: func(command *cobra.Command, args []string) error {
			// Get search filter if provided
			searchQuery := ""
			if len(args) > 0 {
				searchQuery = args[0]
			}

			// Call business logic
			result, err := cmd.List(repo, searchQuery)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for i, todo := range result.Todos {
				fmt.Printf("%d %s\n", i+1, todo.String())
			}
			fmt.Println("--")
			fmt.Printf("TODO: %d of %d tasks shown\n", result.ShownCount, result.TotalCount)
			return nil
		},
	}
}
```

**Step 2: Build to verify it compiles**

```bash
go build
```

Expected: Successful compilation

**Step 3: Commit**

```bash
git add internal/cli/list.go
git commit -m "feat(cli): add List command wrapper

Cobra wrapper that calls cmd.List and formats output.
Matches todo.txt-cli output with separator and summary.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.3: Create do command wrapper

**Files:**
- Create: `internal/cli/do.go`

**Step 1: Create do command wrapper**

Create `internal/cli/do.go`:

```go
package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDoCmd creates a Cobra command for toggling todo completion
func NewDoCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "do [LINE_NUMBER]",
		Short: "Mark a todo item as done (or undone)",
		Long: `Toggles the completion status of a todo item.

# mark task 1 as done
togodo do 1
`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"x"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Call business logic
			result, err := cmd.Do(repo, []int{index})
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for _, todo := range result.ToggledTodos {
				fmt.Printf("%d %s\n", lineNum, todo.String())
				if todo.Done {
					fmt.Printf("TODO: %d marked as done.\n", lineNum)
				} else {
					fmt.Printf("TODO: %d marked as TODO.\n", lineNum)
				}
			}
			return nil
		},
	}
}
```

**Step 2: Build to verify it compiles**

```bash
go build
```

Expected: Successful compilation

**Step 3: Commit**

```bash
git add internal/cli/do.go
git commit -m "feat(cli): add Do command wrapper

Cobra wrapper that calls cmd.Do and formats output.
Handles 1-based line numbers from CLI.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.4: Create pri command wrapper

**Files:**
- Create: `internal/cli/pri.go`

**Step 1: Create pri command wrapper**

Create `internal/cli/pri.go`:

```go
package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewPriCmd creates a Cobra command for setting priority
func NewPriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "pri [LINE_NUMBER] [PRIORITY]",
		Short: "Set priority of a todo item",
		Long: `Sets the priority of a todo item (A, B, C, etc.).

# set task 1 to priority A
togodo pri 1 A
`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"p"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get priority (normalize to uppercase)
			priority := strings.ToUpper(args[1])

			// Call business logic
			result, err := cmd.SetPriority(repo, []int{index}, priority)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for _, todo := range result.UpdatedTodos {
				fmt.Printf("%d %s\n", lineNum, todo.String())
				fmt.Printf("TODO: %d prioritized to (%s).\n", lineNum, priority)
			}
			return nil
		},
	}
}
```

**Step 2: Build to verify it compiles**

```bash
go build
```

Expected: Successful compilation

**Step 3: Commit**

```bash
git add internal/cli/pri.go
git commit -m "feat(cli): add Pri command wrapper

Cobra wrapper that calls cmd.SetPriority and formats output.
Normalizes priority to uppercase.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.5: Create tidy command wrapper

**Files:**
- Create: `internal/cli/tidy.go`

**Step 1: Create tidy command wrapper**

Create `internal/cli/tidy.go`:

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewTidyCmd creates a Cobra command for removing done tasks
func NewTidyCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "Remove completed tasks",
		Long: `Removes all completed tasks from the todo list.

# remove all done tasks
togodo tidy
`,
		Aliases: []string{"clean"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Tidy(repo)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			if len(result.RemovedTodos) == 0 {
				fmt.Println("TODO: No completed tasks to remove.")
			} else {
				for _, todo := range result.RemovedTodos {
					fmt.Println(todo.String())
				}
				fmt.Printf("TODO: %d completed task(s) removed.\n", len(result.RemovedTodos))
			}
			return nil
		},
	}
}
```

**Step 2: Build to verify it compiles**

```bash
go build
```

Expected: Successful compilation

**Step 3: Commit**

```bash
git add internal/cli/tidy.go
git commit -m "feat(cli): add Tidy command wrapper

Cobra wrapper that calls cmd.Tidy and formats output.
Shows removed tasks and count.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.6: Update root.go to use new CLI wrappers

**Files:**
- Modify: `cmd/root.go`

**Step 1: Update imports in root.go**

In `cmd/root.go`, add import for internal/cli:

```go
import (
	// ... existing imports
	"github.com/gkarolyi/togodo/internal/cli"
)
```

**Step 2: Update command registration**

In `cmd/root.go`, find the Execute function and update to use new wrappers:

Before:
```go
rootCmd.AddCommand(NewAddCmd(repository, presenter))
rootCmd.AddCommand(NewListCmd(repository, presenter))
rootCmd.AddCommand(NewDoCmd(repository, presenter))
rootCmd.AddCommand(NewPriCmd(repository, presenter))
rootCmd.AddCommand(NewTidyCmd(repository, presenter))
```

After:
```go
rootCmd.AddCommand(cli.NewAddCmd(repository))
rootCmd.AddCommand(cli.NewListCmd(repository, presenter))
rootCmd.AddCommand(cli.NewDoCmd(repository))
rootCmd.AddCommand(cli.NewPriCmd(repository))
rootCmd.AddCommand(cli.NewTidyCmd(repository))
```

**Step 3: Remove old command files**

Once verified that everything works, remove the old command implementations:

```bash
# Keep these files for now, they have tests
# We'll update them in the next phase
# rm cmd/add.go cmd/list.go cmd/do.go cmd/pri.go cmd/tidy.go
```

Actually, let's keep them for now since they might have tests.

**Step 4: Build and test**

```bash
go build
./togodo add "test task"
./togodo list
```

Expected: Commands work with new output format

**Step 5: Commit**

```bash
git add cmd/root.go
git commit -m "refactor(cli): wire up new CLI command wrappers

Update root.go to use internal/cli wrappers.
Commands now use cmd business logic layer.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 3.7: Run integration tests to check progress

**Files:**
- None (just running tests)

**Step 1: Run integration tests**

```bash
go test ./tests -v
```

**Step 2: Review output**

Note which tests are now passing vs before. Expected improvements:
- Add command should now join args correctly
- List command should show separator and summary
- Output formatting should be closer to todo.txt-cli

**Step 3: Document findings**

Create a quick note of what's improved and what still needs work.

**Step 4: No commit needed**

This is just a checkpoint to see progress.

---

## Phase 4: Migrate TUI to use cmd Functions

### Task 4.1: Update TUI to call cmd.Add

**Files:**
- Modify: `internal/tui/update.go`

**Step 1: Add import**

In `internal/tui/update.go`, add import at the top:

```go
import (
	// ... existing imports
	"github.com/gkarolyi/togodo/cmd"
)
```

**Step 2: Update add logic**

Find the add logic (around line 50) and replace:

Before:
```go
if m.input.Value() != "" {
	m.repository.Add(m.input.Value())
	allTodos, _ := m.repository.ListAll()
	m.choices = allTodos
	m.adding = false
	m.input.Reset()
	m.input.Blur()
}
```

After:
```go
if m.input.Value() != "" {
	_, err := cmd.Add(m.repository, []string{m.input.Value()})
	if err != nil {
		// TODO: Show error in UI
		// For now, just continue
	}

	// Refresh display
	allTodos, _ := m.repository.ListAll()
	m.choices = allTodos
	m.adding = false
	m.input.Reset()
	m.input.Blur()
}
```

**Step 3: Test TUI manually**

```bash
./togodo
# Press 'a', type a task, press Enter
# Verify task appears in list
```

Expected: Add works in TUI

**Step 4: Commit**

```bash
git add internal/tui/update.go
git commit -m "refactor(tui): use cmd.Add for adding tasks

TUI now calls cmd.Add business logic.
Refreshes display with ListAll after add.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4.2: Update TUI to call cmd.Do

**Files:**
- Modify: `internal/tui/update.go`

**Step 1: Update toggle done logic**

Find the toggle done logic (around line 130) and replace:

Before:
```go
case "x":
	for i := range m.selected {
		m.repository.ToggleDone(i)
	}
	allTodos, _ := m.repository.ListAll()
	m.choices = allTodos
```

After:
```go
case "x":
	// Convert selected map to slice of indices
	indices := make([]int, 0, len(m.selected))
	for i := range m.selected {
		indices = append(indices, i)
	}

	if len(indices) > 0 {
		_, err := cmd.Do(m.repository, indices)
		if err != nil {
			// TODO: Show error in UI
		}

		// Refresh display
		allTodos, _ := m.repository.ListAll()
		m.choices = allTodos
		m.selected = make(map[int]struct{}) // Clear selection
	}
```

**Step 2: Test TUI manually**

```bash
./togodo
# Select a task with Space, press 'x'
# Verify task is marked done
```

Expected: Toggle done works in TUI

**Step 3: Commit**

```bash
git add internal/tui/update.go
git commit -m "refactor(tui): use cmd.Do for toggling tasks

TUI now calls cmd.Do business logic.
Clears selection after toggling.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4.3: Update TUI to call cmd.SetPriority

**Files:**
- Modify: `internal/tui/update.go`

**Step 1: Update set priority logic**

Find the set priority logic (around line 30) and replace:

Before:
```go
if m.setting {
	switch msg.String() {
	case "esc":
		m.setting = false
		return m, nil
	case "a", "b", "c", "d", "A", "B", "C", "D":
		priority := strings.ToUpper(msg.String())
		for i := range m.selected {
			m.repository.SetPriority(i, priority)
		}
		allTodos, _ := m.repository.ListAll()
		m.choices = allTodos
		m.setting = false
		return m, nil
	}
	return m, nil
}
```

After:
```go
if m.setting {
	switch msg.String() {
	case "esc":
		m.setting = false
		return m, nil
	case "a", "b", "c", "d", "A", "B", "C", "D":
		priority := strings.ToUpper(msg.String())

		// Convert selected map to slice of indices
		indices := make([]int, 0, len(m.selected))
		for i := range m.selected {
			indices = append(indices, i)
		}

		if len(indices) > 0 {
			_, err := cmd.SetPriority(m.repository, indices, priority)
			if err != nil {
				// TODO: Show error in UI
			}

			// Refresh display
			allTodos, _ := m.repository.ListAll()
			m.choices = allTodos
			m.selected = make(map[int]struct{}) // Clear selection
		}
		m.setting = false
		return m, nil
	}
	return m, nil
}
```

**Step 2: Test TUI manually**

```bash
./togodo
# Select a task with Space, press 'p', press 'A'
# Verify task gets priority A
```

Expected: Set priority works in TUI

**Step 3: Commit**

```bash
git add internal/tui/update.go
git commit -m "refactor(tui): use cmd.SetPriority for setting priority

TUI now calls cmd.SetPriority business logic.
Clears selection after setting priority.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 4.4: Update TUI filter to use Filter instead of Search

**Files:**
- Modify: `internal/tui/update.go`

**Step 1: Update filter logic**

This was already done in Task 1.3, but verify it's correct.

Find the filtering logic and ensure it uses Filter:

```go
filter := todotxtlib.Filter{Text: m.filter}
filteredTodos, _ := m.repository.Filter(filter)
```

**Step 2: Test TUI filtering**

```bash
./togodo
# Press '/', type some text
# Verify filtering works
```

Expected: Filtering works in TUI

**Step 3: No commit needed**

Already committed in Task 1.3.

---

## Phase 5: Clean Up and Final Verification

### Task 5.1: Remove old command implementations

**Files:**
- Modify: `cmd/add.go`, `cmd/list.go`, `cmd/do.go`, `cmd/pri.go`, `cmd/tidy.go`

**Step 1: Check for old NewXxxCmd functions**

These files might still have the old Cobra command constructors. We need to remove them since we're using the ones in internal/cli now.

**Step 2: Keep business logic, remove Cobra wrappers**

In each file (`cmd/add.go`, etc.), remove the old `NewXxxCmd` function if it exists, but keep any helper functions or types that might be used.

Actually, we created new files `cmd/xxx_logic.go` for the business logic, so the old files can be removed entirely if they only contain the old Cobra commands.

**Step 3: Check what's in each old file**

```bash
head -20 cmd/add.go
head -20 cmd/list.go
# etc.
```

**Step 4: Remove old files if they're obsolete**

If the old files only contain the old Cobra commands and nothing else useful:

```bash
# Only remove if they're truly obsolete
# Check each one first
# rm cmd/add.go cmd/list.go cmd/do.go cmd/pri.go cmd/tidy.go
```

**Step 5: Update tests if needed**

Check if `cmd/*_test.go` files need updating to import from the right place.

**Step 6: Commit**

```bash
git add cmd/
git commit -m "refactor(cmd): remove old command implementations

Old Cobra wrappers replaced by internal/cli.
Business logic now in xxx_logic.go files.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

### Task 5.2: Run full test suite

**Files:**
- None (just testing)

**Step 1: Run all tests**

```bash
go test ./... -v
```

**Step 2: Review results**

Check which tests pass/fail:
- todotxtlib tests should all pass
- cmd tests (new business logic tests) should pass
- integration tests should show improvement

**Step 3: Document results**

Note what's working and what still needs implementation.

**Step 4: No commit needed**

---

### Task 5.3: Test CLI manually

**Files:**
- None (manual testing)

**Step 1: Build**

```bash
go build
```

**Step 2: Test commands**

```bash
# Clean slate
rm -f todo.txt

# Add tasks
./togodo add "task one"
./togodo add "task two"

# List
./togodo list

# Toggle done
./togodo do 1

# Set priority
./togodo pri 2 A

# List again
./togodo list

# Tidy
./togodo tidy

# List final
./togodo list
```

**Step 3: Verify output matches expectations**

Check that:
- Add shows line number and confirmation
- List shows separator and summary
- Do shows confirmation
- Pri shows confirmation
- Tidy shows removed tasks

**Step 4: No commit needed**

---

### Task 5.4: Test TUI manually

**Files:**
- None (manual testing)

**Step 1: Launch TUI**

```bash
./togodo
```

**Step 2: Test operations**

- Press 'a', add a task
- Press Space to select, 'x' to toggle done
- Press Space to select, 'p' then 'A' to set priority
- Press '/' to filter
- Press 'q' to quit

**Step 3: Verify all operations work**

All TUI operations should work as before.

**Step 4: No commit needed**

---

### Task 5.5: Update FEATURE_PARITY_PLAN.md

**Files:**
- Modify: `FEATURE_PARITY_PLAN.md`

**Step 1: Update Phase 3 status**

Mark Phase 3 (Gap Analysis) as complete if architecture refactoring is done.

**Step 2: Update Phase 4 readiness**

Note that architecture is now ready for implementing missing commands.

**Step 3: Add architecture refactoring notes**

Add a section noting the architecture refactoring was completed and link to the design document.

**Step 4: Commit**

```bash
git add FEATURE_PARITY_PLAN.md
git commit -m "docs: update feature parity plan for architecture refactoring

Architecture refactoring complete.
Ready to implement missing commands.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Next Steps

After completing this refactoring, the architecture will be in place to implement missing commands. For each new command:

1. Create `cmd/xxx_logic.go` with business logic function
2. Create `cmd/xxx_logic_test.go` with tests
3. Create `internal/cli/xxx.go` with Cobra wrapper
4. Wire up in `cmd/root.go`
5. Run integration tests: `go test ./tests -run TestXxx`
6. Fix until tests pass

The pattern is established and the architecture is clean.

---

## Summary

This plan refactors togodo into a clean three-layer architecture:

**Phase 1:** Simplifies todotxtlib to a 12-method interface
**Phase 2:** Creates cmd business logic functions with result types
**Phase 3:** Migrates CLI to use cmd functions with proper formatting
**Phase 4:** Migrates TUI to use cmd functions
**Phase 5:** Cleans up and verifies everything works

Each task is bite-sized (2-5 minutes) and follows TDD principles where applicable. The refactoring is incremental and safe, with tests validating each step.

After completion, achieving feature parity becomes straightforward: follow the established pattern for each new command, let the integration tests guide implementation.
