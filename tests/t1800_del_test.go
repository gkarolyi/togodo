package tests

import (
	"testing"
)

// TestDelUsage tests del command usage
// Ported from: t1800-del.sh
func TestDelUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("del without args", func(t *testing.T) {
		output, code := env.RunCommand("del")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "usage: togodo del NR [TERM]"
		_ = output
	})

	t.Run("del nonexistent item", func(t *testing.T) {
		output, code := env.RunCommand("del", "42")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "TODO: No task 42."
		_ = output
	})
}

// TestBasicDel tests basic delete functionality
// Ported from: t1800-del.sh
func TestBasicDel(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("delete task by number", func(t *testing.T) {
		output, code := env.RunCommand("del", "3")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// TODO: Should show confirmation "3 stop\nTODO: 3 deleted."
		// TODO: May require -f flag to force delete without confirmation
		_ = output
	})

	t.Run("list after delete", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 (A) notice the sunflowers
1 (B) smell the uppercase Roses +flowers @outside
--
TODO: 2 of 2 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Task 3 "stop" should be gone
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDelMultiple tests deleting multiple tasks
// Ported from: t1800-del.sh
func TestDelMultiple(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("task1\ntask2\ntask3\ntask4")

	t.Run("delete multiple tasks", func(t *testing.T) {
		output, code := env.RunCommand("del", "4", "2")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Should show both deleted tasks
		// Order matches deletion order: 4 first, then 2
		expectedOutput := `4 task4
2 task2
TODO: 2 tasks deleted.`
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after delete", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 task1
2 task3
--
TODO: 2 of 2 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDelWithTerm tests deleting terms from a task
// Ported from: t1800-del.sh
func TestDelWithTerm(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) buy milk and eggs +grocery")

	t.Run("remove term from task", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "eggs")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Should show modified task and confirmation
		expectedOutput := `1 (A) buy milk and +grocery
TODO: Removed 'eggs' from task.`
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after term removal", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 (A) buy milk and +grocery
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
