package cmd

import (
	"testing"
)

func TestAddMultiple(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add multiple tasks from text with newlines
	text := "task one\ntask two\ntask three"
	result, err := AddMultiple(repo, text, false)
	if err != nil {
		t.Fatalf("AddMultiple failed: %v", err)
	}

	// Should have added 3 tasks
	if len(result.Todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(result.Todos))
	}

	// Verify tasks in repository
	todos, _ := repo.ListAll()
	if len(todos) != 3 {
		t.Errorf("Expected 3 todos in repository, got %d", len(todos))
		for i, todo := range todos {
			t.Logf("Todo %d: %s", i, todo.Text)
		}
	}

	// Check task texts
	texts := make(map[string]bool)
	for _, todo := range todos {
		texts[todo.Text] = true
	}

	if !texts["task one"] {
		t.Error("Expected 'task one' in repository")
	}
	if !texts["task two"] {
		t.Error("Expected 'task two' in repository")
	}
	if !texts["task three"] {
		t.Error("Expected 'task three' in repository")
	}
}
