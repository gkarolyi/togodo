package tests

import (
	"testing"
)

// TestReplaceUsage tests replace command usage validation
// Ported from: t1100-replace.sh "replace usage"
func TestReplaceUsage(t *testing.T) {
	env := SetupTestEnv(t)
	env.RunCommand("add", "notice the daisies")

	t.Run("replace with wrong args", func(t *testing.T) {
		output, code := env.RunCommand("replace", "adf", "asdfa")
		expectedCode := 1
		expectedOutput := `usage: togodo replace NR "UPDATED ITEM"`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicReplace tests basic replace functionality
// Ported from: t1100-replace.sh "basic replace"
func TestBasicReplace(t *testing.T) {
	env := SetupTestEnv(t)
	env.RunCommand("add", "notice the daisies")

	t.Run("replace task text", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "smell the cows")
		expectedOutput := `1 notice the daisies
TODO: Replaced task with:
1 smell the cows`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after first replace", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 smell the cows
--
TODO: 1 of 1 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace task text again", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "smell the roses")
		expectedOutput := `1 smell the cows
TODO: Replaced task with:
1 smell the roses`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after second replace", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `1 smell the roses
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

// TestReplaceError tests error handling for non-existent tasks
// Ported from: t1100-replace.sh "replace error"
func TestReplaceError(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("smell the cows\ngrow some corn\nthrash some hay\nchase the chickens")

	t.Run("replace non-existent task", func(t *testing.T) {
		output, code := env.RunCommand("replace", "10", "hej!")
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

// TestReplaceInMultiItemFile tests replacing tasks in multi-item file
// Ported from: t1100-replace.sh "replace in multi-item file"
func TestReplaceInMultiItemFile(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("smell the cows\ngrow some corn\nthrash some hay\nchase the chickens")

	t.Run("replace first task", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "smell the cheese")
		expectedOutput := `1 smell the cows
TODO: Replaced task with:
1 smell the cheese`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace third task", func(t *testing.T) {
		output, code := env.RunCommand("replace", "3", "jump on hay")
		expectedOutput := `3 thrash some hay
TODO: Replaced task with:
3 jump on hay`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace fourth task", func(t *testing.T) {
		output, code := env.RunCommand("replace", "4", "collect the eggs")
		expectedOutput := `4 chase the chickens
TODO: Replaced task with:
4 collect the eggs`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithPriority tests replacing tasks while preserving priority
// Ported from: t1100-replace.sh "replace with priority"
func TestReplaceWithPriority(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) collect the eggs")

	t.Run("replace preserves priority", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "collect the bread")
		expectedOutput := `1 (A) collect the eggs
TODO: Replaced task with:
1 (A) collect the bread`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace again preserves priority", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "collect the eggs")
		expectedOutput := `1 (A) collect the bread
TODO: Replaced task with:
1 (A) collect the eggs`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithAmpersand tests replacing with ampersand character
// Ported from: t1100-replace.sh "replace with &"
func TestReplaceWithAmpersand(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("jump on hay")

	t.Run("replace with ampersand", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "thrash the hay & thrash the wheat")
		expectedOutput := `1 jump on hay
TODO: Replaced task with:
1 thrash the hay & thrash the wheat`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithSpaces tests replacing with multiple spaces
// Ported from: t1100-replace.sh "replace with spaces"
func TestReplaceWithSpaces(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("jump on hay")

	t.Run("replace with multiple spaces", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "notice the   three   spaces")
		expectedOutput := `1 jump on hay
TODO: Replaced task with:
1 notice the   three   spaces`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithSymbols tests replacing with special symbols
// Ported from: t1100-replace.sh "replace with symbols"
func TestReplaceWithSymbols(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("smell the cows\ngrow some corn\nthrash some hay\nchase the chickens")

	t.Run("replace with symbols", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "~@#$%^&*()-_=+[{]}|;:',<.>/?")
		expectedOutput := `1 smell the cows
TODO: Replaced task with:
1 ~@#$%^&*()-_=+[{]}|;:',<.>/?`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace with backtick and quotes", func(t *testing.T) {
		output, code := env.RunCommand("replace", "2", "`!\\\"")
		expectedOutput := `2 grow some corn
TODO: Replaced task with:
2 ` + "`!\\\"" + ``
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after symbol replacements", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `4 chase the chickens
3 thrash some hay
2 ` + "`!\\\"" + `
1 ~@#$%^&*()-_=+[{]}|;:',<.>/?
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

// TestReplaceHandlingPrependedDateOnAdd tests date preservation
// Ported from: t1100-replace.sh "replace handling prepended date on add"
func TestReplaceHandlingPrependedDateOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "new task")
		expectedOutput := "1 2009-02-13 new task\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace preserves date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "this is just a new one")
		expectedOutput := `1 2009-02-13 new task
TODO: Replaced task with:
1 2009-02-13 this is just a new one`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace with new date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "2010-07-04 this also has a new date")
		expectedOutput := `1 2009-02-13 this is just a new one
TODO: Replaced task with:
1 2010-07-04 this also has a new date`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceHandlingPrependedPriorityOnAdd tests adding priority to dated task
// Ported from: t1100-replace.sh "replace handling prepended priority on add"
func TestReplaceHandlingPrependedPriorityOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "new task")
		expectedOutput := "1 2009-02-13 new task\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace with priority", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "(B) this also has a priority now")
		expectedOutput := `1 2009-02-13 new task
TODO: Replaced task with:
1 (B) 2009-02-13 this also has a priority now`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceHandlingPriorityAndPrependedDateOnAdd tests replace preserving both
// Ported from: t1100-replace.sh "replace handling priority and prepended date on add"
func TestReplaceHandlingPriorityAndPrependedDateOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "new task")
		expectedOutput := "1 2009-02-13 new task\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("set priority", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "A")
		expectedOutput := "1 (A) 2009-02-13 new task\nTODO: 1 prioritized (A)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace preserves priority and date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "this is just a new one")
		expectedOutput := `1 (A) 2009-02-13 new task
TODO: Replaced task with:
1 (A) 2009-02-13 this is just a new one`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceHandlingPrependedPriorityAndDateOnAdd tests replacing both priority and date
// Ported from: t1100-replace.sh "replace handling prepended priority and date on add"
func TestReplaceHandlingPrependedPriorityAndDateOnAdd(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	t.Run("add with -t flag", func(t *testing.T) {
		output, code := env.RunCommand("-t", "add", "new task")
		expectedOutput := "1 2009-02-13 new task\nTODO: 1 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("replace with priority and new date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "(C) 2010-07-04 this also has a priority and new date")
		expectedOutput := `1 2009-02-13 new task
TODO: Replaced task with:
1 (C) 2010-07-04 this also has a priority and new date`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithPrependedDateReplacesExistingDate tests date replacement
// Ported from: t1100-replace.sh "replace with prepended date replaces existing date"
func TestReplaceWithPrependedDateReplacesExistingDate(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) 2009-02-13 this is just a new one")

	t.Run("replace date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "2010-07-04 this also has a new date")
		expectedOutput := `1 (A) 2009-02-13 this is just a new one
TODO: Replaced task with:
1 (A) 2010-07-04 this also has a new date`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithPrependedPriorityReplacesExistingPriority tests priority replacement
// Ported from: t1100-replace.sh "replace with prepended priority replaces existing priority"
func TestReplaceWithPrependedPriorityReplacesExistingPriority(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) 2009-02-13 this is just a new one")

	t.Run("replace priority", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "(B) this also has a new priority")
		expectedOutput := `1 (A) 2009-02-13 this is just a new one
TODO: Replaced task with:
1 (B) 2009-02-13 this also has a new priority`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithPrependedPriorityAndDateReplacesExistingDate tests both replacements (no existing priority)
// Ported from: t1100-replace.sh "replace with prepended priority and date replaces existing date"
func TestReplaceWithPrependedPriorityAndDateReplacesExistingDate(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("2009-02-13 this is just a new one")

	t.Run("replace with priority and date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "(B) 2010-07-04 this also has a new date")
		expectedOutput := `1 2009-02-13 this is just a new one
TODO: Replaced task with:
1 (B) 2010-07-04 this also has a new date`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReplaceWithPrependedPriorityAndDateReplacesExistingPriorityAndDate tests both replacements
// Ported from: t1100-replace.sh "replace with prepended priority and date replaces existing priority and date"
func TestReplaceWithPrependedPriorityAndDateReplacesExistingPriorityAndDate(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile("(A) 2009-02-13 this is just a new one")

	t.Run("replace both priority and date", func(t *testing.T) {
		output, code := env.RunCommand("replace", "1", "(B) 2010-07-04 this also has a new prio+date")
		expectedOutput := `1 (A) 2009-02-13 this is just a new one
TODO: Replaced task with:
1 (B) 2010-07-04 this also has a new prio+date`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
