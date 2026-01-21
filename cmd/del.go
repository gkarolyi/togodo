package cmd

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// DelResult contains the result of a Del operation
type DelResult struct {
	DeletedTodo todotxtlib.Todo
	LineNumber  int
}

// Del deletes a todo from the repository
func Del(repo todotxtlib.TodoRepository, index int) (DelResult, error) {
	// Remove todo
	deleted, err := repo.Remove(index)
	if err != nil {
		return DelResult{}, fmt.Errorf("failed to remove todo: %w", err)
	}

	// Save
	if err := repo.Save(); err != nil {
		return DelResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return DelResult{
		DeletedTodo: deleted,
		LineNumber:  index + 1,
	}, nil
}
