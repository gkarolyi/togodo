package tests

import (
	"testing"
)

// TestMultilineSquashItemReplace tests replace with multiline input squashed
// Ported from: t2000-multiline.sh "multiline squash item replace"
func TestMultilineSquashItemReplace(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("smell the cheese")

	t.Run("replace with multiline squashes to one line", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "eat apples\neat oranges\ndrink milk")
		expectedOutput := "1 smell the cheese\nTODO: Replaced task with:\n1 eat apples eat oranges drink milk"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify file has one squashed task", func(t *testing.T) {
		content := env.ReadTodoFile()
		expectedContent := "eat apples eat oranges drink milk"
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}

// TestMultilineSquashItemAdd tests add with multiline input squashed
// Ported from: t2000-multiline.sh "multiline squash item add"
func TestMultilineSquashItemAdd(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("smell the cheese")

	t.Run("add with multiline squashes to one line", func(t *testing.T) {
		output, code := env.RunCommand("add", "eat apples\neat oranges\ndrink milk")
		expectedOutput := "2 eat apples eat oranges drink milk\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify file has two tasks", func(t *testing.T) {
		content := env.ReadTodoFile()
		expectedContent := "smell the cheese\neat apples eat oranges drink milk"
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}

// TestMultilineSquashItemAppend tests append with multiline input squashed
// Ported from: t2000-multiline.sh "multiline squash item append"
func TestMultilineSquashItemAppend(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("smell the cheese")

	t.Run("append with multiline squashes to one line", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "eat apples\neat oranges\ndrink milk")
		expectedOutput := "1 smell the cheese eat apples eat oranges drink milk"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify file has one squashed task", func(t *testing.T) {
		content := env.ReadTodoFile()
		expectedContent := "smell the cheese eat apples eat oranges drink milk"
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}

// TestMultilineSquashItemPrepend tests prepend with multiline input squashed
// Ported from: t2000-multiline.sh "multiline squash item prepend"
func TestMultilineSquashItemPrepend(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("smell the cheese")

	t.Run("prepend with multiline squashes to one line", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "eat apples\neat oranges\ndrink milk")
		expectedOutput := "1 eat apples eat oranges drink milk smell the cheese"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify file has one squashed task", func(t *testing.T) {
		content := env.ReadTodoFile()
		expectedContent := "eat apples eat oranges drink milk smell the cheese"
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}

// TestActualMultilineAdd tests addm command creating multiple tasks
// Ported from: t2000-multiline.sh "actual multiline add"
func TestActualMultilineAdd(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile("smell the cheese")

	t.Run("addm creates multiple separate tasks", func(t *testing.T) {
		output, code := env.RunCommand("addm", "eat apples\neat oranges\ndrink milk")
		expectedOutput := "2 eat apples\nTODO: 2 added.\n3 eat oranges\nTODO: 3 added.\n4 drink milk\nTODO: 4 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("verify file has four separate tasks", func(t *testing.T) {
		content := env.ReadTodoFile()
		expectedContent := "smell the cheese\neat apples\neat oranges\ndrink milk"
		if content != expectedContent {
			t.Errorf("File content mismatch\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
		}
	})
}
