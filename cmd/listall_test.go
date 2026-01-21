package cmd

import (
	"testing"
)

func TestListall(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos
	repo.Add("active task")
	repo.Add("another active")

	// Mark one as done
	_, err := Do(repo, []int{0})
	if err != nil {
		t.Fatalf("Failed to mark todo as done: %v", err)
	}

	// List all
	result, err := Listall(repo)
	if err != nil {
		t.Fatalf("Listall failed: %v", err)
	}

	// Should have 2 todos (including done)
	if len(result.Todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(result.Todos))
	}

	// Verify one is done
	doneCount := 0
	for _, todo := range result.Todos {
		if todo.Done {
			doneCount++
		}
	}
	if doneCount != 1 {
		t.Errorf("Expected 1 done todo, got %d", doneCount)
	}
}
