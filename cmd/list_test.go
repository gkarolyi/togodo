package cmd

import (
	"bytes"
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestExecuteList_AllTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test listing all tasks with empty search query
	output, err := executeListForTest(repo, "")
	assertNoError(t, err)

	expected := `  1 (A) test todo 1 +project2 @context1
  2 (B) test todo 2 +project1 @context2
  3 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterByContext(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering by context
	output, err := executeListForTest(repo, "@context1")
	assertNoError(t, err)

	expected := `  1 (A) test todo 1 +project2 @context1
  2 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterByProject(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering by project
	output, err := executeListForTest(repo, "+project1")
	assertNoError(t, err)

	expected := `  1 (B) test todo 2 +project1 @context2
  2 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterByPriority(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering by priority
	output, err := executeListForTest(repo, "(A)")
	assertNoError(t, err)

	expected := `  1 (A) test todo 1 +project2 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterByKeyword(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering by keyword
	output, err := executeListForTest(repo, "todo")
	assertNoError(t, err)

	expected := `  1 (A) test todo 1 +project2 @context1
  2 (B) test todo 2 +project1 @context2
  3 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_NoMatches(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering with query that matches nothing
	output, err := executeListForTest(repo, "nonexistent")
	assertNoError(t, err)

	expected := ``

	if output != expected {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestExecuteList_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test listing from empty repository
	output, err := executeListForTest(repo, "")
	assertNoError(t, err)

	expected := ``

	if output != expected {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestExecuteList_MultipleFilters(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering with multiple terms
	output, err := executeListForTest(repo, "test todo")
	assertNoError(t, err)

	expected := `  1 (A) test todo 1 +project2 @context1
  2 (B) test todo 2 +project1 @context2
  3 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_CaseSensitiveSearch(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks with different cases
	repo.Add("Task with UPPERCASE")
	repo.Add("task with lowercase")

	// Test case sensitive search
	output, err := executeListForTest(repo, "UPPERCASE")
	assertNoError(t, err)

	expected := `  1 Task with UPPERCASE`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterDoneTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filtering for done tasks
	output, err := executeListForTest(repo, "x ")
	assertNoError(t, err)

	expected := `  1 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterSpecialCharacters(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add task with special characters
	repo.Add("Email user@domain.com about +project due:2024-12-31")

	// Test filtering by email
	output, err := executeListForTest(repo, "user@domain.com")
	assertNoError(t, err)

	expected := `  1 Email user@domain.com about +project due:2024-12-31`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_FilterByDueDate(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks with due dates
	repo.Add("Task 1 due:2024-12-31")
	repo.Add("Task 2 due:2024-01-01")
	repo.Add("Task 3 no due date")

	// Test filtering by due date
	output, err := executeListForTest(repo, "due:2024")
	assertNoError(t, err)

	expected := `  1 Task 1 due:2024-12-31
  2 Task 2 due:2024-01-01`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_WhitespaceInFilter(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test filter with leading/trailing whitespace
	output, err := executeListForTest(repo, "  test todo  ")
	assertNoError(t, err)

	// Should find no matches since the exact string "  test todo  " doesn't exist
	expected := ``

	if output != expected {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestExecuteList_QuotedFilter(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add task with exact phrase
	repo.Add("test todo with exact phrase")
	repo.Add("test different todo phrase")

	// Test filtering for exact phrase
	output, err := executeListForTest(repo, "exact phrase")
	assertNoError(t, err)

	expected := `  1 test todo with exact phrase`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_PriorityOrdering(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks with different priorities
	repo.Add("(C) low priority")
	repo.Add("(A) high priority")
	repo.Add("(B) medium priority")
	repo.Add("no priority")

	// Test listing all tasks to verify ordering
	output, err := executeListForTest(repo, "")
	assertNoError(t, err)

	expected := `  1 (C) low priority
  2 (A) high priority
  3 (B) medium priority
  4 no priority`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_Integration(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Test full integration with complex filter
	output, err := executeListForTest(repo, "+project1")
	assertNoError(t, err)

	expected := `  1 (B) test todo 2 +project1 @context2
  2 x (C) test todo 3 +project1 @context1`

	if output != expected {
		t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestExecuteList_ErrorHandling(t *testing.T) {
	// Note: This test depends on the Repository.Search implementation
	// If it can return errors, we should test that path
	repo, _ := setupTestRepository(t)

	// Test with various potentially problematic inputs
	problematicInputs := []string{
		"",   // empty string
		" ",  // just space
		"\n", // newline
		"\t", // tab
		"@",  // incomplete context
		"+",  // incomplete project
		"(",  // incomplete priority
	}

	for _, input := range problematicInputs {
		_, err := executeListForTest(repo, input)
		assertNoError(t, err) // Should not error, just return filtered results
	}
}

func TestList(t *testing.T) {
	t.Run("lists all tasks", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := List(repo, "")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Todos) != 2 {
			t.Errorf("expected 2 todos, got %d", len(result.Todos))
		}
		if result.TotalCount != 2 {
			t.Errorf("expected total count 2, got %d", result.TotalCount)
		}
		if result.ShownCount != 2 {
			t.Errorf("expected shown count 2, got %d", result.ShownCount)
		}
	})

	t.Run("filters tasks", func(t *testing.T) {
		// Setup
		var buf bytes.Buffer
		buf.WriteString("task one\ntask two\n")
		reader := todotxtlib.NewBufferReader(&buf)
		writer := todotxtlib.NewBufferWriter(&buf)
		repo, _ := todotxtlib.NewFileRepository(reader, writer)

		// Execute
		result, err := List(repo, "one")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(result.Todos))
		}
		if result.TotalCount != 2 {
			t.Errorf("expected total count 2, got %d", result.TotalCount)
		}
		if result.ShownCount != 1 {
			t.Errorf("expected shown count 1, got %d", result.ShownCount)
		}
	})
}
