package cmd

import (
	"strings"

	"github.com/gkarolyi/togodo/internal/cli"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

func executeList(repo todotxtlib.TodoRepository, presenter *cli.Presenter, searchQuery string) error {
	todos, err := repo.Search(searchQuery)
	if err != nil {
		return err
	}

	return presenter.PrintList(todos)
}

// NewListCmd creates a new cobra command for listing todos.
func NewListCmd(repo todotxtlib.TodoRepository, presenter *cli.Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "list [FILTER]",
		Short: "List and filter items in your todo.txt",
		Long: `Lists tasks sorted in order of priority, with done items at the bottom of the list. Tasks can optionally be filtered
by passing an optional [FILTER] argument. If no filter is passed, list shows all items in your todo.txt file. Tasks are shown
with a line number to allow you to easily refer to them. For example:

# list all items in your todo.txt file
togodo list

# list all items in your todo.txt file that contain the string '@work'
togodo list '@work'
`,
		Aliases: []string{"ls", "l"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			searchQuery := strings.Join(args, " ")
			return executeList(repo, presenter, searchQuery)
		},
	}
}
