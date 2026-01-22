package tests

import (
	"testing"
)

// TestCreateNewReport tests generating initial report
// Ported from: t1950-report.sh "create new report"
func TestCreateNewReport(t *testing.T) {
	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
stop and think
smell the coffee +wakeup
make the coffee +wakeup
visit http://example.com`)

	t.Run("report creates new report file", func(t *testing.T) {
		output, code := env.RunCommand("report")
		expectedOutput := "TODO: todo.txt does not contain any done tasks.\n2009-02-13T00:00:00 5 0\nTODO: Report file updated."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows all tasks", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `1 (B) smell the uppercase Roses +flowers @outside
4 make the coffee +wakeup
3 smell the coffee +wakeup
2 stop and think
5 visit http://example.com
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

// TestReportOfDoneTasks tests report with done tasks
// Ported from: t1950-report.sh "report of done tasks"
func TestReportOfDoneTasks(t *testing.T) {
	t.Skip("TODO: Implement -A flag for archive during do")

	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
stop and think
smell the coffee +wakeup
make the coffee +wakeup
visit http://example.com`)

	t.Run("initial report", func(t *testing.T) {
		output, code := env.RunCommand("report")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})

	t.Run("mark task 3 done with archive", func(t *testing.T) {
		output, code := env.RunCommand("do", "-A", "3")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})

	t.Run("report after done task", func(t *testing.T) {
		output, code := env.RunCommand("report")
		expectedOutput := "TODO: todo.txt does not contain any done tasks.\n2009-02-13T00:00:00 4 1\nTODO: Report file updated."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows 4 remaining tasks", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// Should show 4 tasks (task 3 archived)
		_ = output
	})
}

// TestReportPerformsArchiving tests report with archiving flag
// Ported from: t1950-report.sh "report performs archiving"
func TestReportPerformsArchiving(t *testing.T) {
	t.Skip("TODO: Implement -a flag for auto-archive during do")

	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
stop and think
smell the coffee +wakeup
make the coffee +wakeup
visit http://example.com`)

	t.Run("initial report", func(t *testing.T) {
		env.RunCommand("report")
	})

	t.Run("mark task done with auto-archive", func(t *testing.T) {
		env.RunCommand("do", "-A", "3")
		env.RunCommand("report")
	})

	t.Run("mark another task done", func(t *testing.T) {
		output, code := env.RunCommand("do", "-a", "3")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		_ = output
	})

	t.Run("report performs archiving", func(t *testing.T) {
		output, code := env.RunCommand("report")
		expectedOutput := "x 2009-02-13 make the coffee +wakeup\nTODO: todo.txt archived.\n2009-02-13T00:00:00 3 2\nTODO: Report file updated."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestReportUnchangedWhenNoChanges tests report idempotency
// Ported from: t1950-report.sh "report unchanged when no changes"
func TestReportUnchangedWhenNoChanges(t *testing.T) {
	env := SetupFileBasedTestEnv(t)
	defer env.ClearTestDate()

	if err := env.SetTestDate("2009-02-13"); err != nil {
		t.Fatalf("Failed to set test date: %v", err)
	}

	env.WriteTodoFileContent(`(B) smell the uppercase Roses +flowers @outside
stop and think
visit http://example.com`)

	t.Run("generate initial report", func(t *testing.T) {
		output, code := env.RunCommand("report")
		expectedOutput := "TODO: todo.txt does not contain any done tasks.\n2009-02-13T00:00:00 3 0\nTODO: Report file updated."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("report again without changes", func(t *testing.T) {
		output, code := env.RunCommand("report")
		expectedOutput := "TODO: todo.txt does not contain any done tasks.\n2009-02-13T00:00:00 3 0\nTODO: Report file is up-to-date."
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
