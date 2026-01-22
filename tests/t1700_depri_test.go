package tests

import (
	"testing"
)

// TestDepriUsage tests depri command usage errors
// Ported from: t1700-depri.sh "depri usage"
func TestDepriUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("depri with invalid args", func(t *testing.T) {
		output, code := env.RunCommand("depri", "B", "B")
		expectedCode := 1
		expectedOutput := "usage: togodo depri NR [NR ...]"
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDepriNonexistent tests depri with non-existent task
// Ported from: t1700-depri.sh "depri nonexistent item"
func TestDepriNonexistent(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("depri nonexistent item", func(t *testing.T) {
		output, code := env.RunCommand("depri", "42")
		expectedCode := 1
		expectedOutput := "TODO: No task 42."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicDepriority tests basic depriority functionality
// Ported from: t1700-depri.sh "basic depriority"
func TestBasicDepriority(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("list before depri", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
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
		output, code := env.RunCommand("-p", "list")
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
// Ported from: t1700-depri.sh "multiple depriority"
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
		output, code := env.RunCommand("-p", "list")
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
// Ported from: t1700-depri.sh "depriority of unprioritized task"
func TestDepriUnprioritized(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("depri unprioritized task should error", func(t *testing.T) {
		output, code := env.RunCommand("depri", "3", "2")
		expectedCode := 1
		expectedOutput := `TODO: 3 is not prioritized.
2 notice the sunflowers
TODO: 2 deprioritized.`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
