package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DepriResult contains the result of a Depri operation
type DepriResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Depri removes priority from a todo
func Depri(repo todotxtlib.TodoRepository, index int) (DepriResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return DepriResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Check if task has a priority to remove
	if todo.Priority == "" {
		return DepriResult{}, fmt.Errorf("TODO: %d is not prioritized.", todo.LineNumber)
	}

	// Remove priority by setting to empty string
	updated, err := repo.SetPriority(index, "")
	if err != nil {
		return DepriResult{}, fmt.Errorf("failed to remove priority: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return DepriResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DepriResult{
		Todo:       updated,
		LineNumber: updated.LineNumber,
	}, nil
}
