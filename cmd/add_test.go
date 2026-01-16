package cmd

import (
	"strings"
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

	// Check that tasks are in correct sorted order
	output, err := repo.WriteToString()
	assertNoError(t, err)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedOrder := []string{
		"(A) first task +project1 @context1",
		"(B) second task +project2 @context2",
		"third task without priority",
	}

	if len(lines) != len(expectedOrder) {
		t.Fatalf("Expected %d lines, got %d", len(expectedOrder), len(lines))
	}

	for i, expected := range expectedOrder {
		if lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i+1, expected, lines[i])
		}
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

	// Check that the new high priority task is present
	output, err := repo.WriteToString()
	assertNoError(t, err)

	if !strings.Contains(output, "(A) new high priority task") {
		t.Error("New high priority task not found")
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

	// Verify the actual order using output
	output, err := repo.WriteToString()
	assertNoError(t, err)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedOrder := []string{
		"(A) high priority task",
		"(B) medium priority task",
		"(C) low priority task",
		"no priority task",
	}

	if len(lines) != len(expectedOrder) {
		t.Fatalf("Expected %d lines, got %d", len(expectedOrder), len(lines))
	}

	for i, expected := range expectedOrder {
		if lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i+1, expected, lines[i])
		}
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

	// Count duplicates by filtering
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
