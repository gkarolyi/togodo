package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// ReplaceResult contains the result of a Replace operation
type ReplaceResult struct {
	OldTodo    todotxtlib.Todo
	NewTodo    todotxtlib.Todo
	LineNumber int
}

// Replace replaces the entire text of a todo at the given index
func Replace(repo todotxtlib.TodoRepository, index int, newText string) (ReplaceResult, error) {
	// Get existing todo
	oldTodo, err := repo.Get(index)
	if err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Create new todo with replaced text
	newTodo := todotxtlib.NewTodo(newText)
	// Preserve the line number from the old todo
	newTodo.LineNumber = oldTodo.LineNumber

	// Update in repository
	updated, err := repo.Update(index, newTodo)
	if err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return ReplaceResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return ReplaceResult{
		OldTodo:    oldTodo,
		NewTodo:    updated,
		LineNumber: index + 1, // Convert to 1-based
	}, nil
}
