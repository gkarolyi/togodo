package cmd

import (
	"testing"

	"github.com/gkarolyi/togodo/todotxtlib"
)

func TestTidyCmd_WithDoneTasks(t *testing.T) {
	repo, buf := setupTestRepository(t)
	service := todotxtlib.NewTodoService(repo)

	// Mark one more task as done to have multiple done tasks
	repo.ToggleDone(0)

	// Execute tidy
	todos, err := service.RemoveDoneTodos()
	assertNoError(t, err)

	// Verify two todos were removed
	if len(todos) != 2 {
		t.Fatalf("Expected 2 removed todos, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "(B) test todo 2 +project1 @context2\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestTidyCmd_NoDoneTasks(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)
	service := todotxtlib.NewTodoService(repo)

	// Add only undone tasks
	repo.Add("task 1")
	repo.Add("task 2")

	// Execute tidy
	todos, err := service.RemoveDoneTodos()
	assertNoError(t, err)

	// Verify no todos were removed
	if len(todos) != 0 {
		t.Fatalf("Expected 0 removed todos, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := "task 1\n" +
		"task 2\n"

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestTidyCmd_EmptyRepository(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)
	service := todotxtlib.NewTodoService(repo)

	// Execute tidy on empty repository
	todos, err := service.RemoveDoneTodos()
	assertNoError(t, err)

	// Verify no todos were removed
	if len(todos) != 0 {
		t.Fatalf("Expected 0 removed todos, got %d", len(todos))
	}
}

func TestTidyCmd_AllTasksDone(t *testing.T) {
	repo, buf := setupEmptyTestRepository(t)
	service := todotxtlib.NewTodoService(repo)

	// Add tasks and mark them all as done
	repo.Add("x task 1")
	repo.Add("x task 2")
	repo.Add("x task 3")

	// Execute tidy
	todos, err := service.RemoveDoneTodos()
	assertNoError(t, err)

	// Verify all todos were removed
	if len(todos) != 3 {
		t.Fatalf("Expected 3 removed todos, got %d", len(todos))
	}

	output := getRepositoryString(t, repo, buf)

	expectedOutput := ""

	if output != expectedOutput {
		t.Errorf("Expected empty output, got:\n%s", output)
	}
}

func TestTidyCmd_ReturnsRemovedTasks(t *testing.T) {
	repo, _ := setupTestRepository(t)
	service := todotxtlib.NewTodoService(repo)

	// Mark additional task as done
	repo.ToggleDone(0)

	// Get done tasks before tidy
	doneFilter := todotxtlib.Filter{Done: "true"}
	doneTodos, err := repo.Filter(doneFilter)
	assertNoError(t, err)
	doneCount := len(doneTodos)

	// Execute tidy
	removedTodos, err := service.RemoveDoneTodos()
	assertNoError(t, err)

	// Verify that done tasks were returned
	if doneCount == 0 {
		t.Error("Expected to have some done tasks to remove")
	}

	if len(removedTodos) != doneCount {
		t.Errorf("Expected %d removed todos, got %d", doneCount, len(removedTodos))
	}
}
