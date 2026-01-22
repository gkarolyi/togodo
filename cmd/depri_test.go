package cmd

import (
	"testing"
)

func TestDepri(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add prioritized todo
	repo.Add("(A) high priority task")

	// Remove priority
	result, err := Depri(repo, 0)
	if err != nil {
		t.Fatalf("Depri failed: %v", err)
	}

	// Should have no priority
	if result.Todo.Priority != "" {
		t.Errorf("Expected no priority, got '%s'", result.Todo.Priority)
	}

	// Text should not have (A) prefix
	if result.Todo.Text == "(A) high priority task" {
		t.Errorf("Expected priority removed from text, got '%s'", result.Todo.Text)
	}
}
