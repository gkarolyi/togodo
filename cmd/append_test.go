package cmd

import (
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add initial todo
	repo.Add("notice the daisies")

	// Append text
	result, err := Append(repo, 0, "smell the roses")
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Should have appended text
	expected := "notice the daisies smell the roses"
	if result.Todo.Text != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.Todo.Text)
	}

	// Verify it's saved
	todos, _ := repo.ListAll()
	if !strings.Contains(todos[0].Text, "smell the roses") {
		t.Errorf("Expected appended text in repo, got '%s'", todos[0].Text)
	}
}
