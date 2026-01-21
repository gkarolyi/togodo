package todotxtlib

import (
	"io"
	"os"
	"strings"
)

// Reader handles reading todo.txt content and parsing it into Todo structs
type Reader interface {
	Read() ([]Todo, error)
}

// NewFileReader returns a new Reader that reads from a file
func NewFileReader(path string) Reader {
	return &fileReader{
		path: path,
	}
}

// fileReader is a Reader that reads from a file
type fileReader struct {
	path string
}

// Read reads the content of a todo.txt file and returns a slice of Todo structs
func (r *fileReader) Read() (todos []Todo, err error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	file, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readFromReader(file)
}

// NewBufferReader returns a new Reader that reads from an io.Reader
func NewBufferReader(r io.Reader) Reader {
	return &bufferReader{
		reader: r,
	}
}

// bufferReader is a Reader that reads from an io.Reader
type bufferReader struct {
	reader io.Reader
}

// Read reads the content from the io.Reader and returns a slice of Todo structs
func (r *bufferReader) Read() ([]Todo, error) {
	return readFromReader(r.reader)
}

// readFromReader reads todos from any io.Reader
func readFromReader(r io.Reader) ([]Todo, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	todos := make([]Todo, 0, len(lines))

	lineNumber := 1
	for _, line := range lines {
		if line == "" {
			continue
		}
		todo := NewTodo(line)
		todo.LineNumber = lineNumber
		todos = append(todos, todo)
		lineNumber++
	}

	return todos, nil
}
