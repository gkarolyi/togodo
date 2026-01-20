package business

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestTidy(t *testing.T) {
	t.Run("removes done tasks", func(t *testing.T) {
		// Setup
		buf := bytes.Buffer{}
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
