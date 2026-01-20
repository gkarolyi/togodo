package cmd

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
}
