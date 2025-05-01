package todotxtlib

import (
	"bytes"
	"testing"
)

// createTestTodos returns a slice of test todos for use in tests
func createTestTodos() []Todo {
	return []Todo{
		NewTodo("(A) test todo 1 +project2 @context1"),
		NewTodo("(B) test todo 2 +project1 @context2"),
		NewTodo("x (C) test todo 3 +project1 @context1"),
	}
}

// setupTestRepository creates a new repository with a buffer reader and writer
func setupTestRepository(tb testing.TB) *Repository {
	// Create a buffer with test todos
	var buf bytes.Buffer
	for _, todo := range createTestTodos() {
		buf.WriteString(todo.Text + "\n")
	}

	// Create reader and writer that work with the buffer
	reader := NewBufferReader(&buf)
	writer := NewBufferWriter(&buf)

	// Create repository with the buffer-based reader and writer
	repo, err := NewRepository(reader, writer)
	if err != nil {
		tb.Fatalf("Failed to create test repository: %v", err)
	}

	return repo
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
