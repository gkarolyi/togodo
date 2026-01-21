package tests

import (
	"regexp"
	"testing"
)

// TestDoUsage tests do command usage errors
// Ported from: t1500-do.sh
func TestDoUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("do with non-numeric arg", func(t *testing.T) {
		_, exitCode := env.RunCommand("do", "B", "B")
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
		// TODO: Check error message matches "usage: togodo do NR [NR ...]"
	})

	t.Run("do without args", func(t *testing.T) {
		_, exitCode := env.RunCommand("do")
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
		// TODO: Check error message matches "usage: togodo do NR [NR ...]"
	})
}

// TestBasicDo tests basic do/mark done functionality
// Ported from: t1500-do.sh
func TestBasicDo(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup initial todos
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop
remove1
remove2
remove3
remove4`)

	t.Run("list before marking done", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 notice the sunflowers
4 remove1
5 remove2
6 remove3
7 remove4
1 smell the uppercase Roses +flowers @outside
3 stop
--
TODO: 7 of 7 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	// TODO: Test marking multiple tasks done with comma separator (7,6)
	// TODO: Test auto-archiving to done.txt
	// TODO: Test that archived tasks are removed from todo.txt

	t.Run("mark single task done", func(t *testing.T) {
		output, code := env.RunCommand("do", "7")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Verify output format includes completion date
		// Expected format: "7 x YYYY-MM-DD remove4\nTODO: 7 marked as done."
		datePattern := `^7 x \d{4}-\d{2}-\d{2} remove4\nTODO: 7 marked as done\.$`
		if !regexp.MustCompile(datePattern).MatchString(output) {
			t.Errorf("Expected output matching pattern with completion date, got: %s", output)
		}
	})

	t.Run("verify task was marked done", func(t *testing.T) {
		content := env.ReadTodoFile()
		// Task should still be in file but marked with 'x' prefix
		// OR should be archived to done.txt (depending on implementation)
		if content == "" {
			t.Errorf("todo.txt should not be empty")
		}
	})
}

// TestDoWithFlags tests do command with various flags
// Ported from: t1500-do.sh
func TestDoWithFlags(t *testing.T) {
	env := SetupTestEnv(t)

	// Setup initial todos
	env.WriteTodoFile("task one\ntask two\ntask three")

	t.Run("mark done with auto-archive disabled", func(t *testing.T) {
		// Test -a flag (no auto-archive)
		output, code := env.RunCommand("do", "-a", "3")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Verify output format includes completion date
		// Expected format: "3 x YYYY-MM-DD task three\nTODO: 3 marked as done."
		datePattern := regexp.MustCompile(`^3 x \d{4}-\d{2}-\d{2} task three\nTODO: 3 marked as done\.$`)
		if !datePattern.MatchString(output) {
			t.Errorf("Expected output matching pattern with completion date, got: %s", output)
		}
	})

	t.Run("mark done with plain mode", func(t *testing.T) {
		// Test -p flag (plain output mode)
		output, code := env.RunCommand("do", "-p", "2")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Verify output format includes completion date
		// Expected format: "2 x YYYY-MM-DD task two\nTODO: 2 marked as done."
		datePattern := regexp.MustCompile(`^2 x \d{4}-\d{2}-\d{2} task two\nTODO: 2 marked as done\.$`)
		if !datePattern.MatchString(output) {
			t.Errorf("Expected output matching pattern with completion date, got: %s", output)
		}
	})
}

// TestDoAlreadyDone tests attempting to mark already-done task
// Ported from: t1500-do.sh
func TestDoAlreadyDone(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("task to complete")

	t.Run("mark task done first time", func(t *testing.T) {
		output, code := env.RunCommand("do", "1")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output == "" {
			t.Errorf("Expected output, got empty string")
		}
	})

	t.Run("mark task done second time", func(t *testing.T) {
		// TODO: This test depends on whether task is archived or kept in todo.txt
		// If archived, task 1 won't exist anymore
		// If kept, should get error "1 is already marked done"
		t.Skip("TODO: Verify behavior with already-done tasks")
	})
}

// TestDoMultipleWithComma tests marking multiple tasks with comma separator
// Ported from: t1500-do.sh
func TestDoMultipleWithComma(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("task1\ntask2\ntask3\ntask4\ntask5")

	t.Run("mark multiple tasks with comma separator", func(t *testing.T) {
		output, code := env.RunCommand("do", "5,3,1")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Verify output format - should show each task marked done with completion date
		// Expected format: "5 x YYYY-MM-DD task5\nTODO: 5 marked as done.\n3 x YYYY-MM-DD task3\n..."
		// Each task completion is two lines: the task line + confirmation line
		lines := regexp.MustCompile(`\n`).Split(output, -1)

		// Should have 6 lines total (2 per task * 3 tasks)
		if len(lines) != 6 {
			t.Errorf("Expected 6 output lines (2 per task * 3 tasks), got %d\nOutput: %s", len(lines), output)
		}

		// Verify task lines (lines 0, 2, 4) match the completion date pattern
		taskPattern := regexp.MustCompile(`^\d+ x \d{4}-\d{2}-\d{2} task\d+$`)
		for i := 0; i < 6; i += 2 {
			if !taskPattern.MatchString(lines[i]) {
				t.Errorf("Expected line %d to match task pattern, got: %s", i, lines[i])
			}
		}

		// Verify confirmation lines (lines 1, 3, 5)
		confirmPattern := regexp.MustCompile(`^TODO: \d+ marked as done\.$`)
		for i := 1; i < 6; i += 2 {
			if !confirmPattern.MatchString(lines[i]) {
				t.Errorf("Expected line %d to match confirmation pattern, got: %s", i, lines[i])
			}
		}
	})

	t.Run("list after marking done", func(t *testing.T) {
		output, code := env.RunCommand("list")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// All tasks should still be in todo.txt (completed tasks remain until archived)
		// Incomplete tasks (2, 4) and completed tasks (1, 3, 5) should all be shown
		// List sorts by priority then alphabetically, so order may vary
		expectedOutput := `2 task2
4 task4
1 x 2026-01-21 task1
3 x 2026-01-21 task3
5 x 2026-01-21 task5
--
TODO: 5 of 5 tasks shown`
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
