package tests

import (
	"testing"
)

// TestAppendUsage tests append command usage errors
// Ported from: t1600-append.sh "append usage"
func TestAppendUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("append with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("append", "foo", "bar", "baz")
		expectedCode := 1
		expectedOutput := `usage: togodo append NR "TEXT TO APPEND"`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAppendError tests append with non-existent task
// Ported from: t1600-append.sh "append error"
func TestAppendError(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`task1`)

	t.Run("append to non-existent task", func(t *testing.T) {
		output, code := env.RunCommand("append", "10", "something")
		expectedCode := 1
		expectedOutput := "TODO: No task 10."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicAppend tests basic append functionality
// Ported from: t1600-append.sh "basic append"
func TestBasicAppend(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`notice the daisies`)

	t.Run("append text to task", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "smell the roses")
		expectedOutput := "1 notice the daisies smell the roses"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows appended text", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 notice the daisies smell the roses
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

// TestBasicAppendWithAmpersand tests append with & character
// Ported from: t1600-append.sh "basic append with &"
func TestBasicAppendWithAmpersand(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`notice the daisies smell the roses`)

	t.Run("append text with ampersand", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "see the wasps & bees")
		expectedOutput := "1 notice the daisies smell the roses see the wasps & bees"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows appended text with ampersand", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 notice the daisies smell the roses see the wasps & bees
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

// TestAppendWithSpaces tests append preserves multiple spaces
// Ported from: t1600-append.sh "append with spaces"
func TestAppendWithSpaces(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`jump on hay`)

	t.Run("append preserves multiple spaces", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "and notice the   three   spaces")
		expectedOutput := "1 jump on hay and notice the   three   spaces"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAppendWithSymbols tests append with special characters
// Ported from: t1600-append.sh "append with symbols"
func TestAppendWithSymbols(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the cows
grow some corn
thrash some hay
chase the chickens`)

	t.Run("append with various symbols", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "~@#$%^&*()-_=+[{]}|;:',<.>/?")
		expectedOutput := "1 smell the cows ~@#$%^&*()-_=+[{]}|;:',<.>/?"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("append with backtick, exclamation, backslash, quote", func(t *testing.T) {
		output, code := env.RunCommand("append", "2", "`!\\\"")
		expectedOutput := "2 grow some corn `!\\\""
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after appends", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `4 chase the chickens
2 grow some corn ` + "`!\\\"" + `
1 smell the cows ~@#$%^&*()-_=+[{]}|;:',<.>/?
3 thrash some hay
--
TODO: 4 of 4 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAppendOfCurrentSentence tests append adds to current sentence
// Ported from: t1600-append.sh "append of current sentence"
func TestAppendOfCurrentSentence(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`notice the daisies`)

	t.Run("first append", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", ", lilies and roses")
		expectedOutput := "1 notice the daisies, lilies and roses"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("second append", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "; see the wasps")
		expectedOutput := "1 notice the daisies, lilies and roses; see the wasps"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("third append", func(t *testing.T) {
		output, code := env.RunCommand("append", "1", "& bees")
		expectedOutput := "1 notice the daisies, lilies and roses; see the wasps & bees"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestAppendOfCurrentSentenceDelimiters tests append with custom delimiters
// Ported from: t1600-append.sh "append of current sentence SENTENCE_DELIMITERS"
func TestAppendOfCurrentSentenceDelimiters(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SENTENCE_DELIMITERS configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`notice the daisies, lilies and roses; see the wasps & bees`)

	t.Run("append with & delimiter", func(t *testing.T) {
		// TODO: Set TODOTXT_SENTENCE_DELIMITERS='&'
		output, code := env.RunCommand("append", "1", "beans")
		expectedOutput := "1 notice the daisies, lilies and roses; see the wasps & bees&beans"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("append with % delimiter", func(t *testing.T) {
		// TODO: Set TODOTXT_SENTENCE_DELIMITERS='%'
		output, code := env.RunCommand("append", "1", "foo")
		expectedOutput := "1 notice the daisies, lilies and roses; see the wasps & bees&beans %foo"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("append with * delimiter", func(t *testing.T) {
		// TODO: Set TODOTXT_SENTENCE_DELIMITERS='*'
		output, code := env.RunCommand("append", "1", "2")
		expectedOutput := "1 notice the daisies, lilies and roses; see the wasps & bees&beans %foo*2"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
