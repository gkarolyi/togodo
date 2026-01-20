package tests

import (
	"testing"
)

// TestListprojSingle tests listing projects from tasks
// Ported from: t1320-listproj.sh
func TestListprojSingle(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`(A) +1 -- Some project 1 task, whitespace, one char
(A) +p2 -- Some project 2 task, whitespace, two char
+prj03 -- Some project 3 task, no whitespace
+prj04 -- Some project 4 task, no whitespace
+prj05+prj06 -- weird project`)

	t.Run("listproj shows all projects", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		expectedOutput := `+1
+p2
+prj03
+prj04
+prj05+prj06`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojMultiple tests listing projects when tasks have multiple projects
// Ported from: t1320-listproj.sh
func TestListprojMultiple(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`+prj01 -- Some project 1 task
+prj02 -- Some project 2 task
+prj02 +prj03 -- Multi-project task`)

	t.Run("listproj with duplicates", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		// Should list unique projects, possibly sorted
		// +prj02 appears twice but should only be listed once
		expectedOutput := `+prj01
+prj02
+prj03`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
