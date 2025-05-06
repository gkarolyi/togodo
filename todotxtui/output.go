package todotxtui

import (
	"fmt"
	"os"
)

// OutputWriter defines the interface for writing output
type OutputWriter interface {
	WriteLine(line string)
	WriteLines(lines []string)
	WriteError(err error)
	Run() error
}

// StdoutWriter implements OutputWriter for standard output
type StdoutWriter struct{}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func (w *StdoutWriter) WriteLine(line string) {
	fmt.Println(line)
}

func (w *StdoutWriter) WriteLines(lines []string) {
	for _, line := range lines {
		w.WriteLine(line)
	}
}

func (w *StdoutWriter) WriteError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func (w *StdoutWriter) Run() error {
	return nil
}
