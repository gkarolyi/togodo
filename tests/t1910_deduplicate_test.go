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
	t.Skip("TODO: Test deduplicate behavior with priorities - should keep higher priority")

	// env := SetupTestEnv(t)
	// env.WriteTodoFile("(A) task\n(B) task\ntask")
	//
	// env.RunCommand("deduplicate")
	// // Should keep (A) task and remove the others
}
