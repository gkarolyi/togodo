package tests

import (
	"testing"
)

// TestDepriUsage tests depri command usage errors
// Ported from: t1700-depri.sh
func TestDepriUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("depri with invalid args", func(t *testing.T) {
		_, exitCode := env.RunCommand("depri", "B", "B")
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
		// TODO: Check error message matches "usage: togodo depri NR [NR ...]"
	})

	t.Run("depri nonexistent item", func(t *testing.T) {
		output, exitCode := env.RunCommand("depri", "42")
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
		// TODO: Check error message "TODO: No task 42."
		_ = output
	})
}

// TestBasicDepriority tests basic depriority functionality
// Ported from: t1700-depri.sh
func TestBasicDepriority(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("list before depri", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Should show priority-sorted order
		expectedOutput := `2 (A) notice the sunflowers
1 (B) smell the uppercase Roses +flowers @outside
3 stop
--
TODO: 3 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("deprioritize task 1", func(t *testing.T) {
		output, code := env.RunCommand("depri", "1")
		expectedOutput := "1 smell the uppercase Roses +flowers @outside\nTODO: 1 deprioritized."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after depri", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 (A) notice the sunflowers
1 smell the uppercase Roses +flowers @outside
3 stop
--
TODO: 3 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestMultipleDepriority tests deprioritizing multiple tasks
// Ported from: t1700-depri.sh
func TestMultipleDepriority(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
(C) stop`)

	t.Run("depri multiple tasks", func(t *testing.T) {
		output, code := env.RunCommand("depri", "3", "2")
		expectedOutput := `3 stop
TODO: 3 deprioritized.
2 notice the sunflowers
TODO: 2 deprioritized.`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after multiple depri", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 notice the sunflowers
3 stop
--
TODO: 3 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDepriUnprioritized tests deprioritizing an unprioritized task
// Ported from: t1700-depri.sh
func TestDepriUnprioritized(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("depri unprioritized task should error", func(t *testing.T) {
		output, code := env.RunCommand("depri", "3", "2")
		// First command (task 3) should error with exit code 1
		// Second command (task 2) should succeed
		// Overall exit code should be 1 due to first failure
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Should show "TODO: 3 is not prioritized."
		// TODO: Should still process task 2 and show "2 notice the sunflowers\nTODO: 2 deprioritized."
		_ = output
	})
}
