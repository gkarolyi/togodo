package cmd

import (
	"testing"
)

func TestReplace(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add initial todo
	repo.Add("notice the daisies")

	// Replace it
	result, err := Replace(repo, 0, "smell the cows")
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	// Verify old todo
	if result.OldTodo.Text != "notice the daisies" {
		t.Errorf("Expected old todo 'notice the daisies', got '%s'", result.OldTodo.Text)
	}

	// Verify new todo
	if result.NewTodo.Text != "smell the cows" {
		t.Errorf("Expected new todo 'smell the cows', got '%s'", result.NewTodo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].Text != "smell the cows" {
		t.Errorf("Expected 'smell the cows' in repo, got '%s'", todos[0].Text)
	}
}
