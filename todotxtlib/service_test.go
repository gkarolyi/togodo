package todotxtlib

import (
	"testing"
)

// TestService_AddTodos_SingleTask tests adding a single task
func TestService_AddTodos_SingleTask(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	todos, err := service.AddTodos([]string{"(A) task one"})

	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertTodoText(t, todos[0], "(A) task one")

	// Verify saved to repository
	output := buf.String()
	expectedOutput := "(A) task one\n"
	if output != expectedOutput {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

// TestService_AddTodos_MultipleTasks tests adding multiple tasks
func TestService_AddTodos_MultipleTasks(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	todos, err := service.AddTodos([]string{
		"(A) task one",
		"(B) task two",
		"task three",
	})

	assertNoError(t, err)
	assertTodoCount(t, todos, 3)

	// Verify all saved
	output := buf.String()
	assertContains(t, output, "(A) task one")
	assertContains(t, output, "(B) task two")
	assertContains(t, output, "task three")
}

// TestService_AddTodos_AutomaticSorting tests that added todos are sorted
func TestService_AddTodos_AutomaticSorting(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add in "wrong" order
	_, err := service.AddTodos([]string{
		"(C) low priority",
		"(A) high priority",
		"(B) medium priority",
	})
	assertNoError(t, err)

	// Verify sorted correctly
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 3)
	assertTodoPriority(t, allTodos[0], "A")
	assertTodoPriority(t, allTodos[1], "B")
	assertTodoPriority(t, allTodos[2], "C")
}

// TestService_ToggleTodos_SingleTask tests toggling a single task
func TestService_ToggleTodos_SingleTask(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add a task
	service.AddTodos([]string{"(A) task one"})

	// Toggle it
	todos, err := service.ToggleTodos([]int{0})

	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertTodoCompleted(t, todos[0], true)

	// Verify in repository
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 1)
	assertTodoCompleted(t, allTodos[0], true)
}

// TestService_ToggleTodos_MultipleTasks tests toggling multiple tasks
func TestService_ToggleTodos_MultipleTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add tasks
	service.AddTodos([]string{
		"(A) task one",
		"(B) task two",
		"(C) task three",
	})

	// Toggle first and third (before sorting, they're at positions 0 and 2)
	todos, err := service.ToggleTodos([]int{0, 2})

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)
	assertTodoCompleted(t, todos[0], true)
	assertTodoCompleted(t, todos[1], true)

	// Verify in repository
	// After sorting, done tasks move to the end
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 3)

	// Count completed tasks
	completedCount := 0
	for _, todo := range allTodos {
		if todo.Done {
			completedCount++
		}
	}
	if completedCount != 2 {
		t.Errorf("Expected 2 completed tasks, got %d", completedCount)
	}
}

// TestService_ToggleTodos_ToggleBack tests toggling a task back to not done
func TestService_ToggleTodos_ToggleBack(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add and toggle a task
	service.AddTodos([]string{"(A) task one"})
	service.ToggleTodos([]int{0})

	// Toggle it back
	todos, err := service.ToggleTodos([]int{0})

	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertTodoCompleted(t, todos[0], false)
}

// TestService_ToggleTodos_InvalidIndex tests toggling with an invalid index
func TestService_ToggleTodos_InvalidIndex(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{"(A) task one"})

	// Try to toggle non-existent index
	_, err := service.ToggleTodos([]int{99})

	assertError(t, err)
}

// TestService_SetPriorities_SingleTask tests setting priority on a single task
func TestService_SetPriorities_SingleTask(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add a task without priority
	service.AddTodos([]string{"task one"})

	// Set priority
	todos, err := service.SetPriorities([]int{0}, "A")

	assertNoError(t, err)
	assertTodoCount(t, todos, 1)
	assertTodoPriority(t, todos[0], "A")

	// Verify in repository
	allTodos, _ := repo.ListAll()
	assertTodoPriority(t, allTodos[0], "A")
}

// TestService_SetPriorities_MultipleTasks tests setting priority on multiple tasks
func TestService_SetPriorities_MultipleTasks(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add tasks
	service.AddTodos([]string{
		"task one",
		"task two",
		"task three",
	})

	// Set priority on first and third
	todos, err := service.SetPriorities([]int{0, 2}, "B")

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)
	assertTodoPriority(t, todos[0], "B")
	assertTodoPriority(t, todos[1], "B")

	// Verify in repository
	allTodos, _ := repo.ListAll()
	assertTodoPriority(t, allTodos[0], "B")
	assertTodoPriority(t, allTodos[1], "") // No priority set
	assertTodoPriority(t, allTodos[2], "B")
}

