package tests

import (
	"testing"
)

// TestReplaceUsage tests replace command usage
// Ported from: t1100-replace.sh
func TestReplaceUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("replace with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("replace", "adf", "asdfa")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "usage: togodo replace NR \"UPDATED ITEM\""
		_ = output
	})
}

// TestBasicReplace tests basic replace functionality
// Ported from: t1100-replace.sh
func TestBasicReplace(t *testing.T) {
	env := SetupTestEnv(t)
	env.RunCommand("add", "notice the daisies")

	t.Run("replace task text", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "smell the cows")
		expectedOutput := `1 notice the daisies
TODO: Replaced task with:
1 smell the cows`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after replace", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 smell the cows
--
TODO: 1 of 1 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
