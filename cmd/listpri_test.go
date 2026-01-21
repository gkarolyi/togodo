package cmd

import (
	"testing"
)

func TestListpri(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos with different priorities
	addTodos(t, repo, []string{
		"(A) high priority",
		"(B) medium priority",
		"(A) another high",
		"no priority",
	})

	// List priority A
	result, err := Listpri(repo, "A")
	if err != nil {
		t.Fatalf("Listpri failed: %v", err)
	}

	// Should have 2 A-priority todos
	if len(result.Todos) != 2 {
		t.Errorf("Expected 2 A-priority todos, got %d", len(result.Todos))
	}

	// Verify they're the right ones (order may vary after sort)
	foundHigh := false
	foundAnother := false
	for _, todo := range result.Todos {
		if todo.Text == "(A) high priority" {
			foundHigh = true
		}
		if todo.Text == "(A) another high" {
			foundAnother = true
		}
	}
	if !foundHigh {
		t.Errorf("Expected to find '(A) high priority' in results")
	}
	if !foundAnother {
		t.Errorf("Expected to find '(A) another high' in results")
	}
}

func TestListpriAllPriorities(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos with different priorities
	addTodos(t, repo, []string{
		"(A) high priority",
		"(B) medium priority",
		"(C) low priority",
		"no priority",
	})

	// List all prioritized todos (empty priority string)
	result, err := Listpri(repo, "")
	if err != nil {
		t.Fatalf("Listpri failed: %v", err)
	}

	// Should have 3 prioritized todos (A, B, C)
	if len(result.Todos) != 3 {
		t.Errorf("Expected 3 prioritized todos, got %d", len(result.Todos))
	}
}

func TestListpriExcludesDone(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos including done ones
	addTodos(t, repo, []string{
		"(A) high priority",
		"(B) medium priority",
	})

	// Mark one as done
	toggleTodos(t, repo, []int{0})

	// List priority A
	result, err := Listpri(repo, "A")
	if err != nil {
		t.Fatalf("Listpri failed: %v", err)
	}

	// Should have 0 A-priority todos (the one was marked as done)
	if len(result.Todos) != 0 {
		t.Errorf("Expected 0 A-priority todos (done excluded), got %d", len(result.Todos))
	}
}
