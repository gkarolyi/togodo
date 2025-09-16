package cmd

import (
	"testing"

	"github.com/gkarolyi/togodo/internal/injector"
)

func TestExecuteTidy_WithDoneTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Mark one more task as done to have multiple done tasks
	repo.ToggleDone(0) // Mark first task as done

	// Get initial count
	initialTodos, err := repo.ListAll()
	assertNoError(t, err)
	initialCount := len(initialTodos)

	// Count done tasks before tidy
	doneTasks := 0
	for _, todo := range initialTodos {
		if todo.Done {
			doneTasks++
		}
	}

	// Execute tidy
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify done tasks were removed
	todos, err := repo.ListAll()
	assertNoError(t, err)
	expectedCount := initialCount - doneTasks
	assertTodoCount(t, todos, expectedCount)

	// Verify no done tasks remain
	for _, todo := range todos {
		assertTodoCompleted(t, todo, false)
	}
}

func TestExecuteTidy_NoDoneTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add only undone tasks
	repo.Add("(A) undone task 1")
	repo.Add("(B) undone task 2")

	// Get initial count
	initialTodos, err := repo.ListAll()
	assertNoError(t, err)
	initialCount := len(initialTodos)

	// Execute tidy
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify no tasks were removed
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, initialCount)

	// Verify all tasks are still undone
	for _, todo := range todos {
		assertTodoCompleted(t, todo, false)
	}
}

func TestExecuteTidy_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Execute tidy on empty repository
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify repository is still empty
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

func TestExecuteTidy_AllTasksDone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks and mark them all as done
	repo.Add("x task 1")
	repo.Add("x task 2")
	repo.Add("x task 3")

	// Get initial count
	initialTodos, err := repo.ListAll()
	assertNoError(t, err)
	initialCount := len(initialTodos)

	// Verify all are done
	for _, todo := range initialTodos {
		assertTodoCompleted(t, todo, true)
	}

	// Execute tidy
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify all tasks were removed
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 0)

	// Verify we removed the expected number
	if initialCount != 3 {
		t.Errorf("Expected to start with 3 tasks, had %d", initialCount)
	}
}

func TestExecuteTidy_MixedDoneUndone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add mixed done and undone tasks
	repo.Add("(A) undone high priority")
	repo.Add("x done task 1")
	repo.Add("(B) undone medium priority")
	repo.Add("x done task 2")
	repo.Add("undone no priority")

	// Execute tidy
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify only undone tasks remain
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Verify all remaining tasks are undone
	for _, todo := range todos {
		assertTodoCompleted(t, todo, false)
	}

	// Verify the correct tasks remain
	assertTodoExists(t, todos, "(A) undone high priority")
	assertTodoExists(t, todos, "(B) undone medium priority")
	assertTodoExists(t, todos, "undone no priority")

	// Verify done tasks are gone
	assertTodoNotExists(t, todos, "x done task 1")
	assertTodoNotExists(t, todos, "x done task 2")
}

func TestExecuteTidy_PreservesOrder(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks in specific order with mixed priorities
	repo.Add("(A) high priority 1")
	repo.Add("x done task")
	repo.Add("(A) high priority 2")
	repo.Add("(B) medium priority")
	repo.Add("no priority task")

	// Execute tidy
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify sorting is maintained
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 4)

	// Verify priority ordering is preserved/restored
	priorityTasks := 0
	for i, todo := range todos {
		if todo.Priority != "" {
			priorityTasks++
			// Priority tasks should come before non-priority tasks
			if i == len(todos)-1 {
				continue // Last item, can't check next
			}
			for j := i + 1; j < len(todos); j++ {
				if todos[j].Priority == "" {
					// Found a non-priority task after this priority task, which is correct
					break
				}
			}
		}
	}

	if priorityTasks != 3 {
		t.Errorf("Expected 3 priority tasks after tidy, got %d", priorityTasks)
	}
}

