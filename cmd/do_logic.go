package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DoResult contains the result of a Do operation
type DoResult struct {
	ToggledTodos []todotxtlib.Todo
}

// Do toggles the done status of todos at the given indices (0-based)
func Do(repo todotxtlib.TodoRepository, indices []int) (DoResult, error) {
	toggledTodos := make([]todotxtlib.Todo, 0, len(indices))

	for _, index := range indices {
		todo, err := repo.ToggleDone(index)
		if err != nil {
			return DoResult{}, fmt.Errorf("failed to toggle todo at index %d: %w", index, err)
		}
		toggledTodos = append(toggledTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return DoResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DoResult{ToggledTodos: toggledTodos}, nil
}
