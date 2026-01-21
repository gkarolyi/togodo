package tests

import (
	"testing"
)

// TestBasicMove tests moving tasks between files
// Ported from: t1850-move.sh
func TestBasicMove(t *testing.T) {
	env := SetupFileBasedTestEnv(t)

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers`)

	t.Run("move task to done.txt", func(t *testing.T) {
		output, code := env.RunCommand("move", "1", "done.txt")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}

		// Should show moved task and confirmation
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
TODO: 1 moved from 'todo.txt' to 'done.txt'.`
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify task removed from todo.txt", func(t *testing.T) {
		content := env.ReadTodoFileContent()
		expectedContent := `(A) notice the sunflowers`
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})

	t.Run("verify task added to done.txt", func(t *testing.T) {
		content := env.ReadDoneFileContent()
		expectedContent := `(B) smell the uppercase Roses +flowers @outside`
		if content != expectedContent {
			t.Errorf("Done file content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}

// TestMoveUsage tests move command usage
// Ported from: t1850-move.sh
func TestMoveUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("move without args", func(t *testing.T) {
		_, code := env.RunCommand("move")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message matches "usage: togodo move ITEM# DEST_FILE"
	})

	t.Run("move with only one arg", func(t *testing.T) {
		_, code := env.RunCommand("move", "1")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message matches "usage: togodo move ITEM# DEST_FILE"
	})

	t.Run("move with non-existent task", func(t *testing.T) {
		env.WriteTodoFile("task one")
		_, code := env.RunCommand("move", "42", "done.txt")
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		// TODO: Check error message "TODO: No task 42."
	})
}
