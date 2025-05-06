package cmd

import (
	"fmt"
	"os"

	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/gkarolyi/togodo/todotxtui"
)

// BaseCommand provides common functionality for all commands.
// It contains a repository, a formatter, and an output writer.
type BaseCommand struct {
	Repository *todotxtlib.Repository
	Formatter  todotxtui.TodoFormatter
	Output     todotxtui.OutputWriter
}

// Print prints a single todo item to the output.
// It formats the todo item using the formatter and writes the formatted todo item to the output.
func (c *BaseCommand) Print(todo todotxtlib.Todo) error {
	formatted := c.Formatter.Format(todo)
	c.Output.WriteLine(formatted)
	return nil
}

// PrintList prints a list of todo items to the output.
// It formats the list using the formatter and writes the formatted list to the output.
func (c *BaseCommand) PrintList(todos []todotxtlib.Todo) error {
	formatted := c.Formatter.FormatList(todos)
	c.Output.WriteLines(formatted)
	return nil
}

// // PrintUpdated prints the updated todo items with their line numbers to the output.
// func (c *BaseCommand) PrintUpdated(updatedTodoTexts []string) error {
// 	updatedTodos := make([]todotxtlib.Todo, len(updatedTodoTexts))
// 	allTodos, err := c.Repository.ListAll()
// 	if err != nil {
// 		return err
// 	}
// 	for _, todoText := range updatedTodoTexts {
// 		for idx, todo := range allTodos {
// 			if todo.Text == todoText {
// 				c.Output.WriteLine()
// 			}
// 		}
// 	}
// 	c.Output.WriteLines(updated)
// 	return nil
// }

// Write saves the repository using its assigned writer
func (c *BaseCommand) Write() error {
	return c.Repository.Save()
}

// newBaseCommand creates a new base command with the given dependencies
func newBaseCommand(repo *todotxtlib.Repository, formatter todotxtui.TodoFormatter, output todotxtui.OutputWriter) *BaseCommand {
	return &BaseCommand{
		Repository: repo,
		Formatter:  formatter,
		Output:     output,
	}
}

// NewDefaultBaseCommand creates a new base command with the default dependencies.
// It creates a new file repository with the default todo.txt path, the default lipgloss
// formatter, and writes to stdout.
func NewDefaultBaseCommand() *BaseCommand {
	todoTxtPath := getTodoTxtPath()
	repo, err := newFileRepository(todoTxtPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return newBaseCommand(repo, newLipglossFormatter(), newStdoutWriter())
}

func NewTUIBaseCommand() *BaseCommand {
	todoTxtPath := getTodoTxtPath()
	repo, err := newFileRepository(todoTxtPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return newBaseCommand(repo, newLipglossFormatter(), todotxtui.NewTUIWriter(repo))
}

// getTodoTxtPath returns the path to the todo.txt file.
// It tries to open a todo.txt file in the directory specified by the TODO_TXT_PATH
// environment variable, then in the current directory, and finally in the user's home directory.
func getTodoTxtPath() string {
	var todoTxtPath string
	// First try to open a todo.txt file in the directory specified by the TODO_TXT_PATH environment variable
	// if todoTxtPath == "" {
	// 	envtodoTxtPath := os.Getenv("TODO_TXT_PATH")
	// 	if _, err := os.Stat(envtodoTxtPath); os.IsNotExist(err) {
	// 		os.Exit(1)
	// 	} else {
	// 		todoTxtPath = envtodoTxtPath
	// 	}
	// }

	// If that fails, try to open a todo.txt file in the current directory
	if todoTxtPath == "" {
		if _, err := os.Stat("todo.txt"); os.IsNotExist(err) {
			os.Exit(1)
		} else {
			todoTxtPath = "todo.txt"
		}
	}

	// Finally, try to open a todo.txt file in the user's home directory
	if todoTxtPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}
		todoTxtPath = homeDir + "/todo.txt"
		if _, err := os.Stat(todoTxtPath); os.IsNotExist(err) {
			os.Exit(1)
		}
	}

	return todoTxtPath
}

func newFileRepository(path string) (*todotxtlib.Repository, error) {
	reader := todotxtlib.NewFileReader(path)
	writer := todotxtlib.NewFileWriter(path)

	return todotxtlib.NewRepository(reader, writer)
}

func newLipglossFormatter() todotxtui.TodoFormatter {
	return todotxtui.NewLipglossFormatter()
}

func newStdoutWriter() todotxtui.OutputWriter {
	return todotxtui.NewStdoutWriter()
}
