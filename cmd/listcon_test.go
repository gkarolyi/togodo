package cmd

import (
	"testing"
)

func TestListcon(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos with contexts
	addTodos(t, repo, []string{
		"task @home @work",
		"another @home",
		"no context",
	})

	result, err := Listcon(repo)
	if err != nil {
		t.Fatalf("Listcon failed: %v", err)
	}

	// Should have 2 contexts
	if len(result.Contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(result.Contexts))
	}

	// Should include @home and @work
	hasHome := false
	hasWork := false
	for _, ctx := range result.Contexts {
		if ctx == "@home" {
			hasHome = true
		}
		if ctx == "@work" {
			hasWork = true
		}
	}
	if !hasHome || !hasWork {
		t.Errorf("Expected @home and @work contexts, got %v", result.Contexts)
	}
}
