package todotxtlib

import (
	"os"
	"strings"
)

// Reader handles reading todo.txt files and parsing them into a Todo struct
type Reader interface {
	Read(path string) ([]Todo, error)
}

// fileReader is a Reader that reads todo.txt files
type fileReader struct{}

// NewFileReader returns a new FileReader
func NewFileReader() Reader {
	return &fileReader{}
}

// Read reads the content of a todo.txt file and returns a slice of Todo structs
func (r *fileReader) Read(path string) (todos []Todo, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		todo := NewTodo(line)
		todos = append(todos, todo)
	}

	return todos, nil
}
