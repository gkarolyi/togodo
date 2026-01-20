package cmd

import (
	"bytes"
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
	service := todotxtlib.NewTodoService(repo)
	todos, err := service.SearchTodos(searchQuery)
	if err != nil {
		return "", err
	}

	// Format results similar to how the list command would display them
	formatter := cli.NewPlainFormatter()
	formatted := formatter.FormatList(todos)

	return strings.Join(formatted, "\n"), nil
}
