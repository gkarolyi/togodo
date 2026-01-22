package tests

import (
	"testing"
)

// TestBasicListall tests listing all tasks including completed ones
// Ported from: t1350-listall.sh "basic listall"
func TestBasicListall(t *testing.T) {
	t.Skip("TODO: listall requires done.txt support - should show tasks from both todo.txt and done.txt")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)
	// TODO: Also create done.txt with:
	// x 2011-12-01 eat breakfast
	// x 2011-12-05 smell the coffee +wakeup

	t.Run("listall with -p flag", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall")
		// Upstream shows tasks from both todo.txt and done.txt
		// Task numbers: 5,3,1 from todo.txt, 2,0,0,4 for completed tasks
		expectedOutput := `5 (A) stop
3 notice the sunflowers
1 smell the uppercase Roses +flowers @outside
2 x 2011-08-08 tend the garden @outside
0 x 2011-12-01 eat breakfast
0 x 2011-12-05 smell the coffee +wakeup
4 x 2011-12-26 go outside +wakeup
--
TODO: 5 of 5 tasks shown
DONE: 2 of 2 tasks shown
total 7 of 7 tasks shown`
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
	t.Skip("TODO: listall requires done.txt support")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)
	// TODO: Also create done.txt

	t.Run("listall with color codes", func(t *testing.T) {
		output, code := env.RunCommand("listall")
		// Upstream uses ANSI color codes:
		// - Priority (A) in yellow/bold: \033[1;33m
		// - Completed tasks in gray: \033[0;37m
		expectedOutput := "\033[1;33m5 (A) stop\033[0m\n3 notice the sunflowers\n1 smell the uppercase Roses +flowers @outside\n\033[0;37m2 x 2011-08-08 tend the garden @outside\033[0m\n\033[0;37m0 x 2011-12-01 eat breakfast\033[0m\n\033[0;37m0 x 2011-12-05 smell the coffee +wakeup\033[0m\n\033[0;37m4 x 2011-12-26 go outside +wakeup\033[0m\n--\nTODO: 5 of 5 tasks shown\nDONE: 2 of 2 tasks shown\ntotal 7 of 7 tasks shown"
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
	t.Skip("TODO: Implement TODOTXT_VERBOSE configuration and done.txt support")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)
	// TODO: Also create done.txt

	t.Run("listall with TODOTXT_VERBOSE=0", func(t *testing.T) {
		// TODO: Set TODOTXT_VERBOSE=0 environment variable
		output, code := env.RunCommand("-p", "listall")
		// Should show tasks WITHOUT summary statistics
		expectedOutput := `5 (A) stop
3 notice the sunflowers
1 smell the uppercase Roses +flowers @outside
2 x 2011-08-08 tend the garden @outside
0 x 2011-12-01 eat breakfast
0 x 2011-12-05 smell the coffee +wakeup
4 x 2011-12-26 go outside +wakeup`
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
	t.Skip("TODO: listall requires done.txt support")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)
	// TODO: Also create done.txt

	t.Run("listall @outside", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "@outside")
		expectedOutput := `1 smell the uppercase Roses +flowers @outside
2 x 2011-08-08 tend the garden @outside
--
TODO: 2 of 5 tasks shown
DONE: 0 of 2 tasks shown
total 2 of 7 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listall the", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "the")
		expectedOutput := `3 notice the sunflowers
1 smell the uppercase Roses +flowers @outside
2 x 2011-08-08 tend the garden @outside
0 x 2011-12-05 smell the coffee +wakeup
--
TODO: 3 of 5 tasks shown
DONE: 1 of 2 tasks shown
total 4 of 7 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("listall breakfast", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall", "breakfast")
		expectedOutput := `0 x 2011-12-01 eat breakfast
--
TODO: 0 of 5 tasks shown
DONE: 1 of 2 tasks shown
total 1 of 7 tasks shown`
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
DONE: 0 of 2 tasks shown
total 0 of 7 tasks shown`
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
	t.Skip("TODO: listall requires done.txt support")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
x 2011-08-08 tend the garden @outside
notice the sunflowers
x 2011-12-26 go outside +wakeup
(A) stop`)
	// TODO: Also create done.txt with 2 tasks, then add 4 more to done.txt

	t.Run("listall with double-digit task numbers", func(t *testing.T) {
		output, code := env.RunCommand("-p", "listall")
		// Upstream expects tasks numbered 0-9 with proper spacing for alignment
		// Numbers should be right-aligned when reaching double digits
		expectedOutput := ` 5 (A) stop
 3 notice the sunflowers
 1 smell the uppercase Roses +flowers @outside
 2 x 2011-08-08 tend the garden @outside
 0 x 2010-01-01 old task 1
 0 x 2010-01-01 old task 2
 0 x 2010-01-01 old task 3
 0 x 2010-01-01 old task 4
 0 x 2011-12-01 eat breakfast
 0 x 2011-12-05 smell the coffee +wakeup
 4 x 2011-12-26 go outside +wakeup
--
TODO: 5 of 5 tasks shown
DONE: 6 of 6 tasks shown
total 11 of 11 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
