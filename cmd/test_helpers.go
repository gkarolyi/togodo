package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/gkarolyi/togodo/todotxtui"
)

// createTestTodos returns a slice of test todos for use in tests
func createTestTodos() []todotxtlib.Todo {
	return []todotxtlib.Todo{
		todotxtlib.NewTodo("(A) test todo 1 +project2 @context1"),
		todotxtlib.NewTodo("(B) test todo 2 +project1 @context2"),
		todotxtlib.NewTodo("x (C) test todo 3 +project1 @context1"),
	}
}

// setupEmptyTestBaseCommand creates a new BaseCommand with an empty repository for testing
func setupEmptyTestBaseCommand(tb testing.TB) (*BaseCommand, *bytes.Buffer) {
	// Create an empty buffer
	var buf bytes.Buffer

	// Create reader and writer that work with the buffer
	reader := todotxtlib.NewBufferReader(&buf)
	writer := todotxtlib.NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := todotxtlib.NewRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	// Create formatter and output writer for testing
	formatter := todotxtui.NewLipglossFormatter()
	output := todotxtui.NewStdoutWriter()

	// Create base command
	baseCmd := newBaseCommand(repo, formatter, output)

	return baseCmd, &buf
}

// setupTestBaseCommand creates a new BaseCommand with pre-populated test data
func setupTestBaseCommand(tb testing.TB) (*BaseCommand, *bytes.Buffer) {
	// Create a buffer with test todos
	var buf bytes.Buffer
	for _, todo := range createTestTodos() {
		buf.WriteString(todo.Text + "\n")
	}

	// Create reader and writer that work with the buffer
	reader := todotxtlib.NewBufferReader(&buf)
	writer := todotxtlib.NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := todotxtlib.NewRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	// Create formatter and output writer for testing
	formatter := todotxtui.NewLipglossFormatter()
	output := todotxtui.NewStdoutWriter()

	// Create base command
	baseCmd := newBaseCommand(repo, formatter, output)

	return baseCmd, &buf
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

// assertTodoCount asserts that the todo list has the expected number of items
func assertTodoCount(tb testing.TB, todos []todotxtlib.Todo, expectedCount int) {
	tb.Helper()
	if len(todos) != expectedCount {
		tb.Fatalf("Expected %d todos, got %d", expectedCount, len(todos))
	}
}

// assertContains asserts that the text contains the expected substring
func assertContains(tb testing.TB, text, substring string) {
	tb.Helper()
	if !strings.Contains(text, substring) {
		tb.Fatalf("Expected text to contain '%s', but it didn't. Text: '%s'", substring, text)
	}
}

// assertNotContains asserts that the text does not contain the given substring
func assertNotContains(tb testing.TB, text, substring string) {
	tb.Helper()
	if strings.Contains(text, substring) {
		tb.Fatalf("Expected text to not contain '%s', but it did. Text: '%s'", substring, text)
	}
}

// assertTodoText asserts that a todo has the expected text
func assertTodoText(tb testing.TB, todo todotxtlib.Todo, expectedText string) {
	tb.Helper()
	if todo.Text != expectedText {
		tb.Fatalf("Expected todo text to be '%s', got '%s'", expectedText, todo.Text)
	}
}

// assertTodoPriority asserts that a todo has the expected priority
func assertTodoPriority(tb testing.TB, todo todotxtlib.Todo, expectedPriority string) {
	tb.Helper()
	if todo.Priority != expectedPriority {
		tb.Fatalf("Expected todo priority to be '%s', got '%s'", expectedPriority, todo.Priority)
	}
}

// assertTodoCompleted asserts that a todo has the expected completion status
func assertTodoCompleted(tb testing.TB, todo todotxtlib.Todo, expectedCompleted bool) {
	tb.Helper()
	if todo.Done != expectedCompleted {
		tb.Fatalf("Expected todo completed status to be %v, got %v", expectedCompleted, todo.Done)
	}
}

// findTodoByText finds a todo in the list by its text content
func findTodoByText(todos []todotxtlib.Todo, text string) (todotxtlib.Todo, bool) {
	for _, todo := range todos {
		if todo.Text == text {
			return todo, true
		}
	}
	return todotxtlib.Todo{}, false
}

// assertTodoExists asserts that a todo with the given text exists in the list
func assertTodoExists(tb testing.TB, todos []todotxtlib.Todo, text string) {
	tb.Helper()
	_, found := findTodoByText(todos, text)
	if !found {
		tb.Fatalf("Expected to find todo with text '%s' in the list", text)
	}
}

// assertTodoNotExists asserts that a todo with the given text does not exist in the list
func assertTodoNotExists(tb testing.TB, todos []todotxtlib.Todo, text string) {
	tb.Helper()
	_, found := findTodoByText(todos, text)
	if found {
		tb.Fatalf("Expected not to find todo with text '%s' in the list", text)
	}
}

// Helper function to check if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
