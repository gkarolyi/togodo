package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListResult contains the result of a List operation
type ListResult struct {
	Todos       []todotxtlib.Todo
	LineNumbers []int
	TotalCount  int
	ShownCount  int
}

// List lists todos, optionally filtering by search query
func List(repo todotxtlib.TodoRepository, searchQuery string) (ListResult, error) {
	// Get total count
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListResult{}, err
	}
	totalCount := len(allTodos)

	// Filter if search query provided
	var todos []todotxtlib.Todo
	var lineNumbers []int
	if searchQuery != "" {
		filter := todotxtlib.Filter{Text: searchQuery}
		filteredTodos, err := repo.Filter(filter)
		if err != nil {
			return ListResult{}, err
		}
		todos = filteredTodos

		// Find line numbers in original list
		for _, filteredTodo := range filteredTodos {
			for i, originalTodo := range allTodos {
				if originalTodo.Text == filteredTodo.Text {
					lineNumbers = append(lineNumbers, i+1)
					break
				}
			}
		}
	} else {
		todos = allTodos
		// Generate sequential line numbers for unfiltered list
		for i := range allTodos {
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return ListResult{
		Todos:       todos,
		LineNumbers: lineNumbers,
		TotalCount:  totalCount,
		ShownCount:  len(todos),
	}, nil
}
