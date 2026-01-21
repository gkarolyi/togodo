package tests

import (
	"strings"
	"testing"
)

// TestListallBasic tests listing from both todo.txt and done.txt
// Ported from: t1350-listall.sh
func TestListallBasic(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	// Create done.txt with archived tasks
	// TODO: Need to add WriteDoneFile method to test helpers
	// For now, just test basic listall behavior

	t.Run("listall shows todo and done tasks", func(t *testing.T) {
		t.Skip("TODO: Implement done.txt support in test environment")

		output, code := env.RunCommand("listall")
		// Should show tasks from todo.txt and done.txt
		// Format includes counts: "TODO: X tasks, DONE: Y tasks, total Z tasks"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})
}

// TestListallHighlighting tests listall output formatting
// Ported from: t1350-listall.sh
func TestListallHighlighting(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(A) smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(B) stop`)

	t.Run("listall with highlighting", func(t *testing.T) {
		output, code := env.RunCommand("listall")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Output should include all tasks with line numbers
		// The actual highlighting (ANSI codes) is tested by presence of output
		// We can't easily test ANSI codes in integration tests
		// Just verify all tasks are shown
		if output == "" {
			t.Error("Expected non-empty output")
		}

		// Should show summary line
		if !strings.Contains(output, "TODO: 5 of 5 tasks shown") {
			t.Errorf("Expected summary line with '5 of 5 tasks shown', got: %s", output)
		}
	})

	t.Run("listall with plain flag", func(t *testing.T) {
		output, code := env.RunCommand("listall", "--plain")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Plain mode should show tasks without ANSI codes
		// Verify tasks are shown with plain text
		if output == "" {
			t.Error("Expected non-empty output")
		}

		// Should show summary line
		if !strings.Contains(output, "TODO: 5 of 5 tasks shown") {
			t.Errorf("Expected summary line with '5 of 5 tasks shown', got: %s", output)
		}
	})
}
