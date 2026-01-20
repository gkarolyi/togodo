package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// PriResult contains the result of a SetPriority operation
type PriResult struct {
	UpdatedTodos []todotxtlib.Todo
}

// SetPriority sets the priority for todos at the given indices (0-based)
func SetPriority(repo todotxtlib.TodoRepository, indices []int, priority string) (PriResult, error) {
	// STEP 1: Validate all indices first (fail-fast)
	allTodos, err := repo.ListAll()
	if err != nil {
		return PriResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	for _, index := range indices {
		if index < 0 || index >= len(allTodos) {
			return PriResult{}, fmt.Errorf("invalid index: %d", index)
		}
	}

	// STEP 2: All indices are valid, proceed with setting priorities
	updatedTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.SetPriority(index, priority)
		if err != nil {
			// This should not happen since we validated indices
			return PriResult{}, fmt.Errorf("failed to set priority at index %d: %w", index, err)
		}
		updatedTodos = append(updatedTodos, todo)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := repo.Save(); err != nil {
		return PriResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return PriResult{UpdatedTodos: updatedTodos}, nil
}
