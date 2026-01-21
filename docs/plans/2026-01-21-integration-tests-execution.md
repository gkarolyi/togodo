# Integration Tests Execution Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Systematically run all integration tests and implement missing commands to achieve todo.txt-cli feature parity.

**Architecture:** Follow test-driven workflow - run test, categorize failure (missing command/wrong output/missing feature), fix minimally, commit. Focus changes on `internal/cli/` command wrappers, avoid touching `cmd/` or `todotxtlib/` unless feature completely missing.

**Tech Stack:** Go 1.24.4, Cobra CLI framework, buffer-based integration tests

---

## Context

- **Already complete**: t1000_add_list tests passing after Presenter removal
- **Remaining**: 19 test files to process sequentially
- **Skip for now**: t0000_config (more complex, do last)
- **Working directory**: `/Users/gergely.karolyi/Code/gkarolyi/todo-txt-cli-parity`

## CRITICAL: Tool Usage Rules

**YOU MUST follow these tool usage rules strictly. Violating these rules wastes time and causes failures.**

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

### Integration Tests Already Exist

- **DO NOT create new test files**
- **DO NOT write test scripts**
- The integration tests are already ported from bash in `/tests/` directory
- Just run them: `go test -v -run TestName`

### File Operations Examples

**WRONG:**
```bash
cat cmd/list.go  # DON'T USE BASH
echo "package cmd" > cmd/foo.go  # DON'T USE BASH
sed -i 's/old/new/' file.go  # DON'T USE BASH
```

**CORRECT:**
```
Use Read tool on cmd/list.go
Use Write tool to create cmd/foo.go
Use Edit tool to modify file.go
```

### Why This Matters

- Specialized tools are faster and more reliable
- They handle permissions correctly
- They work in all environments (sandbox, remote, etc.)
- Bash is slow and error-prone for file operations

**If you use bash for file operations, you are doing it wrong. Stop and use the correct tool.**

## Test Files to Process (In Order)

1. ✅ t1000_add_list - DONE
2. t1010_add_date - Date support (skip, Phase 2 feature)
3. t1040_add_priority - Priority in add input (skip, Phase 2 feature)
4. t1100_replace - Replace command
5. t1200_pri - Set priority (check if already works)
6. t1250_listpri - List by priority
7. t1310_listcon - List contexts
8. t1320_listproj - List projects
9. t1350_listall - List all including done
10. t1400_prepend - Prepend command
11. t1500_do - Mark done (check if already works)
12. t1600_append - Append command
13. t1700_depri - Remove priority
14. t1800_del - Delete command
15. t1850_move - Move command (skip, Phase 2 feature)
16. t1900_archive - Archive (skip, Phase 2 feature)
17. t1910_deduplicate - Deduplicate (skip, Phase 2 feature)
18. t1950_report - Report (skip, Phase 2 feature)
19. t2000_multiline - Multiline (skip, Phase 2 feature)
20. t2100_help - Help command

---

### Task 1: Fix TestListFiltering

**Files:**
- Modify: `internal/cli/list.go:39`

**Step 1: Run test to identify issue**

Run: `go test -v -run TestListFiltering`
Expected: FAIL - line numbers wrong in filtered output (shows 1 instead of 2)

**Step 2: Understand the issue**

The test shows:
- We have 2 tasks: "notice the daisies" and "smell the roses"
- When filtering by "smell", we get task #2
- But output shows "1 smell the roses" instead of "2 smell the roses"

Issue: The list command uses `i+1` for line numbers, which is correct for the filtered slice but wrong for original line numbers.

**Step 3: Check cmd.List to see what it returns**

Use Read tool on `cmd/list.go` to check if `ListResult` includes original line numbers.

**Step 4: Fix based on findings**

If `cmd.List` doesn't return original line numbers, we need to fix the business logic.
If it does, we need to use them in the CLI output.

**Likely fix in `internal/cli/list.go`:**

```go
// Format output to match todo.txt-cli
for _, todo := range result.Todos {
    // Use the original line number from the todo, not loop index
    fmt.Fprintf(command.OutOrStdout(), "%d %s\n", todo.LineNumber, todo.Text)
}
```

**Step 5: Run test to verify fix**

