package cli

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gkarolyi/togodo/cmd"
	"github.com/gkarolyi/togodo/internal/config"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/spf13/cobra"
)

// NewMoveCmd creates a Cobra command for moving todos
func NewMoveCmd(repo todotxtlib.TodoRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "move ITEM# DEST_FILE",
		Short: "Move a todo item to another file",
		Long: `Moves a todo item from current file to another file.

# move task 1 to done.txt
togodo move 1 done.txt

# move task 3 to archive.txt
togodo move 3 archive.txt
`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"mv"},
		RunE: func(command *cobra.Command, args []string) error {
			// Parse line number
			lineNum, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid line number: %s", args[0])
			}

			// Get destination file path
			destFile := args[1]
			if !filepath.IsAbs(destFile) {
				// If relative path, resolve relative to todo.txt directory
				todoDir := filepath.Dir(config.GetTodoTxtPath())
				destFile = filepath.Join(todoDir, destFile)
			}

			// Create destination repository
			destReader := todotxtlib.NewFileReader(destFile)
			destWriter := todotxtlib.NewFileWriter(destFile)
			destRepo, err := todotxtlib.NewFileRepository(destReader, destWriter)
			if err != nil {
				return fmt.Errorf("failed to create destination repository: %w", err)
			}

			// Get source file name for output
			sourceFile := filepath.Base(config.GetTodoTxtPath())
			destFileName := filepath.Base(destFile)

			// Call business logic
			result, err := cmd.Move(repo, destRepo, lineNum)
			if err != nil {
				return err
			}

			// Format output
			fmt.Fprintf(command.OutOrStdout(), "%d %s\n", result.LineNumber, result.Todo.Text)
			fmt.Fprintf(command.OutOrStdout(), "TODO: %d moved from '%s' to '%s'.\n",
				result.LineNumber, sourceFile, destFileName)

			return nil
		},
	}
}
