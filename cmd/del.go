package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DelResult contains the result of a Del operation
type DelResult struct {
	DeletedTodos []todotxtlib.Todo
}

// Del deletes todos from the repository by indices (0-based)
func Del(repo todotxtlib.TodoRepository, indices []int) (DelResult, error) {
	if len(indices) == 0 {
		return DelResult{}, fmt.Errorf("no indices provided")
	}

	// STEP 1: Validate all indices first (fail-fast)
	allTodos, err := repo.ListAll()
	if err != nil {
		return DelResult{}, fmt.Errorf("failed to list todos: %w", err)
	}

	for _, index := range indices {
		if index < 0 || index >= len(allTodos) {
			return DelResult{}, fmt.Errorf("invalid index: %d", index)
		}
	}

	// STEP 2: Sort indices in descending order to delete from end to beginning
	// This prevents index shifting issues
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	for i := 0; i < len(sortedIndices); i++ {
		for j := i + 1; j < len(sortedIndices); j++ {
			if sortedIndices[i] < sortedIndices[j] {
				sortedIndices[i], sortedIndices[j] = sortedIndices[j], sortedIndices[i]
			}
		}
	}

	// STEP 3: Delete in descending order
	deletedTodos := make([]todotxtlib.Todo, 0, len(sortedIndices))
	for _, index := range sortedIndices {
		deleted, err := repo.Remove(index)
		if err != nil {
			return DelResult{}, fmt.Errorf("failed to remove todo at index %d: %w", index, err)
		}
		deletedTodos = append(deletedTodos, deleted)
	}

	// STEP 4: Save
	if err := repo.Save(); err != nil {
		return DelResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelResult{
		DeletedTodos: deletedTodos,
	}, nil
}