Run: `go test -v -run TestListFiltering`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/cli/list.go cmd/list.go
git commit -m "fix: use original line numbers in filtered list output"
```

---

### Task 2: Implement Replace Command

**Files:**
- Create: `cmd/replace.go`
- Create: `cmd/replace_test.go`
- Create: `internal/cli/replace.go`
- Modify: `internal/cli/root.go:64` (add command)

**Step 1: Run test to verify it fails**

Run: `go test -v -run TestBasicReplace`
Expected: FAIL - "replace" command not found

**Step 2: Write business logic with unit test**

Create: `cmd/replace.go`

```go
package cmd

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// ReplaceResult contains the result of a Replace operation
type ReplaceResult struct {
	OldTodo    todotxtlib.Todo
	NewTodo    todotxtlib.Todo
	LineNumber int
}

// Replace replaces the entire text of a todo at the given index
func Replace(repo todotxtlib.TodoRepository, index int, newText string) (ReplaceResult, error) {
	// Get existing todo
	oldTodo, err := repo.Get(index)
	if err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Create new todo with replaced text
	newTodo := todotxtlib.NewTodo(newText)

	// Update in repository
	updated, err := repo.Update(index, newTodo)
	if err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return ReplaceResult{
		OldTodo:    oldTodo,
		NewTodo:    updated,
		LineNumber: index + 1, // Convert to 1-based
	}, nil
}
```

**Step 3: Write unit test for Replace**

Create: `cmd/replace_test.go`

```go
package cmd

