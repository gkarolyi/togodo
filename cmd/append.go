package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// AppendResult contains the result of an Append operation
type AppendResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Append appends text to the end of a todo
func Append(repo todotxtlib.TodoRepository, index int, text string) (AppendResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return AppendResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Append text
	newText := todo.Text + " " + text

	// Create updated todo
	updatedTodo := todotxtlib.NewTodo(newText)

	// Update in repository
	updated, err := repo.Update(index, updatedTodo)
	if err != nil {
		return AppendResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return AppendResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return AppendResult{
		Todo:       updated,
		LineNumber: index + 1,
	}, nil
}
