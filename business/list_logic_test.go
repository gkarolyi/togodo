package business

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestList(t *testing.T) {
	t.Run("lists all tasks", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
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
		var buf bytes.Buffer
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
