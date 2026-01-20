package tests

import (
	"testing"
)

// TestListconSingle tests listing contexts from tasks
// Ported from: t1310-listcon.sh
func TestListconSingle(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(A) @1 -- Some context 1 task, whitespace, one char
(A) @c2 -- Some context 2 task, whitespace, two char
@con03 -- Some context 3 task, no whitespace
@con04 -- Some context 4 task, no whitespace
@con05@con06 -- weird context`)

	t.Run("listcon shows all contexts", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		expectedOutput := `@1
@c2
@con03
@con04
@con05@con06`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListconMultiple tests listing contexts when tasks have multiple contexts
// Ported from: t1310-listcon.sh
func TestListconMultiple(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`@con01 -- Some context 1 task
@con02 -- Some context 2 task
@con02 @con03 -- Multi-context task`)

	t.Run("listcon with duplicates", func(t *testing.T) {
		output, code := env.RunCommand("listcon")
		// Should list unique contexts, possibly sorted
		// @con02 appears twice but should only be listed once
		expectedOutput := `@con01
@con02
@con03`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
