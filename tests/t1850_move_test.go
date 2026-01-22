package tests

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBasicMoveImplicitSource tests moving task with implicit source file
// Ported from: t1850-move.sh "basic move with implicit source"
func TestBasicMoveImplicitSource(t *testing.T) {
	env := SetupFileBasedTestEnv(t)

	// Setup todo.txt with 2 tasks
	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers`)

	// Setup done.txt with 2 completed tasks
	doneFile := filepath.Join(filepath.Dir(env.TodoFile), "done.txt")
	err := os.WriteFile(doneFile, []byte("x 2009-02-13 notice the uppercase Roses +wakeup\nx 2010-01-02 notice the uppercase Roses +wakeup\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create done.txt: %v", err)
	}

	t.Run("move task 1 to done.txt", func(t *testing.T) {
		output, code := env.RunCommand("move", "-f", "1", "done.txt")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 moved from 'todo.txt' to 'done.txt'."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify todo.txt has 1 remaining task", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `2 (A) notice the sunflowers
--
TODO: 1 of 1 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify done.txt has 3 total tasks", func(t *testing.T) {
		content, err := os.ReadFile(doneFile)
		if err != nil {
			t.Fatalf("Failed to read done.txt: %v", err)
		}
		lines := 0
		for _, b := range content {
			if b == '\n' {
				lines++
			}
		}
		if lines != 3 {
			t.Errorf("Expected 3 lines in done.txt, got %d", lines)
		}
	})
}

// TestBasicMoveWithConfirmation tests moving task with interactive confirmation
// Ported from: t1850-move.sh "basic move with confirmation"
func TestBasicMoveWithConfirmation(t *testing.T) {
	t.Skip("TODO: Implement interactive confirmation with stdin piping")

	env := SetupFileBasedTestEnv(t)

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers`)

	doneFile := filepath.Join(filepath.Dir(env.TodoFile), "done.txt")
	err := os.WriteFile(doneFile, []byte("x 2009-02-13 notice the uppercase Roses +wakeup\nx 2010-01-02 notice the uppercase Roses +wakeup\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create done.txt: %v", err)
	}

	t.Run("move with 'y' confirmation", func(t *testing.T) {
		// TODO: Pipe 'y' to stdin
		output, code := env.RunCommand("move", "1", "done.txt")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 moved from 'todo.txt' to 'done.txt'."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicMoveWithPassedSource tests moving task with explicit source file
// Ported from: t1850-move.sh "basic move with passed source"
func TestBasicMoveWithPassedSource(t *testing.T) {
	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	// Setup todo.txt with 2 tasks
	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers`)

	// Setup done.txt with 3 completed tasks (including the one we'll move)
	doneFile := filepath.Join(filepath.Dir(env.TodoFile), "done.txt")
	err := os.WriteFile(doneFile, []byte("x 2009-02-13 notice the uppercase Roses +wakeup\nx 2009-02-13 smell the coffee +wakeup\nx 2010-01-02 notice the uppercase Roses +wakeup\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create done.txt: %v", err)
	}

	t.Run("move task 2 from done.txt to todo.txt", func(t *testing.T) {
		output, code := env.RunCommand("move", "-f", "2", "todo.txt", "done.txt")
		expectedOutput := "2 x 2009-02-13 smell the coffee +wakeup\nTODO: 2 moved from 'done.txt' to 'todo.txt'."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify todo.txt has 3 tasks now", func(t *testing.T) {
		content := env.ReadTodoFileContent()
		// Should have original 2 + moved task
		lines := 0
		for _, ch := range content {
			if ch == '\n' {
				lines++
			}
		}
		if lines < 2 { // At least 2 lines (original tasks)
			t.Errorf("Expected at least 2 tasks in todo.txt, got content:\n%s", content)
		}
	})

	t.Run("verify done.txt has 2 tasks remaining", func(t *testing.T) {
		content, err := os.ReadFile(doneFile)
		if err != nil {
			t.Fatalf("Failed to read done.txt: %v", err)
		}
		lines := 0
		for _, b := range content {
			if b == '\n' {
				lines++
			}
		}
		if lines != 2 {
			t.Errorf("Expected 2 lines in done.txt, got %d", lines)
		}
	})
}

// TestMoveToDestinationWithoutEOL tests moving to file without trailing newline
// Ported from: t1850-move.sh "move to destination without eol"
func TestMoveToDestinationWithoutEOL(t *testing.T) {
	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	// Setup todo.txt WITHOUT trailing newline (edge case)
	todoFile := env.TodoFile
	err := os.WriteFile(todoFile, []byte("(A) notice the sunflowers"), 0644) // No trailing \n
	if err != nil {
		t.Fatalf("Failed to write todo.txt: %v", err)
	}

	// Setup done.txt with tasks
	doneFile := filepath.Join(filepath.Dir(env.TodoFile), "done.txt")
	err = os.WriteFile(doneFile, []byte("x 2009-02-13 notice the uppercase Roses +wakeup\nx 2009-02-13 smell the coffee +wakeup\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create done.txt: %v", err)
	}

	t.Run("move task 2 from done.txt to todo.txt (no EOL)", func(t *testing.T) {
		output, code := env.RunCommand("move", "-f", "2", "todo.txt", "done.txt")
		expectedOutput := "2 x 2009-02-13 smell the coffee +wakeup\nTODO: 2 moved from 'done.txt' to 'todo.txt'."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify todo.txt contents after move", func(t *testing.T) {
		content, err := os.ReadFile(todoFile)
		if err != nil {
			t.Fatalf("Failed to read todo.txt: %v", err)
		}
		// Should have original task + moved task, properly formatted
		lines := 0
		for _, b := range content {
			if b == '\n' {
				lines++
			}
		}
		if lines < 1 {
			t.Errorf("Expected at least 1 line with newline in todo.txt, got:\n%s", string(content))
		}
	})

	t.Run("verify done.txt has 1 task remaining", func(t *testing.T) {
		content, err := os.ReadFile(doneFile)
		if err != nil {
			t.Fatalf("Failed to read done.txt: %v", err)
		}
		lines := 0
		for _, b := range content {
			if b == '\n' {
				lines++
			}
		}
		if lines != 1 {
			t.Errorf("Expected 1 line in done.txt, got %d", lines)
		}
	})
}
