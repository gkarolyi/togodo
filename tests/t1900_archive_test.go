package tests

import (
	"testing"
)

// TestArchiveWithDuplicates tests archiving done tasks
// Ported from: t1900-archive.sh
func TestArchiveWithDuplicates(t *testing.T) {
	env := SetupFileBasedTestEnv(t)

	env.WriteTodoFileContent(`one
two
three
one
x done
four`)

	t.Run("archive done tasks", func(t *testing.T) {
		output, code := env.RunCommand("archive")
		// Should show archived task and confirmation
		expectedOutput := "x done\nTODO: todo.txt archived."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify done task removed from todo.txt", func(t *testing.T) {
		content := env.ReadTodoFileContent()
		expectedContent := `one
two
three
one
four`
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})

	t.Run("verify done task in done.txt", func(t *testing.T) {
		content := env.ReadDoneFileContent()
		expectedContent := "x done"
		if content != expectedContent {
			t.Errorf("Done file content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})

	t.Run("list after archive", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `5 four
1 one
4 one
3 three
2 two
--
TODO: 5 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestArchiveWarning tests archive when no done tasks exist
// Ported from: t1900-archive.sh
func TestArchiveWarning(t *testing.T) {
	env := SetupFileBasedTestEnv(t)

	env.WriteTodoFileContent(`one
two
three`)

	t.Run("archive with no done tasks", func(t *testing.T) {
		output, code := env.RunCommand("archive")
		expectedOutput := "TODO: todo.txt does not contain any done tasks."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
