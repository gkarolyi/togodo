package cmd

import (
	"testing"
)

func TestListproj(t *testing.T) {
	repo, _ := setupEmptyTestRepository(t)

	// Add todos with projects
	addTodos(t, repo, []string{
		"task +home +work",
		"another +home",
		"no project",
	})

	result, err := Listproj(repo)
	if err != nil {
		t.Fatalf("Listproj failed: %v", err)
	}

	// Should have 2 projects
	if len(result.Projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(result.Projects))
	}

	// Should include +home and +work
	hasHome := false
	hasWork := false
	for _, proj := range result.Projects {
		if proj == "+home" {
			hasHome = true
		}
		if proj == "+work" {
			hasWork = true
		}
	}
	if !hasHome || !hasWork {
		t.Errorf("Expected +home and +work projects, got %v", result.Projects)
	}
}
