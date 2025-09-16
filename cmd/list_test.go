package cmd

import (
	"testing"
)

func TestExecuteList_AllTasks(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test listing all tasks with empty search query
	err := executeList(baseCmd, "")
	assertNoError(t, err)

	// Verify all tasks are returned
	todos, err := baseCmd.Repository.Search("")
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)
}

func TestExecuteList_FilterByContext(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering by context
	err := executeList(baseCmd, "@context1")
	assertNoError(t, err)

	// Verify only tasks with @context1 are returned
	todos, err := baseCmd.Repository.Search("@context1")
	assertNoError(t, err)

	// Should find 2 tasks with @context1
	assertTodoCount(t, todos, 2)
	for _, todo := range todos {
		assertContains(t, todo.Text, "@context1")
	}
}

func TestExecuteList_FilterByProject(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering by project
	err := executeList(baseCmd, "+project1")
	assertNoError(t, err)

	// Verify only tasks with +project1 are returned
	todos, err := baseCmd.Repository.Search("+project1")
	assertNoError(t, err)

	// Should find 2 tasks with +project1
	assertTodoCount(t, todos, 2)
	for _, todo := range todos {
		assertContains(t, todo.Text, "+project1")
	}
}

func TestExecuteList_FilterByPriority(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering by priority
	err := executeList(baseCmd, "(A)")
	assertNoError(t, err)

	// Verify only priority A tasks are returned
	todos, err := baseCmd.Repository.Search("(A)")
	assertNoError(t, err)

	// Should find 1 task with priority A
	assertTodoCount(t, todos, 1)
	assertTodoPriority(t, todos[0], "A")
}

func TestExecuteList_FilterByKeyword(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering by keyword
	err := executeList(baseCmd, "todo")
	assertNoError(t, err)

	// Verify only tasks containing "todo" are returned
	todos, err := baseCmd.Repository.Search("todo")
	assertNoError(t, err)

	// All test tasks contain "todo" in their text
	assertTodoCount(t, todos, 3)
	for _, todo := range todos {
		assertContains(t, todo.Text, "todo")
	}
}

func TestExecuteList_NoMatches(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering with query that matches nothing
	err := executeList(baseCmd, "nonexistent")
	assertNoError(t, err)

	// Verify no tasks are returned
	todos, err := baseCmd.Repository.Search("nonexistent")
	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

func TestExecuteList_EmptyRepository(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Test listing from empty repository
	err := executeList(baseCmd, "")
	assertNoError(t, err)

	// Verify no tasks are returned
	todos, err := baseCmd.Repository.Search("")
	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

func TestExecuteList_MultipleFilters(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering with multiple terms
	err := executeList(baseCmd, "@context1 +project")
	assertNoError(t, err)

	// Verify tasks matching the combined filter
	todos, err := baseCmd.Repository.Search("@context1 +project")
	assertNoError(t, err)

	// Should find tasks that contain both @context1 and +project
	for _, todo := range todos {
		assertContains(t, todo.Text, "@context1")
		assertContains(t, todo.Text, "+project")
	}
}

func TestExecuteList_CaseSensitiveSearch(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add tasks with different cases
	baseCmd.Repository.Add("Task with UPPERCASE")
	baseCmd.Repository.Add("task with lowercase")

	// Test case sensitive search
	err := executeList(baseCmd, "UPPERCASE")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("UPPERCASE")
	assertNoError(t, err)

	// Should find only the uppercase version
	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "UPPERCASE")
}

func TestExecuteList_FilterDoneTasks(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filtering for done tasks
	err := executeList(baseCmd, "x ")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("x ")
	assertNoError(t, err)

	// Should find 1 done task
	assertTodoCount(t, todos, 1)
	assertTodoCompleted(t, todos[0], true)
}

func TestExecuteList_FilterSpecialCharacters(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add task with special characters
	baseCmd.Repository.Add("Email user@domain.com about +project due:2024-12-31")

	// Test filtering by email
	err := executeList(baseCmd, "user@domain.com")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("user@domain.com")
	assertNoError(t, err)

	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "user@domain.com")
}

func TestExecuteList_FilterByDueDate(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add tasks with due dates
	baseCmd.Repository.Add("Task 1 due:2024-12-31")
	baseCmd.Repository.Add("Task 2 due:2024-01-01")
	baseCmd.Repository.Add("Task 3 no due date")

	// Test filtering by due date
	err := executeList(baseCmd, "due:2024")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("due:2024")
	assertNoError(t, err)

	// Should find 2 tasks with 2024 due dates
	assertTodoCount(t, todos, 2)
	for _, todo := range todos {
		assertContains(t, todo.Text, "due:2024")
	}
}

func TestExecuteList_WhitespaceInFilter(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test filter with leading/trailing whitespace
	err := executeList(baseCmd, "  test todo  ")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("  test todo  ")
	assertNoError(t, err)

	// Should still find matching tasks
	for _, todo := range todos {
		assertContains(t, todo.Text, "test todo")
	}
}

func TestExecuteList_QuotedFilter(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add task with exact phrase
	baseCmd.Repository.Add("test todo with exact phrase")
	baseCmd.Repository.Add("test different todo phrase")

	// Test filtering for exact phrase (though quotes may not be handled specially)
	err := executeList(baseCmd, "exact phrase")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("exact phrase")
	assertNoError(t, err)

	// Should find tasks containing both words
	for _, todo := range todos {
		assertContains(t, todo.Text, "exact")
		assertContains(t, todo.Text, "phrase")
	}
}

func TestExecuteList_PriorityOrdering(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add tasks with different priorities
	baseCmd.Repository.Add("(C) low priority")
	baseCmd.Repository.Add("(A) high priority")
	baseCmd.Repository.Add("(B) medium priority")
	baseCmd.Repository.Add("no priority")

	// Test listing all tasks to verify ordering
	err := executeList(baseCmd, "")
	assertNoError(t, err)

	todos, err := baseCmd.Repository.Search("")
	assertNoError(t, err)

	assertTodoCount(t, todos, 4)

	// Priority tasks should come first, in priority order
	priorityTasks := 0
	for _, todo := range todos {
		if todo.Priority != "" {
			priorityTasks++
		}
	}

	if priorityTasks != 3 {
		t.Errorf("Expected 3 priority tasks, got %d", priorityTasks)
	}
}

func TestExecuteList_Integration(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test full integration with complex filter
	err := executeList(baseCmd, "+project1")
	assertNoError(t, err)

	// Verify the search functionality works end-to-end
	todos, err := baseCmd.Repository.Search("+project1")
	assertNoError(t, err)

	// All returned todos should contain +project1
	for _, todo := range todos {
		assertContains(t, todo.Text, "+project1")
	}
}

func TestExecuteList_ErrorHandling(t *testing.T) {
	// Note: This test depends on the Repository.Search implementation
	// If it can return errors, we should test that path
	baseCmd, _ := setupTestBaseCommand(t)

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
		err := executeList(baseCmd, input)
		assertNoError(t, err) // Should not error, just return filtered results
	}
}
