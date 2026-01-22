package tests

import (
	"testing"
)

// TestDoUsage tests do command usage errors
// Ported from: t1500-do.sh "do usage"
func TestDoUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("do with invalid args", func(t *testing.T) {
		output, code := env.RunCommand("do", "B", "B")
		expectedCode := 1
		expectedOutput := "usage: togodo do NR [NR ...]"
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDoMissingNR tests do command without arguments
// Ported from: t1500-do.sh "do missing NR"
func TestDoMissingNR(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("do without args", func(t *testing.T) {
		output, code := env.RunCommand("do")
		expectedCode := 1
		expectedOutput := "usage: togodo do NR [NR ...]"
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicDo tests basic do/mark done functionality with archiving
// Ported from: t1500-do.sh "basic do"
func TestBasicDo(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

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

	t.Run("do 7,6 (comma-separated)", func(t *testing.T) {
		output, code := env.RunCommand("do", "7,6")
		expectedOutput := `7 x 2009-02-13 remove4
TODO: 7 marked as done.
6 x 2009-02-13 remove3
TODO: 6 marked as done.
x 2009-02-13 remove3
x 2009-02-13 remove4
TODO: todo.txt archived.`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after first do", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 notice the sunflowers
4 remove1
5 remove2
1 smell the uppercase Roses +flowers @outside
3 stop
--
TODO: 5 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("do 5 4 (space-separated)", func(t *testing.T) {
		output, code := env.RunCommand("do", "5", "4")
		expectedOutput := `5 x 2009-02-13 remove2
TODO: 5 marked as done.
4 x 2009-02-13 remove1
TODO: 4 marked as done.
x 2009-02-13 remove1
x 2009-02-13 remove2
TODO: todo.txt archived.`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after second do", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 notice the sunflowers
1 smell the uppercase Roses +flowers @outside
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
}

// TestDoMultipleWithComma tests marking multiple tasks done with comma separator
// This is actually part of "basic do" but extracted for clarity
// Ported from: t1500-do.sh "basic do"
func TestDoMultipleWithComma(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFile(`task1
task2
task3
task4
task5`)

	t.Run("do 1,3,5 (comma-separated multiple)", func(t *testing.T) {
		output, code := env.RunCommand("do", "1,3,5")
		expectedOutput := `1 x 2009-02-13 task1
TODO: 1 marked as done.
3 x 2009-02-13 task3
TODO: 3 marked as done.
5 x 2009-02-13 task5
TODO: 5 marked as done.
x 2009-02-13 task1
x 2009-02-13 task3
x 2009-02-13 task5
TODO: todo.txt archived.`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after marking done", func(t *testing.T) {
		output, code := env.RunCommand("list")
		expectedOutput := `2 task2
4 task4
--
TODO: 2 of 5 tasks shown`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestFailMultipleDoAttempts tests error when trying to mark already-done task
// Ported from: t1500-do.sh "fail multiple do attempts"
func TestFailMultipleDoAttempts(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("mark task 3 done first time", func(t *testing.T) {
		output, code := env.RunCommand("do", "3")
		expectedOutput := `3 x 2009-02-13 stop
TODO: 3 marked as done.`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("mark task 3 done second time", func(t *testing.T) {
		output, code := env.RunCommand("do", "3")
		expectedCode := 1
		expectedOutput := "TODO: 3 is already marked done."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
