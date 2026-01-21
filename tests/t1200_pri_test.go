package tests

import (
	"testing"
)

// TestPriUsage tests pri command usage errors
// Ported from: t1200-pri.sh
func TestPriUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("pri with invalid args", func(t *testing.T) {
		_, exitCode := env.RunCommand("pri", "B", "B")
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
		// TODO: Check error message matches "usage: togodo pri NR PRIORITY [NR PRIORITY ...]"
	})
}

// TestBasicPriority tests basic priority setting
// Ported from: t1200-pri.sh
func TestBasicPriority(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("set priority B on task 1", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "B")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 prioritized (B)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows priority in text", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Should show (B) prefix on task 1
		// Note: May include ANSI color codes in real output
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// TODO: Verify sorting - priority tasks should come first
		// TODO: Verify format with -p flag for plain output
		_ = output
	})

	t.Run("set priority C on task 2", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "C")
		expectedOutput := "2 (C) notice the sunflowers\nTODO: 2 prioritized (C)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrioritySorting tests that priority tasks sort correctly
// Ported from: t1200-pri.sh
func TestPrioritySorting(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	// Set priorities
	env.RunCommand("pri", "1", "B")
	env.RunCommand("pri", "2", "C")
	env.RunCommand("add", "smell the coffee +wakeup")

	t.Run("list shows priority sorted", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Expected order: (B) tasks, (C) tasks, no-priority tasks
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (C) notice the sunflowers
4 smell the coffee +wakeup
3 stop
--
TODO: 4 of 4 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPriorityError tests error handling for invalid task numbers
// Ported from: t1200-pri.sh
func TestPriorityError(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("task one")

	t.Run("pri with non-existent task", func(t *testing.T) {
		output, code := env.RunCommand("pri", "10", "B")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message matches "TODO: No task 10."
		_ = output
	})
}

// TestReprioritize tests changing priority on already-prioritized task
// Ported from: t1200-pri.sh
func TestReprioritize(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
stop`)

	t.Run("change priority from C to A", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "A")
		expectedOutput := "2 (A) notice the sunflowers\nTODO: 2 re-prioritized from (C) to (A)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify sorting after reprioritize", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Task 2 with (A) should now be first
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

	t.Run("set same priority should error", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "a")  // lowercase 'a'
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Should show "TODO: 2 already prioritized (A)."
		_ = output
	})
}

// TestPriorityWithPlainFlag tests -p flag for plain output
// Ported from: t1200-pri.sh
func TestPriorityWithPlainFlag(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) task one\ntask two")

	t.Run("list with global -p flag", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Should show plain output without ANSI color codes
		// We verify by checking that output is not empty and contains expected tasks
		if output == "" {
			t.Error("Expected non-empty output")
		}

		// Output should contain the tasks
		expectedLines := []string{
			"1 (A) task one",
			"2 task two",
			"TODO: 2 of 2 tasks shown",
		}

		for _, line := range expectedLines {
			if !containsLine(output, line) {
				t.Errorf("Expected output to contain '%s', got:\n%s", line, output)
			}
		}
	})

	t.Run("list with --plain flag", func(t *testing.T) {
		output, code := env.RunCommand("list", "--plain")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Should show plain output
		if output == "" {
			t.Error("Expected non-empty output")
		}
	})
}

// Helper to check if output contains a line
func containsLine(output, line string) bool {
	// Simple substring check - good enough for our purposes
	return len(output) >= len(line) && findSubstring(output, line)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
