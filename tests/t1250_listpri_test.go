package tests

import (
	"testing"
)

// TestBasicListpri tests basic priority listing functionality
// Ported from: t1250-listpri.sh "basic listpri"
func TestBasicListpri(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
stop`)

	t.Run("listpri A shows 0 tasks", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "A")
		expectedOutput := `--
TODO: 0 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri c (case-insensitive)", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "c")
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

// TestListpriHighlighting tests colored output for prioritized tasks
// Ported from: t1250-listpri.sh "listpri highlighting"
func TestListpriHighlighting(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
stop`)

	t.Run("listpri shows color codes", func(t *testing.T) {
		output, code := env.RunCommand("listpri")
		// Upstream expects ANSI color codes for priorities
		// Priority B is green: \033[0;32m
		// Priority C is blue: \033[0;34m
		expectedOutput := "\033[0;32m1 (B) smell the uppercase Roses +flowers @outside\033[0m\n\033[0;34m2 (C) notice the sunflowers\033[0m\n--\nTODO: 2 of 3 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestFilteringPriorities tests filtering by single priority letter
// Ported from: t1250-listpri.sh "filtering priorities"
func TestFilteringPriorities(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
(m)others will notice this
(n) not a prioritized task
notice the (C)opyright`)

	t.Run("listpri shows all priorities", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri")
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

	t.Run("listpri b", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "b")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
--
TODO: 1 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri c", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "c")
		expectedOutput := `2 (C) notice the sunflowers
--
TODO: 1 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri m (lowercase, invalid)", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "m")
		expectedOutput := `--
TODO: 0 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri n (invalid position)", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "n")
		expectedOutput := `--
TODO: 0 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestFilteringPriorityRanges tests priority range filtering (a-c, c-Z, A-, etc.)
// Ported from: t1250-listpri.sh "filtering priority ranges"
func TestFilteringPriorityRanges(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) one should be here
(C) two should be here also
(n) three should NOT be here
(n)(C) four should NOT be here either
(X) this should show up for priority x`)

	t.Run("listpri a-c", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "a-c")
		expectedOutput := `1 (B) one should be here
2 (C) two should be here also
--
TODO: 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri c-Z", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "c-Z")
		expectedOutput := `2 (C) two should be here also
5 (X) this should show up for priority x
--
TODO: 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri A-", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "A-")
		expectedOutput := `5 (X) this should show up for priority x
--
TODO: 1 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri A-C A-Z (AND logic)", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "A-C", "A-Z")
		expectedOutput := `--
TODO: 0 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri X A-Z", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "X", "A-Z")
		expectedOutput := `5 (X) this should show up for priority x
--
TODO: 1 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestFilteringConcatenation tests concatenating priorities and ranges
// Ported from: t1250-listpri.sh "concatenation of priorities and ranges"
func TestFilteringConcatenation(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) one should be here
(C) two should be here also
(n) three should NOT be here
(n)(C) four should NOT be here either
(X) this should show up for priority x`)

	t.Run("listpri CX", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "CX")
		expectedOutput := `2 (C) two should be here also
5 (X) this should show up for priority x
--
TODO: 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri ABR-Y", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "ABR-Y")
		expectedOutput := `1 (B) one should be here
2 (C) two should be here also
5 (X) this should show up for priority x
--
TODO: 3 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri A-", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "A-")
		expectedOutput := `5 (X) this should show up for priority x
--
TODO: 1 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestFilteringOfTERM tests combining priority filters with text search
// Ported from: t1250-listpri.sh "filtering of TERM"
func TestFilteringOfTERM(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(A) the first should be here
(B) one should be here
(C) two should be here also
(n) three should NOT be here
(n)(C) four should NOT be here either
(X) this should show up for priority x`)

	t.Run("listpri with search term", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "should be")
		expectedOutput := `2 (B) one should be here
3 (C) two should be here also
--
TODO: 2 of 6 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri a with search term", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "a", "should be")
		expectedOutput := `1 (A) the first should be here
--
TODO: 1 of 6 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri b second", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "b", "second")
		expectedOutput := `--
TODO: 0 of 6 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listpri x with search term", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listpri", "x", "should be")
		expectedOutput := `--
TODO: 0 of 6 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
