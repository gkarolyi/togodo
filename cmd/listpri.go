package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListpriResult contains the result of a Listpri operation
type ListpriResult struct {
	Todos       []todotxtlib.Todo
	LineNumbers []int
	Priority    string
	TotalCount  int
}

// Listpri lists todos with a specific priority
func Listpri(repo todotxtlib.TodoRepository, priority string) (ListpriResult, error) {
	// Get all todos
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListpriResult{}, err
	}

	// Filter by priority
	var filtered []todotxtlib.Todo
	var lineNumbers []int
	for i, todo := range allTodos {
		if !todo.Done && todo.Priority != "" {
			// If no priority specified, include all prioritized todos
			// Otherwise, only include matching priority
			if priority == "" || todo.Priority == priority {
				filtered = append(filtered, todo)
				lineNumbers = append(lineNumbers, i+1)
			}
		}
	}

	return ListpriResult{
		Todos:       filtered,
		LineNumbers: lineNumbers,
		Priority:    priority,
		TotalCount:  len(allTodos),
	}, nil
}
