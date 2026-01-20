package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkarolyi/togodo/business"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewPriCmd creates a Cobra command for setting priority
func NewPriCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "pri [LINE_NUMBER] [PRIORITY]",
		Short: "Set priority of a todo item",
		Long: `Sets the priority of a todo item (A, B, C, etc.).

# set task 1 to priority A
togodo pri 1 A
`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"p"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number (1-based)
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Convert to 0-based index
			index := lineNum - 1

			// Get priority (normalize to uppercase)
			priority := strings.ToUpper(args[1])

			// Call business logic
			result, err := business.SetPriority(repo, []int{index}, priority)
			if err != nil {
				return err
			}

			// Format output to match todo.txt-cli
			for _, todo := range result.UpdatedTodos {
				fmt.Printf("%d %s\n", lineNum, todo.Text)
				fmt.Printf("TODO: %d prioritized to (%s).\n", lineNum, priority)
			}
			return nil
		},
	}
}
