package tests

import (
	"testing"
)

// TestReport tests generating task statistics
// Ported from: t1950-report.sh
func TestReport(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(A) task one
(B) task two
task three
x done task`)

	t.Run("report shows statistics", func(t *testing.T) {
		output, code := env.RunCommand("report")
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		// TODO: Should show counts by priority, done vs todo, etc.
		// Example format:
		// 2020-01-01 3 0 3
		// (date, total tasks, done tasks, todo tasks)
		_ = output
	})
}
