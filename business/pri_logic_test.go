package business

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestSetPriority(t *testing.T) {
	t.Run("sets priority on task", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
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

	t.Run("returns error for invalid index", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute with invalid index
		_, err := SetPriority(repo, []int{99}, "A")

		// Assert
		if err == nil {
			t.Fatal("expected error for invalid index, got nil")
		}
	})

	t.Run("validates all indices before setting priority on any", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute with one valid and one invalid index
		_, err := SetPriority(repo, []int{0, 99}, "A")

		// Assert error occurred
		if err == nil {
			t.Fatal("expected error for invalid index, got nil")
		}

		// Verify task 0 priority was NOT set (atomicity preserved)
		todos, _ := repo.ListAll()
		if todos[0].Priority != "" {
			t.Errorf("task 0 priority should not be set when operation fails, got %s", todos[0].Priority)
		}
	})
}
