package cmd

import (
	tui "github.com/gkarolyi/togodo/todotxt-tui"
	"github.com/gkarolyi/togodo/todotxtlib"
	"github.com/gkarolyi/togodo/todotxtui"
)

// Presenter handles formatting and outputting todos
type Presenter struct {
	formatter todotxtui.TodoFormatter
	output    todotxtui.OutputWriter
}

func NewPresenter(formatter todotxtui.TodoFormatter, output todotxtui.OutputWriter) *Presenter {
	return &Presenter{
		formatter: formatter,
		output:    output,
	}
}

// Print prints a single todo item
func (p *Presenter) Print(todo todotxtlib.Todo) error {
	formatted := p.formatter.Format(todo)
	p.output.WriteLine(formatted)
	return nil
}

// PrintList prints a list of todo items
func (p *Presenter) PrintList(todos []todotxtlib.Todo) error {
	formatted := p.formatter.FormatList(todos)
	p.output.WriteLines(formatted)
	return nil
}

// Factory functions for creating dependencies
func createRepository() (*todotxtlib.Repository, error) {
	todoTxtPath := GetTodoTxtPath()
	reader := todotxtlib.NewFileReader(todoTxtPath)
	writer := todotxtlib.NewFileWriter(todoTxtPath)
	return todotxtlib.NewRepository(reader, writer)
}

func createCLIPresenter() *Presenter {
	formatter := todotxtui.NewLipglossFormatter()
	output := todotxtui.NewStdoutWriter()
	return NewPresenter(formatter, output)
}

func createTUIController(repo *todotxtlib.Repository) interface{ Run() error } {
	return tui.NewController(repo)
}
