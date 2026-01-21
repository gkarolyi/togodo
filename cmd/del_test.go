package cmd

import (
	"testing"
)

func TestDel(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos
	repo.Add("task 1")
	repo.Add("task 2")

	// Delete first todo
	result, err := Del(repo, []int{0})
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Should have deleted task 1
	if len(result.DeletedTodos) != 1 {
		t.Fatalf("Expected 1 deleted todo, got %d", len(result.DeletedTodos))
	}
	if result.DeletedTodos[0].Text != "task 1" {
		t.Errorf("Expected deleted 'task 1', got '%s'", result.DeletedTodos[0].Text)
	}

	// Should only have 1 todo left
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo remaining, got %d", len(todos))
	}

	// Remaining should be task 2
	if todos[0].Text != "task 2" {
		t.Errorf("Expected remaining 'task 2', got '%s'", todos[0].Text)
	}
}

func TestDelMultiple(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos
	repo.Add("task 1")
	repo.Add("task 2")
	repo.Add("task 3")
	repo.Add("task 4")

	// Delete multiple todos (indices 3 and 1, i.e., task 4 and task 2)
	result, err := Del(repo, []int{3, 1})
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Should have deleted 2 todos
	if len(result.DeletedTodos) != 2 {
		t.Fatalf("Expected 2 deleted todos, got %d", len(result.DeletedTodos))
	}

	// Should have 2 todos left
	todos, _ := repo.ListAll()
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos remaining, got %d", len(todos))
	}

	// Remaining should be task 1 and task 3 (renumbered)
	if todos[0].Text != "task 1" || todos[1].Text != "task 3" {
		t.Errorf("Expected remaining ['task 1', 'task 3'], got ['%s', '%s']", todos[0].Text, todos[1].Text)
	}

	// Check line numbers are sequential
	if todos[0].LineNumber != 1 || todos[1].LineNumber != 2 {
		t.Errorf("Expected line numbers [1, 2], got [%d, %d]", todos[0].LineNumber, todos[1].LineNumber)
	}
}
