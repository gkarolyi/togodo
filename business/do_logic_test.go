package business

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestDo(t *testing.T) {
	t.Run("toggles task done", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
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

	t.Run("returns error for invalid index", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute with invalid index
		_, err := Do(repo, []int{99})

		// Assert
		if err == nil {
			t.Fatal("expected error for invalid index, got nil")
		}
	})

	t.Run("validates all indices before toggling any", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute with one valid and one invalid index
		_, err := Do(repo, []int{0, 99})

		// Assert error occurred
		if err == nil {
			t.Fatal("expected error for invalid index, got nil")
		}

		// Verify task 0 was NOT toggled (atomicity preserved)
		todos, _ := repo.ListAll()
		if todos[0].Done {
			t.Error("task 0 should not be toggled when operation fails")
		}
	})
}
