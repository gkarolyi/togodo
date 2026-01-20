package tests

import (
	"strings"
	"testing"
)

// TestAppendUsage tests append command usage
// Ported from: t1600-append.sh
func TestAppendUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("append with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("append", "adf", "asdfa")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "usage: togodo append NR \"TEXT TO APPEND\""
		_ = output
	})

	t.Run("append to nonexistent task", func(t *testing.T) {
		output, code := env.RunCommand("append", "10", "hej!")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "TODO: No task 10."
		_ = output
	})
}

// TestBasicAppend tests basic append functionality
// Ported from: t1600-append.sh
func TestBasicAppend(t *testing.T) {
	env := SetupTestEnv(t)
	env.RunCommand("add", "notice the daisies")

	t.Run("append text to task", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "smell the roses")
		expectedOutput := "1 notice the daisies smell the roses"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after append", func(t *testing.T) {
		output, code := env.RunCommand("list")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if !strings.Contains(output, "notice the daisies smell the roses") {
			t.Errorf("List should contain appended text, got:\n%s", output)
		}
	})
}
