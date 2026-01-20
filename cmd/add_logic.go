package cmd

import "github.com/gkarolyi/togodo/todotxtlib"

// AddResult represents the result of an add operation
type AddResult struct {
	Added []todotxtlib.Todo
	Error error
}

// Add adds multiple todos to the repository, sorts, and saves
func Add(repo todotxtlib.TodoRepository, texts []string) AddResult {
	addedTodos := make([]todotxtlib.Todo, 0, len(texts))

	for _, text := range texts {
		todo, err := repo.Add(text)
		if err != nil {
			return AddResult{Error: err}
		}
		addedTodos = append(addedTodos, todo)
	}

	repo.Sort(nil)
	if err := repo.Save(); err != nil {
		return AddResult{Error: err}
	}

	return AddResult{Added: addedTodos}
}
