package tests

import (
	"testing"
)

// TestDelUsage tests del command usage errors
// Ported from: t1800-del.sh "del usage"
func TestDelUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("del without args", func(t *testing.T) {
		output, code := env.RunCommand("del", "B")
		expectedCode := 1
		expectedOutput := "usage: togodo del NR [TERM]"
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDelNonexistent tests del with non-existent task
// Ported from: t1800-del.sh "del nonexistant item"
func TestDelNonexistent(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("del nonexistent item without term", func(t *testing.T) {
		output, code := env.RunCommand("del", "-f", "42")
		expectedCode := 1
		expectedOutput := "TODO: No task 42."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("del nonexistent item with term", func(t *testing.T) {
		output, code := env.RunCommand("del", "-f", "42", "Roses")
		expectedCode := 1
		expectedOutput := "TODO: No task 42."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicDel tests basic delete functionality
// Ported from: t1800-del.sh "basic del"
func TestBasicDel(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("list before delete", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (A) notice the sunflowers
3 stop
--
TODO: 3 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("delete task by number", func(t *testing.T) {
		output, code := env.RunCommand("del", "-f", "1")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 deleted."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after delete", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `2 (A) notice the sunflowers
3 stop
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

// TestDelWithConfirmation tests del command with interactive confirmation
// Ported from: t1800-del.sh "del with confirmation"
func TestDelWithConfirmation(t *testing.T) {
	t.Skip("TODO: Implement interactive confirmation prompts")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("del with 'n' response cancels", func(t *testing.T) {
		// TODO: Pipe 'n' to stdin
		output, code := env.RunCommand("del", "1")
		expectedCode := 1
		expectedOutput := "TODO: No tasks were deleted."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("del with 'x' response cancels", func(t *testing.T) {
		// TODO: Pipe 'x' to stdin
		output, code := env.RunCommand("del", "1")
		expectedCode := 1
		expectedOutput := "TODO: No tasks were deleted."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("del with empty response cancels", func(t *testing.T) {
		// TODO: Pipe empty line to stdin
		output, code := env.RunCommand("del", "1")
		expectedCode := 1
		expectedOutput := "TODO: No tasks were deleted."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("del with 'y' response confirms", func(t *testing.T) {
		// TODO: Pipe 'y' to stdin
		output, code := env.RunCommand("del", "1")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 deleted."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDelPreservingLineNumbers tests that line numbers persist after deletion
// Ported from: t1800-del.sh "del preserving line numbers"
func TestDelPreservingLineNumbers(t *testing.T) {
	t.Skip("TODO: Implement -n flag for renumbering tasks")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("delete task 1 preserves other line numbers", func(t *testing.T) {
		output, code := env.RunCommand("del", "-f", "1")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})

	t.Run("cannot re-delete task 1", func(t *testing.T) {
		output, code := env.RunCommand("del", "-f", "1")
		expectedCode := 1
		expectedOutput := "TODO: No task 1."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("add new tasks", func(t *testing.T) {
		env.RunCommand("add", "new task 1")
		env.RunCommand("add", "new task 2")
	})

	t.Run("del with -n flag renumbers", func(t *testing.T) {
		// TODO: Test -n flag behavior
		output, code := env.RunCommand("del", "-f", "-n", "3")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})
}

// TestBasicDelTerm tests removing terms from tasks
// Ported from: t1800-del.sh "basic del TERM"
func TestBasicDelTerm(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("setup and list", func(t *testing.T) {
		env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (A) notice the sunflowers
3 stop
--
TODO: 3 of 3 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("remove word 'uppercase'", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "uppercase")
		expectedOutput := "1 (B) smell the Roses +flowers @outside\nTODO: Removed 'uppercase' from task."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("remove phrase 'the Roses'", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "the Roses")
		expectedOutput := "1 (B) smell +flowers @outside\nTODO: Removed 'the Roses' from task."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("remove single character 'm'", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "m")
		expectedOutput := "1 (B) sell +flowers @outside\nTODO: Removed 'm' from task."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("remove context @outside", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "@outside")
		expectedOutput := "1 (B) sell +flowers\nTODO: Removed '@outside' from task."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("remove 'sell' (changes to 'l')", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "sell")
		expectedOutput := "1 (B) l +flowers\nTODO: Removed 'sell' from task."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDelNonexistentTerm tests removing a term that doesn't exist
// Ported from: t1800-del.sh "del nonexistant TERM"
func TestDelNonexistentTerm(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(A) notice the sunflowers
stop`)

	t.Run("remove nonexistent term", func(t *testing.T) {
		output, code := env.RunCommand("del", "1", "dung")
		expectedCode := 1
		expectedOutput := "TODO: 'dung' not found; no removal done."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
