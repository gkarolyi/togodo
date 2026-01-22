package tests

import (
	"testing"
)

// TestDeduplicatePreserveLineNumbers tests deduplicate preserves line numbers
// Ported from: t1910-deduplicate.sh "deduplicate and preserve line numbers"
func TestDeduplicatePreserveLineNumbers(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`duplicated
two
x done
duplicated
double task
double task
three`)

	t.Run("deduplicate removes 2 duplicates", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate")
		expectedOutput := "TODO: 2 duplicate task(s) removed"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows 5 tasks with original line numbers", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `5 double task
1 duplicated
7 three
2 two
3 x done
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

// TestDeduplicateWithoutDuplicates tests deduplicate when no duplicates exist
// Ported from: t1910-deduplicate.sh "deduplicate without duplicates"
func TestDeduplicateWithoutDuplicates(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`one
two
three`)

	t.Run("deduplicate with no duplicates", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate")
		expectedCode := 1
		expectedOutput := "TODO: No duplicate tasks found"
		if code != expectedCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}

// TestDeduplicateAndDeleteLines tests deduplicate with -n flag to renumber
// Ported from: t1910-deduplicate.sh "deduplicate and delete lines"
func TestDeduplicateAndDeleteLines(t *testing.T) {
	t.Skip("TODO: Implement -n flag for renumbering tasks")

	env := SetupTestEnv(t)

	env.WriteTodoFile(`duplicated
two
x done
duplicated
double task
two
three`)

	t.Run("deduplicate with -n flag renumbers", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate", "-n")
		expectedOutput := "TODO: 2 duplicate task(s) removed"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows 5 tasks with sequential line numbers", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `4 double task
1 duplicated
5 three
2 two
3 x done
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

// TestDeduplicateMoreThanTwo tests deduplicate with multiple duplicates
// Ported from: t1910-deduplicate.sh "deduplicate more than two occurrences"
func TestDeduplicateMoreThanTwo(t *testing.T) {
	env := SetupTestEnv(t)

	env.WriteTodoFile(`one
duplicated
three
duplicated
duplicated
six
duplicated`)

	t.Run("deduplicate removes 3 occurrences", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate")
		expectedOutput := "TODO: 3 duplicate task(s) removed"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows 4 tasks with first occurrence kept", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := `2 duplicated
1 one
6 six
3 three
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

// TestDeduplicateWithNonPrintable tests deduplicate with ANSI formatting codes
// Ported from: t1910-deduplicate.sh "deduplicate with non-printable duplicates"
func TestDeduplicateWithNonPrintable(t *testing.T) {
	env := SetupTestEnv(t)

	// Tasks with ANSI bold formatting: \033[1m for bold, \033[0m for reset
	env.WriteTodoFile("normal task\na \033[1mbold\033[0m task\nsomething else\na \033[1mbold\033[0m task\nsomething more")

	t.Run("deduplicate removes formatted duplicate", func(t *testing.T) {
		output, code := env.RunCommand("deduplicate")
		expectedOutput := "TODO: 1 duplicate task(s) removed"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})

	t.Run("list shows 4 tasks with formatting preserved", func(t *testing.T) {
		output, code := env.RunCommand("-p", "list")
		expectedOutput := "2 a \033[1mbold\033[0m task\n1 normal task\n3 something else\n5 something more\n--\nTODO: 4 of 4 tasks shown"
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		if output != expectedOutput {
			t.Errorf("Output mismatch\nExpected:\n%s\n\nGot:\n%s", expectedOutput, output)
		}
	})
}
