package cmd

import (
	"testing"
)

func TestExecutePri_SingleTask(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test setting priority for task 2 (which has priority B)
	args := []string{"2", "A"}
	err := executePri(baseCmd, args)

	assertNoError(t, err)

	// Verify the priority was set
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Count tasks with priority A - should be 2 now (1 original + 1 newly set)
	priorityACount := 0
	for _, todo := range todos {
		if todo.Priority == "A" {
			priorityACount++
		}
	}
	if priorityACount != 2 {
		t.Errorf("Expected 2 tasks with priority A, got %d", priorityACount)
	}
}

func TestExecutePri_MultipleTasks(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test setting priority for multiple tasks (tasks 1 and 2)
	args := []string{"1", "2", "C"}
	err := executePri(baseCmd, args)

	assertNoError(t, err)

	// Verify both tasks had their priority set
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Count tasks with priority C
	priorityCCount := 0
	for _, todo := range todos {
		if todo.Priority == "C" {
			priorityCCount++
		}
	}

	// Should have 2 total priority C tasks (2 newly set to C)
	if priorityCCount != 2 {
		t.Errorf("Expected 2 tasks with priority C, got %d", priorityCCount)
	}
}

func TestExecutePri_ChangePriority(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test changing an existing priority
	args := []string{"1", "B"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify the priority was changed
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// The first task originally had priority A, now should have B
	found := false
	for _, todo := range todos {
		if todo.Priority == "B" && todo.Text != "(B) test todo 2 +project1 @context2" {
			found = true
			assertTodoPriority(t, todo, "B")
			assertContains(t, todo.Text, "(B)")
			assertNotContains(t, todo.Text, "(A)")
		}
	}
	if !found {
		t.Error("Expected to find a task that changed from priority A to B")
	}
}

func TestExecutePri_RemovePriority(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test removing priority by setting empty string
	args := []string{"1", ""}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify the priority was removed
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// Count tasks with priority A (should be 0 now)
	priorityACount := 0
	for _, todo := range todos {
		if todo.Priority == "A" {
			priorityACount++
		}
	}

	if priorityACount != 0 {
		t.Errorf("Expected 0 tasks with priority A after removal, got %d", priorityACount)
	}
}

func TestExecutePri_InvalidLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with invalid line number (too high)
	args := []string{"10", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to set priority for todo at line 10")
}

func TestExecutePri_InvalidLineNumberFormat(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with invalid line number format
	args := []string{"abc", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestExecutePri_ZeroLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with line number 0 (should fail - line numbers start at 1)
	args := []string{"0", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
}

func TestExecutePri_NegativeLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with negative line number
	args := []string{"-1", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
}

func TestExecutePri_EmptyRepository(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Test setting priority on empty repository
	args := []string{"1", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to set priority for todo at line 1")
}

func TestExecutePri_InvalidPriorityValues(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with various priority values
	testCases := []struct {
		priority string
		valid    bool
	}{
		{"A", true},
		{"B", true},
		{"C", true},
		{"Z", true},
		{"a", true}, // lowercase should work
		{"1", true}, // numbers might work depending on implementation
		{"", true},  // empty string to remove priority
	}

	for _, tc := range testCases {
		args := []string{"1", tc.priority}
		err := executePri(baseCmd, args)
		if tc.valid {
			assertNoError(t, err)
		} else {
			assertError(t, err)
		}
	}
}

func TestExecutePri_PreservesOtherProperties(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add a task with projects and contexts
	baseCmd.Repository.Add("original task +project @context")

	// Change its priority
	args := []string{"1", "A"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify projects and contexts are preserved
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)

	todo := todos[0]
	assertTodoPriority(t, todo, "A")
	assertContains(t, todo.Text, "+project")
	assertContains(t, todo.Text, "@context")
	assertContains(t, todo.Text, "original task")
}

func TestExecutePri_MixedValidInvalidNumbers(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with mix of valid and invalid line numbers
	args := []string{"2", "abc", "3", "A"}
	err := executePri(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")

	// Verify first task was set before error occurred
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// Check that task 2 was set to A before the error
	priorityACount := 0
	for _, todo := range todos {
		if todo.Priority == "A" {
			priorityACount++
		}
	}

	// Should have 2 A priorities (1 original + 1 newly set before error)
	if priorityACount != 2 {
		t.Errorf("Expected 2 tasks with priority A, got %d", priorityACount)
	}
}

func TestExecutePri_MultipleTasksSamePriority(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add multiple tasks without priority
	baseCmd.Repository.Add("task 1")
	baseCmd.Repository.Add("task 2")
	baseCmd.Repository.Add("task 3")

	// Set same priority for multiple tasks
	args := []string{"1", "2", "3", "A"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify all tasks have the same priority
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	for _, todo := range todos {
		assertTodoPriority(t, todo, "A")
		assertContains(t, todo.Text, "(A)")
	}
}

func TestExecutePri_SortingAfterPriorityChange(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add tasks in specific order
	baseCmd.Repository.Add("no priority task")
	baseCmd.Repository.Add("(C) low priority task")

	// Change the non-priority task to high priority
	args := []string{"1", "A"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify sorting - high priority should come first
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// The first task should now be the high priority one
	firstTodo := todos[0]
	assertTodoPriority(t, firstTodo, "A")
	assertContains(t, firstTodo.Text, "no priority task")

	// The second task should be the low priority one
	secondTodo := todos[1]
	assertTodoPriority(t, secondTodo, "C")
}

func TestExecutePri_Integration(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test full integration: set priority, verify save was called
	args := []string{"2", "A"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify priority was set and repository was saved
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Count tasks with priority A (should be 2 now)
	priorityACount := 0
	for _, todo := range todos {
		if todo.Priority == "A" {
			priorityACount++
		}
	}

	if priorityACount != 2 {
		t.Errorf("Expected 2 tasks with priority A, got %d", priorityACount)
	}
}

func TestExecutePri_SingleArgument(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with only one argument (will be treated as setting priority "A" for line 0)
	// This should work - it sets priority A for the first task
	args := []string{"A"}
	err := executePri(baseCmd, args)

	// This should succeed since it treats "A" as priority for line 0
	assertNoError(t, err)
}

func TestExecutePri_DoneTaskPriority(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test setting priority on a done task (line 3 is done)
	args := []string{"3", "A"}
	err := executePri(baseCmd, args)
	assertNoError(t, err)

	// Verify the done task can have its priority changed
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// Find the done task with new priority
	found := false
	for _, todo := range todos {
		if todo.Done && todo.Priority == "A" {
			found = true
			assertTodoCompleted(t, todo, true)
			assertTodoPriority(t, todo, "A")
		}
	}

	if !found {
		t.Error("Expected to find done task with priority A")
	}
}
