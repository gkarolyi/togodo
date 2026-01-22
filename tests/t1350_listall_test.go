package tests

import (
	"testing"
)

// TestBasicListall tests listing all tasks including completed ones
// Ported from: t1350-listall.sh "basic listall"
func TestBasicListall(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	t.Run("listall with -p flag", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall")
		// Upstream shows: priority tasks first, then active, then completed
		expectedOutput := `1 (A) stop
2 smell the uppercase Roses +flowers @outside
4 notice the sunflowers
3 x 2011-08-08 tend the garden @outside
5 x 2011-12-26 go outside +wakeup
--
TODO: 3 of 5 tasks shown
DONE: 2 of 5 tasks shown
total 5 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListallHighlighting tests ANSI color codes in listall output
// Ported from: t1350-listall.sh "listall highlighting"
func TestListallHighlighting(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	t.Run("listall with color codes", func(t *testing.T) {
		output, code := env.RunCommand("listall")
		// Upstream uses ANSI color codes:
		// - Priority (A) in yellow/bold
		// - Completed tasks in gray
		expectedOutput := "\033[1;33m1 (A) stop\033[0m\n2 smell the uppercase Roses +flowers @outside\n4 notice the sunflowers\n\033[0;37m3 x 2011-08-08 tend the garden @outside\033[0m\n\033[0;37m5 x 2011-12-26 go outside +wakeup\033[0m\n--\nTODO: 3 of 5 tasks shown\nDONE: 2 of 5 tasks shown\ntotal 5 of 5 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListallNonverbose tests listall without summary statistics
// Ported from: t1350-listall.sh "listall nonverbose"
func TestListallNonverbose(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_VERBOSE configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	t.Run("listall with TODOTXT_VERBOSE=0", func(t *testing.T) {
		// TODO: Set TODOTXT_VERBOSE=0 environment variable
		output, code := env.RunCommand("-p", "listall")
		// Should show tasks WITHOUT summary statistics
		expectedOutput := `1 (A) stop
2 smell the uppercase Roses +flowers @outside
4 notice the sunflowers
3 x 2011-08-08 tend the garden @outside
5 x 2011-12-26 go outside +wakeup`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListallFiltering tests filtering listall output
// Ported from: t1350-listall.sh "listall filtering"
func TestListallFiltering(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	t.Run("listall @outside", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "@outside")
		expectedOutput := `2 smell the uppercase Roses +flowers @outside
3 x 2011-08-08 tend the garden @outside
--
TODO: 1 of 5 tasks shown
DONE: 1 of 5 tasks shown
total 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listall the", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "the")
		expectedOutput := `2 smell the uppercase Roses +flowers @outside
3 x 2011-08-08 tend the garden @outside
4 notice the sunflowers
--
TODO: 2 of 5 tasks shown
DONE: 1 of 5 tasks shown
total 3 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listall breakfast", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "breakfast")
		expectedOutput := `--
TODO: 0 of 5 tasks shown
DONE: 0 of 5 tasks shown
total 0 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listall doesnotmatch", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "doesnotmatch")
		expectedOutput := `--
TODO: 0 of 5 tasks shown
DONE: 0 of 5 tasks shown
total 0 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListallNumberWidth tests that task numbers adjust width dynamically
// Ported from: t1350-listall.sh "listall number width"
func TestListallNumberWidth(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)

	t.Run("listall before adding tasks", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall")
		expectedOutput := `1 (A) stop
2 smell the uppercase Roses +flowers @outside
4 notice the sunflowers
3 x 2011-08-08 tend the garden @outside
5 x 2011-12-26 go outside +wakeup
--
TODO: 3 of 5 tasks shown
DONE: 2 of 5 tasks shown
total 5 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	// Add 5 more tasks to reach 10 total
	t.Run("add five more tasks", func(t *testing.T) {
		env.RunCommand("add", "task 6")
		env.RunCommand("add", "task 7")
		env.RunCommand("add", "task 8")
		env.RunCommand("add", "task 9")
		env.RunCommand("add", "task 10")
	})

	t.Run("listall after adding tasks", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall")
		// Now task numbers should be double-digit, verify proper alignment
		expectedOutput := ` 1 (A) stop
 2 smell the uppercase Roses +flowers @outside
 4 notice the sunflowers
 6 task 6
 7 task 7
 8 task 8
 9 task 9
10 task 10
 3 x 2011-08-08 tend the garden @outside
 5 x 2011-12-26 go outside +wakeup
--
TODO: 8 of 10 tasks shown
DONE: 2 of 10 tasks shown
total 10 of 10 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
