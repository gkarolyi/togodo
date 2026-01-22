package tests

import (
	"testing"
)

// TestListprojNoProjects tests listproj with no projects in file
// Ported from: t1320-listproj.sh "listproj no projects"
func TestListprojNoProjects(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`item 1
item 2
item 3`)

	t.Run("listproj with no projects", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		expectedOutput := ""
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected empty output\nGot:\n%s", output)
		}
	})
}

// TestListprojSingle tests listing projects from tasks with single projects
// Ported from: t1320-listproj.sh "Single project per line"
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
// Ported from: t1320-listproj.sh "Multi-project per line"
func TestListprojMultiple(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`+prj01 -- Some project 1 task
+prj02 -- Some project 2 task
+prj02 +prj03 -- Multi-project task`)

	t.Run("listproj with duplicates", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		// Should list unique projects, sorted
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

// TestListprojEmbeddedPlus tests that embedded + signs don't trigger project detection
// Ported from: t1320-listproj.sh "listproj embedded + test"
func TestListprojEmbeddedPlus(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`+prj01 -- Some project 1 task
+prj02 -- Some project 2 task
+prj02 ginatrapani+todo@gmail.com -- Some project 2 task`)

	t.Run("listproj ignores embedded plus in email", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		// "ginatrapani+todo@gmail.com" has + but should NOT be detected as project
		// +prj02 appears twice but should only be listed once
		expectedOutput := `+prj01
+prj02`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestBasicListproj tests basic project listing
// Ported from: t1320-listproj.sh "basic listproj"
func TestBasicListproj(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +roses @outside +shared
(C) notice the sunflowers +sunflowers @garden +shared +landscape
stop`)

	t.Run("basic listproj", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		expectedOutput := `+landscape
+roses
+shared
+sunflowers`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojWithContext tests filtering projects by context
// Ported from: t1320-listproj.sh "listproj with context"
func TestListprojWithContext(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +roses @outside +shared
(C) notice the sunflowers +sunflowers @garden +shared +landscape
stop`)

	t.Run("listproj filtered by context", func(t *testing.T) {
		output, code := env.RunCommand("listproj", "@garden")
		expectedOutput := `+landscape
+shared
+sunflowers`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojVariousPositions tests projects at different positions in tasks
// Ported from: t1320-listproj.sh "listproj of projects at various positions"
func TestListprojVariousPositions(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`+roses +flowers
+roses and only +flowers
(B) +roses prioritized +flowers
(B) 2024-02-21 +roses prioritized +flowers
x 2024-02-21 2024-02-19 +roses done +flowers
+roses +flowers at the front
  +roses +flowers at the front with leading space
my +flowers are +roses in the middle
at the back pick the +roses +flowers
at the back with trailing space +flowers +roses `)

	t.Run("listproj at various positions", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		expectedOutput := `+flowers
+roses`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojDefault tests default project detection patterns
// Ported from: t1320-listproj.sh "listproj with default configuration"
func TestListprojDefault(t *testing.T) {
	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) give a +1 to this project
(C) notice the sunflowers +sunflowers [+gardening] [+landscape]
stop`)

	t.Run("listproj with default config", func(t *testing.T) {
		output, code := env.RunCommand("listproj")
		expectedOutput := `+1
+sunflowers`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojLimitingAlphabetic tests limiting projects to alphabetic characters
// Ported from: t1320-listproj.sh "listproj limiting to alphabetic characters"
func TestListprojLimitingAlphabetic(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SIGIL_VALID_PATTERN configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) give a +1 to this project
(C) notice the sunflowers +sunflowers [+gardening] [+landscape]
stop`)

	t.Run("listproj with alphabetic pattern", func(t *testing.T) {
		// TODO: Set TODOTXT_SIGIL_VALID_PATTERN='[a-zA-Z]\{1,\}'
		output, code := env.RunCommand("listproj")
		expectedOutput := "+sunflowers"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojAllowingBrackets tests allowing brackets around projects
// Ported from: t1320-listproj.sh "listproj allowing brackets around projects"
func TestListprojAllowingBrackets(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SIGIL_BEFORE_PATTERN and TODOTXT_SIGIL_AFTER_PATTERN configuration")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) give a +1 to this project
(C) notice the sunflowers +sunflowers [+gardening] [+landscape]
stop`)

	t.Run("listproj with brackets patterns", func(t *testing.T) {
		// TODO: Set TODOTXT_SIGIL_BEFORE_PATTERN='\[\{0,1\}' TODOTXT_SIGIL_AFTER_PATTERN='\]\{0,1\}'
		output, code := env.RunCommand("listproj")
		expectedOutput := `+1
+gardening
+landscape
+sunflowers`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojWithContextSpecialCases tests context filtering with special options
// Ported from: t1320-listproj.sh "listproj with context special cases"
func TestListprojWithContextSpecialCases(t *testing.T) {
	t.Skip("TODO: Implement -+ and -d flags for custom todo file")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`(B) smell the uppercase Roses +roses @outside +shared
(C) notice the sunflowers +sunflowers @garden +shared +landscape
stop`)

	t.Run("listproj with -+ -d flags", func(t *testing.T) {
		// TODO: Support -+ and -d flags
		output, code := env.RunCommand("listproj", "@garden")
		expectedOutput := `+landscape
+shared
+sunflowers`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojFromDoneTasks tests listing projects from done.txt
// Ported from: t1320-listproj.sh "listproj from done tasks"
func TestListprojFromDoneTasks(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SOURCEVAR configuration for done.txt")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`+prj01 -- Some project 1 task`)
	// TODO: Write to done.txt file

	t.Run("listproj from done file", func(t *testing.T) {
		// TODO: Set TODOTXT_SOURCEVAR=$DONE_FILE
		output, code := env.RunCommand("listproj")
		expectedOutput := `+done01
+done02`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojFromDoneTasksWithFiltering tests listing projects from done.txt with filter
// Ported from: t1320-listproj.sh "listproj from done tasks with filtering"
func TestListprojFromDoneTasksWithFiltering(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SOURCEVAR configuration for done.txt")

	env := SetupTestEnv(t)
	// TODO: Write to done.txt file

	t.Run("listproj from done file with filter", func(t *testing.T) {
		// TODO: Set TODOTXT_SOURCEVAR=$DONE_FILE
		output, code := env.RunCommand("listproj", "Special")
		expectedOutput := "+done01"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojCombinedSources tests listing projects from both todo.txt and done.txt
// Ported from: t1320-listproj.sh "listproj from combined open + done tasks"
func TestListprojCombinedSources(t *testing.T) {
	t.Skip("TODO: Implement TODOTXT_SOURCEVAR configuration for multiple sources")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`+prj01 -- Some project 1 task`)
	// TODO: Write to done.txt file

	t.Run("listproj from combined sources", func(t *testing.T) {
		// TODO: Set TODOTXT_SOURCEVAR='("$TODO_FILE" "$DONE_FILE")'
		output, code := env.RunCommand("listproj")
		expectedOutput := `+done01
+done02
+prj01`
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestListprojWithGrepOptionsDisruption tests that GREP_OPTIONS doesn't affect output
// Ported from: t1320-listproj.sh "listproj with GREP_OPTIONS disruption"
func TestListprojWithGrepOptionsDisruption(t *testing.T) {
	t.Skip("TODO: Verify GREP_OPTIONS doesn't affect output")

	env := SetupTestEnv(t)
	env.WriteTodoFile(`+prj01 -- Some project 1 task`)

	t.Run("listproj with GREP_OPTIONS=-n", func(t *testing.T) {
		// TODO: Set GREP_OPTIONS=-n environment variable
		output, code := env.RunCommand("listproj")
		expectedOutput := "+prj01"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
