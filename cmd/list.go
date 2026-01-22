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
	// Sort todos using default sort (priority + alphabetical)
	repo.Sort(nil)

	// Get total count
	allTodos, err := repo.ListAll()
	if err != nil {
		return ListResult{}, err
	}
	totalCount := len(allTodos)

	// Filter if search query provided
	var todos []todotxtlib.Todo
	if searchQuery != "" {
		filter := todotxtlib.Filter{Text: searchQuery}
		filteredTodos, err := repo.Filter(filter)
		if err != nil {
			return ListResult{}, err
		}
		todos = filteredTodos
	} else {
		todos = allTodos
	}

	// Extract line numbers from todos
	lineNumbers := make([]int, len(todos))
	for i, todo := range todos {
		lineNumbers[i] = todo.LineNumber
	}

	return ListResult{
		Todos:       todos,
		LineNumbers: lineNumbers,
		TotalCount:  totalCount,
		ShownCount:  len(todos),
	}, nil
}
