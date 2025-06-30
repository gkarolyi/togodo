package cmd

import (
	"testing"
)

func TestExecuteDo_SingleTask(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Get initial todos to check state
	initialTodos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, initialTodos, 3)

	// Test toggling a task that's not done (line 1)
	args := []string{"1"}
	err = executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify the task was toggled
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// The done task should be at the end after sorting
	lastTodo := todos[len(todos)-1]
	assertTodoCompleted(t, lastTodo, true)
	assertContains(t, lastTodo.Text, "x ")
}

func TestExecuteDo_MultipleTask(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test toggling multiple tasks
	args := []string{"1", "2"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify both tasks were toggled
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Find the toggled tasks (they should be at the end after sorting)
	doneCount := 0
	for _, todo := range todos {
		if todo.Done {
			doneCount++
		}
	}

	// We should have 2 newly done tasks plus 1 that was already done = 3 total
	if doneCount != 3 {
		t.Errorf("Expected 3 done tasks, got %d", doneCount)
	}
}

func TestExecuteDo_ToggleAlreadyDone(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// The third task in our test data is already done
	// Test toggling it to undone (line 3)
	args := []string{"3"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify the task was toggled back to undone
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Count done tasks - should have 0 since we undid the only done task
	doneCount := 0
	for _, todo := range todos {
		if todo.Done {
			doneCount++
		}
	}

	if doneCount != 0 {
		t.Errorf("Expected 0 done tasks, got %d", doneCount)
	}
}

func TestExecuteDo_InvalidLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with invalid line number (too high)
	args := []string{"10"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to toggle todo at line 10")
}

func TestExecuteDo_InvalidLineNumberFormat(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with invalid line number format
	args := []string{"abc"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")
}

func TestExecuteDo_ZeroLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with line number 0 (should fail - line numbers start at 1)
	args := []string{"0"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
}

func TestExecuteDo_NegativeLineNumber(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with negative line number
	args := []string{"-1"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
}

func TestExecuteDo_EmptyRepository(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Test toggling on empty repository
	args := []string{"1"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to toggle todo at line 1")
}

func TestExecuteDo_SortingAfterToggle(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add some tasks with different priorities
	baseCmd.Repository.Add("(A) high priority task")
	baseCmd.Repository.Add("(B) medium priority task")
	baseCmd.Repository.Add("no priority task")

	// Toggle the high priority task to done
	args := []string{"1"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify sorting - done tasks should be at the end
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// The last task should be the done high priority task
	lastTodo := todos[len(todos)-1]
	assertTodoCompleted(t, lastTodo, true)
	assertContains(t, lastTodo.Text, "high priority task")

	// The first two should be undone
	assertTodoCompleted(t, todos[0], false)
	assertTodoCompleted(t, todos[1], false)
}

func TestExecuteDo_MultipleToggleSameTask(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Toggle the same task twice (should end up back to original state)
	args := []string{"1", "1"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Get the task that was toggled twice
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// The first task should be back to its original state (not done)
	assertTodoCompleted(t, todos[0], false)
	assertNotContains(t, todos[0].Text, "x ")
}

func TestExecuteDo_MixedValidInvalidNumbers(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test with mix of valid and invalid line numbers
	args := []string{"1", "abc", "2"}
	err := executeDo(baseCmd, args)
	assertError(t, err)
	assertContains(t, err.Error(), "failed to convert arg to int")

	// Verify first task was toggled before error occurred
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)

	// Check that first task was marked as done before the error
	doneCount := 0
	for _, todo := range todos {
		if todo.Done {
			doneCount++
		}
	}
	// Should be 2 (the original done task + the first task that was toggled)
	if doneCount != 2 {
		t.Errorf("Expected 2 done tasks, got %d", doneCount)
	}
}

func TestExecuteDo_Integration(t *testing.T) {
	baseCmd, _ := setupTestBaseCommand(t)

	// Test full integration: toggle, verify save was called
	args := []string{"1"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify task was toggled and repository was saved
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Verify the toggled task is now done and sorted to the end
	doneTask := todos[len(todos)-1]
	assertTodoCompleted(t, doneTask, true)
	assertContains(t, doneTask.Text, "x ")
}

func TestExecuteDo_PreservePriorityWhenToggling(t *testing.T) {
	baseCmd, _ := setupEmptyTestBaseCommand(t)

	// Add a priority task
	baseCmd.Repository.Add("(A) important task +project @context")

	// Toggle it to done
	args := []string{"1"}
	err := executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify the task is done but priority and tags are preserved
	todos, err := baseCmd.Repository.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)

	doneTask := todos[0]
	assertTodoCompleted(t, doneTask, true)
	assertTodoPriority(t, doneTask, "A")
	assertContains(t, doneTask.Text, "+project")
	assertContains(t, doneTask.Text, "@context")
	assertContains(t, doneTask.Text, "x ")

	// Toggle it back to undone
	err = executeDo(baseCmd, args)
	assertNoError(t, err)

	// Verify it's back to original state
	todos, err = baseCmd.Repository.ListAll()
	assertNoError(t, err)
	undoneTask := todos[0]
	assertTodoCompleted(t, undoneTask, false)
	assertTodoPriority(t, undoneTask, "A")
	assertContains(t, undoneTask.Text, "+project")
	assertContains(t, undoneTask.Text, "@context")
	assertNotContains(t, undoneTask.Text, "x ")
}
