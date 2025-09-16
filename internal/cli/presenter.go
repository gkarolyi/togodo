package cli

import (
	"github.com/gkarolyi/togodo/todotxtlib"
)

// Presenter handles formatting and outputting todos
type Presenter struct {
	formatter TodoFormatter
	output    OutputWriter
}

func NewPresenter() *Presenter {
	return &Presenter{
		formatter: NewLipglossFormatter(),
		output:    NewStdoutWriter(),
	}
}

// Print prints a single todo item
func (p *Presenter) Print(todo todotxtlib.Todo) error {
	formatted := p.formatter.Format(todo)
	p.WriteLine(formatted)
	return nil
}

// PrintList prints a list of todo items
func (p *Presenter) PrintList(todos []todotxtlib.Todo) error {
	formatted := p.formatter.FormatList(todos)
	p.output.WriteLines(formatted)
	return nil
}

// WriteLine writes a single line to the output
func (p *Presenter) WriteLine(line string) error {
	p.output.WriteLine(line)
	return nil
}
