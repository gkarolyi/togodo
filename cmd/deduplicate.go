package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DeduplicateResult contains the result of a Deduplicate operation
type DeduplicateResult struct {
	RemovedCount int
}

// Deduplicate removes duplicate tasks, keeping only the first occurrence of each unique task text
func Deduplicate(repo todotxtlib.TodoRepository) (DeduplicateResult, error) {
	// Get all todos
	allTodos, err := repo.ListAll()
	if err != nil {
		return DeduplicateResult{}, fmt.Errorf("failed to list all todos: %w", err)
	}

	// Track seen task texts (case-sensitive)
	seen := make(map[string]bool)
	removedCount := 0

	// Process todos backwards to avoid index shifting when removing
	for i := len(allTodos) - 1; i >= 0; i-- {
		text := allTodos[i].Text
		if seen[text] {
			// This is a duplicate, remove it
			if _, err := repo.Remove(i); err != nil {
				return DeduplicateResult{}, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
			}
			removedCount++
		} else {
			// First occurrence, mark as seen
			seen[text] = true
		}
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return DeduplicateResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DeduplicateResult{RemovedCount: removedCount}, nil
}
