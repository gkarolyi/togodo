package tests

import (
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
	t.Skip("TODO: Implement listall command with done.txt support")

	// env := SetupTestEnv(t)
	// listall should show:
	// - Active tasks from todo.txt
	// - Done tasks (with strikethrough or different color)
	// - Summary line showing todo count, done count, total count
}
