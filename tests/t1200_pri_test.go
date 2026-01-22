package tests

import (
	"testing"
)

// TestPriUsage tests pri command usage errors
// Ported from: t1200-pri.sh "priority usage"
func TestPriUsage(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("pri with invalid args", func(t *testing.T) {
		output, code := env.RunCommand("pri", "B", "B")
		expectedCode := 1
		expectedOutput := `usage: togodo pri NR PRIORITY [NR PRIORITY ...]
note: PRIORITY must be anywhere from A to Z.`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicPriority tests basic priority setting
// Ported from: t1200-pri.sh "basic priority"
func TestBasicPriority(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("list before priority", func(t *testing.T) {
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

	t.Run("set priority B on task 1", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "B")
		expectedOutput := "1 (B) smell the uppercase Roses +flowers @outside\nTODO: 1 prioritized (B)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows priority with color", func(t *testing.T) {
		output, code := env.RunCommand("list")
		// Upstream expects ANSI color codes: [0;32m...[0m
		expectedOutput := "\033[0;32m1 (B) smell the uppercase Roses +flowers @outside\033[0m\n2 notice the sunflowers\n3 stop\n--\nTODO: 3 of 3 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p flag shows plain", func(t *testing.T) {
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

	t.Run("set priority C on task 2", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "C")
		expectedOutput := "2 (C) notice the sunflowers\nTODO: 2 prioritized (C)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p after second priority", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (C) notice the sunflowers
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

	t.Run("add new task", func(t *testing.T) {
		output, code := env.RunCommand("add", "smell the coffee +wakeup")
		expectedOutput := "4 smell the coffee +wakeup\nTODO: 4 added."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p shows priorities first", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
2 (C) notice the sunflowers
4 smell the coffee +wakeup
3 stop
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

// TestPriorityError tests error handling for invalid task numbers
// Ported from: t1200-pri.sh "priority error"
func TestPriorityError(t *testing.T) {
	env := SetupTestEnv(t)

	t.Run("pri with non-existent task", func(t *testing.T) {
		output, code := env.RunCommand("pri", "10", "B")
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

// TestReprioritize tests changing priority on already-prioritized task
// Ported from: t1200-pri.sh "reprioritize"
func TestReprioritize(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +flowers @outside
(C) notice the sunflowers
stop`)

	t.Run("change priority from C to A", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "A")
		expectedOutput := "2 (A) notice the sunflowers\nTODO: 2 re-prioritized from (C) to (A)."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list with -p shows new priority order", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `2 (A) notice the sunflowers
1 (B) smell the uppercase Roses +flowers @outside
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

	t.Run("set same priority should error", func(t *testing.T) {
		output, code := env.RunCommand("pri", "2", "a")
		expectedCode := 1
		expectedOutput := "2 (A) notice the sunflowers\nTODO: 2 already prioritized (A)."
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list after error attempt", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `2 (A) notice the sunflowers
1 (B) smell the uppercase Roses +flowers @outside
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

// TestMultiplePriority tests prioritizing multiple tasks in one command
// Ported from: t1200-pri.sh "multiple priority"
func TestMultiplePriority(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	t.Run("pri with multiple tasks", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "A", "2", "B")
		expectedOutput := `1 (A) smell the uppercase Roses +flowers @outside
TODO: 1 prioritized (A).
2 (B) notice the sunflowers
TODO: 2 prioritized (B).`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestMultipleReprioritize tests re-prioritizing multiple tasks in one command
// Ported from: t1200-pri.sh "multiple reprioritize"
func TestMultipleReprioritize(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	// Set initial priorities
	env.RunCommand("pri", "1", "A", "2", "B")

	t.Run("reprioritize multiple tasks", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "Z", "2", "X")
		expectedOutput := `1 (Z) smell the uppercase Roses +flowers @outside
TODO: 1 re-prioritized from (A) to (Z).
2 (X) notice the sunflowers
TODO: 2 re-prioritized from (B) to (X).`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestMultiplePrioritizeError tests partial success when prioritizing multiple tasks
// Ported from: t1200-pri.sh "multiple prioritize error"
func TestMultiplePrioritizeError(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`smell the uppercase Roses +flowers @outside
notice the sunflowers
stop`)

	// Set initial priorities
	env.RunCommand("pri", "1", "A", "2", "B")
	// Change to Z and X
	env.RunCommand("pri", "1", "Z", "2", "X")

	t.Run("pri with mix of valid and invalid tasks (first)", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "B", "4", "B")
		expectedCode := 1
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
TODO: 1 re-prioritized from (Z) to (B).
TODO: No task 4.`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("pri with mix of valid and invalid tasks (second)", func(t *testing.T) {
		output, code := env.RunCommand("pri", "1", "C", "4", "B", "3", "A")
		expectedCode := 1
		expectedOutput := `1 (C) smell the uppercase Roses +flowers @outside
TODO: 1 re-prioritized from (B) to (C).
TODO: No task 4.`
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
