package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
)

// testTodos returns a slice of test todos for use in tests
var testTodos = []todotxtlib.Todo{
	todotxtlib.NewTodo("(A) test todo 1 +project2 @context1"),
	todotxtlib.NewTodo("(B) test todo 2 +project1 @context2"),
	todotxtlib.NewTodo("x (C) test todo 3 +project1 @context1"),
}

// setupEmptyTestRepository creates a new Repository with an empty buffer for testing
func setupEmptyTestRepository(tb testing.TB) (todotxtlib.TodoRepository, *bytes.Buffer) {
	// Create an empty buffer
	var buf bytes.Buffer

	// Create reader and writer that work with the buffer
	reader := todotxtlib.NewBufferReader(&buf)
	writer := todotxtlib.NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	return repo, &buf
}

// setupTestRepository creates a new Repository with pre-populated test data
func setupTestRepository(tb testing.TB) (todotxtlib.TodoRepository, *bytes.Buffer) {
	// Create a buffer with test todos
	var buf bytes.Buffer
	for _, todo := range testTodos {
		buf.WriteString(todo.Text + "\n")
	}

	// Create reader and writer that work with the buffer
	reader := todotxtlib.NewBufferReader(&buf)
	writer := todotxtlib.NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := todotxtlib.NewFileRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	return repo, &buf
}

// assertNoError asserts that the provided error is nil
func assertNoError(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("Expected no error, got: %v", err)
	}
}

// assertError asserts that the provided error is not nil
func assertError(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatal("Expected an error, but got nil")
	}
}

// assertContains asserts that the text contains the expected substring
func assertContains(tb testing.TB, text, substring string) {
	tb.Helper()
	if !strings.Contains(text, substring) {
		tb.Fatalf("Expected text to contain '%s', but it didn't. Text: '%s'", substring, text)
	}
}

// executeListForTest executes the list search and returns formatted output for testing
func executeListForTest(repo todotxtlib.TodoRepository, searchQuery string) (string, error) {
	// Search for todos
	var todos []todotxtlib.Todo
	var err error
	if searchQuery == "" {
		todos, err = repo.ListAll()
	} else {
		filter := todotxtlib.Filter{Text: searchQuery}
		todos, err = repo.Filter(filter)
	}
	if err != nil {
		return "", err
	}

	// Format results similar to how the list command would display them
	formatter := cli.NewPlainFormatter()
	formatted := formatter.FormatList(todos)

	return strings.Join(formatted, "\n"), nil
}

// getRepositoryString saves the repository to the buffer and returns its string content
func getRepositoryString(tb testing.TB, repo todotxtlib.TodoRepository, buf *bytes.Buffer) string {
	tb.Helper()

	// Reset buffer before saving to avoid appending to old content
	buf.Reset()

	// Save to buffer
	err := repo.Save()
	if err != nil {
		tb.Fatalf("Failed to save repository: %v", err)
	}

	return buf.String()
}

// addTodos adds multiple todos, sorts, and saves (mimics old service.AddTodos)
func addTodos(tb testing.TB, repo todotxtlib.TodoRepository, texts []string) []todotxtlib.Todo {
	tb.Helper()
	addedTodos := make([]todotxtlib.Todo, 0, len(texts))

	for _, text := range texts {
		todo, err := repo.Add(text)
		if err != nil {
			tb.Fatalf("Failed to add todo: %v", err)
		}
		addedTodos = append(addedTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		tb.Fatalf("Failed to save todos: %v", err)
	}

	return addedTodos
}

// toggleTodos toggles the done status of todos at the given indices
func toggleTodos(tb testing.TB, repo todotxtlib.TodoRepository, indices []int) []todotxtlib.Todo {
	tb.Helper()
	toggledTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.ToggleDone(index)
		if err != nil {
			tb.Fatalf("Failed to toggle todo at index %d: %v", index, err)
		}
		toggledTodos = append(toggledTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		tb.Fatalf("Failed to save todos: %v", err)
	}

	return toggledTodos
}

// toggleTodosWithError toggles the done status of todos at the given indices (returns error)
func toggleTodosWithError(repo todotxtlib.TodoRepository, indices []int) ([]todotxtlib.Todo, error) {
	toggledTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.ToggleDone(index)
		if err != nil {
			return nil, err
		}
		toggledTodos = append(toggledTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return nil, err
	}

	return toggledTodos, nil
}

// setPriorities sets the priority for todos at the given indices
func setPriorities(tb testing.TB, repo todotxtlib.TodoRepository, indices []int, priority string) []todotxtlib.Todo {
	tb.Helper()
	updatedTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.SetPriority(index, priority)
		if err != nil {
			tb.Fatalf("Failed to set priority for todo at index %d: %v", index, err)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := repo.Save(); err != nil {
		tb.Fatalf("Failed to save todos: %v", err)
	}

	return updatedTodos
}

// setPrioritiesWithError sets the priority for todos at the given indices (returns error)
func setPrioritiesWithError(repo todotxtlib.TodoRepository, indices []int, priority string) ([]todotxtlib.Todo, error) {
	updatedTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.SetPriority(index, priority)
		if err != nil {
			return nil, err
		}
		updatedTodos = append(updatedTodos, todo)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := repo.Save(); err != nil {
		return nil, err
	}

	return updatedTodos, nil
}

// removeDoneTodos removes all completed todos
func removeDoneTodos(tb testing.TB, repo todotxtlib.TodoRepository) []todotxtlib.Todo {
	tb.Helper()

	// Get done todos before removing
	doneFilter := todotxtlib.Filter{Done: "true"}
	doneTodos, err := repo.Filter(doneFilter)
	if err != nil {
		tb.Fatalf("Failed to list done todos: %v", err)
	}

	// Get all todos to iterate
	allTodos, err := repo.ListAll()
	if err != nil {
		tb.Fatalf("Failed to list all todos: %v", err)
	}

	// Remove backwards to avoid index shifting
	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			if _, err := repo.Remove(i); err != nil {
				tb.Fatalf("Failed to remove todo at index %d: %v", i, err)
			}
		}
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		tb.Fatalf("Failed to save todos: %v", err)
	}

	return doneTodos
}

// parseLineNumbers converts command line arguments to 0-based indices
// Line numbers in the command line are 1-based, but we need 0-based indices for the repository
func parseLineNumbers(args []string) ([]int, error) {
	indices := make([]int, 0, len(args))

	for _, arg := range args {
		lineNum, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert arg to int: %w", err)
		}
		if lineNum <= 0 {
			return nil, fmt.Errorf("line number must be positive, got %d", lineNum)
		}
		// Convert 1-based line number to 0-based index
		indices = append(indices, lineNum-1)
	}

	return indices, nil
}

// parsePriorityArgs parses arguments for the priority command
// Format: [LINE_NUMBER...] [PRIORITY]
// Returns: indices (0-based), priority, error
func parsePriorityArgs(args []string) ([]int, string, error) {
	if len(args) < 2 {
		return nil, "", fmt.Errorf("expected at least 2 arguments (line number and priority)")
	}

	// Last argument is the priority
	priority := args[len(args)-1]

	// All other arguments are line numbers
	lineArgs := args[:len(args)-1]
	indices, err := parseLineNumbers(lineArgs)
	if err != nil {
		return nil, "", err
	}

	return indices, priority, nil
}
