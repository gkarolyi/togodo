package cmd

import (
	"strings"
	"testing"
)

func TestPrepend(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add initial todo
	repo.Add("notice the sunflowers")

	// Prepend text
	result, err := Prepend(repo, 0, "really")
	if err != nil {
		t.Fatalf("Prepend failed: %v", err)
	}

	// Should have prepended text
	if !strings.Contains(result.Todo.Text, "really notice the sunflowers") {
		t.Errorf("Expected 'really notice the sunflowers', got '%s'", result.Todo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if !strings.Contains(todos[0].Text, "really notice the sunflowers") {
		t.Errorf("Expected prepended text in repo, got '%s'", todos[0].Text)
	}
}

func TestPrependPreservesPriority(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add prioritized todo
	repo.Add("(A) task with priority")

	// Prepend text
	result, err := Prepend(repo, 0, "important")
	if err != nil {
		t.Fatalf("Prepend failed: %v", err)
	}

	// Should preserve priority
	if result.Todo.Text != "(A) important task with priority" {
		t.Errorf("Expected '(A) important task with priority', got '%s'", result.Todo.Text)
	}
}
