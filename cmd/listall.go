package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListallResult contains the result of listing all todos
type ListallResult struct {
	Todos      []todotxtlib.Todo
	TotalCount int
}

// Listall lists all todos including completed ones
func Listall(repo todotxtlib.TodoRepository) (ListallResult, error) {
	// Get all todos (including done)
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListallResult{}, err
	}

	return ListallResult{
		Todos:      allTodos,
		TotalCount: len(allTodos),
	}, nil
}
