package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// PrependResult contains the result of a Prepend operation
type PrependResult struct {
	Todo       todotxtlib.Todo
	LineNumber int
}

// Prepend prepends text to a todo, preserving priority if present
func Prepend(repo todotxtlib.TodoRepository, index int, text string) (PrependResult, error) {
	// Get existing todo
	todo, err := repo.Get(index)
	if err != nil {
		return PrependResult{}, fmt.Errorf("failed to get todo: %w", err)
	}

	// Extract priority if present
	todoText := todo.Text
	var newText string

	if len(todoText) >= 4 && todoText[0] == '(' && todoText[2] == ')' && todoText[3] == ' ' {
		// Has priority like "(A) task"
		priority := todoText[0:4] // "(A) "
		rest := todoText[4:]
		newText = priority + text + " " + rest
	} else {
		// No priority
		newText = text + " " + todoText
	}

	// Create updated todo
	updatedTodo := todotxtlib.NewTodo(newText)
	// Preserve the line number from the old todo
	updatedTodo.LineNumber = todo.LineNumber

	// Update in repository
	updated, err := repo.Update(index, updatedTodo)
	if err != nil {
		return PrependResult{}, fmt.Errorf("failed to update todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return PrependResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return PrependResult{
		Todo:       updated,
		LineNumber: updated.LineNumber,
	}, nil
}
