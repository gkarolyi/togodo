package business

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListResult contains the result of a List operation
type ListResult struct {
	Todos      []todotxtlib.Todo
	TotalCount int
	ShownCount int
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
	if searchQuery != "" {
		filter := todotxtlib.Filter{Text: searchQuery}
		todos, err = repo.Filter(filter)
		if err != nil {
			return ListResult{}, err
		}
	} else {
		todos = allTodos
	}

	return ListResult{
		Todos:      todos,
		TotalCount: totalCount,
		ShownCount: len(todos),
	}, nil
}
