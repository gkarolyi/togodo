package business

import (
	"fmt"

	"github.com/gkarolyi/togodo/todotxtlib"
)

// TidyResult contains the result of a Tidy operation
type TidyResult struct {
	RemovedTodos []todotxtlib.Todo
}

// Tidy removes all completed todos
func Tidy(repo todotxtlib.TodoRepository) (TidyResult, error) {
	// Get done todos before removing
	doneFilter := todotxtlib.Filter{Done: "true"}
	doneTodos, err := repo.Filter(doneFilter)
	if err != nil {
		return TidyResult{}, fmt.Errorf("failed to filter done todos: %w", err)
	}

	// Get all todos to iterate
	allTodos, err := repo.ListAll()
	if err != nil {
		return TidyResult{}, fmt.Errorf("failed to list all todos: %w", err)
	}

	// Remove backwards to avoid index shifting
	for i := len(allTodos) - 1; i >= 0; i-- {
		if allTodos[i].Done {
			if _, err := repo.Remove(i); err != nil {
				return TidyResult{}, fmt.Errorf("failed to remove todo at index %d: %w", i, err)
			}
		}
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return TidyResult{}, fmt.Errorf("failed to save: %w", err)
	}

	return TidyResult{RemovedTodos: doneTodos}, nil
}
