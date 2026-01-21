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
	result, err := Del(repo, 0)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Should have deleted task 1
	if result.DeletedTodo.Text != "task 1" {
		t.Errorf("Expected deleted 'task 1', got '%s'", result.DeletedTodo.Text)
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
