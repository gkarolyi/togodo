package todotxtlib

import "os"

// Writer handles writing Todo structs to a todo.txt file
type Writer interface {
	Write(path string, todos []Todo) error
}

// NewFileWriter returns a new FileWriter
func NewFileWriter() Writer {
	return &fileWriter{}
}

// fileWriter is a Writer that writes Todo structs to a todo.txt file
type fileWriter struct {
}

// Write writes the given todos to the file at the given path
func (w *fileWriter) Write(path string, todos []Todo) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, todo := range todos {
		if _, err := file.WriteString(todo.Text + "\n"); err != nil {
			return err
		}
	}

	return nil
}
