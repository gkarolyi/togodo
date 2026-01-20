package tests

import (
	"testing"
)

// TestListpriBasic tests listing tasks by priority
// Ported from: t1250-listpri.sh
func TestListpriBasic(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
stop`)

	t.Run("listpri shows all prioritized tasks", func(t *testing.T) {
		output, code := env.RunCommand("listpri")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (C) notice the sunflowers
--
TODO: 2 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri with specific priority", func(t *testing.T) {
		output, code := env.RunCommand("listpri", "C")
		expectedOutput := `2 (C) notice the sunflowers
--
TODO: 1 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri case-insensitive", func(t *testing.T) {
		output, code := env.RunCommand("listpri", "c")
		expectedOutput := `2 (C) notice the sunflowers
--
TODO: 1 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListpriFiltering tests that listpri only shows valid priorities
// Ported from: t1250-listpri.sh
func TestListpriFiltering(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
(m)others will notice this
(n) not a prioritized task
notice the (C)opyright`)

	t.Run("listpri filters invalid priorities", func(t *testing.T) {
		output, code := env.RunCommand("listpri")
		// Should only show (B) and (C), not (m) or (n)
		// (m) and (n) are lowercase so not valid priorities
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (C) notice the sunflowers
--
TODO: 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
