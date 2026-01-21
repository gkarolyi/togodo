package cmd

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// ListprojResult contains the result of listing projects
type ListprojResult struct {
	Projects []string
}

// Listproj lists all projects found in todos
func Listproj(repo todotxtlib.TodoRepository) (ListprojResult, error) {
	projects, err := repo.ListProjects()
	if err != nil {
		return ListprojResult{}, err
	}

	return ListprojResult{
		Projects: projects,
	}, nil
}
