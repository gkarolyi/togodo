package cmd

import (
	"testing"
)

func TestMove(t *testing.T) {
	sourceRepo, _ := setupEmptyTestRepository(t)
	destRepo, _ := setupEmptyTestRepository(t)

	sourceRepo.Add("task to move")
	sourceRepo.Add("task to keep")

	result, err := Move(sourceRepo, destRepo, 1)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	if result.Todo.Text != "task to move" {
		t.Errorf("Expected 'task to move', got '%s'", result.Todo.Text)
	}

	// Verify removed from source
	sourceTodos, _ := sourceRepo.ListAll()
	if len(sourceTodos) != 1 {
		t.Errorf("Expected 1 task in source, got %d", len(sourceTodos))
	}
	if sourceTodos[0].Text != "task to keep" {
		t.Errorf("Expected 'task to keep' remaining in source, got '%s'", sourceTodos[0].Text)
	}

	// Verify added to destination
	destTodos, _ := destRepo.ListAll()
	if len(destTodos) != 1 {
		t.Errorf("Expected 1 task in destination, got %d", len(destTodos))
	}
	if destTodos[0].Text != "task to move" {
		t.Errorf("Expected 'task to move' in destination, got '%s'", destTodos[0].Text)
	}
}

func TestMoveInvalidLineNumber(t *testing.T) {
	sourceRepo, _ := setupEmptyTestRepository(t)
	destRepo, _ := setupEmptyTestRepository(t)

	sourceRepo.Add("task one")

	// Try to move non-existent task 42
	_, err := Move(sourceRepo, destRepo, 42)
	if err == nil {
		t.Error("Expected error for invalid line number")
	}
}
