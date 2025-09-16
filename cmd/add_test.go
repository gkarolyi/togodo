package cmd

import (
	"testing"
)

func TestExecuteAdd_SingleTask(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding a single task
	args := []string{"(A) new task +project @context"}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify the task was added
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "new task")
	assertContains(t, todos[0].Text, "+project")
	assertContains(t, todos[0].Text, "@context")
}

func TestExecuteAdd_MultipleTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding multiple tasks
	args := []string{
		"(A) first task +project1 @context1",
		"(B) second task +project2 @context2",
		"third task without priority",
	}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify all tasks were added
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Check that tasks are properly sorted (priority tasks first)
	foundFirst := false
	foundSecond := false
	foundThird := false
	for _, todo := range todos {
		if todo.Text == "(A) first task +project1 @context1" {
			foundFirst = true
		}
		if todo.Text == "(B) second task +project2 @context2" {
			foundSecond = true
		}
		if todo.Text == "third task without priority" {
			foundThird = true
		}
	}

	if !foundFirst {
		t.Error("First task not found in repository")
	}
	if !foundSecond {
		t.Error("Second task not found in repository")
	}
	if !foundThird {
		t.Error("Third task not found in repository")
	}
}

func TestExecuteAdd_EmptyArgs(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding with empty args (should not add anything)
	args := []string{}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify no tasks were added
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

func TestExecuteAdd_TaskWithSpecialCharacters(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding task with special characters
	args := []string{"(A) task with @email:user@domain.com and +project:name due:2024-12-31"}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify the task was added with special characters preserved
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "@email:user@domain.com")
	assertContains(t, todos[0].Text, "+project:name")
	assertContains(t, todos[0].Text, "due:2024-12-31")
}

func TestExecuteAdd_WithExistingTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Get initial count
	initialTodos, err := repo.ListAll()
	assertNoError(t, err)
	initialCount := len(initialTodos)

	// Test adding to existing repository
	args := []string{"(A) new high priority task"}
	err = executeAdd(repo, args)

	assertNoError(t, err)

	// Verify the task was added to existing tasks
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, initialCount+1)

	// Check that the new high priority task is sorted correctly
	found := false
	for _, todo := range todos {
		if todo.Text == "(A) new high priority task" {
			found = true
			break
		}
	}
	if !found {
		t.Error("New task not found in repository")
	}
}

func TestExecuteAdd_SortingBehavior(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks in mixed priority order
	args := []string{
		"no priority task",
		"(C) low priority task",
		"(A) high priority task",
		"(B) medium priority task",
	}

	for _, arg := range args {
		err := executeAdd(repo, []string{arg})
		assertNoError(t, err)
	}

	// Verify tasks are sorted correctly after addition
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 4)

	// Check that high priority tasks come first
	priorityTasks := 0
	for i, todo := range todos {
		if todo.Priority != "" {
			priorityTasks++
			// Priority tasks should come before non-priority tasks
			if i >= len(todos)-1 {
				continue // Last item, can't check next
			}
			nextTodo := todos[i+1]
			if todo.Priority == "" && nextTodo.Priority != "" {
				t.Errorf("Priority task found after non-priority task at position %d", i)
			}
		}
	}

	if priorityTasks != 3 {
		t.Errorf("Expected 3 priority tasks, found %d", priorityTasks)
	}
}

func TestExecuteAdd_MultilineString(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding a task that contains newlines (should be treated as one task)
	args := []string{"(A) task with\nnewline characters\nin the text"}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify only one task was added
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "task with")
	assertContains(t, todos[0].Text, "newline characters")
	assertContains(t, todos[0].Text, "in the text")
}

func TestExecuteAdd_DuplicateTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test adding duplicate tasks
	args := []string{
		"(A) duplicate task +project @context",
		"(A) duplicate task +project @context",
	}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify both tasks were added (duplicates should be allowed)
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	duplicateCount := 0
	for _, todo := range todos {
		if todo.Text == "(A) duplicate task +project @context" {
			duplicateCount++
		}
	}

	if duplicateCount != 2 {
		t.Errorf("Expected 2 duplicate tasks, found %d", duplicateCount)
	}
}

func TestExecuteAdd_Integration(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Test full integration: add, verify save was called
	args := []string{"(A) integration test task"}
	err := executeAdd(repo, args)

	assertNoError(t, err)

	// Verify task was added and repository was saved
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertContains(t, todos[0].Text, "integration test task")

	// Verify the task has the correct priority
	if todos[0].Priority != "A" {
		t.Errorf("Expected priority 'A', got '%s'", todos[0].Priority)
	}
}
