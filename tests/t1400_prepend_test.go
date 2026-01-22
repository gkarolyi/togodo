package tests

import (
	"testing"
)

// TestPrependUsage tests prepend command usage errors
// Ported from: t1400-prepend.sh "prepend usage"
func TestPrependUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("prepend with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "adf", "asdfa")
		expectedCode := 1
		expectedOutput := `usage: togodo prepend NR "TEXT TO PREPEND"`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicPrepend tests basic prepend functionality
// Ported from: t1400-prepend.sh "basic prepend"
func TestBasicPrepend(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("list before prepend", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 notice the sunflowers
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

	t.Run("prepend text to task 2", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "2", "test")
		expectedOutput := "2 test notice the sunflowers"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("prepend text to task 1", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "test")
		expectedOutput := "1 (B) test smell the uppercase Roses +flowers @outside"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrependWithAmpersand tests prepend with & character
// Ported from: t1400-prepend.sh "prepend with &"
func TestPrependWithAmpersand(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`stop`)

	t.Run("prepend text with ampersand", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "no running & jumping now")
		expectedOutput := "1 no running & jumping now stop"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrependWithSpaces tests prepend preserves multiple spaces
// Ported from: t1400-prepend.sh "prepend with spaces"
func TestPrependWithSpaces(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`notice the   three   spaces and jump on hay`)

	t.Run("prepend preserves multiple spaces", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "really")
		expectedOutput := "1 really notice the   three   spaces and jump on hay"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrependWithSymbols tests prepend with special characters
// Ported from: t1400-prepend.sh "prepend with symbols"
func TestPrependWithSymbols(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the cows
grow some corn
eat some chicken
stop`)

	t.Run("prepend with various symbols", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "~@#$%^&*()-_=+[{]}|;:',<.>/?")
		expectedOutput := "1 ~@#$%^&*()-_=+[{]}|;:',<.>/? smell the cows"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("prepend with backtick, exclamation, backslash", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "2", "`!\\")
		expectedOutput := "2 `!\\ grow some corn"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after prepends", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 ~@#$%^&*()-_=+[{]}|;:',<.>/? smell the cows
2 ` + "`!\\" + ` grow some corn
3 eat some chicken
4 stop
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

// TestPrependHandlingPrependedDateOnAdd tests prepend with auto-added dates
// Ported from: t1400-prepend.sh "prepend handling prepended date on add"
func TestPrependHandlingPrependedDateOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add task with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "new task")
		expectedOutput := "1 2009-02-13 new task\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("prepend to dated task", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "this is just a")
		expectedOutput := "1 2009-02-13 this is just a new task"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrependHandlingPriorityAndPrependedDate tests prepend with priority and date
// Ported from: t1400-prepend.sh "prepend handling priority and prepended date on add"
func TestPrependHandlingPriorityAndPrependedDate(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add task with -pt flags", func(t *testing.T) {
		output, code := env.RunCommand("-pt", "add", "(A) new task")
		expectedOutput := "2 (A) 2009-02-13 new task\nTODO: 2 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("prepend to prioritized dated task", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "this is just a")
		expectedOutput := "1 (A) 2009-02-13 this is just a new task"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestPrependWithPrependedDateKeepsBoth tests prepend preserves manual dates
// Ported from: t1400-prepend.sh "prepend with prepended date keeps both"
func TestPrependWithPrependedDateKeepsBoth(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add task with manual date", func(t *testing.T) {
		output, code := env.RunCommand("add", "2010-07-04 new task")
		expectedOutput := "3 2010-07-04 new task\nTODO: 3 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("prepend keeps both dates", func(t *testing.T) {
		output, code := env.RunCommand("prepend", "1", "this is just a")
		expectedOutput := "1 2010-07-04 this is just a new task"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