import (
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestReplace(t *testing.T) {
	repo := NewTestRepository(t)

	// Add initial todo
	repo.Add("notice the daisies")

	// Replace it
	result, err := Replace(repo, 0, "smell the cows")
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	// Verify old todo
	if result.OldTodo.Text != "notice the daisies" {
		t.Errorf("Expected old todo 'notice the daisies', got '%s'", result.OldTodo.Text)
	}

	// Verify new todo
	if result.NewTodo.Text != "smell the cows" {
		t.Errorf("Expected new todo 'smell the cows', got '%s'", result.NewTodo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].Text != "smell the cows" {
		t.Errorf("Expected 'smell the cows' in repo, got '%s'", todos[0].Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestReplace`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/replace.go`

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

// NewReplaceCmd creates a Cobra command for replacing todos
func NewReplaceCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "replace ITEM# \"UPDATED ITEM\"",
		Short: "Replace a todo item with new text",
		Long: `Replaces the entire text of a todo item.

# replace task 1 with new text
togodo replace 1 "new task text"
`,
		Args: cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get new text (join remaining args)
			newText := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Replace(repo, index, newText)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.OldTodo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: Replaced task with:\n")
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.NewTodo.Text)
			return nil
		},
	}
}
```

**Step 6: Register command in root**

Modify: `internal/cli/root.go`

Add after line 64:
```go
rootCmd.AddCommand(NewReplaceCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestBasicReplace`
Expected: PASS

**Step 8: Run all tests to ensure nothing broke**

Run: `go test ./cmd -v && go test -v -run TestBasicReplace`
Expected: All PASS

**Step 9: Commit**

```bash
git add cmd/replace.go cmd/replace_test.go internal/cli/replace.go internal/cli/root.go
git commit -m "feat: implement replace command

- Add cmd.Replace() business logic with unit tests
- Add cli.NewReplaceCmd() CLI wrapper
- Integration tests passing"
```

---

### Task 3: Check and Fix Pri Command Tests

**Files:**
- Test: `tests/t1200_pri_test.go`
- Potentially modify: `internal/cli/pri.go`

**Step 1: Run test to see current status**

Run: `go test -v -run TestBasicPriority`
Expected: Either PASS or FAIL with output format mismatch

**Step 2: If failing, analyze the diff**

Compare expected vs actual output from test failure.

**Step 3: Fix output format if needed**

Modify `internal/cli/pri.go` to match expected output exactly.

**Step 4: Run test to verify**

Run: `go test -v -run TestBasicPriority`
Expected: PASS

**Step 5: Commit if changes made**

```bash
git add internal/cli/pri.go
git commit -m "fix: adjust pri command output format to match todo.txt-cli"
```

---

### Task 4: Implement Listpri Command

**Files:**
- Create: `cmd/listpri.go`
- Create: `cmd/listpri_test.go`
- Create: `internal/cli/listpri.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test to verify it fails**

Run: `go test -v -run TestListpriBasic`
Expected: FAIL - "listpri" command not found

**Step 2: Write business logic**

Create: `cmd/listpri.go`

```go
package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListpriResult contains the result of a Listpri operation
type ListpriResult struct {
	Todos      []todotxtlib.Todo
	Priority   string
	TotalCount int
}

// Listpri lists todos with a specific priority
func Listpri(repo todotxtlib.TodoRepository, priority string) (ListpriResult, error) {
	// Get all todos
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListpriResult{}, err
	}

	// Filter by priority
	var filtered []todotxtlib.Todo
	for _, todo := range allTodos {
		if !todo.Done && todo.Priority == priority {
			filtered = append(filtered, todo)
		}
	}

	return ListpriResult{
		Todos:      filtered,
		Priority:   priority,
		TotalCount: len(allTodos),
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/listpri_test.go`

```go
package cmd

import (
	"testing"
)

func TestListpri(t *testing.T) {
	repo := NewTestRepository(t)

	// Add todos with different priorities
	repo.Add("(A) high priority")
	repo.Add("(B) medium priority")
	repo.Add("(A) another high")
	repo.Add("no priority")

	// List priority A
	result, err := Listpri(repo, "A")
	if err != nil {
		t.Fatalf("Listpri failed: %v", err)
	}

	// Should have 2 A-priority todos
	if len(result.Todos) != 2 {
		t.Errorf("Expected 2 A-priority todos, got %d", len(result.Todos))
	}

	// Verify they're the right ones
	if result.Todos[0].Text != "(A) high priority" {
		t.Errorf("Expected first todo '(A) high priority', got '%s'", result.Todos[0].Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestListpri`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/listpri.go`

```go
package cli

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListpriCmd creates a Cobra command for listing todos by priority
func NewListpriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listpri [PRIORITY]",
		Short: "List todos with specific priority",
		Long: `Lists all todos with the specified priority.

# list all A-priority tasks
togodo listpri A
`,
		Args: cobra.MaximumNArgs(1),
		Aliases: []string{"lsp"},
		RunE: func(command *cobra.Command, args []string) error {
			// Default to listing all priorities
			priority := ""
			if len(args) > 0 {
				priority = strings.ToUpper(args[0])
			}

			// Call business logic
			result, err := cmd.Listpri(repo, priority)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			// Find original line numbers for each todo
			allTodos, _ := repo.ListAll()
			for _, todo := range result.Todos {
				// Find line number in original list
				lineNum := 1
				for i, t := range allTodos {
					if t.Text == todo.Text {
						lineNum = i + 1
						break
					}
				}
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", lineNum, todo.Text)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", len(result.Todos), result.TotalCount)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewListpriCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestListpriBasic`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/listpri.go cmd/listpri_test.go internal/cli/listpri.go internal/cli/root.go
git commit -m "feat: implement listpri command"
```

---

### Task 5: Implement Listcon Command

**Files:**
- Create: `cmd/listcon.go`
- Create: `cmd/listcon_test.go`
- Create: `internal/cli/listcon.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestListconSingle`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/listcon.go`

```go
package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListconResult contains the result of listing contexts
type ListconResult struct {
	Contexts []string
}

// Listcon lists all contexts found in todos
func Listcon(repo todotxtlib.TodoRepository) (ListconResult, error) {
	contexts, err := repo.ListContexts()
	if err != nil {
		return ListconResult{}, err
	}

	return ListconResult{
		Contexts: contexts,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/listcon_test.go`

```go
package cmd

import (
	"testing"
)

func TestListcon(t *testing.T) {
	repo := NewTestRepository(t)

	// Add todos with contexts
	repo.Add("task @home @work")
	repo.Add("another @home")
	repo.Add("no context")

	result, err := Listcon(repo)
	if err != nil {
		t.Fatalf("Listcon failed: %v", err)
	}

	// Should have 2 contexts
	if len(result.Contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(result.Contexts))
	}

	// Should include @home and @work
	hasHome := false
	hasWork := false
	for _, ctx := range result.Contexts {
		if ctx == "@home" {
			hasHome = true
		}
		if ctx == "@work" {
			hasWork = true
		}
	}
	if !hasHome || !hasWork {
		t.Errorf("Expected @home and @work contexts, got %v", result.Contexts)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestListcon`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/listcon.go`

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListconCmd creates a Cobra command for listing contexts
func NewListconCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listcon",
		Short: "List all contexts",
		Long: `Lists all contexts (@context) found in todos.

# list all contexts
togodo listcon
`,
		Aliases: []string{"lsc"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listcon(repo)
			if err != nil {
				return err
			}

			// Format output
			for _, context := range result.Contexts {
				fmt.Fprintln(command.OutOrStdout(), context)
			}
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewListconCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestListconSingle`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/listcon.go cmd/listcon_test.go internal/cli/listcon.go internal/cli/root.go
git commit -m "feat: implement listcon command"
```

---

### Task 6: Implement Listproj Command

**Files:**
- Create: `cmd/listproj.go`
- Create: `cmd/listproj_test.go`
- Create: `internal/cli/listproj.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestListprojSingle`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/listproj.go`

```go
package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListprojResult contains the result of listing projects
type ListprojResult struct {
	Projects []string
}

// Listproj lists all projects found in todos
func Listproj(repo todotxtlib.TodoRepository) (ListprojResult, error) {
	projects, err := repo.ListProjects()
	if err != nil {
		return ListprojResult{}, err
	}

	return ListprojResult{
		Projects: projects,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/listproj_test.go`

```go
package cmd

import (
	"testing"
)

func TestListproj(t *testing.T) {
	repo := NewTestRepository(t)

	// Add todos with projects
	repo.Add("task +home +work")
	repo.Add("another +home")
	repo.Add("no project")

	result, err := Listproj(repo)
	if err != nil {
		t.Fatalf("Listproj failed: %v", err)
	}

	// Should have 2 projects
	if len(result.Projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(result.Projects))
	}

	// Should include +home and +work
	hasHome := false
	hasWork := false
	for _, proj := range result.Projects {
		if proj == "+home" {
			hasHome = true
		}
		if proj == "+work" {
			hasWork = true
		}
	}
	if !hasHome || !hasWork {
		t.Errorf("Expected +home and +work projects, got %v", result.Projects)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestListproj`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/listproj.go`

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListprojCmd creates a Cobra command for listing projects
func NewListprojCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listproj",
		Short: "List all projects",
		Long: `Lists all projects (+project) found in todos.

# list all projects
togodo listproj
`,
		Aliases: []string{"lsprj"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listproj(repo)
			if err != nil {
				return err
			}

			// Format output
			for _, project := range result.Projects {
				fmt.Fprintln(command.OutOrStdout(), project)
			}
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewListprojCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestListprojSingle`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/listproj.go cmd/listproj_test.go internal/cli/listproj.go internal/cli/root.go
git commit -m "feat: implement listproj command"
```

---

### Task 7: Implement Listall Command

**Files:**
- Create: `cmd/listall.go`
- Create: `cmd/listall_test.go`
- Create: `internal/cli/listall.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestListallBasic`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/listall.go`

```go
package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListallResult contains the result of listing all todos
type ListallResult struct {
	Todos      []todotxtlib.Todo
	TotalCount int
}

// Listall lists all todos including completed ones
func Listall(repo todotxtlib.TodoRepository) (ListallResult, error) {
	// Get all todos (including done)
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListallResult{}, err
	}

	return ListallResult{
		Todos:      allTodos,
		TotalCount: len(allTodos),
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/listall_test.go`

```go
package cmd

import (
	"testing"
)

func TestListall(t *testing.T) {
	repo := NewTestRepository(t)

	// Add todos
	repo.Add("active task")
	repo.Add("another active")

	// Mark one as done
	Do(repo, []int{0})

	// List all
	result, err := Listall(repo)
	if err != nil {
		t.Fatalf("Listall failed: %v", err)
	}

	// Should have 2 todos (including done)
	if len(result.Todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(result.Todos))
	}

	// Verify one is done
	doneCount := 0
	for _, todo := range result.Todos {
		if todo.Done {
			doneCount++
		}
	}
	if doneCount != 1 {
		t.Errorf("Expected 1 done todo, got %d", doneCount)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestListall`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/listall.go`

```go
package cli

import (
	"fmt"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewListallCmd creates a Cobra command for listing all todos
func NewListallCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "listall",
		Short: "List all todos including completed",
		Long: `Lists all todos including completed ones.

# list all tasks
togodo listall
`,
		Aliases: []string{"lsa"},
		RunE: func(command *cobra.Command, args []string) error {
			// Call business logic
			result, err := cmd.Listall(repo)
			if err != nil {
				return err
			}

			// Format output
			for i, todo := range result.Todos {
				fmt.Fprintf(command.OutOrStdout(), "%d %s\n", i+1, todo.Text)
			}
			fmt.Fprintln(command.OutOrStdout(), "--")
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d of %d tasks shown\n", result.TotalCount, result.TotalCount)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewListallCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestListallBasic`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/listall.go cmd/listall_test.go internal/cli/listall.go internal/cli/root.go
git commit -m "feat: implement listall command"
```

---

### Task 8: Implement Prepend Command

**Files:**
- Create: `cmd/prepend.go`
- Create: `cmd/prepend_test.go`
- Create: `internal/cli/prepend.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestBasicPrepend`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/prepend.go`

```go
package cmd

import (
	"fmt"
	"strings"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// PrependResult contains the result of a Prepend operation
type PrependResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Prepend prepends text to a todo, preserving priority if present
func Prepend(repo todotxtlib.TodoRepository, index int, text string) (PrependResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return PrependResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Extract priority if present
	todoText := todo.Text
	var newText string

	if len(todoText) >= 4 && todoText[0] == '(' && todoText[2] == ')' && todoText[3] == ' ' {
		// Has priority like "(A) task"
		priority := todoText[0:4] // "(A) "
		rest := todoText[4:]
		newText = priority + text + " " + rest
	} else {
		// No priority
		newText = text + " " + todoText
	}

	// Create updated todo
	updatedTodo := todotxtlib.NewTodo(newText)

	// Update in repository
	updated, err := repo.Update(index, updatedTodo)
	if err != nil {
		return PrependResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return PrependResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return PrependResult{
		Todo:       updated,
		LineNumber: index + 1,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/prepend_test.go`

```go
package cmd

import (
	"strings"
	"testing"
)

func TestPrepend(t *testing.T) {
	repo := NewTestRepository(t)

	// Add initial todo
	repo.Add("notice the sunflowers")

	// Prepend text
	result, err := Prepend(repo, 0, "really")
	if err != nil {
		t.Fatalf("Prepend failed: %v", err)
	}

	// Should have prepended text
	if !strings.Contains(result.Todo.Text, "really notice the sunflowers") {
		t.Errorf("Expected 'really notice the sunflowers', got '%s'", result.Todo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if !strings.Contains(todos[0].Text, "really notice the sunflowers") {
		t.Errorf("Expected prepended text in repo, got '%s'", todos[0].Text)
	}
}

func TestPrependPreservesPriority(t *testing.T) {
	repo := NewTestRepository(t)

	// Add prioritized todo
	repo.Add("(A) task with priority")

	// Prepend text
	result, err := Prepend(repo, 0, "important")
	if err != nil {
		t.Fatalf("Prepend failed: %v", err)
	}

	// Should preserve priority
	if result.Todo.Text != "(A) important task with priority" {
		t.Errorf("Expected '(A) important task with priority', got '%s'", result.Todo.Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestPrepend`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/prepend.go`

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

// NewPrependCmd creates a Cobra command for prepending to todos
func NewPrependCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "prepend ITEM# \"TEXT TO PREPEND\"",
		Short: "Prepend text to a todo item",
		Long: `Prepends text to the beginning of a todo item, preserving priority.

# prepend text to task 2
togodo prepend 2 "really"
`,
		Args: cobra.MinimumNArgs(2),
		Aliases: []string{"prep"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get text to prepend
			text := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Prepend(repo, index, text)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewPrependCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestBasicPrepend`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/prepend.go cmd/prepend_test.go internal/cli/prepend.go internal/cli/root.go
git commit -m "feat: implement prepend command"
```

---

### Task 9: Check and Fix Do Command Tests

**Files:**
- Test: `tests/t1500_do_test.go`
- Potentially modify: `internal/cli/do.go`

**Step 1: Run test**

Run: `go test -v -run TestBasicDo`
Expected: Either PASS or FAIL with output format mismatch

**Step 2: If failing, fix output format**

Adjust `internal/cli/do.go` to match expected output.

**Step 3: Run test to verify**

Run: `go test -v -run TestBasicDo`
Expected: PASS

**Step 4: Commit if changes made**

```bash
git add internal/cli/do.go
git commit -m "fix: adjust do command output format"
```

---

### Task 10: Implement Append Command

**Files:**
- Create: `cmd/append.go`
- Create: `cmd/append_test.go`
- Create: `internal/cli/append.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestBasicAppend`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/append.go`

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AppendResult contains the result of an Append operation
type AppendResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Append appends text to the end of a todo
func Append(repo todotxtlib.TodoRepository, index int, text string) (AppendResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return AppendResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Append text
	newText := todo.Text + " " + text

	// Create updated todo
	updatedTodo := todotxtlib.NewTodo(newText)

	// Update in repository
	updated, err := repo.Update(index, updatedTodo)
	if err != nil {
		return AppendResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return AppendResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return AppendResult{
		Todo:       updated,
		LineNumber: index + 1,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/append_test.go`

```go
package cmd

import (
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	repo := NewTestRepository(t)

	// Add initial todo
	repo.Add("notice the daisies")

	// Append text
	result, err := Append(repo, 0, "smell the roses")
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Should have appended text
	expected := "notice the daisies smell the roses"
	if result.Todo.Text != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.Todo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if !strings.Contains(todos[0].Text, "smell the roses") {
		t.Errorf("Expected appended text in repo, got '%s'", todos[0].Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestAppend`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/append.go`

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

// NewAppendCmd creates a Cobra command for appending to todos
func NewAppendCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "append ITEM# \"TEXT TO APPEND\"",
		Short: "Append text to a todo item",
		Long: `Appends text to the end of a todo item.

# append text to task 1
togodo append 1 "additional text"
`,
		Args: cobra.MinimumNArgs(2),
		Aliases: []string{"app"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get text to append
			text := strings.Join(args[1:], " ")

			// Call business logic
			result, err := cmd.Append(repo, index, text)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewAppendCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestBasicAppend`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/append.go cmd/append_test.go internal/cli/append.go internal/cli/root.go
git commit -m "feat: implement append command"
```

---

### Task 11: Implement Depri Command

**Files:**
- Create: `cmd/depri.go`
- Create: `cmd/depri_test.go`
- Create: `internal/cli/depri.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestBasicDepriority`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/depri.go`

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DepriResult contains the result of a Depri operation
type DepriResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Depri removes priority from a todo
func Depri(repo todotxtlib.TodoRepository, index int) (DepriResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return DepriResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Remove priority by setting to empty string
	updated, err := repo.SetPriority(index, "")
	if err != nil {
		return DepriResult{}, fmt.Errorf("failed to remove priority: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return DepriResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DepriResult{
		Todo:       updated,
		LineNumber: index + 1,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/depri_test.go`

```go
package cmd

import (
	"testing"
)

func TestDepri(t *testing.T) {
	repo := NewTestRepository(t)

	// Add prioritized todo
	repo.Add("(A) high priority task")

	// Remove priority
	result, err := Depri(repo, 0)
	if err != nil {
		t.Fatalf("Depri failed: %v", err)
	}

	// Should have no priority
	if result.Todo.Priority != "" {
		t.Errorf("Expected no priority, got '%s'", result.Todo.Priority)
	}

	// Text should not have (A) prefix
	if result.Todo.Text == "(A) high priority task" {
		t.Errorf("Expected priority removed from text, got '%s'", result.Todo.Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestDepri`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/depri.go`

```go
package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDepriCmd creates a Cobra command for removing priority
func NewDepriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "depri ITEM#",
		Short: "Remove priority from a todo item",
		Long: `Removes the priority from a todo item.

# remove priority from task 1
togodo depri 1
`,
		Args: cobra.ExactArgs(1),
		Aliases: []string{"dp"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Call business logic
			result, err := cmd.Depri(repo, index)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d deprioritized.\n", result.LineNumber)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewDepriCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestBasicDepriority`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/depri.go cmd/depri_test.go internal/cli/depri.go internal/cli/root.go
git commit -m "feat: implement depri command"
```

---

### Task 12: Implement Del Command

**Files:**
- Create: `cmd/del.go`
- Create: `cmd/del_test.go`
- Create: `internal/cli/del.go`
- Modify: `internal/cli/root.go`

**Step 1: Run test**

Run: `go test -v -run TestBasicDel`
Expected: FAIL - command not found

**Step 2: Write business logic**

Create: `cmd/del.go`

```go
package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DelResult contains the result of a Del operation
type DelResult struct {
	DeletedTodo todotxtlib.Todo
	LineNumber  int
}

// Del deletes a todo from the repository
func Del(repo todotxtlib.TodoRepository, index int) (DelResult, error) {
	// Remove todo
	deleted, err := repo.Remove(index)
	if err != nil {
		return DelResult{}, fmt.Errorf("failed to remove todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return DelResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelResult{
		DeletedTodo: deleted,
		LineNumber:  index + 1,
	}, nil
}
```

**Step 3: Write unit test**

Create: `cmd/del_test.go`

```go
package cmd

import (
	"testing"
)

func TestDel(t *testing.T) {
	repo := NewTestRepository(t)

	// Add todos
	repo.Add("task 1")
	repo.Add("task 2")

	// Delete first todo
	result, err := Del(repo, 0)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Should have deleted task 1
	if result.DeletedTodo.Text != "task 1" {
		t.Errorf("Expected deleted 'task 1', got '%s'", result.DeletedTodo.Text)
	}

	// Should only have 1 todo left
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo remaining, got %d", len(todos))
	}

	// Remaining should be task 2
	if todos[0].Text != "task 2" {
		t.Errorf("Expected remaining 'task 2', got '%s'", todos[0].Text)
	}
}
```

**Step 4: Run unit test**

Run: `go test ./cmd -v -run TestDel`
Expected: PASS

**Step 5: Create CLI command**

Create: `internal/cli/del.go`

```go
package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewDelCmd creates a Cobra command for deleting todos
func NewDelCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "del ITEM#",
		Short: "Delete a todo item",
		Long: `Deletes a todo item from the list.

# delete task 1
togodo del 1
`,
		Args: cobra.ExactArgs(1),
		Aliases: []string{"rm"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Call business logic
			result, err := cmd.Del(repo, index)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.DeletedTodo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d deleted.\n", result.LineNumber)
			return nil
		},
	}
}
```

**Step 6: Register command**

Modify: `internal/cli/root.go`

Add:
```go
rootCmd.AddCommand(NewDelCmd(repo))
```

**Step 7: Run integration test**

Run: `go test -v -run TestBasicDel`
Expected: PASS

**Step 8: Commit**

```bash
git add cmd/del.go cmd/del_test.go internal/cli/del.go internal/cli/root.go
git commit -m "feat: implement del command"
```

---

### Task 13: Run All Integration Tests and Document Results

**Files:**
- None

**Step 1: Run full test suite**

Run: `go test ./tests -v 2>&1 | tee test-results.txt`

**Step 2: Count results**

Run: `grep -E "^(PASS|FAIL|SKIP)" test-results.txt | sort | uniq -c`

**Step 3: Create summary document**

Create a summary of:
- How many tests pass
- How many are skipped (with reasons)
- Any unexpected failures

**Step 4: Identify Phase 2 work**

List all skipped tests and group by:
- Date support needed
- Archive/done.txt support needed
- Other features needed

---

## Phase 2 Features (To be planned separately)

After all commands are implemented, return to skipped tests:

1. **Date Support** (t1010, t1040)
   - Creation dates
   - Completion dates
   - Date formatting

2. **done.txt Archive** (t1900, t1850_move)
   - Archive command
   - Move to done.txt
   - done.txt file handling

3. **Advanced Features** (t1910, t1950, t2000)
   - Deduplication
   - Reporting
   - Multiline support

4. **Config Command** (t0000)
   - Read/write config
   - List settings

---

## Success Criteria

- All command tests either PASS or are SKIPped with documented reason
- Each command has unit tests in `cmd/`
- Each command has CLI wrapper in `internal/cli/`
- No changes to `todotxtlib/` (unless absolutely necessary)
- Frequent commits (one per command)
- All existing tests still pass

## Notes

- Use `cmd.OutOrStdout()` for all output
- Follow TDD: unit test first, then implementation
- Match exact output format from todo.txt-cli
- Keep commits small and focused