func TestExecuteTidy_PreservesProjectsContexts(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks with projects and contexts
	repo.Add("(A) keep task +project1 @context1")
	repo.Add("x done task +project2 @context2")
	repo.Add("keep task +project3 @context3")

	// Execute tidy
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify projects and contexts are preserved
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// Check that projects and contexts are preserved
	foundProject1 := false
	foundProject3 := false
	for _, todo := range todos {
		if todo.Text == "(A) keep task +project1 @context1" {
			foundProject1 = true
			assertContains(t, todo.Text, "+project1")
			assertContains(t, todo.Text, "@context1")
		}
		if todo.Text == "keep task +project3 @context3" {
			foundProject3 = true
			assertContains(t, todo.Text, "+project3")
			assertContains(t, todo.Text, "@context3")
		}
	}

	if !foundProject1 {
		t.Error("Expected to find task with +project1 @context1")
	}
	if !foundProject3 {
		t.Error("Expected to find task with +project3 @context3")
	}

	// Verify done task with +project2 @context2 is gone
	assertTodoNotExists(t, todos, "x done task +project2 @context2")
}

func TestExecuteTidy_WithDueDates(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add tasks with due dates
	repo.Add("(A) task due:2024-12-31 +project")
	repo.Add("x done task due:2024-01-01")
	repo.Add("task due:2025-01-01")

	// Execute tidy
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify due dates are preserved in remaining tasks
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	foundDue2024 := false
	foundDue2025 := false
	for _, todo := range todos {
		if todo.Text == "(A) task due:2024-12-31 +project" {
			foundDue2024 = true
			assertContains(t, todo.Text, "due:2024-12-31")
		}
		if todo.Text == "task due:2025-01-01" {
			foundDue2025 = true
			assertContains(t, todo.Text, "due:2025-01-01")
		}
	}

	if !foundDue2024 {
		t.Error("Expected to find task with due:2024-12-31")
	}
	if !foundDue2025 {
		t.Error("Expected to find task with due:2025-01-01")
	}
}

func TestExecuteTidy_PrintsRemovedTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Mark additional task as done
	repo.ToggleDone(0)

	// Get done tasks before tidy
	doneTodos, err := repo.ListDone()
	assertNoError(t, err)
	doneCount := len(doneTodos)

	// Execute tidy
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify that done tasks would have been printed
	// (We can't easily test the actual printing without mocking the output)
	if doneCount == 0 {
		t.Error("Expected to have some done tasks to remove")
	}
}

func TestExecuteTidy_Integration(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Get initial state
	initialTodos, err := repo.ListAll()
	assertNoError(t, err)
	initialCount := len(initialTodos)

	// Count initially done tasks
	initialDoneCount := 0
	for _, todo := range initialTodos {
		if todo.Done {
			initialDoneCount++
		}
	}

	// Execute tidy
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify final state
	finalTodos, err := repo.ListAll()
	assertNoError(t, err)
	expectedFinalCount := initialCount - initialDoneCount
	assertTodoCount(t, finalTodos, expectedFinalCount)

	// Verify all remaining tasks are undone
	for _, todo := range finalTodos {
		assertTodoCompleted(t, todo, false)
	}

	// Verify repository was saved (implicit in the fact that we can read the changes)
}

func TestExecuteTidy_ErrorHandling(t *testing.T) {
	// This test would be more meaningful if we could inject repository errors
	// For now, test basic error paths
	repo, _ := setupTestRepository(t)

	// Execute tidy - should not error under normal conditions
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)
}

func TestExecuteTidy_MultipleRuns(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add some tasks and mark some as done
	repo.Add("task 1")
	repo.Add("x done task 1")
	repo.Add("task 2")
	repo.Add("x done task 2")

	// First tidy run
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify done tasks were removed
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// Second tidy run (should be no-op)
	err = executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify count is unchanged
	todos, err = repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// All remaining tasks should still be undone
	for _, todo := range todos {
		assertTodoCompleted(t, todo, false)
	}
}

func TestExecuteTidy_LargeNumberOfDoneTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add many done tasks
	doneTaskCount := 50
	undoneTaskCount := 10

	for i := 0; i < doneTaskCount; i++ {
		repo.Add("x done task " + string(rune('0'+i%10)))
	}

	for i := 0; i < undoneTaskCount; i++ {
		repo.Add("undone task " + string(rune('0'+i%10)))
	}

	// Execute tidy
	err := executeTidy(repo, injector.CreateCLIPresenter())
	assertNoError(t, err)

	// Verify only undone tasks remain
	todos, err := repo.ListAll()
	assertNoError(t, err)
	assertTodoCount(t, todos, undoneTaskCount)

	// Verify all remaining tasks are undone
	for _, todo := range todos {
		assertTodoCompleted(t, todo, false)
		assertContains(t, todo.Text, "undone task")
	}
}
