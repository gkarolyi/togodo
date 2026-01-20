package business

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestAdd(t *testing.T) {
	t.Run("adds single task", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
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
