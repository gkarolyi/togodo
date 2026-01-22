package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// PriResult contains the result of a SetPriority operation
type PriResult struct {
	UpdatedTodos  []todotxtlib.Todo
	OldPriorities []string // Old priority for each updated todo (empty string if no priority)
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
	oldPriorities := make([]string, 0, len(indices))

	for _, index := range indices {
		// Get the old priority before changing
		oldTodo := allTodos[index]
		oldPriority := oldTodo.Priority

		// Check if trying to set same priority
		if oldPriority != "" && oldPriority == priority {
			return PriResult{}, fmt.Errorf("TODO: %d already prioritized (%s).", index+1, priority)
		}

		todo, err := repo.SetPriority(index, priority)
		if err != nil {
			// This should not happen since we validated indices
			return PriResult{}, fmt.Errorf("failed to set priority at index %d: %w", index, err)
		}
		updatedTodos = append(updatedTodos, todo)
		oldPriorities = append(oldPriorities, oldPriority)
	}

	// Note: Pri command doesn't sort - preserves user's order
	if err := repo.Save(); err != nil {
		return PriResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return PriResult{
		UpdatedTodos:  updatedTodos,
		OldPriorities: oldPriorities,
	}, nil
}
