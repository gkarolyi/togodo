package tests

import (
	"testing"
)

// TestDeduplicate tests removing duplicate tasks
// Ported from: t1910-deduplicate.sh
func TestDeduplicate(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`task one
task two
task one
task three
task two`)

	t.Run("deduplicate removes duplicates", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// TODO: Should show count of removed duplicates
		_ = output
	})

	t.Run("list after deduplicate", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Should only have 3 tasks now (one, two, three)
		// Each appearing once
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// TODO: Verify exact output
		_ = output
	})
}

// TestDeduplicateWithPriority tests deduplicate with prioritized tasks
// Ported from: t1910-deduplicate.sh
func TestDeduplicateWithPriority(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup: Add same task with different priorities
	env.WriteTodoFile("(A) task\n(B) task\ntask")

	// Run deduplicate
	output, code := env.RunCommand("deduplicate")
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	// Should report 2 duplicates removed
	if !containsSubstring(output, "2") {
		t.Errorf("Expected '2' duplicates removed in output, got: %s", output)
	}

	// Verify only the highest priority task remains
	content := env.ReadTodoFile()

	// Should have only 1 task
	lines := countNonEmptyLines(content)
	if lines != 1 {
		t.Errorf("Expected 1 task in file, got %d\nContent:\n%s", lines, content)
	}

	// Should be the (A) priority task
	if !containsSubstring(content, "(A) task") {
		t.Errorf("Expected '(A) task' in file, got:\n%s", content)
	}
}
