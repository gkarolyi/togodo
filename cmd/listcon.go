package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListconResult contains the result of listing contexts
type ListconResult struct {
	Contexts []string
}

// Listcon lists all contexts found in todos
func Listcon(repo todotxtlib.TodoRepository) (ListconResult, error) {
	contexts, err := repo.ListContexts()
	if err != nil {
		return ListconResult{}, err
	}

	return ListconResult{
		Contexts: contexts,
	}, nil
}