// TestService_SetPriorities_NoSorting tests that SetPriorities doesn't sort
func TestService_SetPriorities_NoSorting(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add tasks - they will be sorted by AddTodos
	// After sorting: A, B, C (by priority)
	service.AddTodos([]string{
		"(C) task one",
		"(B) task two",
		"(A) task three",
	})

	// Get initial order after AddTodos sorting
	beforeTodos, _ := repo.ListAll()
	firstTaskText := beforeTodos[0].Text
	secondTaskText := beforeTodos[1].Text

	// Change priority of third task (C) to AAA
	// This would move it to first if we sorted
	service.SetPriorities([]int{2}, "AAA")

	// Verify order is preserved (no sorting after SetPriorities)
	allTodos, _ := repo.ListAll()
	assertTodoText(t, allTodos[0], firstTaskText)
	assertTodoText(t, allTodos[1], secondTaskText)
	// Third task should have new priority but stay in same position
	assertTodoPriority(t, allTodos[2], "AAA")
}

// TestService_SetPriorities_InvalidIndex tests setting priority with invalid index
func TestService_SetPriorities_InvalidIndex(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{"task one"})

	// Try to set priority on non-existent index
	_, err := service.SetPriorities([]int{99}, "A")

	assertError(t, err)
}

// TestService_RemoveDoneTodos_EmptyList tests removing done todos from empty list
func TestService_RemoveDoneTodos_EmptyList(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	todos, err := service.RemoveDoneTodos()

	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

// TestService_RemoveDoneTodos_NoDoneTodos tests removing when no todos are done
func TestService_RemoveDoneTodos_NoDoneTodos(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one",
		"(B) task two",
	})

	todos, err := service.RemoveDoneTodos()

	assertNoError(t, err)
	assertTodoCount(t, todos, 0)

	// Verify all tasks still present
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 2)
}

// TestService_RemoveDoneTodos_SomeDone tests removing some done todos
func TestService_RemoveDoneTodos_SomeDone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one",
		"(B) task two",
		"(C) task three",
	})

	// Mark first and third as done
	service.ToggleTodos([]int{0, 2})

	// Remove done todos
	todos, err := service.RemoveDoneTodos()

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// Verify only undone task remains
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 1)
	assertTodoText(t, allTodos[0], "(B) task two")
}

// TestService_RemoveDoneTodos_AllDone tests removing when all todos are done
func TestService_RemoveDoneTodos_AllDone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one",
		"(B) task two",
	})

	// Mark all as done
	service.ToggleTodos([]int{0, 1})

	// Remove done todos
	todos, err := service.RemoveDoneTodos()

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)

	// Verify all removed
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 0)
}

// TestService_SearchTodos_NoResults tests searching with no matches
func TestService_SearchTodos_NoResults(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one +project1",
		"(B) task two +project2",
	})

	todos, err := service.SearchTodos("+project3")

	assertNoError(t, err)
	assertTodoCount(t, todos, 0)
}

// TestService_SearchTodos_WithResults tests searching with matches
func TestService_SearchTodos_WithResults(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one +project1 @work",
		"(B) task two +project2 @home",
		"(C) task three +project1 @home",
	})

	// Search for project1
	todos, err := service.SearchTodos("+project1")

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)
	assertContains(t, todos[0].Text, "+project1")
	assertContains(t, todos[1].Text, "+project1")
}

// TestService_SearchTodos_ContextSearch tests searching by context
func TestService_SearchTodos_ContextSearch(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	service.AddTodos([]string{
		"(A) task one @work",
		"(B) task two @home",
		"(C) task three @work",
	})

	// Search for @work
	todos, err := service.SearchTodos("@work")

	assertNoError(t, err)
	assertTodoCount(t, todos, 2)
	assertContains(t, todos[0].Text, "@work")
	assertContains(t, todos[1].Text, "@work")
}

// TestService_Integration_FullWorkflow tests a complete workflow
func TestService_Integration_FullWorkflow(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := NewTodoService(repo)

	// Add some tasks
	addedTodos, err := service.AddTodos([]string{
		"(A) task one +project @work",
		"(B) task two +project",
		"task three",
	})
	assertNoError(t, err)
	assertTodoCount(t, addedTodos, 3)

	// Toggle one task to done
	toggledTodos, err := service.ToggleTodos([]int{1})
	assertNoError(t, err)
	assertTodoCount(t, toggledTodos, 1)
	assertTodoCompleted(t, toggledTodos[0], true)

	// Set priority on a task
	priTodos, err := service.SetPriorities([]int{2}, "C")
	assertNoError(t, err)
	assertTodoCount(t, priTodos, 1)
	assertTodoPriority(t, priTodos[0], "C")

	// Search for project
	searchTodos, err := service.SearchTodos("+project")
	assertNoError(t, err)
	assertTodoCount(t, searchTodos, 2)

	// Remove done todos
	removedTodos, err := service.RemoveDoneTodos()
	assertNoError(t, err)
	assertTodoCount(t, removedTodos, 1)

	// Verify final state
	allTodos, _ := repo.ListAll()
	assertTodoCount(t, allTodos, 2)
}
